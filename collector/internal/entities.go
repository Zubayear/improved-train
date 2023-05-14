package internal

import "time"

type ServerMetricEntity struct {
	ServerIp          string    `json:"server_ip"`
	ServerName        string    `json:"server_name"`
	MetricTime        time.Time `json:"metric_time"`
	CPUUtilization    float32   `json:"cpu_utilization"`
	MemoryUtilization float32   `json:"memory_utilization"`
	DiskUtilization   float32   `json:"disk_utilization"`
}
