package art

import (
	"bytes"
	"iter"
	"strings"
	"unsafe"

	"golang.org/x/text/collate"
)

type alphaLeafNode[K nodeKey, V any] struct {
	key   *byte
	value V
	len   uint32
}

func (n *alphaLeafNode[K, V]) getKey() *byte {
	return n.key
}

func (n *alphaLeafNode[K, V]) getLen() uint32 {
	return n.len
}

func (n *alphaLeafNode[K, V]) getValue() V {
	return n.value
}

type alphaSortedTree[K nodeKey, V any] struct {
	root nodeRef
	end  byte
}

func (t *alphaSortedTree[K, V]) setEnd(b byte) {
	t.end = b
}

func (t *alphaSortedTree[K, V]) setCollator(c *collate.Collator) {}

func NewAlphaSortedTree[K nodeKey, V any](opts ...func(Tree[K, V])) Tree[K, V] {
	t := &alphaSortedTree[K, V]{
		end: '\x00',
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithEndByte[K nodeKey, V any](b byte) func(Tree[K, V]) {
	return func(t Tree[K, V]) {
		t.setEnd(b)
	}
}

func prefixMismatchAlpha[K nodeKey, V any](n nodeRef, key []byte, depth int) int {
	node := n.node()
	maxCmp := min(int(min(maxPrefixLen, node.prefixLen)), len(key)-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if node.prefix[idx] != key[depth+idx] {
			return idx
		}
	}

	if node.prefixLen > maxPrefixLen {
		leaf := minimum[K, V, *alphaLeafNode[K, V]](n)
		leafKeyStr := unsafe.Slice(leaf.key, leaf.len)

		maxCmp = min(int(len(leafKeyStr)), len(key)) - depth
		for ; idx < maxCmp; idx++ {
			realIdx := depth + idx
			if leafKeyStr[realIdx] != key[realIdx] {
				return idx
			}
		}
	}

	return idx
}

func (t *alphaSortedTree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V, *alphaLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.String(leaf.key, leaf.len)
		keyStr = strings.Trim(keyStr, string(t.end))
		return K(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *alphaLeafNode[K, V]](t.root); leaf != nil {
		keyStr := unsafe.String(leaf.key, leaf.len)
		keyStr = strings.Trim(keyStr, string(t.end))
		return K(keyStr), leaf.value, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Insert(key K, val V) {
	keyStr := append([]byte(string(key)), t.end)
	leaf := &alphaLeafNode[K, V]{
		key:   unsafe.SliceData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}
	leafRef := nodeRef{pointer: unsafe.Pointer(leaf), tag: nodeKindLeaf}
	n := t.root

	if t.root.pointer == nil {
		t.root = leafRef
		return
	}

	ref := &t.root
	depth := 0

	for {
		if ref.tag == nodeKindLeaf {
			nl := (*alphaLeafNode[K, V])(ref.pointer)
			leafKeyStr := unsafe.Slice(nl.key, nl.len)

			if bytes.Compare(keyStr, leafKeyStr) == 0 {
				return
			}

			newNode := new(node4)

			longestPrefix := longestCommonPrefix(leafKeyStr, keyStr, depth)
			newNode.prefixLen = uint32(longestPrefix)

			copy(newNode.prefix[:], keyStr[depth:])

			*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

			splitPrefix := depth + longestPrefix
			newNode.addChild(ref, leafKeyStr[splitPrefix], n)
			newNode.addChild(ref, keyStr[splitPrefix], leafRef)
			return
		}

		node := ref.node()
		if node.prefixLen != 0 {
			prefixDiff := prefixMismatchAlpha[K, V](n, keyStr, depth)

			if prefixDiff >= int(node.prefixLen) {
				depth += int(node.prefixLen)
				goto CONTINUE_SEARCH
			}

			newNode := new(node4)

			*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

			newNode.prefixLen = uint32(prefixDiff)
			copy(newNode.prefix[:], node.prefix[:])

			if node.prefixLen <= maxPrefixLen {
				newNode.addChild(ref, node.prefix[prefixDiff], n)
				loLimit := prefixDiff + 1
				node.prefixLen -= uint32(loLimit)
				copy(node.prefix[:], node.prefix[loLimit:])
			} else {
				node.prefixLen -= uint32(prefixDiff + 1)
				leafMin := minimum[K, V, *alphaLeafNode[K, V]](n)
				leafKeyStr := unsafe.Slice(leafMin.key, leafMin.len)

				newNode.addChild(ref, leafKeyStr[depth+prefixDiff], n)
				loLimit := depth + prefixDiff + 1
				copy(node.prefix[:], leafKeyStr[loLimit:])
			}

			newNode.addChild(ref, keyStr[depth+prefixDiff], leafRef)
			return
		}

	CONTINUE_SEARCH:
		child := ref.findChild(keyStr[depth])
		if child != nil {
			n = *child
			ref = child
			depth++
			continue
		}

		ref.addChild(keyStr[depth], leafRef)
		return
	}
}

func (t *alphaSortedTree[K, V]) Search(key K) (V, bool) {
	var notFound V

	keyStr := append([]byte(string(key)), t.end)
	n := t.root
	depth := 0

	for n.pointer != nil {
		if n.tag == nodeKindLeaf {
			leaf := (*alphaLeafNode[K, V])(n.pointer)
			leafKeyStr := unsafe.Slice(leaf.key, leaf.len)
			if bytes.Compare(leafKeyStr, keyStr) == 0 {
				return leaf.value, true
			}
			return notFound, false
		}

		node := n.node()
		if node.prefixLen != 0 {
			prefixLen := node.checkPrefix(keyStr, depth)

			if prefixLen != int(min(maxPrefixLen, node.prefixLen)) {
				return notFound, false
			}

			depth += int(node.prefixLen)
		}

		if child := n.findChild(keyStr[depth]); child != nil {
			n = *child
		} else {
			n = nodeRef{}
		}
		depth++
	}

	return notFound, false
}

func (t *alphaSortedTree[K, V]) Delete(key K) {
	if t.root.pointer == nil {
		return
	}

	keyStr := append([]byte(string(key)), t.end)
	n := t.root
	ref := &t.root
	depth := 0

	for {
		if n.tag == nodeKindLeaf {
			leaf := (*alphaLeafNode[K, V])(n.pointer)
			leafKeyStr := unsafe.Slice(leaf.key, leaf.len)
			if bytes.Compare(leafKeyStr, keyStr) == 0 {
				ref.pointer = nil
				return
			}

			return
		}

		node := n.node()
		if node.prefixLen != 0 {
			prefixLen := node.checkPrefix(keyStr, depth)
			if prefixLen != int(min(maxPrefixLen, node.prefixLen)) {
				return
			}
			depth = depth + int(node.prefixLen)
		}

		child := n.findChild(keyStr[depth])

		if child == nil {
			return
		}

		if child.tag == nodeKindLeaf {
			leaf := (*alphaLeafNode[K, V])(n.pointer)
			leafKeyStr := unsafe.Slice(leaf.key, leaf.len)

			if bytes.Compare(leafKeyStr, keyStr) == 0 {
				ref.deleteChild(keyStr[depth])
				return
			}

			return
		} else {
			n = *child
			ref = child
			depth++
		}
	}
}

// All returns an iterator over the tree in alphabetical order.
func (t *alphaSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all[K, V, *alphaLeafNode[K, V]](t.root, t.end)
}

// Backward returns an iterator over the tree in reverse alphabetical order.
func (t *alphaSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward[K, V, *alphaLeafNode[K, V]](t.root, t.end)
}
