package bitcoin

import (
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

func (l *Listener) Start() error {

	logger.Info("BitcoinListener starting")

	genesisBlockHash, err := l.Client.GetBlockHash(0)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "Error getting genesis block")
	}

	if l.Testnet {
		l.chainParams = &chaincfg.TestNet3Params
	} else {
		l.chainParams = &chaincfg.MainNetParams
	}

	if !genesisBlockHash.IsEqual(l.chainParams.GenesisHash) {
		return errors.New("Invalid genesis hash")
	}

	blockNumber, err := l.Storage.GetBlockToProcess(types.ChainBitcoin)
	if err != nil {
		err = errors.Wrap(err, "Error getting bitcoin block to process from DB")
		logger.Error(err)
		return err
	}

	if blockNumber == 0 {
		blockNumberTmp, err := l.Client.GetBlockCount()
		if err != nil {
			err = errors.Wrap(err, "Error getting the block count from bitcoin-core")
			logger.Error(err)
			return err
		}
		blockNumber = uint64(blockNumberTmp)
	}

	go l.processBlocks(blockNumber)
	return nil
}

func (l *Listener) processBlocks(blockNumber uint64) {
	logger.Infof("Starting from block %d", blockNumber)

	// Time when last new block has been seen
	lastBlockSeen := time.Now()
	missingBlockWarningLogged := false

	for {
		// process next amount of confirmed blocks to make sure
		blockNumberAhead := blockNumber + l.ConfirmedBlockNumber

		block, err := l.getBlock(blockNumberAhead)
		if err != nil {
			logger.Errorf("Error getting blockerr, err:%v, blockNumber: %d", err, blockNumberAhead)
			time.Sleep(time.Second)
			continue
		}

		// Block doesn't exist yet
		if block == nil {
			if time.Since(lastBlockSeen) > 20*time.Minute && !missingBlockWarningLogged {
				logger.Warning("No new block in more than 20 minutes")
				missingBlockWarningLogged = true
			}

			time.Sleep(time.Second)
			continue
		}

		// Reset counter when new block appears
		lastBlockSeen = time.Now()
		missingBlockWarningLogged = false

		// now process the current block after number of block confirmation
		err = l.processBlock(block)
		if err != nil {
			logger.Errorf("Error processing block, err: %v, blockHash: %s", err, block.Header.BlockHash().String())
			time.Sleep(time.Second)
			continue
		}

		// Persist block number
		err = l.Storage.SaveLastProcessedBlock(types.ChainBitcoin, blockNumber)
		if err != nil {
			logger.Errorf("Error saving last processed block, err: %v", err)
			time.Sleep(time.Second)
			// We continue to the next block.
			// The idea behind this is if there was a problem with this single query we want to
			// continue processing because it's safe to reprocess blocks and we don't want a downtime.
		}

		blockNumber++
	}
}

// getBlock returns (nil, nil) if block has not been found (not exists yet)
func (l *Listener) getBlock(blockNumber uint64) (*wire.MsgBlock, error) {
	blockHeight := int64(blockNumber)
	blockHash, err := l.Client.GetBlockHash(blockHeight)
	if err != nil {
		if strings.Contains(err.Error(), "Block height out of range") {
			// Block does not exist yet
			return nil, nil
		}
		err = errors.Wrap(err, "Error getting block hash from bitcoin-core")
		logger.Error(err, "blockHeight", blockHeight)
		return nil, err
	}

	block, err := l.Client.GetBlock(blockHash)
	if err != nil {
		err = errors.Wrap(err, "Error getting block from bitcoin-core")
		logger.Error(err, "blockHash", blockHash.String())
		return nil, err
	}

	return block, nil
}

func (l *Listener) processBlock(block *wire.MsgBlock) error {
	transactions := block.Transactions

	logger.Infof("Processing block: blockHash:%s, blockTime:%v, transactions:%d",
		block.Header.BlockHash().String(),
		block.Header.Timestamp,
		len(transactions),
	)

	for _, transaction := range transactions {
		transactionHash := transaction.TxHash().String()

		for index, output := range transaction.TxOut {
			class, addresses, _, err := txscript.ExtractPkScriptAddrs(output.PkScript, l.chainParams)
			if err != nil {
				// txscript.ExtractPkScriptAddrs returns error on non-standard scripts
				// so this can be Warn.
				logger.Warningf("Error extracting addresses, transactionHash :%s, err:%v", transactionHash, err)

				continue
			}

			// We only support P2PK and P2PKH addresses
			if class != txscript.PubKeyTy && class != txscript.PubKeyHashTy {
				logger.Debugf("Unsupported addresses class, transactionHash :%s, class:%v", transactionHash, class)

				continue
			}

			// Paranoid. We access address[0] later.
			if len(addresses) != 1 {
				logger.Errorf("Invalid addresses length, transactionHash :%s, addresses:%v", transactionHash, addresses)

				continue
			}

			handlerTransaction := Transaction{
				Hash:       transaction.TxHash().String(),
				TxOutIndex: index,
				ValueSat:   output.Value,
				To:         addresses[0].EncodeAddress(),
			}

			err = l.TransactionHandler(handlerTransaction)
			if err != nil {
				return errors.Wrap(err, "Error processing transaction")
			}
		}
	}

	// done process
	logger.Info("Processed block")

	return nil
}
