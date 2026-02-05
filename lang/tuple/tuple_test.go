package tuple

import (
	"testing"
)

func TestT2(t *testing.T) {
	t2 := T2("hello", 42)
	if t2.First != "hello" {
		t.Errorf("expected First to be 'hello', got %v", t2.First)
	}
	if t2.Second != 42 {
		t.Errorf("expected Second to be 42, got %v", t2.Second)
	}
}

func TestTuple2_Unpack(t *testing.T) {
	t2 := T2("hello", 42)
	first, second := t2.Unpack()
	if first != "hello" {
		t.Errorf("expected first to be 'hello', got %v", first)
	}
	if second != 42 {
		t.Errorf("expected second to be 42, got %v", second)
	}
}

func TestTuple2_Swap(t *testing.T) {
	t2 := T2("hello", 42)
	swapped := t2.Swap()
	if swapped.First != 42 {
		t.Errorf("expected First to be 42, got %v", swapped.First)
	}
	if swapped.Second != "hello" {
		t.Errorf("expected Second to be 'hello', got %v", swapped.Second)
	}
}

func TestT3(t *testing.T) {
	t3 := T3("hello", 42, true)
	if t3.First != "hello" {
		t.Errorf("expected First to be 'hello', got %v", t3.First)
	}
	if t3.Second != 42 {
		t.Errorf("expected Second to be 42, got %v", t3.Second)
	}
	if t3.Third != true {
		t.Errorf("expected Third to be true, got %v", t3.Third)
	}
}

func TestTuple3_Unpack(t *testing.T) {
	t3 := T3("hello", 42, true)
	first, second, third := t3.Unpack()
	if first != "hello" || second != 42 || third != true {
		t.Errorf("unexpected values: %v, %v, %v", first, second, third)
	}
}

func TestT4(t *testing.T) {
	t4 := T4("hello", 42, true, 3.14)
	if t4.First != "hello" {
		t.Errorf("expected First to be 'hello', got %v", t4.First)
	}
	if t4.Second != 42 {
		t.Errorf("expected Second to be 42, got %v", t4.Second)
	}
	if t4.Third != true {
		t.Errorf("expected Third to be true, got %v", t4.Third)
	}
	if t4.Fourth != 3.14 {
		t.Errorf("expected Fourth to be 3.14, got %v", t4.Fourth)
	}
}

func TestTuple4_Unpack(t *testing.T) {
	t4 := T4("hello", 42, true, 3.14)
	first, second, third, fourth := t4.Unpack()
	if first != "hello" || second != 42 || third != true || fourth != 3.14 {
		t.Errorf("unexpected values: %v, %v, %v, %v", first, second, third, fourth)
	}
}

func TestFromPair(t *testing.T) {
	pair := FromPair("key", "value")
	if pair.First != "key" || pair.Second != "value" {
		t.Errorf("unexpected pair: %v", pair)
	}
}

