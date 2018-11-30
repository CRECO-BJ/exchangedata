package protocol

import (
	"testing"
	"time"
)

func TestOkexParam(t *testing.T) {
	val := OkexParams{apiKey: "", sign: "123456789abcdefgh", timestamp: time.Now(), passphrase: "test1234"}
	t.Errorf("val is: [%s]", val.String())
}

func TestOkexRequest(t *testing.T) {
	//	val := OkexRequest{}
}
