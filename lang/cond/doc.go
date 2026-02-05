// Package cond 提供条件工具函数，简化条件判断和选择逻辑
//
// 主要功能:
//   - If/IfFunc: 三元表达式替代
//   - IfZero: 零值判断与默认值
//   - Coalesce: 返回第一个非零值
//   - Switch: 类型安全的 switch 表达式
//
// 示例:
//
//	// 三元表达式替代
//	result := cond.If(age >= 18, "成年", "未成年")
//
//	// 延迟求值（避免不必要的计算）
//	result := cond.IfFunc(cached, getCached, fetchFromDB)
//
//	// 返回第一个非零值
//	name := cond.Coalesce(nickname, username, "匿名")
//
//	// Switch 表达式
//	msg := cond.Switch[string, string](status).
//	    Case("active", "活跃").
//	    Case("inactive", "非活跃").
//	    Default("未知")
package cond
