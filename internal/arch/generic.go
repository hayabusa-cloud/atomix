// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build !amd64 && !arm64 && !riscv64 && !loong64

package arch

import (
	"sync/atomic"
	"unsafe"
)

// Generic fallback implementation using sync/atomic.
//
// This file provides atomic operations for architectures without optimized
// assembly implementations (e.g., 386, ppc64, s390x, wasm).
//
// Implementation characteristics:
//   - All operations use sync/atomic which provides sequential consistency
//   - Memory ordering variants (Relaxed, Acquire, Release, AcqRel) are
//     all equivalent since sync/atomic doesn't expose weaker orderings
//   - 128-bit operations are NOT atomic (no hardware support)
//
// WARNING: 128-bit operations on unsupported architectures are implemented
// as non-atomic read/write sequences. They must not be used in concurrent
// contexts without external synchronization.

// =============================================================================
// 32-bit Signed Integer Operations
// =============================================================================

func LoadInt32Relaxed(addr *int32) int32 {
	return atomic.LoadInt32(addr)
}

func LoadInt32Acquire(addr *int32) int32 {
	return atomic.LoadInt32(addr)
}

func StoreInt32Relaxed(addr *int32, val int32) {
	atomic.StoreInt32(addr, val)
}

func StoreInt32Release(addr *int32, val int32) {
	atomic.StoreInt32(addr, val)
}

func SwapInt32Relaxed(addr *int32, new int32) int32 {
	return atomic.SwapInt32(addr, new)
}

func SwapInt32Acquire(addr *int32, new int32) int32 {
	return atomic.SwapInt32(addr, new)
}

func SwapInt32Release(addr *int32, new int32) int32 {
	return atomic.SwapInt32(addr, new)
}

func SwapInt32AcqRel(addr *int32, new int32) int32 {
	return atomic.SwapInt32(addr, new)
}

func CasInt32Relaxed(addr *int32, old, new int32) bool {
	return atomic.CompareAndSwapInt32(addr, old, new)
}

func CasInt32Acquire(addr *int32, old, new int32) bool {
	return atomic.CompareAndSwapInt32(addr, old, new)
}

func CasInt32Release(addr *int32, old, new int32) bool {
	return atomic.CompareAndSwapInt32(addr, old, new)
}

func CasInt32AcqRel(addr *int32, old, new int32) bool {
	return atomic.CompareAndSwapInt32(addr, old, new)
}

func CaxInt32Relaxed(addr *int32, old, new int32) int32 {
	for {
		cur := atomic.LoadInt32(addr)
		if cur != old {
			return cur
		}
		if atomic.CompareAndSwapInt32(addr, old, new) {
			return old
		}
	}
}

func CaxInt32Acquire(addr *int32, old, new int32) int32 {
	return CaxInt32Relaxed(addr, old, new)
}

func CaxInt32Release(addr *int32, old, new int32) int32 {
	return CaxInt32Relaxed(addr, old, new)
}

func CaxInt32AcqRel(addr *int32, old, new int32) int32 {
	return CaxInt32Relaxed(addr, old, new)
}

func AddInt32Relaxed(addr *int32, delta int32) int32 {
	return atomic.AddInt32(addr, delta)
}

func AddInt32Acquire(addr *int32, delta int32) int32 {
	return atomic.AddInt32(addr, delta)
}

func AddInt32Release(addr *int32, delta int32) int32 {
	return atomic.AddInt32(addr, delta)
}

func AddInt32AcqRel(addr *int32, delta int32) int32 {
	return atomic.AddInt32(addr, delta)
}

func AndInt32Relaxed(addr *int32, mask int32) int32 {
	for {
		old := atomic.LoadInt32(addr)
		if atomic.CompareAndSwapInt32(addr, old, old&mask) {
			return old
		}
	}
}

func AndInt32Acquire(addr *int32, mask int32) int32 {
	return AndInt32Relaxed(addr, mask)
}

func AndInt32Release(addr *int32, mask int32) int32 {
	return AndInt32Relaxed(addr, mask)
}

func AndInt32AcqRel(addr *int32, mask int32) int32 {
	return AndInt32Relaxed(addr, mask)
}

func OrInt32Relaxed(addr *int32, mask int32) int32 {
	for {
		old := atomic.LoadInt32(addr)
		if atomic.CompareAndSwapInt32(addr, old, old|mask) {
			return old
		}
	}
}

