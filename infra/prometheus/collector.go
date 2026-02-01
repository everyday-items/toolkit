package prometheus

import (
	"runtime"
	"sync"
	"time"
)

// Collector 指标收集器
type Collector struct {
	registry  *Registry
	namespace string
	subsystem string

	// Go 运行时指标
	goGoroutines   *PrometheusGauge
	goMemAlloc     *PrometheusGauge
	goMemSys       *PrometheusGauge
	goGCPauseTotal *PrometheusCounter

	mu sync.RWMutex
}

// NewCollector 创建收集器
func NewCollector(registry *Registry, namespace, subsystem string) *Collector {
	c := &Collector{
		registry:  registry,
		namespace: namespace,
		subsystem: subsystem,
	}

	c.initRuntimeMetrics()
	c.startRuntimeCollector()

	return c
}

// initRuntimeMetrics 初始化运行时指标
func (c *Collector) initRuntimeMetrics() {
	prefix := c.namespace
	if c.subsystem != "" {
		prefix += "_" + c.subsystem
	}

	// Go 运行时指标
	c.goGoroutines = c.registry.Gauge(
		prefix+"_go_goroutines",
		"Number of goroutines",
	)
	c.goMemAlloc = c.registry.Gauge(
		prefix+"_go_memstats_alloc_bytes",
		"Number of bytes allocated and still in use",
	)
	c.goMemSys = c.registry.Gauge(
		prefix+"_go_memstats_sys_bytes",
		"Number of bytes obtained from system",
	)
	c.goGCPauseTotal = c.registry.Counter(
		prefix+"_go_gc_pause_seconds_total",
		"Total GC pause time in seconds",
	)
}

// startRuntimeCollector 启动运行时指标收集
func (c *Collector) startRuntimeCollector() {
	go func() {
		var lastGCPause uint64
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			c.goGoroutines.Set(float64(runtime.NumGoroutine()))
			c.goMemAlloc.Set(float64(m.Alloc))
			c.goMemSys.Set(float64(m.Sys))

			// GC 暂停时间
			gcPause := m.PauseTotalNs
			if gcPause > lastGCPause {
				c.goGCPauseTotal.Add(float64(gcPause-lastGCPause) / 1e9)
				lastGCPause = gcPause
			}
		}
	}()
}

// Counter 获取自定义 Counter
func (c *Collector) Counter(name, help string, labels ...string) *PrometheusCounter {
	return c.registry.Counter(c.fullName(name), help, labels...)
}

// Gauge 获取自定义 Gauge
func (c *Collector) Gauge(name, help string, labels ...string) *PrometheusGauge {
	return c.registry.Gauge(c.fullName(name), help, labels...)
}

// Histogram 获取自定义 Histogram
func (c *Collector) Histogram(name, help string, buckets []float64, labels ...string) *PrometheusHistogram {
	return c.registry.Histogram(c.fullName(name), help, buckets, labels...)
}

// Summary 获取自定义 Summary
func (c *Collector) Summary(name, help string, quantiles map[float64]float64, labels ...string) *PrometheusSummary {
	return c.registry.Summary(c.fullName(name), help, quantiles, labels...)
}

func (c *Collector) fullName(name string) string {
	if c.namespace == "" {
		return name
	}
	if c.subsystem == "" {
		return c.namespace + "_" + name
	}
	return c.namespace + "_" + c.subsystem + "_" + name
}

// RecordDuration 记录持续时间
func (c *Collector) RecordDuration(name string, duration time.Duration, labels ...string) {
	h := c.Histogram(name+"_seconds", "Duration in seconds", DefaultBuckets, labels...)
	h.Observe(duration.Seconds(), labels...)
}

// RecordCount 记录计数
func (c *Collector) RecordCount(name string, count float64, labels ...string) {
	counter := c.Counter(name+"_total", "Total count", labels...)
	counter.Add(count, labels...)
}

// SetGaugeValue 设置仪表盘值
func (c *Collector) SetGaugeValue(name string, value float64, labels ...string) {
	gauge := c.Gauge(name, "Gauge value", labels...)
	gauge.Set(value, labels...)
}
