package handlers

import (
	"bulletin_board/bulletin"
	"bulletin_board/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		fmt.Println("WebSocket Origin:", origin) // Log the origin
		return origin == "http://localhost:3000" || origin == "http://localhost:5173"
	},
}

// Global state - reference to the BulletinBoard instance
var bulletinBoard = bulletin.NewBulletinBoard()

// HandleWebSocket manages WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Log WebSocket connection request details
	fmt.Println("Upgrading HTTP connection to WebSocket")
	fmt.Println("Request Method:", r.Method)
	fmt.Println("Request URL:", r.URL.Path)
	fmt.Println("Request Origin:", r.Header.Get("Origin"))

	// Set CORS headers explicitly here
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading to WebSocket: %v, Request Method: %s, Request URL: %s\n", err, r.Method, r.URL.Path)
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	fmt.Println("WebSocket connection established with", r.RemoteAddr)
	defer conn.Close()

	// Prompt client for a username
	conn.WriteMessage(websocket.TextMessage, []byte("Enter your username:"))
	_, username, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("Error reading username from %s: %v\n", r.RemoteAddr, err)
		return
	}

	// Ensure the username is not empty
	usernameStr := string(username)
	if usernameStr == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Username cannot be empty. Please provide a valid username."))
		return
	}

	// Create the connection object for the user
	client := &bulletin.Connection{Username: usernameStr, Socket: conn, Board: "public"}

	// Add client to the 'public' board by default
	bulletinBoard.AddUser("public", client)
	fmt.Printf("User %s connected to chat room\n", client.Username)

	// Broadcast the user's join message to the public board
	bulletinBoard.BroadcastMessage("public", fmt.Sprintf("%s joined the chat", client.Username))

	for {
		// Listen for messages
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message from %s: %v\n", client.Username, err)
			bulletinBoard.RemoveUser("public", client)
			bulletinBoard.BroadcastMessage("public", fmt.Sprintf("%s left the chat", client.Username))
			return
		}
		fmt.Printf("Received message from %s: %s\n", client.Username, msg)

		// Handle commands (join, leave, etc.)
		handleClientMessage(client, string(msg))

		// Broadcast new messages (not a command)
		if messageType == websocket.TextMessage {
			timestamp := utils.GetTimestamp()
			fullMessage := fmt.Sprintf("[%s] %s: %s", timestamp, client.Username, msg)
			bulletinBoard.BroadcastMessage(client.Board, fullMessage)
		}
	}
}

// Handle messages and commands from clients
func handleClientMessage(client *bulletin.Connection, message string) {
	// Check if the message is a command
	if len(message) > 0 && message[0] == '/' {
		// Split the command and arguments
		parts := strings.SplitN(message, " ", 2)
		command := parts[0]
		arg := ""
		if len(parts) > 1 {
			arg = parts[1]
		}

		fmt.Printf("Processing command '%s' from %s\n", command, client.Username)

		switch command {
		case "/join":
			if arg == "" {
				client.Socket.WriteMessage(websocket.TextMessage, []byte("Usage: /join <board_name>"))
				return
			}
			joinBoard(client, arg)

		case "/leave":
			if arg == "" {
				client.Socket.WriteMessage(websocket.TextMessage, []byte("Usage: /leave <board_name>"))
				return
			}
			leaveBoard(client, arg)

		case "/list":
			listBoards(client)

		default:
			client.Socket.WriteMessage(websocket.TextMessage, []byte("Unknown command: "+command))
		}
	} else {
		// If not a command, treat it as a message and broadcast
		bulletinBoard.BroadcastMessage(client.Board, fmt.Sprintf("%s: %s", client.Username, message))
	}
}

// Join a new board
func joinBoard(client *bulletin.Connection, boardName string) {
	// Remove user from their current board (if any)
	bulletinBoard.RemoveUser(client.Board, client)

	// Add client to the new board
	client.Board = boardName
	bulletinBoard.AddUser(boardName, client)
	client.Socket.WriteMessage(websocket.TextMessage, []byte("Joined board: "+boardName))

	// Broadcast join message to the new board
	bulletinBoard.BroadcastMessage(boardName, client.Username+" joined the board")
	fmt.Printf("User %s joining board: %s\n", client.Username, boardName)
}

// Leave the current board
func leaveBoard(client *bulletin.Connection, boardName string) {
	if boardName == "" {
		client.Socket.WriteMessage(websocket.TextMessage, []byte("You are not in any board"))
		return
	}

	// Remove client from the board
	bulletinBoard.RemoveUser(boardName, client)
	client.Board = ""
	client.Socket.WriteMessage(websocket.TextMessage, []byte("Left board: "+boardName))

	// Broadcast leave message to the board
	bulletinBoard.BroadcastMessage(boardName, client.Username+" left the board")
	fmt.Printf("User %s leaving board: %s\n", client.Username, boardName)
}

// List all available boards
func listBoards(client *bulletin.Connection) {
	boards := bulletinBoard.ListBoards()
	logMessage := "Available boards: " + strings.Join(boards, ", ")
	client.Socket.WriteMessage(websocket.TextMessage, []byte("Available boards: "+strings.Join(boards, ", ")))

	// Log the boards being returned
	fmt.Printf("User %s requested available boards: %s\n", client.Username, logMessage)
}

// GetBoardUsers handles requests to fetch the list of users in a specific board
func GetBoardUsers(w http.ResponseWriter, r *http.Request) {
	// Extract boardName from the request URL
	boardName := strings.TrimPrefix(r.URL.Path, "/boards/")
	boardName = strings.TrimSuffix(boardName, "/users")

	// Log the board name being requested
	fmt.Printf("Request to get users for board: %s\n", boardName)

	users := bulletinBoard.ListUsers(boardName)
	if users == nil {
		users = []string{} // Return an empty array if no users
	}

	// Log the users in the board
	fmt.Printf("Users in board %s: %v\n", boardName, users)

	// Return the user list as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode user list", http.StatusInternalServerError)
	}
}
