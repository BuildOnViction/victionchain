package endpoints

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/contracts"
	"github.com/tomochain/tomox-sdk/contracts/contractsinterfaces"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/swap"
	"github.com/tomochain/tomox-sdk/swap/bitcoin"
	"github.com/tomochain/tomox-sdk/swap/ethereum"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type depositEndpoint struct {
	depositService interfaces.DepositService
	walletService  interfaces.WalletService
	txService      interfaces.TxService
}

func ServeDepositResource(
	r *mux.Router,
	depositService interfaces.DepositService,
	walletService interfaces.WalletService,
	txService interfaces.TxService,
) {

	e := &depositEndpoint{depositService, walletService, txService}
	r.HandleFunc("/deposit/schema", e.handleGetSchema).Methods("GET")
	r.HandleFunc("/deposit/generate-address", e.handleGenerateAddress).Methods("POST")
	r.HandleFunc("/deposit/history", e.handleGetHistory).Methods("GET")
	r.HandleFunc("/deposit/recovery-transaction", e.handleRecoveryTransaction).Methods("GET")

	// r.HandleFunc("/deposit/testws", e.handleTestWS).Methods("GET")

	// set event handler delegate to this service
	depositService.SetDelegate(e)

	// register websocket channel
	ws.RegisterChannel(ws.DepositChannel, e.ws)
}

// func (e *depositEndpoint) handleTestWS(w http.ResponseWriter, r *http.Request) {
// 	v := r.URL.Query()
// 	address := v.Get("address")
// 	addressAssociation := &types.AddressAssociationRecord{
// 		AssociatedAddress: address,
// 	}

// 	ws.SendDepositMessage(types.UPDATE_STATUS, common.HexToAddress(address), addressAssociation)

// }

func (e *depositEndpoint) handleGetSchema(w http.ResponseWriter, r *http.Request) {
	schemaVersion := e.depositService.GetSchemaVersion()
	schema := map[string]interface{}{
		"version": schemaVersion,
	}
	httputils.WriteJSON(w, http.StatusOK, schema)
}

func (e *depositEndpoint) handleGetHistory(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	chainStr := v.Get("chain")
	addr := v.Get("userAddress")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid User Address")
		return
	}

	var chain types.Chain
	err := chain.Scan([]byte(chainStr))
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Chain is not correct")
		return
	}

	associatedAddress := common.HexToAddress(addr)
	addressAssociation, err := e.depositService.GetAssociationByChainAssociatedAddress(chain, associatedAddress)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Can not get association history")
		return
	}

	// httputils.WriteJSON(w, http.StatusOK, addressAssociation)

	associationTransactions := []*types.AssociationTransactionResponse{}

	if addressAssociation != nil {
		for _, txEnvelope := range addressAssociation.TxEnvelopes {
			bytes := common.Hex2Bytes(txEnvelope)
			var associationTransaction types.AssociationTransaction
			err = rlp.DecodeBytes(bytes, &associationTransaction)
			if err != nil {
				continue
			}

			associationTransactions = append(associationTransactions, associationTransaction.GetJSON())
		}
	}

	// return client
	httputils.WriteJSON(w, http.StatusOK, associationTransactions)

}

