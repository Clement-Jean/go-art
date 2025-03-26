package art

import (
	"encoding/binary"
	"math"
	"math/bits"
	"unsafe"

	"golang.org/x/text/collate"
)

type BinaryComparableKey[K nodeKey] interface {
	Transform(K) ([]byte, []byte)
	Restore([]byte) K
}

type AlphabeticalOrderKey[K chars] struct{}

func (aok AlphabeticalOrderKey[K]) Transform(k K) ([]byte, []byte) {
	b := []byte(k)
	return b, b
}
func (aok AlphabeticalOrderKey[K]) Restore(b []byte) K { return K(b) }

var _ BinaryComparableKey[[]byte] = AlphabeticalOrderKey[[]byte]{}

type CollationOrderKey[K chars | []rune] struct {
	c   *collate.Collator
	buf *collate.Buffer
	src K
}

func (cok *CollationOrderKey[K]) Transform(k K) ([]byte, []byte) {
	cok.src = k
	b := []byte(string(k))
	return b, cok.c.Key(cok.buf, b)
}
func (cok *CollationOrderKey[K]) Restore(b []byte) K { return cok.src }

var _ BinaryComparableKey[[]rune] = &CollationOrderKey[[]rune]{}

type UnsignedBinaryKey[K uints] struct{}

func (ubk UnsignedBinaryKey[K]) Transform(k K) ([]byte, []byte) {
	var b []byte

	switch any(k).(type) {
	case uint8:
		b = []byte{uint8(k)}
	case uint16:
		b = make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(k))
	case uint32:
		b = make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(k))
	case uint64:
		b = make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(k))
	case uint:
		if bits.UintSize == 32 {
			b = make([]byte, 4)
			binary.BigEndian.PutUint32(b, uint32(k))
		} else {
			b = make([]byte, 8)
			binary.BigEndian.PutUint64(b, uint64(k))
		}
	default:
		panic("shouldn't be possible!")
	}
	return b, b
}
func (ubk UnsignedBinaryKey[K]) Restore(b []byte) K {
	var k K
	switch any(k).(type) {
	case uint8:
		return K((uint8)(b[0]))
	case uint16:
		return K(binary.BigEndian.Uint16(b))
	case uint32:
		return K(binary.BigEndian.Uint32(b))
	case uint64:
		return K(binary.BigEndian.Uint64(b))
	case uint:
		if bits.UintSize == 32 {
			return K(binary.BigEndian.Uint32(b))
		} else {
			return K(binary.BigEndian.Uint64(b))
		}
	default:
		panic("shouldn't be possible!")
	}
}

var _ BinaryComparableKey[uint] = UnsignedBinaryKey[uint]{}

type SignedBinaryKey[K ints] struct{}

func (sbk SignedBinaryKey[K]) Transform(k K) ([]byte, []byte) {
	var b []byte

	switch any(k).(type) {
	case int8:
		b = []byte{(*(*uint8)(unsafe.Pointer(&k))) ^ 0x80}
	case int16:
		b = make([]byte, 2)
		binary.BigEndian.PutUint16(b, (*(*uint16)(unsafe.Pointer(&k)))^0x8000)
	case int32:
		b = make([]byte, 4)
		binary.BigEndian.PutUint32(b, (*(*uint32)(unsafe.Pointer(&k)))^0x80000000)
	case int64:
		b = make([]byte, 8)
		binary.BigEndian.PutUint64(b, (*(*uint64)(unsafe.Pointer(&k)))^0x8000000000000000)
	case int:
		if bits.UintSize == 32 {
			b = make([]byte, 4)
			binary.BigEndian.PutUint32(b, (*(*uint32)(unsafe.Pointer(&k)))^0x80000000)
		} else {
			b = make([]byte, 8)
			binary.BigEndian.PutUint64(b, (*(*uint64)(unsafe.Pointer(&k)))^0x8000000000000000)
		}
	default:
		panic("shouldn't be possible!")
	}
	return b, b
}
func (sbk SignedBinaryKey[K]) Restore(b []byte) K {
	var k K
	switch any(k).(type) {
	case int8:
		i := b[0] ^ 0x80
		return K(*(*int8)(unsafe.Pointer(&i)))
	case int16:
		i := binary.BigEndian.Uint16(b) ^ 0x8000
		return K(*(*int16)(unsafe.Pointer(&i)))
	case int32:
		i := binary.BigEndian.Uint32(b) ^ 0x80000000
		return K(*(*int32)(unsafe.Pointer(&i)))
	case int64:
		i := binary.BigEndian.Uint64(b) ^ 0x8000000000000000
		return K(*(*int64)(unsafe.Pointer(&i)))
	case int:
		if bits.UintSize == 32 {
			i := binary.BigEndian.Uint32(b) ^ 0x80000000
			return K(*(*int32)(unsafe.Pointer(&i)))
		} else {
			i := binary.BigEndian.Uint64(b) ^ 0x8000000000000000
			return K(*(*int64)(unsafe.Pointer(&i)))
		}
	default:
		panic("shouldn't be possible!")
	}
}

