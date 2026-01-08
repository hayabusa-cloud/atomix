// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build riscv64

package arch

import "unsafe"

// RISC-V 64-bit atomic operations.
//
// RISC-V has a weak memory model and uses LR/SC (load-reserved/store-conditional)
// for atomic read-modify-write operations. Memory ordering is controlled by
// acquire/release bits on LR/SC and by FENCE instructions.
//
// Implementation strategy:
//   - Load: LD with optional FENCE for acquire
//   - Store: SD with optional FENCE for release
//   - RMW ops: LR/SC loops or AMO instructions with ordering bits
//   - 128-bit: Emulated via LR/SC on low 64 bits with full barriers
//
// RISC-V atomic instructions:
//   - LR.D/LR.W: Load-reserved (doubleword/word)
//   - SC.D/SC.W: Store-conditional (doubleword/word)
//   - AMOADD/AMOSWAP: Atomic memory operations with .aq/.rl/.aqrl suffixes
//   - FENCE: Memory barrier (FENCE r,w / FENCE rw,rw)

// =============================================================================
// 32-bit Signed Integer Operations
// =============================================================================

// All Load/Store require assembly for proper fence handling on RISC-V.
//
//go:noescape
func LoadInt32Relaxed(addr *int32) int32

//go:noescape
func LoadInt32Acquire(addr *int32) int32

//go:noescape
func StoreInt32Relaxed(addr *int32, val int32)

//go:noescape
func StoreInt32Release(addr *int32, val int32)

// Swap atomically stores new and returns the old value.
//
//go:noescape
func SwapInt32Relaxed(addr *int32, new int32) int32

//go:noescape
func SwapInt32Acquire(addr *int32, new int32) int32

//go:noescape
func SwapInt32Release(addr *int32, new int32) int32

//go:noescape
func SwapInt32AcqRel(addr *int32, new int32) int32

// Cas atomically compares *addr with old and swaps if equal.
// Returns true if the swap occurred.
//
//go:noescape
func CasInt32Relaxed(addr *int32, old, new int32) bool

//go:noescape
func CasInt32Acquire(addr *int32, old, new int32) bool

//go:noescape
func CasInt32Release(addr *int32, old, new int32) bool

//go:noescape
func CasInt32AcqRel(addr *int32, old, new int32) bool

// Cax atomically compares *addr with old and swaps if equal.
// Returns the previous value (compare-exchange).
//
//go:noescape
func CaxInt32Relaxed(addr *int32, old, new int32) int32

//go:noescape
func CaxInt32Acquire(addr *int32, old, new int32) int32

//go:noescape
func CaxInt32Release(addr *int32, old, new int32) int32

//go:noescape
func CaxInt32AcqRel(addr *int32, old, new int32) int32

// Add atomically adds delta and returns the new value.
//
//go:noescape
func AddInt32Relaxed(addr *int32, delta int32) int32

//go:noescape
func AddInt32Acquire(addr *int32, delta int32) int32

//go:noescape
func AddInt32Release(addr *int32, delta int32) int32

//go:noescape
func AddInt32AcqRel(addr *int32, delta int32) int32

// =============================================================================
// 32-bit Unsigned Integer Operations
// =============================================================================

//go:noescape
func LoadUint32Relaxed(addr *uint32) uint32

//go:noescape
func LoadUint32Acquire(addr *uint32) uint32

//go:noescape
func StoreUint32Relaxed(addr *uint32, val uint32)

//go:noescape
func StoreUint32Release(addr *uint32, val uint32)

//go:noescape
func SwapUint32Relaxed(addr *uint32, new uint32) uint32

//go:noescape
func SwapUint32Acquire(addr *uint32, new uint32) uint32

//go:noescape
func SwapUint32Release(addr *uint32, new uint32) uint32

//go:noescape
func SwapUint32AcqRel(addr *uint32, new uint32) uint32

//go:noescape
func CasUint32Relaxed(addr *uint32, old, new uint32) bool

//go:noescape
func CasUint32Acquire(addr *uint32, old, new uint32) bool

//go:noescape
func CasUint32Release(addr *uint32, old, new uint32) bool

//go:noescape
func CasUint32AcqRel(addr *uint32, old, new uint32) bool

//go:noescape
func CaxUint32Relaxed(addr *uint32, old, new uint32) uint32

//go:noescape
func CaxUint32Acquire(addr *uint32, old, new uint32) uint32

//go:noescape
func CaxUint32Release(addr *uint32, old, new uint32) uint32

//go:noescape
func CaxUint32AcqRel(addr *uint32, old, new uint32) uint32

//go:noescape
func AddUint32Relaxed(addr *uint32, delta uint32) uint32

//go:noescape
func AddUint32Acquire(addr *uint32, delta uint32) uint32

