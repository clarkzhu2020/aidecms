package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type EmailJob struct {
	ID          string
	Subject     string
	Body        string
	Recipient   string
	Status      string // pending, sent, failed
	Priority    int    // 1-5, 1 is highest
	CreatedAt   time.Time
	SentAt      time.Time
	RetryCount  int
	RateLimited bool
	Template    string            // 模板名称
	Variables   map[string]string // 模板变量
}

var (
	queue           []EmailJob
	lastSentTime    time.Time
	sendInterval    = time.Second * 2 // 初始速率
	minInterval     = time.Second * 1 // 最快1秒1封
	maxInterval     = time.Second * 5 // 最慢5秒1封
	failureCount    int
	successCount    int
	adjustRateAfter = 10        // 每10次发送后调整速率
	isPaused        bool        // 队列是否暂停
	pauseTimer      *time.Timer // 定时暂停计时器
	resumeTimer     *time.Timer // 定时恢复计时器
)

func init() {
	loadQueue()
}

func AddToQueue(subject, body, recipient string, template string, vars map[string]string) {
	job := EmailJob{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Subject:   subject,
		Body:      body,
		Recipient: recipient,
		Status:    "pending",
		Priority:  3, // 默认中等优先级
		CreatedAt: time.Now(),
		Template:  template,
		Variables: vars,
	}

	queue = append(queue, job)
	saveQueue()
}

func ProcessQueue() {
	if isPaused {
		return
	}

	// 按优先级排序
	sort.Slice(queue, func(i, j int) bool {
		return queue[i].Priority < queue[j].Priority
	})

	for i, job := range queue {
		if job.Status == "pending" || job.Status == "failed" {
			// 速率控制
			if time.Since(lastSentTime) < sendInterval {
				queue[i].RateLimited = true
				continue
			}

			// 处理模板变量
			body := job.Body
			if job.Template != "" {
				body = processTemplate(job.Body, job.Variables)
			}

			err := SendAlertEmail(job.Subject, body)
			lastSentTime = time.Now()

			if err != nil {
				queue[i].Status = "failed"
				queue[i].RetryCount++
				failureCount++
			} else {
				queue[i].Status = "sent"
				queue[i].SentAt = time.Now()
				queue[i].RateLimited = false
				successCount++
			}

			// 动态调整发送速率
			if (failureCount+successCount)%adjustRateAfter == 0 {
				adjustSendRate()
			}
			saveQueue()
		}
	}
}

func ShowQueueStatus(args []string) {
	loadQueue()

	fmt.Println("\nEmail Queue Status:")
	if isPaused {
		fmt.Println("⚠️ Queue is currently PAUSED")
	}
	fmt.Printf("%-10s %-20s %-10s %-15s %-10s %-10s\n",
		"ID", "Subject", "Status", "Created", "Retries", "Priority")

	for _, job := range queue {
		status := job.Status
		if job.RateLimited {
			status += "(rate limited)"
		}

		fmt.Printf("%-10s %-20s %-10s %-15s %-10d %-10d\n",
			job.ID[:8],
			truncate(job.Subject, 20),
			status,
			job.CreatedAt.Format("2006-01-02"),
			job.RetryCount,
			job.Priority)
	}
}

func PauseQueue(args []string) {
	duration := parseDuration(args)
	if duration > 0 {
		if pauseTimer != nil {
			pauseTimer.Stop()
		}
		pauseTimer = time.AfterFunc(duration, func() {
			isPaused = true
			fmt.Println("\nQueue processing paused automatically")
		})
		fmt.Printf("Queue will pause after %v\n", duration)
		return
	}

	isPaused = true
	fmt.Println("Queue processing paused")
}

