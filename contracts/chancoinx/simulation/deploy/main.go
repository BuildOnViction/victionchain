package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/chancoin-core/chancoin-gold/accounts/abi/bind"
	"github.com/chancoin-core/chancoin-gold/common"
	"github.com/chancoin-core/chancoin-gold/contracts/chancoinx"
	"github.com/chancoin-core/chancoin-gold/contracts/chancoinx/simulation"
	"github.com/chancoin-core/chancoin-gold/ethclient"
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
	fmt.Println(nonce)
	auth := bind.NewKeyedTransactor(simulation.MainKey)
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4000000) // in units
	auth.GasPrice = big.NewInt(210000000000000)

	// init trc21 issuer
	auth.Nonce = big.NewInt(int64(nonce))
	trc21IssuerAddr, trc21Issuer, err := chancoinx.DeployTRC21Issuer(auth, client, simulation.MinTRC21Apply)
	if err != nil {
		log.Fatal("DeployTRC21Issuer", err)
	}
	trc21Issuer.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fmt.Println("===> trc21 issuer address", trc21IssuerAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : TRC Issuer")
	time.Sleep(2 * time.Second)

	//init CHANCOINX Listing in
	auth.Nonce = big.NewInt(int64(nonce + 1))
	chancoinxListtingAddr, chancoinxListing, err := chancoinx.DeployCHANCOINXListing(auth, client)
	if err != nil {
		log.Fatal("DeployCHANCOINXListing", err)
	}
	chancoinxListing.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fmt.Println("===> chancoinx listing address", chancoinxListtingAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : chancoinx listing ")
	time.Sleep(2 * time.Second)

	// init Relayer Registration
	auth.Nonce = big.NewInt(int64(nonce + 2))
	relayerRegistrationAddr, relayerRegistration, err := chancoinx.DeployRelayerRegistration(auth, client, chancoinxListtingAddr, simulation.MaxRelayers, simulation.MaxTokenList, simulation.MinDeposit)
	if err != nil {
		log.Fatal("DeployRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fmt.Println("===> relayer registration address", relayerRegistrationAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : relayer registration ")
	time.Sleep(2 * time.Second)

	currentNonce := nonce + 3
	tokenList := initTRC21(auth, client, currentNonce, simulation.TokenNameList)

	currentNonce = currentNonce + uint64(len(simulation.TokenNameList)) // init smartcontract

	applyIssuer(trc21Issuer, tokenList, currentNonce)

	currentNonce = currentNonce + uint64(len(simulation.TokenNameList))
	applyChancoinXListing(chancoinxListing, tokenList, currentNonce)

	currentNonce = currentNonce + uint64(len(simulation.TokenNameList))
	airdrop(auth, client, tokenList, simulation.TeamAddresses, currentNonce)

	// relayer registration
	ownerRelayer := bind.NewKeyedTransactor(simulation.OwnerRelayerKey)
	nonce, _ = client.NonceAt(context.Background(), simulation.OwnerRelayerAddr, nil)
	relayerRegistration, err = chancoinx.NewRelayerRegistration(ownerRelayer, relayerRegistrationAddr, client)
	if err != nil {
		log.Fatal("NewRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	relayerRegistration.TransactOpts.Value = simulation.MinDeposit
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fromTokens := []common.Address{}
	toTokens := []common.Address{}

	/*
		for _, token := range tokenList {
			fromTokens = append(fromTokens, token["address"].(common.Address))
			toTokens = append(toTokens, simulation.CHANCOINNative)
		}
	*/

	// CHANCOIN/BTC
	fromTokens = append(fromTokens, simulation.CHANCOINNative)
	toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	// CHANCOIN/USDT
	fromTokens = append(fromTokens, simulation.CHANCOINNative)
	toTokens = append(toTokens, tokenList[9]["address"].(common.Address))

	// ETH/CHANCOIN
	fromTokens = append(fromTokens, tokenList[1]["address"].(common.Address))
	toTokens = append(toTokens, simulation.CHANCOINNative)

	fromTokens = append(fromTokens, tokenList[2]["address"].(common.Address))
	toTokens = append(toTokens, simulation.CHANCOINNative)

	fromTokens = append(fromTokens, tokenList[3]["address"].(common.Address))
	toTokens = append(toTokens, simulation.CHANCOINNative)

	fromTokens = append(fromTokens, tokenList[4]["address"].(common.Address))
	toTokens = append(toTokens, simulation.CHANCOINNative)

	fromTokens = append(fromTokens, tokenList[5]["address"].(common.Address))
	toTokens = append(toTokens, simulation.CHANCOINNative)

	fromTokens = append(fromTokens, tokenList[6]["address"].(common.Address))
	toTokens = append(toTokens, simulation.CHANCOINNative)

	fromTokens = append(fromTokens, tokenList[7]["address"].(common.Address))
	toTokens = append(toTokens, simulation.CHANCOINNative)

	fromTokens = append(fromTokens, tokenList[8]["address"].(common.Address))
	toTokens = append(toTokens, simulation.CHANCOINNative)

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
	fmt.Println("wait 10s to apply token to list chancoinx")
	time.Sleep(2 * time.Second)
}

func initTRC21(auth *bind.TransactOpts, client *ethclient.Client, nonce uint64, tokenNameList []string) []map[string]interface{} {
	tokenListResult := []map[string]interface{}{}
	for _, tokenName := range tokenNameList {
		auth.Nonce = big.NewInt(int64(nonce))
		d := uint8(18)
		if tokenName == "USDT" {
			d = 8
		}
		tokenAddr, _, err := chancoinx.DeployTRC21(auth, client, tokenName, tokenName, d, simulation.TRC21TokenCap, simulation.TRC21TokenFee)
		if err != nil {
			log.Fatal("DeployTRC21 ", tokenName, err)
		}

		fmt.Println(tokenName+" token address", tokenAddr.Hex(), "cap", simulation.TRC21TokenCap)
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

func applyIssuer(trc21Issuer *chancoinx.TRC21Issuer, tokenList []map[string]interface{}, nonce uint64) {
	for _, token := range tokenList {
		trc21Issuer.TransactOpts.Nonce = big.NewInt(int64(nonce))
		trc21Issuer.TransactOpts.Value = simulation.MinTRC21Apply
		_, err := trc21Issuer.Apply(token["address"].(common.Address))
		if err != nil {
			log.Fatal("trc21Issuer Apply  ", token["name"].(string), err)
		}
		fmt.Println("wait 10s to applyIssuer ", token["name"].(string))
		nonce = nonce + 1
	}
	time.Sleep(5 * time.Second)
}

func applyChancoinXListing(chancoinxListing *chancoinx.CHANCOINXListing, tokenList []map[string]interface{}, nonce uint64) {
	for _, token := range tokenList {
		chancoinxListing.TransactOpts.Nonce = big.NewInt(int64(nonce))
		chancoinxListing.TransactOpts.Value = simulation.ChancoinXListingFee
		_, err := chancoinxListing.Apply(token["address"].(common.Address))
		if err != nil {
			log.Fatal("chancoinxListing Apply ", token["name"].(string), err)
		}
		fmt.Println("wait 10s to applyChancoinXListing ", token["name"].(string))
		nonce = nonce + 1
	}
	time.Sleep(5 * time.Second)
}

func airdrop(auth *bind.TransactOpts, client *ethclient.Client, tokenList []map[string]interface{}, addresses []common.Address, nonce uint64) {
	for _, token := range tokenList {
		for _, address := range addresses {
			trc21Contract, _ := chancoinx.NewTRC21(auth, token["address"].(common.Address), client)
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
