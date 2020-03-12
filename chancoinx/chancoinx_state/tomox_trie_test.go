package chancoinx_state

import (
	"fmt"
	"github.com/chancoin-core/chancoin-gold/common"
	"github.com/chancoin-core/chancoin-gold/ethdb"
	"math/big"
	"testing"
)

func TestChancoinxTrieTest(t *testing.T) {
	db, _ := ethdb.NewMemDatabase()
	stateCache := NewDatabase(db)
	trie, _ := stateCache.OpenStorageTrie(EmptyHash, EmptyHash)
	min := common.BigToHash(big.NewInt(1)).Bytes()
	max := common.BigToHash(big.NewInt(2)).Bytes()
	trie.TryUpdate(min, min)
	trie.TryUpdate(max, max)
	left, _, _ := trie.TryGetBestLeftKeyAndValue()
	right, _, _ := trie.TryGetBestRightKeyAndValue()
	fmt.Println(left, right)
	trie.TryDelete(min)
	left, _, _ = trie.TryGetBestLeftKeyAndValue()
	right, _, _ = trie.TryGetBestRightKeyAndValue()
	fmt.Println(left, right)
}
