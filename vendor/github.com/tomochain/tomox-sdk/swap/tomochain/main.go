package tomochain

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/swap/config"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
)

// Status describes status of account processing
type Status string

var logger = utils.Logger

const (
	StatusCreatingAccount    Status = "creating_account"
	StatusWaitingForSigner   Status = "waiting_for_signer"
	StatusConfiguringAccount Status = "configuring_account"
	StatusRemovingSigner     Status = "removing_signer"
)

type LoadAccountHandler func(chain types.Chain, publicKey string) (*types.AddressAssociation, error)

// AccountConfigurator is responsible for configuring new Tomochain accounts that
// participate in ICO.
// Infact, AccountConfigurator will be replaced by smart contract
type AccountConfigurator struct {
	Enabled               bool
	IssuerPublicKey       string
	DistributionPublicKey string

	LockUnixTimestamp uint64
	TokenAssetCode    string
	TokenPriceBTC     string
	TokenPriceETH     string
	StartingBalance   string

	LoadAccountHandler    LoadAccountHandler
	OnSubmitTransaction   func(chain types.Chain, destination string, transaction *types.AssociationTransaction) error
	OnAccountCreated      func(chain types.Chain, destination string)
	OnExchanged           func(chain types.Chain, destination string)
	OnExchangedTimelocked func(chain types.Chain, destination string, transaction *types.AssociationTransaction)

	signerPublicKey  common.Address
	signerPrivateKey *ecdsa.PrivateKey

	accountStatus      map[string]Status
	accountStatusMutex sync.Mutex
}

func NewAccountConfigurator(c *config.TomochainConfig) *AccountConfigurator {
	return &AccountConfigurator{
		IssuerPublicKey:       c.IssuerPublicKey,
		DistributionPublicKey: c.DistributionPublicKey,
		signerPublicKey:       c.GetPublicKey(),
		signerPrivateKey:      c.GetPrivateKey(),
		TokenAssetCode:        c.TokenAssetCode,
		StartingBalance:       c.StartingBalance,
		LockUnixTimestamp:     c.LockUnixTimestamp,
	}
}
