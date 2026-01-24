package shared

import "time"

type CommandStat struct {
	Name      string
	Count     int
	TotalTime time.Duration
	LastUsed  time.Time
}

var Stats = make(map[string]CommandStat)
