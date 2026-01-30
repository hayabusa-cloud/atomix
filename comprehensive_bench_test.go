// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build linux

package atomix_test

import (
	"code.hybscloud.com/atomix"
	"sync/atomic"
	"testing"
)

// Package-level sink variables prevent dead-store elimination.
// Plain MOV stores to stack locals can be eliminated by the compiler,
// but stores to package-level variables cannot (externally observable).
var (
	benchSinkInt32   atomix.Int32
	benchSinkInt64   atomix.Int64
	benchSinkUintptr atomix.Uintptr
	benchSyncInt32   int32
	benchSyncInt64   int64
	benchSyncUintptr uintptr
)

// =============================================================================
// 64-bit Load Operations
// =============================================================================

func BenchmarkCompLoad64_SyncAtomic(b *testing.B) {
	var v int64 = 42
	for range b.N {
		_ = atomic.LoadInt64(&v)
	}
}

func BenchmarkCompLoad64_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int64
	v.Store(42)
	for range b.N {
		_ = v.LoadRelaxed()
	}
}

func BenchmarkCompLoad64_Atomix_Acquire(b *testing.B) {
	var v atomix.Int64
	v.Store(42)
	for range b.N {
		_ = v.LoadAcquire()
	}
}

// =============================================================================
// 64-bit Store Operations
// =============================================================================

func BenchmarkCompStore64_SyncAtomic(b *testing.B) {
	for i := range b.N {
		atomic.StoreInt64(&benchSyncInt64, int64(i))
	}
}

func BenchmarkCompStore64_Atomix_Relaxed(b *testing.B) {
	for i := range b.N {
		benchSinkInt64.StoreRelaxed(int64(i))
	}
}

func BenchmarkCompStore64_Atomix_Release(b *testing.B) {
	for i := range b.N {
		benchSinkInt64.StoreRelease(int64(i))
	}
}

// =============================================================================
// 64-bit Add Operations
// =============================================================================

func BenchmarkCompAdd64_SyncAtomic(b *testing.B) {
	var v int64
	for range b.N {
		atomic.AddInt64(&v, 1)
	}
}

func BenchmarkCompAdd64_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.AddRelaxed(1)
	}
}

func BenchmarkCompAdd64_Atomix_Acquire(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.AddAcquire(1)
	}
}

func BenchmarkCompAdd64_Atomix_Release(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.AddRelease(1)
	}
}

func BenchmarkCompAdd64_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.AddAcqRel(1)
	}
}

// =============================================================================
// 64-bit Swap Operations
// =============================================================================

func BenchmarkCompSwap64_SyncAtomic(b *testing.B) {
	var v int64
	for i := range b.N {
		atomic.SwapInt64(&v, int64(i))
	}
}

func BenchmarkCompSwap64_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.SwapRelaxed(int64(i))
	}
}

func BenchmarkCompSwap64_Atomix_Acquire(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.SwapAcquire(int64(i))
	}
}

func BenchmarkCompSwap64_Atomix_Release(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.SwapRelease(int64(i))
	}
}

func BenchmarkCompSwap64_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.SwapAcqRel(int64(i))
	}
}

// =============================================================================
// 64-bit CAS Operations
// =============================================================================

func BenchmarkCompCAS64_SyncAtomic(b *testing.B) {
	var v int64
	for i := range b.N {
		atomic.CompareAndSwapInt64(&v, int64(i), int64(i+1))
	}
}

func BenchmarkCompCAS64_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareAndSwapRelaxed(int64(i), int64(i+1))
	}
}

func BenchmarkCompCAS64_Atomix_Acquire(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareAndSwapAcquire(int64(i), int64(i+1))
	}
}

func BenchmarkCompCAS64_Atomix_Release(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareAndSwapRelease(int64(i), int64(i+1))
	}
}

func BenchmarkCompCAS64_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareAndSwapAcqRel(int64(i), int64(i+1))
	}
}

// =============================================================================
// 64-bit CompareExchange Operations (atomix-only — returns old value)
// =============================================================================

func BenchmarkCompCax64_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareExchangeRelaxed(int64(i), int64(i+1))
	}
}

func BenchmarkCompCax64_Atomix_Acquire(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareExchangeAcquire(int64(i), int64(i+1))
	}
}

