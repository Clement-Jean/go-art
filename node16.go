//go:build arm64

package art

func insertPosNode16(keys *[16]byte, childrenLen uint8, b byte) int

func searchNode16(keys *[16]byte, childrenLen uint8, b byte) int
