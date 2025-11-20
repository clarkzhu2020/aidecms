package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/event"
)

// EventTest 测试事件系统
func EventTest(args []string) {
	fmt.Println("Testing Event System...")
	fmt.Println()

	// 创建事件分发器
	dispatcher := event.NewDispatcher(5)

	// 注册同步监听器
	fmt.Println("1. Registering synchronous listeners...")
	dispatcher.Listen("user.registered", func(ctx context.Context, e event.Event) error {
		if userEvent, ok := e.(*event.UserRegistered); ok {
			fmt.Printf("   [Sync] Welcome email sent to %s\n", userEvent.Email)
		}
		return nil
	})

	dispatcher.Listen("user.registered", func(ctx context.Context, e event.Event) error {
		if userEvent, ok := e.(*event.UserRegistered); ok {
			fmt.Printf("   [Sync] User profile created for %s\n", userEvent.Username)
		}
		return nil
	})

	// 注册异步监听器
	fmt.Println("\n2. Registering asynchronous listeners...")
	dispatcher.ListenAsync("post.published", func(ctx context.Context, e event.Event) error {
		if postEvent, ok := e.(*event.PostPublished); ok {
			time.Sleep(500 * time.Millisecond) // 模拟耗时操作
			fmt.Printf("   [Async] Sitemap updated for post: %s\n", postEvent.Title)
		}
		return nil
	})

	dispatcher.ListenAsync("post.published", func(ctx context.Context, e event.Event) error {
		if postEvent, ok := e.(*event.PostPublished); ok {
			time.Sleep(300 * time.Millisecond) // 模拟耗时操作
			fmt.Printf("   [Async] Cache cleared for post: %s\n", postEvent.Title)
		}
		return nil
	})

	// 注册带优先级的监听器
	fmt.Println("\n3. Registering priority listeners...")
	dispatcher.ListenWithPriority("order.created", func(ctx context.Context, e event.Event) error {
		fmt.Println("   [Priority 1] Validating order...")
		return nil
	}, 1)

	dispatcher.ListenWithPriority("order.created", func(ctx context.Context, e event.Event) error {
		fmt.Println("   [Priority 2] Calculating shipping...")
		return nil
	}, 2)

	dispatcher.ListenWithPriority("order.created", func(ctx context.Context, e event.Event) error {
		fmt.Println("   [Priority 3] Sending confirmation email...")
		return nil
	}, 3)

	// 测试事件分发
	fmt.Println("\n4. Dispatching events...\n")

	// 用户注册事件
	fmt.Println("Dispatching: UserRegistered")
	userEvent := event.NewUserRegistered(1, "john_doe", "john@example.com")
	dispatcher.Dispatch(userEvent)

	time.Sleep(100 * time.Millisecond)

	// 文章发布事件
	fmt.Println("\nDispatching: PostPublished")
	postEvent := event.NewPostPublished(1, "My First Post", "my-first-post")
	dispatcher.Dispatch(postEvent)

	time.Sleep(100 * time.Millisecond)

	// 订单创建事件
	fmt.Println("\nDispatching: OrderCreated (with priority)")
	orderEvent := event.NewOrderCreated("ORDER-001", 1, 99.99)
	dispatcher.Dispatch(orderEvent)

	// 等待异步监听器完成
	fmt.Println("\nWaiting for async listeners to complete...")
	time.Sleep(1 * time.Second)

	// 显示统计信息
	fmt.Println("\n5. Event Statistics:")
	stats := dispatcher.GetStats()
	fmt.Printf("   Total Events: %d\n", stats["total_events"])
	fmt.Printf("   Total Listeners: %d\n", stats["total_listeners"])
	fmt.Printf("   Total Executions: %d\n", stats["total_executions"])
	fmt.Printf("   Success Rate: %.2f%%\n", stats["success_rate"])
	fmt.Printf("   Queue Size: %d\n", stats["queue_size"])
	fmt.Printf("   Workers: %d\n", stats["workers"])

	// 显示最近的日志
	fmt.Println("\n6. Recent Event Logs:")
	logs := dispatcher.GetLogs("", 10)
	for i, log := range logs {
		status := "✓"
		if !log.Success {
			status = "✗"
		}
		asyncTag := ""
		if log.Async {
			asyncTag = " (async)"
		}
		fmt.Printf("   %d. %s %s -> %s%s [%v]\n",
			i+1, status, log.EventName, log.ListenerName, asyncTag, log.Duration)
	}

	// 停止分发器
	dispatcher.Stop()
	fmt.Println("\n✓ Event system test completed")
}

// EventList 列出所有注册的事件
func EventList(args []string) {
	dispatcher := event.GetDispatcher()

	events := dispatcher.GetAllEvents()
	if len(events) == 0 {
		fmt.Println("No events registered")
		return
	}

	fmt.Printf("Registered Events (%d):\n", len(events))
	for i, eventName := range events {
		listeners := dispatcher.GetListeners(eventName)
		fmt.Printf("%d. %s (%d listeners)\n", i+1, eventName, len(listeners))
		for j, listener := range listeners {
			asyncTag := ""
			if listener.Async {
				asyncTag = " [async]"
			}
			priorityTag := ""
			if listener.Priority != 0 {
				priorityTag = fmt.Sprintf(" [priority: %d]", listener.Priority)
			}
			fmt.Printf("   %d.%d %s%s%s\n", i+1, j+1, listener.Name, asyncTag, priorityTag)
		}
	}
}

// EventStats 显示事件统计信息
func EventStats(args []string) {
	stats := event.GetStats()

	fmt.Println("Event System Statistics:")
	fmt.Printf("  Total Events:     %d\n", stats["total_events"])
	fmt.Printf("  Total Listeners:  %d\n", stats["total_listeners"])
	fmt.Printf("  Total Executions: %d\n", stats["total_executions"])
	fmt.Printf("  Success Count:    %d\n", stats["success_count"])
	fmt.Printf("  Fail Count:       %d\n", stats["fail_count"])
	fmt.Printf("  Success Rate:     %.2f%%\n", stats["success_rate"])
	fmt.Printf("  Queue Size:       %d\n", stats["queue_size"])
	fmt.Printf("  Workers:          %d\n", stats["workers"])
}