func (e *depositEndpoint) handleGenerateAddress(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	chainStr := v.Get("chain")
	addr := v.Get("userAddress")

	pairAddresses := &types.PairAddresses{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(pairAddresses)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid pair addresses")
		return
	}
	defer r.Body.Close()

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid User Address")
		return
	}

	associatedAddress := common.HexToAddress(addr)

	var chain types.Chain
	err = chain.Scan([]byte(chainStr))
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Chain is not correct")
		return
	}

	associationAddress, _ := e.depositService.GetAssociationByChainAssociatedAddress(chain, associatedAddress)
	var addressStr string

	// if we are waiting for this address to be sent, just return it
	if associationAddress != nil && associationAddress.Status == types.PENDING {
		addressStr = associationAddress.Address
	} else {

		address, index, err := e.depositService.GenerateAddress(chain)

		if err != nil {
			logger.Error(err)
			httputils.WriteError(w, http.StatusInternalServerError, "Can not generate Address")
			return
		}

		// update association
		// after receiving and displaying the address generated by server
		err = e.depositService.SaveAssociationByChainAddress(chain, address, index, associatedAddress, pairAddresses)
		if err != nil {
			logger.Error(err)
			httputils.WriteError(w, http.StatusInternalServerError, "Can not save association")
			return
		}

		addressStr = address.Hex()
	}

	response := &types.GenerateAddressResponse{
		ProtocolVersion: swap.ProtocolVersion,
		Chain:           chain.String(),
		Address:         addressStr,
		Signer:          e.depositService.SignerPublicKey().Hex(),
	}

	httputils.WriteJSON(w, http.StatusOK, response)
}

// return Address association for testing first
func (e *depositEndpoint) handleRecoveryTransaction(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("userAddress")
	chainStr := v.Get("chain")
	var chain types.Chain
	err := chain.Scan([]byte(chainStr))
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Chain is not correct")
		return
	}

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid User Address")
		return
	}

	address := common.HexToAddress(addr)

	association, err := e.getAssociationByUserAddress(chain, address)
	// association, err := e.depositService.GetAssociationByChainAddress(chain, address)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Can not get address association")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, association)
}

// ws function handles incoming websocket messages on the order channel
func (e *depositEndpoint) ws(input interface{}, c *ws.Client) {
	msg := &types.WebsocketEvent{}

	bytes, _ := json.Marshal(input)
	if err := json.Unmarshal(bytes, &msg); err != nil {
		logger.Error(err)
		c.SendMessage(ws.DepositChannel, types.ERROR, err.Error())
	}

	switch msg.Type {
	case "NEW_DEPOSIT":
		e.handleNewDeposit(msg, c)
	default:
		log.Print("Response with error")
	}
}

// handleNewDeposit create an association and register user to channel to notify if success
func (e *depositEndpoint) handleNewDeposit(ev *types.WebsocketEvent, c *ws.Client) {
	// set default here
	p := &types.AddressAssociationWebsocketPayload{}

	bytes, err := json.Marshal(ev.Payload)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.DepositChannel, types.ERROR, err.Error())
		return
	}

	logger.Debugf("Payload: %v#", ev.Payload)

	err = json.Unmarshal(bytes, p)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.DepositChannel, types.ERROR, err.Error())
		return
	}

	associationAddress, _ := e.depositService.GetAssociationByChainAssociatedAddress(p.Chain, p.AssociatedAddress)
	var addressStr string

	// if we are waiting for this address to be sent, just return it
	if associationAddress != nil && associationAddress.Status == types.PENDING {
		addressStr = associationAddress.Address
	} else {

		// now do the left
		address, index, err := e.depositService.GenerateAddress(p.Chain)

		if err != nil {
			logger.Error(err)
			c.SendMessage(ws.DepositChannel, types.ERROR, err.Error())
			return
		}

		// update association
		// after receiving and displaying the address generated by server
		err = e.depositService.SaveAssociationByChainAddress(p.Chain, address, index, p.AssociatedAddress, p.PairAddresses)
		if err != nil {
			logger.Error(err)
			c.SendMessage(ws.DepositChannel, types.ERROR, err.Error())
			return
		}

		addressStr = address.Hex()
	}

	response := &types.GenerateAddressResponse{
		ProtocolVersion: swap.ProtocolVersion,
		Chain:           p.Chain.String(),
		Address:         addressStr,
		Signer:          e.depositService.SignerPublicKey().Hex(),
	}

	// register for this user
	ws.RegisterDepositConnection(p.AssociatedAddress, c)

	// send response to channel
	c.SendMessage(ws.DepositChannel, types.UPDATE, response)

}