func OrInt32Acquire(addr *int32, mask int32) int32 {
	return OrInt32Relaxed(addr, mask)
}

func OrInt32Release(addr *int32, mask int32) int32 {
	return OrInt32Relaxed(addr, mask)
}

func OrInt32AcqRel(addr *int32, mask int32) int32 {
	return OrInt32Relaxed(addr, mask)
}

func XorInt32Relaxed(addr *int32, mask int32) int32 {
	for {
		old := atomic.LoadInt32(addr)
		if atomic.CompareAndSwapInt32(addr, old, old^mask) {
			return old
		}
	}
}

func XorInt32Acquire(addr *int32, mask int32) int32 {
	return XorInt32Relaxed(addr, mask)
}

func XorInt32Release(addr *int32, mask int32) int32 {
	return XorInt32Relaxed(addr, mask)
}

func XorInt32AcqRel(addr *int32, mask int32) int32 {
	return XorInt32Relaxed(addr, mask)
}

// =============================================================================
// 32-bit Unsigned Integer Operations
// =============================================================================

func LoadUint32Relaxed(addr *uint32) uint32 {
	return atomic.LoadUint32(addr)
}

func LoadUint32Acquire(addr *uint32) uint32 {
	return atomic.LoadUint32(addr)
}

func StoreUint32Relaxed(addr *uint32, val uint32) {
	atomic.StoreUint32(addr, val)
}

func StoreUint32Release(addr *uint32, val uint32) {
	atomic.StoreUint32(addr, val)
}

func SwapUint32Relaxed(addr *uint32, new uint32) uint32 {
	return atomic.SwapUint32(addr, new)
}

func SwapUint32Acquire(addr *uint32, new uint32) uint32 {
	return atomic.SwapUint32(addr, new)
}

func SwapUint32Release(addr *uint32, new uint32) uint32 {
	return atomic.SwapUint32(addr, new)
}

func SwapUint32AcqRel(addr *uint32, new uint32) uint32 {
	return atomic.SwapUint32(addr, new)
}

func CasUint32Relaxed(addr *uint32, old, new uint32) bool {
	return atomic.CompareAndSwapUint32(addr, old, new)
}

func CasUint32Acquire(addr *uint32, old, new uint32) bool {
	return atomic.CompareAndSwapUint32(addr, old, new)
}

func CasUint32Release(addr *uint32, old, new uint32) bool {
	return atomic.CompareAndSwapUint32(addr, old, new)
}

func CasUint32AcqRel(addr *uint32, old, new uint32) bool {
	return atomic.CompareAndSwapUint32(addr, old, new)
}

func CaxUint32Relaxed(addr *uint32, old, new uint32) uint32 {
	for {
		cur := atomic.LoadUint32(addr)
		if cur != old {
			return cur
		}
		if atomic.CompareAndSwapUint32(addr, old, new) {
			return old
		}
	}
}

func CaxUint32Acquire(addr *uint32, old, new uint32) uint32 {
	return CaxUint32Relaxed(addr, old, new)
}

func CaxUint32Release(addr *uint32, old, new uint32) uint32 {
	return CaxUint32Relaxed(addr, old, new)
}

func CaxUint32AcqRel(addr *uint32, old, new uint32) uint32 {
	return CaxUint32Relaxed(addr, old, new)
}

func AddUint32Relaxed(addr *uint32, delta uint32) uint32 {
	return atomic.AddUint32(addr, delta)
}

func AddUint32Acquire(addr *uint32, delta uint32) uint32 {
	return atomic.AddUint32(addr, delta)
}

func AddUint32Release(addr *uint32, delta uint32) uint32 {
	return atomic.AddUint32(addr, delta)
}

func AddUint32AcqRel(addr *uint32, delta uint32) uint32 {
	return atomic.AddUint32(addr, delta)
}

func AndUint32Relaxed(addr *uint32, mask uint32) uint32 {
	for {
		old := atomic.LoadUint32(addr)
		if atomic.CompareAndSwapUint32(addr, old, old&mask) {
			return old
		}
	}
}

func AndUint32Acquire(addr *uint32, mask uint32) uint32 {
	return AndUint32Relaxed(addr, mask)
}

