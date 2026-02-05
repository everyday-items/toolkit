package reflectx

import (
	"reflect"
)

// DeepCopy 深度拷贝值
//
// 参数:
//   - src: 源值
//
// 返回:
//   - T: 拷贝后的值
//
// 注意: 支持基本类型、结构体、切片、map、指针
// 对于不支持的类型（如 chan、func）返回零值
//
// 示例:
//
//	type User struct { Name string }
//	user := User{Name: "Alice"}
//	copied := reflectx.DeepCopy(user)  // 独立副本
func DeepCopy[T any](src T) T {
	return deepCopyValue(reflect.ValueOf(src)).Interface().(T)
}

// deepCopyValue 递归深拷贝 reflect.Value
func deepCopyValue(src reflect.Value) reflect.Value {
	if !src.IsValid() {
		return src
	}

	switch src.Kind() {
	case reflect.Ptr:
		return deepCopyPtr(src)
	case reflect.Interface:
		return deepCopyInterface(src)
	case reflect.Struct:
		return deepCopyStruct(src)
	case reflect.Slice:
		return deepCopySlice(src)
	case reflect.Map:
		return deepCopyMap(src)
	case reflect.Array:
		return deepCopyArray(src)
	default:
		// 基本类型直接复制
		dst := reflect.New(src.Type()).Elem()
		dst.Set(src)
		return dst
	}
}

// deepCopyPtr 深拷贝指针
func deepCopyPtr(src reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}
	dst := reflect.New(src.Type().Elem())
	dst.Elem().Set(deepCopyValue(src.Elem()))
	return dst
}

// deepCopyInterface 深拷贝接口
func deepCopyInterface(src reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}
	return deepCopyValue(src.Elem())
}

// deepCopyStruct 深拷贝结构体
func deepCopyStruct(src reflect.Value) reflect.Value {
	dst := reflect.New(src.Type()).Elem()
	for i := range src.NumField() {
		srcField := src.Field(i)
		dstField := dst.Field(i)
		if dstField.CanSet() {
			dstField.Set(deepCopyValue(srcField))
		}
	}
	return dst
}

// deepCopySlice 深拷贝切片
func deepCopySlice(src reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}
	dst := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
	for i := range src.Len() {
		dst.Index(i).Set(deepCopyValue(src.Index(i)))
	}
	return dst
}

// deepCopyMap 深拷贝 map
func deepCopyMap(src reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}
	dst := reflect.MakeMap(src.Type())
	for _, key := range src.MapKeys() {
		dst.SetMapIndex(deepCopyValue(key), deepCopyValue(src.MapIndex(key)))
	}
	return dst
}

// deepCopyArray 深拷贝数组
func deepCopyArray(src reflect.Value) reflect.Value {
	dst := reflect.New(src.Type()).Elem()
	for i := range src.Len() {
		dst.Index(i).Set(deepCopyValue(src.Index(i)))
	}
	return dst
}

// Clone 浅拷贝值（仅拷贝顶层）
//
// 参数:
//   - src: 源值
//
// 返回:
//   - T: 拷贝后的值
//
// 注意: 对于指针、切片、map 等引用类型，仅拷贝引用
func Clone[T any](src T) T {
	dst := reflect.New(reflect.TypeOf(src)).Elem()
	dst.Set(reflect.ValueOf(src))
	return dst.Interface().(T)
}
