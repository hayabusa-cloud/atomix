// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64 && lse2

#include "textflag.h"

// ARM64 128-bit operations using LSE2 (ARMv8.4+).
//
// LSE2 guarantees: LDP/STP are single-copy atomic for 16-byte aligned addresses.
// CASP provides atomic 128-bit compare-and-swap.
//
// Build with -tags=lse2 to use this implementation.
// Default (without tag) uses LL/SC (LDXP/STXP) which is faster on some CPUs.
//
// DMB codes:
//   $0x9 = ISHLD (load-load/load-store barrier, acquire)
//   $0xA = ISHST (store-store barrier, release)
//   $0xB = ISH   (full inner-shareable barrier)

// =============================================================================
// 128-bit Load/Store operations
// LSE2: LDP/STP are single-copy atomic for 16-byte aligned addresses
// =============================================================================

// Load relaxed: LDP is atomic with LSE2
TEXT ·LoadUint128Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	LDP	(R8), (R0, R1)
	MOVD	R0, lo+8(FP)
	MOVD	R1, hi+16(FP)
	RET

// Load acquire: LDP + acquire barrier
TEXT ·LoadUint128Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	LDP	(R8), (R0, R1)
	DMB	$0x9
	MOVD	R0, lo+8(FP)
	MOVD	R1, hi+16(FP)
	RET

// Store relaxed: STP is atomic with LSE2
TEXT ·StoreUint128Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	MOVD	lo+8(FP), R0
	MOVD	hi+16(FP), R1
	STP	(R0, R1), (R8)
	RET

// Store release: release barrier + STP
TEXT ·StoreUint128Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	MOVD	lo+8(FP), R0
	MOVD	hi+16(FP), R1
	DMB	$0xA
	STP	(R0, R1), (R8)
	RET

// =============================================================================
// 128-bit Swap operations (CASPD loop)
// =============================================================================

// Swap 128-bit: CAS loop until success, returns old value
// Frame layout: addr+0, newLo+8, newHi+16, oldLo+24, oldHi+32 = 40 bytes
TEXT ·SwapUint128Relaxed(SB), NOSPLIT, $0-40
	MOVD	addr+0(FP), R8
	MOVD	newLo+8(FP), R2
	MOVD	newHi+16(FP), R3
	// Read current value (LDP is atomic with LSE2)
	LDP	(R8), (R0, R1)
swap128_retry:
	MOVD	R0, R4
	MOVD	R1, R5
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	swap128_retry
	MOVD	R4, oldLo+24(FP)
	MOVD	R5, oldHi+32(FP)
	RET

TEXT ·SwapUint128Acquire(SB), NOSPLIT, $0-40
	MOVD	addr+0(FP), R8
	MOVD	newLo+8(FP), R2
	MOVD	newHi+16(FP), R3
	// Read current value (LDP is atomic with LSE2)
	LDP	(R8), (R0, R1)
swap128_acq_retry:
	MOVD	R0, R4
	MOVD	R1, R5
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	swap128_acq_retry
	DMB	$0x9
	MOVD	R4, oldLo+24(FP)
	MOVD	R5, oldHi+32(FP)
	RET

TEXT ·SwapUint128Release(SB), NOSPLIT, $0-40
	MOVD	addr+0(FP), R8
	MOVD	newLo+8(FP), R2
	MOVD	newHi+16(FP), R3
	// Read current value (LDP is atomic with LSE2)
	LDP	(R8), (R0, R1)
	// Release barrier once before store attempts
	DMB	$0xA
swap128_rel_retry:
	MOVD	R0, R4
	MOVD	R1, R5
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	swap128_rel_retry
	MOVD	R4, oldLo+24(FP)
	MOVD	R5, oldHi+32(FP)
	RET

