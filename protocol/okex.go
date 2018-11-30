package protocol

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

// OkexParams ...
type OkexParams struct {
	apiKey     string
	sign       string
	timestamp  time.Time
	passphrase string
}

func (op *OkexParams) String() string {
	return fmt.Sprintf("{'api_key':'%s','sign':'%s','timestamp':'%s','passphrase':'%s'}",
		op.apiKey, op.sign, op.timestamp.String(), op.passphrase)
}

// OkexRequest ...
type OkexRequest struct {
	event      string
	channel    string
	parameters OkexParams
}

func (or *OkexRequest) String() string {
	return fmt.Sprintf("{'event':'%s','channel':'%s','parameters':'%s'}", or.event, or.channel, or.parameters)
}

// OkexResponse ...
type OkexResponse struct {
	channel   string
	success   bool
	errorcode int
	data      interface{}
}

var heartBeatInterval = 30 * time.Second

type WsReceiver interface {
	SubscribeChannel(string) error
	EventHandle(interface{}) error
}

// OkexWsReceiver ...
type OkexWsReceiver struct {
	conn *websocket.Conn

	wSync    sync.Mutex
	channels []string

	exit        chan struct{}
	restartChan []chan struct{}
}

const (
	pingEvent = "{'event':'ping'}"
	pongEvent = "{'event':'pong'}"
)

func (owr *OkexWsReceiver) send(data []byte) (int, error) {
	owr.wSync.Lock()
	owr.wSync.Unlock()
	return owr.conn.Write(data)
}

// dev: heartBeat sends a ping message every a interval, exit if the protocol is closed
// Todo: Does it is better to safely close this routine before the protocol?
func (owr *OkexWsReceiver) heartBeat(exit <-chan struct{}) {
	t := time.NewTimer(heartBeatInterval)
	for {
		select {
		case <-t.C: // heart beat timeout
			_, err := owr.send([]byte(pingEvent))
			if err != nil {
				log.Fatal("cannot send out heartbeat")
			} // just make sure a ping message is sent, the related pong message is processed in the receiver thread
			t.Reset(heartBeatInterval)
		case <-exit:
			t.Stop()
			return
		}
	}
}

const recvBufferLength = 4096

func (owr *OkexWsReceiver) receiver(exit <-chan struct{}) {
	msg := make([]byte, recvBufferLength)
	for {
		select {
		case owr.conn.Read(msg):

			handleReceivedMessage(owr, msg)
		case <-exit:
			return
		}
	}
}
