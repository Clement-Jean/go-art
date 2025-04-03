package art_test

import (
	"fmt"
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

func TestFloatRange(t *testing.T) {
	tests := []struct {
		name           string
		start, end     float64
		keys, expected []float64
	}{
		{
			name:  "start<end",
			start: 0, end: 7,
			keys:     []float64{-0.1234, 1.1234, 2.1234, -3.1234, 4.1234, 5.1234, 6.1234, 7.1234, 8.1234, 9.1234, 10.1234, 11.1234, 12.1234, 13.1234, 14.1234, 15.1234},
			expected: []float64{1.1234, 2.1234, 4.1234, 5.1234, 6.1234},
		},
		{
			name:  "start>end",
			start: 7, end: 0,
			keys:     []float64{-0.1234, 1.1234, 2.1234, -3.1234, 4.1234, 5.1234, 6.1234, 7.1234, 8.1234, 9.1234, 10.1234, 11.1234, 12.1234, 13.1234, 14.1234, 15.1234},
			expected: []float64{1.1234, 2.1234, 4.1234, 5.1234, 6.1234},
		},
		{
			name:  "start==end",
			start: -7.1234, end: -7.1234,
			keys:     []float64{0.1234, 1.1234, 2.1234, 3.1234, 4.1234, 5.1234, 6.1234, -7.1234, 8.1234, 9.1234, 10.1234, 11.1234, 12.1234, 13.1234, 14.1234, 15.1234},
			expected: []float64{-7.1234},
		},
		{
			name:  "outside of range",
			start: 16, end: 20,
			keys:     []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: []float64{},
		},
		{
			name:  "Inf",
			start: math.Inf(-1), end: math.Inf(1),
			keys:     []float64{math.NaN(), 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		// { // Comparing NaN leads to slices not equal
		// 	name:  "NaN",
		// 	start: math.NaN(), end: math.Inf(1),
		// 	keys:     []float64{math.NaN(), 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		// 	expected: []float64{math.NaN(), 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		// },
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("range-%s", tt.name), func(t *testing.T) {
			tr := art.NewFloatBinaryTree[float64, float64]()

			for _, key := range tt.keys {
				tr.Insert(key, key)
			}

			var res []float64
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
