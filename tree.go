package art

import (
	"log"
	"strings"
	"unsafe"
)

type Tree[K nodeKey, V any] struct {
	root nodeRef
}

func New[K nodeKey, V any]() Tree[K, V] {
	return Tree[K, V]{}
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

		if idx >= int(maxNode48) {
			panic("YUP ERROR TOO")
		}

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
	node := (*node)(n.pointer)
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
			if i-1 >= maxNode48 {
				log.Fatalf("%d\n", i-1)
			}

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

func (t *Tree[K, V]) insert(n nodeRef, ref *nodeRef, key K, val V, depth int) {
	keyStr := string(key) + "\000"

	if ref.pointer == nil {
		leaf := &nodeLeaf[K, V]{
			key:   unsafe.StringData(keyStr),
			value: val,
			len:   uint32(len(keyStr)),
		}

		*ref = nodeRef{pointer: unsafe.Pointer(leaf), tag: nodeKindLeaf}
		return
	}

	if ref.tag == nodeKindLeaf {
		nl := (*nodeLeaf[K, V])(ref.pointer)
		leafKeyStr := unsafe.String(nl.key, nl.len)

		if strings.Compare(keyStr, leafKeyStr) == 0 {
			return
		}

		newNode := new(node4)

		longestPrefix := longestCommonPrefix(leafKeyStr, keyStr, depth)
		newNode.prefixLen = uint32(longestPrefix)

		copy(newNode.prefix[:], keyStr[depth:])

		*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

		leaf := &nodeLeaf[K, V]{
			key:   unsafe.StringData(keyStr),
			value: val,
			len:   uint32(len(keyStr)),
		}
		lr := nodeRef{pointer: unsafe.Pointer(leaf), tag: nodeKindLeaf}

		//println(nodeKind(ref.tag()).String())
		splitPrefix := int(depth + longestPrefix)
		//		if splitPrefix < len(leafKeyStr) {
		newNode.addChild(ref, leafKeyStr[splitPrefix], n)
		//}
		//if splitPrefix < len(key) {
		newNode.addChild(ref, keyStr[splitPrefix], lr)
		//}
		//println(unsafe.String(&newNode.prefix[0], newNode.prefixLen))
		return
	}

	node := (*node)(ref.pointer)

	if node.prefixLen != 0 {
		prefixDiff := prefixMismatch[K, V](n, keyStr, len(key), depth)

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

		leaf := &nodeLeaf[K, V]{
			key:   unsafe.StringData(keyStr),
			value: val,
			len:   uint32(len(keyStr)),
		}
		lr := nodeRef{pointer: unsafe.Pointer(leaf), tag: nodeKindLeaf}
		newNode.addChild(ref, keyStr[depth+prefixDiff], lr)
		return
	}

RECURSE_SEARCH:
	if depth >= len(key) {
		log.Fatalf("%s %d\n", keyStr, depth)
	}

	child := ref.findChild(keyStr[depth])
	if child != nil {
		t.insert(*child, child, key, val, depth+1)
		return
	}

	leaf := &nodeLeaf[K, V]{
		key:   unsafe.StringData(keyStr),
		value: val,
		len:   uint32(len(keyStr)),
	}
	lr := nodeRef{pointer: unsafe.Pointer(leaf), tag: nodeKindLeaf}

	ref.addChild(keyStr[depth], lr)
	return
}

func (t *Tree[K, V]) Insert(key K, val V) {
	t.insert(t.root, &t.root, key, val, 0)
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

	keyStr := string(key) + "\000"
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

		node := (*node)(n.pointer)
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
