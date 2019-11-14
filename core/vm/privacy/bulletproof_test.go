package privacy

import (
	"fmt"
	"math/big"
	"testing"
)

func TestInnerProductProveLen1(t *testing.T) {
	fmt.Println("TestInnerProductProve1")
	EC = NewECPrimeGroupKey(1)
	a := make([]*big.Int, 1)
	b := make([]*big.Int, 1)

	a[0] = big.NewInt(2)

	b[0] = big.NewInt(2)

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
