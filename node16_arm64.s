// func searchNode16(keys *[16]byte, childrenLen uint8, b byte) int
TEXT ·searchNode16(SB),$0-24
#define rKeys 		R0
#define rChildrenLen 	R1
#define rB 		R2
#define rIdx 		R3
#define vBitfield 	V0
#define vMask 		V1
	
	MOVD 	keys+0(FP), rKeys
	MOVB 	childrenLen+8(FP), rChildrenLen
	MOVB 	b+9(FP), rB

	VLD1 	(rKeys), [vBitfield.B16]
	VDUP 	rB, vMask.B16

	VCMEQ	vBitfield.B16, vMask.B16, vBitfield.B16

	WORD	$0x0f0c8400 // shrn.8b	v0, v0, #0x4
	VMOV 	vBitfield.D[0], rIdx

	CBZ	rIdx, not_found
	AND	$0x8888888888888888, rIdx, rIdx
	RBIT	rIdx, rIdx
	CLZ	rIdx, rIdx
	ASR	$2, rIdx

	MOVD	rIdx, ret+16(FP)
	RET

not_found:
	MOVD 	$-1, rIdx
	MOVD	rIdx, ret+16(FP)
	RET
