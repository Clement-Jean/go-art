#define rKeys 		R8
#define rChildrenLen 	R9
#define rB 		R10
#define rIdx 		R11
#define rMask		R12
#define vBitfield 	X0
#define vCmp		X1	
#define vMask 		X2
#define zeroes		X3

// func insertPosNode16(keys *[16]byte, childrenLen uint8, b byte) int
TEXT ·insertPosNode16(SB),$0-24
	MOVQ 		keys+0(FP), rKeys
	MOVB 		childrenLen+8(FP), rChildrenLen
	MOVB 		b+9(FP), rB
	PXOR		vBitfield, vBitfield
	PXOR		zeroes, zeroes
	MOVD		$0, rIdx
	MOVD		$0, rMask

	VMOVDQU 	(rKeys), vBitfield
	MOVD		rB, vMask
	PSHUFB		zeroes, vMask

	PCMPGTB		vMask, vBitfield
	PMOVMSKB 	vBitfield, rIdx

	SHRL		$1, rChildrenLen, rMask
	SUBQ		$1, rMask

	ANDQ		rMask, rIdx

	CMPW		rIdx, $0
	JEQ		not_found

	TZCNTW		rIdx, rIdx
	MOVD		rIdx, ret+16(FP)
	RET
	
not_found:
	MOVD 		$-1, rIdx
	MOVD		rIdx, ret+16(FP)
	RET
	
// func searchNode16(keys *[16]byte, childrenLen uint8, b byte) int
TEXT ·searchNode16(SB),$0-24
	MOVQ 		keys+0(FP), rKeys
	MOVB 		childrenLen+8(FP), rChildrenLen
	MOVB 		b+9(FP), rB
	PXOR		vBitfield, vBitfield
	PXOR		zeroes, zeroes
	MOVD		$0, rIdx
	MOVD		$0, rMask

	VMOVDQU 	(rKeys), vBitfield
	MOVD		rB, vMask
	PSHUFB		zeroes, vMask

	PCMPEQB		vBitfield, vMask
	PMOVMSKB 	vMask, rIdx

	SHRL		$1, rChildrenLen, rMask
	SUBQ		$1, rMask

	ANDQ		rMask, rIdx

	CMPW		rIdx, $0
	JEQ		not_found

	TZCNTW		rIdx, rIdx
	MOVD		rIdx, ret+16(FP)
	RET
	
not_found:
	MOVD 		$-1, rIdx
	MOVD		rIdx, ret+16(FP)
	RET
