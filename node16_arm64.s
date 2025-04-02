#define rKeys 		R0
#define rB 		R1
#define rIdx 		R2
#define vBitfield 	V0
#define vMask 		V1

// func insertPosNode16(keys *[16]byte, childrenLen uint8, b byte)
TEXT ·insertPosNode16(SB),$0-24
	MOVD 	keys+0(FP), rKeys
	MOVB 	b+9(FP), rB

	VLD1 	(rKeys), [vBitfield.B16]
	VDUP 	rB, vMask.B16

	WORD	$0x6e213400 // cmhi.16b	v0, v0, v1

	WORD	$0x0f0c8400 // shrn.8b	v0, v0, #0x4
	FMOVD	F0, rIdx
	//VMOV 	vBitfield.D[0], rIdx

	CBNZ	rIdx, found
	MOVD 	$-1, rIdx
	MOVD	rIdx, ret+16(FP)
	RET

found:
	AND	$0x8888888888888888, rIdx, rIdx
	RBIT	rIdx, rIdx
	CLZ	rIdx, rIdx
	ASR	$2, rIdx

	MOVD	rIdx, ret+16(FP)
	RET	
	

// func searchNode16(keys *[16]byte, childrenLen uint8, b byte) int
TEXT ·searchNode16(SB),$0-24	
	MOVD 	keys+0(FP), rKeys
	MOVB 	b+9(FP), rB

	VLD1 	(rKeys), [vBitfield.B16]
	VDUP 	rB, vMask.B16

	VCMEQ	vBitfield.B16, vMask.B16, vBitfield.B16

	WORD	$0x0f0c8400 // shrn.8b	v0, v0, #0x4
	FMOVD	F0, rIdx
	//VMOV 	vBitfield.D[0], rIdx

	CBNZ	rIdx, found
	MOVD 	$-1, rIdx
	MOVD	rIdx, ret+16(FP)
	RET

found:
	AND	$0x8888888888888888, rIdx, rIdx
	RBIT	rIdx, rIdx
	CLZ	rIdx, rIdx
	ASR	$2, rIdx

	MOVD	rIdx, ret+16(FP)
	RET
