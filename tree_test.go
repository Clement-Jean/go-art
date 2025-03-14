package art

import (
	"bufio"
	"bytes"
	"os"
	"testing"
	"unsafe"
)

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		key, other string
		depth      int
		expected   int
	}{
		{"A", "a", 0, 0},
		{"saab", "sad", 0, 2},
		{"saab", "sad", 1, 1},
		{"saab", "sad", 2, 0},
		{"saab", "saab", 0, 4},
	}

	for _, test := range tests {
		if lcp := longestCommonPrefix(test.key, test.other, test.depth); lcp != test.expected {
			t.Errorf("expected %d, got %d (%q, %q)", test.expected, lcp, test.key, test.other)
		}
	}
}

func TestTreeInsertVeryLong(t *testing.T) {
	var tr Tree[[]byte, int]
	key1 := []byte{16, 0, 0, 0, 7, 10, 0, 0, 0, 2, 17, 10, 0, 0, 0, 120, 10, 0, 0, 0, 120, 10, 0,
		0, 0, 216, 10, 0, 0, 0, 202, 10, 0, 0, 0, 194, 10, 0, 0, 0, 224, 10, 0, 0, 0,
		230, 10, 0, 0, 0, 210, 10, 0, 0, 0, 206, 10, 0, 0, 0, 208, 10, 0, 0, 0, 232,
		10, 0, 0, 0, 124, 10, 0, 0, 0, 124, 2, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213,
		251, 173, 202, 211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0,
		8, 18, 182, 92, 236, 147, 171, 101, 150, 195, 112, 185, 218, 108, 246,
		139, 164, 234, 195, 58, 177, 0, 8, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213,
		251, 173, 202, 211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0,
		8, 18, 180, 93, 46, 151, 9, 212, 190, 95, 102, 178, 217, 44, 178, 235,
		29, 190, 218, 8, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213, 251, 173, 202,
		211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0, 8, 18, 180, 93,
		46, 151, 9, 212, 190, 95, 102, 183, 219, 229, 214, 59, 125, 182, 71,
		108, 180, 220, 238, 150, 91, 117, 150, 201, 84, 183, 128, 8, 16, 0, 0,
		0, 2, 12, 185, 89, 44, 213, 251, 173, 202, 211, 95, 185, 89, 110, 118,
		251, 173, 202, 199, 101, 0, 8, 18, 180, 93, 46, 151, 9, 212, 190, 95,
		108, 176, 217, 47, 50, 219, 61, 134, 207, 97, 151, 88, 237, 246, 208,
		8, 18, 255, 255, 255, 219, 191, 198, 134, 5, 223, 212, 72, 44, 208,
		250, 180, 14, 1, 0, 0, 8}
	key2 := []byte{16, 0, 0, 0, 7, 10, 0, 0, 0, 2, 17, 10, 0, 0, 0, 120, 10, 0, 0, 0, 120, 10, 0,
		0, 0, 216, 10, 0, 0, 0, 202, 10, 0, 0, 0, 194, 10, 0, 0, 0, 224, 10, 0, 0, 0,
		230, 10, 0, 0, 0, 210, 10, 0, 0, 0, 206, 10, 0, 0, 0, 208, 10, 0, 0, 0, 232,
		10, 0, 0, 0, 124, 10, 0, 0, 0, 124, 2, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213,
		251, 173, 202, 211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0,
		8, 18, 182, 92, 236, 147, 171, 101, 150, 195, 112, 185, 218, 108, 246,
		139, 164, 234, 195, 58, 177, 0, 8, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213,
		251, 173, 202, 211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0,
		8, 18, 180, 93, 46, 151, 9, 212, 190, 95, 102, 178, 217, 44, 178, 235,
		29, 190, 218, 8, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213, 251, 173, 202,
		211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0, 8, 18, 180, 93,
		46, 151, 9, 212, 190, 95, 102, 183, 219, 229, 214, 59, 125, 182, 71,
		108, 180, 220, 238, 150, 91, 117, 150, 201, 84, 183, 128, 8, 16, 0, 0,
		0, 3, 12, 185, 89, 44, 213, 251, 133, 178, 195, 105, 183, 87, 237, 150,
		155, 165, 150, 229, 97, 182, 0, 8, 18, 161, 91, 239, 50, 10, 61, 150,
		223, 114, 179, 217, 64, 8, 12, 186, 219, 172, 150, 91, 53, 166, 221,
		101, 178, 0, 8, 18, 255, 255, 255, 219, 191, 198, 134, 5, 208, 212, 72,
		44, 208, 250, 180, 14, 1, 0, 0, 8}

	tr.Insert(key1, 299)
	tr.Insert(key2, 302)
}

