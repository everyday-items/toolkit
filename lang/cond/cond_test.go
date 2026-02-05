package cond

import (
	"testing"
)

func TestIf(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		trueVal   string
		falseVal  string
		expected  string
	}{
		{"条件为true", true, "yes", "no", "yes"},
		{"条件为false", false, "yes", "no", "no"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := If(tt.condition, tt.trueVal, tt.falseVal)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIfFunc(t *testing.T) {
	// 测试只有正确的分支被调用
	trueCalled := false
	falseCalled := false

	trueFn := func() string {
		trueCalled = true
		return "true"
	}
	falseFn := func() string {
		falseCalled = true
		return "false"
	}

	// 测试条件为 true
	result := IfFunc(true, trueFn, falseFn)
	if result != "true" {
		t.Errorf("expected 'true', got %v", result)
	}
	if !trueCalled {
		t.Error("trueFn should have been called")
	}
	if falseCalled {
		t.Error("falseFn should not have been called")
	}

	// 重置
	trueCalled = false
	falseCalled = false

	// 测试条件为 false
	result = IfFunc(false, trueFn, falseFn)
	if result != "false" {
		t.Errorf("expected 'false', got %v", result)
	}
	if trueCalled {
		t.Error("trueFn should not have been called")
	}
	if !falseCalled {
		t.Error("falseFn should have been called")
	}
}

func TestIfZero(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		defaultVal string
		expected   string
	}{
		{"非零值", "hello", "default", "hello"},
		{"零值", "", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IfZero(tt.value, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}

	// 测试整数
	if IfZero(0, 42) != 42 {
		t.Error("expected 42 for zero int")
	}
	if IfZero(10, 42) != 10 {
		t.Error("expected 10 for non-zero int")
	}
}

func TestIfZeroFunc(t *testing.T) {
	called := false
	defaultFn := func() int { called = true; return 42 }

	// 非零值，函数不应被调用
	result := IfZeroFunc(10, defaultFn)
	if result != 10 {
		t.Errorf("expected 10, got %v", result)
	}
	if called {
		t.Error("defaultFn should not have been called")
	}

	// 零值，函数应被调用
	result = IfZeroFunc(0, defaultFn)
	if result != 42 {
		t.Errorf("expected 42, got %v", result)
	}
	if !called {
		t.Error("defaultFn should have been called")
	}
}

func TestCoalesce(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected string
	}{
		{"第一个非零", []string{"hello", "world"}, "hello"},
		{"第二个非零", []string{"", "world"}, "world"},
		{"第三个非零", []string{"", "", "third"}, "third"},
		{"全为零", []string{"", "", ""}, ""},
		{"空参数", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Coalesce(tt.values...)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}

	// 测试整数
	if Coalesce(0, 0, 3) != 3 {
		t.Error("expected 3")
	}
}

func TestCoalesceFunc(t *testing.T) {
	callOrder := []int{}

	fn1 := func() string { callOrder = append(callOrder, 1); return "" }
	fn2 := func() string { callOrder = append(callOrder, 2); return "second" }
	fn3 := func() string { callOrder = append(callOrder, 3); return "third" }

	result := CoalesceFunc(fn1, fn2, fn3)
	if result != "second" {
		t.Errorf("expected 'second', got %v", result)
	}
	// 第三个函数不应被调用
	if len(callOrder) != 2 || callOrder[0] != 1 || callOrder[1] != 2 {
		t.Errorf("unexpected call order: %v", callOrder)
	}
}

func TestIfNil(t *testing.T) {
	value := "hello"
	ptr := &value

	// 非 nil 指针
	result := IfNil(ptr, "default")
	if result != "hello" {
		t.Errorf("expected 'hello', got %v", result)
	}

	// nil 指针
	var nilPtr *string
	result = IfNil(nilPtr, "default")
	if result != "default" {
		t.Errorf("expected 'default', got %v", result)
	}
}

func TestIfNilFunc(t *testing.T) {
	called := false
	defaultFn := func() string { called = true; return "default" }

	value := "hello"
	ptr := &value

	// 非 nil 指针，函数不应被调用
	result := IfNilFunc(ptr, defaultFn)
	if result != "hello" {
		t.Errorf("expected 'hello', got %v", result)
	}
	if called {
		t.Error("defaultFn should not have been called")
	}

	// nil 指针，函数应被调用
	var nilPtr *string
	result = IfNilFunc(nilPtr, defaultFn)
	if result != "default" {
		t.Errorf("expected 'default', got %v", result)
	}
	if !called {
		t.Error("defaultFn should have been called")
	}
}

func TestUnless(t *testing.T) {
	// 条件为 false，返回第一个值
	result := Unless(false, "first", "second")
	if result != "first" {
		t.Errorf("expected 'first', got %v", result)
	}

	// 条件为 true，返回第二个值
	result = Unless(true, "first", "second")
	if result != "second" {
		t.Errorf("expected 'second', got %v", result)
	}
}

func TestSwitch(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{"匹配第一个", 200, "OK"},
		{"匹配第二个", 404, "Not Found"},
		{"无匹配", 500, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Switch[int, string](tt.value).
				Case(200, "OK").
				Case(404, "Not Found").
				Default("Unknown")
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSwitch_CaseFunc(t *testing.T) {
	callCount := 0
	fn := func() string {
		callCount++
		return "computed"
	}

	// 匹配的情况，函数被调用
	result := Switch[int, string](1).
		CaseFunc(1, fn).
		Default("default")
	if result != "computed" {
		t.Errorf("expected 'computed', got %v", result)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// 不匹配的情况，函数不被调用
	callCount = 0
	result = Switch[int, string](2).
		CaseFunc(1, fn).
		Default("default")
	if result != "default" {
		t.Errorf("expected 'default', got %v", result)
	}
	if callCount != 0 {
		t.Errorf("expected 0 calls, got %d", callCount)
	}
}

func TestSwitch_CaseIn(t *testing.T) {
	result := Switch[string, string](".jpg").
		CaseIn("图片", ".jpg", ".png", ".gif").
		CaseIn("文档", ".doc", ".pdf").
		Default("未知")

	if result != "图片" {
		t.Errorf("expected '图片', got %v", result)
	}

	result = Switch[string, string](".pdf").
		CaseIn("图片", ".jpg", ".png", ".gif").
		CaseIn("文档", ".doc", ".pdf").
		Default("未知")

	if result != "文档" {
		t.Errorf("expected '文档', got %v", result)
	}

	result = Switch[string, string](".exe").
		CaseIn("图片", ".jpg", ".png", ".gif").
		CaseIn("文档", ".doc", ".pdf").
		Default("未知")

	if result != "未知" {
		t.Errorf("expected '未知', got %v", result)
	}
}

func TestSwitch_Matched(t *testing.T) {
	s := Switch[int, string](200).
		Case(200, "OK")

	if !s.Matched() {
		t.Error("expected Matched to be true")
	}

	s2 := Switch[int, string](500).
		Case(200, "OK")

	if s2.Matched() {
		t.Error("expected Matched to be false")
	}
}

func TestSwitch_Result(t *testing.T) {
	s := Switch[int, string](200).
		Case(200, "OK")

	if s.Result() != "OK" {
		t.Errorf("expected 'OK', got %v", s.Result())
	}

	s2 := Switch[int, string](500).
		Case(200, "OK")

	if s2.Result() != "" {
		t.Errorf("expected empty string, got %v", s2.Result())
	}
}

func TestSwitchTrue(t *testing.T) {
	score := 85
	grade := SwitchTrue[string]().
		When(score >= 90, "A").
		When(score >= 80, "B").
		When(score >= 60, "C").
		Default("F")

	if grade != "B" {
		t.Errorf("expected 'B', got %v", grade)
	}

	// 测试低分
	score = 50
	grade = SwitchTrue[string]().
		When(score >= 90, "A").
		When(score >= 80, "B").
		When(score >= 60, "C").
		Default("F")

	if grade != "F" {
		t.Errorf("expected 'F', got %v", grade)
	}
}

func TestSwitchTrue_WhenFunc(t *testing.T) {
	callCount := 0
	fn := func() string {
		callCount++
		return "computed"
	}

	result := SwitchTrue[string]().
		WhenFunc(true, fn).
		Default("default")

	if result != "computed" {
		t.Errorf("expected 'computed', got %v", result)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// 条件为 false，函数不被调用
	callCount = 0
	result = SwitchTrue[string]().
		WhenFunc(false, fn).
		Default("default")

	if result != "default" {
		t.Errorf("expected 'default', got %v", result)
	}
	if callCount != 0 {
		t.Errorf("expected 0 calls, got %d", callCount)
	}
}

func TestSwitch_FirstMatchWins(t *testing.T) {
	// 验证第一个匹配的 Case 生效
	result := Switch[int, string](1).
		Case(1, "first").
		Case(1, "second").
		Default("default")

	if result != "first" {
		t.Errorf("expected 'first', got %v", result)
	}
}

func TestSwitch_DefaultFunc(t *testing.T) {
	callCount := 0
	defaultFn := func() string {
		callCount++
		return "computed default"
	}

	// 有匹配时，DefaultFunc 不被调用
	result := Switch[int, string](200).
		Case(200, "OK").
		DefaultFunc(defaultFn)
	if result != "OK" {
		t.Errorf("expected 'OK', got %v", result)
	}
	if callCount != 0 {
		t.Errorf("expected 0 calls, got %d", callCount)
	}

	// 无匹配时，DefaultFunc 被调用
	result = Switch[int, string](500).
		Case(200, "OK").
		DefaultFunc(defaultFn)
	if result != "computed default" {
		t.Errorf("expected 'computed default', got %v", result)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestSwitchTrue_Result(t *testing.T) {
	s := SwitchTrue[string]().
		When(true, "matched")

	if s.Result() != "matched" {
		t.Errorf("expected 'matched', got %v", s.Result())
	}

	s2 := SwitchTrue[string]().
		When(false, "not matched")

	if s2.Result() != "" {
		t.Errorf("expected empty string, got %v", s2.Result())
	}
}