func AndUint32Release(addr *uint32, mask uint32) uint32 {
	return AndUint32Relaxed(addr, mask)
}

func AndUint32AcqRel(addr *uint32, mask uint32) uint32 {
	return AndUint32Relaxed(addr, mask)
}

func OrUint32Relaxed(addr *uint32, mask uint32) uint32 {
	for {
		old := atomic.LoadUint32(addr)
		if atomic.CompareAndSwapUint32(addr, old, old|mask) {
			return old
		}
	}
}

func OrUint32Acquire(addr *uint32, mask uint32) uint32 {
	return OrUint32Relaxed(addr, mask)
}

func OrUint32Release(addr *uint32, mask uint32) uint32 {
	return OrUint32Relaxed(addr, mask)
}

func OrUint32AcqRel(addr *uint32, mask uint32) uint32 {
	return OrUint32Relaxed(addr, mask)
}

func XorUint32Relaxed(addr *uint32, mask uint32) uint32 {
	for {
		old := atomic.LoadUint32(addr)
		if atomic.CompareAndSwapUint32(addr, old, old^mask) {
			return old
		}
	}
}

func XorUint32Acquire(addr *uint32, mask uint32) uint32 {
	return XorUint32Relaxed(addr, mask)
}

func XorUint32Release(addr *uint32, mask uint32) uint32 {
	return XorUint32Relaxed(addr, mask)
}

func XorUint32AcqRel(addr *uint32, mask uint32) uint32 {
	return XorUint32Relaxed(addr, mask)
}

// =============================================================================
// 64-bit Signed Integer Operations
// =============================================================================

func LoadInt64Relaxed(addr *int64) int64 {
	return atomic.LoadInt64(addr)
}

func LoadInt64Acquire(addr *int64) int64 {
	return atomic.LoadInt64(addr)
}

func StoreInt64Relaxed(addr *int64, val int64) {
	atomic.StoreInt64(addr, val)
}

func StoreInt64Release(addr *int64, val int64) {
	atomic.StoreInt64(addr, val)
}

func SwapInt64Relaxed(addr *int64, new int64) int64 {
	return atomic.SwapInt64(addr, new)
}

func SwapInt64Acquire(addr *int64, new int64) int64 {
	return atomic.SwapInt64(addr, new)
}

func SwapInt64Release(addr *int64, new int64) int64 {
	return atomic.SwapInt64(addr, new)
}

func SwapInt64AcqRel(addr *int64, new int64) int64 {
	return atomic.SwapInt64(addr, new)
}

func CasInt64Relaxed(addr *int64, old, new int64) bool {
	return atomic.CompareAndSwapInt64(addr, old, new)
}

func CasInt64Acquire(addr *int64, old, new int64) bool {
	return atomic.CompareAndSwapInt64(addr, old, new)
}

func CasInt64Release(addr *int64, old, new int64) bool {
	return atomic.CompareAndSwapInt64(addr, old, new)
}

func CasInt64AcqRel(addr *int64, old, new int64) bool {
	return atomic.CompareAndSwapInt64(addr, old, new)
}

func CaxInt64Relaxed(addr *int64, old, new int64) int64 {
	for {
		cur := atomic.LoadInt64(addr)
		if cur != old {
			return cur
		}
		if atomic.CompareAndSwapInt64(addr, old, new) {
			return old
		}
	}
}

func CaxInt64Acquire(addr *int64, old, new int64) int64 {
	return CaxInt64Relaxed(addr, old, new)
}

func CaxInt64Release(addr *int64, old, new int64) int64 {
	return CaxInt64Relaxed(addr, old, new)
}

func CaxInt64AcqRel(addr *int64, old, new int64) int64 {
	return CaxInt64Relaxed(addr, old, new)
}

func AddInt64Relaxed(addr *int64, delta int64) int64 {
	return atomic.AddInt64(addr, delta)
}

func AddInt64Acquire(addr *int64, delta int64) int64 {
	return atomic.AddInt64(addr, delta)
}

func AddInt64Release(addr *int64, delta int64) int64 {
	return atomic.AddInt64(addr, delta)
}

func AddInt64AcqRel(addr *int64, delta int64) int64 {
	return atomic.AddInt64(addr, delta)
}

