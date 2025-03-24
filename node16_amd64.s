#define rKeys 		R8
#define rChildrenLen 	CL
#define rB 		R10
#define rIdx 		R11
#define rMask		R12
#define rTmp		R13	
#define vBitfield 	X0	
#define vMask 		X1
#define vZeros		X2
#define vTmp		X3	

// func insertPosNode16(keys *[16]byte, childrenLen uint8, b byte) int
TEXT ·insertPosNode16(SB),$0-24
	MOVQ 		keys+0(FP), rKeys
	MOVB 		childrenLen+8(FP), rChildrenLen
	MOVB 		b+9(FP), rB
	PXOR		vZeros, vZeros
	MOVD		$1, rMask

	VMOVDQU 	(rKeys), vBitfield
	MOVD		rB, vMask
	PSHUFB		vZeros, vMask

	MOVB		$0x80, rTmp
	MOVD		rTmp, vTmp
	PSHUFB		vZeros, vTmp
	PXOR		vTmp, vBitfield
	PXOR		vTmp, vMask

	PCMPGTB		vMask, vBitfield
	PMOVMSKB 	vBitfield, rIdx

	SALW		rChildrenLen, rMask
	SUBW		$1, rMask

	ANDW		rMask, rIdx

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
	PXOR		vZeros, vZeros
	MOVD		$1, rMask

	VMOVDQU 	(rKeys), vBitfield
	MOVD		rB, vMask
	PSHUFB		vZeros, vMask

	PCMPEQB		vMask, vBitfield
	PMOVMSKB 	vBitfield, rIdx

	SALW		rChildrenLen, rMask
	SUBW		$1, rMask

	ANDW		rMask, rIdx

	CMPW		rIdx, $0
	JEQ		not_found
 
	TZCNTW		rIdx, rIdx
	MOVD		rIdx, ret+16(FP)
	RET
	
not_found:
	MOVD 		$-1, rIdx
	MOVD		rIdx, ret+16(FP)
	RET
