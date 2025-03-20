package art

import (
	"iter"
	"strings"
	"unsafe"

	"golang.org/x/text/collate"
)

type Tree[K nodeKey, V any] interface {
	setEnd(byte)
	setCollator(*collate.Collator)

	// Insert inserts a key-value pair in the tree.
	Insert(K, V)

	// Search searches for element with the given key.
	// It returns whether the key is present (bool) and its value if it is present.
	Search(K) (V, bool)

	// Delete deletes a element with the given key.
	Delete(K)

	// Minimum find the minimum K/V pair based on the key.
	Minimum() (K, V, bool)

	// Maximum find the maximum K/V pair based on the key.
	Maximum() (K, V, bool)

	All() iter.Seq2[K, V]
	Backward() iter.Seq2[K, V]
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

func minimum[K nodeKey, V any, L nodeLeaf[K, V]](ref nodeRef) L {
	for ref.pointer != nil {
		kind := ref.tag
		if kind == nodeKindLeaf {
			return (L)(ref.pointer)
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

func maximum[K nodeKey, V any, L nodeLeaf[K, V]](ref nodeRef) L {
	for ref.pointer != nil {
		kind := ref.tag
		if kind == nodeKindLeaf {
			return (L)(ref.pointer)
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

func all[K nodeKey, V any, L nodeLeaf[K, V]](root nodeRef, end byte) iter.Seq2[K, V] {
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
				leaf := (L)(n.pointer)
				keyStr := unsafe.String(leaf.getKey(), leaf.getLen())
				keyStr = strings.Trim(keyStr, string(end))

				if !yield(K(keyStr), leaf.getValue()) {
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

func backward[K nodeKey, V any, L nodeLeaf[K, V]](root nodeRef, end byte) iter.Seq2[K, V] {
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
				leaf := (L)(n.pointer)
				keyStr := unsafe.String(leaf.getKey(), leaf.getLen())
				keyStr = strings.Trim(keyStr, string(end))

				if !yield(K(keyStr), leaf.getValue()) {
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
