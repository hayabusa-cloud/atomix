// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix_test

import (
	"sync/atomic"
	"testing"

	"code.hybscloud.com/atomix"
)

// =============================================================================
// sync/atomic vs atomix comparison benchmarks
//
// These benchmarks compare performance of sync/atomic (sequential consistency)
// versus atomix (explicit memory orderings: relaxed, acquire, release).
//
// On x86-64 TSO:
//   - Load: Both use MOV (equivalent performance)
//   - Store: sync/atomic uses XCHG (SC), atomix uses MOV (release)
//
// On ARM64:
//   - Load: sync/atomic uses LDAR (acquire), atomix relaxed uses LDR
//   - Store: sync/atomic uses STLR (release), atomix relaxed uses STR
//   - Relaxed operations should be significantly faster on ARM64
// =============================================================================

// -----------------------------------------------------------------------------
// 64-bit Load benchmarks
// -----------------------------------------------------------------------------

func BenchmarkLoad64_SyncAtomic(b *testing.B) {
	var v atomic.Int64
	v.Store(42)
	for range b.N {
		_ = v.Load()
	}
}

func BenchmarkLoad64_AtomixAcquire(b *testing.B) {
	var v atomix.Int64
	v.Store(42)
	for range b.N {
		_ = v.LoadAcquire()
	}
}

func BenchmarkLoad64_AtomixRelaxed(b *testing.B) {
	var v atomix.Int64
	v.Store(42)
	for range b.N {
		_ = v.LoadRelaxed()
	}
}

// -----------------------------------------------------------------------------
// 64-bit Store benchmarks
// -----------------------------------------------------------------------------

func BenchmarkStore64_SyncAtomic(b *testing.B) {
	var v atomic.Int64
	for i := range b.N {
		v.Store(int64(i))
	}
}

func BenchmarkStore64_AtomixRelease(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.StoreRelease(int64(i))
	}
}

func BenchmarkStore64_AtomixRelaxed(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.StoreRelaxed(int64(i))
	}
}

// -----------------------------------------------------------------------------
// 32-bit Load benchmarks
// -----------------------------------------------------------------------------

func BenchmarkLoad32_SyncAtomic(b *testing.B) {
	var v atomic.Int32
	v.Store(42)
	for range b.N {
		_ = v.Load()
	}
}

func BenchmarkLoad32_AtomixAcquire(b *testing.B) {
	var v atomix.Int32
	v.Store(42)
	for range b.N {
		_ = v.LoadAcquire()
	}
}

func BenchmarkLoad32_AtomixRelaxed(b *testing.B) {
	var v atomix.Int32
	v.Store(42)
	for range b.N {
		_ = v.LoadRelaxed()
	}
}

// -----------------------------------------------------------------------------
// 32-bit Store benchmarks
// -----------------------------------------------------------------------------

func BenchmarkStore32_SyncAtomic(b *testing.B) {
	var v atomic.Int32
	for i := range b.N {
		v.Store(int32(i))
	}
}

func BenchmarkStore32_AtomixRelease(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.StoreRelease(int32(i))
	}
}

func BenchmarkStore32_AtomixRelaxed(b *testing.B) {
	var v atomix.Int32
	for i := range b.N {
		v.StoreRelaxed(int32(i))
	}
}

// -----------------------------------------------------------------------------
// Add benchmarks (RMW operations)
// -----------------------------------------------------------------------------

func BenchmarkAdd64_SyncAtomic(b *testing.B) {
	var v atomic.Int64
	for range b.N {
		v.Add(1)
	}
}

func BenchmarkAdd64_AtomixAcqRel(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.Add(1) // Default is AcqRel
	}
}

func BenchmarkAdd64_AtomixRelaxed(b *testing.B) {
	var v atomix.Int64
	for range b.N {
		v.AddRelaxed(1)
	}
}

// -----------------------------------------------------------------------------
// CAS benchmarks
// -----------------------------------------------------------------------------

func BenchmarkCAS64_SyncAtomic(b *testing.B) {
	var v atomic.Int64
	for i := range b.N {
		v.CompareAndSwap(int64(i), int64(i+1))
	}
}

func BenchmarkCAS64_AtomixAcqRel(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareAndSwap(int64(i), int64(i+1))
	}
}

func BenchmarkCAS64_AtomixRelaxed(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.CompareAndSwapRelaxed(int64(i), int64(i+1))
	}
}

// -----------------------------------------------------------------------------
// Swap benchmarks
// -----------------------------------------------------------------------------

func BenchmarkSwap64_SyncAtomic(b *testing.B) {
	var v atomic.Int64
	for i := range b.N {
		v.Swap(int64(i))
	}
}

func BenchmarkSwap64_AtomixAcqRel(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.Swap(int64(i))
	}
}

func BenchmarkSwap64_AtomixRelaxed(b *testing.B) {
	var v atomix.Int64
	for i := range b.N {
		v.SwapRelaxed(int64(i))
	}
}

// -----------------------------------------------------------------------------
// 128-bit benchmarks (atomix only - sync/atomic doesn't support 128-bit)
// -----------------------------------------------------------------------------

func BenchmarkLoad128_AtomixAcquire(b *testing.B) {
	// Use properly aligned buffer for 128-bit
	var buf [32]byte
	_, v := atomix.PlaceAlignedUint128(buf[:], 0)
	v.Store(42, 0)
	for range b.N {
		_, _ = v.LoadAcquire()
	}
}

func BenchmarkLoad128_AtomixRelaxed(b *testing.B) {
	var buf [32]byte
	_, v := atomix.PlaceAlignedUint128(buf[:], 0)
	v.Store(42, 0)
	for range b.N {
		_, _ = v.LoadRelaxed()
	}
}

func BenchmarkStore128_AtomixRelease(b *testing.B) {
	var buf [32]byte
	_, v := atomix.PlaceAlignedUint128(buf[:], 0)
	for i := range b.N {
		v.StoreRelease(uint64(i), 0)
	}
}

func BenchmarkStore128_AtomixRelaxed(b *testing.B) {
	var buf [32]byte
	_, v := atomix.PlaceAlignedUint128(buf[:], 0)
	for i := range b.N {
		v.StoreRelaxed(uint64(i), 0)
	}
}

func BenchmarkCAS128_AtomixAcqRel(b *testing.B) {
	var buf [32]byte
	_, v := atomix.PlaceAlignedUint128(buf[:], 0)
	for i := range b.N {
		v.CompareAndSwap(uint64(i), 0, uint64(i+1), 0)
	}
}

func BenchmarkCAS128_AtomixRelaxed(b *testing.B) {
	var buf [32]byte
	_, v := atomix.PlaceAlignedUint128(buf[:], 0)
	for i := range b.N {
		v.CompareAndSwapRelaxed(uint64(i), 0, uint64(i+1), 0)
	}
}
