package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/exchangedata/common"
	"github.com/exchangedata/exchanger"
)

const ()

var ( // all lowercase string
	exVar = []common.Exchanger{
		{Name: "bittrex"},
	}
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	exs := []exchanger.ExControl{}
	wg := &sync.WaitGroup{}
	for k := range exVar {
		ex, err := NewExchanger(exVar[k].Name)
		if err != nil {
			log.Fatalf("cannot initialize exchanger, configuration error, %s", exVar[k].Name)
		} else {
			log.Println("Exchanger ", exVar[k].Name, "initialized")
		}
		exs = append(exs, ex)
		ex.Setup()
		wg.Add(1)
		go ex.Start(wg)
	}

	log.Println("All exchangers have been set, waiting for interrupt to terminate exchangedata...")
	<-interrupt // exit only when the application is interrupted
	log.Println("intterupted! Exiting...")
	for _, ex := range exs {
		ex.Stop()
	}
	wg.Wait()
	return
}
