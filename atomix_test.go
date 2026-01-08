// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix_test

import (
	"sync"
	"testing"

	"code.hybscloud.com/atomix"
)

// =============================================================================
// Int32 Tests
// =============================================================================

func TestInt32LoadStore(t *testing.T) {
	var a atomix.Int32

	// Zero value
	if got := a.Load(); got != 0 {
		t.Fatalf("zero value: got %d, want 0", got)
	}

	// Store and Load
	a.Store(42)
	if got := a.Load(); got != 42 {
		t.Fatalf("Load: got %d, want 42", got)
	}

	// Relaxed ordering
	a.StoreRelaxed(-100)
	if got := a.LoadRelaxed(); got != -100 {
		t.Fatalf("LoadRelaxed: got %d, want -100", got)
	}

	// Acquire/Release
	a.StoreRelease(999)
	if got := a.LoadAcquire(); got != 999 {
		t.Fatalf("LoadAcquire: got %d, want 999", got)
	}
}

func TestInt32Swap(t *testing.T) {
	var a atomix.Int32
	a.Store(10)

	old := a.Swap(20)
	if old != 10 {
		t.Fatalf("Swap old: got %d, want 10", old)
	}
	if got := a.Load(); got != 20 {
		t.Fatalf("Swap new: got %d, want 20", got)
	}

	// Test all ordering variants
	a.Store(100)
	if old := a.SwapRelaxed(101); old != 100 {
		t.Fatalf("SwapRelaxed: got %d, want 100", old)
	}
	if old := a.SwapAcquire(102); old != 101 {
		t.Fatalf("SwapAcquire: got %d, want 101", old)
	}
	if old := a.SwapRelease(103); old != 102 {
		t.Fatalf("SwapRelease: got %d, want 102", old)
	}
	if old := a.SwapAcqRel(104); old != 103 {
		t.Fatalf("SwapAcqRel: got %d, want 103", old)
	}
}

func TestInt32CompareAndSwap(t *testing.T) {
	var a atomix.Int32
	a.Store(10)

	// Successful CAS
	if !a.CompareAndSwap(10, 20) {
		t.Fatal("CAS should succeed")
	}
	if got := a.Load(); got != 20 {
		t.Fatalf("CAS result: got %d, want 20", got)
	}

	// Failed CAS
	if a.CompareAndSwap(10, 30) {
		t.Fatal("CAS should fail")
	}
	if got := a.Load(); got != 20 {
		t.Fatalf("CAS unchanged: got %d, want 20", got)
	}

	// Test all ordering variants
	a.Store(100)
	if !a.CompareAndSwapRelaxed(100, 101) {
		t.Fatal("CASRelaxed should succeed")
	}
	if !a.CompareAndSwapAcquire(101, 102) {
		t.Fatal("CASAcquire should succeed")
	}
	if !a.CompareAndSwapRelease(102, 103) {
		t.Fatal("CASRelease should succeed")
	}
	if !a.CompareAndSwapAcqRel(103, 104) {
		t.Fatal("CASAcqRel should succeed")
	}
}

func TestInt32CompareExchange(t *testing.T) {
	var a atomix.Int32
	a.Store(10)

	// Successful exchange
	old := a.CompareExchange(10, 20)
	if old != 10 {
		t.Fatalf("CompareExchange old: got %d, want 10", old)
	}
	if got := a.Load(); got != 20 {
		t.Fatalf("CompareExchange new: got %d, want 20", got)
	}

	// Failed exchange returns current value
	old = a.CompareExchange(10, 30)
	if old != 20 {
		t.Fatalf("CompareExchange failed: got %d, want 20", old)
	}

	// Test ordering variants
	a.Store(100)
	if old := a.CompareExchangeRelaxed(100, 101); old != 100 {
		t.Fatalf("Cax Relaxed: got %d, want 100", old)
	}
	if old := a.CompareExchangeAcquire(101, 102); old != 101 {
		t.Fatalf("Cax Acquire: got %d, want 101", old)
	}
	if old := a.CompareExchangeRelease(102, 103); old != 102 {
		t.Fatalf("Cax Release: got %d, want 102", old)
	}
	if old := a.CompareExchangeAcqRel(103, 104); old != 103 {
		t.Fatalf("Cax AcqRel: got %d, want 103", old)
	}
}

