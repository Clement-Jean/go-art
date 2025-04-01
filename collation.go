package art

import (
	"bytes"
	"iter"
	"strings"
	"unsafe"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type collateLeafNode[V any] struct {
	value     V
	colKey    *byte
	key       *byte
	colKeyLen uint32
	keyLen    uint32
}

func (n *collateLeafNode[V]) getKey() []byte          { return unsafe.Slice(n.key, n.keyLen) }
func (n *collateLeafNode[V]) getTransformKey() []byte { return unsafe.Slice(n.colKey, n.colKeyLen) }
func (n *collateLeafNode[V]) getValue() V             { return n.value }
func (n *collateLeafNode[V]) setValue(val V)          { n.value = val }

type collationSortedTree[K chars | []rune, V any] struct {
	buf  *collate.Buffer
	c    *collate.Collator
	root nodeRef
	size int
}

func NewCollationSortedTree[K chars | []rune, V any](opts ...func(*collationSortedTree[K, V])) Tree[K, V] {
	t := &collationSortedTree[K, V]{
		c:   collate.New(language.Und),
		buf: &collate.Buffer{},
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithCollator[K chars, V any](c *collate.Collator) func(*collationSortedTree[K, V]) {
	return func(t *collationSortedTree[K, V]) {
		t.c = c
	}
}

func (t *collationSortedTree[K, V]) restoreKey(ptr unsafe.Pointer) (K, V) {
	l := (*collateLeafNode[V])(ptr)
	return K(string(l.getKey())), l.value
}

// All returns an iterator over the tree in collation order.
func (t *collationSortedTree[K, V]) All() iter.Seq2[K, V] {
	return all(t.root, t.restoreKey)
}

// Backward returns an iterator over the tree in reverse collation order.
func (t *collationSortedTree[K, V]) Backward() iter.Seq2[K, V] {
	return backward(t.root, t.restoreKey)
}

func (t *collationSortedTree[K, V]) BottomK(k uint) iter.Seq2[K, V] {
	return bottomK(t, k)
}

// Delete deletes a element with the given key.
func (t *collationSortedTree[K, V]) Delete(key K) bool {
	if t.root.pointer == nil {
		return false
	}

	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	keyS, colKey := bck.Transform(key)

	ref := &t.root
	n := *ref
	depth := 0

	for n.pointer != nil {
		if n.tag == nodeKindLeaf {
			leaf := (*collateLeafNode[V])(n.pointer)

			if bytes.Compare(leaf.getKey(), keyS) == 0 {
				*ref = nodeRef{}
				t.size--
				return true
			}

			return false
		}

		node := n.node()
		if node.prefixLen != 0 {
			prefixLen := node.checkPrefix(colKey, depth)
			if prefixLen != int(min(maxPrefixLen, node.prefixLen)) {
				return false
			}
			depth += int(node.prefixLen)
		}

		child := n.findChild(colKey[depth])

		if child == nil {
			return false
		}

		if child.tag == nodeKindLeaf {
			leaf := (*collateLeafNode[V])(child.pointer)

			if bytes.Compare(leaf.getKey(), keyS) == 0 {
				ref.deleteChild(colKey[depth])
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

// Insert inserts a key-value pair in the tree.
func (t *collationSortedTree[K, V]) Insert(key K, val V) {
	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}

	keyS, colKey := bck.Transform(key)

	createLeaf := func() unsafe.Pointer {
		return unsafe.Pointer(&collateLeafNode[V]{
			colKey:    unsafe.SliceData(colKey),
			key:       unsafe.SliceData(keyS),
			value:     val,
			keyLen:    uint32(len(keyS)),
			colKeyLen: uint32(len(colKey)),
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
			nl := (*collateLeafNode[V])(ref.pointer)

			if bytes.Compare(keyS, nl.getKey()) == 0 {
				nl.setValue(val)
				return
			}

			leafKey := nl.getTransformKey()
			newNode := nodePools[nodeKind4].Get().(*node4)

			longestPrefix := longestCommonPrefix(leafKey, colKey, depth)
			newNode.prefixLen = uint32(longestPrefix)

			copy(newNode.prefix[:], colKey[depth:])

			*ref = nodeRef{pointer: unsafe.Pointer(newNode), tag: nodeKind4}

			splitPrefix := int(depth + longestPrefix)
			if splitPrefix < len(leafKey) {
				newNode.addChild(ref, leafKey[splitPrefix], n)
			}

			if splitPrefix < len(colKey) {
				leafRef := nodeRef{pointer: createLeaf(), tag: nodeKindLeaf}
				newNode.addChild(ref, colKey[splitPrefix], leafRef)
			}
			t.size++
			return
		}

		node := ref.node()
		if node.prefixLen != 0 {
			prefixDiff := prefixMismatch[V, *collateLeafNode[V]](n, colKey, depth)

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
				leafMin := (*collateLeafNode[V])(minimum[V](n))
				leafKey := leafMin.getTransformKey()

				newNode.addChild(ref, leafKey[depth+prefixDiff], n)
				loLimit := depth + prefixDiff + 1
				copy(node.prefix[:], leafKey[loLimit:])
			}

			if depth+prefixDiff >= len(colKey) {
				return
			}
			leafRef := nodeRef{pointer: createLeaf(), tag: nodeKindLeaf}
			newNode.addChild(ref, colKey[depth+prefixDiff], leafRef)
			t.size++
			return
		}

	CONTINUE_SEARCH:
		if depth >= len(colKey) {
			return
		}

		child := ref.findChild(colKey[depth])
		if child != nil {
			n = *child
			ref = child
			depth++
			continue
		}

		leafRef := nodeRef{pointer: createLeaf(), tag: nodeKindLeaf}
		ref.addChild(colKey[depth], leafRef)
		t.size++
		return
	}
}

func (t *collationSortedTree[K, V]) Maximum() (K, V, bool) {
	if l := maximum[V](t.root); l != nil {
		k, v := t.restoreKey(l)
		return k, v, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *collationSortedTree[K, V]) Minimum() (K, V, bool) {
	if l := minimum[V](t.root); l != nil {
		k, v := t.restoreKey(l)
		return k, v, true
	}

	var (
		notFoundKey   K
		notFoundValue V
	)
	return notFoundKey, notFoundValue, false
}

func (t *collationSortedTree[K, V]) Prefix(p K) iter.Seq2[K, V] {
	if len(p) == 0 {
		return t.All()
	}

	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	keyS, colKey := bck.Transform(p)

	root := t.root
	if t.root.pointer != nil {
		root = lowestCommonParent[V, *collateLeafNode[V]](root, colKey)
	}

	hasPrefix := func(k K, v V) bool {
		leafKeyS := []byte(string(k))
		return bytes.HasPrefix(leafKeyS, keyS)
	}
	return filter(root, hasPrefix, t.restoreKey)
}

func (t *collationSortedTree[K, V]) Range(start, end K) iter.Seq2[K, V] {
	if len(end) == 0 {
		end, _ = t.restoreKey(maximum[V](t.root))
	}

	if strings.Compare(string(start), string(end)) > 0 { // start > end
		// IDEA: maybe do the iteration in reverse instead?
		start, end = end, start
	}

	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	startKey, startColKey := bck.Transform(start)
	endKey, endColKey := bck.Transform(end)

	return rangeScan[K, V, *collateLeafNode[V]](t.root, startKey, endKey, startColKey, endColKey, t.restoreKey)
}

// Search searches for element with the given key.
// It returns whether the key is present (bool) and its value if it is present.
func (t *collationSortedTree[K, V]) Search(key K) (V, bool) {
	bck := CollationOrderKey[K]{
		buf: t.buf,
		c:   t.c,
	}
	keyS, colKey := bck.Transform(key)

	var notFound V

	n := t.root
	depth := 0

	for n.pointer != nil {
		if n.tag == nodeKindLeaf {
			leaf := (*collateLeafNode[V])(n.pointer)

			if bytes.Compare(leaf.getKey(), keyS) == 0 {
				return leaf.getValue(), true
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

		b := colKey[depth]
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

func (t *collationSortedTree[K, V]) TopK(k uint) iter.Seq2[K, V] {
	return topK(t, k)
}

func (t *collationSortedTree[K, V]) Size() int { return t.size }
