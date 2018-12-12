package poloniex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type client struct {
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	throttle    <-chan time.Time
	httpTimeout time.Duration
	debug       bool
}

var (
	// Technically 6 req/s allowed, but we're being nice / playing it safe.
	reqInterval = 200 * time.Millisecond
)

// NewClient return a new Poloniex HTTP client
func NewClient(apiKey, apiSecret string) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, time.Tick(reqInterval), 30 * time.Second, false}
}

// NewClientWithCustomTimeout returns a new Poloniex HTTP client with custom timeout
func NewClientWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, time.Tick(reqInterval), timeout, false}
}

func (c client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		if c.debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if c.debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from Poloniex API")
	}
}

func (c *client) makeReq(method, resource string, payload map[string]string, authNeeded bool, respCh chan<- []byte, errCh chan<- error) {
	body := []byte{}
	connectTimer := time.NewTimer(c.httpTimeout)

	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s/%s", API_BASE, resource)
	}

	formValues := url.Values{}
	for key, value := range payload {
		formValues.Set(key, value)
	}
	formData := formValues.Encode()

	req, err := http.NewRequest(method, rawurl, strings.NewReader(formData))
	if err != nil {
		respCh <- body
		errCh <- errors.New("You need to set API Key and API Secret to call this method")
		return
	}

	if authNeeded {
		if len(c.apiKey) == 0 || len(c.apiSecret) == 0 {
			respCh <- body
			errCh <- errors.New("You need to set API Key and API Secret to call this method")
			return
		}

		mac := hmac.New(sha512.New, []byte(c.apiSecret))
		_, err := mac.Write([]byte(formData))
		if err != nil {
			errCh <- err
		}
		sig := hex.EncodeToString(mac.Sum(nil))
		req.Header.Add("Key", c.apiKey)
		req.Header.Add("Sign", sig)
	}

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header.Add("Accept", "application/json")

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		respCh <- body
		errCh <- err
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		respCh <- body
		errCh <- err
		return
	}
	if resp.StatusCode != 200 {
		respCh <- body
		errCh <- errors.New(resp.Status)
		return
	}

	respCh <- body
	errCh <- nil
	close(respCh)
	close(errCh)
}

// do prepare and process HTTP request to Poloniex API
func (c *client) do(method, resource string, payload map[string]string, authNeeded bool) (response []byte, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)
	<-c.throttle
	go c.makeReq(method, resource, payload, authNeeded, respCh, errCh)
	response = <-respCh
	err = <-errCh
	return
}

// doCommand prepares an authorized command-request for Poloniex's tradingApi
func (c *client) doCommand(command string, payload map[string]string) (response []byte, err error) {
	if payload == nil {
		payload = make(map[string]string)
	}

	payload["command"] = command
	payload["nonce"] = strconv.FormatInt(time.Now().UnixNano(), 10)

	return c.do("POST", "tradingApi", payload, true)
}

type Currency struct {
	Id                 int     `json:"id"`
	Name               string  `json:"name"`
	MaxDailyWithdrawal string  `json:"maxDailyWithdrawal"`
	TxFee              float64 `json:"txFee,string"`
	MinConf            int     `json:"minConf"`
	Disabled           int     `json:"disabled"`
	Frozen             int     `json:"frozen"`
	Delisted           int     `json:"delisted"`
}

type Currencies struct {
	Pair map[string]Currency
}

type Balance struct {
	Available string `json:"available"`
	BtcValue  string `json:"btcValue"`
	OnOrders  string `json:"onOrders"`
}

type CandleStick struct {
	Date            PoloniexDate `json:"date"`
	High            float64      `json:"high"`
	Low             float64      `json:"low"`
	Open            float64      `json:"open"`
	Close           float64      `json:"close"`
	Volume          float64      `json:"volume"`
	QuoteVolume     float64      `json:"quoteVolume"`
	WeightedAverage float64      `json:"weightedAverage"`
}

type PoloniexDate struct {
	time.Time
}

