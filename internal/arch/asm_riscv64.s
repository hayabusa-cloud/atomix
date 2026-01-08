// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build riscv64

#include "textflag.h"

// RISC-V Atomic Memory Operations:
// AMO instructions: AMOADD, AMOSWAP, AMOAND, AMOOR, AMOXOR
// LR/SC (Load-Reserved/Store-Conditional) for CAS
//
// Go's RISC-V assembler hardcodes ordering bits:
//   LRW/LRD:   funct7=2 → acquire (aq=1)
//   SCW/SCD:   funct7=1 → release (rl=1)
//   AMOSWAP*:  funct7=3 → acquire-release (aq=1, rl=1)
//   AMOADD*:   funct7=3 → acquire-release (aq=1, rl=1)
//
// Therefore:
//   - Swap/Add: All variants alias to Relaxed (AMOSWAP/AMOADD already has aqrl)
//   - CAS/Cax: All variants alias to Relaxed (LR+SC provides aq+rl)
//   - Load: Acquire uses FENCE after (plain load lacks ordering)
//   - Store: Release uses FENCE before (plain store lacks ordering)

// ============================================================================
// 32-bit Signed Integer Operations
// ============================================================================

// func LoadInt32Relaxed(addr *int32) int32
TEXT ·LoadInt32Relaxed(SB), NOSPLIT, $0-12
	MOV	addr+0(FP), A0
	MOVW	(A0), A0
	MOVW	A0, ret+8(FP)
	RET

// func LoadInt32Acquire(addr *int32) int32
TEXT ·LoadInt32Acquire(SB), NOSPLIT, $0-12
	MOV	addr+0(FP), A0
	MOVW	(A0), A0
	FENCE
	MOVW	A0, ret+8(FP)
	RET

// func StoreInt32Relaxed(addr *int32, val int32)
TEXT ·StoreInt32Relaxed(SB), NOSPLIT, $0-12
	MOV	addr+0(FP), A0
	MOVW	val+8(FP), A1
	MOVW	A1, (A0)
	RET

// func StoreInt32Release(addr *int32, val int32)
TEXT ·StoreInt32Release(SB), NOSPLIT, $0-12
	MOV	addr+0(FP), A0
	MOVW	val+8(FP), A1
	FENCE
	MOVW	A1, (A0)
	RET

// func SwapInt32Relaxed(addr *int32, new int32) int32
TEXT ·SwapInt32Relaxed(SB), NOSPLIT, $0-20
	MOV	addr+0(FP), A0
	MOVW	new+8(FP), A1
	AMOSWAPW	A1, (A0), A0
	MOVW	A0, ret+16(FP)
	RET

// func SwapInt32Acquire(addr *int32, new int32) int32
// AMOSWAPW already has aqrl, no FENCE needed
TEXT ·SwapInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Relaxed(SB)

