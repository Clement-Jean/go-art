Adaptive Radix Tree in Go
====

Features:
* Insert/Search/Delete
* Minimum / Maximum value lookups
* Ordered iteration (All)
* Reverse iteration (Backward)
* Support for multiple key types (bytes, floats, unsigned ints, signed ints, collation keys, compound keys)

# Usage

A trivial use case is to order by array of bytes:

```go
package main

import (
	"fmt"
	
	"github.com/Clement-Jean/go-art"
)

func main() {
	// Initialize a new Adaptive Radix Tree
	tree := art.NewAlphaSortedTree[string, string]()

	// Insert key-value pairs into the tree
	tree.Insert("apple", "A sweet red fruit")
	tree.Insert("banana", "A long yellow fruit")
	tree.Insert("cherry", "A small red fruit")
	tree.Insert("date", "A sweet brown fruit")

	// Search for a value by key
	if value, found := tree.Search("banana"); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}

	// Iterate over the tree in ascending order
	fmt.Println("\nAscending order iteration:")
	for key, value := range tree.All() {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}

	// Iterate over the tree in descending order using reverse traversal
	fmt.Println("\nDescending order iteration:")
	for key, value := range tree.Backward() {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
}

// Expected Output:
// Found: A long yellow fruit
//
// Ascending order iteration:
// Key: apple, Value: A sweet red fruit
// Key: banana, Value: A long yellow fruit
// Key: cherry, Value: A small red fruit
// Key: date, Value: A sweet brown fruit
//
// Descending order iteration:
// Key: date, Value: A sweet brown fruit
// Key: cherry, Value: A small red fruit
// Key: banana, Value: A long yellow fruit
// Key: apple, Value: A sweet red fruit
```

But there are more uses cases in the examples folder.

# Performance

go-art outperforms [plar/go-adaptive-radix-tree](https://github.com/plar/go-adaptive-radix-tree) and [kellydunn/go-art](https://github.com/kellydunn/go-art) by taking advantage of SIMD, SWAR and shaving off allocations.

Benchmarks were performed on datasets extracted from different projects:
- The "Words" dataset contains a list of 235,886 english words. [2]
- The "UUIDs" dataset contains 100,000 uuids.                   [2]
- The "HSK Words" dataset contains 4,995 words.                 [5]

To see more benchmarks just run

```
$ go test -run=^$ -bench=. -benchmem -count=10
```

# References

[1] [The Adaptive Radix Tree: ARTful Indexing for Main-Memory Databases (Specification)](http://www-db.in.tum.de/~leis/papers/ART.pdf)

[2] [C99 implementation of the Adaptive Radix Tree](https://github.com/armon/libart)

[3] [go-adaptive-radix-tree](https://github.com/plar/go-adaptive-radix-tree)

[4] [go-art](https://github.com/kellydunn/go-art)

[5] [HSK Words](http://hskhsk.pythonanywhere.com/hskwords). HSK(Hanyu Shuiping Kaoshi) - Standardized test of Standard Mandarin Chinese proficiency.
