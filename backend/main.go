package main

import (
	router "bulletin_board/routers" // Import the router package
	"fmt"
	"net/http"
)

func main() {
	// Set up the server with the router that includes CORS middleware
	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", router.SetupRouter()) // Use the SetupRouter from router package
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
