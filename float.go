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
	bck  FloatBinaryKey[K]
}

func NewFloatBinaryTree[K floats, V any]() Tree[K, V] {
	return &floatSortedTree[K, V]{}
}

func (t *floatSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, func(l *floatLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *floatSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *floatLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *floatSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')

	return delete[K, V, *floatLeafNode[K, V]](&t.root, keyStr, keyStr)
}

func (t *floatSortedTree[K, V]) Insert(key K, val V) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')
	leaf := &floatLeafNode[K, V]{
		key:   unsafe.SliceData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}

	insert[K](&t.root, keyStr, keyStr, leaf)
}

func (t *floatSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *floatLeafNode[K, V]](t.root); leaf != nil {
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
	if leaf := minimum[K, V, *floatLeafNode[K, V]](t.root); leaf != nil {
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
	keyStr = append(keyStr, '\x00')

	return search[K, V, *floatLeafNode[K, V]](t.root, keyStr, keyStr)
}