func (e *depositEndpoint) processNewTransaction(queueTx *types.DepositTransaction, addressAssociation *types.AddressAssociationRecord) error {
	// Add transaction as processing. this should be end point or tokenService
	// if have many nodes, should use queue instead

	// if addressAssociation.Status == types.PENDING || addressAssociation.Status == types.SUCCESS {
	if addressAssociation.Status == types.SUCCESS {
		// is processing or processed, just skip
		logger.Debug("Transaction already processed, skipping")
		return nil
	}

	err := e.processTransaction(addressAssociation)
	if err != nil {
		return err
	}

	// add queue transaction
	err = e.depositService.QueueAdd(queueTx)
	if err != nil {
		return errors.Wrap(err, "Error adding transaction to the processing queue")
	}

	logger.Info("Transaction added to transaction queue: %v", queueTx)
	// Broadcast event to address stream using websocket
	logger.Infof("Broadcasting event: %v", queueTx)
	logger.Info("Transaction processed successfully")
	return nil
}

/***** events from engine ****/

// onNewBitcoinTransaction checks if transaction is valid and adds it to
// the transactions queue for TomochainAccountConfigurator to consume.
//
// Transaction added to transactions queue should be in a format described in
// types.DepositTransaction (especialy amounts). Pooling service should not have to deal with any
// conversions.
func (e *depositEndpoint) OnNewBitcoinTransaction(transaction bitcoin.Transaction) error {
	logger.Infof("Processing transaction: %v", transaction)

	// Let's check if tx is valid first.

	// Check if value is above minimum required
	minimumValueSat := e.depositService.MinimumValueSat()
	if transaction.ValueSat < minimumValueSat {
		logger.Debugf("Value is : %s, below minimum required amount: %s, skipping", transaction.ValueSat, minimumValueSat)
		return nil
	}

	addressTo := common.HexToAddress(transaction.To)

	addressAssociation, err := e.depositService.GetAssociationByChainAddress(types.ChainBitcoin, addressTo)

	if err != nil {
		logger.Errorf("Chain: %s, Got error: %v", types.ChainBitcoin, err)
		return nil
	}

	// there is no address association in the database
	if addressAssociation == nil {
		logger.Info("Transaction not found, skipping")
		return nil
	}

	logger.Infof("Got Association: %v", addressAssociation)

	// Add tx to the processing queue
	queueTx := &types.DepositTransaction{
		Chain:         types.ChainBitcoin,
		TransactionID: transaction.Hash,
		AssetCode:     types.AssetCodeBTC,
		PairName:      addressAssociation.PairName,
		// Amount in the base unit of currency.
		Amount:            transaction.ValueToWei(),
		AssociatedAddress: addressAssociation.AssociatedAddress,
	}

	return e.processNewTransaction(queueTx, addressAssociation)

}

// onNewEthereumTransaction checks if transaction is valid and adds it to
// the transactions queue for TomochainAccountConfigurator to consume.
//
// Transaction added to transactions queue should be in a format described in
// types.DepositTransaction (especialy amounts). Pooling service should not have to deal with any
// conversions.
func (e *depositEndpoint) OnNewEthereumTransaction(transaction ethereum.Transaction) error {
	//logger.Infof("Processing transaction: %v", transaction)

	// Let's check if tx is valid first.

	// Check if value is above minimum required
	minimumValueWei := e.depositService.MinimumValueWei()
	if transaction.ValueWei.Cmp(minimumValueWei) < 0 {
		//logger.Debugf("Value is : %s, below minimum required amount: %s, skipping", transaction.ValueWei, minimumValueWei)
		return nil
	}

	addressTo := common.HexToAddress(transaction.To)

	addressAssociation, err := e.depositService.GetAssociationByChainAddress(types.ChainEthereum, addressTo)

	if err != nil {
		logger.Errorf("Chain: %s, Got error: %v", types.ChainEthereum, err)
		return nil
	}

	// there is no address association in the database
	if addressAssociation == nil {
		//logger.Info("Deposit transaction not found, skipping")
		return nil
	}

	logger.Infof("Got Association: %v", addressAssociation)

	// Add tx to the processing queue
	queueTx := &types.DepositTransaction{
		Chain:         types.ChainEthereum,
		TransactionID: transaction.Hash,
		AssetCode:     types.AssetCodeETH,
		PairName:      addressAssociation.PairName,
		// Amount in the base unit of currency.
		// Amount:            transaction.ValueWei.String(),
		Amount:            transaction.ValueToWei(),
		AssociatedAddress: addressAssociation.AssociatedAddress,
	}

	return e.processNewTransaction(queueTx, addressAssociation)
}

