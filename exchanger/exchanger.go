package exchanger

import (
	"log"
	"net/url"
	"os"

	_ "github.com/exchangedata/common"
)

const (
	ExStop = 0
	ExRun  = iota
)

// Exchanger ...
type Exchanger struct {
	ID            string
	Name          string
	Countries     []string
	WebAPIURL     url.URL
	WebAPIVersion string
	WssURL        url.URL
	WssVersion    string
	WebAPIs       []string
	WssAPIs       []string
	Timeout       int
	RateLimit     int
	UserAgetn     string
	Verbose       bool
	Symbols       []Symbol
	Proxy         url.URL

	Status int

	logger *log.Logger
	db *database.
}

func (e *Exchanger) NewLogger() *log.Logger {
	return log.New(os.Stdout, e.Name+": ", 0)
}

func isValidProxy(url url.URL) bool {
	return len(url.Hostname()) != 0
}

// ExWsControl ...
type ExControl interface {
	Start() error
	Stop() error

	Setup(*Exchanger) error
	HandleData(...interface{}) error
}
