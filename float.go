package art

import (
	"iter"
	"unsafe"
)

type floatLeafNode[K nodeKey, V any] struct {
	key   *byte
	value V
	len   uint32
}

func (n *floatLeafNode[K, V]) getKey() *byte           { return n.key }
func (n *floatLeafNode[K, V]) getTransformKey() *byte  { return n.key }
func (n *floatLeafNode[K, V]) getLen() uint32          { return n.len }
func (n *floatLeafNode[K, V]) getTransformLen() uint32 { return n.len }
func (n *floatLeafNode[K, V]) getValue() V             { return n.value }

type floatSortedTree[K floats, V any] struct {
	root nodeRef
	end  byte
	bck  FloatBinaryKey[K]
}

func NewFloatBinaryTree[K floats, V any](opts ...func(*floatSortedTree[K, V])) Tree[K, V] {
	t := &floatSortedTree[K, V]{
		end: '\x00',
		bck: FloatBinaryKey[K]{},
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func (t *floatSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, func(l *alphaLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *floatSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *alphaLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *floatSortedTree[K, V]) Delete(key K) {
	if t.root.pointer == nil {
		return
	}

	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, t.end)

	delete[K, V, *alphaLeafNode[K, V]](&t.root, keyStr, keyStr)
}

func (t *floatSortedTree[K, V]) Insert(key K, val V) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, t.end)
	leaf := &alphaLeafNode[K, V]{
		key:   unsafe.SliceData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}

	insert[K](&t.root, keyStr, keyStr, leaf)
}

func (t *floatSortedTree[K, V]) Maximum() (K, V, bool) {
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

func (t *floatSortedTree[K, V]) Minimum() (K, V, bool) {
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

func (t *floatSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, t.end)

	return search[K, V, *alphaLeafNode[K, V]](t.root, keyStr, keyStr)
}
