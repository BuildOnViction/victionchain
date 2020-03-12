package chancoinx

import (
	"github.com/chancoin-core/chancoin-gold/accounts/abi/bind"
	"github.com/chancoin-core/chancoin-gold/common"
	"github.com/chancoin-core/chancoin-gold/contracts/chancoinx/contract"
)

type CHANCOINXListing struct {
	*contract.CHANCOINXListingSession
	contractBackend bind.ContractBackend
}

func NewMyCHANCOINXListing(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*CHANCOINXListing, error) {
	smartContract, err := contract.NewCHANCOINXListing(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &CHANCOINXListing{
		&contract.CHANCOINXListingSession{
			Contract:     smartContract,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployCHANCOINXListing(transactOpts *bind.TransactOpts, contractBackend bind.ContractBackend) (common.Address, *CHANCOINXListing, error) {
	contractAddr, _, _, err := contract.DeployCHANCOINXListing(transactOpts, contractBackend)
	if err != nil {
		return contractAddr, nil, err
	}
	smartContract, err := NewMyCHANCOINXListing(transactOpts, contractAddr, contractBackend)
	if err != nil {
		return contractAddr, nil, err
	}

	return contractAddr, smartContract, nil
}
