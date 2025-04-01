package art

import (
	"bytes"
	"iter"
	"unsafe"
)

type alphaLeafNode[K nodeKey, V any] struct {
	key   *byte
	value V
	len   uint32
}

func (n *alphaLeafNode[K, V]) getKey() []byte          { return unsafe.Slice(n.key, n.len) }
func (n *alphaLeafNode[K, V]) getTransformKey() []byte { return unsafe.Slice(n.key, n.len) }
func (n *alphaLeafNode[K, V]) getValue() V             { return n.value }
func (n *alphaLeafNode[K, V]) setValue(val V)          { n.value = val }

type alphaSortedTree[K chars, V any] struct {
	root nodeRef
	bck  AlphabeticalOrderKey[K]
	size int
}

func NewAlphaSortedTree[K chars, V any]() Tree[K, V] {
	return &alphaSortedTree[K, V]{}
}

func (t *alphaSortedTree[K, V]) restoreKey(ptr unsafe.Pointer) (K, V) {
	l := (*alphaLeafNode[K, V])(ptr)
	keyS := l.getKey()
	keyS = keyS[:len(keyS)-1] // drop end byte
	return t.bck.Restore(keyS), l.value
}

// All returns an iterator over the tree in alphabetical order.
func (t *alphaSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, t.restoreKey)
}

// Backward returns an iterator over the tree in reverse alphabetical order.
func (t *alphaSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, t.restoreKey)
}

func (t *alphaSortedTree[K, V]) BottomK(k uint) iter.Seq2[K, V] {
	return bottomK(t, k)
}

func (t *alphaSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	ref := &t.root
	n := *ref
	depth := 0

	for n.pointer != nil {
		if n.tag == nodeKindLeaf {
			leaf := (*alphaLeafNode[K, V])(n.pointer)

			if bytes.Compare(leaf.getKey(), keyS) == 0 {
				*ref = nodeRef{}
				t.size--
				return true
			}

			return false
		}

		node := n.node()
		if node.prefixLen != 0 {
			prefixLen := node.checkPrefix(keyS, depth)
			if prefixLen != int(min(maxPrefixLen, node.prefixLen)) {
				return false
			}
			depth += int(node.prefixLen)
		}

		child := n.findChild(keyS[depth])

		if child == nil {
			return false
		}

		if child.tag == nodeKindLeaf {
			leaf := (*alphaLeafNode[K, V])(child.pointer)

			if bytes.Compare(leaf.getKey(), keyS) == 0 {
				ref.deleteChild(keyS[depth])
				t.size--
				return true
			}

			return false
		}

		n = *child
		ref = child
		depth++
	}
	return false
}

func (t *alphaSortedTree[K, V]) Insert(key K, val V) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	createLeaf := func() unsafe.Pointer {
		return unsafe.Pointer(&alphaLeafNode[K, V]{
			key:   unsafe.SliceData(keyS),
			value: val,
			len:   uint32(len(keyS)),
		})
	}

	if t.root.pointer == nil {
		t.root = nodeRef{pointer: createLeaf(), tag: nodeKindLeaf}
		t.size++
		return
	}

	ref := &t.root
	n := *ref
	depth := 0

	for ref.pointer != nil {
		if ref.tag == nodeKindLeaf {
			nl := (*alphaLeafNode[K, V])(ref.pointer)

			if bytes.Compare(keyS, nl.getKey()) == 0 {
				nl.setValue(val)
				return
			}

			leafKey := nl.getTransformKey()
			newNode := nodePools[nodeKind4].Get().(*node4)

			longestPrefix := longestCommonPrefix(leafKey, keyS, depth)
			newNode.prefixLen = uint32(longestPrefix)

			copy(newNode.prefix[:], keyS[depth:])

			*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

			splitPrefix := int(depth + longestPrefix)
			if splitPrefix < len(leafKey) {
				newNode.addChild(ref, leafKey[splitPrefix], n)
			}

			if splitPrefix < len(keyS) {
				leafRef := nodeRef{pointer: createLeaf(), tag: nodeKindLeaf}
				newNode.addChild(ref, keyS[splitPrefix], leafRef)
			}
			t.size++
			return
		}

		node := ref.node()
		if node.prefixLen != 0 {
			prefixDiff := prefixMismatch[K, V, *alphaLeafNode[K, V]](n, keyS, depth)

			if prefixDiff >= int(node.prefixLen) {
				depth += int(node.prefixLen)
				goto CONTINUE_SEARCH
			}

			newNode := nodePools[nodeKind4].Get().(*node4)

			*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

			newNode.prefixLen = uint32(prefixDiff)
			newNode.prefix = node.prefix

			if node.prefixLen <= maxPrefixLen {
				newNode.addChild(ref, node.prefix[prefixDiff], n)
				loLimit := prefixDiff + 1
				node.prefixLen -= uint32(loLimit)
				copy(node.prefix[:], node.prefix[loLimit:])
			} else {
				node.prefixLen -= uint32(prefixDiff + 1)
				leafMin := (*alphaLeafNode[K, V])(minimum[K, V](n))
				leafKey := leafMin.getTransformKey()

				newNode.addChild(ref, leafKey[depth+prefixDiff], n)
				loLimit := depth + prefixDiff + 1
				copy(node.prefix[:], leafKey[loLimit:])
			}

			if depth+prefixDiff >= len(keyS) {
				return
			}
			leafRef := nodeRef{pointer: createLeaf(), tag: nodeKindLeaf}
			newNode.addChild(ref, keyS[depth+prefixDiff], leafRef)
			return
		}

	CONTINUE_SEARCH:
		if depth >= len(keyS) {
			return
		}

		child := ref.findChild(keyS[depth])
		if child != nil {
			n = *child
			ref = child
			depth++
			continue
		}

		leafRef := nodeRef{pointer: createLeaf(), tag: nodeKindLeaf}
		ref.addChild(keyS[depth], leafRef)
		t.size++
		return
	}
}

