package stats

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/clarkzhu2020/aidecms/cmd/artisan/commands"
)

type StatRecorder interface {
	Record(name string, duration time.Duration)
	GetStats() map[string]CommandStat
}

type CommandStat = commands.CommandStat

type Stats struct {
	data map[string]CommandStat
	mu   sync.RWMutex
}

func New() *Stats {
	return &Stats{
		data: make(map[string]CommandStat),
	}
}

func (s *Stats) Record(name string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat := s.data[name]
	stat.Name = name
	stat.Count++
	stat.TotalTime += duration
	stat.LastUsed = time.Now()
	s.data[name] = stat
}

func (s *Stats) GetStats() map[string]CommandStat {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copy := make(map[string]CommandStat)
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}

func GenerateChart(stats map[string]CommandStat) {
	if len(stats) == 0 {
		fmt.Println("No statistics available to generate chart")
		return
	}

	fmt.Println("\nCommand Usage Chart:")
	fmt.Println("-------------------")

	var statList []CommandStat
	for _, stat := range stats {
		statList = append(statList, stat)
	}

	sort.Slice(statList, func(i, j int) bool {
		return statList[i].Count > statList[j].Count
	})

	maxCount := statList[0].Count
	for _, stat := range statList {
		barLength := int(float64(stat.Count) / float64(maxCount) * 50)
		fmt.Printf("%-20s |%s %d\n", stat.Name, strings.Repeat("=", barLength), stat.Count)
	}
}

func CleanupOldStats(stats map[string]CommandStat, threshold time.Time) bool {
	cleaned := false
	for key, stat := range stats {
		if stat.LastUsed.Before(threshold) {
			delete(stats, key)
			cleaned = true
		}
	}
	return cleaned
}

func CheckForAnomalies(stats map[string]CommandStat, threshold time.Duration) bool {
	hasAnomalies := false
	for _, stat := range stats {
		avgTime := stat.TotalTime / time.Duration(stat.Count)
		if avgTime > threshold {
			fmt.Printf("[WARNING] Command '%s' has long average execution time: %v\n", stat.Name, avgTime)
			hasAnomalies = true
		}
	}
	return hasAnomalies
}
