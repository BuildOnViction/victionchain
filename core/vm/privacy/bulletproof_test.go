package privacy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/btcsuite/btcd/btcec"
)

func genECPrimeGroupKey(n int) CryptoParams {
	// curValue := btcec.S256().Gx
	// s256 := sha256.New()
	gen1Vals := make([]ECPoint, n)
	gen2Vals := make([]ECPoint, n)
	// u := ECPoint{big.NewInt(0), big.NewInt(0)}
	hx, _ := new(big.Int).SetString("50929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0", 16)
	hy, _ := new(big.Int).SetString("31d3c6863973926e049e637cb1b5f40a36dac28af1766968c30c2313f3a38904", 16)
	ch := ECPoint{hx, hy}

	gx, _ := new(big.Int).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	gy, _ := new(big.Int).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	cg := ECPoint{gx, gy}

	i := 0
	for i < n {
		gen2Vals[i] = ch.Mult(
			big.NewInt(int64(i*2 + 1)),
		)
		gen1Vals[i] = cg.Mult(
			big.NewInt(int64(i*2 + 2)),
		)
		i++
	}

	u := cg.Mult(
		big.NewInt(int64(n + 3)),
	)

	return CryptoParams{
		btcec.S256(),
		btcec.S256(),
		gen1Vals,
		gen2Vals,
		btcec.S256().N,
		u,
		n,
		cg,
		ch}
}

func TestInnerProductProveLen1(t *testing.T) {
	fmt.Println("TestInnerProductProve1")
	EC = NewECPrimeGroupKey(1)
	a := make([]*big.Int, 1)
	b := make([]*big.Int, 1)

	a[0] = big.NewInt(1)

	b[0] = big.NewInt(1)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerify(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
	}
}

func TestInnerProductProveLen2(t *testing.T) {
	fmt.Println("TestInnerProductProve2")
	EC = genECPrimeGroupKey(2)
	a := make([]*big.Int, 2)
	b := make([]*big.Int, 2)

	a[0] = big.NewInt(1)
	a[1] = big.NewInt(1)

	b[0] = big.NewInt(1)
	b[1] = big.NewInt(1)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	fmt.Println("P after two vector commitment with gen ", P)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerify(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
	}
}

func TestInnerProductProveLen4(t *testing.T) {
	fmt.Println("TestInnerProductProve4")
	EC = NewECPrimeGroupKey(4)
	a := make([]*big.Int, 4)
	b := make([]*big.Int, 4)

	a[0] = big.NewInt(1)
	a[1] = big.NewInt(1)
	a[2] = big.NewInt(1)
	a[3] = big.NewInt(1)

	b[0] = big.NewInt(1)
	b[1] = big.NewInt(1)
	b[2] = big.NewInt(1)
	b[3] = big.NewInt(1)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerify(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
	}
}

func TestInnerProductProveLen8(t *testing.T) {
	fmt.Println("TestInnerProductProve8")
	EC = NewECPrimeGroupKey(8)
	a := make([]*big.Int, 8)
	b := make([]*big.Int, 8)

	a[0] = big.NewInt(1)
	a[1] = big.NewInt(1)
	a[2] = big.NewInt(1)
	a[3] = big.NewInt(1)
	a[4] = big.NewInt(1)
	a[5] = big.NewInt(1)
	a[6] = big.NewInt(1)
	a[7] = big.NewInt(1)

	b[0] = big.NewInt(2)
	b[1] = big.NewInt(2)
	b[2] = big.NewInt(2)
	b[3] = big.NewInt(2)
	b[4] = big.NewInt(2)
	b[5] = big.NewInt(2)
	b[6] = big.NewInt(2)
	b[7] = big.NewInt(2)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerify(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
	}
}

func TestInnerProductProveLen64Rand(t *testing.T) {
	fmt.Println("TestInnerProductProveLen64Rand")
	EC = NewECPrimeGroupKey(64)
	a := RandVector(64)
	b := RandVector(64)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerify(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
		fmt.Printf("Values Used: \n\ta = %s\n\tb = %s\n", a, b)
	}

}

