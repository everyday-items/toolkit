// Package otel 提供 OpenTelemetry 集成
//
// 支持将追踪数据导出到 Jaeger、Zipkin、OTLP 等后端。
//
// 使用示例:
//
//	tracer := otel.NewOTelTracer(
//	    otel.WithServiceName("my-service"),
//	    otel.WithEndpoint("localhost:4317"),
//	)
//	defer tracer.Shutdown(context.Background())
//
//	ctx, span := tracer.StartSpan(ctx, "operation")
//	defer span.End()
package otel

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/everyday-items/toolkit/infra/observe"
	"github.com/everyday-items/toolkit/util/idgen"
)

// OTelTracer OpenTelemetry 追踪器
type OTelTracer struct {
	// serviceName 服务名称
	serviceName string

	// exporter 导出器
	exporter Exporter

	// sampler 采样器
	sampler Sampler

	// propagator 传播器
	propagator Propagator

	// spans 活跃的 Span
	spans sync.Map

	// config 配置
	config OTelConfig

	mu sync.RWMutex
}

// OTelConfig OpenTelemetry 配置
type OTelConfig struct {
	// ServiceName 服务名称
	ServiceName string

	// ServiceVersion 服务版本
	ServiceVersion string

	// Environment 环境
	Environment string

	// Endpoint 导出端点
	Endpoint string

	// Headers 请求头
	Headers map[string]string

	// SamplingRate 采样率（0-1）
	SamplingRate float64

	// BatchSize 批量大小
	BatchSize int

	// BatchTimeout 批量超时
	BatchTimeout time.Duration

	// MaxQueueSize 最大队列大小
	MaxQueueSize int

	// Insecure 是否使用不安全连接
	Insecure bool
}

// DefaultOTelConfig 返回默认配置
func DefaultOTelConfig() OTelConfig {
	return OTelConfig{
		ServiceName:    "default",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		Endpoint:       "localhost:4317",
		SamplingRate:   1.0,
		BatchSize:      512,
		BatchTimeout:   5 * time.Second,
		MaxQueueSize:   2048,
		Insecure:       true,
	}
}

// OTelOption 配置选项
type OTelOption func(*OTelConfig)

// WithServiceName 设置服务名称
func WithServiceName(name string) OTelOption {
	return func(c *OTelConfig) {
		c.ServiceName = name
	}
}

// WithServiceVersion 设置服务版本
func WithServiceVersion(version string) OTelOption {
	return func(c *OTelConfig) {
		c.ServiceVersion = version
	}
}

// WithEnvironment 设置环境
func WithEnvironment(env string) OTelOption {
	return func(c *OTelConfig) {
		c.Environment = env
	}
}

// WithEndpoint 设置端点
func WithEndpoint(endpoint string) OTelOption {
	return func(c *OTelConfig) {
		c.Endpoint = endpoint
	}
}

// WithSamplingRate 设置采样率
func WithSamplingRate(rate float64) OTelOption {
	return func(c *OTelConfig) {
		c.SamplingRate = rate
	}
}

// WithBatchConfig 设置批量配置
func WithBatchConfig(size int, timeout time.Duration) OTelOption {
	return func(c *OTelConfig) {
		c.BatchSize = size
		c.BatchTimeout = timeout
	}
}

// NewOTelTracer 创建 OpenTelemetry 追踪器
func NewOTelTracer(opts ...OTelOption) *OTelTracer {
	config := DefaultOTelConfig()
	for _, opt := range opts {
		opt(&config)
	}

	t := &OTelTracer{
		serviceName: config.ServiceName,
		sampler:     NewProbabilitySampler(config.SamplingRate),
		propagator:  NewW3CTraceContextPropagator(),
		config:      config,
	}

	return t
}

// SetExporter 设置导出器
func (t *OTelTracer) SetExporter(exporter Exporter) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.exporter = exporter
}