func BenchmarkCompCax64_Atomix_Release(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareExchangeRelease(int64(i), int64(i+1))
	}
}

func BenchmarkCompCax64_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareExchangeAcqRel(int64(i), int64(i+1))
	}
}

// =============================================================================
// 64-bit Or Operations
// =============================================================================

func BenchmarkCompOr64_SyncAtomic(b *testing.B) {
	var v int64
	for range b.N {
		atomic.OrInt64(&v, 0x1)
	}
}

func BenchmarkCompOr64_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.OrRelaxed(0x1)
	}
}

func BenchmarkCompOr64_Atomix_Acquire(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.OrAcquire(0x1)
	}
}

func BenchmarkCompOr64_Atomix_Release(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.OrRelease(0x1)
	}
}

func BenchmarkCompOr64_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.OrAcqRel(0x1)
	}
}

// =============================================================================
// 64-bit And Operations
// =============================================================================

func BenchmarkCompAnd64_SyncAtomic(b *testing.B) {
	var v int64 = -1
	for range b.N {
		atomic.AndInt64(&v, 0x7FFFFFFFFFFFFFFF)
	}
}

func BenchmarkCompAnd64_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int64
	v.Store(-1)
	for range b.N {
		v.AndRelaxed(0x7FFFFFFFFFFFFFFF)
	}
}

func BenchmarkCompAnd64_Atomix_Acquire(b *testing.B) {
	var v atomix.Int64
	v.Store(-1)
	for range b.N {
		v.AndAcquire(0x7FFFFFFFFFFFFFFF)
	}
}

func BenchmarkCompAnd64_Atomix_Release(b *testing.B) {
	var v atomix.Int64
	v.Store(-1)
	for range b.N {
		v.AndRelease(0x7FFFFFFFFFFFFFFF)
	}
}

func BenchmarkCompAnd64_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int64
	v.Store(-1)
	for range b.N {
		v.AndAcqRel(0x7FFFFFFFFFFFFFFF)
	}
}

// =============================================================================
// 64-bit Xor Operations (atomix-only — sync/atomic does not provide Xor)
// =============================================================================

func BenchmarkCompXor64_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.XorRelaxed(0x1)
	}
}

func BenchmarkCompXor64_Atomix_Acquire(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.XorAcquire(0x1)
	}
}

func BenchmarkCompXor64_Atomix_Release(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.XorRelease(0x1)
	}
}

func BenchmarkCompXor64_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.XorAcqRel(0x1)
	}
}

// =============================================================================
// 32-bit Load Operations
// =============================================================================

func BenchmarkCompLoad32_SyncAtomic(b *testing.B) {
	var v int32 = 42
	for range b.N {
		_ = atomic.LoadInt32(&v)
	}
}

func BenchmarkCompLoad32_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int32
	v.Store(42)
	for range b.N {
		_ = v.LoadRelaxed()
	}
}

func BenchmarkCompLoad32_Atomix_Acquire(b *testing.B) {
	var v atomix.Int32
	v.Store(42)
	for range b.N {
		_ = v.LoadAcquire()
	}
}

// =============================================================================
// 32-bit Store Operations
// =============================================================================

func BenchmarkCompStore32_SyncAtomic(b *testing.B) {
	for i := range b.N {
		atomic.StoreInt32(&benchSyncInt32, int32(i))
	}
}

func BenchmarkCompStore32_Atomix_Relaxed(b *testing.B) {
	for i := range b.N {
		benchSinkInt32.StoreRelaxed(int32(i))
	}
}

func BenchmarkCompStore32_Atomix_Release(b *testing.B) {
	for i := range b.N {
		benchSinkInt32.StoreRelease(int32(i))
	}
}

// =============================================================================
// 32-bit Add Operations
// =============================================================================

func BenchmarkCompAdd32_SyncAtomic(b *testing.B) {
	var v int32
	for range b.N {
		atomic.AddInt32(&v, 1)
	}
}

func BenchmarkCompAdd32_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.AddRelaxed(1)
	}
}

func BenchmarkCompAdd32_Atomix_Acquire(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.AddAcquire(1)
	}
}

func BenchmarkCompAdd32_Atomix_Release(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.AddRelease(1)
	}
}

