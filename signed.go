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

func (t *signedSortedTree[K, V]) restoreKey(ptr unsafe.Pointer) (K, V) {
	l := (*signedLeafNode[K, V])(ptr)
	keyS := l.getKey()
	return t.bck.Restore(keyS), l.value
}

func (t *signedSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, t.restoreKey)
}

func (t *signedSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, t.restoreKey)
}

func (t *signedSortedTree[K, V]) BottomK(k uint) iter.Seq2[K, V] {
	return bottomK(t, k)
}

func (t *signedSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')
	ok := delete[K, V, *signedLeafNode[K, V]](&t.root, keyS, keyS)

	if ok {
		t.size--
	}
	return ok
}

func (t *signedSortedTree[K, V]) Insert(key K, val V) {
	_, keyS := t.bck.Transform(key)
	createFn := func() unsafe.Pointer {
		return unsafe.Pointer(&signedLeafNode[K, V]{
			key:   unsafe.SliceData(keyS),
			value: val,
			len:   uint32(len(keyS)),
		})
	}

	if insert[K, V, *signedLeafNode[K, V]](&t.root, keyS, keyS, val, createFn) {
		t.size++
	}
}

func (t *signedSortedTree[K, V]) Maximum() (K, V, bool) {
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

func (t *signedSortedTree[K, V]) Minimum() (K, V, bool) {
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

func (t *signedSortedTree[K, V]) Prefix(key K) iter.Seq2[K, V] { panic("not implemented") }

func (t *signedSortedTree[K, V]) Range(start, end K) iter.Seq2[K, V] { panic("not implemented") }

func (t *signedSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	return search[K, V, *signedLeafNode[K, V]](t.root, keyS, keyS)
}

func (t *signedSortedTree[K, V]) TopK(k uint) iter.Seq2[K, V] {
	return topK(t, k)
}

func (t *signedSortedTree[K, V]) Size() int { return t.size }
