//go:build amd64 || arm64

package art

func searchNode16(keys *[16]byte, childrenLen uint8, b byte) int
