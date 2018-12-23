package main

import (
	"fmt"

	"github.com/exchangedata/exchanger"
	"github.com/exchangedata/exchanger/bittrex"
)

var exchangerSupported []string

func init() {
	exchangerSupported = []string{}
	exchangerSupported = append(exchangerSupported, "bittrex")
}

func isSupported(s string) bool {
	if s == "" {
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
func NewExchanger(s string) (exchanger.ExControl, error) {
	if isSupported(s) == false {
		return nil, fmt.Errorf("not supported exchanger %s", s)
	}

	var nex exchanger.ExControl
	switch s {
	case "bittrex":
		nex = bittrex.NewBittrex()
	}
	return nex, nil
}