func TestTreeInsert(t *testing.T) {
	var tr Tree[string, int]

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tr.Insert(line, len(line))
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}
}

func TestTreeInsertSearch(t *testing.T) {
	var tr Tree[string, int]

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tr.Insert(line, len(line))
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("error seeking file: %s", err)
	}

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if _, ok := tr.Search(line); !ok {
			t.Fatalf("word %s not found", line)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}
}

func TestTree1Insert(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)

	if tr.root.pointer == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag) != nodeKindLeaf {
		t.Errorf("expected a tag of %s, got %s", nodeKindLeaf, nodeKind(tr.root.tag))
	}

	n4 := (*nodeLeaf[string, int])(tr.root.pointer)

	if got := unsafe.String(n4.key, n4.len); got != "hello\000" {
		t.Errorf("expected key to be 'hello', got %q", got)
	}

	if got := n4.value; got != 1 {
		t.Errorf("expected value to be 1, got %d", got)
	}

	if _, ok := tr.Search("hello"); !ok {
		t.Fatalf("didn't find hello after insert")
	}
}

func TestTree2InsertsSameLen(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("hella", 2)

	if tr.root.pointer == nil {
		t.Fatal("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag) != nodeKind4 {
		t.Fatalf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag))
	}

	n4 := (*node4)(tr.root.pointer)

	if got := unsafe.String(&n4.prefix[0], n4.prefixLen); got != "hell" {
		t.Fatalf("expected prefix 'hell', got %q", got)
	}

	if n4.childrenLen != 2 {
		t.Fatalf("expected 2 children, got %d", n4.childrenLen)
	}

	expectedKeys := []string{"hella\000", "hello\000"}
	expectedValues := []int{2, 1}
	for i := range int(n4.childrenLen) {
		if kind := nodeKind(n4.children[i].tag); kind != nodeKindLeaf {
			t.Fatalf("expected a tag of %s, got %s", nodeKindLeaf, kind)
		}

		leaf := (*nodeLeaf[string, int])((n4.children[i].pointer))

		if got := unsafe.String(leaf.key, leaf.len); got != expectedKeys[i] {
			t.Fatalf("expected key to be %q, got %q", expectedKeys[i], got)
		}

		if got := leaf.value; got != expectedValues[i] {
			t.Fatalf("expected value to be %d, got %d", expectedValues[i], got)
		}
	}
}

func TestTree2Inserts(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("hel", 2)

	if tr.root.pointer == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag) != nodeKind4 {
		t.Fatalf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag))
	}

	n4 := (*node4)(tr.root.pointer)

	if got := unsafe.String(&n4.prefix[0], n4.prefixLen); got != "hel" {
		t.Fatalf("expected prefix 'hel', got %q", got)
	}

	if n4.childrenLen != 2 {
		t.Fatalf("expected 1 children, got %d", n4.childrenLen)
	}

	expectedKeys := []string{"hel\000", "hello\000"}
	expectedValues := []int{2, 1}
	for i := range int(n4.childrenLen) {
		if kind := nodeKind(n4.children[i].tag); kind != nodeKindLeaf {
			t.Fatalf("expected a tag of %s, got %s", nodeKindLeaf, kind)
		}

		leaf := (*nodeLeaf[string, int])((n4.children[i].pointer))

		if got := unsafe.String(leaf.key, leaf.len); got != expectedKeys[i] {
			t.Fatalf("expected key to be %q, got %q", expectedKeys[i], got)
		}

		if got := leaf.value; got != expectedValues[i] {
			t.Fatalf("expected value to be %d, got %d", expectedValues[i], got)
		}
	}
}

