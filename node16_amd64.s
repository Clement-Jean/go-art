// func searchNode16(keys *[16]byte, childrenLen uint8, b byte) int
TEXT ·searchNode16(SB),$0-24
#define rKeys 		R0
#define rChildrenLen 	R1
#define rB 		R2
#define rIdx 		R3
#define vBitfield 	X0
#define vMask 		X1
	
	MOVQ 		keys+0(FP), rKeys
	MOVB 		childrenLen+8(FP), rChildrenLen
	MOVB 		b+9(FP), rB

	VMOVDQA 	(data), vBitfield
	VMOVDQA 	rB, vMask // vpbroadcastb if AVX516

	PCMPEQB		vBitefield, vMask, vBitfield
	PMOVMSKB 	vBitfield, rIdx
	
not_found:
	MOVD 		$-1, rIdx
	MOVD		rIdx, ret+16(FP)
	RET
