package art

import (
	"iter"
	"unsafe"
)

type floatLeafNode[K nodeKey, V any] struct {
	value V
	key   *byte
	len   uint32
}

func (n *floatLeafNode[K, V]) getKey() []byte          { return unsafe.Slice(n.key, n.len) }
func (n *floatLeafNode[K, V]) getTransformKey() []byte { return unsafe.Slice(n.key, n.len) }
func (n *floatLeafNode[K, V]) getValue() V             { return n.value }
func (n *floatLeafNode[K, V]) setValue(val V)          { n.value = val }

type floatSortedTree[K floats, V any] struct {
	root nodeRef
	bck  FloatBinaryKey[K]
	size int
}

func NewFloatBinaryTree[K floats, V any]() Tree[K, V] {
	return &floatSortedTree[K, V]{}
}

func (t *floatSortedTree[K, V]) restoreKey(l *floatLeafNode[K, V]) K {
	keyS := l.getKey()
	return t.bck.Restore(keyS)
}

func (t *floatSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, t.restoreKey)
}

func (t *floatSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, t.restoreKey)
}

func (t *floatSortedTree[K, V]) BottomK(k uint) iter.Seq2[K, V] {
	return bottomK(t, k)
}

func (t *floatSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyS := t.bck.Transform(key)
	ok := delete[K, V, *floatLeafNode[K, V]](&t.root, keyS, keyS)

	if ok {
		t.size--
	}
	return ok
}

func (t *floatSortedTree[K, V]) Insert(key K, val V) {
	_, keyS := t.bck.Transform(key)
	createFn := func() *floatLeafNode[K, V] {
		return &floatLeafNode[K, V]{
			key:   unsafe.SliceData(keyS),
			value: val,
			len:   uint32(len(keyS)),
		}
	}

	if insert[K](&t.root, keyS, keyS, val, createFn) {
		t.size++
	}
}

func (t *floatSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V, *floatLeafNode[K, V]](t.root); l != nil {
		return t.restoreKey(l), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *floatSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V, *floatLeafNode[K, V]](t.root); l != nil {
		return t.restoreKey(l), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *floatSortedTree[K, V]) Prefix(key K) iter.Seq2[K, V] { panic("not implemented") }

func (t *floatSortedTree[K, V]) Range(start, end K) iter.Seq2[K, V] { panic("not implemented") }

func (t *floatSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	return search[K, V, *floatLeafNode[K, V]](t.root, keyS, keyS)
}

func (t *floatSortedTree[K, V]) TopK(k uint) iter.Seq2[K, V] {
	return topK(t, k)
}

func (t *floatSortedTree[K, V]) Size() int { return t.size }
