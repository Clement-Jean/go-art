package art

import (
	"iter"
	"unsafe"
)

type compoundLeafNode[K nodeKey, V any] struct {
	key   *byte
	value V
	len   uint32
}

func (n *compoundLeafNode[K, V]) getKey() *byte           { return n.key }
func (n *compoundLeafNode[K, V]) getTransformKey() *byte  { return n.key }
func (n *compoundLeafNode[K, V]) getLen() uint32          { return n.len }
func (n *compoundLeafNode[K, V]) getTransformLen() uint32 { return n.len }
func (n *compoundLeafNode[K, V]) getValue() V             { return n.value }

type compoundSortedTree[K any, V any] struct {
	root nodeRef
	bck  BinaryComparableKey[K]
	size int
}

func NewCompoundTree[K any, V any](bck BinaryComparableKey[K]) Tree[K, V] {
	return &compoundSortedTree[K, V]{
		bck: bck,
	}
}

func (t *compoundSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, func(l *compoundLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *compoundSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, func(l *compoundLeafNode[K, V]) K {
		keyStr := unsafe.Slice(l.key, l.len)
		keyStr = keyStr[:len(keyStr)-1] // drop end byte
		return t.bck.Restore(keyStr)
	})
}

func (t *compoundSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')

	ok := delete[K, V, *compoundLeafNode[K, V]](&t.root, keyStr, keyStr)

	if ok {
		t.size--
	}
	return ok
}

func (t *compoundSortedTree[K, V]) Insert(key K, val V) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')
	leaf := &compoundLeafNode[K, V]{
		key:   unsafe.SliceData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}

	if insert[K](&t.root, keyStr, keyStr, leaf) {
		t.size++
	}
}

func (t *compoundSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *compoundLeafNode[K, V]](t.root); leaf != nil {
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

func (t *compoundSortedTree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V, *compoundLeafNode[K, V]](t.root); leaf != nil {
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

func (t *compoundSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyStr := t.bck.Transform(key)
	keyStr = append(keyStr, '\x00')

	return search[K, V, *compoundLeafNode[K, V]](t.root, keyStr, keyStr)
}

func (t *compoundSortedTree[K, V]) Size() int { return t.size }
