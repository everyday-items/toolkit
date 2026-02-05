package slicex

import (
	"testing"
)

// Partition tests
func TestPartition(t *testing.T) {
	even, odd := Partition([]int{1, 2, 3, 4, 5}, func(n int) bool {
		return n%2 == 0
	})
	if len(even) != 2 || even[0] != 2 || even[1] != 4 {
		t.Errorf("expected [2, 4], got %v", even)
	}
	if len(odd) != 3 || odd[0] != 1 || odd[1] != 3 || odd[2] != 5 {
		t.Errorf("expected [1, 3, 5], got %v", odd)
	}
}

func TestPartition_Empty(t *testing.T) {
	even, odd := Partition([]int{}, func(n int) bool { return n%2 == 0 })
	if even != nil || odd != nil {
		t.Errorf("expected nil, nil for empty slice")
	}
}

func TestPartitionBy(t *testing.T) {
	type User struct {
		Name string
		Dept string
	}
	users := []User{{"Alice", "IT"}, {"Bob", "HR"}, {"Charlie", "IT"}}
	groups := PartitionBy(users, func(u User) string { return u.Dept })

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
	if len(groups["IT"]) != 2 {
		t.Errorf("expected 2 IT users, got %d", len(groups["IT"]))
	}
	if len(groups["HR"]) != 1 {
		t.Errorf("expected 1 HR user, got %d", len(groups["HR"]))
	}
}

