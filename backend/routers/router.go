package routers

import (
	"bulletin_board/handlers"
	"net/http"

	"github.com/rs/cors"
)

// SetupRouter sets up the HTTP router and applies CORS middleware
func SetupRouter() http.Handler {
	// Create a new ServeMux (HTTP router)
	router := http.NewServeMux()

	// WebSocket endpoint
	router.HandleFunc("/ws", handlers.HandleWebSocket)

	// CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Return the router wrapped with the CORS handler
	return corsHandler.Handler(router)
}
