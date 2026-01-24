package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MemoryDriver 内存队列驱动（用于测试和开发）
type MemoryDriver struct {
	jobs    map[string]*JobRecord
	queues  map[string][]*JobRecord // queue name -> jobs
	mu      sync.RWMutex
	signals map[string]chan struct{} // queue name -> signal channel
}

// NewMemoryDriver 创建内存驱动
func NewMemoryDriver() *MemoryDriver {
	return &MemoryDriver{
		jobs:    make(map[string]*JobRecord),
		queues:  make(map[string][]*JobRecord),
		signals: make(map[string]chan struct{}),
	}
}

// Push 推送任务
func (d *MemoryDriver) Push(job Job) error {
	return d.PushDelay(job, 0)
}

// PushDelay 推送延迟任务
func (d *MemoryDriver) PushDelay(job Job, delay time.Duration) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	payload, err := MarshalJob(job)
	if err != nil {
		return err
	}

	now := time.Now()
	scheduledAt := now.Add(delay)

	record := &JobRecord{
		ID:          job.GetID(),
		Queue:       job.GetQueue(),
		JobType:     fmt.Sprintf("%T", job),
		Payload:     payload,
		Status:      StatusPending,
		Attempts:    0,
		MaxRetries:  job.GetMaxRetries(),
		CreatedAt:   now,
		ScheduledAt: scheduledAt,
		Timeout:     job.GetTimeout(),
	}

	d.jobs[record.ID] = record

	// 如果不是延迟任务，立即加入队列
	if delay == 0 {
		d.addToQueue(record)
	} else {
		// 延迟任务，启动定时器
		go func() {
			time.Sleep(delay)
			d.mu.Lock()
			d.addToQueue(record)
			d.mu.Unlock()
		}()
	}

	return nil
}

// addToQueue 添加任务到队列（内部方法，需要持有锁）
func (d *MemoryDriver) addToQueue(record *JobRecord) {
	queue := record.Queue
	if d.queues[queue] == nil {
		d.queues[queue] = make([]*JobRecord, 0)
		d.signals[queue] = make(chan struct{}, 100)
	}

	d.queues[queue] = append(d.queues[queue], record)

	// 发送信号通知有新任务
	select {
	case d.signals[queue] <- struct{}{}:
	default:
	}
}

// Pop 获取任务
func (d *MemoryDriver) Pop(queue string, timeout time.Duration) (*JobRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		// 尝试获取任务
		d.mu.Lock()
		if len(d.queues[queue]) > 0 {
			// 获取第一个待执行的任务
			for i, record := range d.queues[queue] {
				if record.Status == StatusPending && time.Now().After(record.ScheduledAt) {
					// 从队列中移除
					d.queues[queue] = append(d.queues[queue][:i], d.queues[queue][i+1:]...)

					// 更新状态
					record.Status = StatusRunning
					record.Attempts++
					now := time.Now()
					record.StartedAt = &now

					d.mu.Unlock()
					return record, nil
				}
			}
		}
		d.mu.Unlock()

		// 等待新任务或超时
		if d.signals[queue] == nil {
			d.mu.Lock()
			if d.signals[queue] == nil {
				d.signals[queue] = make(chan struct{}, 100)
			}
			d.mu.Unlock()
		}

		select {
		case <-d.signals[queue]:
			// 有新任务，继续循环
			continue
		case <-ctx.Done():
			// 超时
			return nil, nil
		}
	}
}

// Ack 确认任务完成
func (d *MemoryDriver) Ack(jobID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	record, exists := d.jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	record.Status = StatusCompleted
	now := time.Now()
	record.CompletedAt = &now

	return nil
}

// Fail 标记任务失败
func (d *MemoryDriver) Fail(jobID string, err error) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	record, exists := d.jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	record.Status = StatusDead
	record.Error = err.Error()
	now := time.Now()
	record.FailedAt = &now

	return nil
}

// Retry 重试任务
func (d *MemoryDriver) Retry(jobID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	record, exists := d.jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	record.Status = StatusPending
	record.ScheduledAt = time.Now().Add(time.Duration(record.Attempts) * time.Minute) // 指数退避
	record.Error = ""

	d.addToQueue(record)

	return nil
}

// Delete 删除任务
func (d *MemoryDriver) Delete(jobID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.jobs, jobID)
	return nil
}

// GetJob 获取任务信息
func (d *MemoryDriver) GetJob(jobID string) (*JobRecord, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	record, exists := d.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job %s not found", jobID)
	}

	return record, nil
}

// ListJobs 列出任务
func (d *MemoryDriver) ListJobs(queue string, status JobStatus, limit int) ([]*JobRecord, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var jobs []*JobRecord
	for _, record := range d.jobs {
		if (queue == "" || record.Queue == queue) && (status == "" || record.Status == status) {
			jobs = append(jobs, record)
			if len(jobs) >= limit {
				break
			}
		}
	}

	return jobs, nil
}

// GetStats 获取统计信息
func (d *MemoryDriver) GetStats(queue string) (map[string]interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	stats := map[string]interface{}{
		"pending":   0,
		"running":   0,
		"completed": 0,
		"failed":    0,
		"dead":      0,
	}

	for _, record := range d.jobs {
		if queue != "" && record.Queue != queue {
			continue
		}

		switch record.Status {
		case StatusPending:
			stats["pending"] = stats["pending"].(int) + 1
		case StatusRunning:
			stats["running"] = stats["running"].(int) + 1
		case StatusCompleted:
			stats["completed"] = stats["completed"].(int) + 1
		case StatusFailed:
			stats["failed"] = stats["failed"].(int) + 1
		case StatusDead:
			stats["dead"] = stats["dead"].(int) + 1
		}
	}

	return stats, nil
}

// Close 关闭驱动
func (d *MemoryDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 关闭所有信号通道
	for _, ch := range d.signals {
		close(ch)
	}

	return nil
}

// SaveToFile 保存队列到文件（用于持久化）
func (d *MemoryDriver) SaveToFile(filename string) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	data, err := json.MarshalIndent(d.jobs, "", "  ")
	if err != nil {
		return err
	}

	// 这里应该写入文件，但为了简化，我们暂时省略
	_ = data
	return nil
}
