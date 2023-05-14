package internal

import (
	"time"

	"github.com/gocql/gocql"
)

// ConnectToCassandra creates a connection to the Cassandra cluster and returns a session
func ConnectToCassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster("localhost:9042") // Replace with your Cassandra cluster address
	cluster.Keyspace = "collector"                // Replace with your keyspace name
	cluster.NumConns = 10                         // Maximum number of connections in the pool
	cluster.Timeout = 10 * time.Second            // Timeout for queries and other operations
	cluster.PoolConfig.HostSelectionPolicy = gocql.RoundRobinHostPolicy()

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
