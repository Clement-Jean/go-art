package art

import (
	"iter"
	"unsafe"
)

type unsignedLeafNode[K nodeKey, V any] struct {
	key   *byte
	value V
	len   uint32
}

func (n *unsignedLeafNode[K, V]) getKey() *byte           { return n.key }
func (n *unsignedLeafNode[K, V]) getTransformKey() *byte  { return n.key }
func (n *unsignedLeafNode[K, V]) getLen() uint32          { return n.len }
func (n *unsignedLeafNode[K, V]) getTransformLen() uint32 { return n.len }
func (n *unsignedLeafNode[K, V]) getValue() V             { return n.value }

type unsignedSortedTree[K uints, V any] struct {
	root nodeRef
	bck  UnsignedBinaryKey[K]
}

func NewUnsignedBinaryTree[K uints, V any]() Tree[K, V] {
	return &unsignedSortedTree[K, V]{}
}

func (t *unsignedSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, func(l *unsignedLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *unsignedSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *unsignedLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *unsignedSortedTree[K, V]) Delete(key K) {
	if t.root.pointer == nil {
		return
	}

	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')

	delete[K, V, *unsignedLeafNode[K, V]](&t.root, keyStr, keyStr)
}

func (t *unsignedSortedTree[K, V]) Insert(key K, val V) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')
	leaf := &unsignedLeafNode[K, V]{
		key:   unsafe.SliceData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}

	insert[K](&t.root, keyStr, keyStr, leaf)
}

func (t *unsignedSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *unsignedLeafNode[K, V]](t.root); leaf != nil {
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

func (t *unsignedSortedTree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V, *unsignedLeafNode[K, V]](t.root); leaf != nil {
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

func (t *unsignedSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')

	return search[K, V, *unsignedLeafNode[K, V]](t.root, keyStr, keyStr)
}
