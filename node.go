package art

import "unsafe"

const (
	minNode4     = uint8(2)
	maxNode4     = uint8(4)
	minNode16    = uint8(maxNode4 + 1)
	maxNode16    = uint8(16)
	minNode48    = uint8(maxNode16 + 1)
	maxNode48    = uint8(48)
	minNode256   = uint8(maxNode48 + 1)
	maxNode256   = 256
	maxPrefixLen = 10
)

//go:generate stringer -type=nodeKind -linecomment
type nodeKind uint8

const (
	nodeKindUndefined nodeKind = iota // UNDEFINED
	nodeKindLeaf                      // NODE_LEAF
	nodeKind4                         // NODE_4
	nodeKind16                        // NODE_16
	nodeKind48                        // NODE_48
	nodeKind256                       // NODE_256
)

type node struct {
	prefixLen   uint32
	childrenLen uint8
	prefix      [maxPrefixLen]byte
}

type nodeKey interface {
	string | []byte | []rune
}

type nodeRef struct {
	pointer unsafe.Pointer
	tag     nodeKind
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

type node4 struct {
	children [maxNode4]nodeRef
	node
	keys [maxNode4]byte
}

func (n4 *node4) addChild(ref *nodeRef, b byte, child nodeRef) {
	if n4.childrenLen < maxNode4 {
		var idx uint32

		for idx = 0; idx < uint32(n4.childrenLen); idx++ {
			if b < n4.keys[idx] {
				break
			}
		}

		loLimit := idx + 1
		copy(n4.keys[loLimit:], n4.keys[idx:])
		copy(n4.children[loLimit:], n4.children[idx:])

		n4.keys[idx] = b
		n4.children[idx] = child
		n4.childrenLen++
	} else {
		n16 := new(node16)

		copy(n16.keys[:n4.childrenLen], n4.keys[:])
		copy(n16.children[:n4.childrenLen], n4.children[:])

		n16.childrenLen = n4.childrenLen
		n16.prefixLen = n4.prefixLen
		copy(n16.prefix[:], n4.prefix[:])

		*ref = nodeRef{pointer: unsafe.Pointer(n16), tag: nodeKind16}
		n16.addChild(ref, b, child)
	}
}

type node16 struct {
	children [maxNode16]nodeRef
	node
	keys [maxNode16]byte
}

func (n16 *node16) addChild(ref *nodeRef, b byte, child nodeRef) {
	if n16.childrenLen < maxNode16 {
		idx := searchNode16(&n16.keys, n16.childrenLen, b)

		if idx != -1 {
			loLimit := idx + 1
			copy(n16.keys[loLimit:], n16.keys[idx:])
			copy(n16.children[loLimit:], n16.children[idx:])
		} else {
			idx = int(n16.childrenLen)
		}

		n16.keys[idx] = b
		n16.children[idx] = child
		n16.childrenLen++
	} else {
		n48 := new(node48)

		copy(n48.children[:n16.childrenLen], n16.children[:])
		for i := uint8(0); i < n16.childrenLen; i++ {
			n48.keys[n16.keys[i]] = i + 1
		}

		n48.childrenLen = n16.childrenLen
		n48.prefixLen = n16.prefixLen
		copy(n48.prefix[:], n16.prefix[:])

		*ref = nodeRef{pointer: unsafe.Pointer(n48), tag: nodeKind48}
		n48.addChild(ref, b, child)
	}
}

type node48 struct {
	children [maxNode48]nodeRef
	node
	keys [256]byte
}

func (n48 *node48) addChild(ref *nodeRef, b byte, child nodeRef) {
	if n48.childrenLen < maxNode48 {
		pos := uint8(0)
		for n48.children[pos].pointer != nil {
			pos++
		}

		n48.children[pos] = child
		n48.keys[b] = pos + 1
		n48.childrenLen++
	} else {
		n256 := new(node256)

		for i := 0; i < maxNode256; i++ {
			if n48.keys[i] != 0 {
				n256.children[i] = n48.children[n48.keys[i]-1]
			}
		}

		n256.childrenLen = n48.childrenLen
		n256.prefixLen = n48.prefixLen
		copy(n256.prefix[:], n48.prefix[:])

		*ref = nodeRef{pointer: unsafe.Pointer(n256), tag: nodeKind256}
		n256.addChild(b, child)
	}
}

type node256 struct {
	children [maxNode256]nodeRef
	node
}

type nodeLeaf[K nodeKey, V any] struct {
	key   *byte
	value V
	len   uint32
}

func (n256 *node256) addChild(b byte, child nodeRef) {
	n256.childrenLen++
	n256.children[b] = child
}
