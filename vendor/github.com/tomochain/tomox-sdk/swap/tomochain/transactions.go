package tomochain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"strconv"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

func (ac *AccountConfigurator) createAccountTransaction(chain types.Chain, destination string) error {
	transaction, err := ac.buildTransaction(
		ac.signerPublicKey.String(),
		ac.signerPrivateKey,
		types.CreateAccount,
		destination,
		ac.StartingBalance,
	)
	if err != nil {
		return errors.Wrap(err, "Error building transaction")
	}

	err = ac.submitTransaction(chain, destination, transaction)
	if err != nil {
		return errors.Wrap(err, "Error submitting a transaction")
	}

	return nil
}

// configureAccountTransaction is using a signer on an user accounts to configure the account.
func (ac *AccountConfigurator) configureAccountTransaction(depositTransaction *types.DepositTransaction) error {

	var tokenPrice string
	switch depositTransaction.AssetCode {
	case types.AssetCodeETH:
		tokenPrice = ac.TokenPriceETH
		break
	case types.AssetCodeBTC:
		tokenPrice = ac.TokenPriceBTC
		break
	default:
		return errors.Errorf("Invalid intermediateAssetCode: $%s", depositTransaction.AssetCode)
	}

	// // Send WETH token using smart contract
	// exchange by rate from regulator service
	transaction, err := ac.buildTransaction(
		depositTransaction.AssociatedAddress,
		ac.signerPrivateKey,
		types.CreateOffer,
		// 3 fields below are "params"
		depositTransaction.Amount,
		tokenPrice,
		depositTransaction.PairName,
	)
	if err != nil {
		return errors.Wrap(err, "Error building a transaction")
	}

	err = ac.submitTransaction(depositTransaction.Chain, depositTransaction.AssociatedAddress, transaction)
	if err != nil {
		return errors.Wrap(err, "Error submitting a transaction")
	}

	return nil
}

// removeTemporarySigner is removing temporary signer from an account.
func (ac *AccountConfigurator) removeTemporarySigner(chain types.Chain, destination string) error {
	// Remove signer ? need to remove this account wallet? ac.signerPublicKey

	transaction, err := ac.buildTransaction(destination, ac.signerPrivateKey, types.RemoveSigner)
	if err != nil {
		return errors.Wrap(err, "Error building a transaction")
	}

	err = ac.submitTransaction(chain, destination, transaction)
	if err != nil {
		return errors.Wrap(err, "Error submitting a transaction")
	}

	return nil
}

// buildUnlockAccountTransaction creates and returns unlock account transaction.
func (ac *AccountConfigurator) buildUnlockAccountTransaction(source string) (*types.AssociationTransaction, error) {
	// Remove signer, ac.LockUnixTimestamp
	// do not submit and wait for MinTime: LockUnixTimestamp
	return ac.buildTransaction(source, ac.signerPrivateKey, types.RemoveSigner, strconv.FormatUint(ac.LockUnixTimestamp, 10))
}

// this will create hex data of rlp encode data
func (ac *AccountConfigurator) buildTransaction(source string, signer *ecdsa.PrivateKey, transactionType string, params ...string) (*types.AssociationTransaction, error) {

	associationTransaction := &types.AssociationTransaction{
		Source:          source,
		TransactionType: transactionType,
		Params:          params,
	}

	associationTransaction.Hash = associationTransaction.ComputeHash()

	// TODO: If relayer retrieves association from feed to restore sending phrase
	// it has to verify the signature to make sure it is correct
	// and client must update this signature in case of recovery
	signature, err := signer.Sign(rand.Reader, associationTransaction.Hash, nil)
	if err != nil {
		return nil, err
	}

	associationTransaction.Signature = signature

	return associationTransaction, nil
}

func (ac *AccountConfigurator) submitTransaction(chain types.Chain, destination string, transaction *types.AssociationTransaction) error {
	logger.Info("Submitting transaction")

	// no implementation, just return
	if ac.OnSubmitTransaction == nil {
		return nil
	}

	// call storage update transaction for association

	err := ac.OnSubmitTransaction(chain, destination, transaction)

	if err != nil {

		logger.Error("Error submitting transaction")
		return errors.Wrap(err, "Error submitting transaction")
	}

	logger.Info("Transaction successfully submitted")
	return nil
}
