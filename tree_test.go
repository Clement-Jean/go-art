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

// BenchmarkWordsAlphaTreeInsert-10              39          29605098 ns/op
// BenchmarkWordsCollateTreeInsert-10            10         108915104 ns/op
// BenchmarkWordsAlphaTreeSearch-10              61          19196441 ns/op
// BenchmarkWordsCollateTreeSearch-10            61          19225271 ns/op
// BenchmarkUUIDAlphaTreeInsert-10               69          17095396 ns/op
// BenchmarkUUIDCollateTreeInsert-10             12          94671153 ns/op
// BenchmarkUUIDAlphaTreeSearch-10               73          14885850 ns/op
// BenchmarkUUIDCollateTreeSearch-10             13          87197651 ns/op

// BenchmarkWordsAlphaTreeInsert-10              39          29,686,267 ns/op
// BenchmarkWordsCollateTreeInsert-10            10         110,644,838 ns/op
// BenchmarkWordsAlphaTreeSearch-10              61          19,212,364 ns/op
// BenchmarkWordsCollateTreeSearch-10            60          19,238,935 ns/op
// BenchmarkUUIDAlphaTreeInsert-10               68          17,298,688 ns/op
// BenchmarkUUIDCollateTreeInsert-10             12          96,885,066 ns/op
// BenchmarkUUIDAlphaTreeSearch-10               74          14,873,584 ns/op
// BenchmarkUUIDCollateTreeSearch-10             12          91,299,587 ns/op

// Tree Insert Words	9	117,888,698 ns/op	37,942,744 B/op	 1,214,541 allocs/op
// Tree Search Words	26	44,555,608 ns/op	         0 B/op	         0 allocs/op
// Tree Insert UUIDs	18	59,360,135 ns/op	18,375,723 B/op	   485,057 allocs/op
// Tree Search UUIDs	54	21,265,931 ns/op	         0 B/op	         0 allocs/op

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
