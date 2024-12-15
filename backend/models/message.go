package message

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents a message sent by a user in the chat
type Message struct {
	Username  string
	Content   string
	Timestamp string
}

// NewMessage creates a new message instance with a timestamp
func NewMessage(username, content string) *Message {
	return &Message{
		Username:  username,
		Content:   content,
		Timestamp: time.Now().Format(time.RFC3339), // Format timestamp in ISO 8601 format
	}
}

// FormatMessage formats the message as a string to be sent via WebSocket
func (m *Message) FormatMessage() string {
	return fmt.Sprintf("[%s] %s: %s", m.Timestamp, m.Username, m.Content)
}

// BroadcastMessage sends a message to all users in a specific board via WebSocket
func BroadcastMessage(boardName string, message *Message, connections []*websocket.Conn) {
	for _, conn := range connections {
		// Send the formatted message to each connection in the board
		err := conn.WriteMessage(websocket.TextMessage, []byte(message.FormatMessage()))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