func (pd *PoloniexDate) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return errors.New("Timestamp invalid (can't parse int64)")
	}
	pd.Time = time.Unix(i, 0)
	return nil
}

type Deposit struct {
	Currency      string    `json:"currency"`
	Address       string    `json:"address"`
	Amount        float64   `json:"amount,string"`
	Confirmations uint64    `json:"confirmations"`
	TxId          string    `json:"txid"`
	Date          time.Time `json:"timestamp"`
	Status        string    `json:"status"`
}

func (t *Deposit) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Deposit
	aux := &struct {
		Date int64 `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Date = time.Unix(aux.Date, 0)
	return nil
}

type Lending struct {
	Id       uint64    `json:"id"`
	Currency string    `json:"currency"`
	Rate     float64   `json:"rate,string"`
	Amount   float64   `json:"amount,string"`
	Duration float64   `json:"duration,string"`
	Interest float64   `json:"interest,string"`
	Fee      float64   `json:"fee,string"`
	Earned   float64   `json:"earned,string"`
	Open     time.Time `json:"open,string"`
	Close    time.Time `json:"close,string"`
}

func (t *Lending) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Lending
	aux := &struct {
		Open  string `json:"open"`
		Close string `json:"close"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Open, err = time.Parse("2006-01-02 15:04:05", aux.Open)
	t.Close, err = time.Parse("2006-01-02 15:04:05", aux.Close)
	if err != nil {
		return err
	}
	return nil
}

type OrderBook struct {
	Asks     [][]interface{} `json:"asks"`
	Bids     [][]interface{} `json:"bids"`
	IsFrozen int             `json:"isFrozen,string"`
	Error    string          `json:"error"`
}

// This can probably be implemented using UnmarshalJSON
/*
type OrderBook struct {
	Bids     []Orderb `json:"bids"`
	Asks     []Orderb `json:"asks"`
	IsFrozen int      `json:"isFrozen,string"`
}
type Orderb struct {
	Rate     string
	Quantity float64
}
*/

type OpenOrder struct {
	OrderNumber int64   `json:"orderNumber,string"`
	Type        string  `json:"type"`
	Rate        float64 `json:"rate,string"`
	Amount      float64 `json:"amount,string"`
	Total       float64 `json:"total,string"`
}

const (
	API_BASE = "https://poloniex.com" // Poloniex API endpoint
)

// New returns an instantiated poloniex struct
func New(apiKey, apiSecret string) *Poloniex {
	client := NewClient(apiKey, apiSecret)
	return &Poloniex{client}
}

// New returns an instantiated poloniex struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Poloniex {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, timeout)
	return &Poloniex{client}
}

// poloniex represent a poloniex client
type Poloniex struct {
	client *client
}

// set enable/disable http request/response dump
func (c *Poloniex) SetDebug(enable bool) {
	c.client.debug = enable
}

// GetTickers is used to get the ticker for all markets
func (b *Poloniex) GetTickers() (tickers map[string]Ticker, err error) {
	r, err := b.client.do("GET", "public?command=returnTicker", nil, false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &tickers); err != nil {
		return
	}
	return
}

// GetVolumes is used to get the volume for all markets
func (b *Poloniex) GetVolumes() (vc VolumeCollection, err error) {
	r, err := b.client.do("GET", "public?command=return24hVolume", nil, false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &vc); err != nil {
		return
	}
	return
}

func (b *Poloniex) GetCurrencies() (currencies Currencies, err error) {
	r, err := b.client.do("GET", "public?command=returnCurrencies", nil, false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &currencies.Pair); err != nil {
		return
	}
	return
}

