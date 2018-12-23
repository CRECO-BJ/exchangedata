package bittrex

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/exchangedata/common"
	"github.com/exchangedata/database"
	"github.com/urfave/cli"
)

var (
	API_KEY    = ""
	API_SECRET = ""
)

func exampleMain() {
	app := cli.NewApp()

	flags := []cli.Flag{
		cli.StringFlag{Name: "APIKey"},
		cli.StringFlag{Name: "APISecret"},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("yaml ist rad")
		return nil
	}

	app.Flags = flags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	BittrexWebURL               = "https://Bittrex.com"
	BittrexWebTradingEndpoint   = "tradingApi"
	BittrexWebVersion           = "1"
	BittrexTradeHistory         = "returnTradeHistory"
	BittrexBalances             = "returnBalances"
	BittrexBalancesComplete     = "returnCompleteBalances"
	BittrexDepositAddresses     = "returnDepositAddresses"
	BittrexGenerateNewAddress   = "generateNewAddress"
	BittrexDepositsWithdrawals  = "returnDepositsWithdrawals"
	BittrexOrders               = "returnOpenOrders"
	BittrexOrderBuy             = "buy"
	BittrexOrderSell            = "sell"
	BittrexOrderCancel          = "cancelOrder"
	BittrexOrderMove            = "moveOrder"
	BittrexWithdraw             = "withdraw"
	BittrexFeeInfo              = "returnFeeInfo"
	BittrexAvailableBalances    = "returnAvailableAccountBalances"
	BittrexTradableBalances     = "returnTradableBalances"
	BittrexTransferBalance      = "transferBalance"
	BittrexMarginAccountSummary = "returnMarginAccountSummary"
	BittrexMarginBuy            = "marginBuy"
	BittrexMarginSell           = "marginSell"
	BittrexMarginPosition       = "getMarginPosition"
	BittrexMarginPositionClose  = "closeMarginPosition"
	BittrexCreateLoanOffer      = "createLoanOffer"
	BittrexCancelLoanOffer      = "cancelLoanOffer"
	BittrexOpenLoanOffers       = "returnOpenLoanOffers"
	BittrexActiveLoans          = "returnActiveLoans"
	BittrexLendingHistory       = "returnLendingHistory"
	BittrexAutoRenew            = "toggleAutoRenew"

	BittrexAuthRate   = 6
	BittrexUnauthRate = 6
)

// Bittrex struct
type Bittrex struct {
	ex     *common.Exchanger
	ds     *database.DataStore
	client *client
	logger *log.Logger

	done chan struct{} // Closed when the receive rountine received error, then the main exchanger communication routine exit
	// If CloseDone is not closed, the connection should be reconnected...ToDo
	stop chan struct{} // Signal to close connection and exit. Program exiting...
}

func NewBittrex() *Bittrex {
	b := &Bittrex{
		ex: &common.Exchanger{Name: "bittrex"},
	}
	b.NewLogger()
	b.done = make(chan struct{})
	b.stop = make(chan struct{}, 1)
	b.client = NewClient(API_KEY, API_SECRET)
	return b
}

func (b *Bittrex) NewLogger() {
	b.logger = log.New(os.Stdout, "Bittrex:", log.LstdFlags)
}

func (b *Bittrex) Logf(format string, v ...interface{}) {
	b.logger.Printf(format, v)
}

func (b *Bittrex) Logln(v ...interface{}) {
	b.logger.Println(v)
}

func (b *Bittrex) Panicf(format string, v ...interface{}) {
	b.logger.Panicf(format, v)
}

func (b *Bittrex) Panic(v ...interface{}) {
	b.logger.Panic(v)
}

// Setup prepares the basic data for startup and main duty loop
func (b *Bittrex) Setup() error {
	var err error
	b.ds = database.NewDataStore("mysql")
	if err = b.ds.OpenDB(); err != nil {
		b.Panic("open db failed")
	}
	b.ds.AutoMigrate()

	if currencies, err := b.GetCurrencies(); err != nil {
		b.Logln("error get currency ", err)
		return err
	} else {
		if b.ex.Currencies == nil {
			b.ex.Currencies = []*common.Currency{}
		}
		for _, c := range currencies {
			n := &common.Currency{Name: c.CurrencyLong, Abbr: c.Currency} //Exchangers: []*common.Exchanger{b.ex}} // not sure why this panic. ToKnow
			if b.ds.UpdateCurrency(n).Error != nil {
				b.Logln("error update db, currency ", n, b.ds.GetDB().Error)
				return b.ds.GetDB().Error
			}
			b.ex.Currencies = append(b.ex.Currencies, n)
		}
	}

	if markets, err := b.GetMarkets(); err != nil {
		b.Logln("error get market ", err)
		return err
	} else {
		if b.ex.Markets == nil {
			b.ex.Markets = []*common.Market{}
		}
		for _, c := range markets {
			base := b.GetCurrencyByName(c.BaseCurrencyLong)
			if base == nil {
				b.Logf("Error! currency %s should be founded!", c.BaseCurrencyLong)
				continue
			}
			quote := b.GetCurrencyByName(c.MarketCurrencyLong)
			if quote == nil {
				b.Logf("Error! currency %s should be founded!", c.MarketCurrencyLong)
				continue
			}
			sym := &common.Symbol{
				Base:  base,
				Quote: quote,
			}
			m := &common.Market{
				Name:   sym.String(),
				Symbol: sym,
				Active: c.IsActive,
			}
			if b.ds.UpdateMarket(m).Error != nil {
				b.Logln("error update db, market ", c, b.ds.GetDB().Error)
				return b.ds.GetDB().Error
			}
		}
	}

	if b.ds.UpdateExchanger(b.ex).Error != nil {
		b.Logln("error update db, exchanger ", b.ex, b.ds.GetDB().Error)
		return b.ds.GetDB().Error
	}
	return nil
}

// Start ...
func (b *Bittrex) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	b.Logln("bittrex Started ...")
	ticker := time.NewTicker(5 * time.Second) // default is to get ticker every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // timely keepAlive processing
			b.runDataFetcher()
		case <-b.stop:
			close(b.done)
			return
		}
	}
}

// Stop ...
func (b *Bittrex) Stop() {
	b.stop <- struct{}{}
	<-b.done
	close(b.stop)
}

// Done ...
func (b *Bittrex) Done() {
	b.done <- struct{}{}
}

func (b *Bittrex) runDataFetcher() (err error) {
	b.Logln("runDataFetcher ...")
	for _, m := range b.ex.Markets {
		name := m.Symbol.String()
		b.Logln("Get ", name, " data ...")
		// Get Ticker
		ticker, err := b.GetTicker(name)
		if err != nil {
			b.Logln(err, ticker)
		}
		// Get market summary
		marketSummary, err := b.GetMarketSummary(name)
		if err != nil {
			b.Logln(err, marketSummary)
		}

		// Get orders book
		orderBook, err := b.GetOrderBook(name, "both")
		if err != nil {
			b.Logln(err, orderBook)
		}

		// Market history
		marketHistory, err := b.GetMarketHistory(name)
		if err != nil {
			for _, trade := range marketHistory {
				b.Logln(err, trade.Timestamp.String(), trade.OrderUuid, trade.OrderType, trade.FillType, trade.Quantity, trade.Price)
			}
		}

		// Get Distribution
		distribution, err := b.GetDistribution(name)
		if err != nil {
			for _, balance := range distribution.Distribution {
				b.Logln(balance.BalanceD)
			}
		}
	}
	return
}

func (b *Bittrex) GetCurrencyByName(name string) *common.Currency {
	return b.ex.GetCurrencyByName(name)
}
