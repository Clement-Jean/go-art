package art

import (
	"bytes"
	"iter"
	"unsafe"
)

type Tree[K nodeKey, V any] interface {
	// Insert inserts a key-value pair in the tree.
	Insert(K, V)

	// Search searches for element with the given key.
	// It returns whether the key is present (bool) and its value if it is present.
	Search(K) (V, bool)

	// Delete deletes a element with the given key.
	Delete(K) bool

	// Minimum find the minimum K/V pair based on the key.
	Minimum() (K, V, bool)

	// Maximum find the maximum K/V pair based on the key.
	Maximum() (K, V, bool)

	All() iter.Seq2[K, V]
	Backward() iter.Seq2[K, V]
	Prefix(K) iter.Seq2[K, V]
	TopK(uint) iter.Seq2[K, V]
	BottomK(uint) iter.Seq2[K, V]
	Range(K, K) iter.Seq2[K, V]

	Size() int
}

func longestCommonPrefix(key, other []byte, depth int) int {
	maxCmp := min(len(key), len(other))

	idx := depth
	for ; idx < maxCmp; idx++ {
		if key[idx] != other[idx] {
			break
		}
	}

	return idx - depth
}

func prefixMismatch[V any, L nodeLeaf[V]](n nodeRef, key []byte, depth int) int {
	node := n.node()
	maxCmp := min(int(min(maxPrefixLen, node.prefixLen)), len(key)-depth)

	var idx int
	for idx = 0; idx < maxCmp; idx++ {
		if node.prefix[idx] != key[depth+idx] {
			return idx
		}
	}

	if node.prefixLen > maxPrefixLen {
		leaf := (L)(minimum[V](n))
		leafKey := leaf.getTransformKey()

		maxCmp = min(int(len(leafKey)), len(key)) - depth
		for ; idx < maxCmp; idx++ {
			realIdx := depth + idx
			if leafKey[realIdx] != key[realIdx] {
				return idx
			}
		}
	}

	return idx
}

