package services

import (
	"fmt"
	"math/big"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
)

type ValidatorService struct {
	ethereumProvider interfaces.EthereumProvider
	accountDao       interfaces.AccountDao
	orderDao         interfaces.OrderDao
	pairDao          interfaces.PairDao
}

func NewValidatorService(
	ethereumProvider interfaces.EthereumProvider,
	accountDao interfaces.AccountDao,
	orderDao interfaces.OrderDao,
	pairDao interfaces.PairDao,
) *ValidatorService {

	return &ValidatorService{
		ethereumProvider,
		accountDao,
		orderDao,
		pairDao,
	}
}

func (s *ValidatorService) ValidateAvailableBalance(o *types.Order) error {
	pair, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	totalRequiredAmount := o.TotalRequiredSellAmount(pair)

	balanceRecord, err := s.accountDao.GetTokenBalance(o.UserAddress, o.SellToken())
	if err != nil {
		logger.Error(err)
		return err
	}

	var sellTokenBalance *big.Int
	sellTokenBalance = balanceRecord.Balance

	//sellTokenLockedBalance, err := s.orderDao.GetUserLockedBalance(o.UserAddress, o.SellToken(), pair)
	//if err != nil {
	//	logger.Error(err)
	//	return err
	//}

	//availableSellTokenBalance := math.Sub(sellTokenBalance, sellTokenLockedBalance)

	if sellTokenBalance.Cmp(totalRequiredAmount) == -1 {
		return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
	}

	//if availableSellTokenBalance.Cmp(totalRequiredAmount) == -1 {
	//	return fmt.Errorf("Insufficient % available", o.SellTokenSymbol())
	//}

	balanceRecord.Balance.Set(sellTokenBalance)
	//balanceRecord.LockedBalance.Set(totalRequiredAmount)
	err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.SellToken(), balanceRecord)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *ValidatorService) ValidateBalance(o *types.Order) error {
	//exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	pair, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	totalRequiredAmount := o.TotalRequiredSellAmount(pair)

	balanceRecord, err := s.accountDao.GetTokenBalance(o.UserAddress, o.SellToken())
	if err != nil {
		logger.Error(err)
		return err
	}

	var sellTokenBalance *big.Int
	sellTokenBalance = balanceRecord.Balance

	//Sell Token Balance
	if sellTokenBalance.Cmp(totalRequiredAmount) == -1 {
		return fmt.Errorf("Insufficient %v Balance", o.SellTokenSymbol())
	}

	return nil
}
