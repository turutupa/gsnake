package gsnake

import "sync"

var (
	m sync.RWMutex
)

type Stats struct {
	ActiveUsers   int
	RequestCount  int
	ErrorCount    int
	TrafficVolume int64
	CPUUsage      float64
	MemoryUsage   float64
	QueriesCount  int
	APIUsage      map[string]int
	UserActivity  []string
}
