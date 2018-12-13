package database

import (
	"time"

	"github.com/exchangedata/common"
	"github.com/jinzhu/gorm"
)

type Exchanger struct {
	ID         int
	Name       string
	Market     []MarketPair
	Currencies []Currency
	//	WebAPI     WebService
	//	WssAPI     WssService
}

type MarketPair struct {
	ID     int
	Symbol common.Symbol
}

type Symbol struct {
	ID    int
	Base  string `gorm:"size:16"`
	Quote string `gorm:"size:16"`
}

type Trade struct {
	gorm.Model
	Time        time.Time
	ExchangerID int16 `gorm:"index"`
	SymbolID    int32
	OrderID     int64
	Type        int8
	Side        int8
	Price       float64
	Amount      float64
	Total       float64
}

type Currency struct {
	ID        int
	ShortName string `gorm:"size:8"`
	LongName  string `gorm:"size:32"`
}

type WebService struct {
	ID          int
	ExchangerID int    `gorm:"index"`
	Version     string `gorm:"size:32"`
}

type User struct {
	ID int
}
