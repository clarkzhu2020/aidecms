package event

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Event 事件接口
type Event interface {
	// EventName 返回事件名称
	EventName() string
}

// Listener 监听器函数类型
type Listener func(ctx context.Context, event Event) error

// ListenerWrapper 监听器包装器
type ListenerWrapper struct {
	Name     string
	Handler  Listener
	Priority int  // 优先级，数字越小优先级越高
	Async    bool // 是否异步执行
}

// Dispatcher 事件分发器
type Dispatcher struct {
	listeners map[string][]*ListenerWrapper
	mu        sync.RWMutex
	queue     chan *eventJob
	workers   int
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	logs      []EventLog
	logsMu    sync.RWMutex
	maxLogs   int
}

// eventJob 事件任务
type eventJob struct {
	event    Event
	listener *ListenerWrapper
	ctx      context.Context
}

// EventLog 事件日志
type EventLog struct {
	EventName    string
	ListenerName string
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	Success      bool
	Error        string
	Async        bool
}

// NewDispatcher 创建新的事件分发器
func NewDispatcher(workers int) *Dispatcher {
	if workers <= 0 {
		workers = 10
	}

	ctx, cancel := context.WithCancel(context.Background())
	d := &Dispatcher{
		listeners: make(map[string][]*ListenerWrapper),
		queue:     make(chan *eventJob, 1000),
		workers:   workers,
		ctx:       ctx,
		cancel:    cancel,
		logs:      make([]EventLog, 0),
		maxLogs:   1000,
	}

	// 启动工作进程
	d.startWorkers()

	return d
}

// Listen 注册事件监听器
func (d *Dispatcher) Listen(eventName string, listener Listener) *Dispatcher {
	return d.ListenWithOptions(eventName, "", listener, 0, false)
}

// ListenAsync 注册异步事件监听器
func (d *Dispatcher) ListenAsync(eventName string, listener Listener) *Dispatcher {
	return d.ListenWithOptions(eventName, "", listener, 0, true)
}

// ListenWithPriority 注册带优先级的监听器
func (d *Dispatcher) ListenWithPriority(eventName string, listener Listener, priority int) *Dispatcher {
	return d.ListenWithOptions(eventName, "", listener, priority, false)
}

// ListenWithOptions 注册监听器（完整选项）
func (d *Dispatcher) ListenWithOptions(eventName, name string, listener Listener, priority int, async bool) *Dispatcher {
	d.mu.Lock()
	defer d.mu.Unlock()

	if name == "" {
		name = fmt.Sprintf("listener_%d", time.Now().UnixNano())
	}

	wrapper := &ListenerWrapper{
		Name:     name,
		Handler:  listener,
		Priority: priority,
		Async:    async,
	}

	if d.listeners[eventName] == nil {
		d.listeners[eventName] = make([]*ListenerWrapper, 0)
	}

	d.listeners[eventName] = append(d.listeners[eventName], wrapper)

	// 按优先级排序（数字越小优先级越高）
	d.sortListeners(eventName)

	return d
}

// sortListeners 排序监听器
func (d *Dispatcher) sortListeners(eventName string) {
	listeners := d.listeners[eventName]
	for i := 0; i < len(listeners)-1; i++ {
		for j := i + 1; j < len(listeners); j++ {
			if listeners[i].Priority > listeners[j].Priority {
				listeners[i], listeners[j] = listeners[j], listeners[i]
			}
		}
	}
}

// Dispatch 分发事件
func (d *Dispatcher) Dispatch(event Event) error {
	return d.DispatchWithContext(context.Background(), event)
}

// DispatchWithContext 使用自定义 context 分发事件
func (d *Dispatcher) DispatchWithContext(ctx context.Context, event Event) error {
	d.mu.RLock()
	listeners := d.listeners[event.EventName()]
	d.mu.RUnlock()

	if len(listeners) == 0 {
		return nil
	}

	var syncErrors []error

	for _, listener := range listeners {
		if listener.Async {
			// 异步执行
			select {
			case d.queue <- &eventJob{
				event:    event,
				listener: listener,
				ctx:      ctx,
			}:
			default:
				// 队列满了，记录警告
				fmt.Printf("Warning: event queue full, dropping async listener %s for event %s\n",
					listener.Name, event.EventName())
			}
		} else {
			// 同步执行
			if err := d.executeListener(ctx, event, listener); err != nil {
				syncErrors = append(syncErrors, err)
			}
		}
	}

	if len(syncErrors) > 0 {
		return fmt.Errorf("listeners failed: %v", syncErrors)
	}

	return nil
}

