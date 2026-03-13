package event

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBus_SubscribeAndPublish(t *testing.T) {
	bus := New()
	defer bus.Close()

	var received atomic.Bool
	bus.Subscribe(EventAgentStart, func(e Event) {
		received.Store(true)
	})

	bus.Publish(Event{Type: EventAgentStart, Payload: "test"})
	time.Sleep(50 * time.Millisecond)

	if !received.Load() {
		t.Error("应该收到事件")
	}
}

func TestBus_PublishSync(t *testing.T) {
	bus := New()
	defer bus.Close()

	var received bool
	bus.Subscribe(EventToolCall, func(e Event) {
		received = true
	})

	bus.PublishSync(Event{Type: EventToolCall})

	if !received {
		t.Error("同步发布后应该立即收到")
	}
}

func TestBus_SubscribeAll(t *testing.T) {
	bus := New()
	defer bus.Close()

	var count atomic.Int32
	bus.SubscribeAll(func(e Event) {
		count.Add(1)
	})

	bus.Publish(Event{Type: EventAgentStart})
	bus.Publish(Event{Type: EventToolCall})
	time.Sleep(50 * time.Millisecond)

	if count.Load() != 2 {
		t.Errorf("全局订阅应收到 2 个事件, 实际 %d", count.Load())
	}
}

func TestBus_Unsubscribe(t *testing.T) {
	bus := New()
	defer bus.Close()

	var count atomic.Int32
	unsub := bus.Subscribe(EventAgentStart, func(e Event) {
		count.Add(1)
	})

	bus.PublishSync(Event{Type: EventAgentStart})
	unsub()
	bus.PublishSync(Event{Type: EventAgentStart})

	if count.Load() != 1 {
		t.Errorf("取消订阅后不应收到事件, count=%d", count.Load())
	}
}

func TestBus_Close(t *testing.T) {
	bus := New()
	bus.Subscribe(EventAgentStart, func(e Event) {})
	bus.Close()

	if bus.Len() != 0 {
		t.Error("关闭后订阅数应为 0")
	}

	// 关闭后订阅应返回空函数
	unsub := bus.Subscribe(EventAgentStart, func(e Event) {})
	unsub() // 不应 panic
}

func TestBus_Len(t *testing.T) {
	bus := New()
	defer bus.Close()

	if bus.Len() != 0 {
		t.Error("初始订阅数应为 0")
	}

	unsub1 := bus.Subscribe(EventAgentStart, func(e Event) {})
	bus.Subscribe(EventToolCall, func(e Event) {})
	bus.SubscribeAll(func(e Event) {})

	if bus.Len() != 3 {
		t.Errorf("期望 3 个订阅, 实际 %d", bus.Len())
	}

	unsub1()
	if bus.Len() != 2 {
		t.Errorf("取消后期望 2 个订阅, 实际 %d", bus.Len())
	}
}

func TestBus_Concurrent(t *testing.T) {
	bus := New()
	defer bus.Close()

	var count atomic.Int64
	for i := 0; i < 10; i++ {
		bus.Subscribe(EventAgentStart, func(e Event) {
			count.Add(1)
		})
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bus.Publish(Event{Type: EventAgentStart})
		}()
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	if count.Load() != 1000 {
		t.Errorf("期望 1000 次调用, 实际 %d", count.Load())
	}
}

func TestBus_PanicRecovery(t *testing.T) {
	bus := New()
	defer bus.Close()

	bus.Subscribe(EventAgentStart, func(e Event) {
		panic("handler panic")
	})

	// 不应 panic
	bus.PublishSync(Event{Type: EventAgentStart})
}

func TestBus_AutoTimestamp(t *testing.T) {
	bus := New()
	defer bus.Close()

	var received Event
	bus.Subscribe(EventAgentStart, func(e Event) {
		received = e
	})

	bus.PublishSync(Event{Type: EventAgentStart})

	if received.Timestamp.IsZero() {
		t.Error("未设置 Timestamp 时应自动填充")
	}
}