func (e *depositEndpoint) processTokenTransaction(addressAssociation *types.AddressAssociationRecord, depositAmount *big.Int) error {
	// here we update the transaction status and send token via smart contract
	contractAddress := common.HexToAddress(addressAssociation.BaseTokenAddress)
	receiver := common.HexToAddress(addressAssociation.AssociatedAddress)

	token, err := contracts.NewToken(
		e.walletService,
		e.txService,
		contractAddress,
		e.depositService.EthereumClient(),
	)

	if err != nil {
		return errors.Errorf("Could not connect to token address: %s", err.Error())
	}

	logs := []*contractsinterfaces.TokenTransfer{}

	// now calculate tokenAmount
	tokenAmount := depositAmount

	if err != nil {
		return errors.Errorf("Could not convert to token amount to transfer: %s", err.Error())
	}

	done := make(chan bool)

	events, err := token.ListenToTransferEvents()
	if err != nil {
		return errors.Errorf("Could not open transfer events channel")
	}

	go func() {
		for {
			event := <-events
			logs = append(logs, event)
			done <- true
		}
	}()

	_, err = token.Transfer(receiver, tokenAmount)
	if err != nil {
		return errors.Errorf("Could not transfer tokens: %v", err)
	}

	<-done

	if len(logs) != 1 {
		logger.Error("Events log has not the correct length")
		return errors.Errorf("Events log has not the correct length")
	}

	parsedTransfer := logs[0]

	if parsedTransfer.To != receiver {
		logger.Error("Event 'To' field is not correct")
		return errors.Errorf("Event 'To' field is not correct")
	}
	if parsedTransfer.Value.Cmp(tokenAmount) != 0 {
		logger.Error("Event 'Amount' field is not correct")
		return errors.Errorf("Event 'Amount' field is not correct")
	}

	return nil
}

// OnSubmitTransaction when transaction is submitted
func (e *depositEndpoint) OnSubmitTransaction(chain types.Chain, destination string, associationTransaction *types.AssociationTransaction) error {

	logger.Infof("On submit associationTransaction: %s", associationTransaction.TransactionType)

	// Save tx to database
	publicKey := common.HexToAddress(destination)
	addressAssociation, err := e.depositService.GetAssociationByChainAssociatedAddress(chain, publicKey)
	if err != nil {
		logger.Error("Error getting association from DB")
		return err
	}

	address := common.HexToAddress(addressAssociation.Address)
	bytes, err := rlp.EncodeToBytes(associationTransaction)
	if err != nil {
		logger.Error("Error encode transaction to string")
		e.depositService.SaveAssociationStatusByChainAddress(addressAssociation, types.FAILED)
		return err
	}

	transaction := hex.EncodeToString(bytes)

	// only update if found, do not insert
	err = e.depositService.SaveDepositTransaction(chain, address, transaction)
	if err != nil {
		logger.Error("Error saving transaction to DB")
		e.depositService.SaveAssociationStatusByChainAddress(addressAssociation, types.FAILED)
		return err
	}

	if associationTransaction.TransactionType == types.CreateOffer {

		// first param is amount, second param is the price of coin, such as ether
		amount := associationTransaction.Params[0]
		// tokenPrice := associationTransaction.Params[1]
		// TODO: Implement deposit fee here if needed.
		// Update deposit amount based on token price here, token price is <= 1
		// depositAmount = amount * tokenPrice
		// amount * (1 - token price) is the deposit fee that the exchange will receive
		depositAmount := new(big.Int)
		depositAmount, ok := depositAmount.SetString(amount, 10)
		if !ok {
			e.depositService.SaveAssociationStatusByChainAddress(addressAssociation, types.FAILED)
			return errors.Errorf("Could not convert to token amount: %s", amount)
		}

		// now process token transfer
		err = e.processTokenTransaction(addressAssociation, depositAmount)
		// update success if there is no error
		if err != nil {
			e.depositService.SaveAssociationStatusByChainAddress(addressAssociation, types.FAILED)
		} else {
			e.depositService.SaveAssociationStatusByChainAddress(addressAssociation, types.SUCCESS)
		}
		return err
	}

	// just done processing, continue to broadcast
	return nil
}