func TestInnerProductVerifyFastLen1(t *testing.T) {
	fmt.Println("TestInnerProductProve1")
	EC = NewECPrimeGroupKey(1)
	a := make([]*big.Int, 1)
	b := make([]*big.Int, 1)

	a[0] = big.NewInt(2)

	b[0] = big.NewInt(2)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerifyFast(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
	}
}

func TestInnerProductVerifyFastLen2(t *testing.T) {
	fmt.Println("TestInnerProductProve2")
	EC = NewECPrimeGroupKey(2)
	a := make([]*big.Int, 2)
	b := make([]*big.Int, 2)

	a[0] = big.NewInt(2)
	a[1] = big.NewInt(3)

	b[0] = big.NewInt(2)
	b[1] = big.NewInt(3)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerifyFast(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
	}
}

func TestInnerProductVerifyFastLen4(t *testing.T) {
	fmt.Println("TestInnerProductProve4")
	EC = NewECPrimeGroupKey(4)
	a := make([]*big.Int, 4)
	b := make([]*big.Int, 4)

	a[0] = big.NewInt(1)
	a[1] = big.NewInt(1)
	a[2] = big.NewInt(1)
	a[3] = big.NewInt(1)

	b[0] = big.NewInt(1)
	b[1] = big.NewInt(1)
	b[2] = big.NewInt(1)
	b[3] = big.NewInt(1)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerifyFast(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
	}
}

func TestInnerProductVerifyFastLen8(t *testing.T) {
	fmt.Println("TestInnerProductProve8")
	EC = NewECPrimeGroupKey(8)
	a := make([]*big.Int, 8)
	b := make([]*big.Int, 8)

	a[0] = big.NewInt(1)
	a[1] = big.NewInt(1)
	a[2] = big.NewInt(1)
	a[3] = big.NewInt(1)
	a[4] = big.NewInt(1)
	a[5] = big.NewInt(1)
	a[6] = big.NewInt(1)
	a[7] = big.NewInt(1)

	b[0] = big.NewInt(2)
	b[1] = big.NewInt(2)
	b[2] = big.NewInt(2)
	b[3] = big.NewInt(2)
	b[4] = big.NewInt(2)
	b[5] = big.NewInt(2)
	b[6] = big.NewInt(2)
	b[7] = big.NewInt(2)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerifyFast(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
	}
}

func TestInnerProductVerifyFastLen64Rand(t *testing.T) {
	fmt.Println("TestInnerProductProveLen64Rand")
	EC = NewECPrimeGroupKey(64)
	a := RandVector(64)
	b := RandVector(64)

	c := InnerProduct(a, b)

	P := TwoVectorPCommitWithGens(EC.BPG, EC.BPH, a, b)

	ipp := InnerProductProve(a, b, c, P, EC.U, EC.BPG, EC.BPH)

	if InnerProductVerifyFast(c, P, EC.U, EC.BPG, EC.BPH, ipp) {
		fmt.Println("Inner Product Proof correct")
	} else {
		t.Error("Inner Product Proof incorrect")
		fmt.Printf("Values Used: \n\ta = %s\n\tb = %s\n", a, b)
	}

}

func TestMRPProve(t *testing.T) {
	EC = genECPrimeGroupKey(64)
	mRangeProof := MRPProve([]*big.Int{
		new(big.Int).SetInt64(327680),
	})
	// fmt.Printf("Value is : %s %s\n", 0x9999999999, 0x9999999999)
	fmt.Printf("\n\n\n mRangeProof %+v\n\n\n", mRangeProof)

	// expectedMRangeProof := parseTestData("./bulletproof.json")
	// fmt.Printf("\n\n\n expectedMRangeProof %+v\n\n\n", expectedMRangeProof)

	mv := MRPVerify(mRangeProof)
	fmt.Printf("Value is between 1 and 2^%d-1: %t\n", VecLength, mv)

	if mv {
		fmt.Println("MRProof correct")
	} else {
		t.Error("MRProof incorrect")
	}
}

