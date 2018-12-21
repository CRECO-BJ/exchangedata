package bittrex

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type jsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type btAddress struct {
	Currency string `json:"Currency"`
	Address  string `json:"Address"`
}

type btBalance struct {
	Currency      string          `json:"Currency"`
	Balance       decimal.Decimal `json:"Balance"`
	Available     decimal.Decimal `json:"Available"`
	Pending       decimal.Decimal `json:"Pending"`
	CryptoAddress string          `json:"CryptoAddress"`
	Requested     bool            `json:"Requested"`
	Uuid          string          `json:"Uuid"`
}

type btCandle struct {
	TimeStamp  CandleTime      `json:"T"`
	Open       decimal.Decimal `json:"O"`
	Close      decimal.Decimal `json:"C"`
	High       decimal.Decimal `json:"H"`
	Low        decimal.Decimal `json:"L"`
	Volume     decimal.Decimal `json:"V"`
	BaseVolume decimal.Decimal `json:"BV"`
}

type btNewCandles struct {
	Ticks []btCandle `json:"ticks"`
}

var CANDLE_INTERVALS = map[string]bool{
	"oneMin":    true,
	"fiveMin":   true,
	"thirtyMin": true,
	"hour":      true,
	"day":       true,
}

type CandleTime struct {
	time.Time
}

func (t *CandleTime) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("could not parse time %s", string(b))
	}
	// trim enclosing ""
	result, err := time.Parse("2006-01-02T15:04:05", string(b[1:len(b)-1]))
	if err != nil {
		return fmt.Errorf("could not parse time: %v", err)
	}
	t.Time = result
	return nil
}

type btCurrency struct {
	Currency        string          `json:"Currency"`
	CurrencyLong    string          `json:"CurrencyLong"`
	MinConfirmation int             `json:"MinConfirmation"`
	TxFee           decimal.Decimal `json:"TxFee"`
	IsActive        bool            `json:"IsActive"`
	CoinType        string          `json:"CoinType"`
	BaseAddress     string          `json:"BaseAddress"`
	Notice          string          `json:"Notice"`
}

type btDeposit struct {
	Id            int64           `json:"Id"`
	Amount        decimal.Decimal `json:"Amount"`
	Currency      string          `json:"Currency"`
	Confirmations int             `json:"Confirmations"`
	LastUpdated   jTime           `json:"LastUpdated"`
	TxId          string          `json:"TxId"`
	CryptoAddress string          `json:"CryptoAddress"`
}

type btDistribution struct {
	Distribution   []BalanceD      `json:"Distribution"`
	Balances       decimal.Decimal `json:"Balances"`
	AverageBalance decimal.Decimal `json:"AverageBalance"`
}

type BalanceD struct {
	BalanceD decimal.Decimal `json:"Balance"`
}

const TIME_FORMAT = "2006-01-02T15:04:05"

type jTime struct {
	time.Time
}

func (jt *jTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse(TIME_FORMAT, s)
	if err != nil {
		return err
	}
	jt.Time = t
	return nil
}

func (jt jTime) MarshalJSON() ([]byte, error) {
	return json.Marshal((*time.Time)(&jt.Time).Format(TIME_FORMAT))
}

type btMarket struct {
	MarketCurrency     string          `json:"MarketCurrency"`
	BaseCurrency       string          `json:"BaseCurrency"`
	MarketCurrencyLong string          `json:"MarketCurrencyLong"`
	BaseCurrencyLong   string          `json:"BaseCurrencyLong"`
	MinTradeSize       decimal.Decimal `json:"MinTradeSize"`
	MarketName         string          `json:"MarketName"`
	IsActive           bool            `json:"IsActive"`
	IsRestricted       bool            `json:"IsRestricted"`
	Notice             string          `json:"Notice"`
	IsSponsored        bool            `json:"IsSponsored"`
	LogoUrl            string          `json:"LogoUrl"`
	Created            string          `json:"Created"`
}

