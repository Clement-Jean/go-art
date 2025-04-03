package art_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Clement-Jean/go-art"
)

func TestBinarySignedAll(t *testing.T) {
	tr := art.NewSignedBinaryTree[int, int]()
	expected := []int{-1, -9, -100_000, -11}

	slices.Sort(expected)
	tr.Insert(-1, 1)
	tr.Insert(-11, 1)
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
	expected := []int{-1, -9, -100_000, -11}

	slices.Sort(expected)
	slices.Reverse(expected)
	tr.Insert(-1, 1)
	tr.Insert(-11, 1)
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

func TestSignedRange(t *testing.T) {
	tests := []struct {
		name           string
		start, end     int
		keys, expected []int
	}{
		{
			name:  "start<end",
			start: 0, end: 7,
			keys:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:  "start>end",
			start: 7, end: 0,
			keys:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:  "start==end",
			start: 7, end: 7,
			keys:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: []int{7},
		},
		{
			name:  "outside of range",
			start: 16, end: 20,
			keys:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("range-%s", tt.name), func(t *testing.T) {
			tr := art.NewSignedBinaryTree[int, int]()

			for _, key := range tt.keys {
				tr.Insert(key, key)
			}

			var res []int
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
