package cond

// If 根据条件返回两个值中的一个（三元表达式替代）
//
// 参数:
//   - condition: 条件表达式
//   - trueVal: 条件为 true 时返回的值
//   - falseVal: 条件为 false 时返回的值
//
// 返回:
//   - T: 根据条件选择的值
//
// 注意: 两个值都会被求值，如需延迟求值请使用 IfFunc
//
// 示例:
//
//	result := cond.If(age >= 18, "成年", "未成年")
//	max := cond.If(a > b, a, b)
func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// IfFunc 根据条件执行两个函数中的一个（延迟求值）
//
// 参数:
//   - condition: 条件表达式
//   - trueFn: 条件为 true 时执行的函数
//   - falseFn: 条件为 false 时执行的函数
//
// 返回:
//   - T: 执行对应函数的返回值
//
// 使用场景: 当值的计算代价较高时，使用此函数避免不必要的计算
//
// 示例:
//
//	result := cond.IfFunc(cached,
//	    func() Data { return cache.Get(key) },
//	    func() Data { return db.Query(key) },
//	)
func IfFunc[T any](condition bool, trueFn, falseFn func() T) T {
	if condition {
		return trueFn()
	}
	return falseFn()
}

// IfZero 如果值是零值则返回默认值
//
// 参数:
//   - value: 要检查的值
//   - defaultVal: 当 value 是零值时返回的默认值
//
// 返回:
//   - T: value 或 defaultVal
//
// 示例:
//
//	name := cond.IfZero(user.Name, "匿名")
//	port := cond.IfZero(config.Port, 8080)
func IfZero[T comparable](value, defaultVal T) T {
	var zero T
	if value == zero {
		return defaultVal
	}
	return value
}

// IfZeroFunc 如果值是零值则调用函数获取默认值
//
// 参数:
//   - value: 要检查的值
//   - defaultFn: 当 value 是零值时调用的函数
//
// 返回:
//   - T: value 或 defaultFn() 的返回值
//
// 示例:
//
//	id := cond.IfZeroFunc(user.ID, func() int64 { return generateID() })
func IfZeroFunc[T comparable](value T, defaultFn func() T) T {
	var zero T
	if value == zero {
		return defaultFn()
	}
	return value
}

// Coalesce 返回第一个非零值
//
// 参数:
//   - values: 可变参数，要检查的值列表
//
// 返回:
//   - T: 第一个非零值，如果全是零值则返回零值
//
// 示例:
//
//	name := cond.Coalesce(nickname, username, realName, "匿名")
//	// 返回第一个非空字符串
func Coalesce[T comparable](values ...T) T {
	var zero T
	for _, v := range values {
		if v != zero {
			return v
		}
	}
	return zero
}

// CoalesceFunc 返回第一个非零值，使用函数延迟求值
//
// 参数:
//   - fns: 可变参数，返回值的函数列表
//
// 返回:
//   - T: 第一个返回非零值的函数的结果
//
// 示例:
//
//	value := cond.CoalesceFunc(
//	    func() string { return cache.Get(key) },
//	    func() string { return db.Query(key) },
//	    func() string { return "default" },
//	)
func CoalesceFunc[T comparable](fns ...func() T) T {
	var zero T
	for _, fn := range fns {
		if v := fn(); v != zero {
			return v
		}
	}
	return zero
}

// IfNil 如果指针为 nil 则返回默认值
//
// 参数:
//   - ptr: 要检查的指针
//   - defaultVal: 当 ptr 为 nil 时返回的默认值
//
// 返回:
//   - T: *ptr 或 defaultVal
//
// 示例:
//
//	name := cond.IfNil(user.Nickname, "匿名")
func IfNil[T any](ptr *T, defaultVal T) T {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// IfNilFunc 如果指针为 nil 则调用函数获取默认值
//
// 参数:
//   - ptr: 要检查的指针
//   - defaultFn: 当 ptr 为 nil 时调用的函数
//
// 返回:
//   - T: *ptr 或 defaultFn() 的返回值
//
// 示例:
//
//	config := cond.IfNilFunc(customConfig, loadDefaultConfig)
func IfNilFunc[T any](ptr *T, defaultFn func() T) T {
	if ptr == nil {
		return defaultFn()
	}
	return *ptr
}

// Unless If 的反向版本，条件为 false 时返回第一个值
//
// 参数:
//   - condition: 条件表达式
//   - falseVal: 条件为 false 时返回的值
//   - trueVal: 条件为 true 时返回的值
//
// 返回:
//   - T: 根据条件选择的值
//
// 示例:
//
//	msg := cond.Unless(hasError, "操作成功", "操作失败")
func Unless[T any](condition bool, falseVal, trueVal T) T {
	if !condition {
		return falseVal
	}
	return trueVal
}
