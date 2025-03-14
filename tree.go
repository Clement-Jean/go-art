package art

import (
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

func longestCommonPrefix(key, other string, depth int) int {
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
		return minimum[K, V](n4.children[0])
	case nodeKind16:
		n16 := (*node16)(ref.pointer)
		return minimum[K, V](n16.children[0])
	case nodeKind48:
		idx := 0
		n48 := (*node48)(ref.pointer)

		for n48.keys[idx] == 0 {
			idx++
		}
		idx = int(n48.keys[idx]) - 1

		return minimum[K, V](n48.children[idx])
	case nodeKind256:
		idx := 0
		n256 := (*node256)(ref.pointer)

		for n256.children[idx].pointer == nil {
			idx++
		}
		return minimum[K, V](n256.children[idx])
	default:
		panic("shouldn't be possible!")
	}
}

func prefixMismatch[K nodeKey, V any](n nodeRef, key string, keyLen, depth int) int {
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
		leafKeyStr := unsafe.String(leaf.key, leaf.len)

		maxCmp = min(int(leaf.len), keyLen) - depth
		for ; idx < maxCmp; idx++ {
			if leafKeyStr[idx+depth] != key[depth+idx] {
				return idx
			}
		}
	}

	return idx
}

func (ref *nodeRef) node() *node {
	switch nodeKind(ref.tag) {
	case nodeKind4:
		n4 := (*node4)(ref.pointer)
		return &n4.node

	case nodeKind16:
		n16 := (*node16)(ref.pointer)
		return &n16.node

	case nodeKind48:
		n48 := (*node48)(ref.pointer)
		return &n48.node

	case nodeKind256:
		n256 := (*node256)(ref.pointer)
		return &n256.node

	default:
		panic("shouldn't be possible!")
	}
}

func (ref *nodeRef) findChild(b byte) *nodeRef {
	switch nodeKind(ref.tag) {
	case nodeKind4:
		n4 := (*node4)(ref.pointer)

		for i := uint8(0); i < n4.childrenLen; i++ {
			if n4.keys[i] == b {
				return &n4.children[i]
			}
		}

	case nodeKind16:
		n16 := (*node16)(ref.pointer)

		if idx := searchNode16(&n16.keys, n16.childrenLen, b); idx != -1 {
			return &n16.children[idx]
		}

	case nodeKind48:
		n48 := (*node48)(ref.pointer)

		i := n48.keys[b]
		if i != 0 {
			return &n48.children[i-1]
		}

	case nodeKind256:
		n256 := (*node256)(ref.pointer)

		if n256.children[b].pointer != nil {
			return &n256.children[b]
		}

	default:
		panic("shouldn't be possible!")
	}

	return nil
}

func (ptr *nodeRef) addChild(b byte, child nodeRef) {
	switch nodeKind(ptr.tag) {
	case nodeKind4:
		n4 := (*node4)(ptr.pointer)
		n4.addChild(ptr, b, child)

	case nodeKind16:
		n16 := (*node16)(ptr.pointer)
		n16.addChild(ptr, b, child)

	case nodeKind48:
		n48 := (*node48)(ptr.pointer)
		n48.addChild(ptr, b, child)

	case nodeKind256:
		n256 := (*node256)(ptr.pointer)
		n256.addChild(b, child)

	default:
		panic("shouldn't be possible!")
	}
}

func (t *Tree[K, V]) insert(n nodeRef, ref *nodeRef, key string, leaf nodeRef, depth int) {
	if ref.pointer == nil {
		*ref = leaf
		return
	}

	if ref.tag == nodeKindLeaf {
		nl := (*nodeLeaf[K, V])(ref.pointer)
		leafKeyStr := unsafe.String(nl.key, nl.len)

		if strings.Compare(key, leafKeyStr) == 0 {
			return
		}

		newNode := new(node4)

		longestPrefix := longestCommonPrefix(leafKeyStr, key, depth)
		newNode.prefixLen = uint32(longestPrefix)

		copy(newNode.prefix[:], key[depth:])

		*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

		splitPrefix := int(depth + longestPrefix)
		newNode.addChild(ref, leafKeyStr[splitPrefix], n)
		newNode.addChild(ref, key[splitPrefix], leaf)
		return
	}

	node := ref.node()
	if node.prefixLen != 0 {
		prefixDiff := prefixMismatch[K, V](n, key, len(key), depth)

		if prefixDiff >= int(node.prefixLen) {
			depth += int(node.prefixLen)
			goto RECURSE_SEARCH
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
			leafKeyStr := unsafe.String(leafMin.key, leafMin.len)
			newNode.addChild(ref, leafKeyStr[depth+prefixDiff], n)
			loLimit := depth + prefixDiff + 1
			copy(node.prefix[:], leafKeyStr[loLimit:])
		}

		newNode.addChild(ref, key[depth+prefixDiff], leaf)
		return
	}

RECURSE_SEARCH:
	child := ref.findChild(key[depth])
	if child != nil {
		t.insert(*child, child, key, leaf, depth+1)
		return
	}

	ref.addChild(key[depth], leaf)
	return
}

func (t *Tree[K, V]) Insert(key K, val V) {
	keyStr := string(key) + string(t.end)
	leaf := &nodeLeaf[K, V]{
		key:   unsafe.StringData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}
	leafRef := nodeRef{pointer: unsafe.Pointer(leaf), tag: nodeKindLeaf}
	t.insert(t.root, &t.root, keyStr, leafRef, 0)
}

func (n *node) checkPrefix(key string, depth int) int {
	maxCmp := min(int(min(n.prefixLen, maxPrefixLen)), len(key)-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if n.prefix[idx] != key[depth+idx] {
			return idx
		}
	}
	return idx
}

func (l *nodeLeaf[K, V]) leafMatches(key string) int {
	if l.len != uint32(len(key)) {
		return 1
	}

	leafKeyStr := unsafe.String(l.key, l.len)
	return strings.Compare(leafKeyStr, key)
}

func (t *Tree[K, V]) Search(key K) (V, bool) {
	var notFound V

	keyStr := string(key) + string(t.end)
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
