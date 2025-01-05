package main

import (
	"log"
	"net/http"
	"math/rand"
	// "encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
		X:      rand.Float64() * state.Config.MaxX,
		Y:      rand.Float64() * state.Config.MaxY,
		Radius: 20,
		Color:  RandomColor(),
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
		"feedmap": state.FeedMap,
	}

	// if msgJson, err := json.Marshal(message); err == nil {
	// 	log.Println(msgJson)
	// }

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

			log.Println("Removed player", uid)
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
			player.X = Clamp(player.X, 0, state.Config.MaxX)
			player.Y = Clamp(player.Y, 0, state.Config.MaxY)
		}

		var eatenPoints [][2]float64 = FilterPointsInCircle(player.X, player.Y, player.Radius, state.FeedMap)
		if len(eatenPoints) != 0 {
			state.FeedMap = SubtractArrays(state.FeedMap, eatenPoints);
			player.Radius += float64(len(eatenPoints));
			state.FeedMap = append(state.FeedMap, GeneratePoints(len(eatenPoints), state.Config.MaxX, state.Config.MaxY)...)
		}

		state.mu.Unlock()

		var message map[string]interface{}

		if len(eatenPoints) == 0 {
			message = map[string]interface{}{
				"action":  "sync",
				"players": state.Players,
			}
		} else {
			message = map[string]interface{}{
				"action":  "sync",
				"players": state.Players,
				"feedmap": state.FeedMap,
			}
		}

		wss.BroadcastMessage(message)
	}
}

func handleBroadcasts(wss *WebSocketServer) {
	for msg := range wss.Broadcast {
		wss.BroadcastMessage(msg)
	}
}
