package common

import (
	"fmt"
	"strings"
	"time"
)

type Currency struct {
	Name       string       `gorm:"primary_key,not null"`
	Abbr       string       `gorm:"size:16"`
	Exchangers []*Exchanger `gorm:"many2many:currency_exchangers;"`
}

type Exchanger struct {
	Name        string           `gorm:"primary_key,not null"`
	Currencies  []*Currency      `gorm:"many2many:currency_exchangers;"`
	Markets     []*Market        `gorm:"foreignkey:ExRef"`
	CommService CommunicationAPI `gorm:"foreignkey:ExRef"`
}

type Market struct {
	Name        string `gorm:"primary_key, not null"`
	Symbol      Symbol `gorm:"foreignkye:SymbolID"`
	SymbolID    int
	Exchanger   *Exchanger `gorm:"foreignkey:ExchangerID"`
	ExchangerID int
	Active      bool
	info        string
	precision   [3]int
	limits      [3]Limitation
}

type CommunicationAPI struct {
	Version    string
	WebURL     string
	WssURL     string
	Enable     bool
	ExRef      int
	KeySecrect AccessSecret `gorm:"foreignkey:ComAPIRef"`
}

type AccessSecret struct {
	ComAPIRef int
	APIKey    string
	Secret    string
}

type Symbol struct {
	ID      int
	Base    Currency `gorm:"foreignkey:BaseID, not null"`
	BaseID  int
	Quote   Currency `gorm:"foreignkey:QuoteID"`
	QuoteID int
}

type Ticker struct {
	symbol        string
	Info          string
	Time          time.Time `gorm:"primary_key"`
	MarketRef     int       `gorm:"primary_key"`
	High          float64
	Low           float64
	Bid           float64
	BidVolume     float64
	Ask           float64
	AskVolume     float64
	Last          float64
	PreviousClose float64
	Change        float64
	Percentage    float64
	Average       float64
	BaseVolume    float64
	QuoteVolume   float64
	Open          float64
	Close         float64
	Market        *Market `gorm:"foreignkey:MarketRef"`
}

type Trade struct {
	symbol    string
	Time      time.Time
	MarketRef int   `gorm:"primary_key"`
	orderID   int64 `gorm:"primary_key"`
	Type      string
	Side      string
	Price     float64
	Amount    float64
	Total     float64
	Market    *Market `gorm:"foreignkey:MarketRef"`
}

type OrderBook struct {
	symbol    string
	Time      time.Time
	MarketRef int `gorm:"primary_key"`
	bids      []PriceVol
	asks      []PriceVol
	time      time.Time
	Market    *Market `gorm:"foreignkey:MarketRef"`
}

type PriceVol struct {
	price  float64
	volume float64
	occurs []*OrderBook
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
