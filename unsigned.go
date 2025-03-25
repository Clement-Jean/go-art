package art

import (
	"iter"
	"unsafe"
)

type unsignedLeafNode[K nodeKey, V any] struct {
	value V
	key   *byte
	len   uint32
}

func (n *unsignedLeafNode[K, V]) getKey() []byte          { return unsafe.Slice(n.key, n.len) }
func (n *unsignedLeafNode[K, V]) getTransformKey() []byte { return unsafe.Slice(n.key, n.len) }
func (n *unsignedLeafNode[K, V]) getValue() V             { return n.value }
func (n *unsignedLeafNode[K, V]) setValue(val V)          { n.value = val }

type unsignedSortedTree[K uints, V any] struct {
	root nodeRef
	bck  UnsignedBinaryKey[K]
	size int
}

func NewUnsignedBinaryTree[K uints, V any]() Tree[K, V] {
	return &unsignedSortedTree[K, V]{}
}

func (t *unsignedSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, func(l *unsignedLeafNode[K, V]) K {
		keyS := l.getKey()[:l.len-1] // drop end byte
		return t.bck.Restore(keyS)
	})
}

func (t *unsignedSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *unsignedLeafNode[K, V]) K {
		keyS := l.getKey()[:l.len-1] // drop end byte
		return t.bck.Restore(keyS)
	})
}

func (t *unsignedSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	ok := delete[K, V, *unsignedLeafNode[K, V]](&t.root, keyS, keyS)

	if ok {
		t.size--
	}
	return ok
}

func (t *unsignedSortedTree[K, V]) Insert(key K, val V) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')
	leaf := &unsignedLeafNode[K, V]{
		key:   unsafe.SliceData(keyS),
		value: val,
		len:   uint32(len(keyS)),
	}

	if insert[K](&t.root, keyS, keyS, leaf) {
		t.size++
	}
}

func (t *unsignedSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V, *unsignedLeafNode[K, V]](t.root); l != nil {
		keyS := l.getKey()[:l.len-1] // drop end byte
		return t.bck.Restore(keyS), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *unsignedSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V, *unsignedLeafNode[K, V]](t.root); l != nil {
		keyS := l.getKey()[:l.len-1] // drop end byte
		return t.bck.Restore(keyS), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *unsignedSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	return search[K, V, *unsignedLeafNode[K, V]](t.root, keyS, keyS)
}

func (t *unsignedSortedTree[K, V]) Size() int { return t.size }
