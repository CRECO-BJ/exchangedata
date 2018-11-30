package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

const (
	serverURI    = "wss://real.okex.com:10441/websocket"
	serverOrigin = "http://real.okex.com/"
	pingEvent    = "{'event':'ping'}"
	pongEvent    = "{'event':'pong'}"
)

func main() {
	ws, err := websocket.Dial(serverURI, "", serverOrigin)
	if err != nil {
		log.Println("Open websocket error", err)
	}
	defer ws.Close()

	if _, err := ws.Write([]byte(pingEvent)); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
