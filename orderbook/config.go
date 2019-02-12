package orderbook

import "math/big"

type Config struct {
	DataDir string `toml:",omitempty"`
}

var AllowedPairs = map[string]*big.Int{
	"TOMO/WETH": big.NewInt(10e9),
}
