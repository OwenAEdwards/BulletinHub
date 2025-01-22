package bulletin

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection represents a user's WebSocket connection and state
type Connection struct {
	Username string          `bson:"username"`
	Socket   *websocket.Conn `bson:"-"` // WebSocket connections are not persisted
	Board    string          `bson:"board"`
}

// BulletinBoard manages public and private boards
type BulletinBoard struct {
	client *mongo.Client
	db     *mongo.Database
	boards *mongo.Collection
	conns  map[string]*Connection // Map of username to Connection
	mutex  sync.Mutex             // Mutex for connection management
}

// NewBulletinBoard creates a new instance of BulletinBoard
func NewBulletinBoard() *BulletinBoard {
	// Retrieve MongoDB URL from environment variables (provided by Docker), fallback to localhost
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017" // Default to localhost for development
	}

	// Initialize MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Set up database and collection
	db := client.Database("bulletin_board")
	boards := db.Collection("boards")

	return &BulletinBoard{
		client: client,
		db:     db,
		boards: boards,
		conns:  make(map[string]*Connection),
	}
}

// AddUser adds a user to a specific board
func (bb *BulletinBoard) AddUser(boardName string, user *Connection) {
	bb.mutex.Lock()
	bb.conns[user.Username] = user
	defer bb.mutex.Unlock()

	// Insert the user into MongoDB under the board name
	_, err := bb.boards.UpdateOne(
		context.Background(),
		bson.M{"board": boardName},
		bson.M{"$push": bson.M{"users": bson.M{"username": user.Username, "board": user.Board}}},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		log.Println("Error adding user:", err)
	}
}

// RemoveUser removes a user from a specific board
func (bb *BulletinBoard) RemoveUser(boardName string, user *Connection) {
	bb.mutex.Lock()
	delete(bb.conns, user.Username)
	defer bb.mutex.Unlock()

	// Remove the user from the MongoDB board
	_, err := bb.boards.UpdateOne(
		context.Background(),
		bson.M{"board": boardName},
		bson.M{"$pull": bson.M{"users": bson.M{"username": user.Username}}},
	)

	if err != nil {
		log.Println("Error removing user:", err)
	}
}

// ListBoards returns a list of all available board names
func (bb *BulletinBoard) ListBoards() []string {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	var boardNames []string
	cursor, err := bb.boards.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("Error fetching boards:", err)
		return boardNames
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Println("Error decoding board:", err)
			continue
		}
		boardNames = append(boardNames, result["board"].(string))
	}

	return boardNames
}

// ListUsers returns a list of usernames in a specific board
func (bb *BulletinBoard) ListUsers(boardName string) []string {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	var users []string
	var board struct {
		Users []Connection `bson:"users"`
	}

	err := bb.boards.FindOne(context.Background(), bson.M{"board": boardName}).Decode(&board)
	if err != nil {
		log.Println("Error fetching users:", err)
		return users
	}

	for _, conn := range board.Users {
		users = append(users, conn.Username)
	}

	return users
}

// BroadcastMessage sends a message to all users in a specific board
func (bb *BulletinBoard) BroadcastMessage(boardName string, message string) {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	// Retrieve all users in the specified board
	var board struct {
		Users []struct {
			Username string `bson:"username"`
		} `bson:"users"`
	}

	err := bb.boards.FindOne(context.Background(), bson.M{"board": boardName}).Decode(&board)
	if err != nil {
		log.Println("Error fetching board users:", err)
		return
	}

	// Broadcast the message to all active users
	for _, user := range board.Users {
		if conn, ok := bb.conns[user.Username]; ok {
			// Send a message to the websockets
			err := conn.Socket.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("Error sending message to %s: %v", user.Username, err)
			}
		}
	}
}

// BroadcastMessage sends a message to all users in a specific board, excluding the sender
func (bb *BulletinBoard) BroadcastMessageExcludingClient(boardName string, excludeClient *Connection, message string) {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	// Retrieve all users in the specified board
	var board struct {
		Users []struct {
			Username string `bson:"username"`
		} `bson:"users"`
	}

	err := bb.boards.FindOne(context.Background(), bson.M{"board": boardName}).Decode(&board)
	if err != nil {
		log.Println("Error fetching board users:", err)
		return
	}

	// Broadcast the message excluding the client
	for _, user := range board.Users {
		if conn, ok := bb.conns[user.Username]; ok && conn.Username != excludeClient.Username {
			// Send the message to the websockets
			err := conn.Socket.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("Error sending message to %s: %v", user.Username, err)
			}
		}
	}
}
