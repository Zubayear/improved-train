package endpoints

import (
	"fmt"
	"github.com/Zubayear/collector/internal"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime"
	"syscall"
	"time"
)

func SetupRouter(repository internal.ServerMetricsRepository) *gin.Engine {
	router := gin.Default()
	router.GET("api/v1/metrics/:server_ip", getRecentMetricHandler(repository))
	router.POST("api/v1/metrics/:server_ip", getServerMetricWithinTimeRangeHandler(repository))
	router.POST("api/v1/metrics", saveMetricsHandler(repository))
	router.GET("api/v1/metrics/disk", getDiskInfoHandler())
	router.GET("api/v1/metrics/cpu", getCpuInfoHandler())
	router.GET("api/v1/metrics/memory", getMemoryInfoHandler())
	return router
}

func getDiskInfoHandler() gin.HandlerFunc {
	// Get disk usage information
	stat := &syscall.Statfs_t{}
	err := syscall.Statfs("/", stat)
	if err != nil {
		log.Fatal("Failed to get disk usage:", err)
	}

	// Calculate disk utilization
	totalBytes := stat.Blocks * uint64(stat.Bsize)
	usedBytes := (stat.Blocks - stat.Bfree) * uint64(stat.Bsize)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	usedPercent := (float64(usedBytes) / float64(totalBytes)) * 100.0

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"total": bytesToGB(totalBytes), "used": bytesToGB(usedBytes), "free": bytesToGB(freeBytes), "percentage": usedPercent})
	}
}

func getCpuInfoHandler() gin.HandlerFunc {
	// Get CPU utilization information
	prevTime := time.Now()
	prevUsage := runtime.NumCPU()

	time.Sleep(time.Second)

	currTime := time.Now()
	currUsage := runtime.NumCPU()

	elapsedTime := currTime.Sub(prevTime)
	cpuUtilization := (float64(currUsage-prevUsage) / float64(runtime.NumCPU())) * 100.0
	cpuUtilizationPerCore := cpuUtilization / elapsedTime.Seconds()

	fmt.Println("CPU Usage:")
	fmt.Println("Utilization:", cpuUtilization)
	fmt.Println("Utilization per Core:", cpuUtilizationPerCore)
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"utilization": cpuUtilization, "utilization_per_core": cpuUtilizationPerCore})
	}
}

func getMemoryInfoHandler() gin.HandlerFunc {
	// Get memory utilization information
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	totalBytes := memStats.TotalAlloc
	usedBytes := memStats.Alloc
	freeBytes := totalBytes - usedBytes

	fmt.Println("Memory Usage:")
	fmt.Println("Total:", totalBytes)
	fmt.Println("Used:", usedBytes)
	fmt.Println("Free:", freeBytes)
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"total": bytesToGB(totalBytes), "used": bytesToGB(usedBytes), "free": bytesToGB(freeBytes)})
	}
}

func bytesToGB(bytes uint64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.2fGB", gb)
}

func saveMetricsHandler(repository internal.ServerMetricsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var serverMetric internal.ServerMetricDto
		if err := c.ShouldBindJSON(&serverMetric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := repository.Insert(&serverMetric); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metrics"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Metrics saved successfully"})
	}
}

func getRecentMetricHandler(repository internal.ServerMetricsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		serverIp := c.Param("server_ip")
		if valueFormRepo, err := repository.GetRecentMetric(serverIp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"data": valueFormRepo})
		}
	}
}

func getServerMetricWithinTimeRangeHandler(repository internal.ServerMetricsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		serverIp := c.Param("server_ip")
		startTime := c.PostForm("start_time")
		endTime := c.PostForm("end_time")
		if valueFormRepo, err := repository.GetServerMetricWithinTimeRange(serverIp, startTime, endTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"data": valueFormRepo})
		}
	}
}
