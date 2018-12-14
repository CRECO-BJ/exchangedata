package common

import (
	"fmt"
	"strings"
	"time"
)

type Currency struct {
	ID         string       `gorm:"primary_key;size:32"`
	Name       string       `gorm:"unique;size:64"`
	Abbr       string       `gorm:"size:16"`
	Exchangers []*Exchanger `gorm:"many2many:currency_exchangers;"`
}

type Exchanger struct {
	ID          string            `gorm:"primary_key;size:32"`
	Name        string            `gorm:"unique;not null"`
	Currencies  []*Currency       `gorm:"many2many:currency_exchangers;"`
	Markets     []*Market         `gorm:"foreignkey:ExRef"`
	CommService *CommunicationAPI `gorm:"foreignkey:ExRef"`
}

type Market struct {
	ID         string `gorm:"primary_key;size:32"`
	Name       string `gorm:"not null"`
	Symbol     Symbol `gorm:"foreignkey:SymRef"`
	SymRef     int    `gorm:"unique_index:idx_sym_ex"`
	Active     bool
	Info       string
	Precision  int        `gorm:"default:8"`
	Limitation Limitation `gorm:"embedded;embedded_prefix:amount_"`
	MinStep    float64
	ExRef      int `gorm:"unique_index:idx_sym_ex"`
}

type Limitation struct {
	Min, Max int
}

type CommunicationAPI struct {
	ID           string `gorm:"primary_key;size:32"`
	Version      string `gorm:"size:32"`
	WebURL       string
	WssURL       string
	Enable       bool
	ExRef        int
	Timeout      int
	RateLimit    int
	AccessSecret AccessSecret `gorm:"embedded;embedded_prefix:comm_"`
}

type AccessSecret struct {
	ID        string `gorm:"primary_key;size:32"`
	ComAPIRef int
	APIKey    string `gorm:"-"`
	Secret    string `gorm:"-"`
	FilePath  string
	Salt      int
}

type Symbol struct {
	ID      string   `gorm:"primary_key;size:32"`
	Base    Currency `gorm:"foreignkey:BaseID;not null"`
	BaseID  int
	Quote   Currency `gorm:"foreignkey:QuoteID"`
	QuoteID int
}

type Ticker struct {
	ID            string `gorm:"primary_key;size:32"`
	Info          string
	Time          time.Time `gorm:"unique_index:idx_time_market;not null"`
	MarketRef     int       `gorm:"unique_index:idx_time_market"`
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
	ID        string    `gorm:"primary_key;size:32"`
	Time      time.Time `gorm:"not null"`
	MarketRef int       `gorm:"unique_index:idx_market_trade;not null"`
	OrderID   string    `gorm:"unique_index:idx_market_trade;not null"`
	Type      string
	Side      string
	Price     float64
	Amount    float64
	Total     float64
	Market    *Market `gorm:"foreignkey:MarketRef"`
}

type OrderBook struct {
	ID        string     `gorm:"primary_key;size:32"`
	Time      time.Time  `gorm:"unique_index:idx_market_orderbook;not null"`
	MarketRef int        `gorm:"unique_index:idx_market_orderbook"`
	Bids      []PriceVol `gorm:"many2many:bid_pricevols;"`
	Asks      []PriceVol `gorm:"many2many:ask_pricevols;"`
	Market    *Market    `gorm:"foreignkey:MarketRef"`
}

type PriceVol struct {
	ID     string `gorm:"primary_key;size:32"`
	Price  float64
	Volume float64
	Occurs []*OrderBook `gorm:"many2many:bid_pricevols;many2many:ask_pricevols;"`
}

func (s Symbol) ParseString(str string) (Symbol, error) {
	ss := strings.Split(str, "_")
	if len(ss) != 2 {
		return Symbol{}, fmt.Errorf("not a valid traded pair")
	}
	return Symbol{Base: Currency{Abbr: ss[0]}, Quote: Currency{Abbr: ss[1]}}, nil
}

func (s Symbol) String() string {
	return s.Base.Abbr + "_" + s.Quote.Abbr
}

func (s Symbol) MarshallJSON() ([]byte, error) {
	str := s.String()

	return []byte(str), nil
}
