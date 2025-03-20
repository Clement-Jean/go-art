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

type collationSortedTree[K nodeKey, V any] struct {
	buf  *collate.Buffer
	c    *collate.Collator
	root nodeRef
	end  byte
}

func (t *collationSortedTree[K, V]) setEnd(b byte) {
	t.end = b
}

func (t *collationSortedTree[K, V]) setCollator(c *collate.Collator) {
	t.c = c
}

func NewCollationSortedTree[K nodeKey, V any](opts ...func(Tree[K, V])) Tree[K, V] {
	t := &collationSortedTree[K, V]{
		end: '\x00',
		c:   collate.New(language.Und),
		buf: &collate.Buffer{},
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithCollator[K nodeKey, V any](c *collate.Collator) func(Tree[K, V]) {
	return func(t Tree[K, V]) {
		t.setCollator(c)
	}
}

func prefixMismatchCollate[K nodeKey, V any](n nodeRef, key []byte, depth int) int {
	node := n.node()
	maxCmp := min(int(min(maxPrefixLen, node.prefixLen)), len(key)-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if node.prefix[idx] != key[depth+idx] {
			return idx
		}
	}

	if node.prefixLen > maxPrefixLen {
		leaf := minimum[K, V, *collateLeafNode[K, V]](n)
		leafColKey := unsafe.Slice(leaf.colKey, leaf.colKeyLen)

		maxCmp = min(int(len(leafColKey)), len(key)) - depth
		for ; idx < maxCmp; idx++ {
			realIdx := depth + idx
			if leafColKey[realIdx] != key[realIdx] {
				return idx
			}
		}
	}

	return idx
}

func (t *collationSortedTree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V, *collateLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.String(leaf.key, leaf.len)
		keyStr = strings.Trim(keyStr, string(t.end))
		return K(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *collationSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *collateLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.String(leaf.key, leaf.len)
		keyStr = strings.Trim(keyStr, string(t.end))
		return K(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
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

	insert[K](&t.root, keyStr, colKey, leaf)
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

// Delete deletes a element with the given key.
func (t *collationSortedTree[K, V]) Delete(key K) {
	if t.root.pointer == nil {
		return
	}

	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	keyStr, colKey := bck.Transform(key)

	delete[K, V, *collateLeafNode[K, V]](&t.root, keyStr, colKey)
}

// All returns an iterator over the tree in collation order.
func (t *collationSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all[K, V, *collateLeafNode[K, V]](t.root, t.end)
}

// Backward returns an iterator over the tree in reverse collation order.
func (t *collationSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward[K, V, *collateLeafNode[K, V]](t.root, t.end)
}
