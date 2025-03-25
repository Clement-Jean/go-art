package art_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/Clement-Jean/go-art"
)

func loadTestFile(path string) [][]byte {
	var words [][]byte

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("could not open %s", path)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		if line, err := reader.ReadBytes(byte('\n')); err != nil {
			break
		} else if len(line) > 0 {
			words = append(words, line[:len(line)-1])
		}
	}

	return words
}

func TestT(t *testing.T) {
	key := fmt.Appendf(nil, "flydb-key-%09d", 55)
	val := []byte{102, 108, 121, 100, 98, 45, 118, 97, 108, 117, 101, 45, 109, 105, 90, 68, 53, 48, 69, 73, 55, 105, 105, 99, 88, 110, 108, 69, 54, 51, 65, 106, 74, 67, 107, 89, 119, 115, 88, 48, 75, 118, 89, 90, 97, 87, 103, 50, 107, 101, 122, 118, 81, 73, 81, 100, 98, 84, 122, 49, 99, 86, 115, 69, 109, 104, 83, 87, 102, 85, 51, 55, 66, 79, 101, 73, 85, 85, 108, 105, 112, 86, 82, 106, 86, 115, 76, 103, 116, 73, 53, 79, 100, 54, 51, 101, 112, 122, 66, 88, 66, 115, 83, 114, 101, 70, 88, 54, 55, 90, 105, 114, 100, 120, 56, 69, 52, 56, 110, 81, 87, 54, 120, 113, 119, 116, 87, 50, 115, 90, 97, 70, 49, 110, 88, 52, 98, 71, 89, 48}

	tr := art.NewAlphaSortedTree[[]byte, []byte]()

	tr.Insert(key, val)
	got, ok := tr.Search(key)

	if !ok {
		t.Fatal("not found")
	}

	if !slices.Equal(got, val) {
		t.Fatalf("expected %v, got %v", val, got)
	}
}

// Collation Sorted Tree
// Tree Insert Words       12    99,373,198 ns/op   99,781,714 B/op    852,531 allocs/op
// Tree Search Words       13    87,622,356 ns/op   88,598,267 B/op    366,298 allocs/op
// Tree Insert UUIDs       12    87,078,715 ns/op  119,981,510 B/op    405,124 allocs/op
// Tree Search UUIDs       13    82,261,397 ns/op  122,132,361 B/op    300,000 allocs/op

// Alphabetically Sorted Tree (equivalent of go-art and go-adaptive-radix-tree)
// Tree Insert Words       39    29,847,196 ns/op   15,415,973 B/op    508,665 allocs/op
// Tree Search Words       60    19,231,235 ns/op            0 B/op          0 allocs/op
// Tree Insert UUIDs       64    16,935,374 ns/op   14,479,164 B/op    300,664 allocs/op
// Tree Search UUIDs       69    14,798,726 ns/op    9,600,019 B/op    200,000 allocs/op

// go-adaptive-radix-tree (https://github.com/plar/go-adaptive-radix-tree/tree/master)
// Tree Insert Words	    9	117,888,698 ns/op	37,942,744 B/op	 1,214,541 allocs/op
// Tree Search Words	   26	 44,555,608 ns/op	         0 B/op	         0 allocs/op
// Tree Insert UUIDs	   18	 59,360,135 ns/op	18,375,723 B/op	   485,057 allocs/op
// Tree Search UUIDs	   54	 21,265,931 ns/op	         0 B/op	         0 allocs/op

// go-art (https://github.com/kellydunn/go-art)
// Tree Insert Words	    5	272,047,975 ns/op	81,628,987 B/op	 2,547,316 allocs/op
// Tree Search Words	   10	129,011,177 ns/op	13,272,278 B/op	 1,659,033 allocs/op
// Tree Insert UUIDs	   10	140,309,246 ns/op	33,678,160 B/op	   874,561 allocs/op
// Tree Search UUIDs	   20	 82,120,943 ns/op	 3,883,131 B/op	   485,391 allocs/op

func BenchmarkWordsAlphaTreeInsert(b *testing.B) {
	tree := art.NewAlphaSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/words.txt")

	for b.Loop() {
		for _, w := range words {
			tree.Insert(w, w)
		}
	}
}

func BenchmarkWordsCollateTreeInsert(b *testing.B) {
	tree := art.NewCollationSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/words.txt")

	for b.Loop() {
		for _, w := range words {
			tree.Insert(w, w)
		}
	}
}

func BenchmarkWordsAlphaTreeSearch(b *testing.B) {
	tree := art.NewAlphaSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/words.txt")
	for _, w := range words {
		tree.Insert(w, w)
	}

	for b.Loop() {
		for _, w := range words {
			tree.Search(w)
		}
	}
}

func BenchmarkWordsCollateTreeSearch(b *testing.B) {
	tree := art.NewCollationSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/words.txt")
	for _, w := range words {
		tree.Insert(w, w)
	}

	for b.Loop() {
		for _, w := range words {
			tree.Search(w)
		}
	}
}

func BenchmarkUUIDAlphaTreeInsert(b *testing.B) {
	tree := art.NewAlphaSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/uuid.txt")

	for b.Loop() {
		for _, w := range words {
			tree.Insert(w, w)
		}
	}
}

func BenchmarkUUIDCollateTreeInsert(b *testing.B) {
	tree := art.NewCollationSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/uuid.txt")

	for b.Loop() {
		for _, w := range words {
			tree.Insert(w, w)
		}
	}
}

func BenchmarkUUIDAlphaTreeSearch(b *testing.B) {
	tree := art.NewAlphaSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/uuid.txt")
	for _, w := range words {
		tree.Insert(w, w)
	}

	for b.Loop() {
		for _, w := range words {
			tree.Search(w)
		}
	}
}

func BenchmarkUUIDCollateTreeSearch(b *testing.B) {
	tree := art.NewCollationSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/uuid.txt")
	for _, w := range words {
		tree.Insert(w, w)
	}

	for b.Loop() {
		for _, w := range words {
			tree.Search(w)
		}
	}
}