func TestInt32Add(t *testing.T) {
	var a atomix.Int32
	a.Store(10)

	// Add returns NEW value (like sync/atomic)
	got := a.Add(5)
	if got != 15 {
		t.Fatalf("Add: got %d, want 15", got)
	}
	if v := a.Load(); v != 15 {
		t.Fatalf("Load after Add: got %d, want 15", v)
	}

	// Negative delta
	got = a.Add(-20)
	if got != -5 {
		t.Fatalf("Add negative: got %d, want -5", got)
	}

	// Test ordering variants
	a.Store(100)
	if got := a.AddRelaxed(1); got != 101 {
		t.Fatalf("AddRelaxed: got %d, want 101", got)
	}
	if got := a.AddAcquire(1); got != 102 {
		t.Fatalf("AddAcquire: got %d, want 102", got)
	}
	if got := a.AddRelease(1); got != 103 {
		t.Fatalf("AddRelease: got %d, want 103", got)
	}
	if got := a.AddAcqRel(1); got != 104 {
		t.Fatalf("AddAcqRel: got %d, want 104", got)
	}
}

func TestInt32Sub(t *testing.T) {
	var a atomix.Int32
	a.Store(10)

	// Sub returns NEW value (like sync/atomic)
	got := a.Sub(3)
	if got != 7 {
		t.Fatalf("Sub: got %d, want 7", got)
	}
}

func TestInt32BitwiseOps(t *testing.T) {
	var a atomix.Int32
	a.Store(0xFF)

	// And
	old := a.And(0x0F)
	if old != 0xFF {
		t.Fatalf("And old: got %d, want 255", old)
	}
	if got := a.Load(); got != 0x0F {
		t.Fatalf("And new: got %d, want 15", got)
	}

	// Or
	a.Store(0x0F)
	old = a.Or(0xF0)
	if old != 0x0F {
		t.Fatalf("Or old: got %d, want 15", old)
	}
	if got := a.Load(); got != 0xFF {
		t.Fatalf("Or new: got %d, want 255", got)
	}

	// Xor
	a.Store(0xFF)
	old = a.Xor(0x0F)
	if old != 0xFF {
		t.Fatalf("Xor old: got %d, want 255", old)
	}
	if got := a.Load(); got != 0xF0 {
		t.Fatalf("Xor new: got %d, want 240", got)
	}
}

func TestInt32MinMax(t *testing.T) {
	var a atomix.Int32

	// Max
	a.Store(10)
	old := a.Max(20)
	if old != 10 {
		t.Fatalf("Max old: got %d, want 10", old)
	}
	if got := a.Load(); got != 20 {
		t.Fatalf("Max new: got %d, want 20", got)
	}

	// Max when current >= val
	old = a.Max(15)
	if old != 20 {
		t.Fatalf("Max unchanged: got %d, want 20", old)
	}

	// Min
	a.Store(20)
	old = a.Min(10)
	if old != 20 {
		t.Fatalf("Min old: got %d, want 20", old)
	}
	if got := a.Load(); got != 10 {
		t.Fatalf("Min new: got %d, want 10", got)
	}
}

func TestInt32Concurrent(t *testing.T) {
	var a atomix.Int32
	const numGoroutines = 100
	const opsPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				a.Add(1)
			}
		}()
	}

	wg.Wait()

	expected := int32(numGoroutines * opsPerGoroutine)
	if got := a.Load(); got != expected {
		t.Fatalf("concurrent Add: got %d, want %d", got, expected)
	}
}

// =============================================================================
// Uint64 Tests
// =============================================================================

