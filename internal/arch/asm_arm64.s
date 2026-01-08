// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64

#include "textflag.h"

// ARM64 32/64-bit atomic operations using LSE (ARMv8.1+).
//
// Go 1.25 supports ordered LSE variants for 32/64-bit operations:
//   SWPAD/SWPLD/SWPALD, LDADDAD/LDADDLD/LDADDALD, CASAD/CASLD/CASALD
//
// 128-bit operations are in separate files with build tags:
//   asm_arm64_128.s      - LL/SC (LDXP/STXP), default, faster on Graviton4
//   asm_arm64_128_lse2.s - CASP (LSE2), use -tags=lse2
//
// DMB barrier codes:
//   $0x9 = ISHLD (acquire barrier)
//   $0xA = ISHST (release barrier)
//   $0xB = ISH   (full barrier)

// =============================================================================
// 32-bit signed integer operations
// =============================================================================
//
// Relaxed Load/Store are implemented as pure Go in loadstore_arm64.go
// for inlining optimization. Only acquire/release versions need assembly.

// Load acquire (requires LDAR instruction)
TEXT ·LoadInt32Acquire(SB), NOSPLIT, $0-12
	MOVD	addr+0(FP), R0
	LDARW	(R0), R0
	MOVW	R0, ret+8(FP)
	RET

// Store release (requires STLR instruction)
TEXT ·StoreInt32Release(SB), NOSPLIT, $0-12
	MOVD	addr+0(FP), R0
	MOVW	val+8(FP), R1
	STLRW	R1, (R0)
	RET

