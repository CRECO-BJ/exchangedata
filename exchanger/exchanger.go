package exchanger

import (
	"net/url"
)

// Exchanger ...
type Exchanger struct {
	ID         string
	Name       string
	Countries  []string
	apiURL     url.URL
	apiVersion string
	wssURL     url.URL
	wssVersion string
	APIs       []string
	wssAPIs    []string
	Timeout    int
	RateLimit  int
	UserAgetn  string
	Verbose    bool
	Markets    []Market
	Symbols    []Symbol
	Currencies []Currency
	Proxy      url.URL
}
