package main

import (
	"sync"
	"net/url"
	"os"
	"os/signal"
	"log"

	"github.com/gorilla/websocket"
	"github.com/exchangedata/exchanger"
)

const (
)

var (
	exVar = []exchanger.Exchanger{
		{ Name:"OKEX", 
		WssURL:url.URL{Scheme:"wss",Path:"real.okex.com:10441/websocket"}, 
		WebAPIURL:url.URL{Scheme:"https",Path:"www.okex.com/docs/en/"}, 
		Symbols:[]exchanger.Symbol{{Base:"btc",Quote:"usdt"},{Base:"eth",Quote:"usdt"},{Base:"eth",Quote:"btc"}}},
		{ Name:"Poloniex", 
		WssURL:url.URL{Scheme:"wss", Path:"api2.poloniex.com"}, 
		WebAPIURL:url.URL{Scheme:"https",Path:"poloniex.com"},
		Symbols:[]exchanger.Symbol{{Base:"btc",Quote:"usdt"},{Base:"eth",Quote:"usdt"},{Base:"eth",Quote:"btc"}}},
	}
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wg := &sync.WaitGroup{}
	for k := range exVar {
		if exVar[k].UseWss() {
			exVar[k].Dial()
		} else if !exVar[k].UseWeb() {
			continue
		}
		exVar[k].Done = make(chan struct{})
		exVar[k].CloseDone = make(chan struct{})
		wg.Add(1)
		go exVar[k].Run(wg)
		go func () {	// routinue to clean close the connections
			wg.Wait()
			exVar[k].Close()
		}
	}

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")

			for _, ex := range exVar {
				err := ex.CloseDone()
				if err != nil {
					log.Println("write close:", err)
					return
				}
			}
			wg.Wait()
			return
		}
	}	
}