func AndInt64Relaxed(addr *int64, mask int64) int64 {
	for {
		old := atomic.LoadInt64(addr)
		if atomic.CompareAndSwapInt64(addr, old, old&mask) {
			return old
		}
	}
}

func AndInt64Acquire(addr *int64, mask int64) int64 {
	return AndInt64Relaxed(addr, mask)
}

func AndInt64Release(addr *int64, mask int64) int64 {
	return AndInt64Relaxed(addr, mask)
}

func AndInt64AcqRel(addr *int64, mask int64) int64 {
	return AndInt64Relaxed(addr, mask)
}

func OrInt64Relaxed(addr *int64, mask int64) int64 {
	for {
		old := atomic.LoadInt64(addr)
		if atomic.CompareAndSwapInt64(addr, old, old|mask) {
			return old
		}
	}
}

func OrInt64Acquire(addr *int64, mask int64) int64 {
	return OrInt64Relaxed(addr, mask)
}

func OrInt64Release(addr *int64, mask int64) int64 {
	return OrInt64Relaxed(addr, mask)
}

func OrInt64AcqRel(addr *int64, mask int64) int64 {
	return OrInt64Relaxed(addr, mask)
}

func XorInt64Relaxed(addr *int64, mask int64) int64 {
	for {
		old := atomic.LoadInt64(addr)
		if atomic.CompareAndSwapInt64(addr, old, old^mask) {
			return old
		}
	}
}

func XorInt64Acquire(addr *int64, mask int64) int64 {
	return XorInt64Relaxed(addr, mask)
}

func XorInt64Release(addr *int64, mask int64) int64 {
	return XorInt64Relaxed(addr, mask)
}

func XorInt64AcqRel(addr *int64, mask int64) int64 {
	return XorInt64Relaxed(addr, mask)
}

// =============================================================================
// 64-bit Unsigned Integer Operations
// =============================================================================

func LoadUint64Relaxed(addr *uint64) uint64 {
	return atomic.LoadUint64(addr)
}

func LoadUint64Acquire(addr *uint64) uint64 {
	return atomic.LoadUint64(addr)
}

func StoreUint64Relaxed(addr *uint64, val uint64) {
	atomic.StoreUint64(addr, val)
}

func StoreUint64Release(addr *uint64, val uint64) {
	atomic.StoreUint64(addr, val)
}

func SwapUint64Relaxed(addr *uint64, new uint64) uint64 {
	return atomic.SwapUint64(addr, new)
}

func SwapUint64Acquire(addr *uint64, new uint64) uint64 {
	return atomic.SwapUint64(addr, new)
}

func SwapUint64Release(addr *uint64, new uint64) uint64 {
	return atomic.SwapUint64(addr, new)
}

func SwapUint64AcqRel(addr *uint64, new uint64) uint64 {
	return atomic.SwapUint64(addr, new)
}

func CasUint64Relaxed(addr *uint64, old, new uint64) bool {
	return atomic.CompareAndSwapUint64(addr, old, new)
}

func CasUint64Acquire(addr *uint64, old, new uint64) bool {
	return atomic.CompareAndSwapUint64(addr, old, new)
}

func CasUint64Release(addr *uint64, old, new uint64) bool {
	return atomic.CompareAndSwapUint64(addr, old, new)
}

func CasUint64AcqRel(addr *uint64, old, new uint64) bool {
	return atomic.CompareAndSwapUint64(addr, old, new)
}

func CaxUint64Relaxed(addr *uint64, old, new uint64) uint64 {
	for {
		cur := atomic.LoadUint64(addr)
		if cur != old {
			return cur
		}
		if atomic.CompareAndSwapUint64(addr, old, new) {
			return old
		}
	}
}

func CaxUint64Acquire(addr *uint64, old, new uint64) uint64 {
	return CaxUint64Relaxed(addr, old, new)
}

func CaxUint64Release(addr *uint64, old, new uint64) uint64 {
	return CaxUint64Relaxed(addr, old, new)
}

func CaxUint64AcqRel(addr *uint64, old, new uint64) uint64 {
	return CaxUint64Relaxed(addr, old, new)
}

func AddUint64Relaxed(addr *uint64, delta uint64) uint64 {
	return atomic.AddUint64(addr, delta)
}

func AddUint64Acquire(addr *uint64, delta uint64) uint64 {
	return atomic.AddUint64(addr, delta)
}

