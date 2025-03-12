package art

import (
	"strings"
	"unsafe"
)

type Tree[K nodeKey, V any] struct {
	root taggedPointer
}

func New[K nodeKey, V any]() Tree[K, V] {
	return Tree[K, V]{}
}

func longestCommonPrefix(key, other string, depth uint32) uint32 {
	maxCmp := uint32(min(len(key), len(other))) - depth
	var idx uint32

	for idx = range maxCmp {
		if key[depth+idx] != other[depth+idx] {
			return idx
		}
	}

	return idx + 1
}

func minimum[K nodeKey, V any](ptr taggedPointer) *nodeLeaf[K, V] {
	if ptr.pointer() == nil {
		return nil
	}

	if nodeKind(ptr.tag()) == nodeKindLeaf {
		return (*nodeLeaf[K, V])(ptr.pointer())
	}

	switch nodeKind(ptr.tag()) {
	case nodeKind4:
		n4 := (*node4)(ptr.pointer())
		return minimum[K, V](n4.children[0])
	case nodeKind16:
		n16 := (*node16)(ptr.pointer())
		return minimum[K, V](n16.children[0])
	case nodeKind48:
		idx := 0
		n48 := (*node48)(ptr.pointer())

		for n48.keys[idx] == 0 {
			idx++
		}
		idx = int(n48.keys[idx]) - 1
		return minimum[K, V](n48.children[idx])
	case nodeKind256:
		idx := 0
		n256 := (*node256)(ptr.pointer())

		for n256.children[idx] == 0 {
			idx++
		}
		return minimum[K, V](n256.children[idx])
	default:
		panic("shouldn't be possible!")
	}
}

func prefixMismatch[K nodeKey, V any](n *taggedPointer, key string, keyLen, depth uint32) uint32 {
	node := (*node)(n.pointer())
	maxCmp := min(min(maxPrefixLen, node.prefixLen), keyLen-depth)

	var idx uint32
	for idx = range maxCmp {
		if node.prefix[idx] != key[depth+idx] {
			return idx
		}
	}
	idx++

	if node.prefixLen > maxPrefixLen {
		leaf := minimum[K, V](*n)
		leafKeyStr := unsafe.String(leaf.key, leaf.len)

		maxCmp = min(leaf.len, keyLen) - depth
		for ; idx < maxCmp; idx++ {
			if leafKeyStr[idx+depth] != key[depth+idx] {
				return idx
			}
		}
	}

	return idx
}

func (ptr *taggedPointer) findChild(b byte) *taggedPointer {
	switch nodeKind(ptr.tag()) {
	case nodeKind4:
		n4 := (*node4)(ptr.pointer())

		for i := range n4.childrenLen {
			if n4.keys[i] == b {
				return &n4.children[i]
			}
		}

	case nodeKind16:
		n16 := (*node16)(ptr.pointer())

		if idx := searchNode16(&n16.keys, n16.childrenLen, b); idx != -1 {
			return &n16.children[idx]
		}

	case nodeKind48:
		n48 := (*node48)(ptr.pointer())

		i := n48.keys[b]
		if i != 0 {
			return &n48.children[i-1]
		}

	case nodeKind256:
		n256 := (*node256)(ptr.pointer())

		if n256.children[b] != 0 {
			return &n256.children[b]
		}

	default:
		panic("shouldn't be possible!")
	}

	return nil
}

func (ptr *taggedPointer) addChild(ref *taggedPointer, b byte, child taggedPointer) {
	switch nodeKind(ptr.tag()) {
	case nodeKind4:
		n4 := (*node4)(ptr.pointer())
		n4.addChild(ref, b, child)

	case nodeKind16:
		n16 := (*node16)(ptr.pointer())
		n16.addChild(ref, b, child)

	case nodeKind48:
		n48 := (*node48)(ptr.pointer())
		n48.addChild(ref, b, child)

	case nodeKind256:
		n256 := (*node256)(ptr.pointer())
		n256.addChild(b, child)

	default:
		panic("shouldn't be possible!")
	}
}

