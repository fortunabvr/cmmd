// Copyright (c) 2018 The Commercium developers
package cmmutil_test

import (
	"fmt"
	"math"

	"github.com/CommerciumBlockchain/cmmd/cmmutil"
)

func ExampleAmount() {

	a := cmmutil.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = cmmutil.Amount(1e8)
	fmt.Println("100,000,000 Atoms:", a)

	a = cmmutil.Amount(1e5)
	fmt.Println("100,000 Atoms:", a)
	// Output:
	// Zero Atom: 0 CMM
	// 100,000,000 Atoms: 1 CMM
	// 100,000 Atoms: 0.001 CMM
}

func ExampleNewAmount() {
	amountOne, err := cmmutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := cmmutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := cmmutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := cmmutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 CMM
	// 0.01234567 CMM
	// 0 CMM
	// invalid coin amount
}

func ExampleAmount_unitConversions() {
	amount := cmmutil.Amount(44433322211100)

	fmt.Println("Atom to kCoin:", amount.Format(cmmutil.AmountKiloCoin))
	fmt.Println("Atom to Coin:", amount)
	fmt.Println("Atom to MilliCoin:", amount.Format(cmmutil.AmountMilliCoin))
	fmt.Println("Atom to MicroCoin:", amount.Format(cmmutil.AmountMicroCoin))
	fmt.Println("Atom to Atom:", amount.Format(cmmutil.AmountAtom))

	// Output:
	// Atom to kCoin: 444.333222111 kCMM
	// Atom to Coin: 444333.222111 CMM
	// Atom to MilliCoin: 444333222.111 mCMM
	// Atom to MicroCoin: 444333222111 μCMM
	// Atom to Atom: 44433322211100 Atom
}
