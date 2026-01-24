package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisDriver Redis 队列驱动
type RedisDriver struct {
	client *redis.Client
	prefix string
	ctx    context.Context
}

// NewRedisDriver 创建 Redis 驱动
func NewRedisDriver(client *redis.Client, prefix string) *RedisDriver {
	if prefix == "" {
		prefix = "queue"
	}
	return &RedisDriver{
		client: client,
		prefix: prefix,
		ctx:    context.Background(),
	}
}

// Push 推送任务
func (d *RedisDriver) Push(job Job) error {
	return d.PushDelay(job, 0)
}

// PushDelay 推送延迟任务
func (d *RedisDriver) PushDelay(job Job, delay time.Duration) error {
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

	// 保存任务详情
	recordData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	jobKey := d.jobKey(record.ID)
	if err := d.client.Set(d.ctx, jobKey, recordData, 7*24*time.Hour).Err(); err != nil {
		return err
	}

	// 添加到队列或延迟队列
	if delay == 0 {
		// 立即执行的任务，加入列表
		queueKey := d.queueKey(record.Queue)
		return d.client.LPush(d.ctx, queueKey, record.ID).Err()
	} else {
		// 延迟任务，加入有序集合（使用执行时间作为分数）
		delayedKey := d.delayedKey(record.Queue)
		return d.client.ZAdd(d.ctx, delayedKey, redis.Z{
			Score:  float64(scheduledAt.Unix()),
			Member: record.ID,
		}).Err()
	}
}

// Pop 获取任务
func (d *RedisDriver) Pop(queue string, timeout time.Duration) (*JobRecord, error) {
	// 首先检查延迟队列，将到期的任务移到主队列
	d.moveDelayedJobs(queue)

	// 从主队列获取任务（阻塞）
	queueKey := d.queueKey(queue)
	result, err := d.client.BRPop(d.ctx, timeout, queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 超时，没有任务
		}
		return nil, err
	}

	if len(result) < 2 {
		return nil, nil
	}

	jobID := result[1]

	// 获取任务详情
	jobKey := d.jobKey(jobID)
	data, err := d.client.Get(d.ctx, jobKey).Result()
	if err != nil {
		return nil, err
	}

	var record JobRecord
	if err := json.Unmarshal([]byte(data), &record); err != nil {
		return nil, err
	}

	// 更新状态
	record.Status = StatusRunning
	record.Attempts++
	now := time.Now()
	record.StartedAt = &now

	// 保存更新后的状态
	recordData, _ := json.Marshal(record)
	d.client.Set(d.ctx, jobKey, recordData, 7*24*time.Hour)

	// 添加到处理中队列（用于追踪）
	processingKey := d.processingKey(queue)
	d.client.ZAdd(d.ctx, processingKey, redis.Z{
		Score:  float64(now.Unix()),
		Member: jobID,
	})

	return &record, nil
}

// moveDelayedJobs 将到期的延迟任务移到主队列
func (d *RedisDriver) moveDelayedJobs(queue string) {
	delayedKey := d.delayedKey(queue)
	now := float64(time.Now().Unix())

	// 获取所有到期的任务
	jobIDs, err := d.client.ZRangeByScore(d.ctx, delayedKey, &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%f", now),
	}).Result()

	if err != nil || len(jobIDs) == 0 {
		return
	}

	// 移动到主队列
	queueKey := d.queueKey(queue)
	pipe := d.client.Pipeline()
	for _, jobID := range jobIDs {
		pipe.LPush(d.ctx, queueKey, jobID)
		pipe.ZRem(d.ctx, delayedKey, jobID)
	}
	pipe.Exec(d.ctx)
}

// Ack 确认任务完成
func (d *RedisDriver) Ack(jobID string) error {
	jobKey := d.jobKey(jobID)

	// 获取任务
	data, err := d.client.Get(d.ctx, jobKey).Result()
	if err != nil {
		return err
	}

	var record JobRecord
	if err := json.Unmarshal([]byte(data), &record); err != nil {
		return err
	}

	// 更新状态
	record.Status = StatusCompleted
	now := time.Now()
	record.CompletedAt = &now

	// 保存
	recordData, _ := json.Marshal(record)
	d.client.Set(d.ctx, jobKey, recordData, 24*time.Hour) // 完成的任务保留 24 小时

	// 从处理中队列移除
	processingKey := d.processingKey(record.Queue)
	d.client.ZRem(d.ctx, processingKey, jobID)

	return nil
}