TEXT ·SwapUint128AcqRel(SB), NOSPLIT, $0-40
	MOVD	addr+0(FP), R8
	MOVD	newLo+8(FP), R2
	MOVD	newHi+16(FP), R3
	// Read current value (LDP is atomic with LSE2)
	LDP	(R8), (R0, R1)
	// Release barrier once before store attempts
	DMB	$0xA
swap128_aqrl_retry:
	MOVD	R0, R4
	MOVD	R1, R5
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	swap128_aqrl_retry
	DMB	$0x9
	MOVD	R4, oldLo+24(FP)
	MOVD	R5, oldHi+32(FP)
	RET

// =============================================================================
// 128-bit CAS operations (CASPD)
// =============================================================================

TEXT ·CasUint128Relaxed(SB), NOSPLIT, $0-41
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R0
	MOVD	oldHi+16(FP), R1
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
	MOVD	R0, R4
	MOVD	R1, R5
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	CSET	EQ, R6
	MOVB	R6, ret+40(FP)
	RET

TEXT ·CasUint128Acquire(SB), NOSPLIT, $0-41
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R0
	MOVD	oldHi+16(FP), R1
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
	MOVD	R0, R4
	MOVD	R1, R5
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	cas128_acq_fail
	DMB	$0x9
	MOVD	$1, R3
	MOVB	R3, ret+40(FP)
	RET
cas128_acq_fail:
	MOVB	ZR, ret+40(FP)
	RET

TEXT ·CasUint128Release(SB), NOSPLIT, $0-41
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R0
	MOVD	oldHi+16(FP), R1
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
	MOVD	R0, R4
	MOVD	R1, R5
	DMB	$0xA
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	CSET	EQ, R6
	MOVB	R6, ret+40(FP)
	RET

TEXT ·CasUint128AcqRel(SB), NOSPLIT, $0-41
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R0
	MOVD	oldHi+16(FP), R1
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
	MOVD	R0, R4
	MOVD	R1, R5
	DMB	$0xA
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	cas128_ar_fail
	DMB	$0x9
	MOVD	$1, R3
	MOVB	R3, ret+40(FP)
	RET
cas128_ar_fail:
	MOVB	ZR, ret+40(FP)
	RET

// =============================================================================
// 128-bit CAX (Compare-And-Exchange returning old value) operations
// =============================================================================

TEXT ·CaxUint128Relaxed(SB), NOSPLIT, $0-56
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R0
	MOVD	oldHi+16(FP), R1
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
	CASPD	(R0, R1), (R8), (R2, R3)
	MOVD	R0, lo+40(FP)
	MOVD	R1, hi+48(FP)
	RET

TEXT ·CaxUint128Acquire(SB), NOSPLIT, $0-56
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R0
	MOVD	oldHi+16(FP), R1
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
	MOVD	R0, R4
	MOVD	R1, R5
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	cax128_acq_done
	DMB	$0x9
cax128_acq_done:
	MOVD	R0, lo+40(FP)
	MOVD	R1, hi+48(FP)
	RET

TEXT ·CaxUint128Release(SB), NOSPLIT, $0-56
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R0
	MOVD	oldHi+16(FP), R1
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
	DMB	$0xA
	CASPD	(R0, R1), (R8), (R2, R3)
	MOVD	R0, lo+40(FP)
	MOVD	R1, hi+48(FP)
	RET

TEXT ·CaxUint128AcqRel(SB), NOSPLIT, $0-56
	MOVD	addr+0(FP), R8
	MOVD	oldLo+8(FP), R0
	MOVD	oldHi+16(FP), R1
	MOVD	newLo+24(FP), R2
	MOVD	newHi+32(FP), R3
	MOVD	R0, R4
	MOVD	R1, R5
	DMB	$0xA
	CASPD	(R0, R1), (R8), (R2, R3)
	CMP	R0, R4
	CCMP	EQ, R1, R5, $0
	BNE	cax128_ar_done
	DMB	$0x9
cax128_ar_done:
	MOVD	R0, lo+40(FP)
	MOVD	R1, hi+48(FP)
	RET
