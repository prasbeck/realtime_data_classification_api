package main

import (
	"fmt"
	"log"
	"net/http"

	"rtmc/db"
	"rtmc/dsl"
	"rtmc/routes"
	"rtmc/websocket"
)

func main() {
    // Connect to the database
    database, err := db.ConnectDB()
    if err != nil {
        log.Fatalf("Could not connect to the database: %v\n", err)
    }
    defer database.Close()

    // Create user schema
    if err := db.CreateUserSchema(database); err != nil {
        log.Fatalf("Could not create user schema: %v\n", err)
    }

    // Set the maximum number of rules
    dsl.SetMaxRules(10)

    // Initialize WebSocket manager
    wsManager := websocket.NewManager()
    go wsManager.Start()

    // Register routes
    routes.RegisterRoutes(database, wsManager)

    // Define the port to listen on
    port := "8080"
    fmt.Printf("Starting server at port %s\n", port)

    // Start the HTTP server
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
}

