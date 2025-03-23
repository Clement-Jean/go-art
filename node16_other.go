//go:build !amd64 && !arm64

package art

import "math/bits"

func insertPosNode16(keys *[16]byte, childrenLen uint8, b byte) int {
	bitfield := 0
	for i := range 16 {
		if keys[i] > b {
			bitfield |= (1 << i)
		}
	}

	mask := (1 << childrenLen) - 1
	bitfield &= mask

	if bitfield != 0 {
		return bits.TrailingZeros(uint(bitfield))
	}
	return -1
}

func searchNode16(keys *[16]byte, childrenLen uint8, b byte) int {
	bitfield := 0
	for i := range 16 {
		if keys[i] == b {
			bitfield |= (1 << i)
		}
	}

	mask := (1 << childrenLen) - 1
	bitfield &= mask

	if bitfield != 0 {
		return bits.TrailingZeros(uint(bitfield))
	}
	return -1
}
