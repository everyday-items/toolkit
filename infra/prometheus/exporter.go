package prometheus

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/everyday-items/toolkit/infra/observe"
)

// Exporter Prometheus 导出器
type Exporter struct {
	// namespace 命名空间
	namespace string

	// subsystem 子系统
	subsystem string

	// registry 指标注册表
	registry *Registry

	// collector 指标收集器
	collector *Collector

	// server HTTP 服务器
	server *http.Server

	mu sync.RWMutex
}

// ExporterOption 导出器选项
type ExporterOption func(*Exporter)

// NewExporter 创建 Prometheus 导出器
func NewExporter(opts ...ExporterOption) *Exporter {
	e := &Exporter{
		namespace: "app",
		registry:  NewRegistry(),
	}

	for _, opt := range opts {
		opt(e)
	}

	e.collector = NewCollector(e.registry, e.namespace, e.subsystem)

	return e
}

// WithNamespace 设置命名空间
func WithNamespace(namespace string) ExporterOption {
	return func(e *Exporter) {
		e.namespace = namespace
	}
}

// WithSubsystem 设置子系统
func WithSubsystem(subsystem string) ExporterOption {
	return func(e *Exporter) {
		e.subsystem = subsystem
	}
}

// Registry 返回注册表
func (e *Exporter) Registry() *Registry {
	return e.registry
}

// Collector 返回收集器
func (e *Exporter) Collector() *Collector {
	return e.collector
}

// Handler 返回 HTTP 处理器
func (e *Exporter) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		w.Write([]byte(e.registry.Gather()))
	})
}

// ListenAndServe 启动 HTTP 服务器
func (e *Exporter) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", e.Handler())

	e.mu.Lock()
	e.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	e.mu.Unlock()

	return e.server.ListenAndServe()
}

// Shutdown 关闭服务器
func (e *Exporter) Shutdown() error {
	e.mu.RLock()
	server := e.server
	e.mu.RUnlock()

	if server != nil {
		return server.Close()
	}
	return nil
}

// ============== Metrics Interface Adapter ==============

// MetricsAdapter 将 observe.Metrics 适配到 Prometheus
type MetricsAdapter struct {
	registry  *Registry
	namespace string
	subsystem string
}

// NewMetricsAdapter 创建适配器
func NewMetricsAdapter(registry *Registry, namespace, subsystem string) *MetricsAdapter {
	return &MetricsAdapter{
		registry:  registry,
		namespace: namespace,
		subsystem: subsystem,
	}
}

// Counter 获取 Counter
func (a *MetricsAdapter) Counter(name string, tags ...string) observe.Counter {
	fullName := a.fullName(name)
	return &counterAdapter{
		counter: a.registry.Counter(fullName, name),
	}
}

// Histogram 获取 Histogram
func (a *MetricsAdapter) Histogram(name string, tags ...string) observe.Histogram {
	fullName := a.fullName(name)
	return &histogramAdapter{
		histogram: a.registry.Histogram(fullName, name, DefaultBuckets),
	}
}

// Gauge 获取 Gauge
func (a *MetricsAdapter) Gauge(name string, tags ...string) observe.Gauge {
	fullName := a.fullName(name)
	return &gaugeAdapter{
		gauge: a.registry.Gauge(fullName, name),
	}
}

// Timer 获取 Timer
func (a *MetricsAdapter) Timer(name string, tags ...string) observe.Timer {
	fullName := a.fullName(name)
	return &timerAdapter{
		histogram: a.registry.Histogram(fullName+"_seconds", name, DefaultBuckets),
	}
}

func (a *MetricsAdapter) fullName(name string) string {
	parts := []string{}
	if a.namespace != "" {
		parts = append(parts, a.namespace)
	}
	if a.subsystem != "" {
		parts = append(parts, a.subsystem)
	}
	parts = append(parts, name)
	return strings.Join(parts, "_")
}

// 确保实现了 observe.Metrics 接口
var _ observe.Metrics = (*MetricsAdapter)(nil)

// 适配器实现

type counterAdapter struct {
	counter *PrometheusCounter
	value   float64
}

func (c *counterAdapter) Inc() {
	c.counter.Inc()
	c.value++
}

func (c *counterAdapter) Add(v float64) {
	c.counter.Add(v)
	c.value += v
}

func (c *counterAdapter) Value() float64 {
	return c.value
}

type histogramAdapter struct {
	histogram *PrometheusHistogram
	count     uint64
	sum       float64
}

func (h *histogramAdapter) Observe(v float64) {
	h.histogram.Observe(v)
	h.count++
	h.sum += v
}

func (h *histogramAdapter) Count() uint64 {
	return h.count
}

func (h *histogramAdapter) Sum() float64 {
	return h.sum
}

type gaugeAdapter struct {
	gauge *PrometheusGauge
	value float64
}

func (g *gaugeAdapter) Set(v float64) {
	g.gauge.Set(v)
	g.value = v
}

func (g *gaugeAdapter) Inc() {
	g.gauge.Inc()
	g.value++
}

func (g *gaugeAdapter) Dec() {
	g.gauge.Dec()
	g.value--
}

func (g *gaugeAdapter) Add(v float64) {
	g.gauge.Add(v)
	g.value += v
}

func (g *gaugeAdapter) Value() float64 {
	return g.value
}

type timerAdapter struct {
	histogram *PrometheusHistogram
}

func (t *timerAdapter) ObserveDuration(d time.Duration) {
	t.histogram.Observe(d.Seconds())
}

func (t *timerAdapter) Time(fn func()) {
	start := time.Now()
	fn()
	t.ObserveDuration(time.Since(start))
}

func (t *timerAdapter) NewTimer() *observe.TimerContext {
	return observe.NewTimerContext(t)
}
