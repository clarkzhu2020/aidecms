package stats

import (
	"sync"
	"time"
)

type CommandStat struct {
	Name      string
	Count     int
	TotalTime time.Duration
	LastUsed  time.Time
}

var (
	stats   = make(map[string]CommandStat)
	statsMu sync.RWMutex
)

func Record(name string, duration time.Duration) {
	statsMu.Lock()
	defer statsMu.Unlock()

	stat, exists := stats[name]
	if !exists {
		stat = CommandStat{Name: name}
	}

	stat.Count++
	stat.TotalTime += duration
	stat.LastUsed = time.Now()
	stats[name] = stat
}

func GetStats() map[string]CommandStat {
	statsMu.RLock()
	defer statsMu.RUnlock()

	copy := make(map[string]CommandStat)
	for k, v := range stats {
		copy[k] = v
	}
	return copy
}
