package main

import (
	"github.com/Zubayear/collector/endpoints"
	"github.com/Zubayear/collector/internal"
	"log"
)

func main() {
	session, err := internal.ConnectToCassandra()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Create a new instance of the repository
	repository := internal.NewServerMetricsRepository(session)

	// Setup the Gin router
	router := endpoints.SetupRouter(*repository)

	// Start the server
	router.Run(":8080")
}
