package observe

import (
	"context"
	"testing"
)

func TestSpanKind(t *testing.T) {
	kinds := []SpanKind{
		SpanKindInternal,
		SpanKindServer,
		SpanKindClient,
		SpanKindProducer,
		SpanKindConsumer,
	}

	for i, kind := range kinds {
		if int(kind) != i {
			t.Errorf("expected SpanKind %d, got %d", i, kind)
		}
	}
}

func TestStatusCode(t *testing.T) {
	codes := []StatusCode{
		StatusCodeUnset,
		StatusCodeOK,
		StatusCodeError,
	}

	for i, code := range codes {
		if int(code) != i {
			t.Errorf("expected StatusCode %d, got %d", i, code)
		}
	}
}

func TestTokenUsage(t *testing.T) {
	usage := TokenUsage{
		PromptTokens:     100,
		CompletionTokens: 200,
		TotalTokens:      300,
	}

	if usage.PromptTokens != 100 {
		t.Errorf("expected PromptTokens 100, got %d", usage.PromptTokens)
	}

	if usage.CompletionTokens != 200 {
		t.Errorf("expected CompletionTokens 200, got %d", usage.CompletionTokens)
	}

	if usage.TotalTokens != 300 {
		t.Errorf("expected TotalTokens 300, got %d", usage.TotalTokens)
	}
}

func TestSpanOptions(t *testing.T) {
	cfg := &SpanConfig{
		Attributes: make(map[string]any),
	}

	WithSpanKind(SpanKindServer)(cfg)
	if cfg.Kind != SpanKindServer {
		t.Errorf("expected SpanKindServer, got %v", cfg.Kind)
	}

	attrs := map[string]any{"key": "value"}
	WithAttributes(attrs)(cfg)
	if cfg.Attributes["key"] != "value" {
		t.Errorf("expected attribute 'key' = 'value'")
	}
}

func TestContextWithSpan(t *testing.T) {
	ctx := context.Background()

	// 空 context 应该返回 nil
	span := SpanFromContext(ctx)
	if span != nil {
		t.Error("expected nil span from empty context")
	}
}

func TestAttrConstants(t *testing.T) {
	// 验证属性常量不为空
	attrs := []string{
		AttrServiceName,
		AttrServiceVersion,
		AttrAgentID,
		AttrAgentName,
		AttrAgentType,
		AttrLLMProvider,
		AttrLLMModel,
		AttrLLMPromptTokens,
		AttrLLMCompletionTokens,
		AttrLLMTotalTokens,
		AttrToolName,
		AttrToolInput,
		AttrToolOutput,
		AttrErrorType,
		AttrErrorMessage,
		AttrRetrieverType,
		AttrRetrieverTopK,
		AttrDocumentCount,
	}

	for _, attr := range attrs {
		if attr == "" {
			t.Error("attribute constant should not be empty")
		}
	}
}
