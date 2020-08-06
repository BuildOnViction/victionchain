package main

import (
	"context"
	"fmt"
	"github.com/tomochain/tomochain/contracts/tomox/simulation"
	"log"
	"math/big"
	"time"

	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/contracts/tomox"
	"github.com/tomochain/tomochain/contracts/tomox/testnet"
	"github.com/tomochain/tomochain/ethclient"
)

func main() {
	fmt.Println("========================")
	fmt.Println("mainAddr", testnet.MainAddr.Hex())
	fmt.Println("relayerAddr", testnet.RelayerCoinbaseAddr.Hex())
	fmt.Println("ownerRelayerAddr", testnet.OwnerRelayerAddr.Hex())
	fmt.Println("========================")
	client, err := ethclient.Dial(testnet.RpcEndpoint)
	if err != nil {
		fmt.Println(err, client)
	}
	nonce, _ := client.NonceAt(context.Background(), testnet.MainAddr, nil)
	fmt.Println(nonce)
	auth := bind.NewKeyedTransactor(testnet.MainKey)
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4000000) // in units
	auth.GasPrice = big.NewInt(250000000000000)

	// init trc21 issuer
	auth.Nonce = big.NewInt(int64(nonce))
	trc21IssuerAddr, trc21Issuer, err := tomox.DeployTRC21Issuer(auth, client, testnet.MinTRC21Apply)
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
	relayerRegistrationAddr, relayerRegistration, err := tomox.DeployRelayerRegistration(auth, client, tomoxListtingAddr, testnet.MaxRelayers, testnet.MaxTokenList, testnet.MinDeposit)
	if err != nil {
		log.Fatal("DeployRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(250000000000000)

	fmt.Println("===> relayer registration address", relayerRegistrationAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : relayer registration ")
	time.Sleep(2 * time.Second)

	currentNonce := nonce + 3
	tokenList := initTRC21(auth, client, currentNonce, testnet.TokenNameList)

	currentNonce = currentNonce + uint64(len(testnet.TokenNameList)) // init smartcontract

	applyIssuer(trc21Issuer, tokenList, currentNonce)

	currentNonce = currentNonce + uint64(len(testnet.TokenNameList))
	applyTomoXListing(tomoxListing, tokenList, currentNonce)

	currentNonce = currentNonce + uint64(len(testnet.TokenNameList))
	airdrop(auth, client, tokenList, testnet.TeamAddresses, currentNonce)

	// relayer registration
	ownerRelayer := bind.NewKeyedTransactor(testnet.OwnerRelayerKey)
	nonce, _ = client.NonceAt(context.Background(), testnet.OwnerRelayerAddr, nil)
	relayerRegistration, err = tomox.NewRelayerRegistration(ownerRelayer, relayerRegistrationAddr, client)
	if err != nil {
		log.Fatal("NewRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	relayerRegistration.TransactOpts.Value = testnet.MinDeposit
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fromTokens := []common.Address{}
	toTokens := []common.Address{}

	/*
		for _, token := range tokenList {
			fromTokens = append(fromTokens, token["address"].(common.Address))
			toTokens = append(toTokens, testnet.TOMONative)
		}
	*/

	// TOMO/BTC
	fromTokens = append(fromTokens, testnet.TOMONative)
	toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	// TOMO/USDT
	fromTokens = append(fromTokens, testnet.TOMONative)
	toTokens = append(toTokens, tokenList[9]["address"].(common.Address))

	// ETH/TOMO
	fromTokens = append(fromTokens, tokenList[1]["address"].(common.Address))
	toTokens = append(toTokens, testnet.TOMONative)

	/*
		fromTokens = append(fromTokens, tokenList[2]["address"].(common.Address))
		toTokens = append(toTokens, testnet.TOMONative)

		fromTokens = append(fromTokens, tokenList[3]["address"].(common.Address))
		toTokens = append(toTokens, testnet.TOMONative)

		fromTokens = append(fromTokens, tokenList[4]["address"].(common.Address))
		toTokens = append(toTokens, testnet.TOMONative)

		fromTokens = append(fromTokens, tokenList[5]["address"].(common.Address))
		toTokens = append(toTokens, testnet.TOMONative)

		fromTokens = append(fromTokens, tokenList[6]["address"].(common.Address))
		toTokens = append(toTokens, testnet.TOMONative)

		fromTokens = append(fromTokens, tokenList[7]["address"].(common.Address))
		toTokens = append(toTokens, testnet.TOMONative)

		fromTokens = append(fromTokens, tokenList[8]["address"].(common.Address))
		toTokens = append(toTokens, testnet.TOMONative)
	*/

	// ETH/BTC
	fromTokens = append(fromTokens, tokenList[1]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	// XRP/BTC
	// fromTokens = append(fromTokens, tokenList[2]["address"].(common.Address))
	// toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	// BTC/USDT
	fromTokens = append(fromTokens, tokenList[0]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[9]["address"].(common.Address))

	// ETH/USDT
	fromTokens = append(fromTokens, tokenList[1]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[9]["address"].(common.Address))

	_, err = relayerRegistration.Register(testnet.RelayerCoinbaseAddr, testnet.TradeFee, fromTokens, toTokens)
	if err != nil {
		log.Fatal("relayerRegistration Register ", err)
	}
	fmt.Println("wait 10s to apply token to list tomox")
	time.Sleep(2 * time.Second)
}

func initTRC21(auth *bind.TransactOpts, client *ethclient.Client, nonce uint64, tokenNameList []string) []map[string]interface{} {
	tokenListResult := []map[string]interface{}{}
	depositFee := big.NewInt(0)
	withdrawFee := big.NewInt(0)
	for _, tokenName := range tokenNameList {
		auth.Nonce = big.NewInt(int64(nonce))
		d := uint8(18)
		if tokenName == "USDT" {
			d = 8
		}
		tokenAddr, _, err := tomox.DeployTRC21(auth, client, simulation.Owners, simulation.Required, tokenName, tokenName, d, testnet.TRC21TokenCap, testnet.TRC21TokenFee, depositFee, withdrawFee)
		if err != nil {
			log.Fatal("DeployTRC21 ", tokenName, err)
		}

		fmt.Println(tokenName+" token address", tokenAddr.Hex(), "cap", testnet.TRC21TokenCap)
		fmt.Println("wait 10s to execute init smart contract", tokenName)

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
		trc21Issuer.TransactOpts.Value = testnet.MinTRC21Apply
		_, err := trc21Issuer.Apply(token["address"].(common.Address))
		if err != nil {
			log.Fatal("trc21Issuer Apply  ", token["name"].(string), err)
		}
		fmt.Println("wait 10s to applyIssuer ", token["name"].(string))
		nonce = nonce + 1
	}
	time.Sleep(5 * time.Second)
}

func applyTomoXListing(tomoxListing *tomox.TOMOXListing, tokenList []map[string]interface{}, nonce uint64) {
	for _, token := range tokenList {
		tomoxListing.TransactOpts.Nonce = big.NewInt(int64(nonce))
		tomoxListing.TransactOpts.Value = testnet.TomoXListingFee
		_, err := tomoxListing.Apply(token["address"].(common.Address))
		if err != nil {
			log.Fatal("tomoxListing Apply ", token["name"].(string), err)
		}
		fmt.Println("wait 10s to applyTomoXListing ", token["name"].(string))
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
