package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan map[string]interface{}
	Mutex     sync.Mutex
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients:   make(map[*websocket.Conn]bool),
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

	for client := range wss.Clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			delete(wss.Clients, client)
		}
	}
}

