package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/contracts/tomox"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/ethclient"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8501/")
	if err != nil {
		fmt.Println(err, client)
	}

	MainKey, _ := crypto.HexToECDSA(os.Getenv("OWNER_KEY"))
	MainAddr := crypto.PubkeyToAddress(MainKey.PublicKey)

	nonce, _ := client.NonceAt(context.Background(), MainAddr, nil)
	auth := bind.NewKeyedTransactor(MainKey)
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4000000) // in units
	auth.GasPrice = big.NewInt(250000000000000)

	// init trc21 issuer
	auth.Nonce = big.NewInt(int64(nonce))

	price := new(big.Int)
	price.SetString(os.Getenv("PRICE"), 10)

	lendContract, _ := tomox.NewLendingRelayerRegistration(auth, common.HexToAddress(os.Getenv("LENDING_ADDRESS")), client)

	token := common.HexToAddress(os.Getenv("TOKEN_ADDRESS"))

	tx, err := lendContract.SetCollateralPrice(token, price)
	if err != nil {
		fmt.Println("Set price failed!", err)
	}

	time.Sleep(5 * time.Second)
	r, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		fmt.Println("Get receipt failed", err)
	}
	fmt.Println("Done receipt status", r.Status)

}
