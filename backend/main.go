package main

import (
	"bulletin_board/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handlers.HandleWebSocket) // WebSocket endpoint

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil) // Start server
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
