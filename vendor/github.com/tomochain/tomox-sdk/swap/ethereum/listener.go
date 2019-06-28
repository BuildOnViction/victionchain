package ethereum

import (
	"context"
	"math/big"
	"time"

	ethereumTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

const (
	// time out 15 seconds
	timeout = 15
)

func (l *Listener) Start() error {

	logger.Info("EthereumListener starting")

	blockNumber, err := l.Storage.GetBlockToProcess(types.ChainEthereum)
	if err != nil {
		err = errors.Wrap(err, "Error getting ethereum block to process from DB")
		logger.Error(err)
		return err
	}

	// Check if connected to correct network
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout*time.Second))
	defer cancel()
	id, err := l.Client.NetworkID(ctx)
	if err != nil {
		err = errors.Wrap(err, "Error getting ethereum network ID")
		logger.Error(err)
		return err
	}

	if id.String() != l.NetworkID {
		logger.Error("Invalid network ID (have=%s, want=%s)", id.String(), l.NetworkID)
		return errors.Errorf("Invalid network ID (have=%s, want=%s)", id.String(), l.NetworkID)
	}

	go l.processBlocks(blockNumber)
	return nil
}

func (l *Listener) Stop() error {
	ethClient := l.Client.(*ethclient.Client)
	ethClient.Close()
	l.Enabled = false
	return nil
}

func (l *Listener) processBlocks(blockNumber uint64) {
	if blockNumber == 0 {
		logger.Info("Starting from the latest block")
	} else {
		logger.Infof("Starting from block %d", blockNumber)
	}

	// Time when last new block has been seen
	lastBlockSeen := time.Now()
	noBlockWarningLogged := false

	for {
		if l.Enabled == false {
			// stop listener
			break
		}

		// process next amount of confirmed blocks to make sure
		blockNumberAhead := blockNumber + l.ConfirmedBlockNumber

		block, err := l.getBlock(blockNumberAhead)
		if err != nil {
			logger.Errorf("Error getting block, blockNumber: %d", blockNumberAhead)
			time.Sleep(1 * time.Second)
			continue
		}

		// Block doesn't exist yet
		if block == nil {
			if time.Since(lastBlockSeen) > 3*time.Minute && !noBlockWarningLogged {
				logger.Warningf("No new block in more than 3 minutes")
				noBlockWarningLogged = true
			}

			time.Sleep(1 * time.Second)
			continue
		}

		// Reset counter when new block appears
		lastBlockSeen = time.Now()
		noBlockWarningLogged = false

		if block.NumberU64() == 0 {
			logger.Error("Ethereum node is not synced yet. Unable to process blocks")
			time.Sleep(30 * time.Second)
			continue
		}

		if l.TransactionHandler == nil {
			// waiting for handler
			time.Sleep(1 * time.Second)
			continue
		}

		// now process the current block after number of block confirmation
		block, err = l.getBlock(blockNumber)

		err = l.processBlock(block)
		if err != nil {
			logger.Errorf("Error processing block, blockNumber: %d, err: %v", block.NumberU64(), err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Persist block number
		err = l.Storage.SaveLastProcessedBlock(types.ChainEthereum, blockNumber)
		if err != nil {
			logger.Errorf("Error saving last processed block: %s", err)
			time.Sleep(1 * time.Second)
			// We continue to the next block
		}

		blockNumber = block.NumberU64() + 1
	}
}

// getBlock returns (nil, nil) if block has not been found (not exists yet)
func (l *Listener) getBlock(blockNumber uint64) (*ethereumTypes.Block, error) {
	var blockNumberInt *big.Int
	if blockNumber > 0 {
		blockNumberInt = big.NewInt(int64(blockNumber))
	}

	d := time.Now().Add(10 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	block, err := l.Client.BlockByNumber(ctx, blockNumberInt)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		err = errors.Wrap(err, "Error getting block from geth")
		logger.Errorf("Got err: %s, block: %d", err.Error(), blockNumberInt.String())
		return nil, err
	}

	return block, nil
}

func (l *Listener) processBlock(block *ethereumTypes.Block) error {

	// empty block
	if block == nil {
		return errors.New("Block is not found")
	}

	transactions := block.Transactions()

	logger.Infof("Start processing block %d, total transactions: %d", block.Number(), len(transactions))

	//blockTime := time.Unix(block.Time().Int64(), 0)
	//logger.Infof("Processing block: blockNumber:%d, blockTime:%v, transactions:%d",
	//	block.NumberU64(),
	//	blockTime,
	//	len(transactions),
	//)

	for _, transaction := range transactions {
		to := transaction.To()
		if to == nil {
			// Contract creation
			continue
		}

		// this is the address that we need to check in address association
		// server will store associate like ethereumAddress => userAddress
		tx := Transaction{
			Hash:     transaction.Hash().Hex(),
			ValueWei: transaction.Value(),
			To:       to.Hex(),
		}
		err := l.TransactionHandler(tx)
		if err != nil {
			logger.Errorf("Error processing transaction: %s", err.Error())
			return errors.Wrap(err, "Error processing transaction")
		}
	}

	logger.Infof("Processed block %d", block.Number())

	return nil
}
