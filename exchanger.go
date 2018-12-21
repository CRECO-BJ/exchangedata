package main

import (
	"log"

	"github.com/exchangedata/exchanger"
	"github.com/exchangedata/exchanger/okex"
	"github.com/exchangedata/exchanger/poloniex"
	"github.com/exchangedata/exchanger/bittrex"
)

var exchangerSupported []string

func init() {
	exchangerSupported = {"okex", "poloniex", "bittrex"}
}

func isSupported(s string) bool {
	if s=="" {
		return false
	}
	for _, v := range exchangerSupported {
		if v == s {
			return true
		}
	}
	return false
}

// NewExchanger creates a exchanger.Exchanger with the specified name s
func NewExchanger(s string) ( *exchanger.ExControl, error) {
	if isSupported(s)==false {
		return nil, fmt.Errorf("not supported exchanger %s", s)
	}

	var nex *exchagner.ExControl
	switch s {
	case "okex": nex = okex.NewOKEX()
	case "polonex": nex = poloniex.NewPoloniex()
	case "bittrex": nex = bittrex.NewBittrex()
	}
	return nex, nil
}
