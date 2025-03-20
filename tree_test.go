package art_test

import (
	"bufio"
	"log"
	"os"
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

// Collation Sorted Tree
// Tree Insert Words       10   101,629,383 ns/op   94,990,798 B/op    652,958 allocs/op
// Tree Search Words       12    95,836,104 ns/op   95,067,909 B/op    399,721 allocs/op
// Tree Insert UUIDs       12    90,598,632 ns/op  115,181,482 B/op    305,124 allocs/op
// Tree Search UUIDs       12    86,958,861 ns/op  101,243,912 B/op    200,001 allocs/op

// Alphabetically Sorted Tree (equivalent of go-art and go-adaptive-radix-tree)
// Tree Insert Words       39    29,847,196 ns/op   15,415,973 B/op    508,665 allocs/op
// Tree Search Words       60    19,231,235 ns/op            0 B/op          0 allocs/op
// Tree Insert UUIDs       64    16,935,374 ns/op   14,479,164 B/op    300,664 allocs/op
// Tree Search UUIDs       69    14,798,726 ns/op    9,600,019 B/op    200,000 allocs/op

// go-adaptive-radix-tree
// Tree Insert Words	    9	117,888,698 ns/op	37,942,744 B/op	 1,214,541 allocs/op
// Tree Search Words	   26	 44,555,608 ns/op	         0 B/op	         0 allocs/op
// Tree Insert UUIDs	   18	 59,360,135 ns/op	18,375,723 B/op	   485,057 allocs/op
// Tree Search UUIDs	   54	 21,265,931 ns/op	         0 B/op	         0 allocs/op

// go-art
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
