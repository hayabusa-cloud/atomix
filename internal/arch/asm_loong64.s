// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build loong64

#include "textflag.h"

// LoongArch Atomic Operations:
// Uses LL/SC (Load-Linked/Store-Conditional) for CAS.
// Uses AMSWAPDB/AMADDDB instructions for swap/add.
//
// DBAR hints for memory ordering:
// - $0x14 = LoadAcquire barrier
// - $0x12 = StoreRelease barrier
// - $0    = full barrier (DBAR without hint)

// ============================================================================
// 32-bit Signed Integer Operations
// ============================================================================

// func LoadInt32Relaxed(addr *int32) int32
TEXT ·LoadInt32Relaxed(SB), NOSPLIT, $0-12
	MOVV	addr+0(FP), R4
	MOVW	(R4), R4
	MOVW	R4, ret+8(FP)
	RET

// func LoadInt32Acquire(addr *int32) int32
TEXT ·LoadInt32Acquire(SB), NOSPLIT, $0-12
	MOVV	addr+0(FP), R4
	MOVW	(R4), R4
	DBAR	$0x14
	MOVW	R4, ret+8(FP)
	RET

// func StoreInt32Relaxed(addr *int32, val int32)
TEXT ·StoreInt32Relaxed(SB), NOSPLIT, $0-12
	MOVV	addr+0(FP), R4
	MOVW	val+8(FP), R5
	MOVW	R5, (R4)
	RET

// func StoreInt32Release(addr *int32, val int32)
TEXT ·StoreInt32Release(SB), NOSPLIT, $0-12
	MOVV	addr+0(FP), R4
	MOVW	val+8(FP), R5
	DBAR	$0x12
	MOVW	R5, (R4)
	RET

