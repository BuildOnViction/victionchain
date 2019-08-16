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
	fmt.Println("relayerAddr", simulation.ReplayCoinbaseAddr.Hex())
	fmt.Println("ownerRelayerAddr", simulation.OwnerRelayAddr.Hex())
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

	// init TRC21 token : BTC
	auth.Nonce = big.NewInt(int64(nonce + 3))
	BTCTokenAddr, _, err := tomox.DeployTRC21(auth, client, "BTC", "BTC", 18, simulation.TRC21TokenCap, simulation.TRC21TokenFee)
	if err != nil {
		log.Fatal("DeployTRC21 BTC", err)
	}

	fmt.Println("===>  BTC token address", BTCTokenAddr.Hex(), "cap", simulation.TRC21TokenCap)
	fmt.Println("wait 10s to execute init smart contract")
	time.Sleep(10 * time.Second)

	// init TRC21 token : ETH
	auth.Nonce = big.NewInt(int64(nonce + 4))
	ETHTokenAddr, _, err := tomox.DeployTRC21(auth, client, "ETH", "ETH", 18, simulation.TRC21TokenCap, simulation.TRC21TokenFee)
	if err != nil {
		log.Fatal("DeployTRC21 ETH", err)
	}

	fmt.Println("===>  ETH token address", ETHTokenAddr.Hex(), "cap", simulation.TRC21TokenCap)
	fmt.Println("wait 10s to execute init smart contract")
	time.Sleep(10 * time.Second)

	trc21Issuer.TransactOpts.Nonce = big.NewInt(int64(nonce + 5))
	trc21Issuer.TransactOpts.Value = simulation.MinTRC21Apply
	// apply BTC token to trc21 issuer
	_, err = trc21Issuer.Apply(BTCTokenAddr)
	if err != nil {
		log.Fatal("trc21Issuer Apply BTC ", err)
	}
	trc21Issuer.TransactOpts.Nonce = big.NewInt(int64(nonce + 6))
	trc21Issuer.TransactOpts.Value = simulation.MinTRC21Apply

	// apply ETH token to trc21 issuer
	_, err = trc21Issuer.Apply(ETHTokenAddr)
	if err != nil {
		log.Fatal("trc21Issuer Apply ETH ", err)
	}

	fmt.Println("wait 10s to add token to list issuer")
	time.Sleep(10 * time.Second)

	// aplly to list TomoX Token
	tomoxListing.TransactOpts.Nonce = big.NewInt(int64(nonce + 7))
	_, err = tomoxListing.Apply(BTCTokenAddr)
	if err != nil {
		log.Fatal("tomoxListing Apply BTC", err)
	}
	tomoxListing.TransactOpts.Nonce = big.NewInt(int64(nonce + 8))
	_, err = tomoxListing.Apply(ETHTokenAddr)
	if err != nil {
		log.Fatal("tomoxListing Apply ETH", err)
	}
	fmt.Println("wait 10s to apply token to list tomox")
	time.Sleep(10 * time.Second)

	// relayer registration
	ownerRelayer := bind.NewKeyedTransactor(simulation.OwnerRelayerKey)
	nonce, _ = client.NonceAt(context.Background(), simulation.OwnerRelayAddr, nil)
	relayerRegistration, err = tomox.NewRelayerRegistration(ownerRelayer, relayerRegistrationAddr, client)
	if err != nil {
		log.Fatal("NewRelayerRegistration", err)
	}
	relayerRegistration.TransactOpts.Nonce = big.NewInt(int64(nonce))
	relayerRegistration.TransactOpts.Value = big.NewInt(0).Mul(simulation.MinDeposit, simulation.BaseTOMO)
	relayerRegistration.TransactOpts.GasPrice = big.NewInt(210000000000000)

	fromTokens := []common.Address{BTCTokenAddr, ETHTokenAddr}
	toTokens := []common.Address{simulation.TOMONative,simulation.TOMONative}
	_, err = relayerRegistration.Register(simulation.ReplayCoinbaseAddr, simulation.TradeFee, fromTokens, toTokens)
	if err != nil {
		log.Fatal("relayerRegistration Register", err)
	}
	fmt.Println("wait 10s to apply token to list tomox")
	time.Sleep(10 * time.Second)
}
