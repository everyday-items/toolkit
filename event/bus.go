// Package event 提供轻量级事件总线
//
// 支持发布-订阅模式的事件分发，用于系统组件间的松耦合通信。
// 线程安全，支持按类型订阅和全局订阅。
//
// 使用示例:
//
//	bus := event.New()
//	defer bus.Close()
//
//	unsub := bus.Subscribe("agent.start", func(e event.Event) {
//	    fmt.Println("Agent 启动:", e.Payload)
//	})
//	defer unsub()
//
//	bus.Publish(event.Event{Type: "agent.start", Payload: "my-agent"})
package event

import (
	"sync"
	"sync/atomic"
	"time"
)

// 预定义事件类型常量
const (
	// Agent 生命周期事件
	EventAgentStart = "agent.start"
	EventAgentEnd   = "agent.end"
	EventAgentError = "agent.error"

	// 工具调用事件
	EventToolCall   = "tool.call"
	EventToolResult = "tool.result"

	// LLM 调用事件
	EventLLMRequest  = "llm.request"
	EventLLMResponse = "llm.response"
	EventLLMStream   = "llm.stream"

	// Skill 生命周期事件
	EventSkillLoad   = "skill.load"
	EventSkillUnload = "skill.unload"

	// 成本事件
	EventCostUpdate = "cost.update"

	// 安全事件
	EventSecurityAlert = "security.alert"
)

// Event 事件结构
type Event struct {
	// Type 事件类型（如 "agent.start"）
	Type string
	// Payload 事件数据（任意类型）
	Payload any
	// Timestamp 事件发生时间
	Timestamp time.Time
	// Source 事件来源（如 Agent ID）
	Source string
	// ID 事件唯一标识
	ID string
}

// Handler 事件处理函数
type Handler func(Event)

// subscription 订阅记录
type subscription struct {
	id      uint64
	handler Handler
}

// Bus 事件总线
//
// 线程安全的发布-订阅事件分发器。
// 支持按类型订阅和全局订阅（订阅所有事件）。
type Bus struct {
	// mu 保护 subscribers 和 globalSubs
	mu sync.RWMutex
	// subscribers 按事件类型索引的订阅者
	subscribers map[string][]subscription
	// globalSubs 全局订阅者（接收所有事件）
	globalSubs []subscription
	// nextID 递增的订阅 ID
	nextID atomic.Uint64
	// closed 总线是否已关闭
	closed atomic.Bool
}

// New 创建事件总线
func New() *Bus {
	return &Bus{
		subscribers: make(map[string][]subscription),
	}
}

// Subscribe 订阅指定类型的事件
//
// 返回取消订阅函数。调用取消函数后，该处理器不再接收事件。
func (b *Bus) Subscribe(eventType string, handler Handler) (unsubscribe func()) {
	if b.closed.Load() {
		return func() {}
	}

	id := b.nextID.Add(1)
	sub := subscription{id: id, handler: handler}

	b.mu.Lock()
	b.subscribers[eventType] = append(b.subscribers[eventType], sub)
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		subs := b.subscribers[eventType]
		for i, s := range subs {
			if s.id == id {
				b.subscribers[eventType] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}
}

// SubscribeAll 订阅所有事件
//
// 返回取消订阅函数。
func (b *Bus) SubscribeAll(handler Handler) (unsubscribe func()) {
	if b.closed.Load() {
		return func() {}
	}

	id := b.nextID.Add(1)
	sub := subscription{id: id, handler: handler}

	b.mu.Lock()
	b.globalSubs = append(b.globalSubs, sub)
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		for i, s := range b.globalSubs {
			if s.id == id {
				b.globalSubs = append(b.globalSubs[:i], b.globalSubs[i+1:]...)
				break
			}
		}
	}
}

// Publish 异步发布事件
//
// 每个订阅者在独立的 goroutine 中接收事件，
// 不会阻塞发布者。事件处理器中的 panic 会被捕获并忽略。
func (b *Bus) Publish(event Event) {
	if b.closed.Load() {
		return
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	b.mu.RLock()
	// 复制订阅者列表，避免持锁执行 handler
	typeSubs := make([]subscription, len(b.subscribers[event.Type]))
	copy(typeSubs, b.subscribers[event.Type])
	globalSubs := make([]subscription, len(b.globalSubs))
	copy(globalSubs, b.globalSubs)
	b.mu.RUnlock()

	for _, sub := range typeSubs {
		go safeCall(sub.handler, event)
	}
	for _, sub := range globalSubs {
		go safeCall(sub.handler, event)
	}
}

// PublishSync 同步发布事件
//
// 在当前 goroutine 中依次调用所有订阅者，
// 阻塞直到所有处理器执行完毕。
func (b *Bus) PublishSync(event Event) {
	if b.closed.Load() {
		return
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	b.mu.RLock()
	typeSubs := make([]subscription, len(b.subscribers[event.Type]))
	copy(typeSubs, b.subscribers[event.Type])
	globalSubs := make([]subscription, len(b.globalSubs))
	copy(globalSubs, b.globalSubs)
	b.mu.RUnlock()

	for _, sub := range typeSubs {
		safeCall(sub.handler, event)
	}
	for _, sub := range globalSubs {
		safeCall(sub.handler, event)
	}
}

// Close 关闭事件总线
//
// 关闭后不再接受新的订阅和发布。
func (b *Bus) Close() {
	b.closed.Store(true)
	b.mu.Lock()
	b.subscribers = make(map[string][]subscription)
	b.globalSubs = nil
	b.mu.Unlock()
}

// Len 返回当前订阅总数
func (b *Bus) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	count := len(b.globalSubs)
	for _, subs := range b.subscribers {
		count += len(subs)
	}
	return count
}

// safeCall 安全调用 handler，捕获 panic
func safeCall(handler Handler, event Event) {
	defer func() {
		recover() //nolint:errcheck // 忽略 handler panic
	}()
	handler(event)
}
