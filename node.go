package art

import (
	"unsafe"
)

const (
	minNode4     = uint8(2)
	maxNode4     = uint8(4)
	minNode16    = maxNode4 + 1
	maxNode16    = uint8(16)
	minNode48    = maxNode16 + 1
	maxNode48    = uint8(48)
	minNode256   = maxNode48 + 1
	maxNode256   = 256
	maxPrefixLen = 10
)

//go:generate stringer -type=nodeKind -linecomment
type nodeKind uint8

const (
	nodeKind4    nodeKind = iota // NODE_4
	nodeKind16                   // NODE_16
	nodeKind48                   // NODE_48
	nodeKind256                  // NODE_256
	nodeKindLeaf                 // NODE_LEAF
)

type node struct {
	prefixLen   uint32
	childrenLen uint8
	prefix      [maxPrefixLen]byte
}

type chars interface {
	string | []byte
}

type ints interface {
	int | int64 | int32 | int16 | int8
}

type uints interface {
	uint | uint64 | uint32 | uint16 | uint8
}

type floats interface {
	float64 | float32
}

type nodeKey interface {
	chars | uints | ints | floats | any
}

type nodeLeaf[K nodeKey, V any] interface {
	getKey() []byte
	getTransformKey() []byte
	getValue() V
	setValue(V)

	*alphaLeafNode[K, V] |
		*collateLeafNode[K, V] |
		*unsignedLeafNode[K, V] |
		*signedLeafNode[K, V] |
		*floatLeafNode[K, V] |
		*compoundLeafNode[K, V]
}

type nodeRef struct {
	pointer unsafe.Pointer
	tag     nodeKind
}

func (ref *nodeRef) node() *node {
	return (*node)(ref.pointer)
}