func TestUint64LoadStore(t *testing.T) {
	var a atomix.Uint64

	if got := a.Load(); got != 0 {
		t.Fatalf("zero value: got %d, want 0", got)
	}

	a.Store(0xDEADBEEFCAFEBABE)
	if got := a.Load(); got != 0xDEADBEEFCAFEBABE {
		t.Fatalf("Load: got %x, want 0xDEADBEEFCAFEBABE", got)
	}
}

func TestUint64Add(t *testing.T) {
	var a atomix.Uint64
	a.Store(100)

	// Add returns NEW value (like sync/atomic)
	got := a.Add(50)
	if got != 150 {
		t.Fatalf("Add: got %d, want 150", got)
	}
	if v := a.Load(); v != 150 {
		t.Fatalf("Load after Add: got %d, want 150", v)
	}

	// Overflow test
	a.Store(^uint64(0) - 5)
	got = a.Add(10)
	if got != 4 {
		t.Fatalf("Add overflow: got %d, want 4", got)
	}
}

func TestUint64Concurrent(t *testing.T) {
	var a atomix.Uint64
	const numGoroutines = 100
	const opsPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				a.Add(1)
			}
		}()
	}

	wg.Wait()

	expected := uint64(numGoroutines * opsPerGoroutine)
	if got := a.Load(); got != expected {
		t.Fatalf("concurrent Add: got %d, want %d", got, expected)
	}
}

// =============================================================================
// Bool Tests
// =============================================================================

func TestBoolLoadStore(t *testing.T) {
	var a atomix.Bool

	if a.Load() {
		t.Fatal("zero value should be false")
	}

	a.Store(true)
	if !a.Load() {
		t.Fatal("Load after Store(true) should be true")
	}

	a.Store(false)
	if a.Load() {
		t.Fatal("Load after Store(false) should be false")
	}
}

func TestBoolSwap(t *testing.T) {
	var a atomix.Bool
	a.Store(true)

	old := a.Swap(false)
	if !old {
		t.Fatal("Swap old should be true")
	}
	if a.Load() {
		t.Fatal("Swap new should be false")
	}
}

func TestBoolCompareAndSwap(t *testing.T) {
	var a atomix.Bool
	a.Store(true)

	// Successful CAS
	if !a.CompareAndSwap(true, false) {
		t.Fatal("CAS true->false should succeed")
	}
	if a.Load() {
		t.Fatal("after CAS should be false")
	}

	// Failed CAS
	if a.CompareAndSwap(true, true) {
		t.Fatal("CAS should fail (value is false)")
	}
}

func TestBoolConcurrent(t *testing.T) {
	var a atomix.Bool
	var counter int32
	const numGoroutines = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			// Only one goroutine should succeed
			if a.CompareAndSwap(false, true) {
				counter++
			}
		}()
	}

	wg.Wait()

	if counter != 1 {
		t.Fatalf("exactly one CAS should succeed: got %d", counter)
	}
}

// =============================================================================
// Pointer Tests
// =============================================================================

func TestPointerLoadStore(t *testing.T) {
	var a atomix.Pointer[int]

	if a.Load() != nil {
		t.Fatal("zero value should be nil")
	}

	val := 42
	a.Store(&val)
	if got := a.Load(); got != &val {
		t.Fatal("Load should return stored pointer")
	}
	if *a.Load() != 42 {
		t.Fatal("dereferenced value should be 42")
	}
}

func TestPointerSwap(t *testing.T) {
	var a atomix.Pointer[int]
	v1, v2 := 1, 2

	a.Store(&v1)
	old := a.Swap(&v2)
	if old != &v1 {
		t.Fatal("Swap should return old pointer")
	}
	if a.Load() != &v2 {
		t.Fatal("Swap should store new pointer")
	}
}

func TestPointerCompareAndSwap(t *testing.T) {
	var a atomix.Pointer[int]
	v1, v2 := 1, 2

	a.Store(&v1)

	// Successful CAS
	if !a.CompareAndSwap(&v1, &v2) {
		t.Fatal("CAS should succeed")
	}
	if a.Load() != &v2 {
		t.Fatal("CAS should store new pointer")
	}

	// Failed CAS
	if a.CompareAndSwap(&v1, nil) {
		t.Fatal("CAS should fail")
	}
}

