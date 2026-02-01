package prometheus

import (
	"strings"
	"testing"
	"time"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestRegistryCounter(t *testing.T) {
	registry := NewRegistry()

	counter := registry.Counter("test_counter", "Test counter help")
	if counter == nil {
		t.Fatal("expected non-nil counter")
	}

	// 同名 counter 应该返回相同实例
	counter2 := registry.Counter("test_counter", "Test counter help")
	if counter != counter2 {
		t.Error("expected same counter instance")
	}
}

func TestPrometheusCounter(t *testing.T) {
	registry := NewRegistry()
	counter := registry.Counter("requests_total", "Total requests")

	counter.Inc()
	counter.Add(5)

	output := counter.String()
	if !strings.Contains(output, "requests_total") {
		t.Error("expected counter name in output")
	}
	if !strings.Contains(output, "# TYPE requests_total counter") {
		t.Error("expected counter type in output")
	}
}

func TestPrometheusCounterWithLabels(t *testing.T) {
	registry := NewRegistry()
	counter := registry.Counter("http_requests", "HTTP requests", "method", "path")

	counter.Inc("GET", "/api")
	counter.Inc("POST", "/api")
	counter.Add(3, "GET", "/api")

	output := counter.String()
	if !strings.Contains(output, `method="GET"`) {
		t.Error("expected label in output")
	}
}

func TestPrometheusGauge(t *testing.T) {
	registry := NewRegistry()
	gauge := registry.Gauge("active_connections", "Active connections")

	gauge.Set(10)
	gauge.Inc()
	gauge.Dec()
	gauge.Add(5)

	output := gauge.String()
	if !strings.Contains(output, "# TYPE active_connections gauge") {
		t.Error("expected gauge type in output")
	}
}

func TestPrometheusHistogram(t *testing.T) {
	registry := NewRegistry()
	histogram := registry.Histogram("request_duration", "Request duration", nil)

	histogram.Observe(0.1)
	histogram.Observe(0.5)
	histogram.Observe(1.0)

	output := histogram.String()
	if !strings.Contains(output, "# TYPE request_duration histogram") {
		t.Error("expected histogram type in output")
	}
	if !strings.Contains(output, "request_duration_bucket") {
		t.Error("expected bucket in output")
	}
	if !strings.Contains(output, "request_duration_sum") {
		t.Error("expected sum in output")
	}
	if !strings.Contains(output, "request_duration_count") {
		t.Error("expected count in output")
	}
}

func TestPrometheusSummary(t *testing.T) {
	registry := NewRegistry()
	summary := registry.Summary("response_time", "Response time", nil)

	for i := 0; i < 100; i++ {
		summary.Observe(float64(i) * 0.01)
	}

	output := summary.String()
	if !strings.Contains(output, "# TYPE response_time summary") {
		t.Error("expected summary type in output")
	}
}

func TestRegistryGather(t *testing.T) {
	registry := NewRegistry()

	registry.Counter("counter1", "Counter 1").Inc()
	registry.Gauge("gauge1", "Gauge 1").Set(42)

	output := registry.Gather()
	if !strings.Contains(output, "counter1") {
		t.Error("expected counter1 in output")
	}
	if !strings.Contains(output, "gauge1") {
		t.Error("expected gauge1 in output")
	}
}

func TestNewExporter(t *testing.T) {
	exporter := NewExporter()

	if exporter == nil {
		t.Fatal("expected non-nil exporter")
	}

	if exporter.namespace != "app" {
		t.Errorf("expected default namespace 'app', got '%s'", exporter.namespace)
	}
}

func TestExporterWithOptions(t *testing.T) {
	exporter := NewExporter(
		WithNamespace("myapp"),
		WithSubsystem("api"),
	)

	if exporter.namespace != "myapp" {
		t.Errorf("expected namespace 'myapp', got '%s'", exporter.namespace)
	}

	if exporter.subsystem != "api" {
		t.Errorf("expected subsystem 'api', got '%s'", exporter.subsystem)
	}
}

func TestExporterRegistry(t *testing.T) {
	exporter := NewExporter()

	registry := exporter.Registry()
	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestExporterCollector(t *testing.T) {
	exporter := NewExporter()

	collector := exporter.Collector()
	if collector == nil {
		t.Fatal("expected non-nil collector")
	}
}

func TestMetricsAdapter(t *testing.T) {
	registry := NewRegistry()
	adapter := NewMetricsAdapter(registry, "test", "")

	// Counter
	counter := adapter.Counter("requests")
	counter.Inc()
	counter.Add(5)
	if counter.Value() != 6 {
		t.Errorf("expected counter value 6, got %f", counter.Value())
	}

	// Gauge
	gauge := adapter.Gauge("connections")
	gauge.Set(10)
	gauge.Inc()
	gauge.Dec()
	if gauge.Value() != 10 {
		t.Errorf("expected gauge value 10, got %f", gauge.Value())
	}

	// Histogram
	histogram := adapter.Histogram("duration")
	histogram.Observe(1.0)
	histogram.Observe(2.0)
	if histogram.Count() != 2 {
		t.Errorf("expected histogram count 2, got %d", histogram.Count())
	}
	if histogram.Sum() != 3.0 {
		t.Errorf("expected histogram sum 3.0, got %f", histogram.Sum())
	}

	// Timer
	timer := adapter.Timer("latency")
	timer.ObserveDuration(100 * time.Millisecond)
	timer.Time(func() {
		time.Sleep(10 * time.Millisecond)
	})
}

func TestCollector(t *testing.T) {
	registry := NewRegistry()
	collector := NewCollector(registry, "test", "sub")

	// 自定义指标
	counter := collector.Counter("custom_counter", "Custom counter")
	counter.Inc()

	gauge := collector.Gauge("custom_gauge", "Custom gauge")
	gauge.Set(100)

	histogram := collector.Histogram("custom_histogram", "Custom histogram", nil)
	histogram.Observe(1.5)

	// 验证输出
	output := registry.Gather()
	if !strings.Contains(output, "test_sub_custom_counter") {
		t.Error("expected custom counter with prefix")
	}
}

func TestCollectorConvenienceMethods(t *testing.T) {
	registry := NewRegistry()
	collector := NewCollector(registry, "app", "")

	collector.RecordDuration("request", time.Second)
	collector.RecordCount("events", 10)
	collector.SetGaugeValue("active", 5)

	output := registry.Gather()
	if !strings.Contains(output, "app_request_seconds") {
		t.Error("expected duration metric")
	}
	if !strings.Contains(output, "app_events_total") {
		t.Error("expected count metric")
	}
	if !strings.Contains(output, "app_active") {
		t.Error("expected gauge metric")
	}
}

func TestDefaultBuckets(t *testing.T) {
	if len(DefaultBuckets) == 0 {
		t.Error("expected default buckets to be non-empty")
	}
}

func TestDefaultQuantiles(t *testing.T) {
	if len(DefaultQuantiles) == 0 {
		t.Error("expected default quantiles to be non-empty")
	}
}