// GetOrderBook is used to get retrieve the orderbook for a given market
// market: a string literal for the market (ex: BTC_NXT). 'all' not implemented.
// cat: bid, ask or both to identify the type of orderbook to return.
// depth: how deep of an order book to retrieve
func (b *Poloniex) GetOrderBook(market, cat string, depth int) (orderBook OrderBook, err error) {
	// not implemented
	if cat != "bid" && cat != "ask" && cat != "both" {
		cat = "both"
	}
	if depth > 100 {
		depth = 100
	}
	if depth < 1 {
		depth = 1
	}

	r, err := b.client.do("GET", fmt.Sprintf("public?command=returnOrderBook&currencyPair=%s&depth=%d", strings.ToUpper(market), depth), nil, false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &orderBook); err != nil {
		return
	}
	if orderBook.Error != "" {
		err = errors.New(orderBook.Error)
		return
	}
	return
}

// GetOrderTrades is used to get returns all trades involving a given order
// orderNumber: order number.
func (b *Poloniex) GetOrderTrades(orderNumber int) (tradeOrderTransaction []TradeOrderTransaction, err error) {
	r, err := b.client.doCommand("returnOrderTrades", map[string]string{"orderNumber": fmt.Sprintf("%d", orderNumber)})
	if err != nil {
		return
	}
	if string(r) == `{"error":"Order not found, or you are not the person who placed it."}` {
		err = fmt.Errorf("Error : order not found, or you are not the person who placed it.")
		return
	}
	if err = json.Unmarshal(r, &tradeOrderTransaction); err != nil {
		return
	}
	return
}

// Returns candlestick chart data. Required GET parameters are "currencyPair",
// "period" (candlestick period in seconds; valid values are 300, 900, 1800,
// 7200, 14400, and 86400), "start", and "end". "Start" and "end" are given in
// UNIX timestamp format and used to specify the date range for the data
// returned.
func (b *Poloniex) ChartData(currencyPair string, period int, start, end time.Time) (candles []*CandleStick, err error) {
	r, err := b.client.do("GET", fmt.Sprintf(
		"public?command=returnChartData&currencyPair=%s&period=%d&start=%d&end=%d",
		strings.ToUpper(currencyPair),
		period,
		start.Unix(),
		end.Unix(),
	), nil, false)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &candles); err != nil {
		return
	}

	return
}

func (b *Poloniex) GetBalances() (balances map[string]Balance, err error) {
	balances = make(map[string]Balance)
	r, err := b.client.doCommand("returnCompleteBalances", nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &balances); err != nil {
		return
	}

	return
}

func (b *Poloniex) GetTradeHistory(pair string, start uint32) (trades map[string][]Trade, err error) {
	trades = make(map[string][]Trade)
	r, err := b.client.doCommand("returnTradeHistory", map[string]string{"currencyPair": pair, "start": strconv.FormatUint(uint64(start), 10)})
	if err != nil {
		return
	}

	if pair == "all" {
		if err = json.Unmarshal(r, &trades); err != nil {
			return
		}
	} else {
		var pairTrades []Trade
		if err = json.Unmarshal(r, &pairTrades); err != nil {
			return
		}
		trades[pair] = pairTrades
	}

	return
}

type responseDepositsWithdrawals struct {
	Deposits    []Deposit    `json:"deposits"`
	Withdrawals []Withdrawal `json:"withdrawals"`
}

func (b *Poloniex) GetDepositsWithdrawals(start uint32, end uint32) (deposits []Deposit, withdrawals []Withdrawal, err error) {
	deposits = make([]Deposit, 0)
	withdrawals = make([]Withdrawal, 0)
	r, err := b.client.doCommand("returnDepositsWithdrawals", map[string]string{"start": strconv.FormatUint(uint64(start), 10), "end": strconv.FormatUint(uint64(end), 10)})
	if err != nil {
		return
	}
	var response responseDepositsWithdrawals
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}

	return response.Deposits, response.Withdrawals, nil
}

func (b *Poloniex) Buy(pair string, rate float64, amount float64, tradeType string) (TradeOrder, error) {
	reqParams := map[string]string{
		"currencyPair": pair, "rate": strconv.FormatFloat(rate, 'f', -1, 64),
		"amount": strconv.FormatFloat(amount, 'f', -1, 64)}
	if tradeType != "" {
		reqParams[tradeType] = "1"
	}
	r, err := b.client.doCommand("buy", reqParams)
	if err != nil {
		return TradeOrder{}, err
	}
	var orderResponse TradeOrder
	if err = json.Unmarshal(r, &orderResponse); err != nil {
		return TradeOrder{}, err
	}

	if orderResponse.ErrorMessage != "" {
		return TradeOrder{}, errors.New(orderResponse.ErrorMessage)
	}

	return orderResponse, nil
}