// =============================================================================
// Uint128 Tests (require 16-byte alignment)
// =============================================================================

func TestUint128LoadStore(t *testing.T) {
	// Uint128 requires 16-byte alignment, use PlaceAlignedUint128
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)

	lo, hi := a.Load()
	if lo != 0 || hi != 0 {
		t.Fatalf("zero value: got (%d, %d), want (0, 0)", lo, hi)
	}

	a.Store(0xDEADBEEF, 0xCAFEBABE)
	lo, hi = a.Load()
	if lo != 0xDEADBEEF || hi != 0xCAFEBABE {
		t.Fatalf("Load: got (%x, %x), want (0xDEADBEEF, 0xCAFEBABE)", lo, hi)
	}
}

func TestUint128CompareAndSwap(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	a.Store(10, 20)

	// Successful CAS
	if !a.CompareAndSwap(10, 20, 30, 40) {
		t.Fatal("CAS should succeed")
	}
	lo, hi := a.Load()
	if lo != 30 || hi != 40 {
		t.Fatalf("CAS result: got (%d, %d), want (30, 40)", lo, hi)
	}

	// Failed CAS
	if a.CompareAndSwap(10, 20, 50, 60) {
		t.Fatal("CAS should fail")
	}
}

func TestUint128CompareExchange(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	a.Store(10, 20)

	// Successful exchange
	lo, hi := a.CompareExchange(10, 20, 30, 40)
	if lo != 10 || hi != 20 {
		t.Fatalf("Cax old: got (%d, %d), want (10, 20)", lo, hi)
	}
	lo, hi = a.Load()
	if lo != 30 || hi != 40 {
		t.Fatalf("Cax new: got (%d, %d), want (30, 40)", lo, hi)
	}

	// Failed exchange returns current value
	lo, hi = a.CompareExchange(10, 20, 50, 60)
	if lo != 30 || hi != 40 {
		t.Fatalf("Cax failed: got (%d, %d), want (30, 40)", lo, hi)
	}
}

func TestUint128Add(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	a.Store(10, 0)

	newLo, newHi := a.Add(5, 0)
	if newLo != 15 || newHi != 0 {
		t.Fatalf("Add: got (%d, %d), want (15, 0)", newLo, newHi)
	}

	// Test carry from low to high
	a.Store(^uint64(0), 0) // max uint64
	newLo, newHi = a.Add(1, 0)
	if newLo != 0 || newHi != 1 {
		t.Fatalf("Add with carry: got (%d, %d), want (0, 1)", newLo, newHi)
	}
}

func TestUint128Sub(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	a.Store(10, 5)

	newLo, newHi := a.Sub(3, 2)
	if newLo != 7 || newHi != 3 {
		t.Fatalf("Sub: got (%d, %d), want (7, 3)", newLo, newHi)
	}

	// Test borrow from high to low
	a.Store(0, 1)
	newLo, newHi = a.Sub(1, 0)
	if newLo != ^uint64(0) || newHi != 0 {
		t.Fatalf("Sub with borrow: got (%d, %d), want (max, 0)", newLo, newHi)
	}
}

func TestUint128IncDec(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	a.Store(10, 0)

	newLo, _ := a.Inc()
	if newLo != 11 {
		t.Fatalf("Inc: got %d, want 11", newLo)
	}

	newLo, _ = a.Dec()
	if newLo != 10 {
		t.Fatalf("Dec: got %d, want 10", newLo)
	}
}

func TestUint128Comparisons(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)

	// Equal
	a.Store(10, 20)
	if !a.Equal(10, 20) {
		t.Fatal("Equal should return true")
	}
	if a.Equal(10, 21) {
		t.Fatal("Equal should return false")
	}

	// Less
	a.Store(10, 5)
	if !a.Less(10, 6) {
		t.Fatal("Less: (10,5) < (10,6)")
	}
	if !a.Less(11, 5) {
		t.Fatal("Less: (10,5) < (11,5)")
	}
	if a.Less(10, 4) {
		t.Fatal("Not Less: (10,5) not < (10,4)")
	}

	// Greater
	a.Store(10, 5)
	if !a.Greater(10, 4) {
		t.Fatal("Greater: (10,5) > (10,4)")
	}
	if a.Greater(10, 6) {
		t.Fatal("Not Greater: (10,5) not > (10,6)")
	}
}