func minimum[V any](ref nodeRef) unsafe.Pointer {
	for ref.pointer != nil {
		kind := ref.tag
		if kind == nodeKindLeaf {
			return ref.pointer
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

func maximum[V any](ref nodeRef) unsafe.Pointer {
	for ref.pointer != nil {
		kind := ref.tag
		if kind == nodeKindLeaf {
			return ref.pointer
		}

		switch kind {
		case nodeKind4:
			n4 := (*node4)(ref.pointer)
			ref = n4.children[n4.childrenLen-1]
		case nodeKind16:
			n16 := (*node16)(ref.pointer)
			ref = n16.children[n16.childrenLen-1]
		case nodeKind48:
			idx := 255
			n48 := (*node48)(ref.pointer)

			for n48.keys[idx] == 0 {
				idx--
			}
			idx = int(n48.keys[idx]) - 1
			ref = n48.children[idx]
		case nodeKind256:
			idx := 255
			n256 := (*node256)(ref.pointer)

			for n256.children[idx].pointer == nil {
				idx--
			}
			ref = n256.children[idx]
		default:
			panic("shouldn't be possible!")
		}
	}

	return nil
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

func all[K nodeKey, V any](root nodeRef, restore func(unsafe.Pointer) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if root.pointer == nil {
			return
		}

		var q []nodeRef

		q = append(q, root)
		for len(q) != 0 {
			n := q[len(q)-1]
			q = q[:len(q)-1]

			if n.tag == nodeKindLeaf {
				k, v := restore(n.pointer)
				if !yield(k, v) {
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

func backward[K nodeKey, V any](root nodeRef, restore func(unsafe.Pointer) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if root.pointer == nil {
			return
		}

		var q []nodeRef

		q = append(q, root)
		for len(q) != 0 {
			n := q[len(q)-1]
			q = q[:len(q)-1]

			if n.tag == nodeKindLeaf {
				k, v := restore(n.pointer)
				if !yield(k, v) {
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

func topK[K nodeKey, V any](t Tree[K, V], k uint) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if k == 0 {
			return
		}

		for key, val := range t.Backward() {
			if k == 0 {
				return
			}

			if !yield(key, val) {
				break
			}

			k--
		}
	}
}

func bottomK[K nodeKey, V any](t Tree[K, V], k uint) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if k == 0 {
			return
		}

		for key, val := range t.All() {
			if k == 0 {
				return
			}

			if !yield(key, val) {
				break
			}

			k--
		}
	}
}

func lowestCommonParent[V any, L nodeLeaf[V]](root nodeRef, prefix []byte) nodeRef {
	var q []nodeRef

	depth := 0
	q = append(q, root)
	for len(q) != 0 {
		n := q[len(q)-1]
		q = q[:len(q)-1]

		idx := prefixMismatch[V, L](n, prefix, depth)
		if idx == 0 { // no match
			continue
		}

		if idx < min(len(prefix)-depth, maxPrefixLen) {
			root = n
			break
		}

		switch n.tag {
		case nodeKind4:
			n4 := (*node4)(n.pointer)

			for i := int(n4.childrenLen) - 1; i >= 0; i-- {
				if n4.children[i].tag == nodeKindLeaf {
					continue
				}
				q = append(q, n4.children[i])
			}

		case nodeKind16:
			n16 := (*node16)(n.pointer)

			for i := int(n16.childrenLen) - 1; i >= 0; i-- {
				if n16.children[i].tag == nodeKindLeaf {
					continue
				}
				q = append(q, n16.children[i])
			}

		case nodeKind48:
			n48 := (*node48)(n.pointer)

			for i := 255; i >= 0; i-- {
				idx := n48.keys[i]
				if idx == 0 || n48.children[i].tag == nodeKindLeaf {
					continue
				}
				q = append(q, n48.children[idx-1])
			}

		case nodeKind256:
			n256 := (*node256)(n.pointer)

			for i := 255; i >= 0; i-- {
				if n256.children[i].pointer == nil || n256.children[i].tag == nodeKindLeaf {
					continue
				}

				q = append(q, n256.children[i])
			}

		default:
			panic("shouldn't be possible!")
		}

		depth += idx + 1
	}

	return root
}

func filter[K nodeKey, V any](root nodeRef, predicate func(K, V) bool, restore func(unsafe.Pointer) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if root.pointer == nil {
			return
		}

		var q []nodeRef

		q = append(q, root)
		for len(q) != 0 {
			n := q[len(q)-1]
			q = q[:len(q)-1]

			if n.tag == nodeKindLeaf {
				k, v := restore(n.pointer)

				if predicate(k, v) {
					if !yield(k, v) {
						return
					}
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

func rangeScan[K nodeKey, V any, L nodeLeaf[V]](
	root nodeRef,
	start, end []byte,
	transformStart, transformEnd []byte,
	restore func(unsafe.Pointer) (K, V),
) iter.Seq2[K, V] {
	idx := longestCommonPrefix(transformStart, transformEnd, 0)

	var search []byte
	if idx != 0 {
		search = transformStart[:idx]
	}

	return func(yield func(K, V) bool) {
		var q []nodeRef

		depth := 0
		q = append(q, root)
		for len(q) != 0 {
			n := q[len(q)-1]
			q = q[:len(q)-1]

			if n.tag == nodeKindLeaf {
				leaf := (L)(n.pointer)

				if bytes.Compare(leaf.getKey(), start) < 0 {
					continue
				}

				if bytes.Compare(leaf.getKey(), end) > 0 {
					break // no need to go further
				}

				k, v := restore(n.pointer)
				if !yield(k, v) {
					return
				}
				continue
			}

			node := n.node()

			if node.prefixLen > 0 && depth < len(search) {
				nodeKey := unsafe.Slice(&node.prefix[0], min(maxPrefixLen, node.prefixLen))
				idx := longestCommonPrefix(nodeKey, search[depth:depth+min(len(search)-depth, maxPrefixLen)], 0)
				if idx == 0 { // no match
					continue
				}
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

			depth += int(node.prefixLen) + 1
		}
	}
}
