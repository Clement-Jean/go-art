package art

import (
	"bytes"
	"iter"
	"strings"
	"unsafe"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type collateLeafNode[K nodeKey, V any] struct {
	value     V
	colKey    *byte
	key       *byte
	colKeyLen uint32
	keyLen    uint32
}

func (n *collateLeafNode[K, V]) getKey() []byte          { return unsafe.Slice(n.key, n.keyLen) }
func (n *collateLeafNode[K, V]) getTransformKey() []byte { return unsafe.Slice(n.colKey, n.colKeyLen) }
func (n *collateLeafNode[K, V]) getValue() V             { return n.value }
func (n *collateLeafNode[K, V]) setValue(val V)          { n.value = val }

type collationSortedTree[K chars | []rune, V any] struct {
	buf  *collate.Buffer
	c    *collate.Collator
	root nodeRef
	size int
}

func NewCollationSortedTree[K chars | []rune, V any](opts ...func(*collationSortedTree[K, V])) Tree[K, V] {
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

func (t *collationSortedTree[K, V]) restoreKey(l *collateLeafNode[K, V]) K {
	return K(string(l.getKey()))
}

// All returns an iterator over the tree in collation order.
func (t *collationSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, t.restoreKey)
}

// Backward returns an iterator over the tree in reverse collation order.
func (t *collationSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, t.restoreKey)
}

func (t *collationSortedTree[K, V]) BottomK(k uint) iter.Seq2[K, V] {
	return bottomK(t, k)
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
	keyS, colKey := bck.Transform(key)

	ok := delete[K, V, *collateLeafNode[K, V]](&t.root, keyS, colKey)

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

	keyS, colKey := bck.Transform(key)

	createFn := func() *collateLeafNode[K, V] {
		return &collateLeafNode[K, V]{
			colKey:    unsafe.SliceData(colKey),
			key:       unsafe.SliceData(keyS),
			value:     val,
			keyLen:    uint32(len(keyS)),
			colKeyLen: uint32(len(colKey)),
		}
	}

	if insert[K](&t.root, keyS, colKey, val, createFn) {
		t.size++
	}
}

func (t *collationSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V, *collateLeafNode[K, V]](t.root); l != nil {
		return t.restoreKey(l), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *collationSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V, *collateLeafNode[K, V]](t.root); l != nil {
		return t.restoreKey(l), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *collationSortedTree[K, V]) Prefix(p K) iter.Seq2[K, V] {
	if len(p) == 0 {
		return t.All()
	}

	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	keyS, colKey := bck.Transform(p)

	root := t.root
	if t.root.pointer != nil {
		root = lowestCommonParent[K, V, *collateLeafNode[K, V]](root, colKey)
	}

	hasPrefix := func(k K, v V) bool {
		leafKeyS := []byte(string(k))
		return bytes.HasPrefix(leafKeyS, keyS)
	}
	return filter(root, hasPrefix, t.restoreKey)
}

func (t *collationSortedTree[K, V]) Range(start, end K) iter.Seq2[K, V] {
	if len(end) == 0 {
		end = K(string(maximum[K, V, *collateLeafNode[K, V]](t.root).getKey()))
	}

	if strings.Compare(string(start), string(end)) > 0 { // start > end
		// IDEA: maybe do the iteration in reverse instead?
		start, end = end, start
	}

	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	startKey, startColKey := bck.Transform(start)
	endKey, endColKey := bck.Transform(end)

	return rangeScan(t.root, startKey, endKey, startColKey, endColKey, t.restoreKey)
}

// Search searches for element with the given key.
// It returns whether the key is present (bool) and its value if it is present.
func (t *collationSortedTree[K, V]) Search(key K) (V, bool) {
	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	keyS, colKey := bck.Transform(key)

	return search[K, V, *collateLeafNode[K, V]](t.root, keyS, colKey)
}

func (t *collationSortedTree[K, V]) TopK(k uint) iter.Seq2[K, V] {
	return topK(t, k)
}

func (t *collationSortedTree[K, V]) Size() int { return t.size }