func TestZip2(t *testing.T) {
	tests := []struct {
		name     string
		s1       []string
		s2       []int
		expected []Tuple2[string, int]
	}{
		{
			name:     "相同长度",
			s1:       []string{"a", "b", "c"},
			s2:       []int{1, 2, 3},
			expected: []Tuple2[string, int]{{"a", 1}, {"b", 2}, {"c", 3}},
		},
		{
			name:     "第一个切片更短",
			s1:       []string{"a", "b"},
			s2:       []int{1, 2, 3},
			expected: []Tuple2[string, int]{{"a", 1}, {"b", 2}},
		},
		{
			name:     "第二个切片更短",
			s1:       []string{"a", "b", "c"},
			s2:       []int{1, 2},
			expected: []Tuple2[string, int]{{"a", 1}, {"b", 2}},
		},
		{
			name:     "空切片",
			s1:       []string{},
			s2:       []int{1, 2, 3},
			expected: nil,
		},
		{
			name:     "nil切片",
			s1:       nil,
			s2:       []int{1, 2, 3},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Zip2(tt.s1, tt.s2)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: expected %v, got %v", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestZip3(t *testing.T) {
	s1 := []string{"a", "b"}
	s2 := []int{1, 2, 3}
	s3 := []bool{true, false}
	result := Zip3(s1, s2, s3)

	if len(result) != 2 {
		t.Errorf("expected length 2, got %d", len(result))
		return
	}
	if result[0].First != "a" || result[0].Second != 1 || result[0].Third != true {
		t.Errorf("unexpected first tuple: %v", result[0])
	}
	if result[1].First != "b" || result[1].Second != 2 || result[1].Third != false {
		t.Errorf("unexpected second tuple: %v", result[1])
	}
}

func TestZip4(t *testing.T) {
	s1 := []string{"a", "b"}
	s2 := []int{1, 2}
	s3 := []bool{true, false}
	s4 := []float64{1.1, 2.2}
	result := Zip4(s1, s2, s3, s4)

	if len(result) != 2 {
		t.Errorf("expected length 2, got %d", len(result))
		return
	}
	if result[0].First != "a" || result[0].Second != 1 || result[0].Third != true || result[0].Fourth != 1.1 {
		t.Errorf("unexpected first tuple: %v", result[0])
	}
}

func TestUnzip2(t *testing.T) {
	tuples := []Tuple2[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}
	s1, s2 := Unzip2(tuples)

	if len(s1) != 3 || len(s2) != 3 {
		t.Errorf("unexpected lengths: %d, %d", len(s1), len(s2))
		return
	}
	expectedS1 := []string{"a", "b", "c"}
	expectedS2 := []int{1, 2, 3}
	for i := range s1 {
		if s1[i] != expectedS1[i] {
			t.Errorf("s1[%d]: expected %v, got %v", i, expectedS1[i], s1[i])
		}
		if s2[i] != expectedS2[i] {
			t.Errorf("s2[%d]: expected %v, got %v", i, expectedS2[i], s2[i])
		}
	}
}

func TestUnzip2_Empty(t *testing.T) {
	var tuples []Tuple2[string, int]
	s1, s2 := Unzip2(tuples)
	if s1 != nil || s2 != nil {
		t.Errorf("expected nil slices, got %v, %v", s1, s2)
	}
}

func TestUnzip3(t *testing.T) {
	tuples := []Tuple3[string, int, bool]{{"a", 1, true}, {"b", 2, false}}
	s1, s2, s3 := Unzip3(tuples)

	if len(s1) != 2 || len(s2) != 2 || len(s3) != 2 {
		t.Errorf("unexpected lengths: %d, %d, %d", len(s1), len(s2), len(s3))
		return
	}
	if s1[0] != "a" || s2[0] != 1 || s3[0] != true {
		t.Errorf("unexpected first elements: %v, %v, %v", s1[0], s2[0], s3[0])
	}
}

func TestUnzip4(t *testing.T) {
	tuples := []Tuple4[string, int, bool, float64]{{"a", 1, true, 1.1}, {"b", 2, false, 2.2}}
	s1, s2, s3, s4 := Unzip4(tuples)

	if len(s1) != 2 || len(s2) != 2 || len(s3) != 2 || len(s4) != 2 {
		t.Errorf("unexpected lengths")
		return
	}
	if s1[0] != "a" || s2[0] != 1 || s3[0] != true || s4[0] != 1.1 {
		t.Errorf("unexpected first elements")
	}
}

func TestZipWithIndex(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []Tuple2[int, string]
	}{
		{
			name:     "普通切片",
			input:    []string{"a", "b", "c"},
			expected: []Tuple2[int, string]{{0, "a"}, {1, "b"}, {2, "c"}},
		},
		{
			name:     "空切片",
			input:    []string{},
			expected: nil,
		},
		{
			name:     "nil切片",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZipWithIndex(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("index %d: expected %v, got %v", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestZipUnzipRoundTrip(t *testing.T) {
	// 测试 Zip 和 Unzip 的往返
	s1 := []string{"a", "b", "c"}
	s2 := []int{1, 2, 3}

	zipped := Zip2(s1, s2)
	r1, r2 := Unzip2(zipped)

	if len(r1) != len(s1) || len(r2) != len(s2) {
		t.Errorf("lengths don't match after round trip")
		return
	}
	for i := range s1 {
		if r1[i] != s1[i] {
			t.Errorf("r1[%d]: expected %v, got %v", i, s1[i], r1[i])
		}
		if r2[i] != s2[i] {
			t.Errorf("r2[%d]: expected %v, got %v", i, s2[i], r2[i])
		}
	}
}
