package bench_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	art "github.com/Clement-Jean/go-art"
	oart "github.com/kellydunn/go-art"
	gart "github.com/plar/go-adaptive-radix-tree"
)

var (
	words = loadTestFile("../testdata/words.txt")
	sizes = []int{100, 1_000, 10_000, 100_000, len(words)}
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

func BenchmarkGoARTUpdate(b *testing.B) {
	tree := art.NewAlphaSortedTree[[]byte, []byte]()

	for _, word := range words {
		tree.Insert(word, word)
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("update_size_%d", size), func(b *testing.B) {
			i := 0

			for b.Loop() {
				w := words[i]
				tree.Insert(w, w)
				i = (i + 1) % size
			}
		})
	}
}

func BenchmarkOriginalGoARTUpdate(b *testing.B) {
	tree := oart.NewArtTree()

	for _, word := range words {
		tree.Insert(word, word)
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("update_size_%d", size), func(b *testing.B) {
			i := 0

			for b.Loop() {
				w := words[i]
				tree.Insert(w, w)
				i = (i + 1) % size
			}
		})
	}
}

func BenchmarkGoAdaptiveRadixTreeUpdate(b *testing.B) {
	tree := gart.New()

	for _, word := range words {
		tree.Insert(word, word)
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("update_size_%d", size), func(b *testing.B) {
			i := 0

			for b.Loop() {
				w := words[i]
				tree.Insert(w, w)
				i = (i + 1) % size
			}
		})
	}
}

func BenchmarkGoARTInsert(b *testing.B) {
	for _, size := range sizes {
		b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
			i := 0
			tree := art.NewAlphaSortedTree[[]byte, []byte]()

			for b.Loop() {
				w := words[i]
				tree.Insert(w, w)

				i = (i + 1) % size

				if i == 0 {
					tree = art.NewAlphaSortedTree[[]byte, []byte]()
				}
			}
		})
	}
}

func BenchmarkGoAdaptiveRadixTreeInsert(b *testing.B) {
	for _, size := range sizes {
		b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
			i := 0
			tree := gart.New()

			for b.Loop() {
				w := words[i]
				tree.Insert(w, w)

				i = (i + 1) % size

				if i == 0 {
					tree = gart.New()
				}
			}
		})
	}
}

func BenchmarkOriginalGoARTInsert(b *testing.B) {
	for _, size := range sizes {
		b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
			i := 0
			tree := oart.NewArtTree()

			for b.Loop() {
				w := words[i]
				tree.Insert(w, w)

				i = (i + 1) % size

				if i == 0 {
					tree = oart.NewArtTree()
				}
			}
		})
	}
}

func BenchmarkGoARTSearch(b *testing.B) {
	tree := art.NewAlphaSortedTree[[]byte, []byte]()

	for _, word := range words {
		tree.Insert(word, word)
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("search_size_%d", size), func(b *testing.B) {
			i := 0

			for b.Loop() {
				w := words[i]
				tree.Search(w)
				i = (i + 1) % size
			}
		})
	}
}

func BenchmarkGoAdaptiveRadixTreeSearch(b *testing.B) {
	tree := gart.New()

	for _, word := range words {
		tree.Insert(word, word)
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("search_size_%d", size), func(b *testing.B) {
			i := 0

			for b.Loop() {
				w := words[i]
				tree.Search(w)
				i = (i + 1) % size
			}
		})
	}
}

func BenchmarkOriginalGoARTSearch(b *testing.B) {
	tree := oart.NewArtTree()

	for _, word := range words {
		tree.Insert(word, word)
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("search_size_%d", size), func(b *testing.B) {
			i := 0

			for b.Loop() {
				w := words[i]
				tree.Search(w)
				i = (i + 1) % size
			}
		})
	}
}
