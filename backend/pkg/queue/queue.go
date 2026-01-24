package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Job 队列任务接口
type Job interface {
	// Handle 执行任务
	Handle() error
	// GetID 获取任务 ID
	GetID() string
	// GetQueue 获取队列名称
	GetQueue() string
	// GetMaxRetries 获取最大重试次数
	GetMaxRetries() int
	// GetTimeout 获取超时时间
	GetTimeout() time.Duration
}

// BaseJob 基础任务结构
type BaseJob struct {
	ID          string        `json:"id"`
	Queue       string        `json:"queue"`
	MaxRetries  int           `json:"max_retries"`
	Timeout     time.Duration `json:"timeout"`
	Payload     interface{}   `json:"payload"`
	CreatedAt   time.Time     `json:"created_at"`
	ScheduledAt time.Time     `json:"scheduled_at"` // 延迟任务的执行时间
}

// GetID 实现 Job 接口
func (j *BaseJob) GetID() string {
	if j.ID == "" {
		j.ID = fmt.Sprintf("job_%d", time.Now().UnixNano())
	}
	return j.ID
}

// GetQueue 实现 Job 接口
func (j *BaseJob) GetQueue() string {
	if j.Queue == "" {
		return "default"
	}
	return j.Queue
}

// GetMaxRetries 实现 Job 接口
func (j *BaseJob) GetMaxRetries() int {
	if j.MaxRetries == 0 {
		return 3 // 默认重试 3 次
	}
	return j.MaxRetries
}

// GetTimeout 实现 Job 接口
func (j *BaseJob) GetTimeout() time.Duration {
	if j.Timeout == 0 {
		return 5 * time.Minute // 默认 5 分钟超时
	}
	return j.Timeout
}

// JobStatus 任务状态
type JobStatus string

const (
	StatusPending   JobStatus = "pending"   // 等待执行
	StatusRunning   JobStatus = "running"   // 执行中
	StatusCompleted JobStatus = "completed" // 已完成
	StatusFailed    JobStatus = "failed"    // 失败
	StatusRetrying  JobStatus = "retrying"  // 重试中
	StatusDead      JobStatus = "dead"      // 死信（超过最大重试次数）
)