func AddUint64Release(addr *uint64, delta uint64) uint64 {
	return atomic.AddUint64(addr, delta)
}

func AddUint64AcqRel(addr *uint64, delta uint64) uint64 {
	return atomic.AddUint64(addr, delta)
}

func AndUint64Relaxed(addr *uint64, mask uint64) uint64 {
	for {
		old := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, old, old&mask) {
			return old
		}
	}
}

func AndUint64Acquire(addr *uint64, mask uint64) uint64 {
	return AndUint64Relaxed(addr, mask)
}

func AndUint64Release(addr *uint64, mask uint64) uint64 {
	return AndUint64Relaxed(addr, mask)
}

func AndUint64AcqRel(addr *uint64, mask uint64) uint64 {
	return AndUint64Relaxed(addr, mask)
}

func OrUint64Relaxed(addr *uint64, mask uint64) uint64 {
	for {
		old := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, old, old|mask) {
			return old
		}
	}
}

func OrUint64Acquire(addr *uint64, mask uint64) uint64 {
	return OrUint64Relaxed(addr, mask)
}

func OrUint64Release(addr *uint64, mask uint64) uint64 {
	return OrUint64Relaxed(addr, mask)
}

func OrUint64AcqRel(addr *uint64, mask uint64) uint64 {
	return OrUint64Relaxed(addr, mask)
}

func XorUint64Relaxed(addr *uint64, mask uint64) uint64 {
	for {
		old := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, old, old^mask) {
			return old
		}
	}
}

func XorUint64Acquire(addr *uint64, mask uint64) uint64 {
	return XorUint64Relaxed(addr, mask)
}

func XorUint64Release(addr *uint64, mask uint64) uint64 {
	return XorUint64Relaxed(addr, mask)
}

func XorUint64AcqRel(addr *uint64, mask uint64) uint64 {
	return XorUint64Relaxed(addr, mask)
}

// =============================================================================
// Uintptr Operations
// =============================================================================

func LoadUintptrRelaxed(addr *uintptr) uintptr {
	return atomic.LoadUintptr(addr)
}

func LoadUintptrAcquire(addr *uintptr) uintptr {
	return atomic.LoadUintptr(addr)
}

func StoreUintptrRelaxed(addr *uintptr, val uintptr) {
	atomic.StoreUintptr(addr, val)
}

func StoreUintptrRelease(addr *uintptr, val uintptr) {
	atomic.StoreUintptr(addr, val)
}

func SwapUintptrRelaxed(addr *uintptr, new uintptr) uintptr {
	return atomic.SwapUintptr(addr, new)
}

func SwapUintptrAcquire(addr *uintptr, new uintptr) uintptr {
	return atomic.SwapUintptr(addr, new)
}

func SwapUintptrRelease(addr *uintptr, new uintptr) uintptr {
	return atomic.SwapUintptr(addr, new)
}

func SwapUintptrAcqRel(addr *uintptr, new uintptr) uintptr {
	return atomic.SwapUintptr(addr, new)
}

func CasUintptrRelaxed(addr *uintptr, old, new uintptr) bool {
	return atomic.CompareAndSwapUintptr(addr, old, new)
}

func CasUintptrAcquire(addr *uintptr, old, new uintptr) bool {
	return atomic.CompareAndSwapUintptr(addr, old, new)
}

func CasUintptrRelease(addr *uintptr, old, new uintptr) bool {
	return atomic.CompareAndSwapUintptr(addr, old, new)
}

func CasUintptrAcqRel(addr *uintptr, old, new uintptr) bool {
	return atomic.CompareAndSwapUintptr(addr, old, new)
}

func CaxUintptrRelaxed(addr *uintptr, old, new uintptr) uintptr {
	for {
		cur := atomic.LoadUintptr(addr)
		if cur != old {
			return cur
		}
		if atomic.CompareAndSwapUintptr(addr, old, new) {
			return old
		}
	}
}

func CaxUintptrAcquire(addr *uintptr, old, new uintptr) uintptr {
	return CaxUintptrRelaxed(addr, old, new)
}

func CaxUintptrRelease(addr *uintptr, old, new uintptr) uintptr {
	return CaxUintptrRelaxed(addr, old, new)
}

func CaxUintptrAcqRel(addr *uintptr, old, new uintptr) uintptr {
	return CaxUintptrRelaxed(addr, old, new)
}

