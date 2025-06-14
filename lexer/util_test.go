package lexer

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"testing"
)

func TestIsBase16(t *testing.T) {
	numMap := map[string]int64{
		"deadbad":  0xdeadbad,
		"deadbeef": 0xdeadbeef,
		"10c":      0x10c,
		"785c":     0x785c,
	}

	for k, v := range numMap {
		num, err := strconv.ParseInt(k, 16, 0)
		if err != nil {
			t.Error(err)
		}

		if num != v {
			t.Errorf("string key converted is different from int value: %d != %d\n", num, v)
		}

		for i := range k {
			if !IsBase16(k[i]) {
				t.Errorf("IsBase16 failed: IsBase16(%s) = %T\n", k, false)
			}
		}

		rn, err := rand.Int(rand.Reader, big.NewInt(0x07FFFFFFFFFFFFFF))
		if err != nil {
			t.Error(err)
		}

		for j := 0; j < 10; j++ {
			srn := fmt.Sprintf("%x", rn.Int64())
			for i := range srn {
				if !IsBase16(srn[i]) {
					t.Errorf("IsBase16 failed: IsBase16(%s) = %T\n", srn, false)
				}
			}
		}
	}
}

func TestIsBase8(t *testing.T) {
	numMap := map[string]int64{
		"010":   010,
		"020":   020,
		"01000": 01000,
		"075":   075,
	}

	for k, v := range numMap {
		num, err := strconv.ParseInt(k, 8, 0)
		if err != nil {
			t.Error(err)
		}

		if num != v {
			t.Errorf("string key converted is different from int value: %d != %d\n", num, v)
		}

		for i := range k {
			if !IsBase8(k[i]) {
				t.Errorf("IsBase8 failed: IsBase8(%s) = %T\n", k, false)
			}
		}

		rn, err := rand.Int(rand.Reader, big.NewInt(0x07FFFFFFFFFFFFFF))
		if err != nil {
			t.Error(err)
		}

		for j := 0; j < 10; j++ {
			srn := fmt.Sprintf("%o", rn.Int64())
			for i := range srn {
				if !IsBase8(srn[i]) {
					t.Errorf("IsBase8 failed: IsBase8(%s) = %T\n", srn, false)
				}
			}
		}
	}
}

// func TestIsBase2(t *testing.T) {
// 	numMap := map[string]int64{
// 		"0b10110110": 0b10110110,
// 		"0b10001101": 0b10001101,
// 		"0b1000":     0b1000,
// 		"0b0":        0b0,
// 	}

// 	for k, v := range numMap {
// 		num, err := strconv.ParseInt(strings.ReplaceAll(k, "0b", ""), 2, 0)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		if num != v {
// 			t.Errorf("string key converted is different from int value: %d != %d\n", num, v)
// 		}

// 		for i := range k {
// 			if !IsBase2(k[i]) {
// 				t.Errorf("IsBase2 failed: IsBase2(%s) = %t\n", k, false)
// 			}
// 		}

// 		rn, err := rand.Int(rand.Reader, big.NewInt(0x07FFFFFFFFFFFFFF))
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		for j := 0; j < 10; j++ {
// 			srn := strings.ReplaceAll(fmt.Sprintf("%b", rn.Int64()), "0b", "")
// 			for i := range srn {
// 				if !IsBase2(srn[i]) {
// 					t.Errorf("IsBase2 failed: IsBase2(%s) = %T\n", srn, false)
// 				}
// 			}
// 		}
// 	}
// }
