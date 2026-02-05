package mapx

import (
	"testing"
)

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := Keys(m)

	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[k] = true
	}

	for k := range m {
		if !keySet[k] {
			t.Errorf("key %s not found", k)
		}
	}

	// Test nil map
	if Keys[string, int](nil) != nil {
		t.Error("expected nil for nil map")
	}
}

func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	values := Values(m)

	if len(values) != 3 {
		t.Errorf("expected 3 values, got %d", len(values))
	}

	sum := 0
	for _, v := range values {
		sum += v
	}
	if sum != 6 {
		t.Errorf("expected sum 6, got %d", sum)
	}

	// Test nil map
	if Values[string, int](nil) != nil {
		t.Error("expected nil for nil map")
	}
}

func TestEntries(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	entries := Entries(m)

	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	// Test nil map
	if Entries[string, int](nil) != nil {
		t.Error("expected nil for nil map")
	}
}

func TestFromEntries(t *testing.T) {
	entries := []Entry[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
	}
	m := FromEntries(entries)

	if m["a"] != 1 || m["b"] != 2 {
		t.Error("FromEntries failed")
	}

	// Test nil entries
	if FromEntries[string, int](nil) != nil {
		t.Error("expected nil for nil entries")
	}
}

func TestFilter(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	filtered := Filter(m, func(k string, v int) bool {
		return v > 1
	})

	if len(filtered) != 2 {
		t.Errorf("expected 2 items, got %d", len(filtered))
	}

	if _, ok := filtered["a"]; ok {
		t.Error("'a' should not be in filtered map")
	}

	// Test nil map
	if Filter[string, int](nil, func(k string, v int) bool { return true }) != nil {
		t.Error("expected nil for nil map")
	}
}

func TestFilterKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	filtered := FilterKeys(m, func(k string) bool {
		return k != "a"
	})

	if len(filtered) != 2 {
		t.Errorf("expected 2 items, got %d", len(filtered))
	}
}

func TestFilterValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	filtered := FilterValues(m, func(v int) bool {
		return v >= 2
	})

	if len(filtered) != 2 {
		t.Errorf("expected 2 items, got %d", len(filtered))
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := MapValues(m, func(v int) int {
		return v * 2
	})

	if result["a"] != 2 || result["b"] != 4 {
		t.Error("MapValues failed")
	}

	// Test nil map
	if MapValues[string, int, int](nil, func(v int) int { return v }) != nil {
		t.Error("expected nil for nil map")
	}
}

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := MapKeys(m, func(k string) string {
		return k + k
	})

	if result["aa"] != 1 || result["bb"] != 2 {
		t.Error("MapKeys failed")
	}
}

func TestMerge(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}

	result := Merge(m1, m2)

	if result["a"] != 1 || result["b"] != 3 || result["c"] != 4 {
		t.Error("Merge failed")
	}
}

func TestMergeWith(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}

	result := MergeWith(func(v1, v2 int) int {
		return v1 + v2
	}, m1, m2)

	if result["a"] != 1 || result["b"] != 5 || result["c"] != 4 {
		t.Error("MergeWith failed")
	}
}

func TestInvert(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := Invert(m)

	if result[1] != "a" || result[2] != "b" {
		t.Error("Invert failed")
	}

	// Test nil map
	if Invert[string, int](nil) != nil {
		t.Error("expected nil for nil map")
	}
}

func TestPick(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := Pick(m, "a", "c")

	if len(result) != 2 || result["a"] != 1 || result["c"] != 3 {
		t.Error("Pick failed")
	}

	// Test picking non-existent key
	result = Pick(m, "a", "x")
	if len(result) != 1 {
		t.Error("Pick should ignore non-existent keys")
	}
}

func TestOmit(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := Omit(m, "a", "c")

	if len(result) != 1 || result["b"] != 2 {
		t.Error("Omit failed")
	}
}

func TestContains(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	if !Contains(m, "a") {
		t.Error("Contains should return true for existing key")
	}

	if Contains(m, "x") {
		t.Error("Contains should return false for non-existent key")
	}
}

func TestContainsAll(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	if !ContainsAll(m, "a", "b") {
		t.Error("ContainsAll should return true")
	}

	if ContainsAll(m, "a", "x") {
		t.Error("ContainsAll should return false")
	}
}

func TestContainsAny(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	if !ContainsAny(m, "a", "x") {
		t.Error("ContainsAny should return true")
	}

	if ContainsAny(m, "x", "y") {
		t.Error("ContainsAny should return false")
	}
}

