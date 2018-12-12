package poloniex

import (
	"testing"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func TestNew(t *testing.T) {
	// Poloniex client
	poloniex := New(API_KEY, API_SECRET)

	// Get Ticker (BTC-VTC)
	/*
		ticker, err := GetTicker("BTC-DRK")
		fmt.Println(err, ticker)
	*/

	// Get Order Trades (3356534544)
	/*
		tradeOrderTransaction, err := GetOrderTrades(3356534544)
		fmt.Println("Error:", tradeOrderTransaction)
	*/

	// Get Tickers
	/*
		tickers, err := GetTickers()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			for key, ticker := range tickers {
				fmt.Printf("Ticker: %s, Last: %.8f\n", key, ticker.Last)
			}
		}
		tickerName := "BTC_FLO"
		ticker, ok := tickers[tickerName]
		if ok {
			fmt.Printf("BTC_FLO Last: %.8f\n", ticker.Last)
		} else {
			fmt.Println("ticker not found - ", tickerName)
		}
	*/

	// Get Volumes
	/*
		volumes, err := GetVolumes()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			for key, volume := range volumes.Volumes {
				fmt.Printf("Ticker: %s Value: %#v\n", key, volume["BTC"])
			}
		}
	*/

	// Get CandleStick chart data ( OHLCV )
	/*
		candles, err := client.ChartData("BTC_XMR", 300, time.Now().Add(-time.Hour), time.Now())
		if err != nil {
			panic(err)
		}
		for _, candle := range candles {
			fmt.Printf("BTC_XMR %s\tOpened at: %f\tClosed at: %f\n", candle.Date, candle.Open, candle.Close)
		}
	*/

	// Get markets
	/*
		markets, err := GetMarkets()
		fmt.Println(err, markets)
	*/

	// Get orders book
	/*
		orderBook, err := GetOrderBook("BTC-DRK", "both", 100)
		fmt.Println(err, orderBook)
	*/

	// Market history
	/*
		marketHistory, err := GetMarketHistory("BTC-DRK", 100)
		for _, trade := range marketHistory {
			fmt.Println(err, trade.Timestamp.String(), trade.Quantity, trade.Price)
		}
	*/

	// Market

	// BuyLimit
	/*
		uuid, err := BuyLimit("BTC-DOGE", 1000, 0.00000102)
		fmt.Println(err, uuid)
	*/

	// BuyMarket
	/*
		uuid, err := BuyLimit("BTC-DOGE", 1000)
		fmt.Println(err, uuid)
	*/

	// Sell limit
	/*
		uuid, err := SellLimit("BTC-DOGE", 1000, 0.00000115)
		fmt.Println(err, uuid)
	*/

	// Cancel Order
	/*
		err := CancelOrder("e3b4b704-2aca-4b8c-8272-50fada7de474")
		fmt.Println(err)
	*/

	// Get open orders
	/*
		orders, err := GetOpenOrders("BTC-DOGE")
		fmt.Println(err, orders)
	*/

	// Account
	// Get balances
	/*
		balances, err := GetBalances()
		fmt.Println(err, balances)
	*/

	// Get balance
	/*
		balance, err := GetBalance("DOGE")
		fmt.Println(err, balance)
	*/

	// Get address
	/*
		address, err := GetDepositAddress("QBC")
		fmt.Println(err, address)
	*/

	// WithDraw
	/*
		whitdrawUuid, err := Withdraw("QYQeWgSnxwtTuW744z7Bs1xsgszWaFueQc", "QBC", 1.1)
		fmt.Println(err, whitdrawUuid)
	*/

	// Get order history
	/*
		orderHistory, err := GetOrderHistory("BTC-DOGE", 10)
		fmt.Println(err, orderHistory)
	*/

	// Get getwithdrawal history
	/*
		withdrawalHistory, err := GetWithdrawalHistory("all", 0)
		fmt.Println(err, withdrawalHistory)
	*/

	// Get deposit history
	/*
		deposits, err := GetDepositHistory("all", 0)
		fmt.Println(err, deposits)
	*/
}
