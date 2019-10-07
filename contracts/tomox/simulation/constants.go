package simulation

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	BaseTOMO    = big.NewInt(0).Mul(big.NewInt(10), big.NewInt(100000000000000000)) // 1 TOMO
	RpcEndpoint = "http://127.0.0.1:8501/"
	MainKey, _  = crypto.HexToECDSA(os.Getenv("MAIN_ADDRESS_KEY"))
	MainAddr    = crypto.PubkeyToAddress(MainKey.PublicKey) //0x17F2beD710ba50Ed27aEa52fc4bD7Bda5ED4a037

	// TRC21 Token
	MinTRC21Apply = big.NewInt(0).Mul(big.NewInt(100), BaseTOMO) // 100 TOMO
	TRC21TokenCap = big.NewInt(0).Mul(big.NewInt(1000000000000), BaseTOMO)
	TRC21TokenFee = big.NewInt(100)

	// TOMOX
	MaxRelayers  = big.NewInt(100)
	MaxTokenList = big.NewInt(150)
	MinDeposit   = big.NewInt(25000) // 25000 TOMO
	TradeFee     = uint16(1)

	RelayerCoinbaseKey, _ = crypto.HexToECDSA(os.Getenv("RELAYER_COINBASE_KEY")) //
	RelayerCoinbaseAddr   = crypto.PubkeyToAddress(RelayerCoinbaseKey.PublicKey) // 0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e

	OwnerRelayerKey, _ = crypto.HexToECDSA(os.Getenv("RELAYER_OWNER_KEY"))
	OwnerRelayerAddr   = crypto.PubkeyToAddress(OwnerRelayerKey.PublicKey) //0x703c4b2bD70c169f5717101CaeE543299Fc946C7

	TOMONative = common.HexToAddress("0x0000000000000000000000000000000000000001")

	TokenNameList = []string{"BTC", "ETH", "XRP", "LTC", "BNB", "ADA", "ETC", "BCH", "EOS"}
	TeamAddresses = []common.Address{
		// common.HexToAddress("0x8fB1047e874d2e472cd08980FF8383053dd83102"), // MM team
		// common.HexToAddress("0x9ca1514E3Dc4059C29a1608AE3a3E3fd35900888"), // MM team
		// common.HexToAddress("0x15e08dE16f534c890828F2a0D935433aF5B3CE0C"), // CTO
		// common.HexToAddress("0xb68D825655F2fE14C32558cDf950b45beF18D218"), // DEX team
		common.HexToAddress("0xF7349C253FF7747Df661296E0859c44e974fb52E"), // HaiDV
		common.HexToAddress("0x9f6b8fDD3733B099A91B6D70CDC7963ebBbd2684"), // Can
		common.HexToAddress("0x06605B28aab9835be75ca242a8aE58f2e15F2F45"), // Nien
	}
)