func BenchmarkCompAdd32_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.AddAcqRel(1)
	}
}

// =============================================================================
// 32-bit Swap Operations
// =============================================================================

func BenchmarkCompSwap32_SyncAtomic(b *testing.B) {
	var v int32
	for i := range b.N {
		atomic.SwapInt32(&v, int32(i))
	}
}

func BenchmarkCompSwap32_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.SwapRelaxed(int32(i))
	}
}

func BenchmarkCompSwap32_Atomix_Acquire(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.SwapAcquire(int32(i))
	}
}

func BenchmarkCompSwap32_Atomix_Release(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.SwapRelease(int32(i))
	}
}

func BenchmarkCompSwap32_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.SwapAcqRel(int32(i))
	}
}

// =============================================================================
// 32-bit CAS Operations
// =============================================================================

func BenchmarkCompCAS32_SyncAtomic(b *testing.B) {
	var v int32
	for i := range b.N {
		atomic.CompareAndSwapInt32(&v, int32(i), int32(i+1))
	}
}

func BenchmarkCompCAS32_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.CompareAndSwapRelaxed(int32(i), int32(i+1))
	}
}

func BenchmarkCompCAS32_Atomix_Acquire(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.CompareAndSwapAcquire(int32(i), int32(i+1))
	}
}

func BenchmarkCompCAS32_Atomix_Release(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.CompareAndSwapRelease(int32(i), int32(i+1))
	}
}

func BenchmarkCompCAS32_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.CompareAndSwapAcqRel(int32(i), int32(i+1))
	}
}

// =============================================================================
// 32-bit CompareExchange Operations (atomix-only — returns old value)
// =============================================================================

func BenchmarkCompCax32_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.CompareExchangeRelaxed(int32(i), int32(i+1))
	}
}

func BenchmarkCompCax32_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.CompareExchangeAcqRel(int32(i), int32(i+1))
	}
}

// =============================================================================
// 32-bit Or Operations
// =============================================================================

func BenchmarkCompOr32_SyncAtomic(b *testing.B) {
	var v int32
	for range b.N {
		atomic.OrInt32(&v, 0x1)
	}
}

func BenchmarkCompOr32_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.OrRelaxed(0x1)
	}
}

func BenchmarkCompOr32_Atomix_Acquire(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.OrAcquire(0x1)
	}
}

func BenchmarkCompOr32_Atomix_Release(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.OrRelease(0x1)
	}
}

func BenchmarkCompOr32_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.OrAcqRel(0x1)
	}
}

// =============================================================================
// 32-bit And Operations
// =============================================================================

func BenchmarkCompAnd32_SyncAtomic(b *testing.B) {
	var v int32 = -1
	for range b.N {
		atomic.AndInt32(&v, 0x7FFFFFFF)
	}
}

func BenchmarkCompAnd32_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int32
	v.Store(-1)
	for range b.N {
		v.AndRelaxed(0x7FFFFFFF)
	}
}

func BenchmarkCompAnd32_Atomix_Acquire(b *testing.B) {
	var v atomix.Int32
	v.Store(-1)
	for range b.N {
		v.AndAcquire(0x7FFFFFFF)
	}
}

func BenchmarkCompAnd32_Atomix_Release(b *testing.B) {
	var v atomix.Int32
	v.Store(-1)
	for range b.N {
		v.AndRelease(0x7FFFFFFF)
	}
}

func BenchmarkCompAnd32_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int32
	v.Store(-1)
	for range b.N {
		v.AndAcqRel(0x7FFFFFFF)
	}
}

// =============================================================================
// 32-bit Xor Operations (atomix-only — sync/atomic does not provide Xor)
// =============================================================================

func BenchmarkCompXor32_Atomix_Relaxed(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.XorRelaxed(0x1)
	}
}

func BenchmarkCompXor32_Atomix_Acquire(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.XorAcquire(0x1)
	}
}

func BenchmarkCompXor32_Atomix_Release(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.XorRelease(0x1)
	}
}

func BenchmarkCompXor32_Atomix_AcqRel(b *testing.B) {
	var v atomix.Int32
	for range b.N {
		v.XorAcqRel(0x1)
	}
}

// =============================================================================
// 128-bit Operations (atomix only - sync/atomic has no 128-bit)
// =============================================================================

