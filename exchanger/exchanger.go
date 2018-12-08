package exchanger

import (
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
)

// Symbol ...
type Symbol struct {
	Base  string
	Quote string
}

// Exchanger ...
type Exchanger struct {
	ID            string
	Name          string
	Countries     []string
	WebAPIURL     url.URL
	WebAPIVersion string
	WssURL        url.URL
	WssVersion    string
	WebAPIs       []string
	WssAPIs       []string
	Timeout       int
	RateLimit     int
	UserAgetn     string
	Verbose       bool
	Symbols       []Symbol
	Proxy         url.URL

	Conn *websocket.Conn
	Done chan struct{} // Closed when the receive rountine received error, then the main exchanger communication routine exit
	// If CloseDone is not closed, the connection should be reconnected...ToDo
	CloseDone chan struct{} // Signal to close connection and exit. Program exiting...
}

// UseWss ...
func (e *Exchanger) UseWss() bool {
	if e.WssURL.String() == "" {
		return false
	}

	return true
}

// UseWeb ...
func (e *Exchanger) UseWeb() bool {
	if e.WebAPIURL.String() == "" {
		return false
	}

	return true
}

// Dial ...
func (e *Exchanger) Dial() *websocket.Conn {
	var dialer *websocket.Dialer
	if isValidProxy(e.Proxy) {
		//	auth = proxy.Auth(User:, Password:kb109kb109)
		netDialer, err := proxy.SOCKS5("udp", e.Proxy.String(), nil, proxy.Direct)
		if err != nil {
			log.Fatalf("sock5 configuration error %v", err)
		}
		dialer = &websocket.Dialer{NetDial: netDialer.Dial}
	} else {
		dialer = websocket.DefaultDialer
	}
	c, _, err := dialer.Dial(e.WssURL.String(), nil)
	if err != nil {
		log.Fatalf("dial:%s error:%v", e.WssURL.String(), err)
	}
	e.Conn = c
	return e.Conn
}

// Close ...
func (e *Exchanger) Close() {
	e.Conn.Close()
}

// Run ...
func (e *Exchanger) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(e.CloseDone)

	go func() {
		defer close(e.Done)
		for {
			_, message, err := e.Conn.ReadMessage()
			if err != nil { // if read error, rountine exit, redail
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
		case <-e.Done:
			return
		case <-ticker.C:
			err := e.Conn.WriteMessage(websocket.TextMessage, []byte("test"))
			if err != nil {
				log.Println("write:", err)
				return
			}
			log.Println("write:", "test")
		case <-e.CloseDone:
			log.Println("Close connection and to exit")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := e.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-e.Done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func isValidProxy(url url.URL) bool {
	return len(url.Hostname()) != 0
}

type limitation struct {
	min, max int
}

// Market ...
type Market struct {
	ID string
	Symbol
	Active    bool
	info      string
	precision [3]int
	limits    [3]limitation
}

type priceVol struct {
	price  float64
	volume float64
}

type orderBook struct {
	bids []priceVol
	asks []priceVol
	time time.Time
}

type ticker struct {
	symbol                                               Symbol
	info                                                 string
	time                                                 time.Time
	high                                                 float64
	low                                                  float64
	bid                                                  float64
	bidVolume                                            float64
	ask, askVolume                                       float64
	open, close                                          float64
	last, previousClose                                  float64
	change, percentage, average, baseVolume, quoteVolume float64
}

type trade struct {
	info    string
	ID      string
	time    time.Time
	symbol  Symbol
	orderID string
	Type    string
	side    string
	price   float64
	amount  float64
}
