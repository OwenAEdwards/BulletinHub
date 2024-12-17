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

	// Step 1: Wait for the username
	_, usernameMsg, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("Error reading username: %v\n", err)
		return
	}
	username := strings.TrimSpace(string(usernameMsg))
	if username == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Username cannot be empty"))
		return
	}
	fmt.Printf("Received username: %s\n", username)

	// Step 2: Wait for the join command
	_, joinMsg, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("Error reading join command: %v\n", err)
		return
	}
	joinCommand := strings.TrimSpace(string(joinMsg))
	if !strings.HasPrefix(joinCommand, "/join ") {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid join command. Use '/join <boardName>'"))
		return
	}

	// Extract board name
	boardName := strings.TrimSpace(strings.TrimPrefix(joinCommand, "/join "))
	if boardName == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Board name cannot be empty"))
		return
	}
	timestamp := utils.GetTimestamp()
	fmt.Printf("[%s] User %s joining board: %s\n", timestamp, username, boardName)

	// Create the connection object for the user
	client := &bulletin.Connection{Username: username, Socket: conn, Board: boardName}

	// Add client to the specified board
	bulletinBoard.AddUser(boardName, client)
	fmt.Printf("[%s] User %s connected to board: %s\n", timestamp, client.Username, boardName)

	// Broadcast the user's join message
	joinMessage := fmt.Sprintf("[%s] %s joined the board", timestamp, client.Username)
	bulletinBoard.BroadcastMessage(boardName, joinMessage)

	for {
		// Listen for messages
		_, msg, err := conn.ReadMessage()
		timestamp := utils.GetTimestamp() // Obtain a fresh timestamp
		if err != nil {
			fmt.Printf("[%s] Error reading message from %s: %v\n", timestamp, client.Username, err)

			// Broadcast the leave message with the current timestamp
			leaveMessage := fmt.Sprintf("[%s] %s left the chat", utils.GetTimestamp(), client.Username)
			bulletinBoard.RemoveUser(client.Board, client)
			bulletinBoard.BroadcastMessage(client.Board, leaveMessage)
			return
		}
		// Log the received message
		fmt.Printf("[%s] Received message from %s: %s\n", timestamp, client.Username, msg)

		// Handle commands (join, leave, etc.)
		handleClientMessage(client, string(msg))
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
		// If message not a command, then it's a regular message

		// Create the message with a timestamp
		timestamp := utils.GetTimestamp()
		fullMessage := fmt.Sprintf("[%s] %s: %s", timestamp, client.Username, message)

		// Broadcast the message to others
		bulletinBoard.BroadcastMessage(client.Board, fullMessage)
	}
}

// Join a new board
func joinBoard(client *bulletin.Connection, boardName string) {
	// Remove user from their current board (if any)
	bulletinBoard.RemoveUser(client.Board, client)

	// Add client to the new board
	client.Board = boardName
	bulletinBoard.AddUser(boardName, client)

	// Get the current timestamp
	timestamp := utils.GetTimestamp()

	// Send a confirmation message to the client
	joinedMessage := fmt.Sprintf("[%s] Joined board: %s", timestamp, boardName)
	client.Socket.WriteMessage(websocket.TextMessage, []byte(joinedMessage))

	// Broadcast the join message to other users on the new board
	broadcastMessage := fmt.Sprintf("[%s] %s joined the board", timestamp, client.Username)
	bulletinBoard.BroadcastMessage(boardName, broadcastMessage)

	// Log the join event on the server
	fmt.Printf("[%s] User %s joining board: %s\n", timestamp, client.Username, boardName)
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

	// Get the current timestamp
	timestamp := utils.GetTimestamp()

	// Notify the client they've left the board
	client.Socket.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[%s] Left board: %s", timestamp, boardName)))

	// Broadcast leave message to the board
	leaveMessage := fmt.Sprintf("[%s] %s left the board", timestamp, client.Username)
	bulletinBoard.BroadcastMessage(boardName, leaveMessage)

	// Log the leave event with a timestamp
	fmt.Printf("[%s] User %s leaving board: %s\n", timestamp, client.Username, boardName)
}

// List all available boards
func listBoards(client *bulletin.Connection) {
	// Get the current timestamp
	timestamp := utils.GetTimestamp()

	// Retrieve the list of boards
	boards := bulletinBoard.ListBoards()
	boardsList := strings.Join(boards, ", ")

	// Send the available boards list to the client with a timestamp
	client.Socket.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[%s] Available boards: %s", timestamp, boardsList)))

	// Log the request with the timestamp
	fmt.Printf("[%s] User %s requested available boards: %s\n", timestamp, client.Username, boardsList)
}

// GetBoardUsers handles requests to fetch the list of users in a specific board
func GetBoardUsers(w http.ResponseWriter, r *http.Request) {
	// Get the current timestamp
	timestamp := utils.GetTimestamp()

	// Extract boardName from the request URL
	boardName := strings.TrimPrefix(r.URL.Path, "/boards/")
	boardName = strings.TrimSuffix(boardName, "/users")

	// Log the board name being requested
	fmt.Printf("[%s] Request to get users for board: %s\n", timestamp, boardName)

	users := bulletinBoard.ListUsers(boardName)
	if users == nil {
		users = []string{} // Return an empty array if no users
	}

	// Log the users in the board
	fmt.Printf("[%s] Users in board %s: %v\n", timestamp, boardName, users)

	// Return the user list as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode user list", http.StatusInternalServerError)
		fmt.Printf("[%s] Error encoding user list for board %s: %v\n", timestamp, boardName, err)
	}
}