func AddUintptrRelaxed(addr *uintptr, delta uintptr) uintptr {
	return atomic.AddUintptr(addr, delta)
}

func AddUintptrAcquire(addr *uintptr, delta uintptr) uintptr {
	return atomic.AddUintptr(addr, delta)
}

func AddUintptrRelease(addr *uintptr, delta uintptr) uintptr {
	return atomic.AddUintptr(addr, delta)
}

func AddUintptrAcqRel(addr *uintptr, delta uintptr) uintptr {
	return atomic.AddUintptr(addr, delta)
}

func AndUintptrRelaxed(addr *uintptr, mask uintptr) uintptr {
	for {
		old := atomic.LoadUintptr(addr)
		if atomic.CompareAndSwapUintptr(addr, old, old&mask) {
			return old
		}
	}
}

func AndUintptrAcquire(addr *uintptr, mask uintptr) uintptr {
	return AndUintptrRelaxed(addr, mask)
}

func AndUintptrRelease(addr *uintptr, mask uintptr) uintptr {
	return AndUintptrRelaxed(addr, mask)
}

func AndUintptrAcqRel(addr *uintptr, mask uintptr) uintptr {
	return AndUintptrRelaxed(addr, mask)
}

func OrUintptrRelaxed(addr *uintptr, mask uintptr) uintptr {
	for {
		old := atomic.LoadUintptr(addr)
		if atomic.CompareAndSwapUintptr(addr, old, old|mask) {
			return old
		}
	}
}

func OrUintptrAcquire(addr *uintptr, mask uintptr) uintptr {
	return OrUintptrRelaxed(addr, mask)
}

func OrUintptrRelease(addr *uintptr, mask uintptr) uintptr {
	return OrUintptrRelaxed(addr, mask)
}

func OrUintptrAcqRel(addr *uintptr, mask uintptr) uintptr {
	return OrUintptrRelaxed(addr, mask)
}

func XorUintptrRelaxed(addr *uintptr, mask uintptr) uintptr {
	for {
		old := atomic.LoadUintptr(addr)
		if atomic.CompareAndSwapUintptr(addr, old, old^mask) {
			return old
		}
	}
}

func XorUintptrAcquire(addr *uintptr, mask uintptr) uintptr {
	return XorUintptrRelaxed(addr, mask)
}

func XorUintptrRelease(addr *uintptr, mask uintptr) uintptr {
	return XorUintptrRelaxed(addr, mask)
}

func XorUintptrAcqRel(addr *uintptr, mask uintptr) uintptr {
	return XorUintptrRelaxed(addr, mask)
}

// =============================================================================
// Pointer Operations
// =============================================================================

func LoadPointerRelaxed(addr *unsafe.Pointer) unsafe.Pointer {
	return atomic.LoadPointer(addr)
}

func LoadPointerAcquire(addr *unsafe.Pointer) unsafe.Pointer {
	return atomic.LoadPointer(addr)
}

func StorePointerRelaxed(addr *unsafe.Pointer, val unsafe.Pointer) {
	atomic.StorePointer(addr, val)
}

func StorePointerRelease(addr *unsafe.Pointer, val unsafe.Pointer) {
	atomic.StorePointer(addr, val)
}

