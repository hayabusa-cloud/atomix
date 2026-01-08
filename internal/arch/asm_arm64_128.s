// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64 && !lse2

#include "textflag.h"

// ARM64 128-bit operations using LL/SC (Load-Link/Store-Conditional).
//
// This is the default implementation using LDXP/STXP instructions.
// LL/SC is 15-18% faster than CASP on Graviton4 due to lower
// microarchitectural overhead in uncontended cases.
//
// For LSE2/CASP implementation, build with -tags=lse2.
//
// Instructions used:
//   LDXP  - Load exclusive pair (relaxed)
//   LDAXP - Load acquire exclusive pair
//   STXP  - Store exclusive pair (relaxed)
//   STLXP - Store release exclusive pair
//   CLREX - Clear exclusive monitor
//
// DMB codes:
//   $0x9 = ISHLD (acquire barrier)
//   $0xA = ISHST (release barrier)
//   $0xB = ISH   (full barrier)

// =============================================================================
// 128-bit Load operations
// Use LDXP + CLREX for atomic snapshot without LSE2
// =============================================================================

// Load relaxed: LDXP + CLREX
TEXT ·LoadUint128Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	LDXP	(R8), (R0, R1)
	CLREX
	MOVD	R0, lo+8(FP)
	MOVD	R1, hi+16(FP)
	RET

// Load acquire: LDAXP + CLREX
TEXT ·LoadUint128Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	LDAXP	(R8), (R0, R1)
	CLREX
	MOVD	R0, lo+8(FP)
	MOVD	R1, hi+16(FP)
	RET

// =============================================================================
// 128-bit Store operations
// Use LDXP/STXP loop to atomically replace value
// =============================================================================

// Store relaxed: LDXP + STXP loop
TEXT ·StoreUint128Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	MOVD	lo+8(FP), R2
	MOVD	hi+16(FP), R3
store128_rel_retry:
	LDXP	(R8), (R0, R1)
	STXP	(R2, R3), (R8), R4
	CBNZ	R4, store128_rel_retry
	RET

// Store release: LDXP + STLXP loop
TEXT ·StoreUint128Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	MOVD	lo+8(FP), R2
	MOVD	hi+16(FP), R3
store128_release_retry:
	LDXP	(R8), (R0, R1)
	STLXP	(R2, R3), (R8), R4
	CBNZ	R4, store128_release_retry
	RET

// =============================================================================
// 128-bit Swap operations (LL/SC loop)
// =============================================================================

// Swap 128-bit relaxed: LDXP + STXP loop
TEXT ·SwapUint128Relaxed(SB), NOSPLIT, $0-40
	MOVD	addr+0(FP), R8
	MOVD	newLo+8(FP), R2
	MOVD	newHi+16(FP), R3
swap128_llsc_retry:
	LDXP	(R8), (R0, R1)
	STXP	(R2, R3), (R8), R4
	CBNZ	R4, swap128_llsc_retry
	MOVD	R0, oldLo+24(FP)
	MOVD	R1, oldHi+32(FP)
	RET

// Swap 128-bit acquire: LDAXP + STXP loop
TEXT ·SwapUint128Acquire(SB), NOSPLIT, $0-40
	MOVD	addr+0(FP), R8
	MOVD	newLo+8(FP), R2
	MOVD	newHi+16(FP), R3
swap128_acq_llsc_retry:
	LDAXP	(R8), (R0, R1)
	STXP	(R2, R3), (R8), R4
	CBNZ	R4, swap128_acq_llsc_retry
	MOVD	R0, oldLo+24(FP)
	MOVD	R1, oldHi+32(FP)
	RET

// Swap 128-bit release: LDXP + STLXP loop
TEXT ·SwapUint128Release(SB), NOSPLIT, $0-40
	MOVD	addr+0(FP), R8
	MOVD	newLo+8(FP), R2
	MOVD	newHi+16(FP), R3
swap128_rel_llsc_retry:
	LDXP	(R8), (R0, R1)
	STLXP	(R2, R3), (R8), R4
	CBNZ	R4, swap128_rel_llsc_retry
	MOVD	R0, oldLo+24(FP)
	MOVD	R1, oldHi+32(FP)
	RET

// Swap 128-bit acqrel: LDAXP + STLXP loop
TEXT ·SwapUint128AcqRel(SB), NOSPLIT, $0-40
	MOVD	addr+0(FP), R8
	MOVD	newLo+8(FP), R2
	MOVD	newHi+16(FP), R3
swap128_ar_llsc_retry:
	LDAXP	(R8), (R0, R1)
	STLXP	(R2, R3), (R8), R4
	CBNZ	R4, swap128_ar_llsc_retry
	MOVD	R0, oldLo+24(FP)
	MOVD	R1, oldHi+32(FP)
	RET

// =============================================================================
// 128-bit CAS operations (LL/SC)
// =============================================================================

// CAS 128-bit relaxed: LDXP + compare + STXP
TEXT ·CasUint128Relaxed(SB), NOSPLIT, $0-41
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R4
	MOVD	oldHi+16(FP), R5
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
cas128_llsc_retry:
	LDXP	(R8), (R0, R1)
	CMP	R0, R4
	BNE	cas128_llsc_fail
	CMP	R1, R5
	BNE	cas128_llsc_fail
	STXP	(R2, R3), (R8), R6
	CBNZ	R6, cas128_llsc_retry
	MOVD	$1, R0
	MOVB	R0, ret+40(FP)
	RET
