// +build !windows,!freebsd,!darwin

#include "textflag.h"


TEXT ·sys_mprotect(SB),NOSPLIT,$0
	MOVQ	addr+0(FP), DI
	MOVQ	n+8(FP), SI
	MOVL	prot+16(FP), DX
	MOVL	$10, AX			// mprotect
	SYSCALL
	MOVQ	AX, ret+24(FP)
	RET

TEXT ·sys_sigaction(SB),NOSPLIT,$0
	MOVQ	sig+0(FP), DI
	MOVQ	act+8(FP), SI
	MOVL	oact+16(FP), DX
	MOVL	ss+24(FP), R10
	MOVL	$13, AX
	SYSCALL
	MOVQ	AX, ret+32(FP)
	RET
