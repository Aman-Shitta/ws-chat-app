package main

import (
	"net/http"

	"github.com/Aman-Shitta/ws-chat-app/server"
	"golang.org/x/net/websocket"
)

func main() {
	s := server.NewServer()

	http.Handle("/ws", websocket.Handler(s.HandleConnection))
	http.ListenAndServe(":3000", nil)

}
