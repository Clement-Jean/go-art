package art

import (
	"iter"
	"unsafe"
)

type signedLeafNode[K nodeKey, V any] struct {
	value V
	key   *byte
	len   uint32
}

func (n *signedLeafNode[K, V]) getKey() []byte          { return unsafe.Slice(n.key, n.len) }
func (n *signedLeafNode[K, V]) getTransformKey() []byte { return unsafe.Slice(n.key, n.len) }
func (n *signedLeafNode[K, V]) getValue() V             { return n.value }
func (n *signedLeafNode[K, V]) setValue(val V)          { n.value = val }

type signedSortedTree[K ints, V any] struct {
	root nodeRef
	bck  SignedBinaryKey[K]
	size int
}

func NewSignedBinaryTree[K ints, V any]() Tree[K, V] {
	return &signedSortedTree[K, V]{}
}

func (t *signedSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, func(l *signedLeafNode[K, V]) K {
		keyS := l.getKey()
		return t.bck.Restore(keyS)
	})
}

func (t *signedSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *signedLeafNode[K, V]) K {
		keyS := l.getKey()
		return t.bck.Restore(keyS)
	})
}

func (t *signedSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyS := t.bck.Transform(key)
	ok := delete[K, V, *signedLeafNode[K, V]](&t.root, keyS, keyS)

	if ok {
		t.size--
	}
	return ok
}

func (t *signedSortedTree[K, V]) Insert(key K, val V) {
	_, keyS := t.bck.Transform(key)
	leaf := &signedLeafNode[K, V]{
		key:   unsafe.SliceData(keyS),
		value: val,
		len:   uint32(len(keyS)),
	}

	if insert[K](&t.root, keyS, keyS, leaf) {
		t.size++
	}
}

func (t *signedSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V, *signedLeafNode[K, V]](t.root); l != nil {
		keyS := l.getKey()
		return t.bck.Restore(keyS), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *signedSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V, *signedLeafNode[K, V]](t.root); l != nil {
		keyS := l.getKey()
		return t.bck.Restore(keyS), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *signedSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	return search[K, V, *signedLeafNode[K, V]](t.root, keyS, keyS)
}

func (t *signedSortedTree[K, V]) Size() int { return t.size }