// Fail 标记任务失败
func (d *RedisDriver) Fail(jobID string, err error) error {
	jobKey := d.jobKey(jobID)

	// 获取任务
	data, err2 := d.client.Get(d.ctx, jobKey).Result()
	if err2 != nil {
		return err2
	}

	var record JobRecord
	if err2 := json.Unmarshal([]byte(data), &record); err2 != nil {
		return err2
	}

	// 更新状态
	record.Status = StatusDead
	record.Error = err.Error()
	now := time.Now()
	record.FailedAt = &now

	// 保存
	recordData, _ := json.Marshal(record)
	d.client.Set(d.ctx, jobKey, recordData, 7*24*time.Hour) // 失败的任务保留 7 天

	// 添加到死信队列
	deadKey := d.deadKey()
	d.client.LPush(d.ctx, deadKey, jobID)

	// 从处理中队列移除
	processingKey := d.processingKey(record.Queue)
	d.client.ZRem(d.ctx, processingKey, jobID)

	return nil
}

// Retry 重试任务
func (d *RedisDriver) Retry(jobID string) error {
	jobKey := d.jobKey(jobID)

	// 获取任务
	data, err := d.client.Get(d.ctx, jobKey).Result()
	if err != nil {
		return err
	}

	var record JobRecord
	if err := json.Unmarshal([]byte(data), &record); err != nil {
		return err
	}

	// 更新状态和延迟时间（指数退避）
	record.Status = StatusPending
	backoffDelay := time.Duration(record.Attempts) * time.Minute
	record.ScheduledAt = time.Now().Add(backoffDelay)
	record.Error = ""

	// 保存
	recordData, _ := json.Marshal(record)
	d.client.Set(d.ctx, jobKey, recordData, 7*24*time.Hour)

	// 添加到延迟队列
	delayedKey := d.delayedKey(record.Queue)
	d.client.ZAdd(d.ctx, delayedKey, redis.Z{
		Score:  float64(record.ScheduledAt.Unix()),
		Member: jobID,
	})

	// 从处理中队列移除
	processingKey := d.processingKey(record.Queue)
	d.client.ZRem(d.ctx, processingKey, jobID)

	return nil
}

// Delete 删除任务
func (d *RedisDriver) Delete(jobID string) error {
	jobKey := d.jobKey(jobID)
	return d.client.Del(d.ctx, jobKey).Err()
}

// GetJob 获取任务信息
func (d *RedisDriver) GetJob(jobID string) (*JobRecord, error) {
	jobKey := d.jobKey(jobID)
	data, err := d.client.Get(d.ctx, jobKey).Result()
	if err != nil {
		return nil, err
	}

	var record JobRecord
	if err := json.Unmarshal([]byte(data), &record); err != nil {
		return nil, err
	}

	return &record, nil
}

// ListJobs 列出任务
func (d *RedisDriver) ListJobs(queue string, status JobStatus, limit int) ([]*JobRecord, error) {
	// 使用 SCAN 遍历所有任务
	var jobs []*JobRecord
	pattern := d.prefix + ":job:*"

	iter := d.client.Scan(d.ctx, 0, pattern, int64(limit*2)).Iterator()
	for iter.Next(d.ctx) {
		if len(jobs) >= limit {
			break
		}

		data, err := d.client.Get(d.ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var record JobRecord
		if err := json.Unmarshal([]byte(data), &record); err != nil {
			continue
		}

		if (queue == "" || record.Queue == queue) && (status == "" || record.Status == status) {
			jobs = append(jobs, &record)
		}
	}

	return jobs, iter.Err()
}

// GetStats 获取统计信息
func (d *RedisDriver) GetStats(queue string) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"pending":   0,
		"running":   0,
		"completed": 0,
		"failed":    0,
		"dead":      0,
	}

	// 统计各个队列的长度
	if queue != "" {
		queueKey := d.queueKey(queue)
		pending, _ := d.client.LLen(d.ctx, queueKey).Result()
		stats["pending"] = int(pending)

		delayedKey := d.delayedKey(queue)
		delayed, _ := d.client.ZCard(d.ctx, delayedKey).Result()
		stats["delayed"] = int(delayed)

		processingKey := d.processingKey(queue)
		running, _ := d.client.ZCard(d.ctx, processingKey).Result()
		stats["running"] = int(running)
	}

	// 死信队列
	deadKey := d.deadKey()
	dead, _ := d.client.LLen(d.ctx, deadKey).Result()
	stats["dead"] = int(dead)

	return stats, nil
}

// Close 关闭驱动
func (d *RedisDriver) Close() error {
	return nil // Redis 客户端由外部管理
}

// 键名辅助方法
func (d *RedisDriver) queueKey(queue string) string {
	return fmt.Sprintf("%s:queue:%s", d.prefix, queue)
}

func (d *RedisDriver) delayedKey(queue string) string {
	return fmt.Sprintf("%s:delayed:%s", d.prefix, queue)
}

func (d *RedisDriver) processingKey(queue string) string {
	return fmt.Sprintf("%s:processing:%s", d.prefix, queue)
}

func (d *RedisDriver) jobKey(jobID string) string {
	return fmt.Sprintf("%s:job:%s", d.prefix, jobID)
}

func (d *RedisDriver) deadKey() string {
	return fmt.Sprintf("%s:dead", d.prefix)
}
