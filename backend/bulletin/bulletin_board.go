package bulletin

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Connection represents a user's WebSocket connection and state
type Connection struct {
	Username string
	Socket   *websocket.Conn // WebSocket type
}

// BulletinBoard manages public and private boards
type BulletinBoard struct {
	Boards map[string][]*Connection // Key: board name, Value: list of connections
	mutex  sync.Mutex               // Mutex for thread-safe operations
}

// NewBulletinBoard creates a new instance of BulletinBoard
func NewBulletinBoard() *BulletinBoard {
	return &BulletinBoard{
		Boards: make(map[string][]*Connection),
	}
}

// AddUser adds a user to a specific board
func (bb *BulletinBoard) AddUser(boardName string, user *Connection) {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	bb.Boards[boardName] = append(bb.Boards[boardName], user)
}

// RemoveUser removes a user from a specific board
func (bb *BulletinBoard) RemoveUser(boardName string, user *Connection) {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	if _, exists := bb.Boards[boardName]; exists {
		for i, conn := range bb.Boards[boardName] {
			if conn == user {
				// Remove user from the slice
				bb.Boards[boardName] = append(bb.Boards[boardName][:i], bb.Boards[boardName][i+1:]...)
				break
			}
		}
	}
}

// ListUsers returns a list of usernames in a specific board
func (bb *BulletinBoard) ListUsers(boardName string) []string {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	var users []string
	if _, exists := bb.Boards[boardName]; exists {
		for _, conn := range bb.Boards[boardName] {
			users = append(users, conn.Username)
		}
	}
	return users
}

// BroadcastMessage sends a message to all users in a specific board
func (bb *BulletinBoard) BroadcastMessage(boardName string, message string) {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	if _, exists := bb.Boards[boardName]; exists {
		for _, conn := range bb.Boards[boardName] {
			// Send the message to the user's WebSocket
			err := conn.Socket.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
	}
}
