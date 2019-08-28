package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/tomox"
	"github.com/ethereum/go-ethereum/contracts/tomox/simulation"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
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
	trc21IssuerAddr, trc21Issuer, err := tomox.DeployTRC21Issuer(auth, client, simulation.MinTRC21Apply)
	if err != nil {
		log.Fatal("DeployTRC21Issuer", err)
	}
	trc21Issuer.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fmt.Println("===> trc21 issuer address", trc21IssuerAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : TRC Issuer")
	time.Sleep(10 * time.Second)

	//init TOMOX Listing in
	auth.Nonce = big.NewInt(int64(nonce + 1))
	tomoxListtingAddr, tomoxListing, err := tomox.DeployTOMOXListing(auth, client)
	if err != nil {
		log.Fatal("DeployTOMOXListing", err)
	}
	tomoxListing.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fmt.Println("===> tomox listing address", tomoxListtingAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : tomox listing ")
	time.Sleep(10 * time.Second)

	// init Relayer Registration
	auth.Nonce = big.NewInt(int64(nonce + 2))
	relayerRegistrationAddr, relayerRegistration, err := tomox.DeployRelayerRegistration(auth, client, simulation.MaxRelayers, simulation.MaxTokenList, simulation.MinDeposit)
	if err != nil {
		log.Fatal("DeployRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fmt.Println("===> relayer registration address", relayerRegistrationAddr.Hex())
	fmt.Println("wait 10s to execute init smart contract : relayer registration ")
	time.Sleep(10 * time.Second)

	currentNonce := nonce + 3
	tokenList := initTRC21(auth, client, currentNonce, simulation.TokenNameList)

	currentNonce = currentNonce + uint64(len(simulation.TokenNameList)) // init smartcontract

	applyIssuer(trc21Issuer, tokenList, currentNonce)

	currentNonce = currentNonce + uint64(len(simulation.TokenNameList))
	applyTomoXListing(tomoxListing, tokenList, currentNonce)

	currentNonce = currentNonce + uint64(len(simulation.TokenNameList))
	airdrop(auth, client, tokenList, simulation.TeamAddresses, currentNonce)

	// relayer registration
	ownerRelayer := bind.NewKeyedTransactor(simulation.OwnerRelayerKey)
	nonce, _ = client.NonceAt(context.Background(), simulation.OwnerRelayerAddr, nil)
	relayerRegistration, err = tomox.NewRelayerRegistration(ownerRelayer, relayerRegistrationAddr, client)
	if err != nil {
		log.Fatal("NewRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	relayerRegistration.TransactOpts.Value = big.NewInt(0).Mul(simulation.MinDeposit, simulation.BaseTOMO)
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fromTokens := []common.Address{}
	toTokens := []common.Address{}
	for _, token := range tokenList {
		fromTokens = append(fromTokens, token["address"].(common.Address))
		toTokens = append(toTokens, simulation.TOMONative)
	}

	// ETH/BTC
	fromTokens = append(fromTokens, tokenList[1]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	// XRP/BTC
	fromTokens = append(fromTokens, tokenList[2]["address"].(common.Address))
	toTokens = append(toTokens, tokenList[0]["address"].(common.Address))

	_, err = relayerRegistration.Register(simulation.RelayerCoinbaseAddr, simulation.TradeFee, fromTokens, toTokens)
	if err != nil {
		log.Fatal("relayerRegistration Register", err)
	}
	fmt.Println("wait 10s to apply token to list tomox")
	time.Sleep(10 * time.Second)
}

func initTRC21(auth *bind.TransactOpts, client *ethclient.Client, nonce uint64, tokenNameList []string) []map[string]interface{} {
	tokenListResult := []map[string]interface{}{}
	for _, tokenName := range tokenNameList {
		auth.Nonce = big.NewInt(int64(nonce))
		tokenAddr, _, err := tomox.DeployTRC21(auth, client, tokenName, tokenName, 18, simulation.TRC21TokenCap, simulation.TRC21TokenFee)
		if err != nil {
			log.Fatal("DeployTRC21 ", tokenName, err)
		}

		fmt.Println(tokenName+" token address", tokenAddr.Hex(), "cap", simulation.TRC21TokenCap)
		fmt.Println("wait 10s to execute init smart contract", tokenName)
		time.Sleep(10 * time.Second)

		tokenListResult = append(tokenListResult, map[string]interface{}{
			"name":    tokenName,
			"address": tokenAddr,
		})
		nonce = nonce + 1
	}
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
		fmt.Println("wait 10s to applyIssuer ", token["name"].(string))
		time.Sleep(5 * time.Second)
		nonce = nonce + 1

	}
}

func applyTomoXListing(tomoxListing *tomox.TOMOXListing, tokenList []map[string]interface{}, nonce uint64) {
	for _, token := range tokenList {
		tomoxListing.TransactOpts.Nonce = big.NewInt(int64(nonce))
		_, err := tomoxListing.Apply(token["address"].(common.Address))
		if err != nil {
			log.Fatal("tomoxListing Apply ", token["name"].(string), err)
		}
		fmt.Println("wait 10s to applyTomoXListing ", token["name"].(string))
		time.Sleep(5 * time.Second)
		nonce = nonce + 1
	}
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
			time.Sleep(5 * time.Second)
			nonce = nonce + 1
		}
	}
}
