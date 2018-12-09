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

	conn *websocket.Conn
	done chan struct{} // Closed when the receive rountine received error, then the main exchanger communication routine exit
	// If CloseDone is not closed, the connection should be reconnected...ToDo
	closeDone chan struct{} // Signal to close connection and exit. Program exiting...
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
	e.conn = c
	return e.conn
}

// Close ...
func (e *Exchanger) Close() {
	e.conn.Close()
}

// CloseDone ...
func (e *Exchanger) CloseDone() {
	e.closeDone <- struct{}{}
}

// Done ...
func (e *Exchanger) Done() {
	e.done <- struct{}{}
}

// Run ...
func (e *Exchanger) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	e.done = make(chan struct{})
	e.closeDone = make(chan struct{})
	defer close(e.closeDone)
	defer close(e.done)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	exit := false

start:
	if e.UseWss() {
		e.Dial()
	} else if !e.UseWeb() { // no communication url is defined, bypass
		log.Println("no valid communication method", e)
		return
	}

	go func() {
		defer e.Done()
		for {
			_, message, err := e.conn.ReadMessage()
			if err != nil { // if read error, rountine exit, redail
				log.Println(e.Name, "%s read:", err)
				return
			}
			var b []byte
			log.Printf("%s recv: %s", e.Name, string(message))
			for i, x := range message {
				if i%16 == 0 {
					b = append(b, '\n')
				}
				b = strconv.AppendInt(b, int64(x), 16)
				b = append(b, ' ')
			}
			log.Printf("%s recv(hex): %s", e.Name, string(b))
		}
	}()

	for {
		select {
		case <-e.done:
			if exit == true {
				return
			}
			goto start
		case <-ticker.C: // timely keepAlive processing
			//			err := e.conn.WriteMessage(websocket.TextMessage, []byte("test"))
			//			if err != nil {
			//				log.Println(e.Name, "write:", err)
			//				return
			//			}
			log.Println(e.Name, "time to write:", "test")
		case <-e.closeDone:
			log.Println("Close connection and to exit")
			exit = true
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := e.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-e.done:
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
