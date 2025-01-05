package main

import (
	"log"
	"net/http"
)

func main() {
	wss := NewWebSocketServer()
	go handleBroadcasts(wss)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(wss, w, r)
	})

	log.Println("WebSocket server started on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
