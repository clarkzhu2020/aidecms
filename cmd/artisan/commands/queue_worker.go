package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	q "github.com/clarkzhu2020/aidecms/pkg/queue"
)

func QueueWork(args []string) {
	fmt.Println("Starting queue workers...")

	// 创建队列管理器（使用内存驱动，实际应用中应该使用 Redis）
	driver := q.NewMemoryDriver()
	queueMgr := q.NewQueue(driver)

	// 配置工作进程
	queueMgr.SetWorkers(3).SetQueues([]string{"default", "high", "low"})

	// 注册任务处理器
	registerJobHandlers(queueMgr)

	// 启动工作进程
	go func() {
		if err := queueMgr.Work(); err != nil {
			fmt.Printf("Queue worker error: %v\n", err)
		}
	}()

	fmt.Println("✓ Queue workers started")
	fmt.Println("  - Workers: 3")
	fmt.Println("  - Queues: default, high, low")
	fmt.Println("\nPress Ctrl+C to stop...")

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nStopping queue workers...")
	queueMgr.Stop()
	fmt.Println("✓ Queue workers stopped")
}

// registerJobHandlers 注册任务处理器
func registerJobHandlers(queueMgr *q.Queue) {
	// 示例：邮件发送任务
	queueMgr.Register("EmailSendJob", func(payload []byte) error {
		fmt.Printf("[Queue] Processing email job: %s\n", string(payload))
		time.Sleep(1 * time.Second) // 模拟邮件发送
		return nil
	})

	// 示例：数据处理任务
	queueMgr.Register("DataProcessJob", func(payload []byte) error {
		fmt.Printf("[Queue] Processing data job: %s\n", string(payload))
		time.Sleep(2 * time.Second) // 模拟数据处理
		return nil
	})

	// 示例：图片处理任务
	queueMgr.Register("ImageProcessJob", func(payload []byte) error {
		fmt.Printf("[Queue] Processing image job: %s\n", string(payload))
		time.Sleep(3 * time.Second) // 模拟图片处理
		return nil
	})

	fmt.Println("✓ Registered 3 job handlers")
}