cas128_llsc_fail:
	CLREX
	MOVB	ZR, ret+40(FP)
	RET

// CAS 128-bit acquire: LDAXP + compare + STXP
TEXT ·CasUint128Acquire(SB), NOSPLIT, $0-41
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R4
	MOVD	oldHi+16(FP), R5
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
cas128_acq_llsc_retry:
	LDAXP	(R8), (R0, R1)
	CMP	R0, R4
	BNE	cas128_acq_llsc_fail
	CMP	R1, R5
	BNE	cas128_acq_llsc_fail
	STXP	(R2, R3), (R8), R6
	CBNZ	R6, cas128_acq_llsc_retry
	MOVD	$1, R0
	MOVB	R0, ret+40(FP)
	RET
cas128_acq_llsc_fail:
	CLREX
	MOVB	ZR, ret+40(FP)
	RET

// CAS 128-bit release: LDXP + compare + STLXP
TEXT ·CasUint128Release(SB), NOSPLIT, $0-41
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R4
	MOVD	oldHi+16(FP), R5
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
cas128_rel_llsc_retry:
	LDXP	(R8), (R0, R1)
	CMP	R0, R4
	BNE	cas128_rel_llsc_fail
	CMP	R1, R5
	BNE	cas128_rel_llsc_fail
	STLXP	(R2, R3), (R8), R6
	CBNZ	R6, cas128_rel_llsc_retry
	MOVD	$1, R0
	MOVB	R0, ret+40(FP)
	RET
cas128_rel_llsc_fail:
	CLREX
	MOVB	ZR, ret+40(FP)
	RET

// CAS 128-bit acqrel: LDAXP + compare + STLXP
TEXT ·CasUint128AcqRel(SB), NOSPLIT, $0-41
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R4
	MOVD	oldHi+16(FP), R5
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
cas128_ar_llsc_retry:
	LDAXP	(R8), (R0, R1)
	CMP	R0, R4
	BNE	cas128_ar_llsc_fail
	CMP	R1, R5
	BNE	cas128_ar_llsc_fail
	STLXP	(R2, R3), (R8), R6
	CBNZ	R6, cas128_ar_llsc_retry
	MOVD	$1, R0
	MOVB	R0, ret+40(FP)
	RET
cas128_ar_llsc_fail:
	CLREX
	MOVB	ZR, ret+40(FP)
	RET

// =============================================================================
// 128-bit CAX (Compare-And-Exchange returning old value) operations
// =============================================================================

// CAX 128-bit relaxed: LDXP + compare + STXP
TEXT ·CaxUint128Relaxed(SB), NOSPLIT, $0-56
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R4
	MOVD	oldHi+16(FP), R5
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
cax128_llsc_retry:
	LDXP	(R8), (R0, R1)
	CMP	R0, R4
	BNE	cax128_llsc_done
	CMP	R1, R5
	BNE	cax128_llsc_done
	STXP	(R2, R3), (R8), R6
	CBNZ	R6, cax128_llsc_retry
cax128_llsc_done:
	CLREX
	MOVD	R0, lo+40(FP)
	MOVD	R1, hi+48(FP)
	RET

// CAX 128-bit acquire: LDAXP + compare + STXP
TEXT ·CaxUint128Acquire(SB), NOSPLIT, $0-56
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R4
	MOVD	oldHi+16(FP), R5
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
cax128_acq_llsc_retry:
	LDAXP	(R8), (R0, R1)
	CMP	R0, R4
	BNE	cax128_acq_llsc_done
	CMP	R1, R5
	BNE	cax128_acq_llsc_done
	STXP	(R2, R3), (R8), R6
	CBNZ	R6, cax128_acq_llsc_retry
cax128_acq_llsc_done:
	CLREX
	MOVD	R0, lo+40(FP)
	MOVD	R1, hi+48(FP)
	RET

// CAX 128-bit release: LDXP + compare + STLXP
TEXT ·CaxUint128Release(SB), NOSPLIT, $0-56
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R4
	MOVD	oldHi+16(FP), R5
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
cax128_rel_llsc_retry:
	LDXP	(R8), (R0, R1)
	CMP	R0, R4
	BNE	cax128_rel_llsc_done
	CMP	R1, R5
	BNE	cax128_rel_llsc_done
	STLXP	(R2, R3), (R8), R6
	CBNZ	R6, cax128_rel_llsc_retry
cax128_rel_llsc_done:
	CLREX
	MOVD	R0, lo+40(FP)
	MOVD	R1, hi+48(FP)
	RET

// CAX 128-bit acqrel: LDAXP + compare + STLXP
TEXT ·CaxUint128AcqRel(SB), NOSPLIT, $0-56
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R4
	MOVD	oldHi+16(FP), R5
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
cax128_ar_llsc_retry:
	LDAXP	(R8), (R0, R1)
	CMP	R0, R4
	BNE	cax128_ar_llsc_done
	CMP	R1, R5
	BNE	cax128_ar_llsc_done
	STLXP	(R2, R3), (R8), R6
	CBNZ	R6, cax128_ar_llsc_retry
cax128_ar_llsc_done:
	CLREX
	MOVD	R0, lo+40(FP)
	MOVD	R1, hi+48(FP)
	RET
