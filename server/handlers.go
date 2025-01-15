package main

import (
	"log"
	"math"
	"math/rand"
	"net/http"
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
		Name:	nicknames[rand.Intn(len(nicknames))],
		X:      rand.Float64() * float64(state.Config.MaxX),
		Y:      rand.Float64() * float64(state.Config.MaxY),
		Radius: 20,
		Color:  RandomColor(),
		Speed:  2,
		Connection: conn,
	}
	state.mu.Lock()
	state.Players[uid] = player
	state.mu.Unlock()

	log.Println("Added player", player.Name)

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

			log.Println("Removed player", player.Name)
			break
		}

		message := map[string]interface{}{
			"action":  "sync",
			"uid": uid,
		}

		state.mu.Lock()
		player := state.Players[uid]

		if player != nil {
			player.MoveX = msg.MoveX
			player.MoveY = msg.MoveY
			player.ClientTime = msg.ClientTime
			player.X += player.MoveX * player.Speed
			player.Y += player.MoveY * player.Speed
			player.X = Clamp(player.X, 0, float64(state.Config.MaxX))
			player.Y = Clamp(player.Y, 0, float64(state.Config.MaxY))


			for idx, point := range state.FeedMap {
				if IsPointInCircle(player.X, player.Y, player.Radius, point[0], point[1]) {
					player.Radius += 0.5
					player.Speed = 30 / player.Radius + 0.5
					state.FeedMap = append(state.FeedMap[:idx], state.FeedMap[idx+1:]...)
					message["eatenPoint"] = point

					addSize := math.Max(0.0, float64(state.Config.FeedMapSize - len(state.FeedMap)))
					addedPoints := GeneratePoints(int(addSize), state.Config.MaxX, state.Config.MaxY)
					state.FeedMap = append(state.FeedMap, addedPoints...)
					message["addedPoints"] = addedPoints
					break
				}
			}

			for otherUid, other := range state.Players {
				if otherUid != uid && IsCircleInCircle(player.X, player.Y, player.Radius+10, other.X, other.Y, other.Radius) {
					player.Radius += state.Players[otherUid].Radius - 20
					player.Speed = 30 / player.Radius + 0.5
					wss.RemoveClient(state.Players[otherUid].Connection)
					delete(state.Players, otherUid)
					message["eatenPlayer"] = otherUid
					break
				}
			}

			message["player"] = player;
		}
		state.mu.Unlock()

		wss.BroadcastMessage(message)
	}
}

func handleBroadcasts(wss *WebSocketServer) {
	for msg := range wss.Broadcast {
		wss.BroadcastMessage(msg)
	}
}