//go:noescape
func AddUint32Release(addr *uint32, delta uint32) uint32

//go:noescape
func AddUint32AcqRel(addr *uint32, delta uint32) uint32

// =============================================================================
// 64-bit Signed Integer Operations
// =============================================================================

//go:noescape
func LoadInt64Relaxed(addr *int64) int64

//go:noescape
func LoadInt64Acquire(addr *int64) int64

//go:noescape
func StoreInt64Relaxed(addr *int64, val int64)

//go:noescape
func StoreInt64Release(addr *int64, val int64)

//go:noescape
func SwapInt64Relaxed(addr *int64, new int64) int64

//go:noescape
func SwapInt64Acquire(addr *int64, new int64) int64

//go:noescape
func SwapInt64Release(addr *int64, new int64) int64

//go:noescape
func SwapInt64AcqRel(addr *int64, new int64) int64

//go:noescape
func CasInt64Relaxed(addr *int64, old, new int64) bool

//go:noescape
func CasInt64Acquire(addr *int64, old, new int64) bool

//go:noescape
func CasInt64Release(addr *int64, old, new int64) bool

//go:noescape
func CasInt64AcqRel(addr *int64, old, new int64) bool

//go:noescape
func CaxInt64Relaxed(addr *int64, old, new int64) int64

//go:noescape
func CaxInt64Acquire(addr *int64, old, new int64) int64

//go:noescape
func CaxInt64Release(addr *int64, old, new int64) int64

//go:noescape
func CaxInt64AcqRel(addr *int64, old, new int64) int64

//go:noescape
func AddInt64Relaxed(addr *int64, delta int64) int64

//go:noescape
func AddInt64Acquire(addr *int64, delta int64) int64

//go:noescape
func AddInt64Release(addr *int64, delta int64) int64

//go:noescape
func AddInt64AcqRel(addr *int64, delta int64) int64

// =============================================================================
// 64-bit Unsigned Integer Operations
// =============================================================================

//go:noescape
func LoadUint64Relaxed(addr *uint64) uint64

//go:noescape
func LoadUint64Acquire(addr *uint64) uint64

//go:noescape
func StoreUint64Relaxed(addr *uint64, val uint64)

//go:noescape
func StoreUint64Release(addr *uint64, val uint64)

//go:noescape
func SwapUint64Relaxed(addr *uint64, new uint64) uint64

//go:noescape
func SwapUint64Acquire(addr *uint64, new uint64) uint64

//go:noescape
func SwapUint64Release(addr *uint64, new uint64) uint64

//go:noescape
func SwapUint64AcqRel(addr *uint64, new uint64) uint64

//go:noescape
func CasUint64Relaxed(addr *uint64, old, new uint64) bool

//go:noescape
func CasUint64Acquire(addr *uint64, old, new uint64) bool

//go:noescape
func CasUint64Release(addr *uint64, old, new uint64) bool

//go:noescape
func CasUint64AcqRel(addr *uint64, old, new uint64) bool

//go:noescape
func CaxUint64Relaxed(addr *uint64, old, new uint64) uint64

//go:noescape
func CaxUint64Acquire(addr *uint64, old, new uint64) uint64

//go:noescape
func CaxUint64Release(addr *uint64, old, new uint64) uint64

//go:noescape
func CaxUint64AcqRel(addr *uint64, old, new uint64) uint64

//go:noescape
func AddUint64Relaxed(addr *uint64, delta uint64) uint64

//go:noescape
func AddUint64Acquire(addr *uint64, delta uint64) uint64

//go:noescape
func AddUint64Release(addr *uint64, delta uint64) uint64

//go:noescape
func AddUint64AcqRel(addr *uint64, delta uint64) uint64

// =============================================================================
// Uintptr Operations
// =============================================================================

//go:noescape
func LoadUintptrRelaxed(addr *uintptr) uintptr

//go:noescape
func LoadUintptrAcquire(addr *uintptr) uintptr

//go:noescape
func StoreUintptrRelaxed(addr *uintptr, val uintptr)

//go:noescape
func StoreUintptrRelease(addr *uintptr, val uintptr)

//go:noescape
func SwapUintptrRelaxed(addr *uintptr, new uintptr) uintptr

//go:noescape
func SwapUintptrAcquire(addr *uintptr, new uintptr) uintptr

//go:noescape
func SwapUintptrRelease(addr *uintptr, new uintptr) uintptr

//go:noescape
func SwapUintptrAcqRel(addr *uintptr, new uintptr) uintptr

//go:noescape
func CasUintptrRelaxed(addr *uintptr, old, new uintptr) bool

//go:noescape
func CasUintptrAcquire(addr *uintptr, old, new uintptr) bool

