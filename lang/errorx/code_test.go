package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestNewCodedError(t *testing.T) {
	err := NewCodedError(CodeNotFound, DomainGeneral, "用户不存在")
	if err.Code != CodeNotFound {
		t.Errorf("Code 不匹配: %d", err.Code)
	}
	if err.Domain != DomainGeneral {
		t.Errorf("Domain 不匹配: %s", err.Domain)
	}
	if err.Message != "用户不存在" {
		t.Errorf("Message 不匹配: %s", err.Message)
	}
}

func TestCodedError_Error(t *testing.T) {
	err := NewCodedError(CodeNotFound, DomainGeneral, "用户不存在")
	expected := "[GENERAL-1002] 用户不存在"
	if err.Error() != expected {
		t.Errorf("Error() 不匹配: 期望 %q, 实际 %q", expected, err.Error())
	}
}

func TestCodedError_WithCause(t *testing.T) {
	cause := fmt.Errorf("底层错误")
	err := NewCodedError(CodeInternal, DomainGeneral, "操作失败").WithCause(cause)
	if !errors.Is(err, cause) {
		t.Error("WithCause 后 errors.Is 应该匹配")
	}
	if err.Unwrap() != cause {
		t.Error("Unwrap 应返回底层错误")
	}
}

func TestCodedError_WithDetails(t *testing.T) {
	err := NewCodedError(CodeNotFound, DomainGeneral, "未找到").
		WithDetails("id", "123").
		WithDetails("type", "user")

	if err.Details["id"] != "123" {
		t.Error("Details 应包含 id")
	}
	if err.Details["type"] != "user" {
		t.Error("Details 应包含 type")
	}
}

func TestCodedError_HTTPStatus(t *testing.T) {
	tests := []struct {
		code     int
		expected int
	}{
		{CodeOK, http.StatusOK},
		{CodeInvalidInput, http.StatusBadRequest},
		{CodeNotFound, http.StatusNotFound},
		{CodeUnauthorized, http.StatusUnauthorized},
		{CodeForbidden, http.StatusForbidden},
		{CodeTimeout, http.StatusGatewayTimeout},
		{CodeRateLimit, http.StatusTooManyRequests},
		{CodeBudgetExceeded, http.StatusPaymentRequired},
		{CodeInternal, http.StatusInternalServerError},
		{CodeUnknown, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		err := NewCodedError(tt.code, DomainGeneral, "test")
		if err.HTTPStatus() != tt.expected {
			t.Errorf("Code %d: 期望 HTTP %d, 实际 %d", tt.code, tt.expected, err.HTTPStatus())
		}
	}
}

func TestCodedError_ToJSON(t *testing.T) {
	err := NewCodedError(CodeNotFound, DomainGeneral, "未找到").
		WithDetails("id", "123")
	m := err.ToJSON()
	if m["code"] != CodeNotFound {
		t.Error("ToJSON code 不匹配")
	}
	if m["domain"] != DomainGeneral {
		t.Error("ToJSON domain 不匹配")
	}
}

func TestIsCodedError(t *testing.T) {
	ce := ErrNotFound("test")
	wrapped := fmt.Errorf("包装: %w", ce)

	found, ok := IsCodedError(wrapped)
	if !ok {
		t.Fatal("IsCodedError 应该返回 true")
	}
	if found.Code != CodeNotFound {
		t.Errorf("Code 不匹配: %d", found.Code)
	}

	_, ok = IsCodedError(fmt.Errorf("普通错误"))
	if ok {
		t.Error("普通错误不应该是 CodedError")
	}
}

func TestConvenienceConstructors(t *testing.T) {
	tests := []struct {
		fn     func(string) *CodedError
		code   int
		domain string
	}{
		{ErrInvalidInput, CodeInvalidInput, DomainGeneral},
		{ErrNotFound, CodeNotFound, DomainGeneral},
		{ErrUnauthorized, CodeUnauthorized, DomainGeneral},
		{ErrForbidden, CodeForbidden, DomainGeneral},
		{ErrInternal, CodeInternal, DomainGeneral},
		{ErrTimeout, CodeTimeout, DomainGeneral},
		{ErrLLM, CodeLLMError, DomainAI},
		{ErrBudgetExceeded, CodeBudgetExceeded, DomainAI},
		{ErrToolError, CodeToolError, DomainAgent},
		{ErrSkillNotFound, CodeSkillNotFound, DomainAgent},
	}
	for _, tt := range tests {
		err := tt.fn("test")
		if err.Code != tt.code || err.Domain != tt.domain {
			t.Errorf("构造函数 code=%d domain=%s 不匹配", err.Code, err.Domain)
		}
	}
}
