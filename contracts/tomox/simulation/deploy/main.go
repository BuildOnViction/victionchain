package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/contracts/tomox"
	"github.com/tomochain/tomochain/contracts/tomox/simulation"
	"github.com/tomochain/tomochain/ethclient"
)

func main() {
	fmt.Println("========================")
	fmt.Println("mainAddr", simulation.MainAddr.Hex())
	fmt.Println("relayerAddr", simulation.RelayerCoinbaseAddr.Hex())
	fmt.Println("ownerRelayerAddr", simulation.OwnerRelayerAddr.Hex())
	fmt.Println("========================")
	client, err := ethclient.Dial(simulation.RpcEndpoint)
	if err != nil {
		fmt.Println(err, client)
	}
	nonce, _ := client.NonceAt(context.Background(), simulation.MainAddr, nil)
	auth := bind.NewKeyedTransactor(simulation.MainKey)
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(10000000) // in units
	auth.GasPrice = big.NewInt(250000000000000)

	// init trc21 issuer
	auth.Nonce = big.NewInt(int64(nonce))
	trc21IssuerAddr, trc21Issuer, err := tomox.DeployTRC21Issuer(auth, client, simulation.MinTRC21Apply)
	if err != nil {
		log.Fatal("DeployTRC21Issuer", err)
	}
	trc21Issuer.TransactOpts.GasPrice = big.NewInt(250000000000000)

	fmt.Println("===> trc21 issuer address", trc21IssuerAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : TRC Issuer")
	time.Sleep(2 * time.Second)

	//init TOMOX Listing in
	auth.Nonce = big.NewInt(int64(nonce + 1))
	tomoxListtingAddr, tomoxListing, err := tomox.DeployTOMOXListing(auth, client)
	if err != nil {
		log.Fatal("DeployTOMOXListing", err)
	}
	tomoxListing.TransactOpts.GasPrice = big.NewInt(250000000000000)

	fmt.Println("===> tomox listing address", tomoxListtingAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : tomox listing ")
	time.Sleep(2 * time.Second)

	// init Relayer Registration
	auth.Nonce = big.NewInt(int64(nonce + 2))
	relayerRegistrationAddr, relayerRegistration, err := tomox.DeployRelayerRegistration(auth, client, tomoxListtingAddr, simulation.MaxRelayers, simulation.MaxTokenList, simulation.MinDeposit)
	if err != nil {
		log.Fatal("DeployRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(250000000000000)

	fmt.Println("===> relayer registration address", relayerRegistrationAddr.Hex())
	fmt.Println("wait 2s to execute init smart contract : relayer registration ")
	time.Sleep(2 * time.Second)

	auth.Nonce = big.NewInt(int64(nonce + 3))
	lendingRelayerRegistrationAddr, lendingRelayerRegistration, err := tomox.DeployLendingRelayerRegistration(auth, client, relayerRegistrationAddr, tomoxListtingAddr)
	if err != nil {
		log.Fatal("DeployLendingRelayerRegistration", err)
	}
	lendingRelayerRegistration.TransactOpts.GasPrice = big.NewInt(250000000000000)

	fmt.Println("===> lending relayer registration address", lendingRelayerRegistrationAddr.Hex())
	fmt.Println("wait 2s to execute init smart contract : lending relayer registration ")
	time.Sleep(2 * time.Second)

	currentNonce := nonce + 4
	tokenList := initTRC21(auth, client, currentNonce, simulation.TokenNameList)

	currentNonce = currentNonce + uint64(len(simulation.TokenNameList)) // init smartcontract

	applyIssuer(trc21Issuer, tokenList, currentNonce)

	currentNonce = currentNonce + uint64(len(simulation.TokenNameList))
	applyTomoXListing(tomoxListing, tokenList, currentNonce)

	// BTC Collateral
	nonce = currentNonce + uint64(len(simulation.TokenNameList))
	lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	_, err = lendingRelayerRegistration.AddCollateral(tokenList[0]["address"].(common.Address), simulation.CollateralDepositRate, simulation.CollateralLiquidationRate, simulation.CollateralRecallRate)
	if err != nil {
		log.Fatal("Lending add collateral", err)
	}

	// ETH Collateral
	nonce = nonce + 1
	lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	_, err = lendingRelayerRegistration.AddCollateral(tokenList[1]["address"].(common.Address), simulation.CollateralDepositRate, simulation.CollateralLiquidationRate, simulation.CollateralRecallRate)
	if err != nil {
		log.Fatal("Lending add collateral", err)
	}

	// TOMO Collateral
	nonce = nonce + 1
	lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	_, err = lendingRelayerRegistration.AddCollateral(simulation.TOMONative, simulation.CollateralDepositRate, simulation.CollateralLiquidationRate, simulation.CollateralRecallRate)

	if err != nil {
		log.Fatal("Lending add collateral", err)
	}

	// XRP ILO Collateral
	nonce = nonce + 1
	lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	_, err = lendingRelayerRegistration.AddILOCollateral(tokenList[2]["address"].(common.Address), simulation.CollateralDepositRate, simulation.CollateralLiquidationRate, simulation.CollateralRecallRate)
	if err != nil {
		log.Fatal("Lending add ILO collateral", err)
	}

	// USD lending base
	nonce = nonce + 1
	lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	_, err = lendingRelayerRegistration.AddBaseToken(tokenList[9]["address"].(common.Address))
	if err != nil {
		log.Fatal("Lending add base token USD", err)
	}

	// TOMO lending base
	nonce = nonce + 1
	lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	_, err = lendingRelayerRegistration.AddBaseToken(simulation.TOMONative)
	if err != nil {
		log.Fatal("Lending add base token TOMO", err)
	}

	// BTC lending base
	nonce = nonce + 1
	lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	_, err = lendingRelayerRegistration.AddBaseToken(tokenList[0]["address"].(common.Address))
	if err != nil {
		log.Fatal("Lending add base token BTC", err)
	}

	// add term 1 minute for testing
	for i := 0; i < len(simulation.Terms); i++ {
		nonce = nonce + 1
		lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
		_, err = lendingRelayerRegistration.AddTerm(simulation.Terms[i])
		if err != nil {
			log.Fatal("Lending add terms", err)
		}
		fmt.Println("wait 2s to add term lending lending contract", simulation.Terms[i])
		time.Sleep(2 * time.Second)
	}
	fmt.Println("wait 2s to setup lending contract")
	time.Sleep(2 * time.Second)

	currentNonce = nonce + 1
	airdrop(auth, client, tokenList, simulation.TeamAddresses, currentNonce)

	// relayer registration
	ownerRelayer := bind.NewKeyedTransactor(simulation.OwnerRelayerKey)
	nonce, _ = client.NonceAt(context.Background(), simulation.OwnerRelayerAddr, nil)
	relayerRegistration, err = tomox.NewRelayerRegistration(ownerRelayer, relayerRegistrationAddr, client)
	if err != nil {
		log.Fatal("NewRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	relayerRegistration.TransactOpts.Value = simulation.MinDeposit
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(250000000000000)

	fromTokens := []common.Address{}
	toTokens := []common.Address{}

	/*
		for _, token := range tokenList {
			fromTokens = append(fromTokens, token["address"].(common.Address))
			toTokens = append(toTokens, simulation.TOMONative)
		}
	*/

	// TOMO/BTC
	fromTokens = append(fromTokens, simulation.TOMONative)
	toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	// TOMO/USDT
	fromTokens = append(fromTokens, simulation.TOMONative)
	toTokens = append(toTokens, tokenList[9]["address"].(common.Address))

	// ETH/TOMO
	fromTokens = append(fromTokens, tokenList[1]["address"].(common.Address))
	toTokens = append(toTokens, simulation.TOMONative)

	fromTokens = append(fromTokens, tokenList[2]["address"].(common.Address))
	toTokens = append(toTokens, simulation.TOMONative)

	fromTokens = append(fromTokens, tokenList[3]["address"].(common.Address))
	toTokens = append(toTokens, simulation.TOMONative)

	fromTokens = append(fromTokens, tokenList[4]["address"].(common.Address))
	toTokens = append(toTokens, simulation.TOMONative)

	fromTokens = append(fromTokens, tokenList[5]["address"].(common.Address))
	toTokens = append(toTokens, simulation.TOMONative)

	fromTokens = append(fromTokens, tokenList[6]["address"].(common.Address))
	toTokens = append(toTokens, simulation.TOMONative)

	fromTokens = append(fromTokens, tokenList[7]["address"].(common.Address))
	toTokens = append(toTokens, simulation.TOMONative)

	fromTokens = append(fromTokens, tokenList[8]["address"].(common.Address))
	toTokens = append(toTokens, simulation.TOMONative)

	// ETH/BTC
	fromTokens = append(fromTokens, tokenList[1]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	// XRP/BTC
	fromTokens = append(fromTokens, tokenList[2]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	// BTC/USDT
	fromTokens = append(fromTokens, tokenList[0]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[9]["address"].(common.Address))

	// ETH/USDT
	fromTokens = append(fromTokens, tokenList[1]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[9]["address"].(common.Address))

	_, err = relayerRegistration.Register(simulation.RelayerCoinbaseAddr, simulation.TradeFee, fromTokens, toTokens)
	if err != nil {
		log.Fatal("relayerRegistration Register ", err)
	}
	fmt.Println("wait 2s to apply token to list tomox")
	time.Sleep(2 * time.Second)

	// Lending apply
	nonce = nonce + 1
	lendingRelayerRegistration, err = tomox.NewLendingRelayerRegistration(ownerRelayer, lendingRelayerRegistrationAddr, client)
	if err != nil {
		log.Fatal("NewRelayerRegistration", err)
	}
	lendingRelayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	lendingRelayerRegistration.TransactOpts.Value = big.NewInt(0)
	lendingRelayerRegistration.TransactOpts.GasPrice = big.NewInt(250000000000000)
	lendingRelayerRegistration.TransactOpts.GasLimit = uint64(10000000)

	baseTokens := []common.Address{}
	terms := []*big.Int{}
	collaterals := []common.Address{}

	// USD 1 minute for testing
	baseTokens = append(baseTokens, tokenList[9]["address"].(common.Address))
	terms = append(terms, big.NewInt(60))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	baseTokens = append(baseTokens, tokenList[9]["address"].(common.Address))
	terms = append(terms, big.NewInt(60))
	collaterals = append(collaterals, tokenList[2]["address"].(common.Address))

	// USD 1 days
	baseTokens = append(baseTokens, tokenList[9]["address"].(common.Address))
	terms = append(terms, big.NewInt(86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	baseTokens = append(baseTokens, tokenList[9]["address"].(common.Address))
	terms = append(terms, big.NewInt(86400))
	collaterals = append(collaterals, tokenList[2]["address"].(common.Address))

	// USD 7 days
	baseTokens = append(baseTokens, tokenList[9]["address"].(common.Address))
	terms = append(terms, big.NewInt(7*86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// USD 30 days
	baseTokens = append(baseTokens, tokenList[9]["address"].(common.Address))
	terms = append(terms, big.NewInt(30*86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// TOMO 1 min
	baseTokens = append(baseTokens, simulation.TOMONative)
	terms = append(terms, big.NewInt(60))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// TOMO 1 day
	baseTokens = append(baseTokens, simulation.TOMONative)
	terms = append(terms, big.NewInt(86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// TOMO 7 days
	baseTokens = append(baseTokens, simulation.TOMONative)
	terms = append(terms, big.NewInt(7*86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// TOMO 30 days
	baseTokens = append(baseTokens, simulation.TOMONative)
	terms = append(terms, big.NewInt(30*86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// BTC 1 min
	baseTokens = append(baseTokens, tokenList[0]["address"].(common.Address))
	terms = append(terms, big.NewInt(60))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// BTC 1 day
	baseTokens = append(baseTokens, tokenList[0]["address"].(common.Address))
	terms = append(terms, big.NewInt(86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// BTC 7 days
	baseTokens = append(baseTokens, tokenList[0]["address"].(common.Address))
	terms = append(terms, big.NewInt(7*86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	// BTC 30 days
	baseTokens = append(baseTokens, tokenList[0]["address"].(common.Address))
	terms = append(terms, big.NewInt(30*86400))
	collaterals = append(collaterals, common.HexToAddress("0x0"))

	_, err = lendingRelayerRegistration.Update(simulation.RelayerCoinbaseAddr, simulation.LendingTradeFee, baseTokens, terms, collaterals)
	if err != nil {
		log.Fatal("lendingRelayerRegistration Update", err)
	}

	fmt.Println("wait 2s to update lending contract")
	time.Sleep(2 * time.Second)
}

func initTRC21(auth *bind.TransactOpts, client *ethclient.Client, nonce uint64, tokenNameList []string) []map[string]interface{} {
	tokenListResult := []map[string]interface{}{}
	for _, tokenName := range tokenNameList {
		auth.Nonce = big.NewInt(int64(nonce))
		d := uint8(18)
		depositFee := big.NewInt(0)
		withdrawFee := big.NewInt(0)
		tokenCap := simulation.TRC21TokenCap
		if tokenName == "USDT" {
			d = 8
			tokenCap.Div(simulation.TRC21TokenCap, big.NewInt(10000000000))
			withdrawFee = big.NewInt(97000000)
		}
		if tokenName == "BTC" {
			withdrawFee = big.NewInt(400000000000000)
		}
		if tokenName == "ETH" {
			withdrawFee = big.NewInt(3000000000000000)
		}
		tokenAddr, _, err := tomox.DeployTRC21(auth, client, simulation.Owners, simulation.Required, tokenName, tokenName, d, tokenCap, simulation.TRC21TokenFee, depositFee, withdrawFee)
		if err != nil {
			log.Fatal("DeployTRC21 ", tokenName, err)
		}

		fmt.Println(tokenName+" token address", tokenAddr.Hex(), "cap", simulation.TRC21TokenCap)

		tokenListResult = append(tokenListResult, map[string]interface{}{
			"name":    tokenName,
			"address": tokenAddr,
		})
		nonce = nonce + 1
	}
	time.Sleep(5 * time.Second)
	return tokenListResult
}

func applyIssuer(trc21Issuer *tomox.TRC21Issuer, tokenList []map[string]interface{}, nonce uint64) {
	for _, token := range tokenList {
		trc21Issuer.TransactOpts.Nonce = big.NewInt(int64(nonce))
		trc21Issuer.TransactOpts.Value = simulation.MinTRC21Apply
		_, err := trc21Issuer.Apply(token["address"].(common.Address))
		if err != nil {
			log.Fatal("trc21Issuer Apply  ", token["name"].(string), err)
		}
		fmt.Println("applyIssuer ", token["name"].(string))
		nonce = nonce + 1
	}
	time.Sleep(5 * time.Second)
}

func applyTomoXListing(tomoxListing *tomox.TOMOXListing, tokenList []map[string]interface{}, nonce uint64) {
	for _, token := range tokenList {
		tomoxListing.TransactOpts.Nonce = big.NewInt(int64(nonce))
		tomoxListing.TransactOpts.Value = simulation.TomoXListingFee
		_, err := tomoxListing.Apply(token["address"].(common.Address))
		if err != nil {
			log.Fatal("tomoxListing Apply ", token["name"].(string), err)
		}
		fmt.Println("applyTomoXListing ", token["name"].(string))
		nonce = nonce + 1
	}
	time.Sleep(5 * time.Second)
}

func airdrop(auth *bind.TransactOpts, client *ethclient.Client, tokenList []map[string]interface{}, addresses []common.Address, nonce uint64) {
	for _, token := range tokenList {
		for _, address := range addresses {
			trc21Contract, _ := tomox.NewTRC21(auth, token["address"].(common.Address), client)
			trc21Contract.TransactOpts.Nonce = big.NewInt(int64(nonce))
			_, err := trc21Contract.Transfer(address, big.NewInt(0).Mul(common.BasePrice, big.NewInt(1000000)))
			if err == nil {
				fmt.Printf("Transfer %v to %v successfully", token["name"].(string), address.String())
				fmt.Println()
			} else {
				fmt.Printf("Transfer %v to %v failed!", token["name"].(string), address.String())
				fmt.Println()
			}
			nonce = nonce + 1
		}
	}
	time.Sleep(5 * time.Second)
}
