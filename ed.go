package main

import (
	"log"

	"golang.org/x/net/websocket"
)

const (
	serverURI    = "wss://real.okex.com:10441/websocket"
	serverOrigin = "http://real.okex.com:10441/"
)

var ()

func main() {
	ws, err := websocket.Dial(serverURI, "", serverOrigin)
	if err != nil {
		log.Println("Open websocket error", err)
	}
	defer ws.Close()

}