func SwapPointerRelaxed(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer {
	return atomic.SwapPointer(addr, new)
}

func SwapPointerAcquire(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer {
	return atomic.SwapPointer(addr, new)
}

func SwapPointerRelease(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer {
	return atomic.SwapPointer(addr, new)
}

func SwapPointerAcqRel(addr *unsafe.Pointer, new unsafe.Pointer) unsafe.Pointer {
	return atomic.SwapPointer(addr, new)
}

func CasPointerRelaxed(addr *unsafe.Pointer, old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(addr, old, new)
}

func CasPointerAcquire(addr *unsafe.Pointer, old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(addr, old, new)
}

func CasPointerRelease(addr *unsafe.Pointer, old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(addr, old, new)
}

func CasPointerAcqRel(addr *unsafe.Pointer, old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(addr, old, new)
}

func CaxPointerRelaxed(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer {
	for {
		cur := atomic.LoadPointer(addr)
		if cur != old {
			return cur
		}
		if atomic.CompareAndSwapPointer(addr, old, new) {
			return old
		}
	}
}

func CaxPointerAcquire(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer {
	return CaxPointerRelaxed(addr, old, new)
}

func CaxPointerRelease(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer {
	return CaxPointerRelaxed(addr, old, new)
}

func CaxPointerAcqRel(addr *unsafe.Pointer, old, new unsafe.Pointer) unsafe.Pointer {
	return CaxPointerRelaxed(addr, old, new)
}

// =============================================================================
// 128-bit Operations (NOT ATOMIC!)
// =============================================================================

// WARNING: 128-bit operations are NOT concurrency-safe on generic architectures.
// These are simple sequential read/write operations without atomicity guarantees.
// Use external synchronization when accessing 128-bit values concurrently.

func LoadUint128Relaxed(addr *[16]byte) (lo, hi uint64) {
	lo = *(*uint64)(unsafe.Pointer(addr))
	hi = *(*uint64)(unsafe.Add(unsafe.Pointer(addr), 8))
	return
}

func LoadUint128Acquire(addr *[16]byte) (lo, hi uint64) {
	return LoadUint128Relaxed(addr)
}

func StoreUint128Relaxed(addr *[16]byte, lo, hi uint64) {
	*(*uint64)(unsafe.Pointer(addr)) = lo
	*(*uint64)(unsafe.Add(unsafe.Pointer(addr), 8)) = hi
}

func StoreUint128Release(addr *[16]byte, lo, hi uint64) {
	StoreUint128Relaxed(addr, lo, hi)
}

func SwapUint128Relaxed(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64) {
	oldLo = *(*uint64)(unsafe.Pointer(addr))
	oldHi = *(*uint64)(unsafe.Add(unsafe.Pointer(addr), 8))
	*(*uint64)(unsafe.Pointer(addr)) = newLo
	*(*uint64)(unsafe.Add(unsafe.Pointer(addr), 8)) = newHi
	return
}

func SwapUint128Acquire(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64) {
	return SwapUint128Relaxed(addr, newLo, newHi)
}

func SwapUint128Release(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64) {
	return SwapUint128Relaxed(addr, newLo, newHi)
}

func SwapUint128AcqRel(addr *[16]byte, newLo, newHi uint64) (oldLo, oldHi uint64) {
	return SwapUint128Relaxed(addr, newLo, newHi)
}

func CasUint128Relaxed(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool {
	curLo := *(*uint64)(unsafe.Pointer(addr))
	curHi := *(*uint64)(unsafe.Add(unsafe.Pointer(addr), 8))
	if curLo == oldLo && curHi == oldHi {
		*(*uint64)(unsafe.Pointer(addr)) = newLo
		*(*uint64)(unsafe.Add(unsafe.Pointer(addr), 8)) = newHi
		return true
	}
	return false
}

func CasUint128Acquire(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool {
	return CasUint128Relaxed(addr, oldLo, oldHi, newLo, newHi)
}

func CasUint128Release(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool {
	return CasUint128Relaxed(addr, oldLo, oldHi, newLo, newHi)
}

func CasUint128AcqRel(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) bool {
	return CasUint128Relaxed(addr, oldLo, oldHi, newLo, newHi)
}

func CaxUint128Relaxed(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	lo = *(*uint64)(unsafe.Pointer(addr))
	hi = *(*uint64)(unsafe.Add(unsafe.Pointer(addr), 8))
	if lo == oldLo && hi == oldHi {
		*(*uint64)(unsafe.Pointer(addr)) = newLo
		*(*uint64)(unsafe.Add(unsafe.Pointer(addr), 8)) = newHi
	}
	return
}

func CaxUint128Acquire(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	return CaxUint128Relaxed(addr, oldLo, oldHi, newLo, newHi)
}

func CaxUint128Release(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	return CaxUint128Relaxed(addr, oldLo, oldHi, newLo, newHi)
}

func CaxUint128AcqRel(addr *[16]byte, oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	return CaxUint128Relaxed(addr, oldLo, oldHi, newLo, newHi)
}

// =============================================================================
// Memory Barriers
// =============================================================================

// Barriers are no-ops on generic platforms.
// sync/atomic provides sequential consistency, which is stronger than
// any acquire/release barrier requirement.

func BarrierAcquire() {}

func BarrierRelease() {}

func BarrierAcqRel() {}