// func SwapInt32Relaxed(addr *int32, new int32) int32
TEXT ·SwapInt32Relaxed(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	new+8(FP), R5
	AMSWAPDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

// func SwapInt32Acquire(addr *int32, new int32) int32
TEXT ·SwapInt32Acquire(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	new+8(FP), R5
	AMSWAPDBW	R5, (R4), R6
	DBAR	$0x14
	MOVW	R6, ret+16(FP)
	RET

// func SwapInt32Release(addr *int32, new int32) int32
TEXT ·SwapInt32Release(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	new+8(FP), R5
	DBAR	$0x12
	AMSWAPDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

// func SwapInt32AcqRel(addr *int32, new int32) int32
TEXT ·SwapInt32AcqRel(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	new+8(FP), R5
	DBAR	$0x12
	AMSWAPDBW	R5, (R4), R6
	DBAR	$0x14
	MOVW	R6, ret+16(FP)
	RET

// func CasInt32Relaxed(addr *int32, old, new int32) bool
TEXT ·CasInt32Relaxed(SB), NOSPLIT, $0-17
	MOVV	addr+0(FP), R4
	MOVW	old+8(FP), R5
	MOVW	new+12(FP), R6
cas32_relaxed_loop:
	MOVV	R6, R7
	LL	(R4), R8
	BNE	R5, R8, cas32_relaxed_fail
	SC	R7, (R4)
	BEQ	R7, cas32_relaxed_loop
	MOVV	$1, R4
	MOVB	R4, ret+16(FP)
	RET
cas32_relaxed_fail:
	MOVB	R0, ret+16(FP)
	RET

// func CasInt32Acquire(addr *int32, old, new int32) bool
TEXT ·CasInt32Acquire(SB), NOSPLIT, $0-17
	MOVV	addr+0(FP), R4
	MOVW	old+8(FP), R5
	MOVW	new+12(FP), R6
cas32_acq_loop:
	MOVV	R6, R7
	LL	(R4), R8
	BNE	R5, R8, cas32_acq_fail
	SC	R7, (R4)
	BEQ	R7, cas32_acq_loop
	DBAR	$0x14
	MOVV	$1, R4
	MOVB	R4, ret+16(FP)
	RET
cas32_acq_fail:
	MOVB	R0, ret+16(FP)
	RET

// func CasInt32Release(addr *int32, old, new int32) bool
TEXT ·CasInt32Release(SB), NOSPLIT, $0-17
	MOVV	addr+0(FP), R4
	MOVW	old+8(FP), R5
	MOVW	new+12(FP), R6
	DBAR	$0x12
cas32_rel_loop:
	MOVV	R6, R7
	LL	(R4), R8
	BNE	R5, R8, cas32_rel_fail
	SC	R7, (R4)
	BEQ	R7, cas32_rel_loop
	MOVV	$1, R4
	MOVB	R4, ret+16(FP)
	RET
cas32_rel_fail:
	MOVB	R0, ret+16(FP)
	RET

// func CasInt32AcqRel(addr *int32, old, new int32) bool
TEXT ·CasInt32AcqRel(SB), NOSPLIT, $0-17
	MOVV	addr+0(FP), R4
	MOVW	old+8(FP), R5
	MOVW	new+12(FP), R6
	DBAR	$0x12
cas32_aqrl_loop:
	MOVV	R6, R7
	LL	(R4), R8
	BNE	R5, R8, cas32_aqrl_fail
	SC	R7, (R4)
	BEQ	R7, cas32_aqrl_loop
	DBAR	$0x14
	MOVV	$1, R4
	MOVB	R4, ret+16(FP)
	RET
cas32_aqrl_fail:
	MOVB	R0, ret+16(FP)
	RET

// func CaxInt32Relaxed(addr *int32, old, new int32) int32
TEXT ·CaxInt32Relaxed(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	old+8(FP), R5
	MOVW	new+12(FP), R6
cax32_relaxed_loop:
	MOVV	R6, R7
	LL	(R4), R8
	BNE	R5, R8, cax32_relaxed_done
	SC	R7, (R4)
	BEQ	R7, cax32_relaxed_loop
cax32_relaxed_done:
	MOVW	R8, ret+16(FP)
	RET

// func CaxInt32Acquire(addr *int32, old, new int32) int32
TEXT ·CaxInt32Acquire(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	old+8(FP), R5
	MOVW	new+12(FP), R6
cax32_acq_loop:
	MOVV	R6, R7
	LL	(R4), R8
	BNE	R5, R8, cax32_acq_done
	SC	R7, (R4)
	BEQ	R7, cax32_acq_loop
cax32_acq_done:
	DBAR	$0x14
	MOVW	R8, ret+16(FP)
	RET

// func CaxInt32Release(addr *int32, old, new int32) int32
TEXT ·CaxInt32Release(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	old+8(FP), R5
	MOVW	new+12(FP), R6
	DBAR	$0x12
cax32_rel_loop:
	MOVV	R6, R7
	LL	(R4), R8
	BNE	R5, R8, cax32_rel_done
	SC	R7, (R4)
	BEQ	R7, cax32_rel_loop
cax32_rel_done:
	MOVW	R8, ret+16(FP)
	RET

// func CaxInt32AcqRel(addr *int32, old, new int32) int32
TEXT ·CaxInt32AcqRel(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	old+8(FP), R5
	MOVW	new+12(FP), R6
	DBAR	$0x12
cax32_aqrl_loop:
	MOVV	R6, R7
	LL	(R4), R8
	BNE	R5, R8, cax32_aqrl_done
	SC	R7, (R4)
	BEQ	R7, cax32_aqrl_loop
cax32_aqrl_done:
	DBAR	$0x14
	MOVW	R8, ret+16(FP)
	RET

// func AddInt32Relaxed(addr *int32, delta int32) int32
TEXT ·AddInt32Relaxed(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	delta+8(FP), R5
	AMADDDBW	R5, (R4), R6
	ADDV	R5, R6, R4
	MOVW	R4, ret+16(FP)
	RET

// func AddInt32Acquire(addr *int32, delta int32) int32
TEXT ·AddInt32Acquire(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	delta+8(FP), R5
	AMADDDBW	R5, (R4), R6
	DBAR	$0x14
	ADDV	R5, R6, R4
	MOVW	R4, ret+16(FP)
	RET

// func AddInt32Release(addr *int32, delta int32) int32
TEXT ·AddInt32Release(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	delta+8(FP), R5
	DBAR	$0x12
	AMADDDBW	R5, (R4), R6
	ADDV	R5, R6, R4
	MOVW	R4, ret+16(FP)
	RET

// func AddInt32AcqRel(addr *int32, delta int32) int32
TEXT ·AddInt32AcqRel(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	delta+8(FP), R5
	DBAR	$0x12
	AMADDDBW	R5, (R4), R6
	DBAR	$0x14
	ADDV	R5, R6, R4
	MOVW	R4, ret+16(FP)
	RET

// ============================================================================
// 32-bit Unsigned Integer Operations (aliases)
// ============================================================================

TEXT ·LoadUint32Relaxed(SB), NOSPLIT, $0-12
	JMP	·LoadInt32Relaxed(SB)

TEXT ·LoadUint32Acquire(SB), NOSPLIT, $0-12
	JMP	·LoadInt32Acquire(SB)

TEXT ·StoreUint32Relaxed(SB), NOSPLIT, $0-12
	JMP	·StoreInt32Relaxed(SB)

TEXT ·StoreUint32Release(SB), NOSPLIT, $0-12
	JMP	·StoreInt32Release(SB)

TEXT ·SwapUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Relaxed(SB)

TEXT ·SwapUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Acquire(SB)

TEXT ·SwapUint32Release(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Release(SB)

TEXT ·SwapUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

TEXT ·CasUint32Relaxed(SB), NOSPLIT, $0-17
	JMP	·CasInt32Relaxed(SB)

TEXT ·CasUint32Acquire(SB), NOSPLIT, $0-17
	JMP	·CasInt32Acquire(SB)

TEXT ·CasUint32Release(SB), NOSPLIT, $0-17
	JMP	·CasInt32Release(SB)

TEXT ·CasUint32AcqRel(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

TEXT ·CaxUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Relaxed(SB)

TEXT ·CaxUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Acquire(SB)

TEXT ·CaxUint32Release(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Release(SB)

TEXT ·CaxUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

TEXT ·AddUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AddInt32Relaxed(SB)

TEXT ·AddUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·AddInt32Acquire(SB)

TEXT ·AddUint32Release(SB), NOSPLIT, $0-20
	JMP	·AddInt32Release(SB)

TEXT ·AddUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

// ============================================================================
// 64-bit Signed Integer Operations
// ============================================================================

// func LoadInt64Relaxed(addr *int64) int64
TEXT ·LoadInt64Relaxed(SB), NOSPLIT, $0-16
	MOVV	addr+0(FP), R4
	MOVV	(R4), R4
	MOVV	R4, ret+8(FP)
	RET

// func LoadInt64Acquire(addr *int64) int64
TEXT ·LoadInt64Acquire(SB), NOSPLIT, $0-16
	MOVV	addr+0(FP), R4
	MOVV	(R4), R4
	DBAR	$0x14
	MOVV	R4, ret+8(FP)
	RET

// func StoreInt64Relaxed(addr *int64, val int64)
TEXT ·StoreInt64Relaxed(SB), NOSPLIT, $0-16
	MOVV	addr+0(FP), R4
	MOVV	val+8(FP), R5
	MOVV	R5, (R4)
	RET

// func StoreInt64Release(addr *int64, val int64)
TEXT ·StoreInt64Release(SB), NOSPLIT, $0-16
	MOVV	addr+0(FP), R4
	MOVV	val+8(FP), R5
	DBAR	$0x12
	MOVV	R5, (R4)
	RET

// func SwapInt64Relaxed(addr *int64, new int64) int64
TEXT ·SwapInt64Relaxed(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	new+8(FP), R5
	AMSWAPDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

// func SwapInt64Acquire(addr *int64, new int64) int64
TEXT ·SwapInt64Acquire(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	new+8(FP), R5
	AMSWAPDBV	R5, (R4), R6
	DBAR	$0x14
	MOVV	R6, ret+16(FP)
	RET

// func SwapInt64Release(addr *int64, new int64) int64
TEXT ·SwapInt64Release(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	new+8(FP), R5
	DBAR	$0x12
	AMSWAPDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

// func SwapInt64AcqRel(addr *int64, new int64) int64
TEXT ·SwapInt64AcqRel(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	new+8(FP), R5
	DBAR	$0x12
	AMSWAPDBV	R5, (R4), R6
	DBAR	$0x14
	MOVV	R6, ret+16(FP)
	RET

// func CasInt64Relaxed(addr *int64, old, new int64) bool
TEXT ·CasInt64Relaxed(SB), NOSPLIT, $0-25
	MOVV	addr+0(FP), R4
	MOVV	old+8(FP), R5
	MOVV	new+16(FP), R6
cas64_relaxed_loop:
	MOVV	R6, R7
	LLV	(R4), R8
	BNE	R5, R8, cas64_relaxed_fail
	SCV	R7, (R4)
	BEQ	R7, cas64_relaxed_loop
	MOVV	$1, R4
	MOVB	R4, ret+24(FP)
	RET
cas64_relaxed_fail:
	MOVB	R0, ret+24(FP)
	RET

// func CasInt64Acquire(addr *int64, old, new int64) bool
TEXT ·CasInt64Acquire(SB), NOSPLIT, $0-25
	MOVV	addr+0(FP), R4
	MOVV	old+8(FP), R5
	MOVV	new+16(FP), R6
cas64_acq_loop:
	MOVV	R6, R7
	LLV	(R4), R8
	BNE	R5, R8, cas64_acq_fail
	SCV	R7, (R4)
	BEQ	R7, cas64_acq_loop
	DBAR	$0x14
	MOVV	$1, R4
	MOVB	R4, ret+24(FP)
	RET
cas64_acq_fail:
	MOVB	R0, ret+24(FP)
	RET

// func CasInt64Release(addr *int64, old, new int64) bool
TEXT ·CasInt64Release(SB), NOSPLIT, $0-25
	MOVV	addr+0(FP), R4
	MOVV	old+8(FP), R5
	MOVV	new+16(FP), R6
	DBAR	$0x12
cas64_rel_loop:
	MOVV	R6, R7
	LLV	(R4), R8
	BNE	R5, R8, cas64_rel_fail
	SCV	R7, (R4)
	BEQ	R7, cas64_rel_loop
	MOVV	$1, R4
	MOVB	R4, ret+24(FP)
	RET
cas64_rel_fail:
	MOVB	R0, ret+24(FP)
	RET

// func CasInt64AcqRel(addr *int64, old, new int64) bool
TEXT ·CasInt64AcqRel(SB), NOSPLIT, $0-25
	MOVV	addr+0(FP), R4
	MOVV	old+8(FP), R5
	MOVV	new+16(FP), R6
	DBAR	$0x12
cas64_aqrl_loop:
	MOVV	R6, R7
	LLV	(R4), R8
	BNE	R5, R8, cas64_aqrl_fail
	SCV	R7, (R4)
	BEQ	R7, cas64_aqrl_loop
	DBAR	$0x14
	MOVV	$1, R4
	MOVB	R4, ret+24(FP)
	RET
cas64_aqrl_fail:
	MOVB	R0, ret+24(FP)
	RET

// func CaxInt64Relaxed(addr *int64, old, new int64) int64
TEXT ·CaxInt64Relaxed(SB), NOSPLIT, $0-32
	MOVV	addr+0(FP), R4
	MOVV	old+8(FP), R5
	MOVV	new+16(FP), R6
cax64_relaxed_loop:
	MOVV	R6, R7
	LLV	(R4), R8
	BNE	R5, R8, cax64_relaxed_done
	SCV	R7, (R4)
	BEQ	R7, cax64_relaxed_loop
cax64_relaxed_done:
	MOVV	R8, ret+24(FP)
	RET

// func CaxInt64Acquire(addr *int64, old, new int64) int64
TEXT ·CaxInt64Acquire(SB), NOSPLIT, $0-32
	MOVV	addr+0(FP), R4
	MOVV	old+8(FP), R5
	MOVV	new+16(FP), R6
cax64_acq_loop:
	MOVV	R6, R7
	LLV	(R4), R8
	BNE	R5, R8, cax64_acq_done
	SCV	R7, (R4)
	BEQ	R7, cax64_acq_loop
cax64_acq_done:
	DBAR	$0x14
	MOVV	R8, ret+24(FP)
	RET

// func CaxInt64Release(addr *int64, old, new int64) int64
TEXT ·CaxInt64Release(SB), NOSPLIT, $0-32
	MOVV	addr+0(FP), R4
	MOVV	old+8(FP), R5
	MOVV	new+16(FP), R6
	DBAR	$0x12
cax64_rel_loop:
	MOVV	R6, R7
	LLV	(R4), R8
	BNE	R5, R8, cax64_rel_done
	SCV	R7, (R4)
	BEQ	R7, cax64_rel_loop
cax64_rel_done:
	MOVV	R8, ret+24(FP)
	RET

// func CaxInt64AcqRel(addr *int64, old, new int64) int64
TEXT ·CaxInt64AcqRel(SB), NOSPLIT, $0-32
	MOVV	addr+0(FP), R4
	MOVV	old+8(FP), R5
	MOVV	new+16(FP), R6
	DBAR	$0x12
cax64_aqrl_loop:
	MOVV	R6, R7
	LLV	(R4), R8
	BNE	R5, R8, cax64_aqrl_done
	SCV	R7, (R4)
	BEQ	R7, cax64_aqrl_loop
cax64_aqrl_done:
	DBAR	$0x14
	MOVV	R8, ret+24(FP)
	RET

// func AddInt64Relaxed(addr *int64, delta int64) int64
TEXT ·AddInt64Relaxed(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	delta+8(FP), R5
	AMADDDBV	R5, (R4), R6
	ADDV	R5, R6, R4
	MOVV	R4, ret+16(FP)
	RET

// func AddInt64Acquire(addr *int64, delta int64) int64
TEXT ·AddInt64Acquire(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	delta+8(FP), R5
	AMADDDBV	R5, (R4), R6
	DBAR	$0x14
	ADDV	R5, R6, R4
	MOVV	R4, ret+16(FP)
	RET

// func AddInt64Release(addr *int64, delta int64) int64
TEXT ·AddInt64Release(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	delta+8(FP), R5
	DBAR	$0x12
	AMADDDBV	R5, (R4), R6
	ADDV	R5, R6, R4
	MOVV	R4, ret+16(FP)
	RET

// func AddInt64AcqRel(addr *int64, delta int64) int64
TEXT ·AddInt64AcqRel(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	delta+8(FP), R5
	DBAR	$0x12
	AMADDDBV	R5, (R4), R6
	DBAR	$0x14
	ADDV	R5, R6, R4
	MOVV	R4, ret+16(FP)
	RET

// ============================================================================
// 64-bit Unsigned Integer Operations (aliases)
// ============================================================================

TEXT ·LoadUint64Relaxed(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Relaxed(SB)

TEXT ·LoadUint64Acquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

TEXT ·StoreUint64Relaxed(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Relaxed(SB)

TEXT ·StoreUint64Release(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Release(SB)

TEXT ·SwapUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

TEXT ·SwapUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Acquire(SB)

TEXT ·SwapUint64Release(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Release(SB)

TEXT ·SwapUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·CasUint64Relaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

TEXT ·CasUint64Acquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64Acquire(SB)

TEXT ·CasUint64Release(SB), NOSPLIT, $0-25
	JMP	·CasInt64Release(SB)

TEXT ·CasUint64AcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CaxUint64Relaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

TEXT ·CaxUint64Acquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Acquire(SB)

TEXT ·CaxUint64Release(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Release(SB)

TEXT ·CaxUint64AcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·AddUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·AddInt64Relaxed(SB)

TEXT ·AddUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·AddInt64Acquire(SB)

TEXT ·AddUint64Release(SB), NOSPLIT, $0-24
	JMP	·AddInt64Release(SB)

TEXT ·AddUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

// ============================================================================
// Uintptr Operations (aliases to 64-bit)
// ============================================================================

TEXT ·LoadUintptrRelaxed(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Relaxed(SB)

TEXT ·LoadUintptrAcquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

TEXT ·StoreUintptrRelaxed(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Relaxed(SB)

TEXT ·StoreUintptrRelease(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Release(SB)

TEXT ·SwapUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

TEXT ·SwapUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Acquire(SB)

TEXT ·SwapUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Release(SB)

TEXT ·SwapUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·CasUintptrRelaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

TEXT ·CasUintptrAcquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64Acquire(SB)

TEXT ·CasUintptrRelease(SB), NOSPLIT, $0-25
	JMP	·CasInt64Release(SB)

TEXT ·CasUintptrAcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CaxUintptrRelaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

TEXT ·CaxUintptrAcquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Acquire(SB)

TEXT ·CaxUintptrRelease(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Release(SB)

TEXT ·CaxUintptrAcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·AddUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·AddInt64Relaxed(SB)

TEXT ·AddUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·AddInt64Acquire(SB)

TEXT ·AddUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·AddInt64Release(SB)

TEXT ·AddUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

// ============================================================================
// Pointer Operations (aliases to 64-bit)
// ============================================================================

TEXT ·LoadPointerRelaxed(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Relaxed(SB)

TEXT ·LoadPointerAcquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

TEXT ·StorePointerRelaxed(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Relaxed(SB)

TEXT ·StorePointerRelease(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Release(SB)

TEXT ·SwapPointerRelaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

TEXT ·SwapPointerAcquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Acquire(SB)

TEXT ·SwapPointerRelease(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Release(SB)

TEXT ·SwapPointerAcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·CasPointerRelaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

TEXT ·CasPointerAcquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64Acquire(SB)

TEXT ·CasPointerRelease(SB), NOSPLIT, $0-25
	JMP	·CasInt64Release(SB)

TEXT ·CasPointerAcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CaxPointerRelaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

TEXT ·CaxPointerAcquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Acquire(SB)

TEXT ·CaxPointerRelease(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Release(SB)

TEXT ·CaxPointerAcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

// ============================================================================
// 128-bit Operations (emulated via LL/SC on low word)
// ============================================================================

// func LoadUint128Relaxed(addr *[16]byte) (lo, hi uint64)
TEXT ·LoadUint128Relaxed(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
load128_relaxed_loop:
	MOVV	(R4), R5
	MOVV	8(R4), R6
	MOVV	(R4), R7
	BNE	R5, R7, load128_relaxed_loop
	MOVV	R5, lo+8(FP)
	MOVV	R6, hi+16(FP)
	RET

// func LoadUint128Acquire(addr *[16]byte) (lo, hi uint64)
TEXT ·LoadUint128Acquire(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
load128_acq_loop:
	MOVV	(R4), R5
	MOVV	8(R4), R6
	MOVV	(R4), R7
	BNE	R5, R7, load128_acq_loop
	DBAR	$0x14
	MOVV	R5, lo+8(FP)
	MOVV	R6, hi+16(FP)
	RET

// func StoreUint128Relaxed(addr *[16]byte, lo, hi uint64)
TEXT ·StoreUint128Relaxed(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	lo+8(FP), R5
	MOVV	hi+16(FP), R6
store128_relaxed_loop:
	LLV	(R4), R7
	SCV	R5, (R4)
	BEQ	R7, store128_relaxed_loop
	MOVV	R6, 8(R4)
	RET

// func StoreUint128Release(addr *[16]byte, lo, hi uint64)
TEXT ·StoreUint128Release(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	lo+8(FP), R5
	MOVV	hi+16(FP), R6
	DBAR	$0x12
store128_rel_loop:
	LLV	(R4), R7
	SCV	R5, (R4)
	BEQ	R7, store128_rel_loop
	MOVV	R6, 8(R4)
	RET

// func SwapUint128Relaxed(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)
TEXT ·SwapUint128Relaxed(SB), NOSPLIT, $0-40
	MOVV	addr+0(FP), R4
	MOVV	newLo+8(FP), R5
	MOVV	newHi+16(FP), R6
swap128_relaxed_loop:
	LLV	(R4), R7		// load-link old lo
	MOVV	8(R4), R8		// load old hi
	SCV	R5, (R4)		// store-cond new lo
	BEQ	R7, swap128_relaxed_loop
	MOVV	R6, 8(R4)		// store new hi
	MOVV	R7, oldLo+24(FP)
	MOVV	R8, oldHi+32(FP)
	RET

// func SwapUint128Acquire(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)
TEXT ·SwapUint128Acquire(SB), NOSPLIT, $0-40
	JMP	·SwapUint128Relaxed(SB)

// func SwapUint128Release(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)
TEXT ·SwapUint128Release(SB), NOSPLIT, $0-40
	JMP	·SwapUint128Relaxed(SB)

// func SwapUint128AcqRel(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)
TEXT ·SwapUint128AcqRel(SB), NOSPLIT, $0-40
	JMP	·SwapUint128Relaxed(SB)

// func CasUint128Relaxed(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool
TEXT ·CasUint128Relaxed(SB), NOSPLIT, $0-41
	MOVV	addr+0(FP), R4
	MOVV	oldLo+8(FP), R5
	MOVV	oldHi+16(FP), R6
	MOVV	newLo+24(FP), R7
	MOVV	newHi+32(FP), R8
cas128_relaxed_loop:
	LLV	(R4), R9
	BNE	R9, R5, cas128_relaxed_fail
	MOVV	8(R4), R10
	BNE	R10, R6, cas128_relaxed_fail
	SCV	R7, (R4)
	BEQ	R7, cas128_relaxed_loop
	MOVV	R8, 8(R4)
	MOVV	$1, R4
	MOVB	R4, ret+40(FP)
	RET
cas128_relaxed_fail:
	MOVB	R0, ret+40(FP)
	RET

// func CasUint128Acquire(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool
TEXT ·CasUint128Acquire(SB), NOSPLIT, $0-41
	MOVV	addr+0(FP), R4
	MOVV	oldLo+8(FP), R5
	MOVV	oldHi+16(FP), R6
	MOVV	newLo+24(FP), R7
	MOVV	newHi+32(FP), R8
cas128_acq_loop:
	LLV	(R4), R9
	BNE	R9, R5, cas128_acq_fail
	MOVV	8(R4), R10
	BNE	R10, R6, cas128_acq_fail
	SCV	R7, (R4)
	BEQ	R7, cas128_acq_loop
	MOVV	R8, 8(R4)
	DBAR	$0x14
	MOVV	$1, R4
	MOVB	R4, ret+40(FP)
	RET
cas128_acq_fail:
	MOVB	R0, ret+40(FP)
	RET

// func CasUint128Release(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool
TEXT ·CasUint128Release(SB), NOSPLIT, $0-41
	MOVV	addr+0(FP), R4
	MOVV	oldLo+8(FP), R5
	MOVV	oldHi+16(FP), R6
	MOVV	newLo+24(FP), R7
	MOVV	newHi+32(FP), R8
	DBAR	$0x12
cas128_rel_loop:
	LLV	(R4), R9
	BNE	R9, R5, cas128_rel_fail
	MOVV	8(R4), R10
	BNE	R10, R6, cas128_rel_fail
	SCV	R7, (R4)
	BEQ	R7, cas128_rel_loop
	MOVV	R8, 8(R4)
	MOVV	$1, R4
	MOVB	R4, ret+40(FP)
	RET
cas128_rel_fail:
	MOVB	R0, ret+40(FP)
	RET

// func CasUint128AcqRel(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool
TEXT ·CasUint128AcqRel(SB), NOSPLIT, $0-41
	MOVV	addr+0(FP), R4
	MOVV	oldLo+8(FP), R5
	MOVV	oldHi+16(FP), R6
	MOVV	newLo+24(FP), R7
	MOVV	newHi+32(FP), R8
	DBAR	$0x12
cas128_aqrl_loop:
	LLV	(R4), R9
	BNE	R9, R5, cas128_aqrl_fail
	MOVV	8(R4), R10
	BNE	R10, R6, cas128_aqrl_fail
	SCV	R7, (R4)
	BEQ	R7, cas128_aqrl_loop
	MOVV	R8, 8(R4)
	DBAR	$0x14
	MOVV	$1, R4
	MOVB	R4, ret+40(FP)
	RET
cas128_aqrl_fail:
	MOVB	R0, ret+40(FP)
	RET

// func CaxUint128Relaxed(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)
TEXT ·CaxUint128Relaxed(SB), NOSPLIT, $0-56
	MOVV	addr+0(FP), R4
	MOVV	oldLo+8(FP), R5
	MOVV	oldHi+16(FP), R6
	MOVV	newLo+24(FP), R7
	MOVV	newHi+32(FP), R8
cax128_relaxed_loop:
	LLV	(R4), R9
	MOVV	8(R4), R10
	BNE	R9, R5, cax128_relaxed_done
	BNE	R10, R6, cax128_relaxed_done
	SCV	R7, (R4)
	BEQ	R7, cax128_relaxed_loop
	MOVV	R8, 8(R4)
cax128_relaxed_done:
	MOVV	R9, lo+40(FP)
	MOVV	R10, hi+48(FP)
	RET

// func CaxUint128Acquire(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)
TEXT ·CaxUint128Acquire(SB), NOSPLIT, $0-56
	MOVV	addr+0(FP), R4
	MOVV	oldLo+8(FP), R5
	MOVV	oldHi+16(FP), R6
	MOVV	newLo+24(FP), R7
	MOVV	newHi+32(FP), R8
cax128_acq_loop:
	LLV	(R4), R9
	MOVV	8(R4), R10
	BNE	R9, R5, cax128_acq_done
	BNE	R10, R6, cax128_acq_done
	SCV	R7, (R4)
	BEQ	R7, cax128_acq_loop
	MOVV	R8, 8(R4)
cax128_acq_done:
	DBAR	$0x14
	MOVV	R9, lo+40(FP)
	MOVV	R10, hi+48(FP)
	RET

// func CaxUint128Release(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)
TEXT ·CaxUint128Release(SB), NOSPLIT, $0-56
	MOVV	addr+0(FP), R4
	MOVV	oldLo+8(FP), R5
	MOVV	oldHi+16(FP), R6
	MOVV	newLo+24(FP), R7
	MOVV	newHi+32(FP), R8
	DBAR	$0x12
cax128_rel_loop:
	LLV	(R4), R9
	MOVV	8(R4), R10
	BNE	R9, R5, cax128_rel_done
	BNE	R10, R6, cax128_rel_done
	SCV	R7, (R4)
	BEQ	R7, cax128_rel_loop
	MOVV	R8, 8(R4)
cax128_rel_done:
	MOVV	R9, lo+40(FP)
	MOVV	R10, hi+48(FP)
	RET

// func CaxUint128AcqRel(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)
TEXT ·CaxUint128AcqRel(SB), NOSPLIT, $0-56
	MOVV	addr+0(FP), R4
	MOVV	oldLo+8(FP), R5
	MOVV	oldHi+16(FP), R6
	MOVV	newLo+24(FP), R7
	MOVV	newHi+32(FP), R8
	DBAR	$0x12
cax128_aqrl_loop:
	LLV	(R4), R9
	MOVV	8(R4), R10
	BNE	R9, R5, cax128_aqrl_done
	BNE	R10, R6, cax128_aqrl_done
	SCV	R7, (R4)
	BEQ	R7, cax128_aqrl_loop
	MOVV	R8, 8(R4)
cax128_aqrl_done:
	DBAR	$0x14
	MOVV	R9, lo+40(FP)
	MOVV	R10, hi+48(FP)
	RET

// ============================================================================
// Bitwise Operations (And, Or, Xor)
// ============================================================================
//
// LoongArch has AMAND, AMOR, AMXOR instructions that atomically perform
// bitwise operations and return the old value.

// func AndInt32Relaxed(addr *int32, mask int32) int32
TEXT ·AndInt32Relaxed(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMANDW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·AndInt32Acquire(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMANDDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·AndInt32Release(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMANDDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·AndInt32AcqRel(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMANDDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·AndUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

TEXT ·AndUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·AndInt32Acquire(SB)

TEXT ·AndUint32Release(SB), NOSPLIT, $0-20
	JMP	·AndInt32Release(SB)

TEXT ·AndUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

// func AndInt64Relaxed(addr *int64, mask int64) int64
TEXT ·AndInt64Relaxed(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMANDV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·AndInt64Acquire(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMANDDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·AndInt64Release(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMANDDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·AndInt64AcqRel(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMANDDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·AndUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·AndInt64Acquire(SB)

TEXT ·AndUint64Release(SB), NOSPLIT, $0-24
	JMP	·AndInt64Release(SB)

TEXT ·AndUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·AndInt64Acquire(SB)

TEXT ·AndUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·AndInt64Release(SB)

TEXT ·AndUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

// func OrInt32Relaxed(addr *int32, mask int32) int32
TEXT ·OrInt32Relaxed(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMORW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·OrInt32Acquire(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMORDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·OrInt32Release(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMORDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·OrInt32AcqRel(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMORDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·OrUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

TEXT ·OrUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·OrInt32Acquire(SB)

TEXT ·OrUint32Release(SB), NOSPLIT, $0-20
	JMP	·OrInt32Release(SB)

TEXT ·OrUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

// func OrInt64Relaxed(addr *int64, mask int64) int64
TEXT ·OrInt64Relaxed(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMORV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·OrInt64Acquire(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMORDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·OrInt64Release(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMORDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·OrInt64AcqRel(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMORDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·OrUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·OrInt64Acquire(SB)

TEXT ·OrUint64Release(SB), NOSPLIT, $0-24
	JMP	·OrInt64Release(SB)

TEXT ·OrUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·OrInt64Acquire(SB)

TEXT ·OrUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·OrInt64Release(SB)

TEXT ·OrUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

// func XorInt32Relaxed(addr *int32, mask int32) int32
TEXT ·XorInt32Relaxed(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMXORW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·XorInt32Acquire(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMXORDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·XorInt32Release(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMXORDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·XorInt32AcqRel(SB), NOSPLIT, $0-20
	MOVV	addr+0(FP), R4
	MOVW	mask+8(FP), R5
	AMXORDBW	R5, (R4), R6
	MOVW	R6, ret+16(FP)
	RET

TEXT ·XorUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

TEXT ·XorUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·XorInt32Acquire(SB)

TEXT ·XorUint32Release(SB), NOSPLIT, $0-20
	JMP	·XorInt32Release(SB)

TEXT ·XorUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

// func XorInt64Relaxed(addr *int64, mask int64) int64
TEXT ·XorInt64Relaxed(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMXORV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·XorInt64Acquire(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMXORDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·XorInt64Release(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMXORDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·XorInt64AcqRel(SB), NOSPLIT, $0-24
	MOVV	addr+0(FP), R4
	MOVV	mask+8(FP), R5
	AMXORDBV	R5, (R4), R6
	MOVV	R6, ret+16(FP)
	RET

TEXT ·XorUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·XorInt64Acquire(SB)

TEXT ·XorUint64Release(SB), NOSPLIT, $0-24
	JMP	·XorInt64Release(SB)

TEXT ·XorUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·XorInt64Acquire(SB)

TEXT ·XorUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·XorInt64Release(SB)

TEXT ·XorUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

// ============================================================================
// Barrier Operations
// ============================================================================

// func BarrierAcquire()
TEXT ·BarrierAcquire(SB), NOSPLIT, $0-0
	DBAR	$0x14
	RET

// func BarrierRelease()
TEXT ·BarrierRelease(SB), NOSPLIT, $0-0
	DBAR	$0x12
	RET

// func BarrierAcqRel()
TEXT ·BarrierAcqRel(SB), NOSPLIT, $0-0
	DBAR
	RET
