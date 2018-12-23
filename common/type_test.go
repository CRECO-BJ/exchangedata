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

var testExchanger = &Exchanger{
	Name: "testExchanger",
	Currencies: []*Currency{
		&Currency{Name: "Bitcoin", Abbr: "BTC"},
		&Currency{Name: "Cardano", Abbr: "ADA"},
		&Currency{Name: "Tether", Abbr: "USDT"},
		&Currency{Name: "Verge", Abbr: "XVG"},
		&Currency{Name: "NXT", Abbr: "NXT"},
		&Currency{Name: "UnikoinGold", Abbr: "UKG"},
		&Currency{Name: "Ethereum", Abbr: "ETH"},
		&Currency{Name: "Ignis", Abbr: "IGNIS"},
		&Currency{Name: "Sirin", Abbr: "SRN"},
		&Currency{Name: "Worldwide Asset Exchange", Abbr: "WAX"},
		&Currency{Name: "0x Protocol", Abbr: "ZRX"},
		&Currency{Name: "BLOCKv", Abbr: "VEE"},
	},
}

func TestGetCurrencyByName(t *testing.T) {
	for _, x := range testExchanger.Currencies {
		if testExchanger.GetCurrencyByName(x.Name) == nil {
			t.Fatalf("%s should be found", x.Name)
		}
	}
}
