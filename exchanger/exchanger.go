package exchanger

import (
	"net/url"
	"sync"

	_ "github.com/exchangedata/common"
)

const (
	ExStop = 0
	ExRun  = iota
)

// Exchanger communication configuration
type ExchangerConf struct {
	WebAPIURL     url.URL
	WebAPIVersion string
	WssURL        url.URL
	WssVersion    string
	WebAPIs       []string
	WssAPIs       []string
	Timeout       int
	RateLimit     int
	Verbose       bool
	Proxy         url.URL

	Status int
}

// ExWsControl ...
type ExControl interface {
	Setup() error
	Start(wg *sync.WaitGroup)
	Stop()
}

func isValidProxy(url url.URL) bool {
	return len(url.Hostname()) != 0
}
