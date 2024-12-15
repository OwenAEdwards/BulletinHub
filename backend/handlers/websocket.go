package handlers

import (
	"bulletin_board/utils"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for testing
	},
}

// Connection struct to track user state
type Connection struct {
	Username string
	Board    string
	Socket   *websocket.Conn
}

// Global server state
var (
	clients        = make(map[*Connection]bool)     // Active connections
	bulletinBoards = make(map[string][]*Connection) // Public and private boards
	mutex          sync.Mutex                       // Mutex for safe concurrency
)

// HandleWebSocket manages WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Prompt client for a username
	conn.WriteMessage(websocket.TextMessage, []byte("Enter your username:"))
	_, username, err := conn.ReadMessage()
	if err != nil {
		return
	}
	client := &Connection{Username: string(username), Socket: conn}

	mutex.Lock()
	clients[client] = true
	bulletinBoards["public"] = append(bulletinBoards["public"], client) // Join public board by default
	mutex.Unlock()

	fmt.Println("User connected:", client.Username)
	broadcastMessage(fmt.Sprintf("%s joined the chat", client.Username), "public")

	for {
		// Listen for messages
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			removeClient(client)
			return
		}

		// Handle commands (join, leave, etc.)
		handleClientMessage(client, string(msg))

		// Broadcast new messages
		if messageType == websocket.TextMessage {
			timestamp := utils.GetTimestamp()
			fullMessage := fmt.Sprintf("[%s] %s: %s", timestamp, client.Username, msg)
			broadcastMessage(fullMessage, client.Board)
		}
	}
}

// Handle messages and commands from clients
func handleClientMessage(client *Connection, message string) {
	// Check if the message is a command
	if len(message) > 0 && message[0] == '/' {
		// Split the command and arguments
		parts := strings.SplitN(message, " ", 2)
		command := parts[0]
		arg := ""
		if len(parts) > 1 {
			arg = parts[1]
		}

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
		broadcastMessage(fmt.Sprintf("%s: %s", client.Username, message), client.Board)
	}
}

// Broadcast a message to all clients in a specific board
func broadcastMessage(message, board string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := bulletinBoards[board]; exists {
		for _, conn := range bulletinBoards[board] {
			err := conn.Socket.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Error broadcasting message:", err)
			}
		}
	}
}

// Remove a client from all boards and close connection
func removeClient(client *Connection) {
	mutex.Lock()
	defer mutex.Unlock()

	for board, connections := range bulletinBoards {
		for i, conn := range connections {
			if conn == client {
				bulletinBoards[board] = append(connections[:i], connections[i+1:]...)
			}
		}
	}
	delete(clients, client)
	fmt.Printf("User %s disconnected\n", client.Username)
}

func joinBoard(client *Connection, boardName string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Leave the current board if any
	leaveBoard(client, client.Board)

	// Add client to the new board
	bulletinBoards[boardName] = append(bulletinBoards[boardName], client)
	client.Board = boardName
	client.Socket.WriteMessage(websocket.TextMessage, []byte("Joined board: "+boardName))
	broadcastMessage(client.Username+" joined the board", boardName)
}

func leaveBoard(client *Connection, boardName string) {
	mutex.Lock()
	defer mutex.Unlock()

	if boardName == "" {
		return
	}

	if connections, exists := bulletinBoards[boardName]; exists {
		for i, conn := range connections {
			if conn == client {
				// Remove client from the board
				bulletinBoards[boardName] = append(connections[:i], connections[i+1:]...)
				break
			}
		}
		client.Socket.WriteMessage(websocket.TextMessage, []byte("Left board: "+boardName))
		broadcastMessage(client.Username+" left the board", boardName)
	}
	client.Board = ""
}

func listBoards(client *Connection) {
	mutex.Lock()
	defer mutex.Unlock()

	var boards []string
	for board := range bulletinBoards {
		boards = append(boards, board)
	}

	client.Socket.WriteMessage(websocket.TextMessage, []byte("Available boards: "+strings.Join(boards, ", ")))
}