func TestGetOrDefault(t *testing.T) {
	m := map[string]int{"a": 1}

	if GetOrDefault(m, "a", 0) != 1 {
		t.Error("GetOrDefault should return existing value")
	}

	if GetOrDefault(m, "x", 99) != 99 {
		t.Error("GetOrDefault should return default value")
	}
}

func TestGetOrCompute(t *testing.T) {
	m := map[string]int{"a": 1}

	// Existing key
	v := GetOrCompute(m, "a", func() int { return 99 })
	if v != 1 {
		t.Error("GetOrCompute should return existing value")
	}

	// New key
	v = GetOrCompute(m, "b", func() int { return 2 })
	if v != 2 || m["b"] != 2 {
		t.Error("GetOrCompute should compute and store new value")
	}
}

func TestClone(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	cloned := Clone(m)

	if cloned["a"] != 1 || cloned["b"] != 2 {
		t.Error("Clone failed")
	}

	// Modify original
	m["a"] = 99
	if cloned["a"] != 1 {
		t.Error("Clone should create independent copy")
	}

	// Test nil map
	if Clone[string, int](nil) != nil {
		t.Error("expected nil for nil map")
	}
}

func TestEqual(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}
	m3 := map[string]int{"a": 1, "b": 3}
	m4 := map[string]int{"a": 1}

	if !Equal(m1, m2) {
		t.Error("Equal should return true for equal maps")
	}

	if Equal(m1, m3) {
		t.Error("Equal should return false for different values")
	}

	if Equal(m1, m4) {
		t.Error("Equal should return false for different lengths")
	}
}

func TestIsEmpty(t *testing.T) {
	if !IsEmpty(map[string]int{}) {
		t.Error("IsEmpty should return true for empty map")
	}

	if IsEmpty(map[string]int{"a": 1}) {
		t.Error("IsEmpty should return false for non-empty map")
	}
}

func TestForEach(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	sum := 0
	ForEach(m, func(k string, v int) {
		sum += v
	})

	if sum != 3 {
		t.Error("ForEach failed")
	}
}

func TestAny(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	if !Any(m, func(k string, v int) bool { return v > 2 }) {
		t.Error("Any should return true")
	}

	if Any(m, func(k string, v int) bool { return v > 10 }) {
		t.Error("Any should return false")
	}
}

func TestAll(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	if !All(m, func(k string, v int) bool { return v > 0 }) {
		t.Error("All should return true")
	}

	if All(m, func(k string, v int) bool { return v > 1 }) {
		t.Error("All should return false")
	}
}

func TestNone(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	if !None(m, func(k string, v int) bool { return v > 10 }) {
		t.Error("None should return true")
	}

	if None(m, func(k string, v int) bool { return v > 2 }) {
		t.Error("None should return false")
	}
}

func TestCount(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	count := Count(m, func(k string, v int) bool { return v > 1 })
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestDiff(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"b": 2, "c": 4}
	diff := Diff(m1, m2)

	if len(diff) != 1 {
		t.Errorf("expected 1 element, got %d", len(diff))
	}
	if diff["a"] != 1 {
		t.Errorf("expected diff['a']=1, got %d", diff["a"])
	}
}

func TestDiff_Nil(t *testing.T) {
	var m1 map[string]int
	m2 := map[string]int{"a": 1}
	diff := Diff(m1, m2)
	if diff != nil {
		t.Errorf("expected nil, got %v", diff)
	}
}

func TestDiffValues(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 3}
	diff := DiffValues(m1, m2)

	if len(diff) != 1 {
		t.Errorf("expected 1 element, got %d", len(diff))
	}
	if diff["b"] != 2 {
		t.Errorf("expected diff['b']=2, got %d", diff["b"])
	}
}

func TestIntersection(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"b": 20, "c": 30, "d": 40}
	inter := Intersection(m1, m2)

	if len(inter) != 2 {
		t.Errorf("expected 2 elements, got %d", len(inter))
	}
	if inter["b"] != 2 {
		t.Errorf("expected inter['b']=2, got %d", inter["b"])
	}
	if inter["c"] != 3 {
		t.Errorf("expected inter['c']=3, got %d", inter["c"])
	}
}

func TestIntersection_Nil(t *testing.T) {
	var m1 map[string]int
	m2 := map[string]int{"a": 1}
	inter := Intersection(m1, m2)
	if inter != nil {
		t.Errorf("expected nil, got %v", inter)
	}
}

func TestIntersectionValues(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 3}
	inter := IntersectionValues(m1, m2)

	if len(inter) != 1 {
		t.Errorf("expected 1 element, got %d", len(inter))
	}
	if inter["a"] != 1 {
		t.Errorf("expected inter['a']=1, got %d", inter["a"])
	}
}