// StartSpan 开始新的 Span
func (t *OTelTracer) StartSpan(ctx context.Context, name string, opts ...observe.SpanOption) (context.Context, observe.Span) {
	// 应用选项
	cfg := &observe.SpanConfig{
		Attributes: make(map[string]any),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// 生成 ID
	traceID := t.ExtractTraceID(ctx)
	if traceID == "" {
		traceID = idgen.NanoID()
	}
	spanID := idgen.NanoID()

	// 获取父 Span
	var parentSpanID string
	if parentSpan := observe.SpanFromContext(ctx); parentSpan != nil {
		parentSpanID = parentSpan.SpanID()
	}

	// 采样决策
	shouldSample := t.sampler.ShouldSample(traceID, name)

	// 创建 Span
	span := &OTelSpan{
		tracer:       t,
		traceID:      traceID,
		spanID:       spanID,
		parentSpanID: parentSpanID,
		name:         name,
		kind:         cfg.Kind,
		startTime:    time.Now(),
		attributes:   make(map[string]any),
		events:       make([]SpanEvent, 0),
		status:       observe.StatusCodeUnset,
		recording:    shouldSample,
	}

	// 设置初始属性
	for k, v := range cfg.Attributes {
		span.attributes[k] = v
	}

	// 添加资源属性
	span.attributes["service.name"] = t.serviceName
	span.attributes["service.version"] = t.config.ServiceVersion
	span.attributes["deployment.environment"] = t.config.Environment

	// 存储 Span
	t.spans.Store(spanID, span)

	// 更新 context
	ctx = observe.ContextWithSpan(ctx, span)
	ctx = t.InjectTraceID(ctx, traceID)

	return ctx, span
}

// ExtractTraceID 提取 Trace ID
func (t *OTelTracer) ExtractTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey{}).(string); ok {
		return traceID
	}
	return ""
}

// InjectTraceID 注入 Trace ID
func (t *OTelTracer) InjectTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// Shutdown 关闭追踪器
func (t *OTelTracer) Shutdown(ctx context.Context) error {
	// 导出所有剩余 Span
	if t.exporter != nil {
		spans := make([]*SpanData, 0)
		t.spans.Range(func(key, value any) bool {
			if span, ok := value.(*OTelSpan); ok {
				spans = append(spans, span.toSpanData())
			}
			return true
		})

		if len(spans) > 0 {
			t.exporter.ExportSpans(ctx, spans)
		}

		return t.exporter.Shutdown(ctx)
	}
	return nil
}

type traceIDKey struct{}

// OTelSpan OpenTelemetry Span 实现
type OTelSpan struct {
	tracer       *OTelTracer
	traceID      string
	spanID       string
	parentSpanID string
	name         string
	kind         observe.SpanKind
	startTime    time.Time
	endTime      time.Time
	attributes   map[string]any
	events       []SpanEvent
	status       observe.StatusCode
	statusMsg    string
	input        any
	output       any
	tokenUsage   observe.TokenUsage
	recording    bool
	ended        bool

	mu sync.RWMutex
}

// SpanEvent Span 事件
type SpanEvent struct {
	Name       string         `json:"name"`
	Timestamp  time.Time      `json:"timestamp"`
	Attributes map[string]any `json:"attributes"`
}

// SpanID 返回 Span ID
func (s *OTelSpan) SpanID() string {
	return s.spanID
}

// TraceID 返回 Trace ID
func (s *OTelSpan) TraceID() string {
	return s.traceID
}

// SetName 设置名称
func (s *OTelSpan) SetName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.name = name
}

// SetInput 设置输入
func (s *OTelSpan) SetInput(input any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.input = input
	s.attributes["input"] = input
}

// SetOutput 设置输出
func (s *OTelSpan) SetOutput(output any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.output = output
	s.attributes["output"] = output
}