func ResumeQueue(args []string) {
	duration := parseDuration(args)
	if duration > 0 {
		if resumeTimer != nil {
			resumeTimer.Stop()
		}
		resumeTimer = time.AfterFunc(duration, func() {
			isPaused = false
			fmt.Println("\nQueue processing resumed automatically")
		})
		fmt.Printf("Queue will resume after %v\n", duration)
		return
	}

	isPaused = false
	fmt.Println("Queue processing resumed")
}

func parseDuration(args []string) time.Duration {
	if len(args) == 0 {
		return 0
	}

	duration, err := time.ParseDuration(args[0])
	if err != nil {
		return 0
	}
	return duration
}

func SetPriority(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: queue:priority <jobID> <priority(1-5)> [jobID2 priority2 ...]")
		return
	}

	// 批量设置优先级
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}

		id := args[i]
		priority := 0
		fmt.Sscanf(args[i+1], "%d", &priority)

		if priority < 1 || priority > 5 {
			fmt.Printf("Invalid priority %d for job %s\n", priority, id)
			continue
		}

		found := false
		for j := range queue {
			if queue[j].ID == id || strings.HasPrefix(queue[j].ID, id) {
				queue[j].Priority = priority
				saveQueue()
				fmt.Printf("Priority updated for job %s\n", id)
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Job %s not found\n", id)
		}
	}
}

func ShowQueueStats(args []string) {
	loadQueue()

	total := len(queue)
	var pending, sent, failed int
	var totalTime time.Duration

	for _, job := range queue {
		switch job.Status {
		case "pending":
			pending++
		case "sent":
			sent++
			totalTime += job.SentAt.Sub(job.CreatedAt)
		case "failed":
			failed++
		}
	}

	avgTime := time.Duration(0)
	if sent > 0 {
		avgTime = totalTime / time.Duration(sent)
	}

	fmt.Println("\nQueue Statistics:")
	fmt.Printf("Total jobs:    %d\n", total)
	fmt.Printf("Pending:       %d\n", pending)
	fmt.Printf("Sent:          %d\n", sent)
	fmt.Printf("Failed:        %d\n", failed)
	fmt.Printf("Avg send time: %v\n", avgTime.Round(time.Second))
	fmt.Printf("Current rate:  %v per email\n", sendInterval)
}

func adjustSendRate() {
	successRate := float64(successCount) / float64(successCount+failureCount)

	switch {
	case successRate > 0.9: // 成功率>90%，加快发送
		sendInterval = max(minInterval, sendInterval/2)
	case successRate < 0.7: // 成功率<70%，减慢发送
		sendInterval = min(maxInterval, sendInterval*2)
	}

	// 重置计数器
	failureCount = 0
	successCount = 0
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func max(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func RetryFailedJobs(args []string) {
	loadQueue()
	count := 0

	for i, job := range queue {
		if job.Status == "failed" && job.RetryCount < 3 {
			queue[i].Status = "pending"
			count++
		}
	}

	saveQueue()
	fmt.Printf("Marked %d failed jobs for retry\n", count)
}

func CleanQueue(args []string) {
	loadQueue()
	threshold := time.Now().AddDate(0, 0, -7) // 保留7天
	count := 0

	var newQueue []EmailJob
	for _, job := range queue {
		if job.Status == "sent" && job.SentAt.Before(threshold) {
			count++
		} else {
			newQueue = append(newQueue, job)
		}
	}

	queue = newQueue
	saveQueue()
	fmt.Printf("Cleaned up %d old jobs\n", count)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func processTemplate(template string, vars map[string]string) string {
	for k, v := range vars {
		template = strings.ReplaceAll(template, "{{"+k+"}}", v)
	}
	return template
}

func saveQueue() {
	filePath := filepath.Join("storage", "queue", "email_queue.json")
	os.MkdirAll(filepath.Dir(filePath), 0755)

	data, err := json.MarshalIndent(queue, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(filePath, data, 0644)
}

func loadQueue() {
	filePath := filepath.Join("storage", "queue", "email_queue.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	json.Unmarshal(data, &queue)
}
