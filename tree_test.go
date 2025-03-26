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

func BenchmarkHSKAlphaTreeInsert(b *testing.B) {
	tree := art.NewAlphaSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/hsk.txt")

	for b.Loop() {
		for _, w := range words {
			tree.Insert(w, w)
		}
	}
}

func BenchmarkHSKAlphaTreeSearch(b *testing.B) {
	tree := art.NewAlphaSortedTree[[]byte, []byte]()
	words := loadTestFile("testdata/hsk.txt")
	for _, w := range words {
		tree.Insert(w, w)
	}

	for b.Loop() {
		for _, w := range words {
			tree.Search(w)
		}
	}
}
