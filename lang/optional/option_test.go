package optional

import (
	"testing"
)

func TestSome(t *testing.T) {
	opt := Some(42)
	if !opt.IsSome() {
		t.Error("expected IsSome() to be true")
	}
	if opt.IsNone() {
		t.Error("expected IsNone() to be false")
	}
	if opt.Unwrap() != 42 {
		t.Errorf("expected Unwrap() to return 42, got %v", opt.Unwrap())
	}
}

func TestNone(t *testing.T) {
	opt := None[int]()
	if opt.IsSome() {
		t.Error("expected IsSome() to be false")
	}
	if !opt.IsNone() {
		t.Error("expected IsNone() to be true")
	}
	if opt.Unwrap() != 0 {
		t.Errorf("expected Unwrap() to return 0 for None, got %v", opt.Unwrap())
	}
}

func TestFromPtr(t *testing.T) {
	value := 42
	opt := FromPtr(&value)
	if !opt.IsSome() || opt.Unwrap() != 42 {
		t.Error("FromPtr should create Some for non-nil pointer")
	}

	var nilPtr *int
	opt = FromPtr(nilPtr)
	if opt.IsSome() {
		t.Error("FromPtr should create None for nil pointer")
	}
}

func TestFromValue(t *testing.T) {
	opt := FromValue(42, true)
	if !opt.IsSome() || opt.Unwrap() != 42 {
		t.Error("FromValue with ok=true should create Some")
	}

	opt = FromValue(42, false)
	if opt.IsSome() {
		t.Error("FromValue with ok=false should create None")
	}
}

func TestFromZero(t *testing.T) {
	opt := FromZero("")
	if opt.IsSome() {
		t.Error("FromZero should create None for empty string")
	}

	opt = FromZero("hello")
	if !opt.IsSome() || opt.Unwrap() != "hello" {
		t.Error("FromZero should create Some for non-empty string")
	}

	optInt := FromZero(0)
	if optInt.IsSome() {
		t.Error("FromZero should create None for 0")
	}

	optInt = FromZero(42)
	if !optInt.IsSome() || optInt.Unwrap() != 42 {
		t.Error("FromZero should create Some for non-zero int")
	}
}

func TestUnwrapOr(t *testing.T) {
	opt := Some(42)
	if opt.UnwrapOr(0) != 42 {
		t.Error("UnwrapOr should return value for Some")
	}

	opt = None[int]()
	if opt.UnwrapOr(99) != 99 {
		t.Error("UnwrapOr should return default for None")
	}
}

func TestUnwrapOrElse(t *testing.T) {
	opt := Some(42)
	called := false
	value := opt.UnwrapOrElse(func() int { called = true; return 99 })
	if value != 42 || called {
		t.Error("UnwrapOrElse should return value without calling fn for Some")
	}

	opt = None[int]()
	value = opt.UnwrapOrElse(func() int { return 99 })
	if value != 99 {
		t.Error("UnwrapOrElse should call fn for None")
	}
}

func TestExpect(t *testing.T) {
	opt := Some(42)
	if opt.Expect("should exist") != 42 {
		t.Error("Expect should return value for Some")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expect should panic for None")
		}
	}()
	None[int]().Expect("should panic")
}

func TestToPtr(t *testing.T) {
	opt := Some(42)
	ptr := opt.ToPtr()
	if ptr == nil || *ptr != 42 {
		t.Error("ToPtr should return pointer to value for Some")
	}

	opt = None[int]()
	if opt.ToPtr() != nil {
		t.Error("ToPtr should return nil for None")
	}
}

func TestFilter(t *testing.T) {
	opt := Some(42)

	filtered := opt.Filter(func(n int) bool { return n > 50 })
	if filtered.IsSome() {
		t.Error("Filter should return None when predicate is false")
	}

	filtered = opt.Filter(func(n int) bool { return n > 40 })
	if !filtered.IsSome() || filtered.Unwrap() != 42 {
		t.Error("Filter should return original Option when predicate is true")
	}

	none := None[int]()
	filtered = none.Filter(func(n int) bool { return true })
	if filtered.IsSome() {
		t.Error("Filter should return None for None input")
	}
}

func TestOr(t *testing.T) {
	some := Some(42)
	none := None[int]()
	other := Some(99)

	if some.Or(other).Unwrap() != 42 {
		t.Error("Or should return original for Some")
	}

	if none.Or(other).Unwrap() != 99 {
		t.Error("Or should return other for None")
	}
}

