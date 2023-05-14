package internal

type ServerMetricDto struct {
	ServerIp          string  `json:"server_ip"`
	ServerName        string  `json:"server_name"`
	CPUUtilization    float32 `json:"cpu_utilization"`
	MemoryUtilization float32 `json:"memory_utilization"`
	DiskUtilization   float32 `json:"disk_utilization"`
}
