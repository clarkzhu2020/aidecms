package schedule

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task 表示一个调度任务
type Task struct {
	ID          string
	Name        string
	Schedule    string // Cron 表达式或预定义调度
	Handler     func() error
	LastRunAt   time.Time
	NextRunAt   time.Time
	RunCount    int
	FailCount   int
	IsRunning   bool
	Description string
	cronExpr    *CronExpression
	mu          sync.RWMutex
}

// Scheduler 任务调度器
type Scheduler struct {
	tasks      map[string]*Task
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	ticker     *time.Ticker
	isRunning  bool
	runningMu  sync.RWMutex
	logs       []TaskLog
	logsMu     sync.RWMutex
	maxLogSize int
}

// TaskLog 任务执行日志
type TaskLog struct {
	TaskID    string
	TaskName  string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Success   bool
	Error     string
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		tasks:      make(map[string]*Task),
		ctx:        ctx,
		cancel:     cancel,
		logs:       make([]TaskLog, 0),
		maxLogSize: 1000, // 最多保留 1000 条日志
	}
}

// AddTask 添加任务
func (s *Scheduler) AddTask(task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.ID == "" {
		task.ID = fmt.Sprintf("task_%d", time.Now().UnixNano())
	}

	// 解析 Cron 表达式
	if task.Schedule != "" {
		cronExpr, err := ParseCron(task.Schedule)
		if err != nil {
			return fmt.Errorf("invalid cron expression: %w", err)
		}
		task.cronExpr = cronExpr
		task.NextRunAt = cronExpr.Next(time.Now())
	}

	s.tasks[task.ID] = task
	return nil
}

// RemoveTask 移除任务
func (s *Scheduler) RemoveTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	delete(s.tasks, taskID)
	return nil
}

// GetTask 获取任务
func (s *Scheduler) GetTask(taskID string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	return task, nil
}

// ListTasks 列出所有任务
func (s *Scheduler) ListTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.runningMu.Lock()
	if s.isRunning {
		s.runningMu.Unlock()
		return
	}
	s.isRunning = true
	s.runningMu.Unlock()

	s.ticker = time.NewTicker(time.Second)
	go s.run()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()

	if !s.isRunning {
		return
	}

	s.cancel()
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.isRunning = false
}

// IsRunning 检查调度器是否运行中
func (s *Scheduler) IsRunning() bool {
	s.runningMu.RLock()
	defer s.runningMu.RUnlock()
	return s.isRunning
}

// run 调度器主循环
func (s *Scheduler) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case now := <-s.ticker.C:
			s.checkAndRunTasks(now)
		}
	}
}

// checkAndRunTasks 检查并运行到期任务
func (s *Scheduler) checkAndRunTasks(now time.Time) {
	s.mu.RLock()
	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	s.mu.RUnlock()

	for _, task := range tasks {
		if s.shouldRun(task, now) {
			go s.runTask(task)
		}
	}
}

// shouldRun 检查任务是否应该运行
func (s *Scheduler) shouldRun(task *Task, now time.Time) bool {
	task.mu.RLock()
	defer task.mu.RUnlock()

	if task.IsRunning {
		return false
	}

	if task.cronExpr == nil {
		return false
	}

	// 检查是否到达下次运行时间（精确到分钟）
	return now.Unix() >= task.NextRunAt.Unix()
}

// runTask 运行任务
func (s *Scheduler) runTask(task *Task) {
	task.mu.Lock()
	if task.IsRunning {
		task.mu.Unlock()
		return
	}
	task.IsRunning = true
	task.mu.Unlock()

	log := TaskLog{
		TaskID:    task.ID,
		TaskName:  task.Name,
		StartTime: time.Now(),
	}

	// 运行任务
	err := task.Handler()

	log.EndTime = time.Now()
	log.Duration = log.EndTime.Sub(log.StartTime)

	task.mu.Lock()
	task.IsRunning = false
	task.LastRunAt = log.StartTime
	task.RunCount++

	if err != nil {
		task.FailCount++
		log.Success = false
		log.Error = err.Error()
	} else {
		log.Success = true
	}

	// 计算下次运行时间
	if task.cronExpr != nil {
		task.NextRunAt = task.cronExpr.Next(time.Now())
	}
	task.mu.Unlock()

	// 保存日志
	s.addLog(log)
}

// addLog 添加日志
func (s *Scheduler) addLog(log TaskLog) {
	s.logsMu.Lock()
	defer s.logsMu.Unlock()

	s.logs = append(s.logs, log)

	// 限制日志大小
	if len(s.logs) > s.maxLogSize {
		s.logs = s.logs[len(s.logs)-s.maxLogSize:]
	}
}

// GetLogs 获取日志
func (s *Scheduler) GetLogs(taskID string, limit int) []TaskLog {
	s.logsMu.RLock()
	defer s.logsMu.RUnlock()

	var logs []TaskLog
	count := 0

	// 从后往前遍历（最新的日志在前）
	for i := len(s.logs) - 1; i >= 0 && count < limit; i-- {
		if taskID == "" || s.logs[i].TaskID == taskID {
			logs = append(logs, s.logs[i])
			count++
		}
	}

	return logs
}

// RunNow 立即运行任务
func (s *Scheduler) RunNow(taskID string) error {
	task, err := s.GetTask(taskID)
	if err != nil {
		return err
	}

	go s.runTask(task)
	return nil
}

// GetStats 获取统计信息
func (s *Scheduler) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalTasks := len(s.tasks)
	runningTasks := 0
	totalRuns := 0
	totalFails := 0

	for _, task := range s.tasks {
		task.mu.RLock()
		if task.IsRunning {
			runningTasks++
		}
		totalRuns += task.RunCount
		totalFails += task.FailCount
		task.mu.RUnlock()
	}

	return map[string]interface{}{
		"total_tasks":   totalTasks,
		"running_tasks": runningTasks,
		"total_runs":    totalRuns,
		"total_fails":   totalFails,
		"success_rate":  calculateSuccessRate(totalRuns, totalFails),
	}
}

func calculateSuccessRate(total, fails int) float64 {
	if total == 0 {
		return 0
	}
	return float64(total-fails) / float64(total) * 100
}