func TestInsert3Times(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("hella", 2)
	tr.Insert("hellu", 3)

	if tr.root.pointer == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag) != nodeKind4 {
		t.Fatalf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag))
	}

	n4 := (*node4)(tr.root.pointer)

	if got := unsafe.String(&n4.prefix[0], n4.prefixLen); got != "hell" {
		t.Fatalf("expected prefix 'hell', got %q", got)
	}

	if n4.childrenLen != 3 {
		t.Fatalf("expected 3 children, got %d", n4.childrenLen)
	}

	expectedKeys := []string{"hella\000", "hello\000", "hellu\000"}
	expectedValues := []int{2, 1, 3}
	for i := range int(n4.childrenLen) {
		if kind := nodeKind(n4.children[i].tag); kind != nodeKindLeaf {
			t.Fatalf("expected a tag of %s, got %s", nodeKindLeaf, kind)
		}

		leaf := (*nodeLeaf[string, int])((n4.children[i].pointer))

		if got := unsafe.String(leaf.key, leaf.len); got != expectedKeys[i] {
			t.Fatalf("expected key to be %q, got %q", expectedKeys[i], got)
		}

		if got := leaf.value; got != expectedValues[i] {
			t.Fatalf("expected value to be %d, got %d", expectedValues[i], got)
		}
	}
}

func TestInsert2TimesWithPrefix(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("olleh", 2)

	if tr.root.pointer == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag) != nodeKind4 {
		t.Fatalf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag))
	}

	n4 := (*node4)(tr.root.pointer)

	if got := unsafe.String(&n4.prefix[0], n4.prefixLen); got != "" {
		t.Fatalf("expected prefix '', got %s", got)
	}

	println(n4.childrenLen)
}

func TestInsert5TimesWithPrefix(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("helle", 2)
	tr.Insert("hellu", 3)
	tr.Insert("hella", 4)
	tr.Insert("helli", 5)

	if tr.root.pointer == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag) != nodeKind16 {
		t.Fatalf("expected a tag of %s, got %s", nodeKind16, nodeKind(tr.root.tag))
	}

	n16 := (*node16)(tr.root.pointer)

	if got := unsafe.String(&n16.prefix[0], n16.prefixLen); got != "hell" {
		t.Fatalf("expected prefix 'hell', got %s", got)
	}

	for i := range int(n16.childrenLen) {
		if kind := nodeKind(n16.children[i].tag); kind != nodeKindLeaf {
			t.Fatalf("expected a tag of %s, got %s", nodeKindLeaf, kind)
		}

		leaf := (*nodeLeaf[string, int])((n16.children[i].pointer))
		keyStr := unsafe.String(leaf.key, leaf.len-1)
		val, ok := tr.Search(keyStr)

		if !ok {
			t.Errorf("failed to find %q", keyStr)
		}

		println(keyStr, val)
	}
}

