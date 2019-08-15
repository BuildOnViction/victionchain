package simulation

import (
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var (
	BaseTOMO    = big.NewInt(0).Mul(big.NewInt(10), big.NewInt(100000000000000000)) // 1 TOMO
	RpcEndpoint = "http://127.0.0.1:8501/"
	MainKey, _  = crypto.HexToECDSA("65ec4d4dfbcac594a14c36baa462d6f73cd86134840f6cf7b80a1e1cd33473e2")
	MainAddr    = crypto.PubkeyToAddress(MainKey.PublicKey) //0x17F2beD710ba50Ed27aEa52fc4bD7Bda5ED4a037

	// TRC21 Token
	MinTRC21Apply = big.NewInt(0).Mul(big.NewInt(100), BaseTOMO) // 100 TOMO
	TRC21TokenCap = big.NewInt(0).Mul(big.NewInt(1000000000000), BaseTOMO)
	TRC21TokenFee = big.NewInt(100)

	// TOMOX
	MaxRelayers  = big.NewInt(200)
	MaxTokenList = big.NewInt(200)
	MinDeposit   = big.NewInt(25000) // 25000 TOMO
	TradeFee     = uint16(1)

	ReplayKey, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	ReplayAddr   = crypto.PubkeyToAddress(ReplayKey.PublicKey) // 0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e

	OwnerRelayerKey, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	OwnerRelayAddr     = crypto.PubkeyToAddress(OwnerRelayerKey.PublicKey) //0x703c4b2bD70c169f5717101CaeE543299Fc946C7
)
