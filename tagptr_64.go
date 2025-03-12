// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 || arm64 || loong64 || mips64 || mips64le || ppc64 || ppc64le || riscv64 || s390x || wasm

package art

import (
	"runtime"
	"unsafe"
)

const (
	addrBits = 48
	tagBits  = 64 - addrBits + 3
)

const taggedPointerBits = tagBits

// taggedPointerPack created a taggedPointer from a pointer and a tag.
// Tag bits that don't fit in the result are discarded.
func taggedPointerPack(ptr unsafe.Pointer, tag uintptr) taggedPointer {
	return taggedPointer(uint64(uintptr(ptr))<<(64-addrBits) | uint64(tag&(1<<tagBits-1)))
}

// Pointer returns the pointer from a taggedPointer.
func (tp taggedPointer) pointer() unsafe.Pointer {
	if runtime.GOARCH == "amd64" {
		// amd64 systems can place the stack above the VA hole, so we need to sign extend
		// val before unpacking.
		return unsafe.Pointer(uintptr(int64(tp) >> tagBits << 3))
	}
	return unsafe.Pointer(uintptr(tp >> tagBits << 3))
}

// Tag returns the tag from a taggedPointer.
func (tp taggedPointer) tag() uintptr {
	return uintptr(tp & (1<<taggedPointerBits - 1))
}
