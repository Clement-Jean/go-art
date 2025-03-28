package art_test

import (
	"bufio"
	"fmt"
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

type testOp = uint8

const (
	insert testOp = iota
	search
	delete
)

func BenchmarkWordsAlpha(b *testing.B) {
	words := loadTestFile("testdata/words.txt")
	sizes := []int{100, 1_000, 10_000, 100_000, len(words)}

	for _, op := range []testOp{insert, search, delete} {
		tree := art.NewAlphaSortedTree[[]byte, []byte]()

		if op == search || op == delete {
			for _, w := range words {
				tree.Insert(w, w)
			}
		}

		switch op {
		case insert:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						tree.Insert(w, w)
						i = (i + 1) % size
					}
				})
			}

		case search:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("search_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						if _, ok := tree.Search(w); !ok {
							b.Fatalf("word %q not found", w)
						}
						i = (i + 1) % size
					}
				})
			}

		case delete:
			// TODO i = (i + 1) % size leads to try deleting already deleted items
			//      and thus the benchmark appears to be faster and faster which is
			//      weird.
			for _, size := range sizes {
				b.Run(fmt.Sprintf("delete_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						tree.Delete(w)
						i = (i + 1) % size
					}
				})
			}
		}
	}
}

func BenchmarkWordsCollate(b *testing.B) {
	words := loadTestFile("testdata/words.txt")
	sizes := []int{100, 1_000, 10_000, 100_000, len(words)}

	for _, op := range []testOp{insert, search, delete} {
		tree := art.NewCollationSortedTree[[]byte, []byte]()

		if op == search || op == delete {
			for _, w := range words {
				tree.Insert(w, w)
			}
		}

		switch op {
		case insert:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						tree.Insert(w, w)
						i = (i + 1) % size
					}
				})
			}

		case search:
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

		case delete:
			// TODO i = (i + 1) % size leads to try deleting already deleted items
			//      and thus the benchmark appears to be faster and faster which is
			//      weird.
			for _, size := range sizes {
				b.Run(fmt.Sprintf("delete_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						tree.Delete(w)
						i = (i + 1) % size
					}
				})
			}
		}
	}
}

func BenchmarkUUIDAlpha(b *testing.B) {
	uuids := loadTestFile("testdata/uuid.txt")
	sizes := []int{100, 1_000, 10_000, 100_000}

	for _, op := range []testOp{insert, search, delete} {
		tree := art.NewAlphaSortedTree[[]byte, []byte]()

		if op == search || op == delete {
			for _, u := range uuids {
				tree.Insert(u, u)
			}
		}

		switch op {
		case insert:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						u := uuids[i]
						tree.Insert(u, u)
						i = (i + 1) % size
					}
				})
			}

		case search:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("search_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						u := uuids[i]
						tree.Search(u)
						i = (i + 1) % size
					}
				})
			}

		case delete:
			// TODO i = (i + 1) % size leads to try deleting already deleted items
			//      and thus the benchmark appears to be faster and faster which is
			//      weird.
			for _, size := range sizes {
				b.Run(fmt.Sprintf("delete_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						u := uuids[i]
						tree.Delete(u)
						i = (i + 1) % size
					}
				})
			}
		}
	}
}

func BenchmarkUUIDCollate(b *testing.B) {
	uuids := loadTestFile("testdata/uuid.txt")
	sizes := []int{100, 1_000, 10_000, 100_000}

	for _, op := range []testOp{insert, search, delete} {
		tree := art.NewCollationSortedTree[[]byte, []byte]()

		if op == search || op == delete {
			for _, u := range uuids {
				tree.Insert(u, u)
			}
		}

		switch op {
		case insert:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						u := uuids[i]
						tree.Insert(u, u)
						i = (i + 1) % size
					}
				})
			}

		case search:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("search_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						u := uuids[i]
						tree.Search(u)
						i = (i + 1) % size
					}
				})
			}

		case delete:
			// TODO i = (i + 1) % size leads to try deleting already deleted items
			//      and thus the benchmark appears to be faster and faster which is
			//      weird.
			for _, size := range sizes {
				b.Run(fmt.Sprintf("delete_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						u := uuids[i]
						tree.Delete(u)
						i = (i + 1) % size
					}
				})
			}
		}
	}
}

func BenchmarkHSKAlpha(b *testing.B) {
	words := loadTestFile("testdata/hsk.txt")
	sizes := []int{100, 1_000, len(words)}

	for _, op := range []testOp{insert, search, delete} {
		tree := art.NewAlphaSortedTree[[]byte, []byte]()

		if op == search || op == delete {
			for _, w := range words {
				tree.Insert(w, w)
			}
		}

		switch op {
		case insert:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						tree.Insert(w, w)
						i = (i + 1) % size
					}
				})
			}

		case search:
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

		case delete:
			// TODO i = (i + 1) % size leads to try deleting already deleted items
			//      and thus the benchmark appears to be faster and faster which is
			//      weird.
			for _, size := range sizes {
				b.Run(fmt.Sprintf("delete_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						tree.Delete(w)
						i = (i + 1) % size
					}
				})
			}
		}
	}
}

func BenchmarkHSKCollate(b *testing.B) {
	words := loadTestFile("testdata/hsk.txt")
	sizes := []int{100, 1_000, len(words)}

	for _, op := range []testOp{insert, search, delete} {
		tree := art.NewCollationSortedTree[[]byte, []byte]()

		if op == search || op == delete {
			for _, w := range words {
				tree.Insert(w, w)
			}
		}

		switch op {
		case insert:
			for _, size := range sizes {
				b.Run(fmt.Sprintf("insert_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						tree.Insert(w, w)
						i = (i + 1) % size
					}
				})
			}

		case search:
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

		case delete:
			// TODO i = (i + 1) % size leads to try deleting already deleted items
			//      and thus the benchmark appears to be faster and faster which is
			//      weird.
			for _, size := range sizes {
				b.Run(fmt.Sprintf("delete_size_%d", size), func(b *testing.B) {
					i := 0

					for b.Loop() {
						w := words[i]
						tree.Delete(w)
						i = (i + 1) % size
					}
				})
			}
		}
	}
}

func FuzzAlphaTreeInsert(f *testing.F) {
	words := loadTestFile("testdata/words.txt")

	for _, word := range words {
		f.Add(word, len(word))
	}

	tree := art.NewAlphaSortedTree[[]byte, int]()
	f.Fuzz(func(t *testing.T, a []byte, b int) {
		tree.Insert(a, b)
	})
}
