package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins (you can tighten this later)
	},
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WS upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Println("✅ WebSocket connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ WS read error:", err)
			break
		}
		log.Printf("📩 Received: %s\n", msg)
		conn.WriteMessage(websocket.TextMessage, []byte("Echo: "+string(msg)))
	}
}
