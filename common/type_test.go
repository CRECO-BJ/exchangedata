package common

import (
	"testing"
)

type testData struct {
	in, want, info string
}

var td = []testData{
	{"BTC", "", "no seperator"},
	{"", "", ""},
}

func TestSymbolPars(t *testing.T) {
	_ = td
}
