package art

import "math/bits"

const (
	lo7BitsMask = uint32(0x1010101) * 0x7F
	hiBitMask   = uint32(0x1010101) * 0x80
)

// searchNode4 tries to find b in keys
// see: http://0x80.pl/notesen/2023-03-06-swar-find-any.html
func searchNode4(keys uint32, b byte) int {
	bitMask := uint32(0x1010101) * uint32(b)
	xor1 := keys ^ bitMask
	isMatch := ((xor1 - 0x1010101) & ^(xor1) & 0x80808080)

	if isMatch != 0 {
		return int((((isMatch - 1) & 0x1010101) * 0x1010101 >> 24) - 1)
	}
	return -1
}

// insertPosNode4 tries to find an element in keys which is smaller than b
// see: https://stackoverflow.com/a/68717720/11269045
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

// shiftLeftClear moves the elements in keys by 1 to the left at pos
//
// e.g.
// shiftLeftClear(0b00000000_00000000_11111111_11111111, 1)
// -> 0b00000000_11111111_00000000_11111111
func shiftLeftClear(keys *uint32, pos int) {
	bitPos := pos << 3
	mask := uint32(0xFFFFFFFF) << bitPos
	backup := *keys & mask // snapshot
	*keys &= ^mask         // clear
	backup <<= 8           // shift
	*keys |= backup        // set
}

// shiftRightClear moves the elements in keys by 1 to the right at pos
//
// shiftRightClear(0b00000000_00000000_11111111_11111111, 1)
// -> 0b00000000_00000000_00000000_11111111
func shiftRightClear(keys *uint32, pos int) {
	bitPos := pos << 3
	mask := uint32(0xFFFFFFFF) << bitPos
	backup := *keys & mask // snapshot
	*keys &= ^(mask >> 8)  // clear
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
