// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tomox_state

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

type DumpExchange struct {
	Nonce        uint64            `json:"nonce"`
	AskRoot      string            `json:"askRoot"`
	BidRoot      string            `json:"bidRoot"`
	AskStorage   map[string]string `json:"askStorage"`
	BidStorage   map[string]string `json:"bidStorage"`
}

type Dump struct {
	Root     string                  `json:"root"`
	Accounts map[string]DumpExchange `json:"exchanges"`
}

func (self *TomoXStateDB) RawDump() Dump {
	dump := Dump{
		Root:     fmt.Sprintf("%x", self.trie.Hash()),
		Accounts: make(map[string]DumpExchange),
	}

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		var data exchangeObject
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			panic(err)
		}

		obj := newStateExchanges(nil, common.BytesToHash(addr), data, nil)
		account := DumpExchange{
			Nonce:        data.Nonce,
			AskRoot:      common.Bytes2Hex(data.AskRoot[:]),
			BidRoot:      common.Bytes2Hex(data.BidRoot[:]),
			AskStorage:   make(map[string]string),
			BidStorage:   make(map[string]string),
		}
		storageIt := trie.NewIterator(obj.getAsksTrie(self.db).NodeIterator(nil))
		for storageIt.Next() {
			account.AskStorage[common.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = common.Bytes2Hex(storageIt.Value)
		}
		storageIt = trie.NewIterator(obj.getBidsTrie(self.db).NodeIterator(nil))
		for storageIt.Next() {
			account.BidStorage[common.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = common.Bytes2Hex(storageIt.Value)
		}
		dump.Accounts[common.Bytes2Hex(addr)] = account
	}
	return dump
}

func (self *TomoXStateDB) Dump() []byte {
	json, err := json.MarshalIndent(self.RawDump(), "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}

	return json
}
