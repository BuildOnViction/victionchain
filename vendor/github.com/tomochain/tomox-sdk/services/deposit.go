package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/ethereum"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/swap"
	swapConfig "github.com/tomochain/tomox-sdk/swap/config"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/ws"
)

// need to refractor using interface.SwappEngine and only expose neccessary methods
type DepositService struct {
	configDao      interfaces.ConfigDao
	associationDao interfaces.AssociationDao
	pairDao        interfaces.PairDao
	orderDao       interfaces.OrderDao
	swapEngine     *swap.Engine
	engine         interfaces.Engine
	broker         *rabbitmq.Connection
}

// NewAddressService returns a new instance of accountService
func NewDepositService(
	configDao interfaces.ConfigDao,
	associationDao interfaces.AssociationDao,
	pairDao interfaces.PairDao,
	orderDao interfaces.OrderDao,
	swapEngine *swap.Engine,
	engine interfaces.Engine,
	broker *rabbitmq.Connection,
) *DepositService {

	depositService := &DepositService{configDao, associationDao, pairDao, orderDao, swapEngine, engine, broker}

	// set storage engine to this service
	swapEngine.SetStorage(depositService)

	swapEngine.SetQueue(depositService)

	// run watching
	swapEngine.Start()

	return depositService
}

func (s *DepositService) EthereumClient() interfaces.EthereumClient {
	provider := s.engine.Provider().(*ethereum.EthereumProvider)
	return provider.Client
}

func (s *DepositService) SetDelegate(handler interfaces.SwapEngineHandler) {
	// set event handler delegate to this service
	s.swapEngine.SetDelegate(handler)
}

func (s *DepositService) GenerateAddress(chain types.Chain) (common.Address, uint64, error) {

	err := s.configDao.IncrementAddressIndex(chain)
	if err != nil {
		return swapConfig.EmptyAddress, 0, err
	}
	index, err := s.configDao.GetAddressIndex(chain)
	if err != nil {
		return swapConfig.EmptyAddress, 0, err
	}
	logger.Infof("Current index: %d", index)
	address, err := s.swapEngine.EthereumAddressGenerator().Generate(index)
	return address, index, err
}

func (s *DepositService) SignerPublicKey() common.Address {
	return s.swapEngine.SignerPublicKey()
}

func (s *DepositService) GetSchemaVersion() uint64 {
	return s.configDao.GetSchemaVersion()
}

func (s *DepositService) RecoveryTransaction(chain types.Chain, address common.Address) error {
	return nil
}

/***** implement Storage interface ***/
func (s *DepositService) GetBlockToProcess(chain types.Chain) (uint64, error) {
	return s.configDao.GetBlockToProcess(chain)
}

func (s *DepositService) SaveLastProcessedBlock(chain types.Chain, block uint64) error {
	return s.configDao.SaveLastProcessedBlock(chain, block)
}

func (s *DepositService) SaveDepositTransaction(chain types.Chain, sourceAccount common.Address, txEnvelope string) error {
	return s.associationDao.SaveDepositTransaction(chain, sourceAccount, txEnvelope)
}

func (s *DepositService) QueueAdd(transaction *types.DepositTransaction) error {
	err := s.broker.PublishDepositTransaction(transaction)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// QueuePool receives and removes the head of this queue. Returns nil if no elements found.
func (s *DepositService) QueuePool() (<-chan *types.DepositTransaction, error) {
	return s.broker.QueuePoolDepositTransactions()
}

func (s *DepositService) MinimumValueWei() *big.Int {
	return s.swapEngine.MinimumValueWei()
}

func (s *DepositService) MinimumValueSat() int64 {
	return s.swapEngine.MinimumValueSat()
}

func (s *DepositService) GetAssociationByChainAddress(chain types.Chain, userAddress common.Address) (*types.AddressAssociationRecord, error) {
	return s.associationDao.GetAssociationByChainAddress(chain, userAddress)
}

func (s *DepositService) GetAssociationByChainAssociatedAddress(chain types.Chain, associatedAddress common.Address) (*types.AddressAssociationRecord, error) {
	return s.associationDao.GetAssociationByChainAssociatedAddress(chain, associatedAddress)
}

func (s *DepositService) SaveAssociationByChainAddress(chain types.Chain, address common.Address, index uint64, associatedAddress common.Address, pairAddresses *types.PairAddresses) error {

	association := &types.AddressAssociationRecord{
		ID:                bson.NewObjectId(),
		Chain:             chain,
		Address:           address.Hex(),
		AddressIndex:      index,
		Status:            types.PENDING,
		AssociatedAddress: associatedAddress.Hex(),
		PairName:          pairAddresses.Name,
		BaseTokenAddress:  pairAddresses.BaseToken.Hex(),
		QuoteTokenAddress: pairAddresses.QuoteToken.Hex(),
	}

	return s.associationDao.SaveAssociation(association)
}

func (s *DepositService) SaveAssociationStatusByChainAddress(addressAssociation *types.AddressAssociationRecord, status string) error {

	if addressAssociation == nil {
		return errors.New("AddressAssociationRecord is nil")
	}

	userAddress := common.HexToAddress(addressAssociation.AssociatedAddress)
	address := common.HexToAddress(addressAssociation.Address)

	// send message to channel deposit to noti the status, should limit the txEnvelope < 100
	// if not it would be very slow
	if status == types.SUCCESS {
		ws.SendDepositMessage(types.SUCCESS_EVENT, userAddress, addressAssociation)
	} else if status != types.PENDING {
		// just pending and return the status
		ws.SendDepositMessage(types.PENDING, userAddress, status)
	}

	return s.associationDao.SaveAssociationStatus(addressAssociation.Chain, address, status)
}
