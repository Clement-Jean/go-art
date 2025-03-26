package main

import (
	"fmt"
	"math"
	"math/bits"

	"github.com/Clement-Jean/go-art"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func alphaTree() {
	tree := art.NewAlphaSortedTree[string, string]()

	tree.Insert("apple", "A sweet red fruit")
	tree.Insert("banana", "A long yellow fruit")
	tree.Insert("cherry", "A small red fruit")
	tree.Insert("date", "A sweet brown fruit")

	if value, found := tree.Search("banana"); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("\nAscending order iteration:")
	for key, value := range tree.All() {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}

	fmt.Println("\nDescending order iteration:")
	for key, value := range tree.Backward() {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
}

func collateTree() {
	col := collate.New(language.English, collate.Numeric)
	tree := art.NewCollationSortedTree(art.WithCollator[string, int](col))

	tree.Insert("11", 1)
	tree.Insert("1", 2)
	tree.Insert("9", 3)
	tree.Insert("100", 4)

	if value, found := tree.Search("11"); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("\nAscending order iteration:")
	for key, value := range tree.All() {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}

	fmt.Println("\nDescending order iteration:")
	for key, value := range tree.Backward() {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}
}

func unsignedTree() {
	tree := art.NewUnsignedBinaryTree[uint, int]()

	tree.Insert(11, 1)
	tree.Insert(1, 2)
	tree.Insert(9, 3)
	tree.Insert(100, 4)

	if value, found := tree.Search(11); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("\nAscending order iteration:")
	for key, value := range tree.All() {
		fmt.Printf("Key: %d, Value: %d\n", key, value)
	}

	fmt.Println("\nDescending order iteration:")
	for key, value := range tree.Backward() {
		fmt.Printf("Key: %d, Value: %d\n", key, value)
	}
}

func signedTree() {
	tree := art.NewSignedBinaryTree[int, int]()

	tree.Insert(-11, 1)
	tree.Insert(1, 2)
	tree.Insert(-9, 3)
	tree.Insert(100, 4)

	if value, found := tree.Search(-11); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("\nAscending order iteration:")
	for key, value := range tree.All() {
		fmt.Printf("Key: %d, Value: %d\n", key, value)
	}

	fmt.Println("\nDescending order iteration:")
	for key, value := range tree.Backward() {
		fmt.Printf("Key: %d, Value: %d\n", key, value)
	}
}

func floatTree() {
	tree := art.NewFloatBinaryTree[float64, int]()

	tree.Insert(math.Inf(-1), 1)
	tree.Insert(math.Inf(1), 2)
	tree.Insert(math.NaN(), 3)
	tree.Insert(0, 4)
	tree.Insert(0.1234, 5)

	if value, found := tree.Search(0.1234); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("\nAscending order iteration:")
	for key, value := range tree.All() {
		fmt.Printf("Key: %f, Value: %d\n", key, value)
	}

	fmt.Println("\nDescending order iteration:")
	for key, value := range tree.Backward() {
		fmt.Printf("Key: %f, Value: %d\n", key, value)
	}
}

type Account struct {
	ID   uint
	name string
}

type AccountKey struct{}

func (ak AccountKey) Transform(a Account) ([]byte, []byte) {
	var (
		ubk art.UnsignedBinaryKey[uint]
		aok art.AlphabeticalOrderKey[string]
		b   []byte
	)

	_, c := ubk.Transform(a.ID)
	_, d := aok.Transform(a.name)
	b = append(b, c...)
	b = append(b, d...)
	return b, b
}
func (ak AccountKey) Restore(b []byte) Account {
	var (
		a   Account
		ubk art.UnsignedBinaryKey[uint]
		aok art.AlphabeticalOrderKey[string]
	)

	id := b[:(bits.UintSize / 8)]
	name := b[(bits.UintSize / 8):]

	a.ID = ubk.Restore(id)
	a.name = aok.Restore(name)
	return a
}

func compoundTree() {
	var key AccountKey
	tree := art.NewCompoundTree[Account, bool](key)

	tree.Insert(Account{ID: 1, name: "Clement"}, true)
	tree.Insert(Account{ID: 3, name: "Elisabeth"}, true)
	tree.Insert(Account{ID: 2, name: "Matt"}, false)

	if value, found := tree.Search(Account{ID: 1, name: "Clement"}); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}

	fmt.Println("\nAscending order iteration:")
	for key, value := range tree.All() {
		fmt.Printf("Key: %v, Value: %t\n", key, value)
	}

	fmt.Println("\nDescending order iteration:")
	for key, value := range tree.Backward() {
		fmt.Printf("Key: %v, Value: %t\n", key, value)
	}
}

func main() {
	alphaTree()
	println()
	collateTree()
	println()
	unsignedTree()
	println()
	signedTree()
	println()
	floatTree()
	println()
	compoundTree()
}
