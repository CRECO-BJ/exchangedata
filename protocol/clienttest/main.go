package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
)

const (
	addr      = "real.okex.com:10441"
	pingEvent = "{'event':'ping'}"
	pongEvent = "{'event':'pong'}"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/websocket"}
	log.Printf("connecting to %s", u.String())

	netDialer, err := proxy.SOCKS5("tcp", "localhost:1080", nil, proxy.Direct)
	if err != nil {
		log.Fatalf("sock5 configuration error %v", err)
	}
	dialer := websocket.Dialer{NetDial: netDialer.Dial}
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("dial:%s error:%v", u.String(), err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(pingEvent))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
