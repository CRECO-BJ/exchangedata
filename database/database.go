package database

import (
	"log"

	"github.com/exchangedata/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DataStore ...
type DataStore struct {
	Dialect  string
	Name     string
	User     string
	Password string

	db *gorm.DB
}

func NewDataStore(Dialet string) *DataStore {
	return &DataStore{Dialect: Dialet,
		Name:     "exdata",
		User:     "root",
		Password: "office98"} //@TODO later should change to from the config
}

func (d *DataStore) OpenDB() error {
	//	user+":"+passwd+"@tcp(127.0.0.1:3306)/"+name+"?charset=utf8&parseTime=True&loc=Local"
	uri := d.User + ":" + d.Password + "@/" + d.Name + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(d.Dialect, uri)
	if err != nil {
		log.Println("DB open with", err)
		return err
	}
	d.db = db

	return nil
}

func (d *DataStore) CloseDB() {
	d.db.Close()
}

func (d *DataStore) GetDB() *gorm.DB {
	return d.db
}

func (d *DataStore) AutoMigrate() *gorm.DB {
	return d.db.AutoMigrate(&common.Exchanger{}, &common.Market{}, &common.Currency{},
		&common.CommunicationAPI{}, &common.AccessSecret{}, &common.Symbol{},
		&common.Ticker{}, &common.Trade{}, &common.OrderBook{}, &common.PriceVol{})
}

func (d *DataStore) UpdateCurrency(c *common.Currency) *gorm.DB {
	t := &common.Currency{}
	if !d.db.Where("name = ?", c.Name).First(t).RecordNotFound() {
		if !t.AbbrFinal {
			c.Abbr = t.Abbr
		}
		c.ID = t.ID
	}

	return d.db.Save(c) //Set("gorm:save_associations", false).
}

func (d *DataStore) UpdateExchanger(c *common.Exchanger) *gorm.DB {
	return d.db.Where("name = ?", c.Name).FirstOrCreate(c)
}

func (d *DataStore) UpdateSymbol(c *common.Symbol) *gorm.DB {
	if d.UpdateCurrency(c.Base).RecordNotFound() {
		log.Println("update currency error:", c.Base, d.db.Error)
	}
	if d.UpdateCurrency(c.Quote).RecordNotFound() {
		log.Println("update currency error:", c.Quote, d.db.Error)
	}
	return d.db.Where("quote_id = ? and base_id = ?", c.Quote.ID, c.Base.ID).FirstOrCreate(c)
}

func (d *DataStore) UpdateMarket(c *common.Market) *gorm.DB {
	t := &common.Market{}
	if c.ExRef == 0 && c.Exchanger != nil {
		if c.Exchanger.ID == 0 {
			if d.UpdateExchanger(c.Exchanger).RecordNotFound() {
				log.Println("update exchanger error:", c.Exchanger, d.db.Error)
			}
		}
		c.ExRef = c.Exchanger.ID
	}
	if c.SymRef == 0 && c.Symbol != nil {
		if c.Symbol.ID == 0 {
			if d.UpdateSymbol(c.Symbol).RecordNotFound() {
				log.Println("update symbol error:", c.Symbol, d.db.Error)
			}
		}
		c.SymRef = c.Symbol.ID
	}
	if !d.db.Where("name = ? AND ex_ref = ?", c.Name, c.ExRef).First(t).RecordNotFound() {
		c.ID = t.ID
	}
	return d.db.Save(c)
}

func (d *DataStore) UpdateTicker(c *common.Ticker) *gorm.DB {
	if c.MarketRef == 0 && c.Market == nil {
		return nil
	} else if c.MarketRef == 0 { // tickers from the net may not have the MarketRef set, sometimes from a new market
		if d.UpdateMarket(c.Market).RecordNotFound() {
			return d.db
		}
		c.MarketRef = c.Market.ID
	} else if c.Market == nil {
		c.Market = &common.Market{}
		d.db.First(c.Market, c.MarketRef)
	}

	return d.db.Where("time = ? and market_ref = ?", c.Time, c.MarketRef).FirstOrCreate(c)
}

func (d *DataStore) UpdateTrade(c *common.Trade) *gorm.DB {
	if c.MarketRef == 0 && c.Market == nil {
		return nil
	} else if c.MarketRef == 0 {
		if d.UpdateMarket(c.Market).RecordNotFound() {
			return d.db
		}
		c.MarketRef = c.Market.ID
	} else if c.Market == nil {
		c.Market = &common.Market{}
		d.db.First(c.Market, c.MarketRef)
	}

	return d.db.Where("order_id = ? and market_ref = ?", c.OrderID, c.MarketRef).FirstOrCreate(c)
}

func (d *DataStore) UpdateOrderBook(c *common.OrderBook) *gorm.DB {
	if c.MarketRef == 0 && c.Market == nil {
		return nil
	} else if c.MarketRef == 0 {
		if d.UpdateMarket(c.Market).RecordNotFound() {
			return d.db
		}
		c.MarketRef = c.Market.ID
	} else if c.Market == nil {
		c.Market = &common.Market{}
		d.db.First(c.Market, c.MarketRef)
	}

	for _, p := range c.Asks {
		if d.UpdatePriceVol(p).Error != nil {
			log.Println("error update pricevol:", p, d.db.Error)
		}
	}

	for _, p := range c.Bids {
		if d.UpdatePriceVol(p).Error != nil {
			log.Println("error update pricevol:", p, d.db.Error)
		}
	}

	return d.db.Where("time = ? and market_ref = ?", c.Time, c.MarketRef).FirstOrCreate(c)
}

func (d *DataStore) UpdatePriceVol(c *common.PriceVol) *gorm.DB {
	return d.db.Where("price = ? and volume = ?", c.Price, c.Volume).FirstOrCreate(c)
}
