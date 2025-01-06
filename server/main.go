package main

import (
	"log"
	"net/http"
)

func main() {
	wss := NewWebSocketServer()
	go handleBroadcasts(wss)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handleConnections(wss, w, r)
	})

	log.Println("WebSocket server started on http://0.0.0.0:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