func (ref *nodeRef) findChild(b byte) *nodeRef {
	switch ref.tag {
	case nodeKind4:
		n4 := (*node4)(ref.pointer)

		if i := searchNode4(n4.keys, b); i != -1 && i < int(n4.childrenLen) {
			return &n4.children[i]
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
	switch ptr.tag {
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

func (ptr *nodeRef) deleteChild(b byte) {
	switch ptr.tag {
	case nodeKind4:
		n4 := (*node4)(ptr.pointer)
		n4.deleteChild(ptr, b)

	case nodeKind16:
		n16 := (*node16)(ptr.pointer)
		n16.deleteChild(ptr, b)

	case nodeKind48:
		n48 := (*node48)(ptr.pointer)
		n48.deleteChild(ptr, b)

	case nodeKind256:
		n256 := (*node256)(ptr.pointer)
		n256.deleteChild(ptr, b)

	default:
		panic("shouldn't be possible!")
	}
}

type node4 struct {
	node
	children [maxNode4]nodeRef
	keys     uint32
}

func (n4 *node4) clear() {
	clear(n4.children[:])
	n4.node = node{}
	n4.keys = 0
}

func (n4 *node4) addChild(ref *nodeRef, b byte, child nodeRef) {
	if n4.childrenLen < maxNode4 {
		var idx int

		if i := insertPosNode4(n4.keys, b); i != -1 {
			idx = i
			loLimit := idx + 1
			shiftLeftClear(&n4.keys, idx)
			copy(n4.children[loLimit:], n4.children[idx:])
		} else {
			idx = int(n4.childrenLen)
		}

		setAtPos(&n4.keys, idx, b)
		n4.children[idx] = child
		n4.childrenLen++
	} else {
		n16 := nodePools[nodeKind16].Get().(*node16)

		copy(n16.keys[:], deconstruct(n4.keys)[:])
		copy(n16.children[:], n4.children[:])

		n16.childrenLen = n4.childrenLen
		n16.prefixLen = n4.prefixLen
		copy(n16.prefix[:], n4.prefix[:])

		*ref = nodeRef{pointer: unsafe.Pointer(n16), tag: nodeKind16}
		n16.addChild(ref, b, child)

		n4.clear()
		nodePools[nodeKind4].Put(n4)
	}
}

func (n4 *node4) deleteChild(ref *nodeRef, b byte) {
	if i := searchNode4(n4.keys, b); i != -1 {
		shiftRightClear(&n4.keys, i+1)
		copy(n4.children[i:], n4.children[i+1:])
		n4.childrenLen--
	}

	if n4.childrenLen == 1 {
		child := n4.children[0]

		if child.tag != nodeKindLeaf {
			prefix := n4.prefixLen
			childNode := child.node()

			if prefix < maxPrefixLen {
				n4.prefix[prefix] = getAtPos(n4.keys, 0)
				prefix++
			}

			if prefix < maxPrefixLen {
				subPrefix := min(childNode.prefixLen, maxPrefixLen-prefix)
				copy(n4.prefix[prefix:], childNode.prefix[:])
				prefix += subPrefix
			}

			hiLimit := min(maxPrefixLen, prefix)
			copy(childNode.prefix[:], n4.prefix[:hiLimit])
			childNode.prefixLen += n4.prefixLen + 1
		}
		*ref = child

		n4.clear()
		nodePools[nodeKind4].Put(n4)
	}
}

type node16 struct {
	node
	children [maxNode16]nodeRef
	keys     [maxNode16]byte
}

func (n16 *node16) clear() {
	clear(n16.children[:])
	n16.node = node{}
	clear(n16.keys[:])
}

func (n16 *node16) addChild(ref *nodeRef, b byte, child nodeRef) {
	if n16.childrenLen < maxNode16 {
		idx := insertPosNode16(&n16.keys, n16.childrenLen, b)

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
		n48 := nodePools[nodeKind48].Get().(*node48)

		copy(n48.children[:n16.childrenLen], n16.children[:])
		for i := uint8(0); i < n16.childrenLen; i++ {
			n48.keys[n16.keys[i]] = i + 1
		}

		n48.childrenLen = n16.childrenLen
		n48.prefixLen = n16.prefixLen
		copy(n48.prefix[:], n16.prefix[:])

		*ref = nodeRef{pointer: unsafe.Pointer(n48), tag: nodeKind48}
		n48.addChild(ref, b, child)

		n16.clear()
		nodePools[nodeKind16].Put(n16)
	}
}

func (n16 *node16) deleteChild(ref *nodeRef, b byte) {
	pos := searchNode16(&n16.keys, n16.childrenLen, b)

	copy(n16.keys[pos:], n16.keys[pos+1:])
	copy(n16.children[pos:], n16.children[pos+1:])
	n16.childrenLen--

	if n16.childrenLen == 3 {
		n4 := nodePools[nodeKind4].Get().(*node4)
		*ref = nodeRef{
			pointer: unsafe.Pointer(n4),
			tag:     nodeKind4,
		}

		n4.childrenLen = n16.childrenLen
		n4.prefixLen = n16.prefixLen
		copy(n4.prefix[:], n16.prefix[:])

		n4.keys = construct(n16.keys[0], n16.keys[1], n16.keys[2], n16.keys[3])
		copy(n4.children[:], n16.children[:])

		n16.clear()
		nodePools[nodeKind16].Put(n16)
	}
}

type node48 struct {
	node
	children [maxNode48]nodeRef
	keys     [256]byte
}

func (n48 *node48) clear() {
	clear(n48.children[:])
	n48.node = node{}
	clear(n48.keys[:])
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
		n256 := nodePools[nodeKind256].Get().(*node256)

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

		n48.clear()
		nodePools[nodeKind48].Put(n48)
	}
}

func (n48 *node48) deleteChild(ref *nodeRef, b byte) {
	pos := n48.keys[b]
	n48.keys[b] = 0
	n48.children[pos-1].pointer = nil
	n48.childrenLen--

	if n48.childrenLen == 12 {
		n16 := nodePools[nodeKind16].Get().(*node16)
		*ref = nodeRef{
			pointer: unsafe.Pointer(n16),
			tag:     nodeKind16,
		}

		n16.childrenLen = n48.childrenLen
		n16.prefixLen = n48.prefixLen
		copy(n16.prefix[:], n48.prefix[:])

		children := 0
		for i := 0; i < 256; i++ {
			pos = n48.keys[i]
			if pos != 0 {
				n16.keys[children] = uint8(i)
				n16.children[children] = n48.children[pos-1]
				children++
			}
		}

		n48.clear()
		nodePools[nodeKind48].Put(n48)
	}
}

type node256 struct {
	node
	children [maxNode256]nodeRef
}

func (n256 *node256) clear() {
	clear(n256.children[:])
	n256.node = node{}
}

func (n256 *node256) addChild(b byte, child nodeRef) {
	n256.childrenLen++
	n256.children[b] = child
}

func (n256 *node256) deleteChild(ref *nodeRef, b byte) {
	n256.children[b].pointer = nil
	n256.childrenLen--

	if n256.childrenLen == 37 {
		n48 := nodePools[nodeKind48].Get().(*node48)
		*ref = nodeRef{
			pointer: unsafe.Pointer(n48),
			tag:     nodeKind48,
		}

		n48.childrenLen = n256.childrenLen
		n48.prefixLen = n256.prefixLen
		copy(n48.prefix[:], n256.prefix[:])

		pos := 0
		for i := 0; i < 256; i++ {
			if n256.children[i].pointer != nil {
				n48.children[pos] = n256.children[i]
				n48.keys[i] = uint8(pos + 1)
				pos++
			}
		}

		n256.clear()
		nodePools[nodeKind256].Put(n256)
	}
}
