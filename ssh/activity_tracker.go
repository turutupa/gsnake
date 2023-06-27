package ssh

type ActivityTracker struct {
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