// SetTokenUsage 设置 Token 使用量
func (s *OTelSpan) SetTokenUsage(usage observe.TokenUsage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokenUsage = usage
	s.attributes[observe.AttrLLMPromptTokens] = usage.PromptTokens
	s.attributes[observe.AttrLLMCompletionTokens] = usage.CompletionTokens
	s.attributes[observe.AttrLLMTotalTokens] = usage.TotalTokens
}

// SetAttribute 设置属性
func (s *OTelSpan) SetAttribute(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attributes[key] = value
}

// SetAttributes 批量设置属性
func (s *OTelSpan) SetAttributes(attrs map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range attrs {
		s.attributes[k] = v
	}
}

// AddEvent 添加事件
func (s *OTelSpan) AddEvent(name string, attrs ...any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event := SpanEvent{
		Name:       name,
		Timestamp:  time.Now(),
		Attributes: make(map[string]any),
	}

	// 解析属性
	for i := 0; i < len(attrs)-1; i += 2 {
		if key, ok := attrs[i].(string); ok {
			event.Attributes[key] = attrs[i+1]
		}
	}

	s.events = append(s.events, event)
}

// RecordError 记录错误
func (s *OTelSpan) RecordError(err error) {
	if err == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.attributes[observe.AttrErrorType] = fmt.Sprintf("%T", err)
	s.attributes[observe.AttrErrorMessage] = err.Error()

	s.events = append(s.events, SpanEvent{
		Name:      "exception",
		Timestamp: time.Now(),
		Attributes: map[string]any{
			"exception.type":    fmt.Sprintf("%T", err),
			"exception.message": err.Error(),
		},
	})
}

// SetStatus 设置状态
func (s *OTelSpan) SetStatus(code observe.StatusCode, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = code
	s.statusMsg = message
}

// End 结束 Span
func (s *OTelSpan) End() {
	s.mu.Lock()
	if s.ended {
		s.mu.Unlock()
		return
	}
	s.ended = true
	s.endTime = time.Now()
	s.mu.Unlock()

	// 从追踪器中移除
	s.tracer.spans.Delete(s.spanID)

	// 导出
	if s.tracer.exporter != nil && s.recording {
		s.tracer.exporter.ExportSpans(context.Background(), []*SpanData{s.toSpanData()})
	}
}

// EndWithError 结束并记录错误
func (s *OTelSpan) EndWithError(err error) {
	s.RecordError(err)
	s.SetStatus(observe.StatusCodeError, err.Error())
	s.End()
}

// IsRecording 是否正在记录
func (s *OTelSpan) IsRecording() bool {
	return s.recording
}

// toSpanData 转换为导出数据
func (s *OTelSpan) toSpanData() *SpanData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]SpanEvent, len(s.events))
	copy(events, s.events)

	attrs := make(map[string]any)
	for k, v := range s.attributes {
		attrs[k] = v
	}

	return &SpanData{
		TraceID:      s.traceID,
		SpanID:       s.spanID,
		ParentSpanID: s.parentSpanID,
		Name:         s.name,
		Kind:         s.kind,
		StartTime:    s.startTime,
		EndTime:      s.endTime,
		Attributes:   attrs,
		Events:       events,
		Status:       s.status,
		StatusMsg:    s.statusMsg,
	}
}

// SpanData 导出数据
type SpanData struct {
	TraceID      string             `json:"trace_id"`
	SpanID       string             `json:"span_id"`
	ParentSpanID string             `json:"parent_span_id,omitempty"`
	Name         string             `json:"name"`
	Kind         observe.SpanKind   `json:"kind"`
	StartTime    time.Time          `json:"start_time"`
	EndTime      time.Time          `json:"end_time"`
	Attributes   map[string]any     `json:"attributes"`
	Events       []SpanEvent        `json:"events"`
	Status       observe.StatusCode `json:"status"`
	StatusMsg    string             `json:"status_message,omitempty"`
}

// 确保实现了接口
var _ observe.Tracer = (*OTelTracer)(nil)
var _ observe.Span = (*OTelSpan)(nil)
