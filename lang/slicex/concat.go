package slicex

// Concat 连接多个切片
//
// 参数:
//   - slices: 要连接的切片列表
//
// 返回:
//   - []T: 连接后的新切片
//
// 示例:
//
//	result := slicex.Concat([]int{1, 2}, []int{3, 4}, []int{5})
//	// [1, 2, 3, 4, 5]
func Concat[T any](slices ...[]T) []T {
	totalLen := 0
	for _, s := range slices {
		totalLen += len(s)
	}
	if totalLen == 0 {
		return nil
	}
	result := make([]T, 0, totalLen)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// Flatten 将二维切片展平为一维切片
//
// 参数:
//   - slices: 二维切片
//
// 返回:
//   - []T: 展平后的一维切片
//
// 示例:
//
//	result := slicex.Flatten([][]int{{1, 2}, {3, 4}, {5}})
//	// [1, 2, 3, 4, 5]
func Flatten[T any](slices [][]T) []T {
	return Concat(slices...)
}

// Prepend 在切片前面添加元素
//
// 参数:
//   - slice: 原切片
//   - items: 要添加的元素
//
// 返回:
//   - []T: 新切片
//
// 示例:
//
//	result := slicex.Prepend([]int{3, 4}, 1, 2)
//	// [1, 2, 3, 4]
func Prepend[T any](slice []T, items ...T) []T {
	if len(items) == 0 {
		return slice
	}
	result := make([]T, 0, len(slice)+len(items))
	result = append(result, items...)
	result = append(result, slice...)
	return result
}

// Append 在切片后面添加元素（与内置 append 相似但返回新切片）
//
// 参数:
//   - slice: 原切片
//   - items: 要添加的元素
//
// 返回:
//   - []T: 新切片
//
// 示例:
//
//	result := slicex.Append([]int{1, 2}, 3, 4)
//	// [1, 2, 3, 4]
func Append[T any](slice []T, items ...T) []T {
	if len(items) == 0 {
		return slice
	}
	result := make([]T, 0, len(slice)+len(items))
	result = append(result, slice...)
	result = append(result, items...)
	return result
}
