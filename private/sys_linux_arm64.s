#include "textflag.h"


TEXT ·sys_mprotect(SB),NOSPLIT,$0
	MOVD	addr+0(FP), R0
	MOVD	n+8(FP), R1
	MOVD	prot+16(FP), R2		// little-endian
	MOVD	$226, R8			// mprotect
	SVC
	MOVD	R0, ret+24(FP)
	RET

