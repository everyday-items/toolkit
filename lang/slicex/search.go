package slicex

// LastIndexOf 查找元素在切片中最后出现的索引
//
// 参数:
//   - slice: 要查找的切片
//   - item: 要查找的元素
//
// 返回:
//   - int: 最后出现的索引，未找到返回 -1
//
// 示例:
//
//	index := slicex.LastIndexOf([]int{1, 2, 3, 2, 1}, 2)  // 3
func LastIndexOf[T comparable](slice []T, item T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// LastIndexOfFunc 使用函数查找最后一个满足条件的元素索引
//
// 参数:
//   - slice: 要查找的切片
//   - fn: 匹配函数
//
// 返回:
//   - int: 最后满足条件的索引，未找到返回 -1
//
// 示例:
//
//	index := slicex.LastIndexOfFunc([]int{1, 2, 3, 4}, func(n int) bool {
//	    return n%2 == 0
//	})  // 3 (元素 4 的索引)
func LastIndexOfFunc[T any](slice []T, fn func(T) bool) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if fn(slice[i]) {
			return i
		}
	}
	return -1
}

// FindLastIndex 与 LastIndexOfFunc 相同（别名）
func FindLastIndex[T any](slice []T, fn func(T) bool) int {
	return LastIndexOfFunc(slice, fn)
}

// First 返回切片的第一个元素
//
// 参数:
//   - slice: 输入切片
//
// 返回:
//   - T: 第一个元素
//   - bool: 是否存在
//
// 示例:
//
//	first, ok := slicex.First([]int{1, 2, 3})  // 1, true
//	first, ok := slicex.First([]int{})        // 0, false
func First[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[0], true
}

// Last 返回切片的最后一个元素
//
// 参数:
//   - slice: 输入切片
//
// 返回:
//   - T: 最后一个元素
//   - bool: 是否存在
//
// 示例:
//
//	last, ok := slicex.Last([]int{1, 2, 3})  // 3, true
//	last, ok := slicex.Last([]int{})        // 0, false
func Last[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[len(slice)-1], true
}

// FirstOr 返回切片的第一个元素，空切片返回默认值
//
// 参数:
//   - slice: 输入切片
//   - defaultVal: 默认值
//
// 返回:
//   - T: 第一个元素或默认值
//
// 示例:
//
//	first := slicex.FirstOr([]int{1, 2, 3}, 0)  // 1
//	first := slicex.FirstOr([]int{}, 0)        // 0
func FirstOr[T any](slice []T, defaultVal T) T {
	if len(slice) == 0 {
		return defaultVal
	}
	return slice[0]
}

// LastOr 返回切片的最后一个元素，空切片返回默认值
//
// 参数:
//   - slice: 输入切片
//   - defaultVal: 默认值
//
// 返回:
//   - T: 最后一个元素或默认值
//
// 示例:
//
//	last := slicex.LastOr([]int{1, 2, 3}, 0)  // 3
//	last := slicex.LastOr([]int{}, 0)        // 0
func LastOr[T any](slice []T, defaultVal T) T {
	if len(slice) == 0 {
		return defaultVal
	}
	return slice[len(slice)-1]
}

// Nth 返回切片的第 n 个元素（支持负索引）
//
// 参数:
//   - slice: 输入切片
//   - n: 索引（负数表示从末尾开始，-1 为最后一个）
//
// 返回:
//   - T: 第 n 个元素
//   - bool: 是否存在
//
// 示例:
//
//	v, ok := slicex.Nth([]int{1, 2, 3}, 1)   // 2, true
//	v, ok := slicex.Nth([]int{1, 2, 3}, -1)  // 3, true
//	v, ok := slicex.Nth([]int{1, 2, 3}, 10)  // 0, false
func Nth[T any](slice []T, n int) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	if n < 0 {
		n = len(slice) + n
	}
	if n < 0 || n >= len(slice) {
		var zero T
		return zero, false
	}
	return slice[n], true
}

// NthOr 返回切片的第 n 个元素，不存在返回默认值
//
// 参数:
//   - slice: 输入切片
//   - n: 索引（支持负索引）
//   - defaultVal: 默认值
//
// 返回:
//   - T: 第 n 个元素或默认值
func NthOr[T any](slice []T, n int, defaultVal T) T {
	v, ok := Nth(slice, n)
	if !ok {
		return defaultVal
	}
	return v
}
