// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64 && !atomix_llsc

#include "textflag.h"

// Go assembly spells the relaxed instruction CASPD. The ordered encodings
// use the same operands: X0:X1 expected/result, X2:X3 desired, [X8] address.
#define CASP_RELAXED CASPD (R0, R1), (R8), (R2, R3)
#define CASP_ACQUIRE WORD $0x48607d02 // CASPA
#define CASP_RELEASE WORD $0x4820fd02 // CASPL
#define CASP_ACQ_REL WORD $0x4860fd02 // CASPAL

// On contention, pause before retrying with the pair returned by CASP.
#define SWAP_PAIR(OP) \
	MOVD addr+0(FP), R8; \
	LDP newLo+8(FP), (R2, R3); \
	LDP (R8), (R0, R1); \
	retry: MOVD R0, R4; \
	MOVD R1, R5; \
	OP; \
	CMP R0, R4; \
	CCMP EQ, R1, R5, $0; \
	BNE contended; \
	STP (R0, R1), oldLo+24(FP); \
	RET; \
	contended: YIELD; \
	B retry

#define LOAD_CAS_ARGS \
	MOVD addr+0(FP), R8; \
	LDP oldLo+8(FP), (R0, R1); \
	LDP newLo+24(FP), (R2, R3)

#define CAS_PAIR(OP) \
	LOAD_CAS_ARGS; \
	MOVD R0, R4; \
	MOVD R1, R5; \
	OP; \
	CMP R0, R4; \
	CCMP EQ, R1, R5, $0; \
	CSET EQ, R6; \
	MOVB R6, ret+40(FP); \
	RET

#define CAX_PAIR(OP) \
	LOAD_CAS_ARGS; \
	OP; \
	STP (R0, R1), lo+40(FP); \
	RET

// FEAT_LSE2 makes aligned 16-byte LDP/STP atomic in Normal Inner and Outer
// Write-Back memory. FEAT_LSE supplies CASP and its ordered variants.
TEXT ·LoadUint128Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	LDP	(R8), (R0, R1)
	STP	(R0, R1), lo+8(FP)
	RET

TEXT ·LoadUint128Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	LDP	(R8), (R0, R1)
	DMB	$0x9 // ISHLD
	STP	(R0, R1), lo+8(FP)
	RET

TEXT ·StoreUint128Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	LDP	lo+8(FP), (R0, R1)
	STP	(R0, R1), (R8)
	RET

TEXT ·StoreUint128Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R8
	LDP	lo+8(FP), (R0, R1)
	DMB	$0xA // ISHST publishes earlier stores.
	STP	(R0, R1), (R8)
	RET

TEXT ·SwapUint128Relaxed(SB), NOSPLIT, $0-40
	SWAP_PAIR(CASP_RELAXED)

TEXT ·SwapUint128Acquire(SB), NOSPLIT, $0-40
	SWAP_PAIR(CASP_ACQUIRE)

TEXT ·SwapUint128Release(SB), NOSPLIT, $0-40
	SWAP_PAIR(CASP_RELEASE)

TEXT ·SwapUint128AcqRel(SB), NOSPLIT, $0-40
	SWAP_PAIR(CASP_ACQ_REL)

TEXT ·CasUint128Relaxed(SB), NOSPLIT, $0-41
	CAS_PAIR(CASP_RELAXED)

TEXT ·CasUint128Acquire(SB), NOSPLIT, $0-41
	CAS_PAIR(CASP_ACQUIRE)

TEXT ·CasUint128Release(SB), NOSPLIT, $0-41
	CAS_PAIR(CASP_RELEASE)

TEXT ·CasUint128AcqRel(SB), NOSPLIT, $0-41
	CAS_PAIR(CASP_ACQ_REL)

TEXT ·CaxUint128Relaxed(SB), NOSPLIT, $0-56
	CAX_PAIR(CASP_RELAXED)

TEXT ·CaxUint128Acquire(SB), NOSPLIT, $0-56
	CAX_PAIR(CASP_ACQUIRE)

TEXT ·CaxUint128Release(SB), NOSPLIT, $0-56
	CAX_PAIR(CASP_RELEASE)

TEXT ·CaxUint128AcqRel(SB), NOSPLIT, $0-56
	CAX_PAIR(CASP_ACQ_REL)
