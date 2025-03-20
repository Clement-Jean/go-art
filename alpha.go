package art

import (
	"iter"
	"unsafe"

	"golang.org/x/text/collate"
)

type alphaLeafNode[K nodeKey, V any] struct {
	key   *byte
	value V
	len   uint32
}

func (n *alphaLeafNode[K, V]) getKey() *byte           { return n.key }
func (n *alphaLeafNode[K, V]) getTransformKey() *byte  { return n.key }
func (n *alphaLeafNode[K, V]) getLen() uint32          { return n.len }
func (n *alphaLeafNode[K, V]) getTransformLen() uint32 { return n.len }
func (n *alphaLeafNode[K, V]) getValue() V             { return n.value }

type alphaSortedTree[K nodeKey, V any] struct {
	root nodeRef
	end  byte
	bck  AlphabeticalOrderKey[K]
}

func (t *alphaSortedTree[K, V]) setEnd(b byte) {
	t.end = b
}

func (t *alphaSortedTree[K, V]) setCollator(c *collate.Collator) {}

func NewAlphaSortedTree[K nodeKey, V any](opts ...func(Tree[K, V])) Tree[K, V] {
	t := &alphaSortedTree[K, V]{
		end: '\x00',
		bck: AlphabeticalOrderKey[K]{},
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithEndByte[K nodeKey, V any](b byte) func(Tree[K, V]) {
	return func(t Tree[K, V]) {
		t.setEnd(b)
	}
}

func (t *alphaSortedTree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V, *alphaLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.Slice(leaf.key, leaf.len)
		keyStr = keyStr[:len(keyStr)-1] //strings.Trim(keyStr, string(t.end))
		return t.bck.Restore(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *alphaLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.Slice(leaf.key, leaf.len)
		keyStr = keyStr[:len(keyStr)-1] //strings.Trim(keyStr, string(t.end))
		return t.bck.Restore(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Insert(key K, val V) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, t.end)
	leaf := &alphaLeafNode[K, V]{
		key:   unsafe.SliceData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}

	insert[K](&t.root, keyStr, keyStr, leaf)
}

func (t *alphaSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, t.end)

	return search[K, V, *alphaLeafNode[K, V]](t.root, keyStr, keyStr)
}

func (t *alphaSortedTree[K, V]) Delete(key K) {
	if t.root.pointer == nil {
		return
	}

	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, t.end)

	delete[K, V, *alphaLeafNode[K, V]](&t.root, keyStr, keyStr)
}

// All returns an iterator over the tree in alphabetical order.
func (t *alphaSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all[K, V, *alphaLeafNode[K, V]](t.root, t.end)
}

// Backward returns an iterator over the tree in reverse alphabetical order.
func (t *alphaSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward[K, V, *alphaLeafNode[K, V]](t.root, t.end)
}