// func SwapInt32Release(addr *int32, new int32) int32
// AMOSWAPW already has aqrl, no FENCE needed
TEXT ·SwapInt32Release(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Relaxed(SB)

// func SwapInt32AcqRel(addr *int32, new int32) int32
// AMOSWAPW already has aqrl, no FENCE needed
TEXT ·SwapInt32AcqRel(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Relaxed(SB)

// func CasInt32Relaxed(addr *int32, old, new int32) bool
TEXT ·CasInt32Relaxed(SB), NOSPLIT, $0-17
	MOV	addr+0(FP), A0
	MOVW	old+8(FP), A1
	MOVW	new+12(FP), A2
cas32_relaxed_loop:
	LRW	(A0), A3
	BNE	A3, A1, cas32_relaxed_fail
	SCW	A2, (A0), A4
	BNE	A4, ZERO, cas32_relaxed_loop
	MOV	$1, A0
	MOVB	A0, ret+16(FP)
	RET
cas32_relaxed_fail:
	MOVB	ZERO, ret+16(FP)
	RET

// func CasInt32Acquire(addr *int32, old, new int32) bool
// LRW has acquire, SCW has release, no FENCE needed
TEXT ·CasInt32Acquire(SB), NOSPLIT, $0-17
	JMP	·CasInt32Relaxed(SB)

// func CasInt32Release(addr *int32, old, new int32) bool
// LRW has acquire, SCW has release, no FENCE needed
TEXT ·CasInt32Release(SB), NOSPLIT, $0-17
	JMP	·CasInt32Relaxed(SB)

// func CasInt32AcqRel(addr *int32, old, new int32) bool
// LRW has acquire, SCW has release, no FENCE needed
TEXT ·CasInt32AcqRel(SB), NOSPLIT, $0-17
	JMP	·CasInt32Relaxed(SB)

// func CaxInt32Relaxed(addr *int32, old, new int32) int32
TEXT ·CaxInt32Relaxed(SB), NOSPLIT, $0-20
	MOV	addr+0(FP), A0
	MOVW	old+8(FP), A1
	MOVW	new+12(FP), A2
cax32_relaxed_loop:
	LRW	(A0), A3
	BNE	A3, A1, cax32_relaxed_done
	SCW	A2, (A0), A4
	BNE	A4, ZERO, cax32_relaxed_loop
cax32_relaxed_done:
	MOVW	A3, ret+16(FP)
	RET

// func CaxInt32Acquire(addr *int32, old, new int32) int32
// LRW has acquire, SCW has release, no FENCE needed
TEXT ·CaxInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Relaxed(SB)

// func CaxInt32Release(addr *int32, old, new int32) int32
// LRW has acquire, SCW has release, no FENCE needed
TEXT ·CaxInt32Release(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Relaxed(SB)

// func CaxInt32AcqRel(addr *int32, old, new int32) int32
// LRW has acquire, SCW has release, no FENCE needed
TEXT ·CaxInt32AcqRel(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Relaxed(SB)

// func AddInt32Relaxed(addr *int32, delta int32) int32
TEXT ·AddInt32Relaxed(SB), NOSPLIT, $0-20
	MOV	addr+0(FP), A0
	MOVW	delta+8(FP), A1
	AMOADDW	A1, (A0), A2
	ADDW	A1, A2, A0
	MOVW	A0, ret+16(FP)
	RET

// func AddInt32Acquire(addr *int32, delta int32) int32
// AMOADDW already has aqrl, no FENCE needed
TEXT ·AddInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·AddInt32Relaxed(SB)

// func AddInt32Release(addr *int32, delta int32) int32
// AMOADDW already has aqrl, no FENCE needed
TEXT ·AddInt32Release(SB), NOSPLIT, $0-20
	JMP	·AddInt32Relaxed(SB)

// func AddInt32AcqRel(addr *int32, delta int32) int32
// AMOADDW already has aqrl, no FENCE needed
TEXT ·AddInt32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AddInt32Relaxed(SB)

// ============================================================================
// 32-bit Unsigned Integer Operations (aliases to signed versions)
// ============================================================================

// func LoadUint32Relaxed(addr *uint32) uint32
TEXT ·LoadUint32Relaxed(SB), NOSPLIT, $0-12
	JMP	·LoadInt32Relaxed(SB)

// func LoadUint32Acquire(addr *uint32) uint32
TEXT ·LoadUint32Acquire(SB), NOSPLIT, $0-12
	JMP	·LoadInt32Acquire(SB)

// func StoreUint32Relaxed(addr *uint32, val uint32)
TEXT ·StoreUint32Relaxed(SB), NOSPLIT, $0-12
	JMP	·StoreInt32Relaxed(SB)

// func StoreUint32Release(addr *uint32, val uint32)
TEXT ·StoreUint32Release(SB), NOSPLIT, $0-12
	JMP	·StoreInt32Release(SB)

// func SwapUint32Relaxed(addr *uint32, new uint32) uint32
TEXT ·SwapUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Relaxed(SB)

// func SwapUint32Acquire(addr *uint32, new uint32) uint32
TEXT ·SwapUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Acquire(SB)

// func SwapUint32Release(addr *uint32, new uint32) uint32
TEXT ·SwapUint32Release(SB), NOSPLIT, $0-20
	JMP	·SwapInt32Release(SB)

// func SwapUint32AcqRel(addr *uint32, new uint32) uint32
TEXT ·SwapUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

// func CasUint32Relaxed(addr *uint32, old, new uint32) bool
TEXT ·CasUint32Relaxed(SB), NOSPLIT, $0-17
	JMP	·CasInt32Relaxed(SB)

// func CasUint32Acquire(addr *uint32, old, new uint32) bool
TEXT ·CasUint32Acquire(SB), NOSPLIT, $0-17
	JMP	·CasInt32Acquire(SB)

// func CasUint32Release(addr *uint32, old, new uint32) bool
TEXT ·CasUint32Release(SB), NOSPLIT, $0-17
	JMP	·CasInt32Release(SB)

// func CasUint32AcqRel(addr *uint32, old, new uint32) bool
TEXT ·CasUint32AcqRel(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

// func CaxUint32Relaxed(addr *uint32, old, new uint32) uint32
TEXT ·CaxUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Relaxed(SB)

// func CaxUint32Acquire(addr *uint32, old, new uint32) uint32
TEXT ·CaxUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Acquire(SB)

// func CaxUint32Release(addr *uint32, old, new uint32) uint32
TEXT ·CaxUint32Release(SB), NOSPLIT, $0-20
	JMP	·CaxInt32Release(SB)

// func CaxUint32AcqRel(addr *uint32, old, new uint32) uint32
TEXT ·CaxUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

// func AddUint32Relaxed(addr *uint32, delta uint32) uint32
TEXT ·AddUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AddInt32Relaxed(SB)

// func AddUint32Acquire(addr *uint32, delta uint32) uint32
TEXT ·AddUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·AddInt32Acquire(SB)

// func AddUint32Release(addr *uint32, delta uint32) uint32
TEXT ·AddUint32Release(SB), NOSPLIT, $0-20
	JMP	·AddInt32Release(SB)

// func AddUint32AcqRel(addr *uint32, delta uint32) uint32
TEXT ·AddUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

// ============================================================================
// 64-bit Signed Integer Operations
// ============================================================================

// func LoadInt64Relaxed(addr *int64) int64
TEXT ·LoadInt64Relaxed(SB), NOSPLIT, $0-16
	MOV	addr+0(FP), A0
	MOV	(A0), A0
	MOV	A0, ret+8(FP)
	RET

// func LoadInt64Acquire(addr *int64) int64
TEXT ·LoadInt64Acquire(SB), NOSPLIT, $0-16
	MOV	addr+0(FP), A0
	MOV	(A0), A0
	FENCE
	MOV	A0, ret+8(FP)
	RET

// func StoreInt64Relaxed(addr *int64, val int64)
TEXT ·StoreInt64Relaxed(SB), NOSPLIT, $0-16
	MOV	addr+0(FP), A0
	MOV	val+8(FP), A1
	MOV	A1, (A0)
	RET

// func StoreInt64Release(addr *int64, val int64)
TEXT ·StoreInt64Release(SB), NOSPLIT, $0-16
	MOV	addr+0(FP), A0
	MOV	val+8(FP), A1
	FENCE
	MOV	A1, (A0)
	RET

// func SwapInt64Relaxed(addr *int64, new int64) int64
TEXT ·SwapInt64Relaxed(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
	MOV	new+8(FP), A1
	AMOSWAPD	A1, (A0), A0
	MOV	A0, ret+16(FP)
	RET

// func SwapInt64Acquire(addr *int64, new int64) int64
// AMOSWAPD already has aqrl, no FENCE needed
TEXT ·SwapInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

// func SwapInt64Release(addr *int64, new int64) int64
// AMOSWAPD already has aqrl, no FENCE needed
TEXT ·SwapInt64Release(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

// func SwapInt64AcqRel(addr *int64, new int64) int64
// AMOSWAPD already has aqrl, no FENCE needed
TEXT ·SwapInt64AcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

// func CasInt64Relaxed(addr *int64, old, new int64) bool
TEXT ·CasInt64Relaxed(SB), NOSPLIT, $0-25
	MOV	addr+0(FP), A0
	MOV	old+8(FP), A1
	MOV	new+16(FP), A2
cas64_relaxed_loop:
	LRD	(A0), A3
	BNE	A3, A1, cas64_relaxed_fail
	SCD	A2, (A0), A4
	BNE	A4, ZERO, cas64_relaxed_loop
	MOV	$1, A0
	MOVB	A0, ret+24(FP)
	RET
cas64_relaxed_fail:
	MOVB	ZERO, ret+24(FP)
	RET

// func CasInt64Acquire(addr *int64, old, new int64) bool
// LRD has acquire, SCD has release, no FENCE needed
TEXT ·CasInt64Acquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

// func CasInt64Release(addr *int64, old, new int64) bool
// LRD has acquire, SCD has release, no FENCE needed
TEXT ·CasInt64Release(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

// func CasInt64AcqRel(addr *int64, old, new int64) bool
// LRD has acquire, SCD has release, no FENCE needed
TEXT ·CasInt64AcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

// func CaxInt64Relaxed(addr *int64, old, new int64) int64
TEXT ·CaxInt64Relaxed(SB), NOSPLIT, $0-32
	MOV	addr+0(FP), A0
	MOV	old+8(FP), A1
	MOV	new+16(FP), A2
cax64_relaxed_loop:
	LRD	(A0), A3
	BNE	A3, A1, cax64_relaxed_done
	SCD	A2, (A0), A4
	BNE	A4, ZERO, cax64_relaxed_loop
cax64_relaxed_done:
	MOV	A3, ret+24(FP)
	RET

// func CaxInt64Acquire(addr *int64, old, new int64) int64
// LRD has acquire, SCD has release, no FENCE needed
TEXT ·CaxInt64Acquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

// func CaxInt64Release(addr *int64, old, new int64) int64
// LRD has acquire, SCD has release, no FENCE needed
TEXT ·CaxInt64Release(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

// func CaxInt64AcqRel(addr *int64, old, new int64) int64
// LRD has acquire, SCD has release, no FENCE needed
TEXT ·CaxInt64AcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

// func AddInt64Relaxed(addr *int64, delta int64) int64
TEXT ·AddInt64Relaxed(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
	MOV	delta+8(FP), A1
	AMOADDD	A1, (A0), A2
	ADD	A1, A2, A0
	MOV	A0, ret+16(FP)
	RET

// func AddInt64Acquire(addr *int64, delta int64) int64
// AMOADDD already has aqrl, no FENCE needed
TEXT ·AddInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·AddInt64Relaxed(SB)

// func AddInt64Release(addr *int64, delta int64) int64
// AMOADDD already has aqrl, no FENCE needed
TEXT ·AddInt64Release(SB), NOSPLIT, $0-24
	JMP	·AddInt64Relaxed(SB)

// func AddInt64AcqRel(addr *int64, delta int64) int64
// AMOADDD already has aqrl, no FENCE needed
TEXT ·AddInt64AcqRel(SB), NOSPLIT, $0-24
	JMP	·AddInt64Relaxed(SB)

// ============================================================================
// 64-bit Unsigned Integer Operations (aliases to signed versions)
// ============================================================================

// func LoadUint64Relaxed(addr *uint64) uint64
TEXT ·LoadUint64Relaxed(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Relaxed(SB)

// func LoadUint64Acquire(addr *uint64) uint64
TEXT ·LoadUint64Acquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

// func StoreUint64Relaxed(addr *uint64, val uint64)
TEXT ·StoreUint64Relaxed(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Relaxed(SB)

// func StoreUint64Release(addr *uint64, val uint64)
TEXT ·StoreUint64Release(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Release(SB)

// func SwapUint64Relaxed(addr *uint64, new uint64) uint64
TEXT ·SwapUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

// func SwapUint64Acquire(addr *uint64, new uint64) uint64
TEXT ·SwapUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Acquire(SB)

// func SwapUint64Release(addr *uint64, new uint64) uint64
TEXT ·SwapUint64Release(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Release(SB)

// func SwapUint64AcqRel(addr *uint64, new uint64) uint64
TEXT ·SwapUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

// func CasUint64Relaxed(addr *uint64, old, new uint64) bool
TEXT ·CasUint64Relaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

// func CasUint64Acquire(addr *uint64, old, new uint64) bool
TEXT ·CasUint64Acquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64Acquire(SB)

// func CasUint64Release(addr *uint64, old, new uint64) bool
TEXT ·CasUint64Release(SB), NOSPLIT, $0-25
	JMP	·CasInt64Release(SB)

// func CasUint64AcqRel(addr *uint64, old, new uint64) bool
TEXT ·CasUint64AcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

// func CaxUint64Relaxed(addr *uint64, old, new uint64) uint64
TEXT ·CaxUint64Relaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

// func CaxUint64Acquire(addr *uint64, old, new uint64) uint64
TEXT ·CaxUint64Acquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Acquire(SB)

// func CaxUint64Release(addr *uint64, old, new uint64) uint64
TEXT ·CaxUint64Release(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Release(SB)

// func CaxUint64AcqRel(addr *uint64, old, new uint64) uint64
TEXT ·CaxUint64AcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

// func AddUint64Relaxed(addr *uint64, delta uint64) uint64
TEXT ·AddUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·AddInt64Relaxed(SB)

// func AddUint64Acquire(addr *uint64, delta uint64) uint64
TEXT ·AddUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·AddInt64Acquire(SB)

// func AddUint64Release(addr *uint64, delta uint64) uint64
TEXT ·AddUint64Release(SB), NOSPLIT, $0-24
	JMP	·AddInt64Release(SB)

// func AddUint64AcqRel(addr *uint64, delta uint64) uint64
TEXT ·AddUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

// ============================================================================
// Uintptr Operations (aliases to 64-bit on riscv64)
// ============================================================================

// func LoadUintptrRelaxed(addr *uintptr) uintptr
TEXT ·LoadUintptrRelaxed(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Relaxed(SB)

// func LoadUintptrAcquire(addr *uintptr) uintptr
TEXT ·LoadUintptrAcquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

// func StoreUintptrRelaxed(addr *uintptr, val uintptr)
TEXT ·StoreUintptrRelaxed(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Relaxed(SB)

// func StoreUintptrRelease(addr *uintptr, val uintptr)
TEXT ·StoreUintptrRelease(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Release(SB)

// func SwapUintptrRelaxed(addr *uintptr, new uintptr) uintptr
TEXT ·SwapUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

// func SwapUintptrAcquire(addr *uintptr, new uintptr) uintptr
TEXT ·SwapUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Acquire(SB)

// func SwapUintptrRelease(addr *uintptr, new uintptr) uintptr
TEXT ·SwapUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Release(SB)

// func SwapUintptrAcqRel(addr *uintptr, new uintptr) uintptr
TEXT ·SwapUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

// func CasUintptrRelaxed(addr *uintptr, old, new uintptr) bool
TEXT ·CasUintptrRelaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

// func CasUintptrAcquire(addr *uintptr, old, new uintptr) bool
TEXT ·CasUintptrAcquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64Acquire(SB)

// func CasUintptrRelease(addr *uintptr, old, new uintptr) bool
TEXT ·CasUintptrRelease(SB), NOSPLIT, $0-25
	JMP	·CasInt64Release(SB)

// func CasUintptrAcqRel(addr *uintptr, old, new uintptr) bool
TEXT ·CasUintptrAcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

// func CaxUintptrRelaxed(addr *uintptr, old, new uintptr) uintptr
TEXT ·CaxUintptrRelaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

// func CaxUintptrAcquire(addr *uintptr, old, new uintptr) uintptr
TEXT ·CaxUintptrAcquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Acquire(SB)

// func CaxUintptrRelease(addr *uintptr, old, new uintptr) uintptr
TEXT ·CaxUintptrRelease(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Release(SB)

// func CaxUintptrAcqRel(addr *uintptr, old, new uintptr) uintptr
TEXT ·CaxUintptrAcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

// func AddUintptrRelaxed(addr *uintptr, delta uintptr) uintptr
TEXT ·AddUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·AddInt64Relaxed(SB)

// func AddUintptrAcquire(addr *uintptr, delta uintptr) uintptr
TEXT ·AddUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·AddInt64Acquire(SB)

// func AddUintptrRelease(addr *uintptr, delta uintptr) uintptr
TEXT ·AddUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·AddInt64Release(SB)

// func AddUintptrAcqRel(addr *uintptr, delta uintptr) uintptr
TEXT ·AddUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

// ============================================================================
// Pointer Operations (aliases to 64-bit on riscv64)
// ============================================================================

// func LoadPointerRelaxed(addr *unsafe.Pointer) unsafe.Pointer
TEXT ·LoadPointerRelaxed(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Relaxed(SB)

// func LoadPointerAcquire(addr *unsafe.Pointer) unsafe.Pointer
TEXT ·LoadPointerAcquire(SB), NOSPLIT, $0-16
	JMP	·LoadInt64Acquire(SB)

// func StorePointerRelaxed(addr *unsafe.Pointer, val unsafe.Pointer)
TEXT ·StorePointerRelaxed(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Relaxed(SB)

// func StorePointerRelease(addr *unsafe.Pointer, val unsafe.Pointer)
TEXT ·StorePointerRelease(SB), NOSPLIT, $0-16
	JMP	·StoreInt64Release(SB)

// func SwapPointerRelaxed(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer
TEXT ·SwapPointerRelaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Relaxed(SB)

// func SwapPointerAcquire(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer
TEXT ·SwapPointerAcquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Acquire(SB)

// func SwapPointerRelease(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer
TEXT ·SwapPointerRelease(SB), NOSPLIT, $0-24
	JMP	·SwapInt64Release(SB)

// func SwapPointerAcqRel(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer
TEXT ·SwapPointerAcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

// func CasPointerRelaxed(addr *unsafe.Pointer, old, new unsafe.Pointer) bool
TEXT ·CasPointerRelaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64Relaxed(SB)

// func CasPointerAcquire(addr *unsafe.Pointer, old, new unsafe.Pointer) bool
TEXT ·CasPointerAcquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64Acquire(SB)

// func CasPointerRelease(addr *unsafe.Pointer, old, new unsafe.Pointer) bool
TEXT ·CasPointerRelease(SB), NOSPLIT, $0-25
	JMP	·CasInt64Release(SB)

// func CasPointerAcqRel(addr *unsafe.Pointer, old, new unsafe.Pointer) bool
TEXT ·CasPointerAcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

// func CaxPointerRelaxed(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer
TEXT ·CaxPointerRelaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Relaxed(SB)

// func CaxPointerAcquire(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer
TEXT ·CaxPointerAcquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Acquire(SB)

// func CaxPointerRelease(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer
TEXT ·CaxPointerRelease(SB), NOSPLIT, $0-32
	JMP	·CaxInt64Release(SB)

// func CaxPointerAcqRel(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer
TEXT ·CaxPointerAcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

// ============================================================================
// 128-bit Operations
// RISC-V does not have native 128-bit atomics.
// We emulate using a spinlock on the low word via LR/SC.
// This is NOT truly atomic 128-bit but provides mutual exclusion.
// The address MUST be 16-byte aligned.
// ============================================================================

// func LoadUint128Relaxed(addr *[16]byte) (lo, hi uint64)
// Emulated: retry-load until consistent
TEXT ·LoadUint128Relaxed(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
load128_relaxed_loop:
	MOV	(A0), A1		// load lo
	MOV	8(A0), A2		// load hi
	MOV	(A0), A3		// reload lo
	BNE	A1, A3, load128_relaxed_loop
	MOV	A1, lo+8(FP)
	MOV	A2, hi+16(FP)
	RET

// func LoadUint128Acquire(addr *[16]byte) (lo, hi uint64)
TEXT ·LoadUint128Acquire(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
load128_acq_loop:
	MOV	(A0), A1
	MOV	8(A0), A2
	MOV	(A0), A3
	BNE	A1, A3, load128_acq_loop
	FENCE
	MOV	A1, lo+8(FP)
	MOV	A2, hi+16(FP)
	RET

// func StoreUint128Relaxed(addr *[16]byte, lo, hi uint64)
// Use LR/SC on lo word as a lock
TEXT ·StoreUint128Relaxed(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
	MOV	lo+8(FP), A1
	MOV	hi+16(FP), A2
store128_relaxed_loop:
	LRD	(A0), A3
	SCD	A1, (A0), A4
	BNE	A4, ZERO, store128_relaxed_loop
	MOV	A2, 8(A0)
	RET

// func StoreUint128Release(addr *[16]byte, lo, hi uint64)
TEXT ·StoreUint128Release(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
	MOV	lo+8(FP), A1
	MOV	hi+16(FP), A2
	FENCE
store128_rel_loop:
	LRD	(A0), A3
	SCD	A1, (A0), A4
	BNE	A4, ZERO, store128_rel_loop
	MOV	A2, 8(A0)
	RET

// func SwapUint128Relaxed(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)
// Use LR/SC on lo word, returns old value
TEXT ·SwapUint128Relaxed(SB), NOSPLIT, $0-40
	MOV	addr+0(FP), A0
	MOV	newLo+8(FP), A1
	MOV	newHi+16(FP), A2
swap128_relaxed_loop:
	LRD	(A0), A3		// load-reserve old lo
	MOV	8(A0), A4		// load old hi
	SCD	A1, (A0), A5		// store-conditional new lo
	BNE	A5, ZERO, swap128_relaxed_loop
	MOV	A2, 8(A0)		// store new hi
	MOV	A3, oldLo+24(FP)
	MOV	A4, oldHi+32(FP)
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
	MOV	addr+0(FP), A0
	MOV	oldLo+8(FP), A1
	MOV	oldHi+16(FP), A2
	MOV	newLo+24(FP), A3
	MOV	newHi+32(FP), A4
cas128_relaxed_loop:
	LRD	(A0), A5		// load-reserve lo
	BNE	A5, A1, cas128_relaxed_fail
	MOV	8(A0), A6		// load hi
	BNE	A6, A2, cas128_relaxed_fail
	SCD	A3, (A0), A7		// store-conditional new lo
	BNE	A7, ZERO, cas128_relaxed_loop
	MOV	A4, 8(A0)		// store new hi
	MOV	$1, A0
	MOVB	A0, ret+40(FP)
	RET
cas128_relaxed_fail:
	MOVB	ZERO, ret+40(FP)
	RET

// func CasUint128Acquire(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool
TEXT ·CasUint128Acquire(SB), NOSPLIT, $0-41
	MOV	addr+0(FP), A0
	MOV	oldLo+8(FP), A1
	MOV	oldHi+16(FP), A2
	MOV	newLo+24(FP), A3
	MOV	newHi+32(FP), A4
cas128_acq_loop:
	LRD	(A0), A5
	BNE	A5, A1, cas128_acq_fail
	MOV	8(A0), A6
	BNE	A6, A2, cas128_acq_fail
	SCD	A3, (A0), A7
	BNE	A7, ZERO, cas128_acq_loop
	MOV	A4, 8(A0)
	FENCE
	MOV	$1, A0
	MOVB	A0, ret+40(FP)
	RET
cas128_acq_fail:
	MOVB	ZERO, ret+40(FP)
	RET

// func CasUint128Release(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool
TEXT ·CasUint128Release(SB), NOSPLIT, $0-41
	MOV	addr+0(FP), A0
	MOV	oldLo+8(FP), A1
	MOV	oldHi+16(FP), A2
	MOV	newLo+24(FP), A3
	MOV	newHi+32(FP), A4
	FENCE
cas128_rel_loop:
	LRD	(A0), A5
	BNE	A5, A1, cas128_rel_fail
	MOV	8(A0), A6
	BNE	A6, A2, cas128_rel_fail
	SCD	A3, (A0), A7
	BNE	A7, ZERO, cas128_rel_loop
	MOV	A4, 8(A0)
	MOV	$1, A0
	MOVB	A0, ret+40(FP)
	RET
cas128_rel_fail:
	MOVB	ZERO, ret+40(FP)
	RET

// func CasUint128AcqRel(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool
TEXT ·CasUint128AcqRel(SB), NOSPLIT, $0-41
	MOV	addr+0(FP), A0
	MOV	oldLo+8(FP), A1
	MOV	oldHi+16(FP), A2
	MOV	newLo+24(FP), A3
	MOV	newHi+32(FP), A4
	FENCE
cas128_aqrl_loop:
	LRD	(A0), A5
	BNE	A5, A1, cas128_aqrl_fail
	MOV	8(A0), A6
	BNE	A6, A2, cas128_aqrl_fail
	SCD	A3, (A0), A7
	BNE	A7, ZERO, cas128_aqrl_loop
	MOV	A4, 8(A0)
	FENCE
	MOV	$1, A0
	MOVB	A0, ret+40(FP)
	RET
cas128_aqrl_fail:
	MOVB	ZERO, ret+40(FP)
	RET

// func CaxUint128Relaxed(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)
TEXT ·CaxUint128Relaxed(SB), NOSPLIT, $0-56
	MOV	addr+0(FP), A0
	MOV	oldLo+8(FP), A1
	MOV	oldHi+16(FP), A2
	MOV	newLo+24(FP), A3
	MOV	newHi+32(FP), A4
cax128_relaxed_loop:
	LRD	(A0), A5
	MOV	8(A0), A6
	BNE	A5, A1, cax128_relaxed_done
	BNE	A6, A2, cax128_relaxed_done
	SCD	A3, (A0), A7
	BNE	A7, ZERO, cax128_relaxed_loop
	MOV	A4, 8(A0)
cax128_relaxed_done:
	MOV	A5, lo+40(FP)
	MOV	A6, hi+48(FP)
	RET

// func CaxUint128Acquire(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)
TEXT ·CaxUint128Acquire(SB), NOSPLIT, $0-56
	MOV	addr+0(FP), A0
	MOV	oldLo+8(FP), A1
	MOV	oldHi+16(FP), A2
	MOV	newLo+24(FP), A3
	MOV	newHi+32(FP), A4
cax128_acq_loop:
	LRD	(A0), A5
	MOV	8(A0), A6
	BNE	A5, A1, cax128_acq_done
	BNE	A6, A2, cax128_acq_done
	SCD	A3, (A0), A7
	BNE	A7, ZERO, cax128_acq_loop
	MOV	A4, 8(A0)
cax128_acq_done:
	FENCE
	MOV	A5, lo+40(FP)
	MOV	A6, hi+48(FP)
	RET

// func CaxUint128Release(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)
TEXT ·CaxUint128Release(SB), NOSPLIT, $0-56
	MOV	addr+0(FP), A0
	MOV	oldLo+8(FP), A1
	MOV	oldHi+16(FP), A2
	MOV	newLo+24(FP), A3
	MOV	newHi+32(FP), A4
	FENCE
cax128_rel_loop:
	LRD	(A0), A5
	MOV	8(A0), A6
	BNE	A5, A1, cax128_rel_done
	BNE	A6, A2, cax128_rel_done
	SCD	A3, (A0), A7
	BNE	A7, ZERO, cax128_rel_loop
	MOV	A4, 8(A0)
cax128_rel_done:
	MOV	A5, lo+40(FP)
	MOV	A6, hi+48(FP)
	RET

// func CaxUint128AcqRel(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)
TEXT ·CaxUint128AcqRel(SB), NOSPLIT, $0-56
	MOV	addr+0(FP), A0
	MOV	oldLo+8(FP), A1
	MOV	oldHi+16(FP), A2
	MOV	newLo+24(FP), A3
	MOV	newHi+32(FP), A4
	FENCE
cax128_aqrl_loop:
	LRD	(A0), A5
	MOV	8(A0), A6
	BNE	A5, A1, cax128_aqrl_done
	BNE	A6, A2, cax128_aqrl_done
	SCD	A3, (A0), A7
	BNE	A7, ZERO, cax128_aqrl_loop
	MOV	A4, 8(A0)
cax128_aqrl_done:
	FENCE
	MOV	A5, lo+40(FP)
	MOV	A6, hi+48(FP)
	RET

// ============================================================================
// Bitwise Operations (And, Or, Xor)
// ============================================================================
//
// RISC-V has AMOAND, AMOOR, AMOXOR instructions with acquire-release semantics.
// Like AMOSWAP/AMOADD, Go's assembler emits these with aqrl bits set.

// func AndInt32Relaxed(addr *int32, mask int32) int32
TEXT ·AndInt32Relaxed(SB), NOSPLIT, $0-20
	MOV	addr+0(FP), A0
	MOVW	mask+8(FP), A1
	AMOANDW	A1, (A0), A0
	MOVW	A0, ret+16(FP)
	RET

TEXT ·AndInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

TEXT ·AndInt32Release(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

TEXT ·AndInt32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

TEXT ·AndUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

TEXT ·AndUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

TEXT ·AndUint32Release(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

TEXT ·AndUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AndInt32Relaxed(SB)

// func AndInt64Relaxed(addr *int64, mask int64) int64
TEXT ·AndInt64Relaxed(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
	MOV	mask+8(FP), A1
	AMOANDD	A1, (A0), A0
	MOV	A0, ret+16(FP)
	RET

TEXT ·AndInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndInt64Release(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndInt64AcqRel(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUint64Release(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

TEXT ·AndUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·AndInt64Relaxed(SB)

// func OrInt32Relaxed(addr *int32, mask int32) int32
TEXT ·OrInt32Relaxed(SB), NOSPLIT, $0-20
	MOV	addr+0(FP), A0
	MOVW	mask+8(FP), A1
	AMOORW	A1, (A0), A0
	MOVW	A0, ret+16(FP)
	RET

TEXT ·OrInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

TEXT ·OrInt32Release(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

TEXT ·OrInt32AcqRel(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

TEXT ·OrUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

TEXT ·OrUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

TEXT ·OrUint32Release(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

TEXT ·OrUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·OrInt32Relaxed(SB)

// func OrInt64Relaxed(addr *int64, mask int64) int64
TEXT ·OrInt64Relaxed(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
	MOV	mask+8(FP), A1
	AMOORD	A1, (A0), A0
	MOV	A0, ret+16(FP)
	RET

TEXT ·OrInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrInt64Release(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrInt64AcqRel(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUint64Release(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

TEXT ·OrUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·OrInt64Relaxed(SB)

// func XorInt32Relaxed(addr *int32, mask int32) int32
TEXT ·XorInt32Relaxed(SB), NOSPLIT, $0-20
	MOV	addr+0(FP), A0
	MOVW	mask+8(FP), A1
	AMOXORW	A1, (A0), A0
	MOVW	A0, ret+16(FP)
	RET

TEXT ·XorInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

TEXT ·XorInt32Release(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

TEXT ·XorInt32AcqRel(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

TEXT ·XorUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

TEXT ·XorUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

TEXT ·XorUint32Release(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

TEXT ·XorUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·XorInt32Relaxed(SB)

// func XorInt64Relaxed(addr *int64, mask int64) int64
TEXT ·XorInt64Relaxed(SB), NOSPLIT, $0-24
	MOV	addr+0(FP), A0
	MOV	mask+8(FP), A1
	AMOXORD	A1, (A0), A0
	MOV	A0, ret+16(FP)
	RET

TEXT ·XorInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorInt64Release(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorInt64AcqRel(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUint64Release(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

TEXT ·XorUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·XorInt64Relaxed(SB)

// ============================================================================
// Barrier Operations
// ============================================================================

// func BarrierAcquire()
TEXT ·BarrierAcquire(SB), NOSPLIT, $0-0
	FENCE
	RET

// func BarrierRelease()
TEXT ·BarrierRelease(SB), NOSPLIT, $0-0
	FENCE
	RET

// func BarrierAcqRel()
TEXT ·BarrierAcqRel(SB), NOSPLIT, $0-0
	FENCE
	RET