type btMarketSummary struct {
	MarketName     string          `json:"MarketName"`
	High           decimal.Decimal `json:"High"`
	Low            decimal.Decimal `json:"Low"`
	Ask            decimal.Decimal `json:"Ask"`
	Bid            decimal.Decimal `json:"Bid"`
	OpenBuyOrders  int             `json:"OpenBuyOrders"`
	OpenSellOrders int             `json:"OpenSellOrders"`
	Volume         decimal.Decimal `json:"Volume"`
	Last           decimal.Decimal `json:"Last"`
	BaseVolume     decimal.Decimal `json:"BaseVolume"`
	PrevDay        decimal.Decimal `json:"PrevDay"`
	TimeStamp      string          `json:"TimeStamp"`
}

type btOrder struct {
	OrderUuid         string          `json:"OrderUuid"`
	Exchange          string          `json:"Exchange"`
	TimeStamp         jTime           `json:"TimeStamp"`
	OrderType         string          `json:"OrderType"`
	Limit             decimal.Decimal `json:"Limit"`
	Quantity          decimal.Decimal `json:"Quantity"`
	QuantityRemaining decimal.Decimal `json:"QuantityRemaining"`
	Commission        decimal.Decimal `json:"Commission"`
	Price             decimal.Decimal `json:"Price"`
	PricePerUnit      decimal.Decimal `json:"PricePerUnit"`
}

// For getorder
type btOrder2 struct {
	AccountId                  string
	OrderUuid                  string `json:"OrderUuid"`
	Exchange                   string `json:"Exchange"`
	Type                       string
	Quantity                   decimal.Decimal `json:"Quantity"`
	QuantityRemaining          decimal.Decimal `json:"QuantityRemaining"`
	Limit                      decimal.Decimal `json:"Limit"`
	Reserved                   decimal.Decimal
	ReserveRemaining           decimal.Decimal
	CommissionReserved         decimal.Decimal
	CommissionReserveRemaining decimal.Decimal
	CommissionPaid             decimal.Decimal
	Price                      decimal.Decimal `json:"Price"`
	PricePerUnit               decimal.Decimal `json:"PricePerUnit"`
	Opened                     string
	Closed                     string
	IsOpen                     bool
	Sentinel                   string
	CancelInitiated            bool
	ImmediateOrCancel          bool
	IsConditional              bool
	Condition                  string
	ConditionTarget            decimal.Decimal
}

type btOrderBook struct {
	Buy  []Orderb `json:"buy"`
	Sell []Orderb `json:"sell"`
}

type Orderb struct {
	Quantity decimal.Decimal `json:"Quantity"`
	Rate     decimal.Decimal `json:"Rate"`
}

type btTicker struct {
	Bid  decimal.Decimal `json:"Bid"`
	Ask  decimal.Decimal `json:"Ask"`
	Last decimal.Decimal `json:"Last"`
}

// Used in getmarkethistory
type btTrade struct {
	OrderUuid int64           `json:"Id"`
	Timestamp jTime           `json:"TimeStamp"`
	Quantity  decimal.Decimal `json:"Quantity"`
	Price     decimal.Decimal `json:"Price"`
	Total     decimal.Decimal `json:"Total"`
	FillType  string          `json:"FillType"`
	OrderType string          `json:"OrderType"`
}

type Uuid struct {
	Id string `json:"uuid"`
}

type btWithdrawal struct {
	PaymentUuid    string          `json:"PaymentUuid"`
	Currency       string          `json:"Currency"`
	Amount         decimal.Decimal `json:"Amount"`
	Address        string          `json:"Address"`
	Opened         jTime           `json:"Opened"`
	Authorized     bool            `json:"Authorized"`
	PendingPayment bool            `json:"PendingPayment"`
	TxCost         decimal.Decimal `json:"TxCost"`
	TxId           string          `json:"TxId"`
	Canceled       bool            `json:"Canceled"`
}
