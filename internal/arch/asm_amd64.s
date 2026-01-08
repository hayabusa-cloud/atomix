// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build amd64

#include "textflag.h"

// On x86-64 TSO, all memory orderings are equivalent for LOCK-prefixed
// instructions. We implement one canonical version and use JMP for aliases.

// =============================================================================
// 32-bit signed integer operations
// =============================================================================
//
// Load/Store operations are implemented as pure Go in loadstore_amd64.go
// for inlining optimization. Only RMW operations below require assembly.

// Swap: XCHG is implicitly locked
TEXT ·SwapInt32AcqRel(SB), NOSPLIT, $0-20
	MOVQ	addr+0(FP), DX
	MOVL	new+8(FP), AX
	XCHGL	AX, (DX)
	MOVL	AX, ret+16(FP)
	RET

TEXT ·SwapInt32Relaxed(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

TEXT ·SwapInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

TEXT ·SwapInt32Release(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

// CAS: LOCK CMPXCHG
TEXT ·CasInt32AcqRel(SB), NOSPLIT, $0-17
	MOVQ	addr+0(FP), DX
	MOVL	old+8(FP), AX
	MOVL	new+12(FP), CX
	LOCK
	CMPXCHGL	CX, (DX)
	SETEQ	ret+16(FP)
	RET

TEXT ·CasInt32Relaxed(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

TEXT ·CasInt32Acquire(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

TEXT ·CasInt32Release(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

// CompareExchange: LOCK CMPXCHG (returns old value)
TEXT ·CaxInt32AcqRel(SB), NOSPLIT, $0-20
	MOVQ	addr+0(FP), DX
	MOVL	old+8(FP), AX
	MOVL	new+12(FP), CX
	LOCK
	CMPXCHGL	CX, (DX)
	MOVL	AX, ret+16(FP)
	RET

TEXT ·CaxInt32Relaxed(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

TEXT ·CaxInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

TEXT ·CaxInt32Release(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

// Add: LOCK XADD + ADD (returns new value)
TEXT ·AddInt32AcqRel(SB), NOSPLIT, $0-20
	MOVQ	addr+0(FP), DX
	MOVL	delta+8(FP), AX
	LOCK
	XADDL	AX, (DX)
	ADDL	delta+8(FP), AX
	MOVL	AX, ret+16(FP)
	RET

TEXT ·AddInt32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

TEXT ·AddInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

TEXT ·AddInt32Release(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

// =============================================================================
// 32-bit unsigned integer operations (aliases to signed)
// =============================================================================
//
// Load/Store operations are implemented as pure Go in loadstore_amd64.go

TEXT ·SwapUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

TEXT ·SwapUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

TEXT ·SwapUint32Release(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

TEXT ·SwapUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·SwapInt32AcqRel(SB)

TEXT ·CasUint32Relaxed(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

TEXT ·CasUint32Acquire(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

TEXT ·CasUint32Release(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

TEXT ·CasUint32AcqRel(SB), NOSPLIT, $0-17
	JMP	·CasInt32AcqRel(SB)

TEXT ·CaxUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

TEXT ·CaxUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

TEXT ·CaxUint32Release(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

TEXT ·CaxUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·CaxInt32AcqRel(SB)

TEXT ·AddUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

TEXT ·AddUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

TEXT ·AddUint32Release(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

TEXT ·AddUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AddInt32AcqRel(SB)

// =============================================================================
// 64-bit signed integer operations
// =============================================================================
//
// Load/Store operations are implemented as pure Go in loadstore_amd64.go

TEXT ·SwapInt64AcqRel(SB), NOSPLIT, $0-24
	MOVQ	addr+0(FP), DX
	MOVQ	new+8(FP), AX
	XCHGQ	AX, (DX)
	MOVQ	AX, ret+16(FP)
	RET

TEXT ·SwapInt64Relaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapInt64Release(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·CasInt64AcqRel(SB), NOSPLIT, $0-25
	MOVQ	addr+0(FP), DX
	MOVQ	old+8(FP), AX
	MOVQ	new+16(FP), CX
	LOCK
	CMPXCHGQ	CX, (DX)
	SETEQ	ret+24(FP)
	RET

TEXT ·CasInt64Relaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasInt64Acquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasInt64Release(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CaxInt64AcqRel(SB), NOSPLIT, $0-32
	MOVQ	addr+0(FP), DX
	MOVQ	old+8(FP), AX
	MOVQ	new+16(FP), CX
	LOCK
	CMPXCHGQ	CX, (DX)
	MOVQ	AX, ret+24(FP)
	RET

TEXT ·CaxInt64Relaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxInt64Acquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxInt64Release(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

// Add: LOCK XADD + ADD (returns new value)
TEXT ·AddInt64AcqRel(SB), NOSPLIT, $0-24
	MOVQ	addr+0(FP), DX
	MOVQ	delta+8(FP), AX
	LOCK
	XADDQ	AX, (DX)
	ADDQ	delta+8(FP), AX
	MOVQ	AX, ret+16(FP)
	RET

TEXT ·AddInt64Relaxed(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

TEXT ·AddInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

TEXT ·AddInt64Release(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

// =============================================================================
// 64-bit unsigned integer operations (aliases to signed)
// =============================================================================
//
// Load/Store operations are implemented as pure Go in loadstore_amd64.go

TEXT ·SwapUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapUint64Release(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·CasUint64Relaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasUint64Acquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasUint64Release(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasUint64AcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CaxUint64Relaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxUint64Acquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxUint64Release(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxUint64AcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·AddUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

TEXT ·AddUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

TEXT ·AddUint64Release(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

TEXT ·AddUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

// =============================================================================
// Uintptr operations (aliases to 64-bit on amd64)
// =============================================================================
//
// Load/Store operations are implemented as pure Go in loadstore_amd64.go

TEXT ·SwapUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·CasUintptrRelaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasUintptrAcquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasUintptrRelease(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasUintptrAcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CaxUintptrRelaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxUintptrAcquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxUintptrRelease(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxUintptrAcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·AddUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

TEXT ·AddUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

TEXT ·AddUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

TEXT ·AddUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·AddInt64AcqRel(SB)

// =============================================================================
// Pointer operations (aliases to 64-bit on amd64)
// =============================================================================
//
// Load/Store operations are implemented as pure Go in loadstore_amd64.go

TEXT ·SwapPointerRelaxed(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapPointerAcquire(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapPointerRelease(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·SwapPointerAcqRel(SB), NOSPLIT, $0-24
	JMP	·SwapInt64AcqRel(SB)

TEXT ·CasPointerRelaxed(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasPointerAcquire(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasPointerRelease(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CasPointerAcqRel(SB), NOSPLIT, $0-25
	JMP	·CasInt64AcqRel(SB)

TEXT ·CaxPointerRelaxed(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxPointerAcquire(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxPointerRelease(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

TEXT ·CaxPointerAcqRel(SB), NOSPLIT, $0-32
	JMP	·CaxInt64AcqRel(SB)

// =============================================================================
// 128-bit operations
// =============================================================================

// Load 128-bit: Use LOCK CMPXCHG16B with same old/new to read atomically
TEXT ·LoadUint128Relaxed(SB), NOSPLIT, $16-24
	MOVQ	BX, 0(SP)
	MOVQ	addr+0(FP), DI
	XORQ	AX, AX
	XORQ	DX, DX
	XORQ	BX, BX
	XORQ	CX, CX
	LOCK
	CMPXCHG16B	(DI)
	MOVQ	AX, lo+8(FP)
	MOVQ	DX, hi+16(FP)
	MOVQ	0(SP), BX
	RET

TEXT ·LoadUint128Acquire(SB), NOSPLIT, $0-24
	JMP	·LoadUint128Relaxed(SB)

// Store 128-bit: CAS loop until success
TEXT ·StoreUint128Relaxed(SB), NOSPLIT, $16-24
	MOVQ	BX, 0(SP)
	MOVQ	addr+0(FP), DI
	MOVQ	lo+8(FP), BX
	MOVQ	hi+16(FP), CX
	// Load current value
	XORQ	AX, AX
	XORQ	DX, DX
	LOCK
	CMPXCHG16B	(DI)
	JE	store128_done
store128_retry:
	LOCK
	CMPXCHG16B	(DI)
	JNE	store128_retry
store128_done:
	MOVQ	0(SP), BX
	RET

TEXT ·StoreUint128Release(SB), NOSPLIT, $0-24
	JMP	·StoreUint128Relaxed(SB)

// Swap 128-bit: CAS loop until success, returns old value
// Frame layout: addr+0, newLo+8, newHi+16, oldLo+24, oldHi+32 = 40 bytes
TEXT ·SwapUint128AcqRel(SB), NOSPLIT, $16-40
	MOVQ	BX, 0(SP)
	MOVQ	addr+0(FP), DI
	MOVQ	newLo+8(FP), BX
	MOVQ	newHi+16(FP), CX
	// Load current value
	XORQ	AX, AX
	XORQ	DX, DX
	LOCK
	CMPXCHG16B	(DI)
	JE	swap128_done
swap128_retry:
	LOCK
	CMPXCHG16B	(DI)
	JNE	swap128_retry
swap128_done:
	MOVQ	AX, oldLo+24(FP)
	MOVQ	DX, oldHi+32(FP)
	MOVQ	0(SP), BX
	RET

TEXT ·SwapUint128Relaxed(SB), NOSPLIT, $0-40
	JMP	·SwapUint128AcqRel(SB)

TEXT ·SwapUint128Acquire(SB), NOSPLIT, $0-40
	JMP	·SwapUint128AcqRel(SB)

TEXT ·SwapUint128Release(SB), NOSPLIT, $0-40
	JMP	·SwapUint128AcqRel(SB)

// CAS 128-bit: LOCK CMPXCHG16B
TEXT ·CasUint128AcqRel(SB), NOSPLIT, $16-41
	MOVQ	BX, 0(SP)
	MOVQ	addr+0(FP), DI
	MOVQ	oldLo+8(FP), AX
	MOVQ	oldHi+16(FP), DX
	MOVQ	newLo+24(FP), BX
	MOVQ	newHi+32(FP), CX
	LOCK
	CMPXCHG16B	(DI)
	SETEQ	ret+40(FP)
	MOVQ	0(SP), BX
	RET

TEXT ·CasUint128Relaxed(SB), NOSPLIT, $0-41
	JMP	·CasUint128AcqRel(SB)

TEXT ·CasUint128Acquire(SB), NOSPLIT, $0-41
	JMP	·CasUint128AcqRel(SB)

TEXT ·CasUint128Release(SB), NOSPLIT, $0-41
	JMP	·CasUint128AcqRel(SB)

// CompareExchange 128-bit: returns old value
TEXT ·CaxUint128AcqRel(SB), NOSPLIT, $16-56
	MOVQ	BX, 0(SP)
	MOVQ	addr+0(FP), DI
	MOVQ	oldLo+8(FP), AX
	MOVQ	oldHi+16(FP), DX
	MOVQ	newLo+24(FP), BX
	MOVQ	newHi+32(FP), CX
	LOCK
	CMPXCHG16B	(DI)
	MOVQ	AX, lo+40(FP)
	MOVQ	DX, hi+48(FP)
	MOVQ	0(SP), BX
	RET

TEXT ·CaxUint128Relaxed(SB), NOSPLIT, $0-56
	JMP	·CaxUint128AcqRel(SB)

TEXT ·CaxUint128Acquire(SB), NOSPLIT, $0-56
	JMP	·CaxUint128AcqRel(SB)

TEXT ·CaxUint128Release(SB), NOSPLIT, $0-56
	JMP	·CaxUint128AcqRel(SB)

// =============================================================================
// Bitwise operations (And, Or, Xor) using CAS loops
// =============================================================================
//
// x86-64 has LOCK AND/OR/XOR but they don't return the old value.
// We use CAS loops to atomically read-modify-write and return old value.

// And 32-bit: CAS loop
TEXT ·AndInt32AcqRel(SB), NOSPLIT, $0-20
	MOVQ	addr+0(FP), DX
	MOVL	mask+8(FP), SI
and32_retry:
	MOVL	(DX), AX
	MOVL	AX, CX
	ANDL	SI, CX
	LOCK
	CMPXCHGL	CX, (DX)
	JNE	and32_retry
	MOVL	AX, ret+16(FP)
	RET

TEXT ·AndInt32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

TEXT ·AndInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

TEXT ·AndInt32Release(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

TEXT ·AndUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

TEXT ·AndUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

TEXT ·AndUint32Release(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

TEXT ·AndUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·AndInt32AcqRel(SB)

// And 64-bit: CAS loop
TEXT ·AndInt64AcqRel(SB), NOSPLIT, $0-24
	MOVQ	addr+0(FP), DX
	MOVQ	mask+8(FP), SI
and64_retry:
	MOVQ	(DX), AX
	MOVQ	AX, CX
	ANDQ	SI, CX
	LOCK
	CMPXCHGQ	CX, (DX)
	JNE	and64_retry
	MOVQ	AX, ret+16(FP)
	RET

TEXT ·AndInt64Relaxed(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndInt64Release(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndUint64Release(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

// And uintptr (alias to 64-bit on amd64)
TEXT ·AndUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

TEXT ·AndUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·AndInt64AcqRel(SB)

// Or 32-bit: CAS loop
TEXT ·OrInt32AcqRel(SB), NOSPLIT, $0-20
	MOVQ	addr+0(FP), DX
	MOVL	mask+8(FP), SI
or32_retry:
	MOVL	(DX), AX
	MOVL	AX, CX
	ORL	SI, CX
	LOCK
	CMPXCHGL	CX, (DX)
	JNE	or32_retry
	MOVL	AX, ret+16(FP)
	RET

TEXT ·OrInt32Relaxed(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

TEXT ·OrInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

TEXT ·OrInt32Release(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

TEXT ·OrUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

TEXT ·OrUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

TEXT ·OrUint32Release(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

TEXT ·OrUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·OrInt32AcqRel(SB)

// Or 64-bit: CAS loop
TEXT ·OrInt64AcqRel(SB), NOSPLIT, $0-24
	MOVQ	addr+0(FP), DX
	MOVQ	mask+8(FP), SI
or64_retry:
	MOVQ	(DX), AX
	MOVQ	AX, CX
	ORQ	SI, CX
	LOCK
	CMPXCHGQ	CX, (DX)
	JNE	or64_retry
	MOVQ	AX, ret+16(FP)
	RET

TEXT ·OrInt64Relaxed(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrInt64Release(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrUint64Release(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

// Or uintptr (alias to 64-bit on amd64)
TEXT ·OrUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

TEXT ·OrUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·OrInt64AcqRel(SB)

// Xor 32-bit: CAS loop
TEXT ·XorInt32AcqRel(SB), NOSPLIT, $0-20
	MOVQ	addr+0(FP), DX
	MOVL	mask+8(FP), SI
xor32_retry:
	MOVL	(DX), AX
	MOVL	AX, CX
	XORL	SI, CX
	LOCK
	CMPXCHGL	CX, (DX)
	JNE	xor32_retry
	MOVL	AX, ret+16(FP)
	RET

TEXT ·XorInt32Relaxed(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

TEXT ·XorInt32Acquire(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

TEXT ·XorInt32Release(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

TEXT ·XorUint32Relaxed(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

TEXT ·XorUint32Acquire(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

TEXT ·XorUint32Release(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

TEXT ·XorUint32AcqRel(SB), NOSPLIT, $0-20
	JMP	·XorInt32AcqRel(SB)

// Xor 64-bit: CAS loop
TEXT ·XorInt64AcqRel(SB), NOSPLIT, $0-24
	MOVQ	addr+0(FP), DX
	MOVQ	mask+8(FP), SI
xor64_retry:
	MOVQ	(DX), AX
	MOVQ	AX, CX
	XORQ	SI, CX
	LOCK
	CMPXCHGQ	CX, (DX)
	JNE	xor64_retry
	MOVQ	AX, ret+16(FP)
	RET

TEXT ·XorInt64Relaxed(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorInt64Acquire(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorInt64Release(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorUint64Relaxed(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorUint64Acquire(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorUint64Release(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorUint64AcqRel(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

// Xor uintptr (alias to 64-bit on amd64)
TEXT ·XorUintptrRelaxed(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorUintptrAcquire(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorUintptrRelease(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

TEXT ·XorUintptrAcqRel(SB), NOSPLIT, $0-24
	JMP	·XorInt64AcqRel(SB)

// =============================================================================
// Barrier operations
// =============================================================================

// On x86-64 TSO, acquire and release are compiler barriers only.
// MFENCE is only needed for sequential consistency.
TEXT ·BarrierAcquire(SB), NOSPLIT, $0-0
	RET

TEXT ·BarrierRelease(SB), NOSPLIT, $0-0
	RET

// Full fence using MFENCE
TEXT ·BarrierAcqRel(SB), NOSPLIT, $0-0
	MFENCE
	RET
