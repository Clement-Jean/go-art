package art

import (
	"bytes"
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

func (t *alphaSortedTree[K, V]) restoreKey(l *alphaLeafNode[K, V]) K {
	keyS := l.getKey()
	keyS = keyS[:len(keyS)-1] // drop end byte
	return t.bck.Restore(keyS)
}

// All returns an iterator over the tree in alphabetical order.
func (t *alphaSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, t.restoreKey)
}

// Backward returns an iterator over the tree in reverse alphabetical order.
func (t *alphaSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, t.restoreKey)
}

func (t *alphaSortedTree[K, V]) BottomK(k uint) iter.Seq2[K, V] {
	return bottomK(t, k)
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

	createFn := func() *alphaLeafNode[K, V] {
		return &alphaLeafNode[K, V]{
			key:   unsafe.SliceData(keyS),
			value: val,
			len:   uint32(len(keyS)),
		}
	}

	if insert[K](&t.root, keyS, keyS, val, createFn) {
		t.size++
	}
}

func (t *alphaSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V, *alphaLeafNode[K, V]](t.root); l != nil {
		return t.restoreKey(l), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V, *alphaLeafNode[K, V]](t.root); l != nil {
		return t.restoreKey(l), l.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Prefix(p K) iter.Seq2[K, V] {
	if len(p) == 0 {
		return t.All()
	}

	root := t.root
	if t.root.pointer != nil {
		root = lowestCommonParent[K, V, *alphaLeafNode[K, V]](root, []byte(p))
	}

	hasPrefix := func(k K, v V) bool { return bytes.HasPrefix([]byte(k), []byte(p)) }
	return filter(root, hasPrefix, t.restoreKey)
}

func (t *alphaSortedTree[K, V]) Range(start, end K) iter.Seq2[K, V] {
	if len(end) == 0 {
		end = K(maximum[K, V, *alphaLeafNode[K, V]](t.root).getKey())
	}

	if bytes.Compare([]byte(start), []byte(end)) > 0 { // start > end
		// IDEA: maybe do the iteration in reverse instead?
		start, end = end, start
	}

	startKey := append([]byte(start), '\x00')
	endKey := append([]byte(end), '\x00')
	return rangeScan(t.root, startKey, endKey, startKey, endKey, t.restoreKey)
}

func (t *alphaSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	return search[K, V, *alphaLeafNode[K, V]](t.root, keyS, keyS)
}

func (t *alphaSortedTree[K, V]) TopK(k uint) iter.Seq2[K, V] {
	return topK(t, k)
}

func (t *alphaSortedTree[K, V]) Size() int { return t.size }
