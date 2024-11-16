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