//go:noescape
func CasUintptrRelease(addr *uintptr, old, new uintptr) bool

//go:noescape
func CasUintptrAcqRel(addr *uintptr, old, new uintptr) bool

//go:noescape
func CaxUintptrRelaxed(addr *uintptr, old, new uintptr) uintptr

//go:noescape
func CaxUintptrAcquire(addr *uintptr, old, new uintptr) uintptr

//go:noescape
func CaxUintptrRelease(addr *uintptr, old, new uintptr) uintptr

//go:noescape
func CaxUintptrAcqRel(addr *uintptr, old, new uintptr) uintptr

//go:noescape
func AddUintptrRelaxed(addr *uintptr, delta uintptr) uintptr

//go:noescape
func AddUintptrAcquire(addr *uintptr, delta uintptr) uintptr

//go:noescape
func AddUintptrRelease(addr *uintptr, delta uintptr) uintptr

//go:noescape
func AddUintptrAcqRel(addr *uintptr, delta uintptr) uintptr

// =============================================================================
// Pointer Operations
// =============================================================================

//go:noescape
func LoadPointerRelaxed(addr *unsafe.Pointer) unsafe.Pointer

//go:noescape
func LoadPointerAcquire(addr *unsafe.Pointer) unsafe.Pointer

//go:noescape
func StorePointerRelaxed(addr *unsafe.Pointer, val unsafe.Pointer)

//go:noescape
func StorePointerRelease(addr *unsafe.Pointer, val unsafe.Pointer)

