package main

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
	"github.com/toorop/go-bittrex"
)

var (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {
	app := cli.NewApp()
	
	flags := []cli.Flag{
		altsrc.NewIntFlag(cli.IntFlag{Name: "test"}),
		cli.StringFlag{Name: "load"},
	}
	
	app.Action = func(c *cli.Context) error {
		fmt.Println("yaml ist rad")
		return nil
	}
	
	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load"))
	app.Flags = flags
	
	err := app.Run(os.Args)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Bittrex client
	bittrex := bittrex.New(API_KEY, API_SECRET)

	// Get markets
	markets, err := bittrex.GetMarkets()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, market := range markets {
			fmt.Println(market)
		}
	}

	// Get Ticker (BTC-USDT)
	ticker, err := bittrex.GetTicker("BTC-USDT")
	fmt.Println(err, ticker)

	// Get Distribution (JBS)
	distribution, err := bittrex.GetDistribution("BTC-USDS")
	for _, balance := range distribution.Distribution {
		fmt.Println(balance.BalanceD)
	}

	// Get market summaries
	//marketSummaries, err := bittrex.GetMarketSummaries()
	//fmt.Println(err, marketSummaries)

	// Get market summary
	marketSummary, err := bittrex.GetMarketSummary("BTC-ETH")
	fmt.Println(err, marketSummary)

	// Get orders book
	orderBook, err := bittrex.GetOrderBook("BTC-USD", "both")
	fmt.Println(err, orderBook)

	// Get order book buy or sell side only
	orderb, err := bittrex.GetOrderBookBuySell("BTC-USD", "buy")
	fmt.Println(err, orderb)

	// Market history
	marketHistory, err := bittrex.GetMarketHistory("BTC-XRP")
	for _, trade := range marketHistory {
		fmt.Println(err, trade.Timestamp.String(), trade.OrderUuid, trade.OrderType, trade.FillType, trade.Quantity, trade.Price)
	}

	// Market

	// BuyLimit
	/*
		uuid, err := bittrex.BuyLimit("BTC-DOGE", 1000, 0.00000102)
		fmt.Println(err, uuid)
	*/

	// Sell limit
	/*
		uuid, err := bittrex.SellLimit("BTC-DOGE", 1000, 0.00000115)
		fmt.Println(err, uuid)
	*/

	// Cancel Order
	/*
		err := bittrex.CancelOrder("e3b4b704-2aca-4b8c-8272-50fada7de474")
		fmt.Println(err)
	*/

	// Get open orders
	/*
		orders, err := bittrex.GetOpenOrders("BTC-DOGE")
		fmt.Println(err, orders)
	*/

	// Account
	// Get balances
	/*
		balances, err := bittrex.GetBalances()
		fmt.Println(err, balances)
	*/

	// Get balance
	/*
		balance, err := bittrex.GetBalance("DOGE")
		fmt.Println(err, balance)
	*/

	// Get address
	/*
		address, err := bittrex.GetDepositAddress("QBC")
		fmt.Println(err, address)
	*/

	// WithDraw
	/*
		whitdrawUuid, err := bittrex.Withdraw("QYQeWgSnxwtTuW744z7Bs1xsgszWaFueQc", "QBC", 1.1)
		fmt.Println(err, whitdrawUuid)
	*/

	// Get order history
	/*
		orderHistory, err := bittrex.GetOrderHistory("BTC-DOGE")
		fmt.Println(err, orderHistory)
	*/

	// Get getwithdrawal history
	/*
		withdrawalHistory, err := bittrex.GetWithdrawalHistory("all")
		fmt.Println(err, withdrawalHistory)
	*/

	// Get deposit history
	/*
		deposits, err := bittrex.GetDepositHistory("all")
		fmt.Println(err, deposits)
	*/

}