func (b *Poloniex) Sell(pair string, rate float64, amount float64, tradeType string) (TradeOrder, error) {
	reqParams := map[string]string{
		"currencyPair": pair, "rate": strconv.FormatFloat(rate, 'f', -1, 64),
		"amount": strconv.FormatFloat(amount, 'f', -1, 64)}
	if tradeType != "" {
		reqParams[tradeType] = "1"
	}
	r, err := b.client.doCommand("sell", reqParams)
	if err != nil {
		return TradeOrder{}, err
	}
	var orderResponse TradeOrder
	if err = json.Unmarshal(r, &orderResponse); err != nil {
		return TradeOrder{}, err
	}

	if orderResponse.ErrorMessage != "" {
		return TradeOrder{}, errors.New(orderResponse.ErrorMessage)
	}

	return orderResponse, nil
}

func (b *Poloniex) GetOpenOrders(pair string) (openOrders map[string][]OpenOrder, err error) {
	openOrders = make(map[string][]OpenOrder)
	r, err := b.client.doCommand("returnOpenOrders", map[string]string{"currencyPair": pair})
	if err != nil {
		return
	}
	if pair == "all" {
		if err = json.Unmarshal(r, &openOrders); err != nil {
			return
		}
	} else {
		var onePairOrders []OpenOrder
		if err = json.Unmarshal(r, &onePairOrders); err != nil {
			return
		}
		openOrders[pair] = onePairOrders
	}
	return
}

func (b *Poloniex) CancelOrder(orderNumber string) error {
	_, err := b.client.doCommand("cancelOrder", map[string]string{"orderNumber": orderNumber})
	if err != nil {
		return err
	}
	return nil
}

// Returns whole lending history chart data. Required GET parameters are "start",
// "end" (UNIX timestamp format and used to specify the date range for the data returned)
// and optionally limit (<0 for no limit, poloniex automatically limits to 500 records)
func (b *Poloniex) LendingHistory(start, end time.Time, limit int) (lendings []Lending, err error) {
	lendings = make([]Lending, 0)
	reqParams := map[string]string{
		"start": strconv.FormatUint(uint64(start.Unix()), 10),
		"end":   strconv.FormatUint(uint64(end.Unix()), 10)}
	if limit >= 0 {
		reqParams["limit"] = string(limit)
	}

	r, err := b.client.doCommand("returnLendingHistory", reqParams)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &lendings); err != nil {
		return
	}

	return
}

type Tickers struct {
	Pair map[string]Ticker
}

type Ticker struct {
	Id            int     `json:"id"`
	Last          float64 `json:"last,string"`
	LowestAsk     float64 `json:"lowestAsk,string"`
	HighestBid    float64 `json:"highestBid,string"`
	PercentChange float64 `json:"percentChange,string"`
	BaseVolume    float64 `json:"baseVolume,string"`
	QuoteVolume   float64 `json:"quoteVolume,string"`
	IsFrozen      int     `json:"isFrozen,string"`
	High24Hr      float64 `json:"high24hr,string"`
	Low24Hr       float64 `json:"low24hr,string"`
}

const (
	TRADE_FILL_OR_KILL        = "fillOrKill"
	TRADE_IMMEDIATE_OR_CANCEL = "immediateOrCancel"
	TRADE_POST_ONLY           = "postOnly"
)

type Trade struct {
	GlobalTradeID uint64    `json:"globalTradeID"`
	TradeID       uint64    `json:"tradeID,string"`
	Date          time.Time `json:"date,string"`
	Type          string    `json:"type"`
	Category      string    `json:"category"`
	Rate          float64   `json:"rate,string"`
	Amount        float64   `json:"amount,string"`
	Total         float64   `json:"total,string"`
	Fee           float64   `json:"fee,string"`
}