func TestTransform(t *testing.T) {
	m := map[int]string{1: "one", 2: "two"}
	result := Transform(m, func(k int, v string) (string, int) {
		return v, k
	})

	if len(result) != 2 {
		t.Errorf("expected 2 elements, got %d", len(result))
	}
	if result["one"] != 1 {
		t.Errorf("expected result['one']=1, got %d", result["one"])
	}
	if result["two"] != 2 {
		t.Errorf("expected result['two']=2, got %d", result["two"])
	}
}

func TestTransform_Nil(t *testing.T) {
	var m map[int]string
	result := Transform(m, func(k int, v string) (string, int) {
		return v, k
	})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestPop(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	v, ok := Pop(m, "a")

	if !ok {
		t.Error("expected ok to be true")
	}
	if v != 1 {
		t.Errorf("expected v=1, got %d", v)
	}
	if len(m) != 1 {
		t.Errorf("expected map length 1, got %d", len(m))
	}
	if _, ok := m["a"]; ok {
		t.Error("expected key 'a' to be deleted")
	}
}

func TestPop_NotFound(t *testing.T) {
	m := map[string]int{"a": 1}
	v, ok := Pop(m, "b")

	if ok {
		t.Error("expected ok to be false")
	}
	if v != 0 {
		t.Errorf("expected v=0, got %d", v)
	}
	if len(m) != 1 {
		t.Errorf("expected map length 1, got %d", len(m))
	}
}

func TestPopOr(t *testing.T) {
	m := map[string]int{"a": 1}
	v := PopOr(m, "a", 0)
	if v != 1 {
		t.Errorf("expected v=1, got %d", v)
	}

	v = PopOr(m, "b", 99)
	if v != 99 {
		t.Errorf("expected v=99, got %d", v)
	}
}

func TestUpdate(t *testing.T) {
	m := map[string]int{"count": 5}
	ok := Update(m, "count", func(v int) int { return v + 1 })

	if !ok {
		t.Error("expected ok to be true")
	}
	if m["count"] != 6 {
		t.Errorf("expected m['count']=6, got %d", m["count"])
	}
}

func TestUpdate_NotFound(t *testing.T) {
	m := map[string]int{"count": 5}
	ok := Update(m, "other", func(v int) int { return v + 1 })

	if ok {
		t.Error("expected ok to be false")
	}
	if len(m) != 1 {
		t.Errorf("expected map length 1, got %d", len(m))
	}
}

func TestUpdateOrInsert(t *testing.T) {
	m := map[string]int{}

	// 插入新值
	v := UpdateOrInsert(m, "count", 0, func(v int) int { return v + 1 })
	if v != 1 {
		t.Errorf("expected v=1, got %d", v)
	}
	if m["count"] != 1 {
		t.Errorf("expected m['count']=1, got %d", m["count"])
	}

	// 更新已有值
	v = UpdateOrInsert(m, "count", 0, func(v int) int { return v + 1 })
	if v != 2 {
		t.Errorf("expected v=2, got %d", v)
	}
	if m["count"] != 2 {
		t.Errorf("expected m['count']=2, got %d", m["count"])
	}
}

func TestSymmetricDiff(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 2, "c": 3}
	diff := SymmetricDiff(m1, m2)

	if len(diff) != 2 {
		t.Errorf("expected 2 elements, got %d", len(diff))
	}
	if diff["a"] != 1 {
		t.Errorf("expected diff['a']=1, got %d", diff["a"])
	}
	if diff["c"] != 3 {
		t.Errorf("expected diff['c']=3, got %d", diff["c"])
	}
}

func TestCollect(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
	byID := Collect(users, func(u User) (int, string) {
		return u.ID, u.Name
	})

	if len(byID) != 2 {
		t.Errorf("expected 2 elements, got %d", len(byID))
	}
	if byID[1] != "Alice" {
		t.Errorf("expected byID[1]='Alice', got %s", byID[1])
	}
	if byID[2] != "Bob" {
		t.Errorf("expected byID[2]='Bob', got %s", byID[2])
	}
}

func TestCollect_Nil(t *testing.T) {
	var users []struct{ ID int }
	result := Collect(users, func(u struct{ ID int }) (int, int) {
		return u.ID, u.ID
	})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestCollectByKey(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
	byID := CollectByKey(users, func(u User) int { return u.ID })

	if len(byID) != 2 {
		t.Errorf("expected 2 elements, got %d", len(byID))
	}
	if byID[1].Name != "Alice" {
		t.Errorf("expected byID[1].Name='Alice', got %s", byID[1].Name)
	}
}