var _ BinaryComparableKey[int] = SignedBinaryKey[int]{}

type FloatBinaryKey[K floats] struct{}

func (fbk FloatBinaryKey[K]) Transform(k K) ([]byte, []byte) {
	var b []byte

	switch any(k).(type) {
	case float32:
		var i uint32
		f64 := float64(k)

		if math.IsInf(f64, 1) {
			i = math.MaxUint32 - 1
		} else if math.IsInf(f64, -1) {
			i = 1
		} else if math.IsNaN(f64) {
			i = 0
		} else {
			// http://stereopsis.com/radix.html

			i = *(*uint32)(unsafe.Pointer(&k))

			t := i >> 31
			mask := -(*(*int32)(unsafe.Pointer(&t)))
			mask2 := *(*uint32)(unsafe.Pointer(&mask)) | 0x80000000
			i ^= mask2
			i += 2
		}

		b = make([]byte, 4)
		binary.BigEndian.PutUint32(b, i)

	case float64:
		var i uint64
		f64 := float64(k)

		if math.IsInf(f64, 1) {
			i = math.MaxUint64 - 1
		} else if math.IsInf(f64, -1) {
			i = 1
		} else if math.IsNaN(f64) {
			i = 0
		} else {
			// http://stereopsis.com/radix.html

			i = *(*uint64)(unsafe.Pointer(&f64))

			t := i >> 63
			mask := -(*(*int64)(unsafe.Pointer(&t)))
			mask2 := *(*uint64)(unsafe.Pointer(&mask)) | 0x8000000000000000
			i ^= mask2
			i += 2
		}

		b = make([]byte, 8)
		binary.BigEndian.PutUint64(b, i)
	}
	return b, b
}
func (fbk FloatBinaryKey[K]) Restore(b []byte) K {
	var k K

	switch any(k).(type) {
	case float32:
		i := binary.BigEndian.Uint32(b)
		if i == math.MaxUint32-1 {
			return K(math.Inf(1))
		} else if i == 1 {
			return K(math.Inf(-1))
		} else if i == 0 {
			return K(math.NaN())
		} else if i == 2 {
			return K(0)
		}

		i -= 2
		mask := ((i >> 31) - 1) | 0x80000000
		i ^= mask

		return *(*K)(unsafe.Pointer(&i))

	case float64:
		i := binary.BigEndian.Uint64(b)

		if i == math.MaxUint64-1 {
			return K(math.Inf(1))
		} else if i == 1 {
			return K(math.Inf(-1))
		} else if i == 0 {
			return K(math.NaN())
		} else if i == 2 {
			return K(0)
		}

		i -= 2
		mask := ((i >> 63) - 1) | 0x8000000000000000
		i ^= mask

		return *(*K)(unsafe.Pointer(&i))
	default:
		panic("shouldn't be possible!")
	}
}

var _ BinaryComparableKey[float64] = FloatBinaryKey[float64]{}