func TestInsert17TimesWithPrefix(t *testing.T) {
	var tr Tree[[]byte, int]

	key := []byte{1, 2, 3}

	for i := 1; i <= 17; i++ {
		tr.Insert(append(key, byte(i)), i)
	}

	if tr.root.pointer == nil {
		t.Error("expected a root node, got nil")
	}

	if tr.root.tag != nodeKind48 {
		t.Fatalf("expected a tag of %s, got %s", nodeKind48, tr.root.tag)
	}

	n48 := (*node48)(tr.root.pointer)

	if got := unsafe.Slice(&n48.prefix[0], n48.prefixLen); bytes.Compare(got, key) != 0 {
		t.Fatalf("expected prefix 'hell', got %s", got)
	}

	if n48.childrenLen != 17 {
		t.Fatalf("expected 17 children, got %d", n48.childrenLen)
	}

	for i := range n48.childrenLen {
		leaf := (*nodeLeaf[[]byte, int])(n48.children[i].pointer)

		if kind := n48.children[i].tag; kind != nodeKindLeaf {
			t.Fatalf("expected a tag of %s, got %s", nodeKindLeaf, kind)
		}

		gotKey := unsafe.Slice(leaf.key, leaf.len)

		if bytes.Compare(gotKey, append(append(key, byte(i+1)), 0)) != 0 {
			t.Fatalf("expected %v, got %v", append(key, byte(i+1)), gotKey)
		}

		if leaf.value != int(i+1) {
			t.Fatalf("expected value %d, got %d", i+1, leaf.value)
		}
	}
}

func TestInsert49TimesWithPrefix(t *testing.T) {
	var tr Tree[[]byte, int]

	key := []byte{1, 2, 3}

	for i := range 49 {
		tr.Insert(append(key, byte(i)), i)
	}

	if tr.root.pointer == nil {
		t.Error("expected a root node, got nil")
	}

	if tr.root.tag != nodeKind256 {
		t.Fatalf("expected a tag of %s, got %s", nodeKind256, tr.root.tag)
	}

	n256 := (*node256)(tr.root.pointer)

	if got := unsafe.Slice(&n256.prefix[0], n256.prefixLen); bytes.Compare(got, key) != 0 {
		t.Fatalf("expected prefix 'hell', got %s", got)
	}

	if n256.childrenLen != 49 {
		t.Fatalf("expected 49 children, got %d", n256.childrenLen)
	}

	for i := range n256.childrenLen {
		leaf := (*nodeLeaf[[]byte, int])(n256.children[i].pointer)

		if kind := n256.children[i].tag; kind != nodeKindLeaf {
			t.Fatalf("expected a tag of %s, got %s", nodeKindLeaf, kind)
		}

		gotKey := unsafe.Slice(leaf.key, leaf.len)

		if bytes.Compare(gotKey, append(append(key, byte(i)), 0)) != 0 {
			t.Fatalf("expected %v, got %v", append(key, byte(i)), gotKey)
		}

		if leaf.value != int(i) {
			t.Fatalf("expected value %d, got %d", i, leaf.value)
		}
	}
}

func TestRecursiveInsert(t *testing.T) {
	var tr Tree[string, int]

	tr.Insert("abacate", 1)
	tr.Insert("abacinate", 2)
	tr.Insert("abacination", 3)
	tr.Insert("abacinations", 4)

	if tr.root.pointer == nil {
		t.Fatal("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag) != nodeKind4 {
		t.Fatalf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag))
	}

	n4 := (*node4)(tr.root.pointer)

	for i := range n4.childrenLen {
		kind := nodeKind(n4.children[i].tag)

		switch kind {
		case nodeKindLeaf:
			leaf := (*nodeLeaf[string, int])(n4.children[i].pointer)
			println("1.", unsafe.String(leaf.key, leaf.len))

		case nodeKind4:
			n4 = (*node4)(n4.children[i].pointer)

			println("2.", n4.childrenLen)
			for i := range n4.childrenLen {
				kind = nodeKind(n4.children[i].tag)

				switch kind {
				case nodeKindLeaf:
					leaf := (*nodeLeaf[string, int])(n4.children[i].pointer)
					println("3.", unsafe.String(leaf.key, leaf.len))

				case nodeKind4:
					n4 = (*node4)(n4.children[i].pointer)
					println("4.", n4.childrenLen)
				}
			}
		}
	}
}

func BenchmarkSearch16(b *testing.B) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("helle", 2)
	tr.Insert("hellu", 3)
	tr.Insert("hella", 4)
	tr.Insert("helli", 5)

	for b.Loop() {
		tr.Search("hellu")
	}
}