func BenchmarkComp128Load_Atomix_Relaxed(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	v.Store(0x1234, 0xABCD)
	for range b.N {
		_, _ = v.LoadRelaxed()
	}
}

func BenchmarkComp128Load_Atomix_Acquire(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	v.Store(0x1234, 0xABCD)
	for range b.N {
		_, _ = v.LoadAcquire()
	}
}

func BenchmarkComp128Store_Atomix_Relaxed(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	for i := range b.N {
		v.StoreRelaxed(uint64(i), uint64(i>>32))
	}
}

func BenchmarkComp128Store_Atomix_Release(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	for i := range b.N {
		v.StoreRelease(uint64(i), uint64(i>>32))
	}
}

func BenchmarkComp128Swap_Atomix_Relaxed(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	for i := range b.N {
		_, _ = v.SwapRelaxed(uint64(i), uint64(i>>32))
	}
}

func BenchmarkComp128Swap_Atomix_AcqRel(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	for i := range b.N {
		_, _ = v.SwapAcqRel(uint64(i), uint64(i>>32))
	}
}

func BenchmarkComp128CAS_Atomix_Relaxed(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	for i := range b.N {
		v.CompareAndSwapRelaxed(uint64(i), 0, uint64(i+1), 0)
	}
}

func BenchmarkComp128CAS_Atomix_Acquire(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	for i := range b.N {
		v.CompareAndSwapAcquire(uint64(i), 0, uint64(i+1), 0)
	}
}

func BenchmarkComp128CAS_Atomix_Release(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	for i := range b.N {
		v.CompareAndSwapRelease(uint64(i), 0, uint64(i+1), 0)
	}
}

func BenchmarkComp128CAS_Atomix_AcqRel(b *testing.B) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	for i := range b.N {
		v.CompareAndSwapAcqRel(uint64(i), 0, uint64(i+1), 0)
	}
}

// =============================================================================
// Uintptr Operations
// =============================================================================

func BenchmarkCompLoadPtr_SyncAtomic(b *testing.B) {
	var v uintptr = 0x1234
	for range b.N {
		_ = atomic.LoadUintptr(&v)
	}
}

func BenchmarkCompLoadPtr_Atomix_Relaxed(b *testing.B) {
	var v atomix.Uintptr
	v.Store(0x1234)
	for range b.N {
		_ = v.LoadRelaxed()
	}
}

func BenchmarkCompLoadPtr_Atomix_Acquire(b *testing.B) {
	var v atomix.Uintptr
	v.Store(0x1234)
	for range b.N {
		_ = v.LoadAcquire()
	}
}

func BenchmarkCompStorePtr_SyncAtomic(b *testing.B) {
	for i := range b.N {
		atomic.StoreUintptr(&benchSyncUintptr, uintptr(i))
	}
}

func BenchmarkCompStorePtr_Atomix_Relaxed(b *testing.B) {
	for i := range b.N {
		benchSinkUintptr.StoreRelaxed(uintptr(i))
	}
}

func BenchmarkCompStorePtr_Atomix_Release(b *testing.B) {
	for i := range b.N {
		benchSinkUintptr.StoreRelease(uintptr(i))
	}
}

func BenchmarkCompAddPtr_SyncAtomic(b *testing.B) {
	var v uintptr
	for range b.N {
		atomic.AddUintptr(&v, 1)
	}
}

func BenchmarkCompAddPtr_Atomix_Relaxed(b *testing.B) {
	var v atomix.Uintptr
	for range b.N {
		v.AddRelaxed(1)
	}
}

func BenchmarkCompAddPtr_Atomix_AcqRel(b *testing.B) {
	var v atomix.Uintptr
	for range b.N {
		v.AddAcqRel(1)
	}
}

// =============================================================================
// Memory Barriers
// =============================================================================

func BenchmarkCompBarrier_Atomix_Acquire(b *testing.B) {
	for range b.N {
		atomix.BarrierAcquire()
	}
}

func BenchmarkCompBarrier_Atomix_Release(b *testing.B) {
	for range b.N {
		atomix.BarrierRelease()
	}
}

func BenchmarkCompBarrier_Atomix_AcqRel(b *testing.B) {
	for range b.N {
		atomix.BarrierAcqRel()
	}
}
