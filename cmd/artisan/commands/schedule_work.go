package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chenyusolar/aidecms/pkg/schedule"
)

// ScheduleWork 启动调度器工作进程
func ScheduleWork(args []string) {
	fmt.Println("Starting scheduler...")

	// 创建调度器
	scheduler := schedule.NewScheduler()

	// 注册示例任务
	registerTasks(scheduler)

	// 启动调度器
	scheduler.Start()
	fmt.Println("✓ Scheduler started successfully")
	fmt.Printf("✓ %d tasks registered\n", len(scheduler.ListTasks()))
	fmt.Println("\nPress Ctrl+C to stop...")

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// 停止调度器
	fmt.Println("\nStopping scheduler...")
	scheduler.Stop()
	fmt.Println("✓ Scheduler stopped")
}

// registerTasks 注册示例任务
func registerTasks(scheduler *schedule.Scheduler) {
	// 示例1: 每分钟执行
	scheduler.NewTask("test-every-minute").
		EveryMinute().
		Description("Test task that runs every minute").
		Do(func() error {
			fmt.Println("[Task] Every minute task executed")
			return nil
		})

	// 示例2: 每5分钟执行
	scheduler.NewTask("cleanup-cache").
		EveryFiveMinutes().
		Description("Clean up expired cache").
		Do(func() error {
			fmt.Println("[Task] Cleanup cache task executed")
			// 这里可以调用缓存清理逻辑
			return nil
		})

	// 示例3: 每天凌晨2点执行
	scheduler.NewTask("backup-database").
		DailyAt(2, 0).
		Description("Backup database daily at 2:00 AM").
		Do(func() error {
			fmt.Println("[Task] Database backup task executed")
			// 这里可以调用数据库备份逻辑
			return nil
		})

	// 示例4: 每周日凌晨执行
	scheduler.NewTask("weekly-report").
		Weekly().
		Description("Generate weekly report").
		Do(func() error {
			fmt.Println("[Task] Weekly report task executed")
			// 这里可以调用报告生成逻辑
			return nil
		})

	// 示例5: 自定义 Cron 表达式
	scheduler.NewTask("custom-task").
		Cron("0 */2 * * *"). // 每2小时执行一次
		Description("Custom cron task").
		Do(func() error {
			fmt.Println("[Task] Custom cron task executed")
			return nil
		})
}