func (e *depositEndpoint) OnTomochainAccountCreated(chain types.Chain, destination string) {
	publicKey := common.HexToAddress(destination)
	association, err := e.getAssociationByUserAddress(chain, publicKey)
	if err != nil {
		logger.Error("Error getting association")
		return
	}

	if association == nil {
		logger.Error("Association not found")
		return
	}
	// broast cast event association
	logger.Infof("Broasting event: %v", association)
}

func (e *depositEndpoint) OnExchanged(chain types.Chain, destination string) {
	publicKey := common.HexToAddress(destination)
	association, err := e.getAssociationByUserAddress(chain, publicKey)
	if err != nil {
		logger.Error("Error getting association")
		return
	}

	if association == nil {
		logger.Error("Association not found")
		return
	}

	logger.Infof("Broasting event: %v", association)
}

func (e *depositEndpoint) OnExchangedTimelocked(chain types.Chain, destination string, associationTransaction *types.AssociationTransaction) {
	publicKey := common.HexToAddress(destination)
	association, err := e.getAssociationByUserAddress(chain, publicKey)
	if err != nil {
		logger.Error("Error getting association")
		return
	}

	if association == nil {
		logger.Error("Association not found")
		return
	}

	// Save tx to database
	bytes, err := rlp.EncodeToBytes(associationTransaction)
	if err != nil {
		logger.Error("Error encode transaction to string")
		return
	}

	transaction := hex.EncodeToString(bytes)
	err = e.depositService.SaveDepositTransaction(chain, publicKey, transaction)
	if err != nil {
		logger.Error("Error saving unlock transaction to DB")
		return
	}

	logger.Infof("Broasting event: %v", association)
}

func (e *depositEndpoint) LoadAccountHandler(chain types.Chain, destination string) (*types.AddressAssociation, error) {
	publicKey := common.HexToAddress(destination)
	return e.getAssociationByUserAddress(chain, publicKey)
}

func (e *depositEndpoint) getAssociationByUserAddress(chain types.Chain, publicKey common.Address) (*types.AddressAssociation, error) {

	// get from associated address
	logger.Infof("Get association chain :%s, associatedAddress: %s", chain, publicKey.Hex())
	record, err := e.depositService.GetAssociationByChainAssociatedAddress(chain, publicKey)

	if err != nil {
		return nil, err
	}

	if record == nil {
		return nil, nil
	}

	addressAssociation, err := record.GetJSON()
	if err == nil {
		addressAssociation.TomochainPublicKey = e.depositService.SignerPublicKey()
	}

	return addressAssociation, err
}

func (e *depositEndpoint) processTransaction(addressAssociation *types.AddressAssociationRecord) error {

	// update status is pending, if is pending or success then return true, to tell it stop
	// otherwise return false, then update the status is pending, on success update it success,
	// on faile update it fail, so next time it will return false again to re-process
	// call smart contract and listen to event, then update to mongodb via DAO

	// re-set status as pending and wait for confirmation, this is even for retrying a failed transaction
	return e.depositService.SaveAssociationStatusByChainAddress(addressAssociation, types.PENDING)
}
