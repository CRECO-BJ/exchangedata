package common

import "testing"

func TestSymbol(t *testing.T) {
	testString1 := "BTC_LTC"
	sym1 := Symbol{}
	err := sym1.ParseString(testString1)
	if err != nil {
		t.Fatal("ParseString error ", testString1)
	}
	if sym1.String() != testString1 {
		t.Fatal("Parse then string error ", testString1)
	}
}
