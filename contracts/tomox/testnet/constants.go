package testnet

import (
	"math/big"
	"os"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/crypto"
)

var (
	BaseTOMO    = big.NewInt(0).Mul(big.NewInt(10), big.NewInt(100000000000000000)) // 1 TOMO
	RpcEndpoint = "http://127.0.0.1:8545/"
	MainKey, _  = crypto.HexToECDSA(os.Getenv("MAIN_ADDRESS_KEY"))
	MainAddr    = crypto.PubkeyToAddress(MainKey.PublicKey)

	// TRC21 Token
	MinTRC21Apply   = big.NewInt(0).Mul(big.NewInt(10), BaseTOMO) // 10 TOMO
	TRC21TokenCap   = big.NewInt(0).Mul(big.NewInt(1000000000000), BaseTOMO)
	TRC21TokenFee   = big.NewInt(100)
	TomoXListingFee = big.NewInt(0).Mul(big.NewInt(1000), BaseTOMO) // 1000 TOMO

	// TOMOX
	MaxRelayers  = big.NewInt(200)
	MaxTokenList = big.NewInt(200)
	MinDeposit   = big.NewInt(0).Mul(big.NewInt(25000), BaseTOMO) // 25000 TOMO
	TradeFee     = uint16(10)                                     // trade fee decimals 10

	RelayerCoinbaseKey, _ = crypto.HexToECDSA(os.Getenv("RELAYER_COINBASE_KEY")) //
	RelayerCoinbaseAddr   = crypto.PubkeyToAddress(RelayerCoinbaseKey.PublicKey)

	OwnerRelayerKey, _ = crypto.HexToECDSA(os.Getenv("RELAYER_OWNER_KEY"))
	OwnerRelayerAddr   = crypto.PubkeyToAddress(OwnerRelayerKey.PublicKey)

	TOMONative = common.HexToAddress("0x0000000000000000000000000000000000000001")

	TokenNameList = []string{"BTC", "ETH", "XRP", "LTC", "BNB", "ADA", "ETC", "BCH", "EOS", "USD"}
	TeamAddresses = []common.Address{
		common.HexToAddress("0xE3584D2D430eF34FF9fEeCBEBE6E0f6980082F05"), // Test1
		common.HexToAddress("0x16a73f3a64eca79e117258e66dfd7071cc8312a9"), // BTCTOMO
		common.HexToAddress("0xac177441ac2237b2f79ecff1b8f6bca39e27ef9f"), // ETHTOMO
		common.HexToAddress("0x4215250e55984c75bbce8ae639b86a6cad8ec126"), // XRPTOMO
		common.HexToAddress("0x6b70ca959814866dd5c426d63d47dde9cc6c32d2"), // LTCTOMO
		common.HexToAddress("0x33df079fe9b9cd7fb23a1085e4eaaa8eb6952cb3"), // BNBTOMO
		common.HexToAddress("0x3cab8292137804688714670640d19f9d7a60c472"), // ADATOMO
		common.HexToAddress("0x9415d953d47c5f155cac9de7b24a756f352eafbf"), // ETCTOMO
		common.HexToAddress("0xe32d2e7c8e8809e45c8e2332830b48d9e231e3f2"), // BCHTOMO
		common.HexToAddress("0xf76ddbda664ea47088937e1cf9ff15036714dee3"), // EOSTOMO
		common.HexToAddress("0xc465ee82440dada9509feb235c7cd7d896acf13c"), // ETHBTC
		common.HexToAddress("0xb95bdc136c579dc3fd2b2424a8e925a90228d2c2"), // XRPBTC
		common.HexToAddress("0xe36c1842365595D44854eEcd64B11c8115E133EF"), // USDTOMO
		common.HexToAddress("0xaaC1959F6F0fb539F653409079Ec4146267B7555"), // BTCUSD
		common.HexToAddress("0x726DA688e2e09f01A2e1aB4c10F25B7CEdD4a0f3"), // ETHUSD
	}
)
