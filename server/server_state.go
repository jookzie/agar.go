package main

import (
	"sync"
	
	"github.com/gorilla/websocket"
)

type Player struct {
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Radius     float64 `json:"radius"`
	Color      string  `json:"color"`
	MoveX      float64 `json:"moveX"`
	MoveY      float64 `json:"moveY"`
	Speed      float64 `json:"speed"`
	ClientTime int64   `json:"clientTime"`
	Connection *websocket.Conn
}

type ServerConfig struct {
	MaxX        int `json:"maxX"`
	MaxY        int `json:"maxY"`
	FeedMapSize int `json:"feedmapsize"`
}

type ServerState struct {
	Config  ServerConfig       `json:"config"`
	Players map[string]*Player `json:"players"`
	FeedMap [][2]int           `json:"feedmap"`
	mu      sync.Mutex
}

type ClientMessage struct {
	UID        string  `json:"uid"`
	MoveX      float64 `json:"moveX"`
	MoveY      float64 `json:"moveY"`
	ClientTime int64   `json:"clientTime"`
}

var config ServerConfig = ServerConfig{
	MaxX: 5000,
	MaxY: 5000,
	FeedMapSize: 1000,
}

var state = ServerState{
	Config: config,
	Players: make(map[string]*Player),
	FeedMap: GeneratePoints(config.FeedMapSize, int(config.MaxX), int(config.MaxY)),
}