func (t *Tree[K, V]) insert(n *taggedPointer, key string, leaf *nodeLeaf[K, V], depth uint32) *node {
	if n.pointer() == nil {
		*n = taggedPointerPack(unsafe.Pointer(leaf), uintptr(nodeKindLeaf))
		return nil
	}

	if nodeKind(n.tag()) == nodeKindLeaf { // leaf
		nl := (*nodeLeaf[K, V])(n.pointer())
		leafKeyStr := unsafe.String(nl.key, nl.len)

		if strings.Compare(key, leafKeyStr) == 0 {
			return nil
		}

		newNode := new(node4)

		longestPrefix := longestCommonPrefix(key, leafKeyStr, depth)
		newNode.prefixLen = longestPrefix

		copy(newNode.prefix[:], key[depth:depth+min(maxPrefixLen, longestPrefix)])

		splitPrefix := int(depth + longestPrefix)
		if splitPrefix < len(leafKeyStr) {
			newNode.addChild(n, leafKeyStr[splitPrefix], *n)
		}
		if splitPrefix < len(key) {
			newNode.addChild(n, key[splitPrefix], taggedPointerPack(unsafe.Pointer(leaf), uintptr(nodeKindLeaf)))
		}
		*n = taggedPointerPack(unsafe.Pointer(newNode), uintptr(nodeKind4))
		return &newNode.node
	}

	node := (*node)(n.pointer())
	if node.prefixLen != 0 {
		prefixDiff := prefixMismatch[K, V](n, key, uint32(len(key)), depth)

		if prefixDiff >= node.prefixLen {
			depth += node.prefixLen
			goto RECURSE_SEARCH
		}

		newNode := new(node4)
		*n = taggedPointerPack(unsafe.Pointer(newNode), uintptr(nodeKind4))

		newNode.prefixLen = prefixDiff
		copy(newNode.prefix[:], node.prefix[depth:depth+min(maxPrefixLen, prefixDiff)])

		if node.prefixLen <= maxPrefixLen {
			newNode.addChild(n, node.prefix[prefixDiff], *n)
			loLimit := prefixDiff + 1
			hiLimit := depth + loLimit + min(maxPrefixLen, node.prefixLen)
			node.prefixLen -= loLimit
			copy(node.prefix[:], node.prefix[loLimit:hiLimit])
		} else {
			node.prefixLen -= prefixDiff + 1
			leaf = minimum[K, V](*n)
			leafKeyStr := unsafe.String(leaf.key, leaf.len)
			newNode.addChild(n, leafKeyStr[depth+prefixDiff], *n)
			loLimit := depth + prefixDiff + 1
			hiLimit := loLimit + min(maxPrefixLen, node.prefixLen)
			copy(node.prefix[:], leafKeyStr[loLimit:hiLimit])
		}

		newNode.addChild(n, key[depth+prefixDiff], taggedPointerPack(unsafe.Pointer(leaf), uintptr(nodeKindLeaf)))
		return &newNode.node
	}

RECURSE_SEARCH:
	child := (*n).findChild(key[depth])
	if child != nil {
		return t.insert(child, key, leaf, depth+1)
	}

	(*n).addChild(n, key[depth], taggedPointerPack(unsafe.Pointer(leaf), uintptr(nodeKindLeaf)))
	return nil
}

func (t *Tree[K, V]) Insert(key K, val V) *nodeLeaf[K, V] {
	keyStr := string(key)
	leaf := &nodeLeaf[K, V]{
		key:   unsafe.StringData(keyStr),
		value: val,
		len:   uint32(len(key)),
	}
	_ = t.insert(&t.root, keyStr, leaf, 0)
	return leaf
}

func (n *node) checkPrefix(key string, depth uint32) uint32 {
	maxCmp := min(min(n.prefixLen, maxPrefixLen), uint32(len(key))-depth)

	var idx uint32
	for idx = range maxCmp {
		if n.prefix[idx] != key[depth+idx] {
			return idx
		}
	}
	return idx + 1
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

	keyStr := string(key)
	n := t.root
	depth := uint32(0)

	for n.pointer() != nil {
		if nodeKind(n.tag()) == nodeKindLeaf {
			leaf := (*nodeLeaf[K, V])(n.pointer())

			if leaf.leafMatches(keyStr) == 0 {
				return leaf.value, true
			}
			return notFound, false
		}

		node := (*node)(n.pointer())
		if node.prefixLen != 0 {
			prefixLen := node.checkPrefix(keyStr, depth)

			if prefixLen != min(maxPrefixLen, node.prefixLen) {
				return notFound, false
			}

			depth += node.prefixLen
		}

		if child := n.findChild(keyStr[depth]); child != nil {
			n = *child
		} else {
			n = 0
		}
		depth++
	}

	return notFound, false
}