func TestOrElse(t *testing.T) {
	some := Some(42)
	none := None[int]()

	called := false
	result := some.OrElse(func() Option[int] { called = true; return Some(99) })
	if result.Unwrap() != 42 || called {
		t.Error("OrElse should return original without calling fn for Some")
	}

	result = none.OrElse(func() Option[int] { return Some(99) })
	if result.Unwrap() != 99 {
		t.Error("OrElse should call fn for None")
	}
}

func TestAnd(t *testing.T) {
	some := Some(42)
	none := None[int]()
	other := Some("hello")

	result := And(some, other)
	if !result.IsSome() || result.Unwrap() != "hello" {
		t.Error("And should return other for Some")
	}

	result2 := And(none, other)
	if result2.IsSome() {
		t.Error("And should return None for None input")
	}
}

func TestMap(t *testing.T) {
	opt := Some(42)
	result := Map(opt, func(n int) string { return "number" })
	if !result.IsSome() || result.Unwrap() != "number" {
		t.Error("Map should transform value for Some")
	}

	none := None[int]()
	result = Map(none, func(n int) string { return "number" })
	if result.IsSome() {
		t.Error("Map should return None for None input")
	}
}

func TestMapOr(t *testing.T) {
	opt := Some(42)
	result := MapOr(opt, "default", func(n int) string { return "number" })
	if result != "number" {
		t.Error("MapOr should transform value for Some")
	}

	none := None[int]()
	result = MapOr(none, "default", func(n int) string { return "number" })
	if result != "default" {
		t.Error("MapOr should return default for None")
	}
}

func TestFlatMap(t *testing.T) {
	opt := Some(42)
	result := FlatMap(opt, func(n int) Option[string] {
		if n > 0 {
			return Some("positive")
		}
		return None[string]()
	})
	if !result.IsSome() || result.Unwrap() != "positive" {
		t.Error("FlatMap should chain Options for Some")
	}

	none := None[int]()
	result = FlatMap(none, func(n int) Option[string] { return Some("never") })
	if result.IsSome() {
		t.Error("FlatMap should return None for None input")
	}
}

func TestFlatten(t *testing.T) {
	nested := Some(Some(42))
	flat := Flatten(nested)
	if !flat.IsSome() || flat.Unwrap() != 42 {
		t.Error("Flatten should unwrap nested Some")
	}

	nestedNone := Some(None[int]())
	flat = Flatten(nestedNone)
	if flat.IsSome() {
		t.Error("Flatten should return None for nested None")
	}

	outer := None[Option[int]]()
	flat = Flatten(outer)
	if flat.IsSome() {
		t.Error("Flatten should return None for outer None")
	}
}

func TestZip(t *testing.T) {
	opt1 := Some(42)
	opt2 := Some("hello")

	result := Zip(opt1, opt2)
	if !result.IsSome() {
		t.Error("Zip should return Some when both are Some")
	}
	zipped := result.Unwrap()
	if zipped.First != 42 || zipped.Second != "hello" {
		t.Error("Zip should contain both values")
	}

	none := None[int]()
	result2 := Zip(none, opt2)
	if result2.IsSome() {
		t.Error("Zip should return None when first is None")
	}

	result3 := Zip(opt1, None[string]())
	if result3.IsSome() {
		t.Error("Zip should return None when second is None")
	}
}

func TestZipWith(t *testing.T) {
	opt1 := Some(10)
	opt2 := Some(5)

	result := ZipWith(opt1, opt2, func(a, b int) int { return a + b })
	if !result.IsSome() || result.Unwrap() != 15 {
		t.Error("ZipWith should combine values")
	}

	none := None[int]()
	result = ZipWith(none, opt2, func(a, b int) int { return a + b })
	if result.IsSome() {
		t.Error("ZipWith should return None when first is None")
	}
}

func TestContains(t *testing.T) {
	opt := Some(42)
	if !opt.Contains(42) {
		t.Error("Contains should return true for equal value")
	}
	if opt.Contains(43) {
		t.Error("Contains should return false for different value")
	}

	none := None[int]()
	if none.Contains(42) {
		t.Error("Contains should return false for None")
	}
}

func TestString(t *testing.T) {
	some := Some(42)
	if some.String() != "Some(?)" {
		t.Errorf("unexpected String output: %s", some.String())
	}

	none := None[int]()
	if none.String() != "None" {
		t.Errorf("unexpected String output: %s", none.String())
	}
}

func TestUnwrapOrZero(t *testing.T) {
	opt := Some(42)
	if opt.UnwrapOrZero() != 42 {
		t.Error("UnwrapOrZero should return value for Some")
	}

	none := None[int]()
	if none.UnwrapOrZero() != 0 {
		t.Error("UnwrapOrZero should return zero value for None")
	}
}