type Point struct {
	x string
	y string
}

type IPP struct {
	L          []map[string]string `json:"L"`
	R          []map[string]string `json:"R"`
	A          string              `json:"A"`
	B          string              `json:"B"`
	Challenges []string            `json:"Challenges"`
}

type BulletProof struct {
	Comms []string `json:"Comms"`
	A     string   `json:"A"`
	S     string   `json:"S"`
	Cx    string   `json:"Cx"`
	Cy    string   `json:"Cy"`
	Cz    string   `json:"Cz"`
	T1    string   `json:"T1"`
	T2    string   `json:"T2"`
	Th    string   `json:"Th"`
	Tau   string   `json:"Tau"`
	Mu    string   `json:"Mu"`
	Ipp   IPP      `json:"Ipp"`
}

func parseTestData(filePath string) MultiRangeProof {
	jsonFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	// var result map[string]interface{}
	result := BulletProof{}

	json.Unmarshal([]byte(byteValue), &result)

	fmt.Println("result ", result.Tau)
	fmt.Println("result ", result.Th)
	fmt.Println("result.Ipp ", result.Ipp)

	ipp := result.Ipp

	proof := MultiRangeProof{
		Comms: MapECPointFromHex(result.Comms, ECPointFromHex),
		A:     ECPointFromHex(result.A),
		S:     ECPointFromHex(result.S),
		T1:    ECPointFromHex(result.T1),
		T2:    ECPointFromHex(result.T2),
		Th:    bigIFromHex(result.Th),
		Tau:   bigIFromHex(result.Tau),
		Mu:    bigIFromHex(result.Mu),
		Cx:    bigIFromHex(result.Cx),
		Cy:    bigIFromHex(result.Cy),
		Cz:    bigIFromHex(result.Cz),
		IPP: InnerProdArg{
			L:          MapECPoint(ipp.L, ECPointFromPoint),
			R:          MapECPoint(ipp.R, ECPointFromPoint),
			A:          bigIFromHex(ipp.A),
			B:          bigIFromHex(ipp.B),
			Challenges: MapBigI(ipp.Challenges, bigIFromHex),
		},
	}

	fmt.Println(proof)
	return proof
}

func MapBigI(list []string, f func(string) *big.Int) []*big.Int {
	result := make([]*big.Int, len(list))

	for i, item := range list {
		result[i] = f(item)
	}
	return result
}

func MapECPointFromHex(list []string, f func(string) ECPoint) []ECPoint {
	result := make([]ECPoint, len(list))

	for i, item := range list {
		result[i] = f(item)
	}
	return result
}

func MapECPoint(list []map[string]string, f func(Point) ECPoint) []ECPoint {
	result := make([]ECPoint, len(list))

	for i, item := range list {
		result[i] = f(Point{
			x: item["x"],
			y: item["y"],
		})
	}
	return result
}

func bigIFromHex(hex string) *big.Int {
	tmp, _ := new(big.Int).SetString(hex, 16)
	return tmp
}

func ECPointFromHex(hex string) ECPoint {
	Px, _ := new(big.Int).SetString(hex[:64], 16)
	Py, _ := new(big.Int).SetString(hex[64:], 16)
	P := ECPoint{Px, Py}
	return P
}

func ECPointFromPoint(ecpoint Point) ECPoint {
	Px, _ := new(big.Int).SetString(ecpoint.x, 16)
	Py, _ := new(big.Int).SetString(ecpoint.y, 16)
	P := ECPoint{Px, Py}
	return P
}

func TestMRPProveFromJS(t *testing.T) {
	EC = genECPrimeGroupKey(64)
	mRangeProof := parseTestData("./bulletproof.json")
	fmt.Printf("\n\n\n%+v\n\n\n", mRangeProof)
	mv := MRPVerify(mRangeProof)
	fmt.Printf("Value is between 1 and 2^%d-1: %t\n", VecLength, mv)

	if mv {
		fmt.Println("MRProof synced JS correct")
	} else {
		t.Error("MRProof synced JS incorrect")
	}
}
