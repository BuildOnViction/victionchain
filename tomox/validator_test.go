package tomox

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"math/big"
	"testing"
)

var fakeDb, _ = ethdb.NewMemDatabase()

func TestIsValidRelayer(t *testing.T) {
	order := &tomox_state.OrderItem{
		ExchangeAddress: common.HexToAddress("relayer1"),
	}
	var stateDb, _ = state.New(common.Hash{}, state.NewDatabase(fakeDb))
	slotKec := crypto.Keccak256(order.ExchangeAddress.Hash().Bytes(), common.BigToHash(new(big.Int).SetUint64(tomox_state.RelayerMappingSlot["RELAYER_LIST"])).Bytes())
	locRelayerState := new(big.Int).SetBytes(slotKec)
	stateDb.SetState(common.HexToAddress(common.RelayerRegistrationSMC), common.BigToHash(locRelayerState), common.BigToHash(new(big.Int).SetUint64(0)))
	if valid := IsValidRelayer(stateDb, order.ExchangeAddress); valid {
		t.Error("TestIsValidRelayer FAILED. It should be invalid relayer", "ExchangeAddress", order.ExchangeAddress)
	}

	stateDb.SetState(common.HexToAddress(common.RelayerRegistrationSMC), common.BigToHash(locRelayerState), common.BigToHash(new(big.Int).SetUint64(2500)))
	if valid := IsValidRelayer(stateDb, order.ExchangeAddress); !valid {
		t.Error("TestIsValidRelayer FAILED. This address should be a valid relayer", "ExchangeAddress", order.ExchangeAddress)
	}
}

// test full-verify orderItem
// VerifyOrder consists of some partial tests:
// verifyTimestamp, verifyOrderSide, verifyOrderType, verifyRelayer, verifyBalance, verifySignature
// in this test, we sequentially make each test PASS
func TestOrderItem_VerifyBasicOrderInfo(t *testing.T) {
	addr := common.HexToAddress("0x0332d186212b04E6933682b3bed8e232b6b3361a")

	order := &tomox_state.OrderItem{
		PairName:    "TOMO/WETH",
		BaseToken:   common.HexToAddress(common.TomoNativeAddress),
		QuoteToken:  common.HexToAddress("0x0aaad186212b04E6933682b3bed8e232b6b3361a"),
		UserAddress: addr,
		Nonce:       big.NewInt(1),
		Quantity:    big.NewInt(1000),
		Price:       big.NewInt(1100),
	}

	if err := order.VerifyBasicOrderInfo(); err != tomox_state.ErrInvalidOrderSide {
		t.Error(err)
	}

	// set valid orderSide: Ask
	order.Side = Ask

	// after verifyOrderSide PASS, order should fail the next step: verifyOrderType
	if err := order.VerifyBasicOrderInfo(); err != tomox_state.ErrInvalidOrderType {
		t.Error(err)
	}

	// set valid orderType
	order.Type = Limit

	// wrong hash
	if err := order.VerifyBasicOrderInfo(); err != tomox_state.ErrWrongHash {
		t.Error(err)
	}

	// set a wrong signature
	privKey, _ := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	pubKey := privKey.Public()
	pubKeyECDSA, _ := pubKey.(*ecdsa.PublicKey)
	pubKeyBytes := crypto.FromECDSAPub(pubKeyECDSA)

	copy(addr[:], crypto.Keccak256(pubKeyBytes[1:])[12:])
	order.UserAddress = addr

	// set valid hash
	order.Hash = order.ComputeHash()

	signatureBytes, _ := crypto.Sign(common.StringToHash("invalid hash").Bytes(), privKey)
	sig := &tomox_state.Signature{
		R: common.BytesToHash(signatureBytes[0:32]),
		S: common.BytesToHash(signatureBytes[32:64]),
		V: signatureBytes[64] + 27,
	}
	order.Signature = sig

	// wrong signature
	if err := order.VerifyBasicOrderInfo(); err != tomox_state.ErrInvalidSignature {
		t.Error(err)
	}

	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		order.Hash.Bytes(),
	)

	// set valid signature
	signatureBytes, _ = crypto.Sign(message, privKey)
	sig = &tomox_state.Signature{
		R: common.BytesToHash(signatureBytes[0:32]),
		S: common.BytesToHash(signatureBytes[32:64]),
		V: signatureBytes[64] + 27,
	}
	order.Signature = sig

	// Finally, we have a valid order
	if err := order.VerifyBasicOrderInfo(); err != nil {
		t.Error(err)
	}

}

func TestTxDataMatch_DecodeOrder(t *testing.T) {
	txDataMatch := &TxDataMatch{
		Order: []byte("abc"),
	}
	var err error
	if _, err = txDataMatch.DecodeOrder(); err == nil {
		t.Error("It should fail")
	}

	orderItem := &tomox_state.OrderItem{
		PairName:        "TOMO/WETH",
		Price:           big.NewInt(1),
		Quantity:        big.NewInt(100),
		Type:            Limit,
		Side:            Bid,
		UserAddress:     common.HexToAddress("aaa"),
		ExchangeAddress: common.HexToAddress("bbb"),
		Signature:       &tomox_state.Signature{},
	}
	b, err := EncodeBytesItem(orderItem)
	if err != nil {
		t.Error(err)
	}
	txDataMatch.Order = b
	if _, err = txDataMatch.DecodeOrder(); err != nil {
		t.Error(err)
	}
}
