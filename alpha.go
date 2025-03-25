package art

import (
	"iter"
	"unsafe"
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

type alphaSortedTree[K chars, V any] struct {
	root nodeRef
	bck  AlphabeticalOrderKey[K]
	size int
}

func NewAlphaSortedTree[K chars, V any]() Tree[K, V] {
	return &alphaSortedTree[K, V]{}
}

// All returns an iterator over the tree in alphabetical order.
func (t *alphaSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, func(l *alphaLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

// Backward returns an iterator over the tree in reverse alphabetical order.
func (t *alphaSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *alphaLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *alphaSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')

	ok := delete[K, V, *alphaLeafNode[K, V]](&t.root, keyStr, keyStr)

	if ok {
		t.size--
	}
	return ok
}

func (t *alphaSortedTree[K, V]) Insert(key K, val V) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')
	leaf := &alphaLeafNode[K, V]{
		key:   unsafe.SliceData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}

	if insert[K](&t.root, keyStr, keyStr, leaf) {
		t.size++
	}
}

func (t *alphaSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *alphaLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.Slice(leaf.key, leaf.len)
		keyStr = keyStr[:len(keyStr)-1]
		return t.bck.Restore(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V, *alphaLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.Slice(leaf.key, leaf.len)
		keyStr = keyStr[:len(keyStr)-1]
		return t.bck.Restore(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')

	return search[K, V, *alphaLeafNode[K, V]](t.root, keyStr, keyStr)
}

func (t *alphaSortedTree[K, V]) Size() int { return t.size }