//go:noescape
func SwapPointerRelaxed(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer

//go:noescape
func SwapPointerAcquire(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer

//go:noescape
func SwapPointerRelease(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer

//go:noescape
func SwapPointerAcqRel(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer

//go:noescape
func CasPointerRelaxed(addr *unsafe.Pointer, old, new unsafe.Pointer) bool

//go:noescape
func CasPointerAcquire(addr *unsafe.Pointer, old, new unsafe.Pointer) bool

//go:noescape
func CasPointerRelease(addr *unsafe.Pointer, old, new unsafe.Pointer) bool

//go:noescape
func CasPointerAcqRel(addr *unsafe.Pointer, old, new unsafe.Pointer) bool

//go:noescape
func CaxPointerRelaxed(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer

//go:noescape
func CaxPointerAcquire(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer

//go:noescape
func CaxPointerRelease(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer

//go:noescape
func CaxPointerAcqRel(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer

// =============================================================================
// 128-bit Operations
// =============================================================================

// 128-bit operations are emulated on RISC-V using LR/SC on the low 64 bits
// with full memory barriers. Requires 16-byte alignment.

//go:noescape
func LoadUint128Relaxed(addr *[16]byte) (lo, hi uint64)

//go:noescape
func LoadUint128Acquire(addr *[16]byte) (lo, hi uint64)

//go:noescape
func StoreUint128Relaxed(addr *[16]byte, lo, hi uint64)

//go:noescape
func StoreUint128Release(addr *[16]byte, lo, hi uint64)

//go:noescape
func SwapUint128Relaxed(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)

//go:noescape
func SwapUint128Acquire(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)

//go:noescape
func SwapUint128Release(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)

//go:noescape
func SwapUint128AcqRel(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64)

//go:noescape
func CasUint128Relaxed(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool

//go:noescape
func CasUint128Acquire(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool

//go:noescape
func CasUint128Release(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool

//go:noescape
func CasUint128AcqRel(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool

//go:noescape
func CaxUint128Relaxed(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)

//go:noescape
func CaxUint128Acquire(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)

//go:noescape
func CaxUint128Release(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)

//go:noescape
func CaxUint128AcqRel(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64)

// =============================================================================
// Bitwise Operations (And, Or, Xor)
// =============================================================================

// And atomically performs *addr &= mask and returns the old value.
//
//go:noescape
func AndInt32Relaxed(addr *int32, mask int32) int32

//go:noescape
func AndInt32Acquire(addr *int32, mask int32) int32

//go:noescape
func AndInt32Release(addr *int32, mask int32) int32

//go:noescape
func AndInt32AcqRel(addr *int32, mask int32) int32

//go:noescape
func AndUint32Relaxed(addr *uint32, mask uint32) uint32

//go:noescape
func AndUint32Acquire(addr *uint32, mask uint32) uint32

//go:noescape
func AndUint32Release(addr *uint32, mask uint32) uint32

//go:noescape
func AndUint32AcqRel(addr *uint32, mask uint32) uint32

//go:noescape
func AndInt64Relaxed(addr *int64, mask int64) int64

//go:noescape
func AndInt64Acquire(addr *int64, mask int64) int64

//go:noescape
func AndInt64Release(addr *int64, mask int64) int64

//go:noescape
func AndInt64AcqRel(addr *int64, mask int64) int64

//go:noescape
func AndUint64Relaxed(addr *uint64, mask uint64) uint64

//go:noescape
func AndUint64Acquire(addr *uint64, mask uint64) uint64

//go:noescape
func AndUint64Release(addr *uint64, mask uint64) uint64

//go:noescape
func AndUint64AcqRel(addr *uint64, mask uint64) uint64

//go:noescape
func AndUintptrRelaxed(addr *uintptr, mask uintptr) uintptr

//go:noescape
func AndUintptrAcquire(addr *uintptr, mask uintptr) uintptr

//go:noescape
func AndUintptrRelease(addr *uintptr, mask uintptr) uintptr

//go:noescape
func AndUintptrAcqRel(addr *uintptr, mask uintptr) uintptr

// Or atomically performs *addr |= mask and returns the old value.
//
//go:noescape
func OrInt32Relaxed(addr *int32, mask int32) int32

//go:noescape
func OrInt32Acquire(addr *int32, mask int32) int32

//go:noescape
func OrInt32Release(addr *int32, mask int32) int32

//go:noescape
func OrInt32AcqRel(addr *int32, mask int32) int32

//go:noescape
func OrUint32Relaxed(addr *uint32, mask uint32) uint32

//go:noescape
func OrUint32Acquire(addr *uint32, mask uint32) uint32

//go:noescape
func OrUint32Release(addr *uint32, mask uint32) uint32

//go:noescape
func OrUint32AcqRel(addr *uint32, mask uint32) uint32

//go:noescape
func OrInt64Relaxed(addr *int64, mask int64) int64

//go:noescape
func OrInt64Acquire(addr *int64, mask int64) int64

//go:noescape
func OrInt64Release(addr *int64, mask int64) int64

//go:noescape
func OrInt64AcqRel(addr *int64, mask int64) int64

//go:noescape
func OrUint64Relaxed(addr *uint64, mask uint64) uint64

//go:noescape
func OrUint64Acquire(addr *uint64, mask uint64) uint64

//go:noescape
func OrUint64Release(addr *uint64, mask uint64) uint64

//go:noescape
func OrUint64AcqRel(addr *uint64, mask uint64) uint64

//go:noescape
func OrUintptrRelaxed(addr *uintptr, mask uintptr) uintptr

//go:noescape
func OrUintptrAcquire(addr *uintptr, mask uintptr) uintptr

//go:noescape
func OrUintptrRelease(addr *uintptr, mask uintptr) uintptr

//go:noescape
func OrUintptrAcqRel(addr *uintptr, mask uintptr) uintptr

// Xor atomically performs *addr ^= mask and returns the old value.
//
//go:noescape
func XorInt32Relaxed(addr *int32, mask int32) int32

//go:noescape
func XorInt32Acquire(addr *int32, mask int32) int32

//go:noescape
func XorInt32Release(addr *int32, mask int32) int32

//go:noescape
func XorInt32AcqRel(addr *int32, mask int32) int32

//go:noescape
func XorUint32Relaxed(addr *uint32, mask uint32) uint32

//go:noescape
func XorUint32Acquire(addr *uint32, mask uint32) uint32

//go:noescape
func XorUint32Release(addr *uint32, mask uint32) uint32

//go:noescape
func XorUint32AcqRel(addr *uint32, mask uint32) uint32

//go:noescape
func XorInt64Relaxed(addr *int64, mask int64) int64

//go:noescape
func XorInt64Acquire(addr *int64, mask int64) int64

//go:noescape
func XorInt64Release(addr *int64, mask int64) int64

//go:noescape
func XorInt64AcqRel(addr *int64, mask int64) int64

//go:noescape
func XorUint64Relaxed(addr *uint64, mask uint64) uint64

//go:noescape
func XorUint64Acquire(addr *uint64, mask uint64) uint64

//go:noescape
func XorUint64Release(addr *uint64, mask uint64) uint64

//go:noescape
func XorUint64AcqRel(addr *uint64, mask uint64) uint64

//go:noescape
func XorUintptrRelaxed(addr *uintptr, mask uintptr) uintptr

//go:noescape
func XorUintptrAcquire(addr *uintptr, mask uintptr) uintptr

//go:noescape
func XorUintptrRelease(addr *uintptr, mask uintptr) uintptr

//go:noescape
func XorUintptrAcqRel(addr *uintptr, mask uintptr) uintptr

// =============================================================================
// Memory Barriers
// =============================================================================

// FENCE instructions with different read/write combinations.

//go:noescape
func BarrierAcquire()

//go:noescape
func BarrierRelease()

//go:noescape
func BarrierAcqRel()
