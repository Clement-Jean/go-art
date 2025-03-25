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

func (n *alphaLeafNode[K, V]) getKey() []byte          { return unsafe.Slice(n.key, n.len) }
func (n *alphaLeafNode[K, V]) getTransformKey() []byte { return unsafe.Slice(n.key, n.len) }
func (n *alphaLeafNode[K, V]) getValue() V             { return n.value }
func (n *alphaLeafNode[K, V]) setValue(val V)          { n.value = val }

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
		keyS := l.getKey()
		keyS = keyS[:len(keyS)-1] // drop end byte
		return t.bck.Restore(keyS)
	})
}

// Backward returns an iterator over the tree in reverse alphabetical order.
func (t *alphaSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *alphaLeafNode[K, V]) K {
		keyS := l.getKey()
		keyS = keyS[:len(keyS)-1] // drop end byte
		return t.bck.Restore(keyS)
	})
}

func (t *alphaSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	ok := delete[K, V, *alphaLeafNode[K, V]](&t.root, keyS, keyS)

	if ok {
		t.size--
	}
	return ok
}

func (t *alphaSortedTree[K, V]) Insert(key K, val V) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')
	leaf := &alphaLeafNode[K, V]{
		key:   unsafe.SliceData(keyS),
		value: val,
		len:   uint32(len(keyS)),
	}

	if insert[K](&t.root, keyS, keyS, leaf) {
		t.size++
	}
}

func (t *alphaSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V, *alphaLeafNode[K, V]](t.root); l != nil {
		keyS := l.getKey()
		keyS = keyS[:len(keyS)-1]
		return t.bck.Restore(keyS), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V, *alphaLeafNode[K, V]](t.root); l != nil {
		keyS := l.getKey()
		keyS = keyS[:len(keyS)-1]
		return t.bck.Restore(keyS), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	return search[K, V, *alphaLeafNode[K, V]](t.root, keyS, keyS)
}

func (t *alphaSortedTree[K, V]) Size() int { return t.size }
