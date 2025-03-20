package art

import (
	"bytes"
	"iter"
	"strings"
	"unsafe"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type collateLeafNode[K nodeKey, V any] struct {
	colKey    *byte
	key       *byte
	value     V
	colKeyLen uint32
	len       uint32
}

func (n *collateLeafNode[K, V]) getKey() *byte {
	return n.key
}

func (n *collateLeafNode[K, V]) getLen() uint32 {
	return n.len
}

func (n *collateLeafNode[K, V]) getValue() V {
	return n.value
}

type collationSortedTree[K nodeKey, V any] struct {
	c    *collate.Collator
	root nodeRef
	buf  collate.Buffer
	end  byte
}

func (t *collationSortedTree[K, V]) setEnd(b byte) {
	t.end = b
}

func (t *collationSortedTree[K, V]) setCollator(c *collate.Collator) {
	t.c = c
}

func NewCollationSortedTree[K nodeKey, V any](opts ...func(Tree[K, V])) Tree[K, V] {
	t := &collationSortedTree[K, V]{
		end: '\x00',
		c:   collate.New(language.Und),
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithCollator[K nodeKey, V any](c *collate.Collator) func(Tree[K, V]) {
	return func(t Tree[K, V]) {
		t.setCollator(c)
	}
}

func prefixMismatchCollate[K nodeKey, V any](n nodeRef, key []byte, depth int) int {
	node := n.node()
	maxCmp := min(int(min(maxPrefixLen, node.prefixLen)), len(key)-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if node.prefix[idx] != key[depth+idx] {
			return idx
		}
	}

	if node.prefixLen > maxPrefixLen {
		leaf := minimum[K, V, *collateLeafNode[K, V]](n)
		leafColKey := unsafe.Slice(leaf.colKey, leaf.colKeyLen)

		maxCmp = min(int(len(leafColKey)), len(key)) - depth
		for ; idx < maxCmp; idx++ {
			realIdx := depth + idx
			if leafColKey[realIdx] != key[realIdx] {
				return idx
			}
		}
	}

	return idx
}

func (t *collationSortedTree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V, *collateLeafNode[K, V]](t.root); leaf != nil {
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

func (t *collationSortedTree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V, *collateLeafNode[K, V]](t.root); leaf != nil {
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

// Insert inserts a key-value pair in the tree.
func (t *collationSortedTree[K, V]) Insert(key K, val V) {
	keyStr := append([]byte(string(key)), t.end)
	colKey := t.c.Key(&t.buf, keyStr)
	leaf := &collateLeafNode[K, V]{
		colKey:    unsafe.SliceData(colKey),
		key:       unsafe.SliceData(keyStr),
		value:     val,
		colKeyLen: uint32(len(colKey)),
		len:       uint32(len(keyStr)),
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
			nl := (*collateLeafNode[K, V])(ref.pointer)
			leafKeyStr := unsafe.Slice(nl.key, nl.len)

			if bytes.Compare(keyStr, leafKeyStr) == 0 {
				return
			}

			leafColKey := unsafe.Slice(nl.colKey, nl.colKeyLen)
			newNode := new(node4)

			longestPrefix := longestCommonPrefix(leafColKey, colKey, depth)
			newNode.prefixLen = uint32(longestPrefix)

			copy(newNode.prefix[:], colKey[depth:])

			*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

			splitPrefix := int(depth + longestPrefix)
			newNode.addChild(ref, leafColKey[splitPrefix], n)
			newNode.addChild(ref, colKey[splitPrefix], leafRef)
			return
		}

		node := ref.node()
		if node.prefixLen != 0 {
			prefixDiff := prefixMismatchCollate[K, V](n, colKey, depth)

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
				leafMin := minimum[K, V, *collateLeafNode[K, V]](n)
				leafColKey := unsafe.Slice(leafMin.colKey, leafMin.colKeyLen)

				newNode.addChild(ref, leafColKey[depth+prefixDiff], n)
				loLimit := depth + prefixDiff + 1
				copy(node.prefix[:], leafColKey[loLimit:])
			}

			newNode.addChild(ref, colKey[depth+prefixDiff], leafRef)
			return
		}

	CONTINUE_SEARCH:
		child := ref.findChild(colKey[depth])
		if child != nil {
			n = *child
			ref = child
			depth++
			continue
		}

		ref.addChild(colKey[depth], leafRef)
		return
	}
}

// Search searches for element with the given key.
// It returns whether the key is present (bool) and its value if it is present.
func (t *collationSortedTree[K, V]) Search(key K) (V, bool) {
	var notFound V

	keyStr := append([]byte(string(key)), t.end)
	colKey := t.c.Key(&t.buf, keyStr)
	n := t.root
	depth := 0

	for n.pointer != nil {
		if n.tag == nodeKindLeaf {
			leaf := (*collateLeafNode[K, V])(n.pointer)
			leafKeyStr := unsafe.Slice(leaf.key, leaf.len)
			if bytes.Compare(leafKeyStr, keyStr) == 0 {
				return leaf.value, true
			}
			return notFound, false
		}

		node := n.node()
		if node.prefixLen != 0 {
			prefixLen := node.checkPrefix(colKey, depth)

			if prefixLen != int(min(maxPrefixLen, node.prefixLen)) {
				return notFound, false
			}

			depth += int(node.prefixLen)
		}

		if child := n.findChild(colKey[depth]); child != nil {
			n = *child
		} else {
			n = nodeRef{}
		}
		depth++
	}

	return notFound, false
}

// Delete deletes a element with the given key.
func (t *collationSortedTree[K, V]) Delete(key K) {
	if t.root.pointer == nil {
		return
	}

	keyStr := append([]byte(string(key)), t.end)
	colKey := t.c.Key(&t.buf, keyStr)
	n := t.root
	ref := &t.root
	depth := 0

	for {
		if n.tag == nodeKindLeaf {
			leaf := (*collateLeafNode[K, V])(n.pointer)
			leafKeyStr := unsafe.Slice(leaf.key, leaf.len)
			if bytes.Compare(leafKeyStr, keyStr) == 0 {
				ref.pointer = nil
				return
			}

			return
		}

		node := n.node()
		if node.prefixLen != 0 {
			prefixLen := node.checkPrefix(colKey, depth)
			if prefixLen != int(min(maxPrefixLen, node.prefixLen)) {
				return
			}
			depth = depth + int(node.prefixLen)
		}

		child := n.findChild(colKey[depth])

		if child == nil {
			return
		}

		if child.tag == nodeKindLeaf {
			leaf := (*collateLeafNode[K, V])(n.pointer)
			leafKeyStr := unsafe.Slice(leaf.key, leaf.len)

			if bytes.Compare(leafKeyStr, keyStr) == 0 {
				ref.deleteChild(colKey[depth])
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

// All returns an iterator over the tree in collation order.
func (t *collationSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all[K, V, *collateLeafNode[K, V]](t.root, t.end)
}

// Backward returns an iterator over the tree in reverse collation order.
func (t *collationSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward[K, V, *collateLeafNode[K, V]](t.root, t.end)
}
