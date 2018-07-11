package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 5 * time.Second
	// Time allowed to read the next pong message from the client.
	pongWait = 10 * time.Second
	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = 5 * time.Second
	// Poll cache for changes with this period.
	refreshPeriod = 5 * time.Second
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func wsStatsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	go writer(ws, CacheGet())
	reader(ws)
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	// read and discard/drain input
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, data interface{}) {
	pingTicker := time.NewTicker(pingPeriod)
	refreshTicker := time.NewTicker(refreshPeriod)
	defer func() {
		pingTicker.Stop()
		refreshTicker.Stop()
		ws.Close()
	}()
	// every refresh period, send updated stats to client
	for {
		select {
		case <-refreshTicker.C: // use go channels to signal when ready
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			stats, err := json.Marshal(CacheGet())
			if err != nil {
				log.Println(err)
				ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
				return
			}
			if err := ws.WriteMessage(websocket.TextMessage, stats); err != nil {
				log.Println(err)
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