func TestUint128Concurrent(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	const numGoroutines = 50
	const opsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				a.Inc()
			}
		}()
	}

	wg.Wait()

	expected := uint64(numGoroutines * opsPerGoroutine)
	lo, hi := a.Load()
	if lo != expected || hi != 0 {
		t.Fatalf("concurrent Inc: got (%d, %d), want (%d, 0)", lo, hi, expected)
	}
}

// =============================================================================
// Barrier Tests
// =============================================================================

func TestBarriers(t *testing.T) {
	// Barriers should not panic
	atomix.BarrierAcquire()
	atomix.BarrierRelease()
	atomix.BarrierAcqRel()
}

// =============================================================================
// Uintptr Tests
// =============================================================================

func TestUintptrBasic(t *testing.T) {
	var a atomix.Uintptr

	if got := a.Load(); got != 0 {
		t.Fatalf("zero value: got %d, want 0", got)
	}

	a.Store(0xDEADBEEF)
	if got := a.Load(); got != 0xDEADBEEF {
		t.Fatalf("Load: got %x, want 0xDEADBEEF", got)
	}

	old := a.Swap(0xCAFEBABE)
	if old != 0xDEADBEEF {
		t.Fatalf("Swap old: got %x, want 0xDEADBEEF", old)
	}
}

// =============================================================================
// Int64 Tests
// =============================================================================

func TestInt64Basic(t *testing.T) {
	var a atomix.Int64

	if got := a.Load(); got != 0 {
		t.Fatalf("zero value: got %d, want 0", got)
	}

	a.Store(-9223372036854775808) // min int64
	if got := a.Load(); got != -9223372036854775808 {
		t.Fatalf("Load min: got %d", got)
	}

	a.Store(9223372036854775807) // max int64
	if got := a.Load(); got != 9223372036854775807 {
		t.Fatalf("Load max: got %d", got)
	}
}

// =============================================================================
// Uint32 Tests
// =============================================================================

func TestUint32Basic(t *testing.T) {
	var a atomix.Uint32

	if got := a.Load(); got != 0 {
		t.Fatalf("zero value: got %d, want 0", got)
	}

	a.Store(0xFFFFFFFF)
	if got := a.Load(); got != 0xFFFFFFFF {
		t.Fatalf("Load max: got %x, want 0xFFFFFFFF", got)
	}

	// Add returns NEW value (like sync/atomic); 0xFFFFFFFF + 1 = 0 (wrap)
	got := a.Add(1)
	if got != 0 {
		t.Fatalf("Add overflow: got %x, want 0", got)
	}
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkInt32LoadRelaxed(b *testing.B) {
	var a atomix.Int32
	a.Store(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.LoadRelaxed()
	}
}

func BenchmarkInt32LoadAcquire(b *testing.B) {
	var a atomix.Int32
	a.Store(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.LoadAcquire()
	}
}

func BenchmarkInt32AddAcqRel(b *testing.B) {
	var a atomix.Int32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.AddAcqRel(1)
	}
}

func BenchmarkInt32CASAcqRel(b *testing.B) {
	var a atomix.Int32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.CompareAndSwapAcqRel(int32(i), int32(i+1))
	}
}

func BenchmarkUint64AddContended(b *testing.B) {
	var a atomix.Uint64
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.Add(1)
		}
	})
}

func BenchmarkUint128CASContended(b *testing.B) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for {
				lo, hi := a.LoadRelaxed()
				if a.CompareAndSwapAcqRel(lo, hi, lo+1, hi) {
					break
				}
			}
		}
	})
}

func BenchmarkBarrierAcqRel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		atomix.BarrierAcqRel()
	}
}