func (t *Trade) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Trade
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Date, err = time.Parse("2006-01-02 15:04:05", aux.Date)
	if err != nil {
		return err
	}
	return nil
}

type ResultingTrade struct {
	Amount  float64 `json:"amount,string"`
	Date    string  `json:"date"`
	Rate    float64 `json:"rate,string"`
	Total   float64 `json:"total,string"`
	TradeID string  `json:"tradeID"`
	Type    string  `json:"type"`
}

type TradeOrder struct {
	OrderNumber     string           `json:"orderNumber"`
	ResultingTrades []ResultingTrade `json:"resultingTrades"`
	ErrorMessage    string           `json:"error"`
}

type TradeOrderTransaction struct {
	GlobalTradeID uint64    `json:"globalTradeID"`
	TradeID       uint64    `json:"tradeID"`
	CurrencyPair  string    `json:"currencyPair"`
	Type          string    `json:"type"`
	Rate          float64   `json:"rate,string"`
	Amount        float64   `json:"amount,string"`
	Total         float64   `json:"total,string"`
	Fee           float64   `json:"fee,string"`
	Date          time.Time `json:"date,string"`
}

type Volume map[string]float64

type VolumeCollection struct {
	TotalBTC  float64 `json:"totalBTC,string"`
	TotalETH  float64 `json:"totalETH,string"`
	TotalUSDC float64 `json:"totalUSDC,string"`
	TotalUSDT float64 `json:"totalUSDT,string"`
	TotalXMR  float64 `json:"totalXMR,string"`
	TotalXUSD float64 `json:"totalXUSD,string"`
	Volumes   map[string]Volume
}

func (tc *VolumeCollection) UnmarshalJSON(b []byte) error {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	tc.Volumes = make(map[string]Volume)
	for k, v := range m {
		switch k {
		case "totalBTC":
			f, err := parseJSONFloatString(v)
			if err != nil {
				return err
			}
			tc.TotalBTC = f
		case "totalETH":
			f, err := parseJSONFloatString(v)
			if err != nil {
				return err
			}
			tc.TotalETH = f
		case "totalUSDC":
			f, err := parseJSONFloatString(v)
			if err != nil {
				return err
			}
			tc.TotalUSDC = f
		case "totalUSDT":
			f, err := parseJSONFloatString(v)
			if err != nil {
				return err
			}
			tc.TotalUSDT = f
		case "totalXMR":
			f, err := parseJSONFloatString(v)
			if err != nil {
				return err
			}
			tc.TotalXMR = f
		case "totalXUSD":
			f, err := parseJSONFloatString(v)
			if err != nil {
				return err
			}
			tc.TotalXUSD = f
		default:
			t := make(Volume)
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			tc.Volumes[k] = t
		}
	}
	return nil
}

func (t *Volume) UnmarshalJSON(b []byte) error {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	for k, v := range m {
		f, err := parseJSONFloatString(v)
		if err != nil {
			return err
		}
		(*t)[k] = f
	}
	return nil
}

func parseJSONFloatString(b json.RawMessage) (float64, error) {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

type Withdrawal struct {
	WithdrawalNumber uint64    `json:"withdrawalNumber"`
	Currency         string    `json:"currency"`
	Address          string    `json:"address"`
	Amount           float64   `json:"amount,string"`
	Date             time.Time `json:"timestamp"`
	Status           string    `json:"status"`
	TxId             string    `json:"txid"`
	IpAddress        string    `json:"ipAddress"`
}

func (t *Withdrawal) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Withdrawal
	aux := &struct {
		Date int64 `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Date = time.Unix(aux.Date, 0)
	if strings.HasPrefix(aux.Status, "COMPLETE") {
		t.TxId = strings.TrimPrefix(t.Status, "COMPLETE: ")
		t.Status = "COMPLETE"
	}
	return nil
}
