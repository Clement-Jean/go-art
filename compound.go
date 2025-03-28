package art

import (
	"iter"
	"unsafe"
)

type compoundLeafNode[K nodeKey, V any] struct {
	value V
	key   *byte
	len   uint32
}

func (n *compoundLeafNode[K, V]) getKey() []byte          { return unsafe.Slice(n.key, n.len) }
func (n *compoundLeafNode[K, V]) getTransformKey() []byte { return unsafe.Slice(n.key, n.len) }
func (n *compoundLeafNode[K, V]) getValue() V             { return n.value }
func (n *compoundLeafNode[K, V]) setValue(val V)          { n.value = val }

type compoundSortedTree[K any, V any] struct {
	bck  BinaryComparableKey[K]
	root nodeRef
	size int
}

func NewCompoundTree[K any, V any](bck BinaryComparableKey[K]) Tree[K, V] {
	return &compoundSortedTree[K, V]{
		bck: bck,
	}
}

func (t *compoundSortedTree[K, V]) restoreKey(l *compoundLeafNode[K, V]) K {
	keyS := l.getKey()[:l.len-1]
	return t.bck.Restore(keyS)
}

func (t *compoundSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, t.restoreKey)
}

func (t *compoundSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, t.restoreKey)
}

func (t *compoundSortedTree[K, V]) BottomK(k uint) iter.Seq2[K, V] {
	return bottomK(t, k)
}

func (t *compoundSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	ok := delete[K, V, *compoundLeafNode[K, V]](&t.root, keyS, keyS)

	if ok {
		t.size--
	}
	return ok
}

func (t *compoundSortedTree[K, V]) Insert(key K, val V) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')
	createFn := func() *compoundLeafNode[K, V] {
		return &compoundLeafNode[K, V]{
			key:   unsafe.SliceData(keyS),
			value: val,
			len:   uint32(len(keyS)),
		}
	}

	if insert[K](&t.root, keyS, keyS, val, createFn) {
		t.size++
	}
}

func (t *compoundSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V, *compoundLeafNode[K, V]](t.root); l != nil {
		return t.restoreKey(l), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *compoundSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V, *compoundLeafNode[K, V]](t.root); l != nil {
		return t.restoreKey(l), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *compoundSortedTree[K, V]) Prefix(key K) iter.Seq2[K, V] { panic("not implemented") }

func (t *compoundSortedTree[K, V]) Range(start, end K) iter.Seq2[K, V] { panic("not implemented") }

func (t *compoundSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	return search[K, V, *compoundLeafNode[K, V]](t.root, keyS, keyS)
}

func (t *compoundSortedTree[K, V]) TopK(k uint) iter.Seq2[K, V] {
	return topK(t, k)
}

func (t *compoundSortedTree[K, V]) Size() int { return t.size }
