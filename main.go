package main

import (
	"log"

	"github.com/gorilla/websocket"
)

const (
	serverURI    = "wss://real.okex.com:10441/websocket"
	serverOrigin = "http://real.okex.com:10441/"
)

var (
)

func ExchangerComm() {
	for ex := range exchangers {
		if ex.useWss() {
			ex.Dial()
		} else !ex.useWeb() {
			return
		} 
		go ex.ReceiveFunc()
		go ex.HandleFunc()
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := addr["okex"]
	log.Printf("connecting to %s", u.String())

	var dialer *websocket.Dialer
	if viaProxy {
		//	auth = proxy.Auth(User:, Password:kb109kb109)
		netDialer, err := proxy.SOCKS5("udp", "localhost:1080", nil, proxy.Direct)
		if err != nil {
			log.Fatalf("sock5 configuration error %v", err)
		}
		dialer = &websocket.Dialer{NetDial: netDialer.Dial}
	} else {
		dialer = websocket.DefaultDialer
	}
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
			var b []byte
			log.Printf("recv: %s", string(message))
			for i, x := range message {
				if i%16 == 0 {
					b = append(b, '\n')
				}
				b = strconv.AppendInt(b, int64(x), 16)
				b = append(b, ' ')
			}
			log.Printf("recv(hex): %s", string(b))
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
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
			log.Println("write:", pingEvent)
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
