package commands

// 注意：这个文件包含了命令统计功能，与stats包中的统计功能不同

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	stats      map[string]CommandStat
	statsMutex sync.Mutex
)

func GetCommandStats() map[string]CommandStat {
	loadStats()
	statsMutex.Lock()
	defer statsMutex.Unlock()

	copy := make(map[string]CommandStat)
	for k, v := range stats {
		copy[k] = v
	}
	return copy
}

func RecordCommandUsage(name string, duration time.Duration) {
	loadStats()

	if stat, exists := stats[name]; exists {
		stat.Count++
		stat.TotalTime += duration
		stat.LastUsed = time.Now()
		stats[name] = stat
	} else {
		stats[name] = CommandStat{
			Name:      name,
			Count:     1,
			TotalTime: duration,
			LastUsed:  time.Now(),
		}
	}

	saveStats()
}

func ShowStats(args []string) {
	loadStats()

	if len(stats) == 0 {
		fmt.Println("No command statistics available")
		return
	}

	var statList []CommandStat
	for _, stat := range stats {
		statList = append(statList, stat)
	}

	sort.Slice(statList, func(i, j int) bool {
		return statList[i].Count > statList[j].Count
	})

	fmt.Println("Command Usage Statistics:")
	fmt.Printf("%-20s %-10s %-15s %-20s\n", "Command", "Count", "Avg Time", "Last Used")
	for _, stat := range statList {
		avgTime := stat.TotalTime / time.Duration(stat.Count)
		fmt.Printf("%-20s %-10d %-15v %-20s\n",
			stat.Name,
			stat.Count,
			avgTime.Round(time.Millisecond),
			stat.LastUsed.Format("2006-01-02 15:04:05"))
	}
}

func ResetStats(args []string) {
	stats = make(map[string]CommandStat)
	saveStats()
	fmt.Println("Command statistics have been reset")
}

func ExportStats(args []string) {
	loadStats()
	if len(args) < 1 {
		fmt.Println("Format is required (json or csv)")
		return
	}

	format := args[0]
	filePath := filepath.Join("storage", "stats", "command_stats_export."+format)

	switch format {
	case "json":
		exportJSON(filePath)
	case "csv":
		exportCSV(filePath)
	default:
		fmt.Println("Unsupported format. Use json or csv")
	}
}

func loadStats() {
	stats = make(map[string]CommandStat)

	filePath := filepath.Join("storage", "stats", "command_stats.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	json.Unmarshal(data, &stats)
}

func saveStats() {
	filePath := filepath.Join("storage", "stats", "command_stats.json")
	os.MkdirAll(filepath.Dir(filePath), 0755)

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(filePath, data, 0644)
}

func exportJSON(filePath string) {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		fmt.Printf("Failed to export stats: %v\n", err)
		return
	}

	os.WriteFile(filePath, data, 0644)
	fmt.Printf("Statistics exported to %s\n", filePath)
}

func exportCSV(filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Command", "Count", "Total Time (ms)", "Avg Time (ms)", "Last Used"})

	for _, stat := range stats {
		avgTime := stat.TotalTime / time.Duration(stat.Count)
		writer.Write([]string{
			stat.Name,
			fmt.Sprintf("%d", stat.Count),
			fmt.Sprintf("%d", stat.TotalTime.Milliseconds()),
			fmt.Sprintf("%d", avgTime.Milliseconds()),
			stat.LastUsed.Format("2006-01-02 15:04:05"),
		})
	}

	fmt.Printf("Statistics exported to %s\n", filePath)
}
