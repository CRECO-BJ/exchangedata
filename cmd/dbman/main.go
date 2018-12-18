package main

import (
	"log"

	"github.com/exchangedata/common"
	"github.com/exchangedata/database"
)

var cur = []common.Currency{
	common.Currency{Name: "bitcoin", Abbr: "BTC"},
	common.Currency{Name: "Litecoin", Abbr: "LTC"},
	common.Currency{Name: "bitcoin cash", Abbr: "BCH"},
	common.Currency{Name: "Dogecoin", Abbr: "DOGE"},
}

var ex = []common.Exchanger{
	common.Exchanger{Name: "Poloniex"},
	common.Exchanger{Name: "bittrex"},
	common.Exchanger{Name: "okcoin", Info: "Okex is located in japan"},
}

var mkt = []*common.Market{
	&common.Market{Name: "LTC-BTC",
		Symbol: &common.Symbol{Base: &common.Currency{Name: "Litecoin", Abbr: "LTC"},
			Quote: &common.Currency{Name: "bitcoin", Abbr: "BTC"}},
		Active:     true,
		Info:       "common market test",
		Precision:  5,
		Limitation: common.Limitation{Min: 0.01, Max: 10000},
		MinStep:    0.000001,
		Exchanger:  &ex[1]},
	&common.Market{Name: "DOGE-BTC",
		Symbol: &common.Symbol{
			Base:  &cur[3],
			Quote: &cur[0]},
		Active:     true,
		Info:       "extiguish test",
		Precision:  8,
		Limitation: common.Limitation{Min: 100, Max: 10000000},
		MinStep:    0.00000000001,
		Exchanger:  &ex[1]},
	&common.Market{Name: "DOGE-BTC",
		Symbol: &common.Symbol{
			Base:  &cur[3],
			Quote: &cur[0]},
		Active:     true,
		Info:       "extiguish test",
		Precision:  8,
		Limitation: common.Limitation{Min: 100, Max: 10000000},
		MinStep:    0.00000000001,
		Exchanger:  &ex[0]},
}

var (
	Currencs  map[string]*common.Currency
	Exchanges map[string]*common.Exchanger
	Markes    map[uint]*common.Market
)

func init() {
	Currencs = make(map[string]*common.Currency)
	Exchanges = make(map[string]*common.Exchanger)
	Markes = make(map[uint]*common.Market)
	for k, m := range mkt {
		if Currencs[m.Symbol.Base.Name] == nil || Currencs[m.Symbol.Base.Name].Name != m.Symbol.Base.Name {
			Currencs[m.Symbol.Base.Name] = m.Symbol.Base
		} else {
			m.Symbol.Base = Currencs[m.Symbol.Base.Name]
		}
		if Currencs[m.Symbol.Quote.Name] == nil || Currencs[m.Symbol.Quote.Name].Name != m.Symbol.Quote.Name {
			Currencs[m.Symbol.Quote.Name] = m.Symbol.Quote
		} else {
			m.Symbol.Quote = Currencs[m.Symbol.Quote.Name]
		}
		if Exchanges[m.Exchanger.Name] == nil || Exchanges[m.Exchanger.Name].Name != m.Exchanger.Name {
			Exchanges[m.Exchanger.Name] = m.Exchanger
		} else {
			m.Exchanger = Exchanges[m.Exchanger.Name]
		}
		Markes[uint(k)] = m
	}
	log.Println(Currencs)
	log.Println(Exchanges)
	log.Println(Markes)
}

func main() {
	var err error
	ds := database.NewDataStore("mysql")
	if err = ds.OpenDB(); err != nil {
		panic("open db failed")
	}
	defer ds.CloseDB()

	ds.AutoMigrate()
	for _, c := range Currencs {
		if ds.UpdateCurrency(c).Error != nil {
			log.Println("error save", ds.GetDB().Error)
		}
	}
	for _, c := range Exchanges {
		if ds.UpdateExchanger(c).Error != nil {
			log.Println("error save", ds.GetDB().Error)
		}
	}
	for _, c := range Markes {
		if ds.UpdateMarket(c).Error != nil {
			log.Println("error save", ds.GetDB().Error)
		}
	}
}

/*
func CompareDiffer(src interface{}, dest interface{}) (differ []string, err error) {
	vst := reflect.TypeOf(src)
	vdt := reflect.TypeOf(dest)
	if vst != vdt {
		err = fmt.Errorf("error non comparable data")
		return
	}
	vsv := reflect.ValueOf(src)
	vdv := reflect.ValueOf(dest)
	for k := 0; k < vst.NumField(); k = k + 1 {
		if !reflect.DeepEqual(vsv.Field(k), vdv.Field(k)) {
			if vst.Field(k).Name != "ID" {
				reflect.Copy(vdv.Field(k), vsv.Field(k))
				differ = append(differ, vst.Field(k).Name)
			}
		}
	}
	return
}
*/
