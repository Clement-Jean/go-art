package art_test

import (
	"math"
	"slices"
	"testing"

	"github.com/Clement-Jean/go-art"
)

func TestBinaryFloatAll(t *testing.T) {
	tr := art.NewFloatBinaryTree[float64, int]()
	expected := []float64{-1, -2, 0 /*math.NaN()*/, math.Inf(1), math.Inf(-1), 1, 2}

	slices.Sort(expected)
	tr.Insert(1, 1)
	tr.Insert(2, 1)
	// tr.Insert(math.NaN(), 1)
	tr.Insert(math.Inf(1), 1)
	tr.Insert(math.Inf(-1), 1)
	tr.Insert(-1, 1)
	tr.Insert(-2, 1)
	tr.Insert(-0, 1)
	tr.Insert(+0, 1)

	var got []float64
	for k, _ := range tr.All() {
		got = append(got, k)
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestBinaryFloatBackward(t *testing.T) {
	tr := art.NewFloatBinaryTree[float64, int]()
	expected := []float64{-1, -2, 0 /*math.NaN()*/, math.Inf(1), math.Inf(-1), 1, 2}

	slices.Sort(expected)
	tr.Insert(1, 1)
	tr.Insert(2, 1)
	// tr.Insert(math.NaN(), 1)
	tr.Insert(math.Inf(1), 1)
	tr.Insert(math.Inf(-1), 1)
	tr.Insert(-1, 1)
	tr.Insert(-2, 1)
	tr.Insert(-0, 1)
	tr.Insert(+0, 1)

	var got []float64
	for k, _ := range tr.Backward() {
		got = append(got, k)
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