// executeListener 执行监听器
func (d *Dispatcher) executeListener(ctx context.Context, event Event, listener *ListenerWrapper) error {
	log := EventLog{
		EventName:    event.EventName(),
		ListenerName: listener.Name,
		StartTime:    time.Now(),
		Async:        listener.Async,
	}

	err := listener.Handler(ctx, event)

	log.EndTime = time.Now()
	log.Duration = log.EndTime.Sub(log.StartTime)

	if err != nil {
		log.Success = false
		log.Error = err.Error()
	} else {
		log.Success = true
	}

	d.addLog(log)

	return err
}

// startWorkers 启动工作进程
func (d *Dispatcher) startWorkers() {
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go d.worker(i)
	}
}

// worker 工作进程
func (d *Dispatcher) worker(id int) {
	defer d.wg.Done()

	for {
		select {
		case <-d.ctx.Done():
			return
		case job := <-d.queue:
			d.executeListener(job.ctx, job.event, job.listener)
		}
	}
}

// Stop 停止事件分发器
func (d *Dispatcher) Stop() {
	d.cancel()
	d.wg.Wait()
	close(d.queue)
}

// Forget 移除事件监听器
func (d *Dispatcher) Forget(eventName, listenerName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	listeners := d.listeners[eventName]
	newListeners := make([]*ListenerWrapper, 0)

	for _, listener := range listeners {
		if listener.Name != listenerName {
			newListeners = append(newListeners, listener)
		}
	}

	d.listeners[eventName] = newListeners
}

// ForgetAll 移除事件的所有监听器
func (d *Dispatcher) ForgetAll(eventName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.listeners, eventName)
}

// GetListeners 获取事件的所有监听器
func (d *Dispatcher) GetListeners(eventName string) []*ListenerWrapper {
	d.mu.RLock()
	defer d.mu.RUnlock()

	listeners := d.listeners[eventName]
	result := make([]*ListenerWrapper, len(listeners))
	copy(result, listeners)
	return result
}

// HasListeners 检查事件是否有监听器
func (d *Dispatcher) HasListeners(eventName string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return len(d.listeners[eventName]) > 0
}

// GetAllEvents 获取所有注册的事件名称
func (d *Dispatcher) GetAllEvents() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	events := make([]string, 0, len(d.listeners))
	for eventName := range d.listeners {
		events = append(events, eventName)
	}
	return events
}

// addLog 添加日志
func (d *Dispatcher) addLog(log EventLog) {
	d.logsMu.Lock()
	defer d.logsMu.Unlock()

	d.logs = append(d.logs, log)

	// 限制日志大小
	if len(d.logs) > d.maxLogs {
		d.logs = d.logs[len(d.logs)-d.maxLogs:]
	}
}

// GetLogs 获取事件日志
func (d *Dispatcher) GetLogs(eventName string, limit int) []EventLog {
	d.logsMu.RLock()
	defer d.logsMu.RUnlock()

	var logs []EventLog
	count := 0

	// 从后往前遍历（最新的日志在前）
	for i := len(d.logs) - 1; i >= 0 && count < limit; i-- {
		if eventName == "" || d.logs[i].EventName == eventName {
			logs = append(logs, d.logs[i])
			count++
		}
	}

	return logs
}

// GetStats 获取统计信息
func (d *Dispatcher) GetStats() map[string]interface{} {
	d.mu.RLock()
	d.logsMu.RLock()
	defer d.mu.RUnlock()
	defer d.logsMu.RUnlock()

	totalEvents := len(d.listeners)
	totalListeners := 0
	for _, listeners := range d.listeners {
		totalListeners += len(listeners)
	}

	totalExecutions := len(d.logs)
	successCount := 0
	for _, log := range d.logs {
		if log.Success {
			successCount++
		}
	}

	successRate := 0.0
	if totalExecutions > 0 {
		successRate = float64(successCount) / float64(totalExecutions) * 100
	}

	return map[string]interface{}{
		"total_events":     totalEvents,
		"total_listeners":  totalListeners,
		"total_executions": totalExecutions,
		"success_count":    successCount,
		"fail_count":       totalExecutions - successCount,
		"success_rate":     successRate,
		"queue_size":       len(d.queue),
		"workers":          d.workers,
	}
}

// Subscribe 订阅多个事件到同一个监听器
func (d *Dispatcher) Subscribe(events []string, listener Listener) *Dispatcher {
	for _, eventName := range events {
		d.Listen(eventName, listener)
	}
	return d
}

// Until 分发事件直到第一个监听器返回非 nil 结果
func (d *Dispatcher) Until(event Event) (interface{}, error) {
	d.mu.RLock()
	listeners := d.listeners[event.EventName()]
	d.mu.RUnlock()

	for _, listener := range listeners {
		if err := listener.Handler(context.Background(), event); err != nil {
			return nil, err
		}
		// 如果需要返回值，可以从 event 中获取
		// 这里简化处理
	}

	return nil, nil
}

// GetEventType 获取事件类型名称
func GetEventType(event Event) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
