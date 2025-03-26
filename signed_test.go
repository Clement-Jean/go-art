package art_test

import (
	"slices"
	"testing"

	"github.com/Clement-Jean/go-art"
)

func TestBinarySignedAll(t *testing.T) {
	tr := art.NewSignedBinaryTree[int, int]()
	expected := []int{-1, -9, -100_000, 11}

	slices.Sort(expected)
	tr.Insert(-1, 1)
	tr.Insert(11, 1)
	tr.Insert(-100_000, 1)
	tr.Insert(-9, 1)

	var got []int
	for k, _ := range tr.All() {
		got = append(got, k)
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestBinarySignedBackward(t *testing.T) {
	tr := art.NewSignedBinaryTree[int, int]()
	expected := []int{-1, -9, -100_000, 11}

	slices.Sort(expected)
	slices.Reverse(expected)
	tr.Insert(-1, 1)
	tr.Insert(11, 1)
	tr.Insert(-100_000, 1)
	tr.Insert(-9, 1)

	var got []int
	for k, _ := range tr.Backward() {
		got = append(got, k)
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
