package syncx

import (
	"sync"
)

// Once 泛型版的 sync.Once，可以返回值
//
// 与 sync.Once 不同，Once[T] 可以返回初始化的值
type Once[T any] struct {
	once  sync.Once
	value T
}

// Do 执行初始化函数（只执行一次）
//
// 参数:
//   - fn: 初始化函数
//
// 返回:
//   - T: 初始化的值
//
// 示例:
//
//	var o syncx.Once[*Config]
//	config := o.Do(func() *Config {
//	    return loadConfig()
//	})
func (o *Once[T]) Do(fn func() T) T {
	o.once.Do(func() {
		o.value = fn()
	})
	return o.value
}

// Value 返回已初始化的值（如果尚未初始化则返回零值）
//
// 返回:
//   - T: 值
//
// 注意: 建议先调用 Do 确保值已初始化
func (o *Once[T]) Value() T {
	return o.value
}

// OnceValue 创建一个只执行一次的函数
//
// 参数:
//   - fn: 要执行的函数
//
// 返回:
//   - func() T: 包装后的函数，多次调用返回相同结果
//
// 示例:
//
//	getConfig := syncx.OnceValue(func() *Config {
//	    return loadConfig()
//	})
//	config1 := getConfig()  // 执行 loadConfig
//	config2 := getConfig()  // 返回缓存结果
func OnceValue[T any](fn func() T) func() T {
	var o Once[T]
	return func() T {
		return o.Do(fn)
	}
}

// OnceValueErr 创建一个只执行一次的函数（可能返回错误）
//
// 参数:
//   - fn: 要执行的函数
//
// 返回:
//   - func() (T, error): 包装后的函数
//
// 示例:
//
//	getDB := syncx.OnceValueErr(func() (*DB, error) {
//	    return connectDB()
//	})
//	db1, err := getDB()  // 执行 connectDB
//	db2, err := getDB()  // 返回缓存结果
func OnceValueErr[T any](fn func() (T, error)) func() (T, error) {
	var once sync.Once
	var value T
	var err error
	return func() (T, error) {
		once.Do(func() {
			value, err = fn()
		})
		return value, err
	}
}

// OnceFunc 创建一个只执行一次的函数（无返回值）
//
// 参数:
//   - fn: 要执行的函数
//
// 返回:
//   - func(): 包装后的函数
//
// 示例:
//
//	initOnce := syncx.OnceFunc(func() {
//	    initialize()
//	})
//	initOnce()  // 执行 initialize
//	initOnce()  // 不执行
func OnceFunc(fn func()) func() {
	var once sync.Once
	return func() {
		once.Do(fn)
	}
}

// OnceErr 泛型版的 sync.Once，可以返回值和错误
type OnceErr[T any] struct {
	once  sync.Once
	value T
	err   error
}

// Do 执行初始化函数（只执行一次）
//
// 参数:
//   - fn: 初始化函数
//
// 返回:
//   - T: 初始化的值
//   - error: 初始化错误
//
// 示例:
//
//	var o syncx.OnceErr[*DB]
//	db, err := o.Do(func() (*DB, error) {
//	    return connectDB()
//	})
func (o *OnceErr[T]) Do(fn func() (T, error)) (T, error) {
	o.once.Do(func() {
		o.value, o.err = fn()
	})
	return o.value, o.err
}

// Value 返回已初始化的值
//
// 返回:
//   - T: 值
//   - error: 错误
func (o *OnceErr[T]) Value() (T, error) {
	return o.value, o.err
}

// IsInitialized 检查是否已初始化
//
// 返回:
//   - bool: 如果已初始化返回 true
//
// 注意: 这不是原子操作，仅供参考
func (o *OnceErr[T]) IsInitialized() bool {
	// 尝试执行空操作来检查状态
	initialized := true
	o.once.Do(func() {
		initialized = false
	})
	return initialized
}
