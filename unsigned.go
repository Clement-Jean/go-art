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

func (t *unsignedSortedTree[K, V]) restoreKey(ptr unsafe.Pointer) (K, V) {
	l := (*unsignedLeafNode[K, V])(ptr)
	keyS := l.getKey()
	return t.bck.Restore(keyS), l.value
}

func (t *unsignedSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, t.restoreKey)
}

func (t *unsignedSortedTree[K, V]) BottomK(k uint) iter.Seq2[K, V] {
	return bottomK(t, k)
}

func (t *unsignedSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, t.restoreKey)
}

func (t *unsignedSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyS := t.bck.Transform(key)
	ok := delete[K, V, *unsignedLeafNode[K, V]](&t.root, keyS, keyS)

	if ok {
		t.size--
	}
	return ok
}

func (t *unsignedSortedTree[K, V]) Insert(key K, val V) {
	_, keyS := t.bck.Transform(key)
	createFn := func() unsafe.Pointer {
		return unsafe.Pointer(&unsignedLeafNode[K, V]{
			key:   unsafe.SliceData(keyS),
			value: val,
			len:   uint32(len(keyS)),
		})
	}

	if insert[K, V, *unsignedLeafNode[K, V]](&t.root, keyS, keyS, val, createFn) {
		t.size++
	}
}

func (t *unsignedSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V](t.root); l != nil {
		k, v := t.restoreKey(l)
		return k, v, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *unsignedSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V](t.root); l != nil {
		k, v := t.restoreKey(l)
		return k, v, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *unsignedSortedTree[K, V]) Prefix(key K) iter.Seq2[K, V] { panic("not implemented") }

func (t *unsignedSortedTree[K, V]) Range(start, end K) iter.Seq2[K, V] { panic("not implemented") }

func (t *unsignedSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	return search[K, V, *unsignedLeafNode[K, V]](t.root, keyS, keyS)
}

func (t *unsignedSortedTree[K, V]) TopK(k uint) iter.Seq2[K, V] {
	return topK(t, k)
}

func (t *unsignedSortedTree[K, V]) Size() int { return t.size }