func (t *alphaSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[K, V](t.root); l != nil {
		k, v := t.restoreKey(l)
		return k, v, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[K, V](t.root); l != nil {
		k, v := t.restoreKey(l)
		return k, v, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *alphaSortedTree[K, V]) Prefix(p K) iter.Seq2[K, V] {
	if len(p) == 0 {
		return t.All()
	}

	root := t.root
	if t.root.pointer != nil {
		root = lowestCommonParent[K, V, *alphaLeafNode[K, V]](root, []byte(p))
	}

	hasPrefix := func(k K, v V) bool { return bytes.HasPrefix([]byte(k), []byte(p)) }
	return filter(root, hasPrefix, t.restoreKey)
}

func (t *alphaSortedTree[K, V]) Range(start, end K) iter.Seq2[K, V] {
	if len(end) == 0 {
		end, _ = t.restoreKey(maximum[K, V](t.root))
	}

	if bytes.Compare([]byte(start), []byte(end)) > 0 { // start > end
		// IDEA: maybe do the iteration in reverse instead?
		start, end = end, start
	}

	startKey := append([]byte(start), '\x00')
	endKey := append([]byte(end), '\x00')
	return rangeScan[K, V, *alphaLeafNode[K, V]](t.root, startKey, endKey, startKey, endKey, t.restoreKey)
}

func (t *alphaSortedTree[K, V]) Search(key K) (V, bool) {
	_, keyS := t.bck.Transform(key)
	keyS = append(keyS, '\x00')

	var notFound V

	n := t.root
	depth := 0

	for n.pointer != nil {
		if n.tag == nodeKindLeaf {
			leaf := (*alphaLeafNode[K, V])(n.pointer)

			if bytes.Compare(leaf.getKey(), keyS) == 0 {
				return leaf.getValue(), true
			}
			return notFound, false
		}

		node := n.node()
		if node.prefixLen != 0 {
			prefixLen := node.checkPrefix(keyS, depth)

			if prefixLen != int(min(maxPrefixLen, node.prefixLen)) {
				return notFound, false
			}

			depth += int(node.prefixLen)
		}

		b := keyS[depth]
		switch n.tag {
		case nodeKind4:
			n4 := (*node4)(n.pointer)

			if i := searchNode4(n4.keys, b); i != -1 && i < int(n4.childrenLen) {
				n = n4.children[i]
				depth++
				continue
			}

		case nodeKind16:
			n16 := (*node16)(n.pointer)

			if idx := searchNode16(&n16.keys, n16.childrenLen, b); idx != -1 {
				n = n16.children[idx]
				depth++
				continue
			}

		case nodeKind48:
			n48 := (*node48)(n.pointer)

			if i := n48.keys[b]; i != 0 {
				n = n48.children[i-1]
				depth++
				continue
			}

		case nodeKind256:
			n256 := (*node256)(n.pointer)

			if n256.children[b].pointer != nil {
				n = n256.children[b]
				depth++
				continue
			}

		default:
			panic("shouldn't be possible!")
		}

		break
	}

	return notFound, false
}

func (t *alphaSortedTree[K, V]) TopK(k uint) iter.Seq2[K, V] {
	return topK(t, k)
}

func (t *alphaSortedTree[K, V]) Size() int { return t.size }
