package main

import (
	// "encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

// Player represents a single player's state
type Player struct {
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Radius     float64 `json:"radius"`
	Color      string  `json:"color"`
	MoveX      float64 `json:"moveX"`
	MoveY      float64 `json:"moveY"`
	Speed      float64 `json:"speed"`
	ClientTime int64   `json:"clientTime"`
}

// ServerConfig represents the game configuration
type ServerConfig struct {
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
}

// ClientMessage represents incoming messages from clients
type ClientMessage struct {
	UID        string  `json:"uid"`
	MoveX      float64 `json:"moveX"`
	MoveY      float64 `json:"moveY"`
	ClientTime int64   `json:"clientTime"`
}

// WebSocketServer manages WebSocket connections and broadcasts messages
type WebSocketServer struct {
	Clients map[*websocket.Conn]bool
	Broadcast chan map[string]interface{}
	Mutex sync.Mutex
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients: make(map[*websocket.Conn]bool),
		Broadcast: make(chan map[string]interface{}),
	}
}

func (wss *WebSocketServer) AddClient(conn *websocket.Conn) {
	wss.Mutex.Lock()
	defer wss.Mutex.Unlock()
	wss.Clients[conn] = true
}

func (wss *WebSocketServer) RemoveClient(conn *websocket.Conn) {
	wss.Mutex.Lock()
	defer wss.Mutex.Unlock()
	delete(wss.Clients, conn)
	conn.Close()
}

func (wss *WebSocketServer) BroadcastMessage(msg map[string]interface{}) {
	wss.Mutex.Lock()
	defer wss.Mutex.Unlock()

	// if value, err := json.Marshal(msg); err == nil {
	// 	 log.Println("Broadcasting message:", string(value))
	// }

	for client := range wss.Clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			delete(wss.Clients, client)
		}
	}
}

// ServerState holds the current state of the game server
type ServerState struct {
	Config  ServerConfig                `json:"config"`
	Players map[string]*Player          `json:"players"`
	mu      sync.Mutex
}

var state = ServerState{
	Config: ServerConfig{
		MaxX: 1000,
		MaxY: 1000,
	},
	Players: make(map[string]*Player),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	wss := NewWebSocketServer()
	go handleBroadcasts(wss)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(wss, w, r)
	})

	log.Println("WebSocket server started on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(wss *WebSocketServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer wss.RemoveClient(conn)

	wss.AddClient(conn)

	uid := uuid.New().String()
	player := &Player{
		X:      400,
		Y:      300,
		Radius: 20,
		Color:  randomColor(),
		Speed:  20,
	}
	state.mu.Lock()
	state.Players[uid] = player
	state.mu.Unlock()

	log.Println("Added player", uid)

	message := map[string]interface{}{
		"action":  "join",
		"uid":     uid,
		"config":  state.Config,
		"players": state.Players,
	}

	if err := conn.WriteJSON(message); err != nil {
		return
	} 	

	for {
		var msg ClientMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			state.mu.Lock()
			delete(state.Players, uid)
			state.mu.Unlock()
			
			var message map[string]interface{} = map[string]interface{}{
				"action": "disconnect",
				"player": uid,
			}

			log.Println("Removed player", uid)
			wss.BroadcastMessage(message)
			break
		} 

		state.mu.Lock()
		player := state.Players[uid]
		if player != nil {
			player.MoveX = msg.MoveX
			player.MoveY = msg.MoveY
			player.ClientTime = msg.ClientTime
			player.X += player.MoveX * player.Speed
			player.Y += player.MoveY * player.Speed
			player.X = clamp(player.X, 0, state.Config.MaxX)
			player.Y = clamp(player.Y, 0, state.Config.MaxY)
		}
		state.mu.Unlock()

		var message map[string]interface{} = map[string]interface{}{
			"action":  "sync",
			"players": state.Players,
		}

		wss.BroadcastMessage(message)
	}
}

func handleBroadcasts(wss *WebSocketServer) {
	for msg := range wss.Broadcast {
		wss.BroadcastMessage(msg)
	}
}

func randomColor() string {
	rand.Seed(time.Now().UnixNano())
	red := 128 + rand.Intn(128)
	green := 128 + rand.Intn(128)
	blue := 128 + rand.Intn(128)

	// Return as a hex color string in #RRGGBB format
	return fmt.Sprintf("#%02X%02X%02X", red, green, blue)
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

