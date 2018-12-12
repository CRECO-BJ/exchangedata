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

var ( // all lowercase string
	exVar = []exchanger.Exchanger{
		{Name: "okex",
			WssURL:    url.URL{Scheme: "wss", Path: "real.okex.com:10441/websocket"},
			WebAPIURL: url.URL{Scheme: "https", Path: "www.okex.com/docs/en/"},
			Symbols:   []exchanger.Symbol{{Base: "btc", Quote: "usdt"}, {Base: "eth", Quote: "usdt"}, {Base: "eth", Quote: "btc"}}},
		{Name: "poloniex",
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
		ex, err := NewExchanger(exVar[k])
		if err != nil {
			log.Fatalf("cannot initialize exchanger, configuration error, %s", exVar[k].Name)
		}
		ex.Setup(exVar[k])
		wg.Add(1)
		if ex.Start(wg) != nil {
			log.Fatalf("cannot start exchanger, %s", exVar[k].Name)
		}
	}

	<-interrupt // exit only when the application is interrupted
	for _, ex := range exVar {
		ex.Stop()
	}
	wg.Wait()
	return
}
