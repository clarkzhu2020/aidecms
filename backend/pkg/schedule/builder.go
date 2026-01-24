package schedule

import (
	"fmt"
)

// TaskBuilder 任务构建器
type TaskBuilder struct {
	task      *Task
	scheduler *Scheduler
}

// NewTask 创建新任务
func (s *Scheduler) NewTask(name string) *TaskBuilder {
	return &TaskBuilder{
		task: &Task{
			Name: name,
		},
		scheduler: s,
	}
}

// Cron 设置 Cron 表达式
func (tb *TaskBuilder) Cron(expr string) *TaskBuilder {
	tb.task.Schedule = expr
	return tb
}

// EveryMinute 每分钟执行
func (tb *TaskBuilder) EveryMinute() *TaskBuilder {
	tb.task.Schedule = "* * * * *"
	return tb
}

// EveryFiveMinutes 每5分钟执行
func (tb *TaskBuilder) EveryFiveMinutes() *TaskBuilder {
	tb.task.Schedule = "*/5 * * * *"
	return tb
}

// EveryTenMinutes 每10分钟执行
func (tb *TaskBuilder) EveryTenMinutes() *TaskBuilder {
	tb.task.Schedule = "*/10 * * * *"
	return tb
}

// EveryFifteenMinutes 每15分钟执行
func (tb *TaskBuilder) EveryFifteenMinutes() *TaskBuilder {
	tb.task.Schedule = "*/15 * * * *"
	return tb
}

// EveryThirtyMinutes 每30分钟执行
func (tb *TaskBuilder) EveryThirtyMinutes() *TaskBuilder {
	tb.task.Schedule = "*/30 * * * *"
	return tb
}

// Hourly 每小时执行
func (tb *TaskBuilder) Hourly() *TaskBuilder {
	tb.task.Schedule = "0 * * * *"
	return tb
}

// HourlyAt 每小时在指定分钟执行
func (tb *TaskBuilder) HourlyAt(minute int) *TaskBuilder {
	tb.task.Schedule = fmt.Sprintf("%d * * * *", minute)
	return tb
}

// Daily 每天执行
func (tb *TaskBuilder) Daily() *TaskBuilder {
	tb.task.Schedule = "0 0 * * *"
	return tb
}

// DailyAt 每天在指定时间执行
func (tb *TaskBuilder) DailyAt(hour, minute int) *TaskBuilder {
	tb.task.Schedule = fmt.Sprintf("%d %d * * *", minute, hour)
	return tb
}

// Weekly 每周执行
func (tb *TaskBuilder) Weekly() *TaskBuilder {
	tb.task.Schedule = "0 0 * * 0" // 周日
	return tb
}

// WeeklyOn 每周在指定日期和时间执行
// weekday: 0=Sunday, 1=Monday, ..., 6=Saturday
func (tb *TaskBuilder) WeeklyOn(weekday, hour, minute int) *TaskBuilder {
	tb.task.Schedule = fmt.Sprintf("%d %d * * %d", minute, hour, weekday)
	return tb
}

// Monthly 每月执行
func (tb *TaskBuilder) Monthly() *TaskBuilder {
	tb.task.Schedule = "0 0 1 * *" // 每月1号
	return tb
}

// MonthlyOn 每月在指定日期和时间执行
func (tb *TaskBuilder) MonthlyOn(day, hour, minute int) *TaskBuilder {
	tb.task.Schedule = fmt.Sprintf("%d %d %d * *", minute, hour, day)
	return tb
}

// Yearly 每年执行
func (tb *TaskBuilder) Yearly() *TaskBuilder {
	tb.task.Schedule = "0 0 1 1 *" // 每年1月1日
	return tb
}

// Description 设置描述
func (tb *TaskBuilder) Description(desc string) *TaskBuilder {
	tb.task.Description = desc
	return tb
}

// Do 设置处理函数并注册任务
func (tb *TaskBuilder) Do(handler func() error) error {
	tb.task.Handler = handler
	return tb.scheduler.AddTask(tb.task)
}

// Weekdays 工作日执行
func (tb *TaskBuilder) Weekdays() *TaskBuilder {
	tb.task.Schedule = "0 0 * * 1-5" // 周一到周五
	return tb
}

// Weekends 周末执行
func (tb *TaskBuilder) Weekends() *TaskBuilder {
	tb.task.Schedule = "0 0 * * 0,6" // 周日和周六
	return tb
}

// At 在指定时间执行（需要先设置 Daily/Weekly 等）
func (tb *TaskBuilder) At(timeStr string) *TaskBuilder {
	// 解析时间字符串，格式: "HH:MM"
	var hour, minute int
	fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)

	// 更新 cron 表达式的时间部分
	// 这是一个简化版本，实际使用时需要更完善的实现
	tb.task.Schedule = fmt.Sprintf("%d %d * * *", minute, hour)
	return tb
}
