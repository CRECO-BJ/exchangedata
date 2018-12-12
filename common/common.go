package common

import (
	"fmt"
	"strings"
	"time"
)

// Symbol ...
type Symbol struct {
	Base  string
	Quote string
}

func (s Symbol) ParseString(str string) (Symbol, error) {
	ss := strings.Split(str, "_")
	if len(ss) != 2 {
		return Symbol{}, fmt.Errorf("not a valid traded pair")
	}
	return Symbol{Base: ss[0], Quote: ss[1]}, nil
}

func (s Symbol) String() string {
	return s.Base + "_" + s.Quote
}

func (s Symbol) MarshallJSON() ([]byte, error) {
	str := s.String()

	return []byte(str), nil
}

type Limitation struct {
	min, max int
}

// Market ...
type Market struct {
	ID string
	Symbol
	Active    bool
	info      string
	precision [3]int
	limits    [3]Limitation
}

type PriceVol struct {
	price  float64
	volume float64
}

type OrderBook struct {
	bids []PriceVol
	asks []PriceVol
	time time.Time
}

type Ticker struct {
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

type Trade struct {
	Time        time.Time
	ExchangerID int64
	Symbol      Symbol
	orderID     int64
	Type        string
	Side        string
	Price       float64
	Amount      float64
}
