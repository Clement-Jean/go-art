package art

import (
	"testing"
	"unsafe"
)

func TestTree1Insert(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)

	if tr.root.pointer() == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag()) != nodeKindLeaf {
		t.Errorf("expected a tag of %s, got %s", nodeKindLeaf, nodeKind(tr.root.tag()))
	}

	n4 := (*nodeLeaf[string, int])(tr.root.pointer())

	if got := unsafe.String(n4.key, n4.len); got != "hello" {
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

	if tr.root.pointer() == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag()) != nodeKind4 {
		t.Errorf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag()))
	}

	n4 := (*node4)(tr.root.pointer())

	if got := unsafe.String(&n4.prefix[0], n4.prefixLen); got != "hell" {
		t.Errorf("expected prefix 'hell', got %s", got)
	}

	if n4.childrenLen != 2 {
		t.Errorf("expected 2 children, got %d", n4.childrenLen)
	}

	expectedKeys := []string{"hella", "hello"}
	expectedValues := []int{2, 1}
	for i := range int(n4.childrenLen) {
		if nodeKind(n4.children[i].tag()) != nodeKindLeaf {
			t.Errorf("expected a tag of %s, got %s", nodeKindLeaf, nodeKind(n4.children[i].tag()))
		}

		leaf := (*nodeLeaf[string, int])((n4.children[i].pointer()))

		if got := unsafe.String(leaf.key, leaf.len); got != expectedKeys[i] {
			t.Errorf("expected key to be %q, got %q", expectedKeys[i], got)
		}

		if got := leaf.value; got != expectedValues[i] {
			t.Errorf("expected value to be %d, got %d", expectedValues[i], got)
		}
	}
}

func TestTree2Inserts(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("hel", 2)

	if tr.root.pointer() == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag()) != nodeKind4 {
		t.Errorf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag()))
	}

	n4 := (*node4)(tr.root.pointer())

	if got := unsafe.String(&n4.prefix[0], n4.prefixLen); got != "hel" {
		t.Errorf("expected prefix 'hel', got %s", got)
	}

	if n4.childrenLen != 1 {
		t.Errorf("expected 1 children, got %d", n4.childrenLen)
	}

	expectedKeys := []string{"hello"}
	expectedValues := []int{1}
	for i := range int(n4.childrenLen) {
		if nodeKind(n4.children[i].tag()) != nodeKindLeaf {
			t.Errorf("expected a tag of %s, got %s", nodeKindLeaf, nodeKind(n4.children[i].tag()))
		}

		leaf := (*nodeLeaf[string, int])((n4.children[i].pointer()))

		if got := unsafe.String(leaf.key, leaf.len); got != expectedKeys[i] {
			t.Errorf("expected key to be %q, got %q", expectedKeys[i], got)
		}

		if got := leaf.value; got != expectedValues[i] {
			t.Errorf("expected value to be %d, got %d", expectedValues[i], got)
		}
	}
}

func TestInsert3Times(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("hella", 2)
	tr.Insert("hellu", 3)

	if tr.root.pointer() == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag()) != nodeKind4 {
		t.Errorf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag()))
	}

	n4 := (*node4)(tr.root.pointer())

	if got := unsafe.String(&n4.prefix[0], n4.prefixLen); got != "hell" {
		t.Errorf("expected prefix 'hell', got %s", got)
	}

	if n4.childrenLen != 3 {
		t.Errorf("expected 3 children, got %d", n4.childrenLen)
	}

	expectedKeys := []string{"hella", "hellu", "hello"}
	expectedValues := []int{2, 3, 1}
	for i := range int(n4.childrenLen) {
		if nodeKind(n4.children[i].tag()) != nodeKindLeaf {
			t.Errorf("expected a tag of %s, got %s", nodeKindLeaf, nodeKind(n4.children[i].tag()))
		}

		leaf := (*nodeLeaf[string, int])((n4.children[i].pointer()))

		if got := unsafe.String(leaf.key, leaf.len); got != expectedKeys[i] {
			t.Errorf("expected key to be %q, got %q", expectedKeys[i], got)
		}

		if got := leaf.value; got != expectedValues[i] {
			t.Errorf("expected value to be %d, got %d", expectedValues[i], got)
		}
	}
}

func TestInsert2TimesWithPrefix(t *testing.T) {
	var tr Tree[string, int]
	tr.Insert("hello", 1)
	tr.Insert("olleh", 2)

	if tr.root.pointer() == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag()) != nodeKind4 {
		t.Errorf("expected a tag of %s, got %s", nodeKind4, nodeKind(tr.root.tag()))
	}

	n4 := (*node4)(tr.root.pointer())

	if got := unsafe.String(&n4.prefix[0], n4.prefixLen); got != "" {
		t.Errorf("expected prefix '', got %s", got)
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

	if tr.root.pointer() == nil {
		t.Error("expected a root node, got nil")
	}

	if nodeKind(tr.root.tag()) != nodeKind16 {
		t.Errorf("expected a tag of %s, got %s", nodeKind16, nodeKind(tr.root.tag()))
	}

	n16 := (*node16)(tr.root.pointer())

	if got := unsafe.String(&n16.prefix[0], n16.prefixLen); got != "hell" {
		t.Errorf("expected prefix 'hell', got %s", got)
	}

	for i := range int(n16.childrenLen) {
		if nodeKind(n16.children[i].tag()) != nodeKindLeaf {
			t.Errorf("expected a tag of %s, got %s", nodeKindLeaf, nodeKind(n16.children[i].tag()))
		}

		leaf := (*nodeLeaf[string, int])((n16.children[i].pointer()))
		keyStr := unsafe.String(leaf.key, leaf.len)
		val, ok := tr.Search(keyStr)

		if !ok {
			t.Errorf("failed to find %q", keyStr)
		}

		println(val)
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
