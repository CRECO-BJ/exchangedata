package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"

	"github.com/exchangedata/exchanger"
)

const ()

var (
	exVar = []exchanger.Exchanger{
		{Name: "OKEX",
			WssURL:    url.URL{Scheme: "wss", Path: "real.okex.com:10441/websocket"},
			WebAPIURL: url.URL{Scheme: "https", Path: "www.okex.com/docs/en/"},
			Symbols:   []exchanger.Symbol{{Base: "btc", Quote: "usdt"}, {Base: "eth", Quote: "usdt"}, {Base: "eth", Quote: "btc"}}},
		{Name: "Poloniex",
			WssURL:    url.URL{Scheme: "wss", Path: "api2.poloniex.com"},
			WebAPIURL: url.URL{Scheme: "https", Path: "poloniex.com"},
			Symbols:   []exchanger.Symbol{{Base: "btc", Quote: "usdt"}, {Base: "eth", Quote: "usdt"}, {Base: "eth", Quote: "btc"}}},
	}
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wg := &sync.WaitGroup{}
	for k := range exVar {
		wg.Add(1)
		go exVar[k].Run(wg)
		go func(e *exchanger.Exchanger) { // routinue to clean close the connections
			wg.Wait()
			e.Close()
		}(&exVar[k])
	}

	<-interrupt // exit only when the application is interrupted
	log.Println("interrupt")
	for _, ex := range exVar {
		ex.CloseDone()
	}
	wg.Wait()
	return
}
