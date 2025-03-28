package art_test

import (
	"fmt"
	"math/bits"
	"testing"

	"github.com/Clement-Jean/go-art"
)

type Account struct {
	ID uint
	name string
}

func (acc *Account) String() string {
	return fmt.Sprintf("{id: %d, name: %q}", acc.ID, acc.name)
}

type AccountKey struct{}

func (ak AccountKey) Transform(a Account) ([]byte, []byte) {
	var (
		ubk art.UnsignedBinaryKey[uint]
		aok art.AlphabeticalOrderKey[string]
		b []byte
	)

	_, c := ubk.Transform(a.ID)
	_, d := aok.Transform(a.name)
	b = append(b, c...)
	b = append(b, d...)
	return b, b
}
func (ak AccountKey) Restore(b []byte) Account {
	var (
		a Account
		ubk art.UnsignedBinaryKey[uint]
		aok art.AlphabeticalOrderKey[string]
	)

	id := b[:(bits.UintSize / 8)]
	name := b[(bits.UintSize / 8):]

	a.ID = ubk.Restore(id)
	a.name = aok.Restore(name)
	return a
}

func TestCompound(t *testing.T) {
	var ak AccountKey

	acc := Account{ID: 1, name: "Clement"}
	tr := art.NewCompoundTree[Account, int](ak)
	tr.Insert(acc, 1)

	for k, v := range tr.All() {
		if k.ID != 1 {
			t.Fatalf("expected ID of 1, got %d", k.ID)
		}

		if k.name != "Clement" {
			t.Fatalf("expected name of Clement, got %s", k.name)
		}

		if v != 1 {
			t.Fatalf("expected value of 1, got %d", v)
		}
	}
}