// Concat tests
func TestConcat(t *testing.T) {
	result := Concat([]int{1, 2}, []int{3, 4}, []int{5})
	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

func TestConcat_Empty(t *testing.T) {
	result := Concat[int]()
	if result != nil {
		t.Errorf("expected nil for empty concat, got %v", result)
	}
}

func TestFlatten(t *testing.T) {
	result := Flatten([][]int{{1, 2}, {3, 4}, {5}})
	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestPrepend(t *testing.T) {
	result := Prepend([]int{3, 4}, 1, 2)
	expected := []int{1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

// Range tests
func TestRange(t *testing.T) {
	result := Range(0, 5, 1)
	expected := []int{0, 1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

func TestRange_Step2(t *testing.T) {
	result := Range(0, 10, 2)
	expected := []int{0, 2, 4, 6, 8}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestRange_Negative(t *testing.T) {
	result := Range(5, 0, -1)
	expected := []int{5, 4, 3, 2, 1}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], result[i])
		}
	}
}

func TestRangeN(t *testing.T) {
	result := RangeN(5)
	expected := []int{0, 1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestRangeFrom(t *testing.T) {
	result := RangeFrom(5, 3)
	expected := []int{5, 6, 7}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

func TestRepeat(t *testing.T) {
	result := Repeat("hello", 3)
	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
	for _, v := range result {
		if v != "hello" {
			t.Errorf("expected 'hello', got %v", v)
		}
	}
}

func TestRepeatFunc(t *testing.T) {
	result := RepeatFunc(3, func(i int) int { return i * 2 })
	expected := []int{0, 2, 4}
	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}
}

// Aggregate tests
func TestMin(t *testing.T) {
	result := Min([]int{3, 1, 4, 1, 5})
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}
}

func TestMax(t *testing.T) {
	result := Max([]int{3, 1, 4, 1, 5})
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

func TestMinMax(t *testing.T) {
	min, max := MinMax([]int{3, 1, 4, 1, 5})
	if min != 1 || max != 5 {
		t.Errorf("expected 1, 5, got %d, %d", min, max)
	}
}

func TestMinBy(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	users := []User{{"Alice", 30}, {"Bob", 20}, {"Charlie", 25}}
	youngest, ok := MinBy(users, func(a, b User) bool { return a.Age < b.Age })
	if !ok || youngest.Name != "Bob" {
		t.Errorf("expected Bob, got %v", youngest)
	}
}

func TestMaxBy(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	users := []User{{"Alice", 30}, {"Bob", 20}, {"Charlie", 25}}
	oldest, ok := MaxBy(users, func(a, b User) bool { return a.Age < b.Age })
	if !ok || oldest.Name != "Alice" {
		t.Errorf("expected Alice, got %v", oldest)
	}
}

func TestMinByKey(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	users := []User{{"Alice", 30}, {"Bob", 20}, {"Charlie", 25}}
	youngest, ok := MinByKey(users, func(u User) int { return u.Age })
	if !ok || youngest.Name != "Bob" {
		t.Errorf("expected Bob, got %v", youngest)
	}
}

func TestSum(t *testing.T) {
	result := Sum([]int{1, 2, 3, 4, 5})
	if result != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

func TestSumBy(t *testing.T) {
	type Item struct {
		Price int
	}
	items := []Item{{10}, {20}, {30}}
	total := SumBy(items, func(i Item) int { return i.Price })
	if total != 60 {
		t.Errorf("expected 60, got %d", total)
	}
}

func TestAverage(t *testing.T) {
	result := Average([]int{1, 2, 3, 4, 5})
	if result != 3.0 {
		t.Errorf("expected 3.0, got %f", result)
	}
}

func TestProduct(t *testing.T) {
	result := Product([]int{1, 2, 3, 4})
	if result != 24 {
		t.Errorf("expected 24, got %d", result)
	}
}

func TestCountValue(t *testing.T) {
	result := CountValue([]int{1, 2, 2, 3, 2}, 2)
	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

// Shuffle tests
func TestShuffle(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	slice := make([]int, len(original))
	copy(slice, original)
	Shuffle(slice)

	// 验证元素完整性
	sum := 0
	for _, v := range slice {
		sum += v
	}
	if sum != 15 {
		t.Errorf("shuffle changed elements, sum: %d", sum)
	}
}

func TestShuffleCopy(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	shuffled := ShuffleCopy(original)

	// 验证原切片未变
	for i, v := range original {
		if v != i+1 {
			t.Errorf("original slice was modified")
		}
	}

	// 验证元素完整性
	sum := 0
	for _, v := range shuffled {
		sum += v
	}
	if sum != 15 {
		t.Errorf("shuffle changed elements, sum: %d", sum)
	}
}

func TestSample(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	sample := Sample(slice, 3)
	if len(sample) != 3 {
		t.Errorf("expected length 3, got %d", len(sample))
	}

	// 验证无重复
	seen := make(map[int]bool)
	for _, v := range sample {
		if seen[v] {
			t.Errorf("duplicate value in sample: %d", v)
		}
		seen[v] = true
	}
}

func TestSampleOne(t *testing.T) {
	slice := []int{1, 2, 3}
	v, ok := SampleOne(slice)
	if !ok {
		t.Error("expected ok to be true")
	}
	if v < 1 || v > 3 {
		t.Errorf("unexpected value: %d", v)
	}

	_, ok = SampleOne([]int{})
	if ok {
		t.Error("expected ok to be false for empty slice")
	}
}

// Channel tests
func TestToChannel(t *testing.T) {
	slice := []int{1, 2, 3}
	ch := ToChannel(slice)
	result := FromChannel(ch)

	if len(result) != len(slice) {
		t.Errorf("expected length %d, got %d", len(slice), len(result))
	}
	for i := range result {
		if result[i] != slice[i] {
			t.Errorf("index %d: expected %d, got %d", i, slice[i], result[i])
		}
	}
}

func TestFromChannelN(t *testing.T) {
	ch := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)

	result := FromChannelN(ch, 5)
	if len(result) != 5 {
		t.Errorf("expected length 5, got %d", len(result))
	}
}

// Search tests
func TestLastIndexOf(t *testing.T) {
	result := LastIndexOf([]int{1, 2, 3, 2, 1}, 2)
	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}

	result = LastIndexOf([]int{1, 2, 3}, 4)
	if result != -1 {
		t.Errorf("expected -1, got %d", result)
	}
}

func TestLastIndexOfFunc(t *testing.T) {
	result := LastIndexOfFunc([]int{1, 2, 3, 4}, func(n int) bool {
		return n%2 == 0
	})
	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

func TestFirst(t *testing.T) {
	v, ok := First([]int{1, 2, 3})
	if !ok || v != 1 {
		t.Errorf("expected 1, true, got %d, %v", v, ok)
	}

	_, ok = First([]int{})
	if ok {
		t.Error("expected false for empty slice")
	}
}

func TestLast(t *testing.T) {
	v, ok := Last([]int{1, 2, 3})
	if !ok || v != 3 {
		t.Errorf("expected 3, true, got %d, %v", v, ok)
	}

	_, ok = Last([]int{})
	if ok {
		t.Error("expected false for empty slice")
	}
}

func TestFirstOr(t *testing.T) {
	v := FirstOr([]int{1, 2, 3}, 0)
	if v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	v = FirstOr([]int{}, 0)
	if v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}

func TestLastOr(t *testing.T) {
	v := LastOr([]int{1, 2, 3}, 0)
	if v != 3 {
		t.Errorf("expected 3, got %d", v)
	}

	v = LastOr([]int{}, 0)
	if v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}

func TestNth(t *testing.T) {
	// 正索引
	v, ok := Nth([]int{1, 2, 3}, 1)
	if !ok || v != 2 {
		t.Errorf("expected 2, true, got %d, %v", v, ok)
	}

	// 负索引
	v, ok = Nth([]int{1, 2, 3}, -1)
	if !ok || v != 3 {
		t.Errorf("expected 3, true, got %d, %v", v, ok)
	}

	// 越界
	_, ok = Nth([]int{1, 2, 3}, 10)
	if ok {
		t.Error("expected false for out of bounds")
	}
}
