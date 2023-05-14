package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type ServerMetrics struct {
	ServerIP          string  `json:"server_ip"`
	ServerName        string  `json:"server_name"`
	DiskUtilization   float32 `json:"disk_utilization"`
	CPUUtilization    float32 `json:"cpu_utilization"`
	MemoryUtilization float32 `json:"memory_utilization"`
}

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Failed to get hostname:", err)
		return
	}
	fmt.Println("Hostname:", hostname)

	hostIP, err := getHostIP()
	if err != nil {
		fmt.Println("Failed to get host IP:", err)
		return
	}
	fmt.Println("Host IP:", hostIP)

	// Run the monitoring loop indefinitely
	for {
		// Get disk usage information
		diskUsage, err := disk.Usage("/")
		if err != nil {
			fmt.Println("Failed to get disk usage:", err)
		} else {
			fmt.Println("Disk Usage:")
			fmt.Println("Total:", diskUsage.Total)
			fmt.Println("Used:", diskUsage.Used)
			fmt.Println("Free:", diskUsage.Free)
			fmt.Println("Percentage:", diskUsage.UsedPercent)
			fmt.Println()
		}

		// Get CPU usage information
		cpuUsage, err := cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Println("Failed to get CPU usage:", err)
		} else {
			fmt.Println("CPU Usage:")
			fmt.Println("Usage:", cpuUsage)
			fmt.Println()
		}

		// Get memory usage information
		memUsage, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Failed to get memory usage:", err)
		} else {
			fmt.Println("Memory Usage:")
			fmt.Println("Total:", memUsage.Total)
			fmt.Println("Used:", memUsage.Used)
			fmt.Println("Free:", memUsage.Free)
			fmt.Println("Percentage:", memUsage.UsedPercent)
			fmt.Println()
		}

		// Create a new instance of ServerMetrics and populate its fields
		metrics := ServerMetrics{
			ServerIP:          hostIP,
			ServerName:        hostname,
			DiskUtilization:   float32(diskUsage.UsedPercent),
			CPUUtilization:    float32(cpuUsage[0]),
			MemoryUtilization: float32(memUsage.UsedPercent),
		}

		// Encode the struct into JSON
		jsonData, err := json.Marshal(metrics)
		if err != nil {
			fmt.Println("Failed to encode JSON:", err)
			return
		}

		go postMetrics("http://127.0.0.1:8080/api/v1/metrics", string(jsonData))

		// Wait for some time before fetching the metrics again
		time.Sleep(time.Second * 25)
	}
}

func postMetrics(url, jsonData string) {

	payload := strings.NewReader(jsonData)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

// getHostIP retrieves the host IP address
func getHostIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("No host IP address found")
}
