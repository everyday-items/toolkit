package timex

import (
	"testing"
	"time"
)

// Duration tests
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"零", 0, "0s"},
		{"秒", time.Second * 5, "5s"},
		{"分钟", time.Minute * 3, "3m"},
		{"小时", time.Hour * 2, "2h"},
		{"天", time.Hour * 24, "1d"},
		{"组合", time.Hour*2 + time.Minute*30 + time.Second*5, "2h30m5s"},
		{"天和小时", time.Hour * 26, "1d2h"},
		{"毫秒", time.Millisecond * 500, "500ms"},
		{"负数", -time.Hour, "-1h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatDurationShort(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"零", 0, "0s"},
		{"秒", time.Second * 90, "1m30s"},
		{"天和小时", time.Hour*26 + time.Minute*30, "1d2h"},
		{"小时和分钟", time.Hour*2 + time.Minute*30 + time.Second*45, "2h30m"},
		{"只有天", time.Hour * 48, "2d"},
		{"只有小时", time.Hour * 5, "5h"},
		{"负数", -time.Hour, "-1h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDurationShort(tt.duration)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{"天", "1d", time.Hour * 24, false},
		{"天和小时", "1d2h", time.Hour * 26, false},
		{"完整格式", "1d2h3m4s", time.Hour*26 + time.Minute*3 + time.Second*4, false},
		{"标准Go格式", "2h30m", time.Hour*2 + time.Minute*30, false},
		{"毫秒", "500ms", time.Millisecond * 500, false},
		{"天和毫秒", "1d500ms", time.Hour*24 + time.Millisecond*500, false},
		{"负数", "-1d", -time.Hour * 24, false},
		{"空字符串", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDuration(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMustParseDuration(t *testing.T) {
	d := MustParseDuration("1d")
	if d != time.Hour*24 {
		t.Errorf("expected 24h, got %v", d)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	MustParseDuration("invalid")
}

func TestDurationRound(t *testing.T) {
	d := time.Second * 90
	rounded := DurationRound(d, time.Minute)
	if rounded != time.Minute*2 {
		t.Errorf("expected 2m, got %v", rounded)
	}
}

func TestDurationTruncate(t *testing.T) {
	d := time.Second * 90
	truncated := DurationTruncate(d, time.Minute)
	if truncated != time.Minute {
		t.Errorf("expected 1m, got %v", truncated)
	}
}

// Location tests
func TestShanghai(t *testing.T) {
	loc := Shanghai()
	if loc == nil {
		t.Error("expected non-nil location")
	}

	// 多次调用应返回相同实例
	loc2 := Shanghai()
	if loc != loc2 {
		t.Error("expected same instance")
	}
}

func TestBeijing(t *testing.T) {
	loc := Beijing()
	if loc == nil {
		t.Error("expected non-nil location")
	}

	// 北京和上海应该是相同时区
	if loc != Shanghai() {
		t.Error("expected Beijing and Shanghai to be same location")
	}
}

func TestTokyo(t *testing.T) {
	loc := Tokyo()
	if loc == nil {
		t.Error("expected non-nil location")
	}
}

func TestNewYork(t *testing.T) {
	loc := NewYork()
	if loc == nil {
		t.Error("expected non-nil location")
	}
}

func TestLondon(t *testing.T) {
	loc := London()
	if loc == nil {
		t.Error("expected non-nil location")
	}
}

func TestUTC(t *testing.T) {
	loc := UTC()
	if loc != time.UTC {
		t.Error("expected time.UTC")
	}
}

func TestInShanghai(t *testing.T) {
	utcTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	shanghaiTime := InShanghai(utcTime)

	// 上海比 UTC 快 8 小时
	expected := time.Date(2024, 1, 1, 8, 0, 0, 0, Shanghai())
	if !shanghaiTime.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, shanghaiTime)
	}
}

func TestNowShanghai(t *testing.T) {
	now := NowShanghai()
	if now.Location() != Shanghai() {
		t.Error("expected Shanghai timezone")
	}
}

func TestNowBeijing(t *testing.T) {
	now := NowBeijing()
	if now.Location() != Beijing() {
		t.Error("expected Beijing timezone")
	}
}

func TestNowTokyo(t *testing.T) {
	now := NowTokyo()
	if now.Location() != Tokyo() {
		t.Error("expected Tokyo timezone")
	}
}

func TestNowNewYork(t *testing.T) {
	now := NowNewYork()
	if now.Location() != NewYork() {
		t.Error("expected NewYork timezone")
	}
}

func TestNowUTC(t *testing.T) {
	now := NowUTC()
	if now.Location() != time.UTC {
		t.Error("expected UTC timezone")
	}
}

func TestParseInShanghai(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	value := "2024-01-29 10:00:00"

	tm, err := ParseInShanghai(layout, value)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if tm.Location() != Shanghai() {
		t.Error("expected Shanghai timezone")
	}

	if tm.Hour() != 10 {
		t.Errorf("expected hour 10, got %d", tm.Hour())
	}
}

func TestFixedZone(t *testing.T) {
	cst := FixedZone("CST", 8)
	if cst == nil {
		t.Error("expected non-nil location")
	}

	// 测试偏移
	utcTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	cstTime := utcTime.In(cst)

	if cstTime.Hour() != 8 {
		t.Errorf("expected hour 8, got %d", cstTime.Hour())
	}
}

func TestFixedZone_Negative(t *testing.T) {
	est := FixedZone("EST", -5)
	if est == nil {
		t.Error("expected non-nil location")
	}

	utcTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	estTime := utcTime.In(est)

	if estTime.Hour() != 5 {
		t.Errorf("expected hour 5, got %d", estTime.Hour())
	}
}
