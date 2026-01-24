package event

import (
	"sync"
)

// globalDispatcher 全局事件分发器
var (
	globalDispatcher *Dispatcher
	once             sync.Once
)

// GetDispatcher 获取全局事件分发器
func GetDispatcher() *Dispatcher {
	once.Do(func() {
		globalDispatcher = NewDispatcher(10)
	})
	return globalDispatcher
}

// Listen 注册全局事件监听器
func Listen(eventName string, listener Listener) {
	GetDispatcher().Listen(eventName, listener)
}

// ListenAsync 注册全局异步事件监听器
func ListenAsync(eventName string, listener Listener) {
	GetDispatcher().ListenAsync(eventName, listener)
}

// ListenWithPriority 注册全局带优先级的监听器
func ListenWithPriority(eventName string, listener Listener, priority int) {
	GetDispatcher().ListenWithPriority(eventName, listener, priority)
}

// Dispatch 分发事件到全局分发器
func Dispatch(event Event) error {
	return GetDispatcher().Dispatch(event)
}

// Forget 移除全局事件监听器
func Forget(eventName, listenerName string) {
	GetDispatcher().Forget(eventName, listenerName)
}

// ForgetAll 移除事件的所有全局监听器
func ForgetAll(eventName string) {
	GetDispatcher().ForgetAll(eventName)
}

// Subscribe 订阅多个事件到同一个全局监听器
func Subscribe(events []string, listener Listener) {
	GetDispatcher().Subscribe(events, listener)
}

// HasListeners 检查事件是否有全局监听器
func HasListeners(eventName string) bool {
	return GetDispatcher().HasListeners(eventName)
}

// GetStats 获取全局事件统计
func GetStats() map[string]interface{} {
	return GetDispatcher().GetStats()
}
