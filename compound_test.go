package art_test

import (
	"fmt"
	"math/bits"
	"slices"
	"testing"

	"github.com/Clement-Jean/go-art"
)

type Account struct {
	ID   uint
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
		b   []byte
	)

	_, c := ubk.Transform(a.ID)
	_, d := aok.Transform(a.name)
	b = append(b, c...)
	b = append(b, d...)
	return b, b
}
func (ak AccountKey) Restore(b []byte) Account {
	var (
		a   Account
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

func TestCompoundRange(t *testing.T) {
	tests := []struct {
		name           string
		start, end     Account
		keys, expected []Account
	}{
		{
			name:  "start<end",
			start: Account{ID: 0}, end: Account{ID: 7},
			keys: []Account{
				{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7},
				{ID: 8}, {ID: 9}, {ID: 10}, {ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15},
			},
			expected: []Account{
				{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7},
			},
		},
		{
			name:  "start>end",
			start: Account{ID: 7}, end: Account{ID: 0},
			keys: []Account{
				{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7},
				{ID: 8}, {ID: 9}, {ID: 10}, {ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15},
			},
			expected: []Account{
				{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7},
			},
		},
		{
			name:  "start==end",
			start: Account{ID: 7}, end: Account{ID: 7},
			keys: []Account{
				{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7},
				{ID: 8}, {ID: 9}, {ID: 10}, {ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15},
			},
			expected: []Account{{ID: 7}},
		},
		{
			name:  "outside of range",
			start: Account{ID: 16}, end: Account{ID: 20},
			keys: []Account{
				{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7},
				{ID: 8}, {ID: 9}, {ID: 10}, {ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15},
			},
			expected: []Account{},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("range-%s", tt.name), func(t *testing.T) {
			acc := AccountKey{}
			tr := art.NewCompoundTree[Account, Account](acc)

			for _, key := range tt.keys {
				tr.Insert(key, key)
			}

			var res []Account
			for key, _ := range tr.Range(tt.start, tt.end) {
				res = append(res, key)
			}

			if !slices.Equal(tt.expected, res) {
				fmt.Printf("%v %v\n", tt.expected, res)
				t.Fatal("slices are not the same")
			}
		})
	}
}
