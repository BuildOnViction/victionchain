package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/math"
)

type AccountService struct {
	AccountDao interfaces.AccountDao
	TokenDao   interfaces.TokenDao
}

// NewAddressService returns a new instance of accountService
func NewAccountService(
	accountDao interfaces.AccountDao,
	tokenDao interfaces.TokenDao,
) *AccountService {
	return &AccountService{
		AccountDao: accountDao,
		TokenDao:   tokenDao,
	}
}

func (s *AccountService) Create(a *types.Account) error {
	addr := a.Address

	acc, err := s.AccountDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return err
	}

	if acc != nil {
		return ErrAccountExists
	}

	tokens, err := s.TokenDao.GetAll()
	if err != nil {
		logger.Error(err)
		return err
	}

	a.IsBlocked = false
	a.TokenBalances = make(map[common.Address]*types.TokenBalance)

	ten := big.NewInt(10)

	// currently by default, the tokens balances are set to 0
	for _, token := range tokens {
		decimals := big.NewInt(int64(token.Decimals))
		a.TokenBalances[token.ContractAddress] = &types.TokenBalance{
			Address:          token.ContractAddress,
			Symbol:           token.Symbol,
			Balance:          math.Mul(big.NewInt(types.DefaultTestBalance()), math.Exp(ten, decimals)),
			InOrderBalance:   big.NewInt(types.DefaultTestInOrderBalance()),
			AvailableBalance: big.NewInt(types.DefaultTestAvailableBalance()),
		}
	}

	nativeCurrency := types.GetNativeCurrency()

	a.TokenBalances[nativeCurrency.Address] = &types.TokenBalance{
		Address:          nativeCurrency.Address,
		Symbol:           nativeCurrency.Symbol,
		Balance:          math.Mul(big.NewInt(types.DefaultTestBalance()), math.Exp(ten, big.NewInt(int64(nativeCurrency.Decimals)))),
		InOrderBalance:   big.NewInt(types.DefaultTestInOrderBalance()),
		AvailableBalance: big.NewInt(types.DefaultTestAvailableBalance()),
	}

	if a != nil {
		err = s.AccountDao.Create(a)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (s *AccountService) FindOrCreate(addr common.Address) (*types.Account, error) {
	a, err := s.AccountDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if a != nil {
		return a, nil
	}

	tokens, err := s.TokenDao.GetAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	a = &types.Account{
		Address:       addr,
		IsBlocked:     false,
		TokenBalances: make(map[common.Address]*types.TokenBalance),
	}

	ten := big.NewInt(10)

	// currently by default, the tokens balances are set to 0
	for _, t := range tokens {
		decimals := big.NewInt(int64(t.Decimals))
		a.TokenBalances[t.ContractAddress] = &types.TokenBalance{
			Address:          t.ContractAddress,
			Symbol:           t.Symbol,
			Balance:          math.Mul(big.NewInt(types.DefaultTestBalance()), math.Exp(ten, decimals)),
			InOrderBalance:   big.NewInt(types.DefaultTestInOrderBalance()),
			AvailableBalance: big.NewInt(types.DefaultTestAvailableBalance()),
		}
	}

	nativeCurrency := types.GetNativeCurrency()

	a.TokenBalances[nativeCurrency.Address] = &types.TokenBalance{
		Address:          nativeCurrency.Address,
		Symbol:           nativeCurrency.Symbol,
		Balance:          math.Mul(big.NewInt(types.DefaultTestBalance()), math.Exp(ten, big.NewInt(int64(nativeCurrency.Decimals)))),
		InOrderBalance:   big.NewInt(types.DefaultTestInOrderBalance()),
		AvailableBalance: big.NewInt(types.DefaultTestAvailableBalance()),
	}

	err = s.AccountDao.Create(a)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return a, nil
}

func (s *AccountService) GetByID(id bson.ObjectId) (*types.Account, error) {
	return s.AccountDao.GetByID(id)
}

func (s *AccountService) GetAll() ([]types.Account, error) {
	return s.AccountDao.GetAll()
}

func (s *AccountService) GetByAddress(a common.Address) (*types.Account, error) {
	return s.AccountDao.GetByAddress(a)
}

func (s *AccountService) GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error) {
	return s.AccountDao.GetTokenBalance(owner, token)
}

func (s *AccountService) GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error) {
	return s.AccountDao.GetTokenBalances(owner)
}

func (s *AccountService) Transfer(token common.Address, fromAddress common.Address, toAddress common.Address, amount *big.Int) error {
	return s.AccountDao.Transfer(token, fromAddress, toAddress, amount)
}

func (s *AccountService) GetFavoriteTokens(owner common.Address) (map[common.Address]bool, error) {
	return s.AccountDao.GetFavoriteTokens(owner)
}

func (s *AccountService) AddFavoriteToken(owner, token common.Address) error {
	return s.AccountDao.AddFavoriteToken(owner, token)
}

func (s *AccountService) DeleteFavoriteToken(owner, token common.Address) error {
	return s.AccountDao.DeleteFavoriteToken(owner, token)
}
