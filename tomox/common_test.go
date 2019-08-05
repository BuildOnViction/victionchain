package tomox

import (
	"github.com/ethereum/go-ethereum/common"
	"reflect"
	"testing"
)

// Testing scenario:
// encode originalTxMatchesBatch -> byteData
// decode byteData -> txMatchesBatch
// compare originalTxMatchesBatch and txMatchesBatch
func TestTxMatchesBatch(t *testing.T) {
	originalTxMatchesBatch := []TxDataMatch{
		{
			Order:       []byte("order1"),
			Trades:      []map[string]string{{"takerOrderHash": "hash11"}, {"takerOrderHash": "hash12"}},
			OrderInBook: []byte("order1'"),
			ObOld:       common.StringToHash("ObOld1"),
			ObNew:       common.StringToHash("ObNew1"),
			AskOld:      common.StringToHash("AskOld1"),
			AskNew:      common.StringToHash("AskNew1"),
			BidOld:      common.StringToHash("BidOld1"),
			BidNew:      common.StringToHash("BidNew1"),
		},
		{
			Order:       []byte("order2"),
			Trades:      []map[string]string{{"takerOrderHash": "hash21"}, {"takerOrderHash": "hash22"}},
			OrderInBook: []byte("order2'"),
			ObOld:       common.StringToHash("ObOld2"),
			ObNew:       common.StringToHash("ObNew2"),
			AskOld:      common.StringToHash("AskOld2"),
			AskNew:      common.StringToHash("AskNew2"),
			BidOld:      common.StringToHash("BidOld2"),
			BidNew:      common.StringToHash("BidNew2"),
		},
		{
			Order:       []byte("order3"),
			Trades:      []map[string]string{{"takerOrderHash": "hash31"}, {"takerOrderHash": "hash32"}},
			OrderInBook: []byte("order3'"),
			ObOld:       common.StringToHash("ObOld3"),
			ObNew:       common.StringToHash("ObNew3"),
			AskOld:      common.StringToHash("AskOld3"),
			AskNew:      common.StringToHash("AskNew3"),
			BidOld:      common.StringToHash("BidOld3"),
			BidNew:      common.StringToHash("BidNew3"),
		},
	}

	encodedData, err := EncodeTxMatchesBatch(originalTxMatchesBatch)
	if err != nil {
		t.Error("Failed to encode", err.Error())
	}

	txMatchesBatch, err := DecodeTxMatchesBatch(encodedData)
	if err != nil {
		t.Error("Failed to decode", err.Error())
	}

	eq := reflect.DeepEqual(originalTxMatchesBatch, txMatchesBatch)
	if eq {
		t.Log("Awesome, encode and decode txMatchesBatch are correct")
	} else {
		t.Error("txMatchesBatch is different from originalTxMatchesBatch", "txMatchesBatch", txMatchesBatch, "originalTxMatchesBatch", originalTxMatchesBatch)
	}
}
