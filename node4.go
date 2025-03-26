package art

import "math/bits"

const (
	lo7BitsMask = uint32(0x1010101) * 0x7F
	hiBitMask   = uint32(0x1010101) * 0x80
)

func searchNode4(keys uint32, b byte) int {
	bitMask := uint32(0x1010101) * uint32(b)
	xor1 := keys ^ bitMask

	if isMatch := ((xor1 - 0x1010101) & ^(xor1) & 0x80808080); isMatch != 0 {
		return int((((isMatch - 1) & 0x1010101) * 0x1010101 >> 24) - 1)
	}
	return -1
}

func insertPosNode4(keys uint32, b byte) int {
	bitMask := uint32(0x1010101) * uint32(b) // broadcast
	t0 := (((keys | hiBitMask) - (bitMask & ^hiBitMask)) | (keys ^ bitMask)) ^ (keys | ^bitMask)
	t1 := t0 & hiBitMask
	t2 := t1 + t1 - (t1 >> 7)
	t2 = ^t2

	if t2 != 0 {
		return bits.TrailingZeros32(t2) >> 3
	}
	return -1
}

func getAtPos(keys uint32, pos int) byte {
	return byte(keys>>(pos<<3)) & 0xFF
}

func setAtPos(keys *uint32, pos int, b byte) {
	bitPos := pos << 3
	*keys &= ^(0xFF << bitPos)     // clear
	*keys |= (uint32(b) << bitPos) // set
}

func shiftLeftClear(keys *uint32, pos int) {
	bitPos := pos << 3
	mask := uint32(0xFFFFFFFF) << bitPos
	backup := *keys & mask // snapshot
	*keys &= ^mask         // clear
	backup <<= 8           // shift
	*keys |= backup        // set
}

func shiftRightClear(keys *uint32, pos int) {
	bitPos := pos << 3
	mask := uint32(0xFFFFFFFF) << bitPos
	mask2 := uint32(0xFFFFFFFF) << (bitPos - 8)
	backup := *keys & mask // snapshot
	*keys &= ^mask2        // clear
	backup >>= 8           // shift
	*keys |= backup        // set
}

func construct(a, b, c, d byte) uint32 {
	return uint32(d)<<24 | uint32(c)<<16 | uint32(b)<<8 | uint32(a)
}

func deconstruct(keys uint32) []byte {
	return []byte{
		byte(keys) & 0xFF,
		byte((keys >> 8) & 0xFF),
		byte((keys >> 16) & 0xFF),
		byte((keys >> 24) & 0xFF),
	}
}
