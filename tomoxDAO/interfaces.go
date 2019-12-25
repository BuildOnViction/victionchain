// Copyright 2019 The Tomochain Authors
// This file is part of the Core Tomochain infrastructure
// https://tomochain.com
// Package tomoxDAO provides an interface to work with tomox database, including leveldb for masternode and mongodb for SDK node
package tomoxDAO

import (
	"github.com/globalsign/mgo"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
)

const defaultCacheLimit = 1024

type TomoXDAO interface {
	// for both leveldb and mongodb
	IsEmptyKey(key []byte) bool
	Close()

	// mongodb methods
	HasObject(hash common.Hash, val interface{}) (bool, error)
	GetObject(hash common.Hash, val interface{}) (interface{}, error)
	PutObject(hash common.Hash, val interface{}) error
	DeleteObject(hash common.Hash, val interface{}) error // won't return error if key not found

		// basic tomox
		InitBulk() *mgo.Session
		CommitBulk() error
		GetOrderByTxHash(txhash common.Hash) []*tradingstate.OrderItem
		GetListOrderByHashes(hashes []string) []*tradingstate.OrderItem
		DeleteTradeByTxHash(txhash common.Hash)

		// tomox lending
		InitLendingBulk() *mgo.Session
		CommitLendingBulk() error
		GetLendingItemByTxHash(txhash common.Hash) []*lendingstate.LendingItem
		GetListLendingItemByHashes(hashes []string) []*lendingstate.LendingItem
		DeleteLendingTradeByTxHash(txhash common.Hash)

	// leveldb methods
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	NewBatch() ethdb.Batch
}