// JobRecord 任务记录
type JobRecord struct {
	ID          string        `json:"id"`
	Queue       string        `json:"queue"`
	JobType     string        `json:"job_type"`
	Payload     string        `json:"payload"` // JSON 编码的任务数据
	Status      JobStatus     `json:"status"`
	Attempts    int           `json:"attempts"`
	MaxRetries  int           `json:"max_retries"`
	CreatedAt   time.Time     `json:"created_at"`
	ScheduledAt time.Time     `json:"scheduled_at"` // 延迟任务
	StartedAt   *time.Time    `json:"started_at,omitempty"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	FailedAt    *time.Time    `json:"failed_at,omitempty"`
	Error       string        `json:"error,omitempty"`
	Timeout     time.Duration `json:"timeout"`
}

// Driver 队列驱动接口
type Driver interface {
	// Push 推送任务到队列
	Push(job Job) error
	// PushDelay 推送延迟任务
	PushDelay(job Job, delay time.Duration) error
	// Pop 从队列获取任务
	Pop(queue string, timeout time.Duration) (*JobRecord, error)
	// Ack 确认任务完成
	Ack(jobID string) error
	// Fail 标记任务失败
	Fail(jobID string, err error) error
	// Retry 重试任务
	Retry(jobID string) error
	// Delete 删除任务
	Delete(jobID string) error
	// GetJob 获取任务信息
	GetJob(jobID string) (*JobRecord, error)
	// ListJobs 列出任务
	ListJobs(queue string, status JobStatus, limit int) ([]*JobRecord, error)
	// GetStats 获取统计信息
	GetStats(queue string) (map[string]interface{}, error)
	// Close 关闭驱动
	Close() error
}

// Queue 队列管理器
type Queue struct {
	driver       Driver
	handlers     map[string]JobHandler
	ctx          context.Context
	cancel       context.CancelFunc
	workers      int
	workerQueues []string
}

// JobHandler 任务处理函数
type JobHandler func(payload []byte) error

// NewQueue 创建新的队列管理器
func NewQueue(driver Driver) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	return &Queue{
		driver:       driver,
		handlers:     make(map[string]JobHandler),
		ctx:          ctx,
		cancel:       cancel,
		workers:      1,
		workerQueues: []string{"default"},
	}
}

// SetWorkers 设置工作进程数量
func (q *Queue) SetWorkers(workers int) *Queue {
	q.workers = workers
	return q
}

// SetQueues 设置监听的队列
func (q *Queue) SetQueues(queues []string) *Queue {
	q.workerQueues = queues
	return q
}

// Register 注册任务处理器
func (q *Queue) Register(jobType string, handler JobHandler) {
	q.handlers[jobType] = handler
}

// Push 推送任务
func (q *Queue) Push(job Job) error {
	return q.driver.Push(job)
}

// PushDelay 推送延迟任务
func (q *Queue) PushDelay(job Job, delay time.Duration) error {
	return q.driver.PushDelay(job, delay)
}

// Work 启动队列工作进程
func (q *Queue) Work() error {
	fmt.Printf("Starting %d workers for queues: %v\n", q.workers, q.workerQueues)

	// 启动多个工作进程
	for i := 0; i < q.workers; i++ {
		go q.worker(i)
	}

	// 等待取消信号
	<-q.ctx.Done()
	return nil
}

// Stop 停止队列工作进程
func (q *Queue) Stop() {
	q.cancel()
	if q.driver != nil {
		q.driver.Close()
	}
}

// worker 工作进程
func (q *Queue) worker(id int) {
	fmt.Printf("Worker #%d started\n", id)

	for {
		select {
		case <-q.ctx.Done():
			fmt.Printf("Worker #%d stopped\n", id)
			return
		default:
			// 轮询所有队列
			for _, queueName := range q.workerQueues {
				q.processQueue(queueName)
			}
			time.Sleep(1 * time.Second) // 避免空轮询消耗 CPU
		}
	}
}

// processQueue 处理队列中的任务
func (q *Queue) processQueue(queueName string) {
	// 获取任务（阻塞等待 5 秒）
	jobRecord, err := q.driver.Pop(queueName, 5*time.Second)
	if err != nil || jobRecord == nil {
		return
	}

	// 查找处理器
	handler, exists := q.handlers[jobRecord.JobType]
	if !exists {
		q.driver.Fail(jobRecord.ID, fmt.Errorf("no handler for job type: %s", jobRecord.JobType))
		return
	}

	// 执行任务
	err = q.executeJob(jobRecord, handler)
	if err != nil {
		// 任务失败，检查是否需要重试
		if jobRecord.Attempts < jobRecord.MaxRetries {
			q.driver.Retry(jobRecord.ID)
		} else {
			// 超过最大重试次数，进入死信队列
			q.driver.Fail(jobRecord.ID, err)
		}
		return
	}

	// 任务成功，确认完成
	q.driver.Ack(jobRecord.ID)
}

// executeJob 执行任务
func (q *Queue) executeJob(jobRecord *JobRecord, handler JobHandler) error {
	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), jobRecord.Timeout)
	defer cancel()

	// 在 goroutine 中执行任务
	errChan := make(chan error, 1)
	go func() {
		errChan <- handler([]byte(jobRecord.Payload))
	}()

	// 等待完成或超时
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("job timeout after %v", jobRecord.Timeout)
	}
}

// GetStats 获取队列统计
func (q *Queue) GetStats(queue string) (map[string]interface{}, error) {
	return q.driver.GetStats(queue)
}

// RetryFailed 重试所有失败的任务
func (q *Queue) RetryFailed(queue string) error {
	jobs, err := q.driver.ListJobs(queue, StatusFailed, 1000)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err := q.driver.Retry(job.ID); err != nil {
			fmt.Printf("Failed to retry job %s: %v\n", job.ID, err)
		}
	}

	return nil
}

// PurgeQueue 清空队列
func (q *Queue) PurgeQueue(queue string, status JobStatus) error {
	jobs, err := q.driver.ListJobs(queue, status, 10000)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err := q.driver.Delete(job.ID); err != nil {
			fmt.Printf("Failed to delete job %s: %v\n", job.ID, err)
		}
	}

	return nil
}

// MarshalJob 序列化任务
func MarshalJob(job interface{}) (string, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalJob 反序列化任务
func UnmarshalJob(data string, job interface{}) error {
	return json.Unmarshal([]byte(data), job)
}