// Swap relaxed: SWP
TEXT ·SwapInt32Relaxed(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	new+8(FP), R1
	SWPW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Swap acquire: SWPAW
TEXT ·SwapInt32Acquire(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	new+8(FP), R1
	SWPAW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Swap release: SWPLW
TEXT ·SwapInt32Release(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	new+8(FP), R1
	SWPLW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Swap acqrel: SWPALW
TEXT ·SwapInt32AcqRel(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	new+8(FP), R1
	SWPALW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// CAS relaxed
// ARM64 CASW: R1 (expected) is overwritten with loaded value
// Must save expected before CAS, then compare loaded vs saved expected
TEXT ·CasInt32Relaxed(SB), NOSPLIT, $0-21
	MOVD	addr+0(FP), R0
	MOVW	old+8(FP), R1
	MOVW	new+12(FP), R2
	MOVW	R1, R3		// Save expected value
	CASW	R1, (R0), R2	// R1 = loaded value from memory
	CMPW	R1, R3		// Compare loaded with expected (32-bit)
	CSET	EQ, R4
	MOVB	R4, ret+16(FP)
	RET

// CAS acquire: CASAW
TEXT ·CasInt32Acquire(SB), NOSPLIT, $0-21
	MOVD	addr+0(FP), R0
	MOVW	old+8(FP), R1
	MOVW	new+12(FP), R2
	MOVW	R1, R3		// Save expected value
	CASAW	R1, (R0), R2	// R1 = loaded value
	CMPW	R1, R3		// Compare loaded with expected (32-bit)
	CSET	EQ, R4
	MOVB	R4, ret+16(FP)
	RET

// CAS release: CASLW
TEXT ·CasInt32Release(SB), NOSPLIT, $0-21
	MOVD	addr+0(FP), R0
	MOVW	old+8(FP), R1
	MOVW	new+12(FP), R2
	MOVW	R1, R3		// Save expected value
	CASLW	R1, (R0), R2	// R1 = loaded value
	CMPW	R1, R3		// Compare loaded with expected (32-bit)
	CSET	EQ, R4
	MOVB	R4, ret+16(FP)
	RET

// CAS acqrel: CASALW
TEXT ·CasInt32AcqRel(SB), NOSPLIT, $0-21
	MOVD	addr+0(FP), R0
	MOVW	old+8(FP), R1
	MOVW	new+12(FP), R2
	MOVW	R1, R3		// Save expected value
	CASALW	R1, (R0), R2	// R1 = loaded value
	CMPW	R1, R3		// Compare loaded with expected (32-bit)
	CSET	EQ, R4
	MOVB	R4, ret+16(FP)
	RET

// CompareExchange relaxed (returns old value)
// ARM64 CAS: R1 is overwritten with loaded value from memory
TEXT ·CaxInt32Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVW	old+8(FP), R1
	MOVW	new+12(FP), R2
	CASW	R1, (R0), R2	// R1 = loaded value (old)
	MOVW	R1, ret+16(FP)	// Return loaded value
	RET

// CompareExchange acquire: CASAW
TEXT ·CaxInt32Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVW	old+8(FP), R1
	MOVW	new+12(FP), R2
	CASAW	R1, (R0), R2	// R1 = loaded value
	MOVW	R1, ret+16(FP)	// Return loaded value
	RET

// CompareExchange release: CASLW
TEXT ·CaxInt32Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVW	old+8(FP), R1
	MOVW	new+12(FP), R2
	CASLW	R1, (R0), R2	// R1 = loaded value
	MOVW	R1, ret+16(FP)	// Return loaded value
	RET

// CompareExchange acqrel: CASALW
TEXT ·CaxInt32AcqRel(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVW	old+8(FP), R1
	MOVW	new+12(FP), R2
	CASALW	R1, (R0), R2	// R1 = loaded value
	MOVW	R1, ret+16(FP)	// Return loaded value
	RET

// Add relaxed: LDADD + ADD (returns new value)
TEXT ·AddInt32Relaxed(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	delta+8(FP), R1
	LDADDW	R1, (R0), R2
	ADDW	R1, R2, R2
	MOVW	R2, ret+16(FP)
	RET

// Add acquire: LDADDAW + ADD (returns new value)
TEXT ·AddInt32Acquire(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	delta+8(FP), R1
	LDADDAW	R1, (R0), R2
	ADDW	R1, R2, R2
	MOVW	R2, ret+16(FP)
	RET

// Add release: LDADDLW + ADD (returns new value)
TEXT ·AddInt32Release(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	delta+8(FP), R1
	LDADDLW	R1, (R0), R2
	ADDW	R1, R2, R2
	MOVW	R2, ret+16(FP)
	RET

// Add acqrel: LDADDALW + ADD (returns new value)
TEXT ·AddInt32AcqRel(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	delta+8(FP), R1
	LDADDALW	R1, (R0), R2
	ADDW	R1, R2, R2
	MOVW	R2, ret+16(FP)
	RET

// =============================================================================
// 32-bit unsigned integer operations (aliases)
// =============================================================================
//
// Relaxed Load/Store are in loadstore_arm64.go

TEXT ·LoadUint32Acquire(SB), NOSPLIT, $0-12
	JMP	·LoadInt32Acquire(SB)

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

TEXT ·CasUint32Relaxed(SB), NOSPLIT, $0-21
	JMP	·CasInt32Relaxed(SB)

TEXT ·CasUint32Acquire(SB), NOSPLIT, $0-21
	JMP	·CasInt32Acquire(SB)

TEXT ·CasUint32Release(SB), NOSPLIT, $0-21
	JMP	·CasInt32Release(SB)

TEXT ·CasUint32AcqRel(SB), NOSPLIT, $0-21
	JMP	·CasInt32AcqRel(SB)

TEXT ·CaxUint32Relaxed(SB), NOSPLIT, $0-24
	JMP	·CaxInt32Relaxed(SB)

TEXT ·CaxUint32Acquire(SB), NOSPLIT, $0-24
	JMP	·CaxInt32Acquire(SB)

TEXT ·CaxUint32Release(SB), NOSPLIT, $0-24
	JMP	·CaxInt32Release(SB)

TEXT ·CaxUint32AcqRel(SB), NOSPLIT, $0-24
	JMP	·CaxInt32AcqRel(SB)

TEXT ·AddUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AddInt32Relaxed(SB)

TEXT ·AddUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·AddInt32Acquire(SB)

TEXT ·AddUint32Release(SB), NOSPLIT, $0-20
	JMP	·AddInt32Release(SB)

TEXT ·AddUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

// =============================================================================
// 64-bit signed integer operations
// =============================================================================
//
// Relaxed Load/Store are in loadstore_arm64.go

TEXT ·LoadInt64Acquire(SB), NOSPLIT, $0-16
	MOVD	addr+0(FP), R0
	LDAR	(R0), R0
	MOVD	R0, ret+8(FP)
	RET

TEXT ·StoreInt64Release(SB), NOSPLIT, $0-16
	MOVD	addr+0(FP), R0
	MOVD	val+8(FP), R1
	STLR	R1, (R0)
	RET

TEXT ·SwapInt64Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	new+8(FP), R1
	SWPD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

TEXT ·SwapInt64Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	new+8(FP), R1
	SWPAD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

TEXT ·SwapInt64Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	new+8(FP), R1
	SWPLD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

TEXT ·SwapInt64AcqRel(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	new+8(FP), R1
	SWPALD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

TEXT ·CasInt64Relaxed(SB), NOSPLIT, $0-25
	MOVD	addr+0(FP), R0
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	MOVD	R1, R3		// Save expected value
	CASD	R1, (R0), R2	// R1 = loaded value
	CMP	R1, R3		// Compare loaded with expected
	CSET	EQ, R4
	MOVB	R4, ret+24(FP)
	RET

TEXT ·CasInt64Acquire(SB), NOSPLIT, $0-25
	MOVD	addr+0(FP), R0
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	MOVD	R1, R3		// Save expected value
	CASAD	R1, (R0), R2	// R1 = loaded value
	CMP	R1, R3		// Compare loaded with expected
	CSET	EQ, R4
	MOVB	R4, ret+24(FP)
	RET

TEXT ·CasInt64Release(SB), NOSPLIT, $0-25
	MOVD	addr+0(FP), R0
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	MOVD	R1, R3		// Save expected value
	CASLD	R1, (R0), R2	// R1 = loaded value
	CMP	R1, R3		// Compare loaded with expected
	CSET	EQ, R4
	MOVB	R4, ret+24(FP)
	RET

TEXT ·CasInt64AcqRel(SB), NOSPLIT, $0-25
	MOVD	addr+0(FP), R0
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	MOVD	R1, R3		// Save expected value
	CASALD	R1, (R0), R2	// R1 = loaded value
	CMP	R1, R3		// Compare loaded with expected
	CSET	EQ, R4
	MOVB	R4, ret+24(FP)
	RET

TEXT ·CaxInt64Relaxed(SB), NOSPLIT, $0-32
	MOVD	addr+0(FP), R0
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	CASD	R1, (R0), R2	// R1 = loaded value (old)
	MOVD	R1, ret+24(FP)	// Return loaded value
	RET

TEXT ·CaxInt64Acquire(SB), NOSPLIT, $0-32
	MOVD	addr+0(FP), R0
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	CASAD	R1, (R0), R2	// R1 = loaded value
	MOVD	R1, ret+24(FP)
	RET

TEXT ·CaxInt64Release(SB), NOSPLIT, $0-32
	MOVD	addr+0(FP), R0
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	CASLD	R1, (R0), R2	// R1 = loaded value
	MOVD	R1, ret+24(FP)	// Return loaded value
	RET

TEXT ·CaxInt64AcqRel(SB), NOSPLIT, $0-32
	MOVD	addr+0(FP), R0
	MOVD	old+8(FP), R1
	MOVD	new+16(FP), R2
	CASALD	R1, (R0), R2	// R1 = loaded value
	MOVD	R1, ret+24(FP)	// Return loaded value
	RET

// Add relaxed: LDADDD + ADD (returns new value)
TEXT ·AddInt64Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	delta+8(FP), R1
	LDADDD	R1, (R0), R2
	ADD	R1, R2, R2
	MOVD	R2, ret+16(FP)
	RET

// Add acquire: LDADDAD + ADD (returns new value)
TEXT ·AddInt64Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	delta+8(FP), R1
	LDADDAD	R1, (R0), R2
	ADD	R1, R2, R2
	MOVD	R2, ret+16(FP)
	RET

// Add release: LDADDLD + ADD (returns new value)
TEXT ·AddInt64Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	delta+8(FP), R1
	LDADDLD	R1, (R0), R2
	ADD	R1, R2, R2
	MOVD	R2, ret+16(FP)
	RET

// Add acqrel: LDADDALD + ADD (returns new value)
TEXT ·AddInt64AcqRel(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	delta+8(FP), R1
	LDADDALD	R1, (R0), R2
	ADD	R1, R2, R2
	MOVD	R2, ret+16(FP)
	RET

// =============================================================================
// 64-bit unsigned integer operations (aliases)
// =============================================================================
//
// Relaxed Load/Store are in loadstore_arm64.go

TEXT ·LoadUint64Acquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

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

// =============================================================================
// Uintptr operations (aliases to 64-bit on arm64)
// =============================================================================
//
// Relaxed Load/Store are in loadstore_arm64.go

TEXT ·LoadUintptrAcquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

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

// =============================================================================
// Pointer operations (aliases to 64-bit on arm64)
// =============================================================================
//
// Relaxed Load/Store are in loadstore_arm64.go

TEXT ·LoadPointerAcquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

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

// =============================================================================
// 128-bit operations are in separate files:
//   asm_arm64_128.s      - LL/SC implementation (default, faster on most CPUs)
//   asm_arm64_128_lse2.s - CASP implementation (use -tags=lse2)
// =============================================================================

// =============================================================================
// Barrier operations
// =============================================================================

TEXT ·BarrierAcquire(SB), NOSPLIT, $0-0
	DMB	$0x9
	RET

TEXT ·BarrierRelease(SB), NOSPLIT, $0-0
	DMB	$0xA
	RET

TEXT ·BarrierAcqRel(SB), NOSPLIT, $0-0
	DMB	$0xB
	RET

// =============================================================================
// Bitwise OR operations using LDOR (LSE)
// =============================================================================
//
// Or atomically performs *addr |= mask and returns the old value.
// Uses LDOR instruction (Go's name for ARM LDSET) which sets bits where mask is 1.

// Or32 relaxed: LDORW
TEXT ·OrInt32Relaxed(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	LDORW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Or32 acquire: LDORAW
TEXT ·OrInt32Acquire(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	LDORAW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Or32 release: LDORLW
TEXT ·OrInt32Release(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	LDORLW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Or32 acqrel: LDORALW
TEXT ·OrInt32AcqRel(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	LDORALW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Or64 relaxed: LDORD
TEXT ·OrInt64Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	LDORD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Or64 acquire: LDORAD
TEXT ·OrInt64Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	LDORAD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Or64 release: LDORLD
TEXT ·OrInt64Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	LDORLD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Or64 acqrel: LDORALD
TEXT ·OrInt64AcqRel(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	LDORALD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Unsigned Or operations (aliases to signed)
TEXT ·OrUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

TEXT ·OrUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·OrInt32Acquire(SB)

TEXT ·OrUint32Release(SB), NOSPLIT, $0-20
	JMP	·OrInt32Release(SB)

TEXT ·OrUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

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

// =============================================================================
// Bitwise AND operations using LDCLR (LSE)
// =============================================================================
//
// And atomically performs *addr &= mask and returns the old value.
// Uses LDCLR instruction: LDCLR clears bits where operand is 1.
// So And(addr, mask) = LDCLR(addr, ~mask)
// We compute ~mask with MVN (bitwise NOT).

// And32 relaxed: MVN + LDCLRW
TEXT ·AndInt32Relaxed(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	MVNW	R1, R1		// R1 = ~mask
	LDCLRW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// And32 acquire: MVN + LDCLRAW
TEXT ·AndInt32Acquire(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	MVNW	R1, R1		// R1 = ~mask
	LDCLRAW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// And32 release: MVN + LDCLRLW
TEXT ·AndInt32Release(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	MVNW	R1, R1		// R1 = ~mask
	LDCLRLW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// And32 acqrel: MVN + LDCLRALW
TEXT ·AndInt32AcqRel(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	MVNW	R1, R1		// R1 = ~mask
	LDCLRALW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// And64 relaxed: MVN + LDCLRD
TEXT ·AndInt64Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	MVN	R1, R1		// R1 = ~mask
	LDCLRD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// And64 acquire: MVN + LDCLRAD
TEXT ·AndInt64Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	MVN	R1, R1		// R1 = ~mask
	LDCLRAD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// And64 release: MVN + LDCLRLD
TEXT ·AndInt64Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	MVN	R1, R1		// R1 = ~mask
	LDCLRLD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// And64 acqrel: MVN + LDCLRALD
TEXT ·AndInt64AcqRel(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	MVN	R1, R1		// R1 = ~mask
	LDCLRALD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Unsigned And operations (aliases to signed)
TEXT ·AndUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

TEXT ·AndUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·AndInt32Acquire(SB)

TEXT ·AndUint32Release(SB), NOSPLIT, $0-20
	JMP	·AndInt32Release(SB)

TEXT ·AndUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

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

// =============================================================================
// Bitwise XOR operations using LDEOR (LSE)
// =============================================================================
//
// Xor atomically performs *addr ^= mask and returns the old value.
// Uses LDEOR instruction which XORs bits.

// Xor32 relaxed: LDEORW
TEXT ·XorInt32Relaxed(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	LDEORW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Xor32 acquire: LDEORAW
TEXT ·XorInt32Acquire(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	LDEORAW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Xor32 release: LDEORLW
TEXT ·XorInt32Release(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	LDEORLW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Xor32 acqrel: LDEORALW
TEXT ·XorInt32AcqRel(SB), NOSPLIT, $0-20
	MOVD	addr+0(FP), R0
	MOVW	mask+8(FP), R1
	LDEORALW	R1, (R0), R2
	MOVW	R2, ret+16(FP)
	RET

// Xor64 relaxed: LDEORD
TEXT ·XorInt64Relaxed(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	LDEORD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Xor64 acquire: LDEORAD
TEXT ·XorInt64Acquire(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	LDEORAD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Xor64 release: LDEORLD
TEXT ·XorInt64Release(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	LDEORLD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Xor64 acqrel: LDEORALD
TEXT ·XorInt64AcqRel(SB), NOSPLIT, $0-24
	MOVD	addr+0(FP), R0
	MOVD	mask+8(FP), R1
	LDEORALD	R1, (R0), R2
	MOVD	R2, ret+16(FP)
	RET

// Unsigned Xor operations (aliases to signed)
TEXT ·XorUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

TEXT ·XorUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·XorInt32Acquire(SB)

TEXT ·XorUint32Release(SB), NOSPLIT, $0-20
	JMP	·XorInt32Release(SB)

TEXT ·XorUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

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
