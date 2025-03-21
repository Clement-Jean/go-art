package art_test

import (
	"slices"
	"testing"

	"github.com/Clement-Jean/go-art"
)

func TestBinaryUnsignedAll(t *testing.T) {
	tr := art.NewUnsignedBinaryTree[uint, int]()
	expected := []uint{1, 9, 100_000, 11}

	slices.Sort(expected)
	tr.Insert(1, 1)
	tr.Insert(11, 1)
	tr.Insert(100_000, 1)
	tr.Insert(9, 1)

	var got []uint
	for k, _ := range tr.All() {
		got = append(got, k)
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestBinaryUnsignedBackward(t *testing.T) {
	tr := art.NewUnsignedBinaryTree[uint, int]()
	expected := []uint{1, 9, 100_000, 11}

	slices.Sort(expected)
	slices.Reverse(expected)
	tr.Insert(1, 1)
	tr.Insert(11, 1)
	tr.Insert(100_000, 1)
	tr.Insert(9, 1)

	var got []uint
	for k, _ := range tr.Backward() {
		got = append(got, k)
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
