package event

import (
	"context"
	"testing"
	"time"
)

func TestDispatcher(t *testing.T) {
	dispatcher := NewDispatcher(2)
	defer dispatcher.Stop()

	executed := false

	// 注册监听器
	dispatcher.Listen("test.event", func(ctx context.Context, event Event) error {
		executed = true
		return nil
	})

	// 分发事件
	testEvent := &BaseEvent{Name: "test.event"}
	err := dispatcher.Dispatch(testEvent)
	if err != nil {
		t.Fatalf("Dispatch error: %v", err)
	}

	if !executed {
		t.Error("Listener was not executed")
	}
}

func TestAsyncListener(t *testing.T) {
	dispatcher := NewDispatcher(2)
	defer dispatcher.Stop()

	done := make(chan bool, 1)

	// 注册异步监听器
	dispatcher.ListenAsync("test.async", func(ctx context.Context, event Event) error {
		time.Sleep(100 * time.Millisecond)
		done <- true
		return nil
	})

	// 分发事件
	testEvent := &BaseEvent{Name: "test.async"}
	dispatcher.Dispatch(testEvent)

	// 等待异步执行完成
	select {
	case <-done:
		// 成功
	case <-time.After(1 * time.Second):
		t.Error("Async listener did not execute in time")
	}
}

func TestListenerPriority(t *testing.T) {
	dispatcher := NewDispatcher(2)
	defer dispatcher.Stop()

	var order []int

	// 注册不同优先级的监听器
	dispatcher.ListenWithPriority("test.priority", func(ctx context.Context, event Event) error {
		order = append(order, 2)
		return nil
	}, 2)

	dispatcher.ListenWithPriority("test.priority", func(ctx context.Context, event Event) error {
		order = append(order, 1)
		return nil
	}, 1)

	dispatcher.ListenWithPriority("test.priority", func(ctx context.Context, event Event) error {
		order = append(order, 3)
		return nil
	}, 3)

	// 分发事件
	testEvent := &BaseEvent{Name: "test.priority"}
	dispatcher.Dispatch(testEvent)

	// 验证执行顺序
	if len(order) != 3 {
		t.Fatalf("Expected 3 listeners to execute, got %d", len(order))
	}

	if order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Errorf("Listeners executed in wrong order: %v", order)
	}
}

func TestForget(t *testing.T) {
	dispatcher := NewDispatcher(2)
	defer dispatcher.Stop()

	executed := false

	// 注册监听器
	dispatcher.ListenWithOptions("test.forget", "test-listener", func(ctx context.Context, event Event) error {
		executed = true
		return nil
	}, 0, false)

	// 移除监听器
	dispatcher.Forget("test.forget", "test-listener")

	// 分发事件
	testEvent := &BaseEvent{Name: "test.forget"}
	dispatcher.Dispatch(testEvent)

	if executed {
		t.Error("Forgotten listener was executed")
	}
}

func TestMultipleListeners(t *testing.T) {
	dispatcher := NewDispatcher(2)
	defer dispatcher.Stop()

	count := 0

	// 注册多个监听器
	for i := 0; i < 5; i++ {
		dispatcher.Listen("test.multiple", func(ctx context.Context, event Event) error {
			count++
			return nil
		})
	}

	// 分发事件
	testEvent := &BaseEvent{Name: "test.multiple"}
	dispatcher.Dispatch(testEvent)

	if count != 5 {
		t.Errorf("Expected 5 listeners to execute, got %d", count)
	}
}

func TestGetStats(t *testing.T) {
	dispatcher := NewDispatcher(2)
	defer dispatcher.Stop()

	// 注册监听器
	dispatcher.Listen("test.stats", func(ctx context.Context, event Event) error {
		return nil
	})

	// 分发事件
	testEvent := &BaseEvent{Name: "test.stats"}
	dispatcher.Dispatch(testEvent)

	// 获取统计
	stats := dispatcher.GetStats()

	if stats["total_events"].(int) == 0 {
		t.Error("Expected at least one event")
	}

	if stats["total_listeners"].(int) == 0 {
		t.Error("Expected at least one listener")
	}

	if stats["total_executions"].(int) == 0 {
		t.Error("Expected at least one execution")
	}
}

func TestGlobalDispatcher(t *testing.T) {
	executed := false

	// 使用全局分发器
	Listen("test.global", func(ctx context.Context, event Event) error {
		executed = true
		return nil
	})

	// 分发事件
	testEvent := &BaseEvent{Name: "test.global"}
	err := Dispatch(testEvent)
	if err != nil {
		t.Fatalf("Dispatch error: %v", err)
	}

	// 给异步监听器一些时间执行
	time.Sleep(100 * time.Millisecond)

	if !executed {
		t.Error("Global listener was not executed")
	}

	// 清理
	ForgetAll("test.global")
}

func TestUserRegisteredEvent(t *testing.T) {
	dispatcher := NewDispatcher(2)
	defer dispatcher.Stop()

	var capturedEvent *UserRegistered

	// 注册监听器
	dispatcher.Listen("user.registered", func(ctx context.Context, event Event) error {
		if userEvent, ok := event.(*UserRegistered); ok {
			capturedEvent = userEvent
		}
		return nil
	})

	// 创建并分发事件
	event := NewUserRegistered(1, "testuser", "test@example.com")
	dispatcher.Dispatch(event)

	if capturedEvent == nil {
		t.Fatal("Event was not captured")
	}

	if capturedEvent.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", capturedEvent.UserID)
	}

	if capturedEvent.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", capturedEvent.Username)
	}
}
