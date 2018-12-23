package common

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// @Dev: all the Name should be lowercase, Abbr should be uppercase, the leading and tail space should be trimmed
// @User: sometimes exchanger api has mismatching currency and market setting. Such mismatching can be addressed by predefined setting in configuration or database.

type Currency struct {
	ID         uint         `gorm:"primary_key"`
	Name       string       `gorm:"unique;size:64;not null"`
	Abbr       string       `gorm:"size:16"`
	AbbrFinal  bool         `gorm:"default:false"` // if set, Abbr cannot be modified by a network response message.
	Exchangers []*Exchanger `gorm:"many2many:currency_exchangers"`
	Info       string
}

type Exchanger struct {
	ID          uint              `gorm:"primary_key"`
	Name        string            `gorm:"unique;not null"`
	Currencies  []*Currency       `gorm:"many2many:currency_exchangers"`
	Markets     []*Market         `gorm:"foreignkey:ExRef"`
	CommService *CommunicationAPI `gorm:"foreignkey:ExRef"`
	Info        string
}

type Market struct {
	ID         uint    `gorm:"primary_key"`
	Name       string  `gorm:"not null"`
	Symbol     *Symbol `gorm:"foreignkey:SymRef"`
	SymRef     uint    `gorm:"unique_index:idx_sym_ex"`
	Active     bool    `gorm:"default:true"`
	Info       string
	Precision  uint       `gorm:"default:8"`
	Limitation Limitation `gorm:"embedded;embedded_prefix:amount_"`
	MinStep    float64
	ExRef      uint       `gorm:"unique_index:idx_sym_ex"`
	Exchanger  *Exchanger `gorm:"foreignkey:ExRef"`
}

type Limitation struct {
	Min, Max float64
}

type CommunicationAPI struct {
	ID           uint   `gorm:"primary_key"`
	Version      string `gorm:"size:32"`
	WebURL       string `gorm:"size:32;not null"`
	WssURL       string
	Enable       bool
	ExRef        uint
	Timeout      int
	RateLimit    int
	AccessSecret AccessSecret `gorm:"embedded;embedded_prefix:comm_"`
}

type AccessSecret struct {
	ComAPIRef uint
	APIKey    string `gorm:"-"`
	Secret    string `gorm:"-"`
	FilePath  string
	Salt      int
}

type Symbol struct {
	ID      uint      `gorm:"primary_key"`
	Base    *Currency `gorm:"foreignkey:BaseID;not null"`
	BaseID  uint      `gorm:"unique_index:idx_base_quote"`
	Quote   *Currency `gorm:"foreignkey:QuoteID"`
	QuoteID uint      `gorm:"unique_index:idx_base_quote"`
}

type Ticker struct {
	ID            uint      `gorm:"primary_key"`
	Time          time.Time `gorm:"unique_index:idx_time_market;not null"`
	MarketRef     uint      `gorm:"unique_index:idx_time_market;not null"`
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
	Market        *Market `gorm:"foreignkey:MarketRef;association_save_reference:false"`
}

type Trade struct {
	ID        uint      `gorm:"primary_key"`
	Time      time.Time `gorm:"not null"`
	MarketRef uint      `gorm:"unique_index:idx_market_trade;not null"`
	OrderID   string    `gorm:"unique_index:idx_market_trade;not null"`
	Type      string
	Side      string
	Price     float64
	Amount    float64
	Total     float64
	Market    *Market `gorm:"foreignkey:MarketRef;association_save_reference:false"`
}

type OrderBook struct {
	ID        uint        `gorm:"primary_key"`
	Time      time.Time   `gorm:"unique_index:idx_market_orderbook;not null"`
	MarketRef uint        `gorm:"unique_index:idx_market_orderbook;not null"`
	Bids      []*PriceVol `gorm:"many2many:bid_pricevols"`
	Asks      []*PriceVol `gorm:"many2many:ask_pricevols"`
	Market    *Market     `gorm:"foreignkey:MarketRef;association_save_reference:false"`
}

type PriceVol struct {
	ID     uint         `gorm:"primary_key"`
	Price  float64      `gorm:"unique_index:idx_price_volume"`
	Volume float64      `gorm:"unique_index:idx_price_volume"`
	Occurs []*OrderBook `gorm:"many2many:bid_pricevols;many2many:ask_pricevols"`
}

func (s *Symbol) ParseString(str string) error {
	ss := strings.Split(str, "_")
	if len(ss) != 2 {
		return fmt.Errorf("not a valid traded pair")
	}

	if s.Base == nil {
		s.Base = &Currency{Abbr: ss[0]}
	} else {
		s.Base.Abbr = ss[0]
	}

	if s.Quote == nil {
		s.Quote = &Currency{Abbr: ss[1]}
	} else {
		s.Quote.Abbr = ss[1]
	}
	return nil
}

func (s *Symbol) String() string {
	return s.Base.Abbr + "_" + s.Quote.Abbr
}

func (s *Symbol) MarshallJSON() ([]byte, error) {
	str := s.String()

	return []byte(str), nil
}

func (e Exchanger) GetCurrencyByName(name string) *Currency {
	if e.Currencies == nil {
		panic("currencies used but not initialized")
	}
	dname := strings.ToLower(name)
	for _, c := range e.Currencies {
		if c == nil {
			log.Println("nil currency pointer in exchanger")
			continue
		}
		cname := strings.ToLower(c.Name)
		if strings.Compare(cname, dname) == 0 {
			return c
		}
	}
	return nil
}
