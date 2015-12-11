#include "textflag.h"


TEXT ·sys_mprotect(SB),NOSPLIT,$0
	MOVQ	addr+0(FP), DI
	MOVQ	n+8(FP), SI
	MOVQ	prot+16(FP), DX
	MOVQ	$74, AX			// mprotect
	SYSCALL
	MOVQ	AX, ret+24(FP)
	RET
