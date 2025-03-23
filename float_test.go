package art_test

import (
	"math"
	"slices"
	"testing"

	"github.com/Clement-Jean/go-art"
)

const float64EqualityThreshold = 1e-9

func TestBinaryFloatAll32(t *testing.T) {
	tr := art.NewFloatBinaryTree[float32, int]()
	expected := []float32{-1, -2, 0, float32(math.NaN()), float32(math.Inf(1)), float32(math.Inf(-1)), 1, 2, math.MaxFloat32, math.SmallestNonzeroFloat32}

	slices.Sort(expected)
	tr.Insert(1, 1)
	tr.Insert(2, 1)
	tr.Insert(float32(math.NaN()), 1)
	tr.Insert(float32(math.Inf(1)), 1)
	tr.Insert(float32(math.Inf(-1)), 1)
	tr.Insert(-1, 1)
	tr.Insert(-2, 1)
	tr.Insert(math.MaxFloat32, 1)
	tr.Insert(math.SmallestNonzeroFloat32, 1)
	tr.Insert(0, 1)
	tr.Insert(-0, 1)
	tr.Insert(+0, 1)

	var got []float32
	for k, _ := range tr.All() {
		got = append(got, k)
	}

	for i, f := range got {
		if math.Abs(float64(f-expected[i])) > float64EqualityThreshold {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	}
}

func TestBinaryFloatBackward32(t *testing.T) {
	tr := art.NewFloatBinaryTree[float32, int]()
	expected := []float32{-1, -2, 0, float32(math.NaN()), float32(math.Inf(1)), float32(math.Inf(-1)), 1, 2, math.MaxFloat32, math.SmallestNonzeroFloat32}

	slices.Sort(expected)
	slices.Reverse(expected)
	tr.Insert(1, 1)
	tr.Insert(2, 1)
	tr.Insert(float32(math.NaN()), 1)
	tr.Insert(float32(math.Inf(1)), 1)
	tr.Insert(float32(math.Inf(-1)), 1)
	tr.Insert(-1, 1)
	tr.Insert(-2, 1)
	tr.Insert(math.MaxFloat32, 1)
	tr.Insert(math.SmallestNonzeroFloat32, 1)
	tr.Insert(0, 1)
	tr.Insert(-0, 1)
	tr.Insert(+0, 1)

	var got []float32
	for k, _ := range tr.Backward() {
		got = append(got, k)
	}

	for i, f := range got {
		if math.Abs(float64(f-expected[i])) > float64EqualityThreshold {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	}
}

func TestBinaryFloatAll64(t *testing.T) {
	tr := art.NewFloatBinaryTree[float64, int]()
	expected := []float64{-1, -2, 0, math.NaN(), math.Inf(1), math.Inf(-1), 1, 2, math.MaxFloat64, math.SmallestNonzeroFloat64}

	slices.Sort(expected)
	tr.Insert(1, 1)
	tr.Insert(2, 1)
	tr.Insert(math.NaN(), 1)
	tr.Insert(math.Inf(1), 1)
	tr.Insert(math.Inf(-1), 1)
	tr.Insert(-1, 1)
	tr.Insert(-2, 1)
	tr.Insert(math.MaxFloat64, 1)
	tr.Insert(math.SmallestNonzeroFloat64, 1)
	tr.Insert(0, 1)
	tr.Insert(-0, 1)
	tr.Insert(+0, 1)

	var got []float64
	for k, _ := range tr.All() {
		got = append(got, k)
	}

	for i, f := range got {
		if math.Abs(f-expected[i]) > float64EqualityThreshold {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	}
}

func TestBinaryFloatBackward64(t *testing.T) {
	tr := art.NewFloatBinaryTree[float64, int]()
	expected := []float64{-1, -2, 0, math.NaN(), math.Inf(1), math.Inf(-1), 1, 2, math.MaxFloat64, math.SmallestNonzeroFloat64}

	slices.Sort(expected)
	slices.Reverse(expected)
	tr.Insert(1, 1)
	tr.Insert(2, 1)
	tr.Insert(math.NaN(), 1)
	tr.Insert(math.Inf(1), 1)
	tr.Insert(math.Inf(-1), 1)
	tr.Insert(-1, 1)
	tr.Insert(-2, 1)
	tr.Insert(math.MaxFloat64, 1)
	tr.Insert(math.SmallestNonzeroFloat64, 1)
	tr.Insert(0, 1)
	tr.Insert(-0, 1)
	tr.Insert(+0, 1)

	var got []float64
	for k, _ := range tr.Backward() {
		got = append(got, k)
	}

	for i, f := range got {
		if math.Abs(f-expected[i]) > float64EqualityThreshold {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	}
}
