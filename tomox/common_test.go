package tomox

import (
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
		},
		{
			Order:       []byte("order2"),
			Trades:      []map[string]string{{"takerOrderHash": "hash21"}, {"takerOrderHash": "hash22"}},
		},
		{
			Order:       []byte("order3"),
			Trades:      []map[string]string{{"takerOrderHash": "hash31"}, {"takerOrderHash": "hash32"}},
		},
	}

	encodedData, err := EncodeTxMatchesBatch(TxMatchBatch{
		Data:      originalTxMatchesBatch,
	})
	if err != nil {
		t.Error("Failed to encode", err.Error())
	}

	txMatchesBatch, err := DecodeTxMatchesBatch(encodedData)
	if err != nil {
		t.Error("Failed to decode", err.Error())
	}

	eq := reflect.DeepEqual(originalTxMatchesBatch, txMatchesBatch.Data)
	if eq {
		t.Log("Awesome, encode and decode txMatchesBatch are correct")
	} else {
		t.Error("txMatchesBatch is different from originalTxMatchesBatch", "txMatchesBatch", txMatchesBatch, "originalTxMatchesBatch", originalTxMatchesBatch)
	}
}
