package art

import (
	"bytes"
	"iter"
	"strings"
	"unsafe"
)

type Tree[K nodeKey, V any] struct {
	root nodeRef
	end  byte
}

func New[K nodeKey, V any](opts ...func(*Tree[K, V])) *Tree[K, V] {
	t := &Tree[K, V]{
		end: '\x00',
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithEndByte[K nodeKey, V any](b byte) func(*Tree[K, V]) {
	return func(t *Tree[K, V]) {
		t.end = b
	}
}

func longestCommonPrefix(key, other []byte, depth int) int {
	maxCmp := min(len(key), len(other)) - depth

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if key[depth+idx] != other[depth+idx] {
			return idx
		}
	}

	return idx
}

func minimum[K nodeKey, V any](ref nodeRef) *nodeLeaf[K, V] {
	for ref.pointer != nil {
		kind := ref.tag
		if kind == nodeKindLeaf {
			return (*nodeLeaf[K, V])(ref.pointer)
		}

		switch kind {
		case nodeKind4:
			n4 := (*node4)(ref.pointer)
			ref = n4.children[0]
		case nodeKind16:
			n16 := (*node16)(ref.pointer)
			ref = n16.children[0]
		case nodeKind48:
			idx := 0
			n48 := (*node48)(ref.pointer)

			for n48.keys[idx] == 0 {
				idx++
			}
			idx = int(n48.keys[idx]) - 1
			ref = n48.children[idx]
		case nodeKind256:
			idx := 0
			n256 := (*node256)(ref.pointer)

			for n256.children[idx].pointer == nil {
				idx++
			}
			ref = n256.children[idx]
		default:
			panic("shouldn't be possible!")
		}
	}

	return nil
}

func (t *Tree[K, V]) Minimum() (K, V, bool) {
	if leaf := minimum[K, V](t.root); leaf != nil {
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

func maximum[K nodeKey, V any](ref nodeRef) *nodeLeaf[K, V] {
	if ref.pointer == nil {
		return nil
	}

	kind := ref.tag
	if kind == nodeKindLeaf {
		return (*nodeLeaf[K, V])(ref.pointer)
	}

	switch kind {
	case nodeKind4:
		n4 := (*node4)(ref.pointer)
		return maximum[K, V](n4.children[n4.childrenLen-1])
	case nodeKind16:
		n16 := (*node16)(ref.pointer)
		return maximum[K, V](n16.children[n16.childrenLen-1])
	case nodeKind48:
		idx := 255
		n48 := (*node48)(ref.pointer)

		for n48.keys[idx] == 0 {
			idx--
		}
		idx = int(n48.keys[idx]) - 1

		return maximum[K, V](n48.children[idx])
	case nodeKind256:
		idx := 255
		n256 := (*node256)(ref.pointer)

		for n256.children[idx].pointer == nil {
			idx--
		}
		return maximum[K, V](n256.children[idx])
	default:
		panic("shouldn't be possible!")
	}
}

func (t *Tree[K, V]) Maximum() (K, V, bool) {
	if leaf := maximum[K, V](t.root); leaf != nil {
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

func prefixMismatch[K nodeKey, V any](n nodeRef, key []byte, keyLen, depth int) int {
	node := n.node()
	maxCmp := min(int(min(maxPrefixLen, node.prefixLen)), keyLen-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if node.prefix[idx] != key[depth+idx] {
			return idx
		}
	}

	if node.prefixLen > maxPrefixLen {
		leaf := minimum[K, V](n)
		leafKeyStr := unsafe.Slice(leaf.key, leaf.len)

		maxCmp = min(int(leaf.len), keyLen) - depth
		for ; idx < maxCmp; idx++ {
			if leafKeyStr[idx+depth] != key[depth+idx] {
				return idx
			}
		}
	}

	return idx
}

// Insert inserts a key-value pair in the tree.
func (t *Tree[K, V]) Insert(key K, val V) {
	keyStr := []byte(string(key) + string(t.end))
	leaf := &nodeLeaf[K, V]{
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
			nl := (*nodeLeaf[K, V])(ref.pointer)
			leafKeyStr := unsafe.Slice(nl.key, nl.len)

			if bytes.Compare(keyStr, leafKeyStr) == 0 {
				return
			}

			newNode := new(node4)

			longestPrefix := longestCommonPrefix(leafKeyStr, keyStr, depth)
			newNode.prefixLen = uint32(longestPrefix)

			copy(newNode.prefix[:], keyStr[depth:])

			*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

			splitPrefix := int(depth + longestPrefix)
			newNode.addChild(ref, leafKeyStr[splitPrefix], n)
			newNode.addChild(ref, keyStr[splitPrefix], leafRef)
			return
		}

		node := ref.node()
		if node.prefixLen != 0 {
			prefixDiff := prefixMismatch[K, V](n, keyStr, len(key), depth)

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
				leafMin := minimum[K, V](n)
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

func (n *node) checkPrefix(key []byte, depth int) int {
	maxCmp := min(int(min(n.prefixLen, maxPrefixLen)), len(key)-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if n.prefix[idx] != key[depth+idx] {
			return idx
		}
	}
	return idx
}

func (l *nodeLeaf[K, V]) leafMatches(key []byte) int {
	if l.len != uint32(len(key)) {
		return 1
	}

	leafKeyStr := unsafe.Slice(l.key, l.len)
	return bytes.Compare(leafKeyStr, key)
}

// Search searches for element with the given key.
// It returns whether the key is present (bool) and its value if it is present.
func (t *Tree[K, V]) Search(key K) (V, bool) {
	var notFound V

	keyStr := []byte(string(key) + string(t.end))
	n := t.root
	depth := 0

	for n.pointer != nil {
		if n.tag == nodeKindLeaf {
			leaf := (*nodeLeaf[K, V])(n.pointer)

			if leaf.leafMatches(keyStr) == 0 {
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

// Delete deletes a element with the given key.
func (t *Tree[K, V]) Delete(key K) {
	keyStr := []byte(string(key) + string(t.end))

	if t.root.pointer == nil {
		return
	}

	n := t.root
	ref := &t.root
	depth := 0

	for {
		if n.tag == nodeKindLeaf {
			leaf := (*nodeLeaf[K, V])(n.pointer)

			if leaf.leafMatches(keyStr) == 0 {
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
			leaf := (*nodeLeaf[K, V])(n.pointer)

			if leaf.leafMatches(keyStr) == 0 {
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
func (t *Tree[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if t.root.pointer == nil {
			return
		}

		var q []nodeRef

		q = append(q, t.root)
		for len(q) != 0 {
			n := q[len(q)-1]
			q = q[:len(q)-1]

			if n.tag == nodeKindLeaf {
				leaf := (*nodeLeaf[K, V])(n.pointer)
				keyStr := unsafe.String(leaf.key, leaf.len)
				keyStr = strings.Trim(keyStr, string(t.end))

				if !yield(K(keyStr), leaf.value) {
					return
				}
				continue
			}

			switch n.tag {
			case nodeKind4:
				n4 := (*node4)(n.pointer)

				for i := int(n4.childrenLen) - 1; i >= 0; i-- {
					q = append(q, n4.children[i])
				}

			case nodeKind16:
				n16 := (*node16)(n.pointer)

				for i := int(n16.childrenLen) - 1; i >= 0; i-- {
					q = append(q, n16.children[i])
				}

			case nodeKind48:
				n48 := (*node48)(n.pointer)

				for i := 255; i >= 0; i-- {
					idx := n48.keys[i]
					if idx == 0 {
						continue
					}
					q = append(q, n48.children[idx-1])
				}

			case nodeKind256:
				n256 := (*node256)(n.pointer)

				for i := 255; i >= 0; i-- {
					if n256.children[i].pointer == nil {
						continue
					}
					q = append(q, n256.children[i])
				}

			default:
				panic("shouldn't be possible!")
			}
		}
	}
}

// Backward returns an iterator over the tree in reverse alphabetical order.
func (t *Tree[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if t.root.pointer == nil {
			return
		}

		var q []nodeRef

		q = append(q, t.root)
		for len(q) != 0 {
			n := q[len(q)-1]
			q = q[:len(q)-1]

			if n.tag == nodeKindLeaf {
				leaf := (*nodeLeaf[K, V])(n.pointer)
				keyStr := unsafe.String(leaf.key, leaf.len)
				keyStr = strings.Trim(keyStr, string(t.end))

				if !yield(K(keyStr), leaf.value) {
					return
				}
				continue
			}

			switch n.tag {
			case nodeKind4:
				n4 := (*node4)(n.pointer)

				for i := uint8(0); i < n4.childrenLen; i++ {
					q = append(q, n4.children[i])
				}

			case nodeKind16:
				n16 := (*node16)(n.pointer)

				for i := uint8(0); i < n16.childrenLen; i++ {
					q = append(q, n16.children[i])
				}

			case nodeKind48:
				n48 := (*node48)(n.pointer)

				for i := 0; i < 256; i++ {
					idx := n48.keys[i]
					if idx == 0 {
						continue
					}
					q = append(q, n48.children[idx-1])
				}

			case nodeKind256:
				n256 := (*node256)(n.pointer)

				for i := 0; i < 256; i++ {
					if n256.children[i].pointer == nil {
						continue
					}
					q = append(q, n256.children[i])
				}

			default:
				panic("shouldn't be possible!")
			}
		}
	}
}
