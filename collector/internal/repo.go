package internal

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"time"
)

type ServerMetricRepo interface {
	Insert(metrics *ServerMetricDto) error
	GetRecentMetric(serverIp string) (*ServerMetricEntity, error)
	GetServerMetricWithinTimeRange(serverIp, metricTime string) ([]*ServerMetricEntity, error)
}

// ServerMetricsRepository provides methods to interact with the server_metrics table
type ServerMetricsRepository struct {
	db *gocql.Session
}

// NewServerMetricsRepository creates a new instance of ServerMetricsRepository
func NewServerMetricsRepository(session *gocql.Session) *ServerMetricsRepository {
	return &ServerMetricsRepository{
		db: session,
	}
}

func (s *ServerMetricsRepository) Insert(metrics *ServerMetricDto) error {

	// Prepare the insert statement
	stmt := s.db.Query(`
		INSERT INTO server_metrics (
			server_ip, server_name, metric_time, 
			disk_utilization, cpu_utilization, memory_utilization
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	// Bind the values to the prepared statement
	if err := stmt.Bind(
		metrics.ServerIp,
		metrics.ServerName,
		time.Now(),
		metrics.DiskUtilization,
		metrics.CPUUtilization,
		metrics.MemoryUtilization,
	).Exec(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *ServerMetricsRepository) GetRecentMetric(serverIp string) (*ServerMetricEntity, error) {
	//query := "SELECT * FROM server_metrics WHERE server_ip = ? LIMIT 1 ALLOW FILTERING"
	query := `SELECT server_ip, server_name, metric_time, cpu_utilization, memory_utilization, disk_utilization FROM server_metrics WHERE server_ip = ? ORDER BY metric_time DESC LIMIT 1`

	var metric ServerMetricEntity

	if err := s.db.Query(query, serverIp).Scan(&metric.ServerIp,
		&metric.ServerName,
		&metric.MetricTime,
		&metric.DiskUtilization,
		&metric.CPUUtilization,
		&metric.MemoryUtilization); err != nil {
		return nil, err
	} else {
		return &metric, nil
	}
}

func (s *ServerMetricsRepository) GetServerMetricWithinTimeRange(serverIp, startTime, endTime string) ([]*ServerMetricEntity, error) {
	startTimeParsed := parseTime(startTime)
	endTimeParsed := parseTime(endTime)
	query := `SELECT server_ip, server_name, metric_time, cpu_utilization, memory_utilization, disk_utilization FROM server_metrics WHERE server_ip = ? AND metric_time >= ? AND metric_time <= ?`

	iter := s.db.Query(query, serverIp, startTimeParsed, endTimeParsed).Iter()

	var metrics []*ServerMetricEntity
	var metric ServerMetricEntity

	for iter.Scan(
		&metric.ServerIp,
		&metric.ServerName,
		&metric.MetricTime,
		&metric.CPUUtilization,
		&metric.MemoryUtilization,
		&metric.DiskUtilization,
	) {
		metrics = append(metrics, &metric)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func parseTime(timestampStr string) time.Time {

	layout := "2006-01-02 15:04:05.000-0700"
	timestamp, err := time.Parse(layout, timestampStr)
	if err != nil {
		fmt.Println("Failed to parse timestamp:", err)
		return time.Now()
	}
	return timestamp
}
