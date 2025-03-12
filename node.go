package art

import (
	"unsafe"
)

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

type node4 struct {
	node
	keys     [maxNode4]byte
	children [maxNode4]taggedPointer
}

func (n4 *node4) addChild(ref *taggedPointer, b byte, child taggedPointer) *node16 {
	if n4.childrenLen < maxNode4 {
		var idx uint32

		for idx = range uint32(n4.childrenLen) {
			if b < n4.keys[idx] {
				break
			}
		}

		loLimit := idx + 1
		hiLimit := idx + uint32(n4.childrenLen)
		copy(n4.keys[loLimit:], n4.keys[idx:hiLimit])
		copy(n4.children[loLimit:], n4.children[idx:hiLimit])

		n4.keys[idx] = b
		n4.children[idx] = child
		n4.childrenLen++
	} else {
		n16 := new(node16)

		copy(n16.keys[:], n4.keys[:])
		copy(n16.children[:], n4.children[:])

		n16.childrenLen = n4.childrenLen
		n16.prefixLen = n4.prefixLen
		copy(n16.prefix[:], n4.prefix[:])

		*ref = taggedPointerPack(unsafe.Pointer(n16), uintptr(nodeKind16))
		n16.addChild(ref, b, child)
		return n16
	}
	return nil
}

type node16 struct {
	node
	keys     [maxNode16]byte
	children [maxNode16]taggedPointer
}

func (n16 *node16) addChild(ref *taggedPointer, b byte, child taggedPointer) *node48 {
	if n16.childrenLen < maxNode16 {
		idx := searchNode16(&n16.keys, n16.childrenLen, b)

		if idx != -1 {
			loLimit := idx + 1
			hiLimit := idx + int(n16.childrenLen)
			copy(n16.keys[loLimit:], n16.keys[idx:hiLimit])
			copy(n16.children[loLimit:], n16.children[idx:hiLimit])
		} else {
			idx = int(n16.childrenLen)
		}

		n16.keys[idx] = b
		n16.children[idx] = child
		n16.childrenLen++
	} else {
		n48 := new(node48)

		copy(n48.children[:], n16.children[:])
		for i := uint8(0); i < n16.childrenLen; i++ {
			n48.keys[n16.keys[i]] = i + 1
		}

		n48.childrenLen = n16.childrenLen
		n48.prefixLen = n16.prefixLen
		copy(n48.prefix[:], n16.prefix[:])

		*ref = taggedPointerPack(unsafe.Pointer(n48), uintptr(nodeKind48))
		n48.addChild(ref, b, child)
		return n48
	}
	return nil
}

type node48 struct {
	node
	keys     [256]byte
	children [maxNode48]taggedPointer
}

func (n48 *node48) addChild(ref *taggedPointer, b byte, child taggedPointer) *node256 {
	if n48.childrenLen < maxNode48 {
		pos := uint8(0)
		for n48.children[pos].pointer() != nil {
			pos++
		}

		n48.children[pos] = child
		n48.keys[b] = pos + 1
		n48.childrenLen++
	} else {
		n256 := new(node256)

		for i := range maxNode256 {
			if n48.keys[i] != 0 {
				n256.children[i] = n48.children[n48.keys[i]-1]
			}
		}

		n256.childrenLen = n48.childrenLen
		n256.prefixLen = n48.prefixLen
		copy(n256.prefix[:], n48.prefix[:])

		*ref = taggedPointerPack(unsafe.Pointer(n256), uintptr(nodeKind256))
		n256.addChild(b, child)
		return n256
	}

	return nil
}

type node256 struct {
	node
	children [maxNode256]taggedPointer
}

type nodeLeaf[K nodeKey, V any] struct {
	key   *byte
	value V
	len   uint32
}

func (n256 *node256) addChild(b byte, child taggedPointer) {
	if n256.childrenLen == 255 {
		panic("cannot grow anymore...")
	}

	n256.childrenLen++
	n256.children[b] = child
}
