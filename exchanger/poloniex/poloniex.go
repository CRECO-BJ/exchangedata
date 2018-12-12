package poloniex

import (
	"sync"
	"github.com/exchangedata/exchanger"
	"github.com/gorilla/websocket"
)

const (
	poloniexWebURL               = "https://poloniex.com"
	poloniexWebTradingEndpoint   = "tradingApi"
	poloniexWebVersion           = "1"
	poloniexTradeHistory         = "returnTradeHistory"
	poloniexBalances             = "returnBalances"
	poloniexBalancesComplete     = "returnCompleteBalances"
	poloniexDepositAddresses     = "returnDepositAddresses"
	poloniexGenerateNewAddress   = "generateNewAddress"
	poloniexDepositsWithdrawals  = "returnDepositsWithdrawals"
	poloniexOrders               = "returnOpenOrders"
	poloniexOrderBuy             = "buy"
	poloniexOrderSell            = "sell"
	poloniexOrderCancel          = "cancelOrder"
	poloniexOrderMove            = "moveOrder"
	poloniexWithdraw             = "withdraw"
	poloniexFeeInfo              = "returnFeeInfo"
	poloniexAvailableBalances    = "returnAvailableAccountBalances"
	poloniexTradableBalances     = "returnTradableBalances"
	poloniexTransferBalance      = "transferBalance"
	poloniexMarginAccountSummary = "returnMarginAccountSummary"
	poloniexMarginBuy            = "marginBuy"
	poloniexMarginSell           = "marginSell"
	poloniexMarginPosition       = "getMarginPosition"
	poloniexMarginPositionClose  = "closeMarginPosition"
	poloniexCreateLoanOffer      = "createLoanOffer"
	poloniexCancelLoanOffer      = "cancelLoanOffer"
	poloniexOpenLoanOffers       = "returnOpenLoanOffers"
	poloniexActiveLoans          = "returnActiveLoans"
	poloniexLendingHistory       = "returnLendingHistory"
	poloniexAutoRenew            = "toggleAutoRenew"

	poloniexAuthRate   = 6
	poloniexUnauthRate = 6
)


// Poloniex struct
type Poloniex struct {
	exchanger.Exchanger

	conn *websocket.Conn

	wsEx *sync.WaitGroup
	done chan struct{} // Closed when the receive rountine received error, then the main exchanger communication routine exit
	// If CloseDone is not closed, the connection should be reconnected...ToDo
	stop chan struct{} // Signal to close connection and exit. Program exiting...
}

func NewPoloniex() *Poloniex {
	p := &Poloniex{}
	p.wsEx = &sync.WaitGroup
	p.done = make( chan struct{})
	p.stop = make( chan struct{})
	return p
}

// Start ...
func (p *Poloniex) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(p.keepAlive * time.Second)
	defer ticker.Stop()
	exit := false
	logger := p.NewLogger()
start:
	if p.UseWss() {
		p.Connect()
	} else if !p.UseWeb() { // no communication url is defined, bypass
		logger.Println("no valid communication method", e)
		return
	}

	go func() {
		defer p.Done()
		for {
			_, message, err := p.conn.ReadMessage()
			if err != nil { // if read error, rountine exit, redail
				logger.Println("read:", err)
				return
			}
			p.HandleMessage(message)
		}
	}()

	p.SubScribeWss()

	for {
		select {
		case <-p.done:
			if exit == true {
				logger.Println("rountine exited")
				return
			}
			goto start
		case <-ticker.C: // timely keepAlive processing
			//			err := p.conn.WriteMessage(websocket.TextMessage, []byte("test"))
			//			if err != nil {
			//				logger.Println( "write:", err)
			//				return
			//			}
			logger.Println("time to write:", "test")
		case <-p.stop:
			logger.Println("Close connection and to exit")
			exit = true
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := p.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Println("write close with error:", err)
				return
			}
			select {
			case <-p.done:
			case <-time.After(time.Second):
			}
			logger.Println("rountine exited")
			return
		}
	}
}

// Connect ...
func (p *Poloniex) Connect() *websocket.Conn {
	var dialer *websocket.Dialer
	loger := p.NewLogger()

	if isValidProxy(p.Proxy) {
		//	auth = proxy.Auth(User:, Password:kb109kb109)
		netDialer, err := proxy.SOCKS5("udp", p.Proxy.String(), nil, proxy.Direct)
		if err != nil {
			loger.Fatalf("sock5 configuration error %v", err)
		}
		dialer = &websocket.Dialer{NetDial: netDialer.Dial}
	} else {
		dialer = websocket.DefaultDialer
	}
	c, _, err := dialer.Dial(p.WssURL.String(), nil)
	if err != nil {
		loger.Fatalf("dial:%s error:%v", p.WssURL.String(), err)
	}
	p.conn = c
	return p.conn
}

// Stop ...
func (e *Poloniex) Stop() {
	e.stop <- struct{}{}
}

// Done ...
func (e *Poloniex) Done() {
	e.done <- struct{}{}
}

