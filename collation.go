package art

import (
	"iter"
	"strings"
	"unsafe"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type collateLeafNode[K nodeKey, V any] struct {
	colKey    *byte
	key       *byte
	value     V
	colKeyLen uint32
	len       uint32
}

func (n *collateLeafNode[K, V]) getKey() *byte           { return n.key }
func (n *collateLeafNode[K, V]) getTransformKey() *byte  { return n.colKey }
func (n *collateLeafNode[K, V]) getLen() uint32          { return n.len }
func (n *collateLeafNode[K, V]) getTransformLen() uint32 { return n.colKeyLen }
func (n *collateLeafNode[K, V]) getValue() V             { return n.value }

type collationSortedTree[K chars, V any] struct {
	buf  *collate.Buffer
	c    *collate.Collator
	root nodeRef
	size int
}

func NewCollationSortedTree[K chars, V any](opts ...func(*collationSortedTree[K, V])) Tree[K, V] {
	t := &collationSortedTree[K, V]{
		c:   collate.New(language.Und),
		buf: &collate.Buffer{},
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithCollator[K chars, V any](c *collate.Collator) func(*collationSortedTree[K, V]) {
	return func(t *collationSortedTree[K, V]) {
		t.c = c
	}
}

// All returns an iterator over the tree in collation order.
func (t *collationSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, func(l *collateLeafNode[K, V]) K {
		return K(unsafe.String(l.key, l.len))
	})
}

// Backward returns an iterator over the tree in reverse collation order.
func (t *collationSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *collateLeafNode[K, V]) K {
		return K(unsafe.String(l.key, l.len))
	})
}

// Delete deletes a element with the given key.
func (t *collationSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	keyStr, colKey := bck.Transform(key)

	ok := delete[K, V, *collateLeafNode[K, V]](&t.root, keyStr, colKey)

	if ok {
		t.size--
	}
	return ok
}

// Insert inserts a key-value pair in the tree.
func (t *collationSortedTree[K, V]) Insert(key K, val V) {
	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}

	keyStr, colKey := bck.Transform(key)
	leaf := &collateLeafNode[K, V]{
		colKey:    unsafe.SliceData(colKey),
		key:       unsafe.SliceData(keyStr),
		value:     val,
		colKeyLen: uint32(len(colKey)),
		len:       uint32(len(keyStr)),
	}

	if insert[K](&t.root, keyStr, colKey, leaf) {
		t.size++
	}
}

func (t *collationSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *collateLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.String(leaf.key, leaf.len)
		keyStr = strings.Trim(keyStr, string('\x00'))
		return K(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *collationSortedTree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V, *collateLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.String(leaf.key, leaf.len)
		keyStr = strings.Trim(keyStr, string('\x00'))
		return K(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

// Search searches for element with the given key.
// It returns whether the key is present (bool) and its value if it is present.
func (t *collationSortedTree[K, V]) Search(key K) (V, bool) {
	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	keyStr, colKey := bck.Transform(key)

	return search[K, V, *collateLeafNode[K, V]](t.root, keyStr, colKey)
}

func (t *collationSortedTree[K, V]) Size() int { return t.size }
