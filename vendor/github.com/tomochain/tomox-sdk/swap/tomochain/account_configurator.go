package tomochain

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

func (ac *AccountConfigurator) Start() error {

	logger.Info("TomochainAccountConfigurator starting")

	if !common.IsHexAddress(ac.IssuerPublicKey) {
		return errors.New("Invalid IssuerPublicKey")
	}

	if !common.IsHexAddress(ac.DistributionPublicKey) {
		return errors.New("Invalid DistributionPublicKey")
	}

	ac.accountStatus = make(map[string]Status)

	go ac.logStats()
	return nil
}

func (ac *AccountConfigurator) Stop() error {
	ac.Enabled = false
	return nil
}

func (ac *AccountConfigurator) logStats() {
	for {
		if ac.Enabled == false {
			// stop logging
			break
		}
		logger.Infof("statuses: %v", ac.accountStatus)
		time.Sleep(15 * time.Second)
	}
}

// ConfigureAccount configures a new account that participated in ICO.
// * First it creates a new account.
// * Once a signer is replaced on the account, it creates trust lines and exchanges assets.
// from coinmarket place trust lines we get exchange rate and call OnExchange event to collect the exchange rate
func (ac *AccountConfigurator) ConfigureAccount(depositTransaction *types.DepositTransaction) {

	logger.Info("Configuring Tomochain account")

	ac.setAccountStatus(depositTransaction.AssociatedAddress, StatusCreatingAccount)
	defer func() {
		ac.removeAccountStatus(depositTransaction.AssociatedAddress)
	}()

	// Check if account exists. If it is, skip creating it.
	for {
		if ac.Enabled == false {
			break
		}

		_, exists, err := ac.getAccount(depositTransaction.Chain, depositTransaction.AssociatedAddress)
		if err != nil {
			logger.Errorf("Error loading account from Tomochain: %s", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}

		if exists {
			break
		}

		logger.Info("Creating Tomochain account")
		err = ac.createAccountTransaction(depositTransaction.Chain, depositTransaction.AssociatedAddress)
		if err != nil {
			logger.Error("Error creating Tomochain account")
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	if ac.OnAccountCreated != nil {
		ac.OnAccountCreated(depositTransaction.Chain, depositTransaction.AssociatedAddress)
	}

	ac.setAccountStatus(depositTransaction.AssociatedAddress, StatusWaitingForSigner)

	// Wait for signer changes...
	// check if association is correct
	for {
		if ac.Enabled == false {
			break
		}
		account, err := ac.LoadAccount(depositTransaction.Chain, depositTransaction.AssociatedAddress)
		if err != nil {
			logger.Error("Error loading account to check trustline")
			time.Sleep(2 * time.Second)
			continue
		}

		if ac.signerExistsOnly(account) {
			break
		}

		time.Sleep(2 * time.Second)
	}

	logger.Info("Signer found")

	ac.setAccountStatus(depositTransaction.AssociatedAddress, StatusConfiguringAccount)

	// When signer was created we can configure account in Bifrost without requiring
	// the user to share the account's secret key.
	logger.Info("Sending token")
	err := ac.configureAccountTransaction(depositTransaction)
	if err != nil {
		logger.Error("Error configuring an account")
		return
	}

	ac.setAccountStatus(depositTransaction.AssociatedAddress, StatusRemovingSigner)

	if ac.LockUnixTimestamp == 0 {
		logger.Info("Removing temporary signer")
		err = ac.removeTemporarySigner(depositTransaction.Chain, depositTransaction.AssociatedAddress)
		if err != nil {
			logger.Error("Error removing temporary signer")
			return
		}

		if ac.OnExchanged != nil {
			ac.OnExchanged(depositTransaction.Chain, depositTransaction.AssociatedAddress)
		}
	} else {
		logger.Info("Creating unlock transaction to remove temporary signer")
		transaction, err := ac.buildUnlockAccountTransaction(depositTransaction.AssociatedAddress)
		if err != nil {
			logger.Error("Error creating unlock transaction")
			return
		}

		if ac.OnExchangedTimelocked != nil {
			ac.OnExchangedTimelocked(depositTransaction.Chain, depositTransaction.AssociatedAddress, transaction)
		}
	}

	logger.Info("Account successully configured")
}

func (ac *AccountConfigurator) LoadAccount(chain types.Chain, publicKey string) (*types.AddressAssociation, error) {
	if ac.LoadAccountHandler != nil {
		return ac.LoadAccountHandler(chain, publicKey)
	}

	return nil, nil
}

func (ac *AccountConfigurator) setAccountStatus(account string, status Status) {
	ac.accountStatusMutex.Lock()
	defer ac.accountStatusMutex.Unlock()
	ac.accountStatus[account] = status
}

func (ac *AccountConfigurator) removeAccountStatus(account string) {
	ac.accountStatusMutex.Lock()
	defer ac.accountStatusMutex.Unlock()
	delete(ac.accountStatus, account)
}

func (ac *AccountConfigurator) getAccount(chain types.Chain, account string) (*types.AddressAssociation, bool, error) {
	hAccount, err := ac.LoadAccount(chain, account)
	return hAccount, true, err
}

// signerExistsOnly returns true if account has exactly one signer and it's
// equal to `signerPublicKey`.
func (ac *AccountConfigurator) signerExistsOnly(account *types.AddressAssociation) bool {
	tempSignerFound := false

	logger.Debugf("account :%v, signerPublicKey: %s", account, ac.signerPublicKey.Hex())

	if account != nil && account.TomochainPublicKey == ac.signerPublicKey {
		tempSignerFound = true
	}

	return tempSignerFound
}
