// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix_test

import (
	"sync"
	"testing"
	"unsafe"

	"code.hybscloud.com/atomix"
)

// Test basic Load/Store for each ordering
func TestMemoryOrderInt32LoadStore(t *testing.T) {
	var v int32 = 42

	// Relaxed
	if got := atomix.Relaxed.LoadInt32(&v); got != 42 {
		t.Fatalf("Relaxed.LoadInt32: got %d, want 42", got)
	}
	atomix.Relaxed.StoreInt32(&v, 100)
	if v != 100 {
		t.Fatalf("Relaxed.StoreInt32: got %d, want 100", v)
	}

	// Acquire
	if got := atomix.Acquire.LoadInt32(&v); got != 100 {
		t.Fatalf("Acquire.LoadInt32: got %d, want 100", got)
	}

	// Release
	atomix.Release.StoreInt32(&v, 200)
	if v != 200 {
		t.Fatalf("Release.StoreInt32: got %d, want 200", v)
	}

	// Unknown ordering fallback
	if got := atomix.MemoryOrder(99).LoadInt32(&v); got != 200 {
		t.Fatalf("Unknown.LoadInt32: got %d, want 200", got)
	}
	atomix.MemoryOrder(99).StoreInt32(&v, 300)
	if v != 300 {
		t.Fatalf("Unknown.StoreInt32: got %d, want 300", v)
	}
}

func TestMemoryOrderUint32LoadStore(t *testing.T) {
	var v uint32 = 42

	if got := atomix.Relaxed.LoadUint32(&v); got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
	atomix.Release.StoreUint32(&v, 100)
	if got := atomix.Acquire.LoadUint32(&v); got != 100 {
		t.Fatalf("got %d, want 100", got)
	}
}

func TestMemoryOrderInt64LoadStore(t *testing.T) {
	var v int64 = 42

	if got := atomix.Relaxed.LoadInt64(&v); got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
	atomix.Release.StoreInt64(&v, 100)
	if got := atomix.Acquire.LoadInt64(&v); got != 100 {
		t.Fatalf("got %d, want 100", got)
	}
}

func TestMemoryOrderUint64LoadStore(t *testing.T) {
	var v uint64 = 42

	if got := atomix.Relaxed.LoadUint64(&v); got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
	atomix.Release.StoreUint64(&v, 100)
	if got := atomix.Acquire.LoadUint64(&v); got != 100 {
		t.Fatalf("got %d, want 100", got)
	}
}

func TestMemoryOrderUintptrLoadStore(t *testing.T) {
	var v uintptr = 42

	if got := atomix.Relaxed.LoadUintptr(&v); got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
	atomix.Release.StoreUintptr(&v, 100)
	if got := atomix.Acquire.LoadUintptr(&v); got != 100 {
		t.Fatalf("got %d, want 100", got)
	}
}

// Test Swap operations
func TestMemoryOrderSwap(t *testing.T) {
	var v int32 = 42

	old := atomix.Relaxed.SwapInt32(&v, 100)
	if old != 42 {
		t.Fatalf("Relaxed.SwapInt32: got old=%d, want 42", old)
	}
	if v != 100 {
		t.Fatalf("Relaxed.SwapInt32: got v=%d, want 100", v)
	}

	old = atomix.AcqRel.SwapInt32(&v, 200)
	if old != 100 {
		t.Fatalf("AcqRel.SwapInt32: got old=%d, want 100", old)
	}
}

// Test CompareAndSwap operations
func TestMemoryOrderCompareAndSwap(t *testing.T) {
	var v int32 = 42

	// Success case
	if !atomix.Relaxed.CompareAndSwapInt32(&v, 42, 100) {
		t.Fatal("Relaxed.CompareAndSwapInt32: should succeed")
	}
	if v != 100 {
		t.Fatalf("got %d, want 100", v)
	}

	// Failure case
	if atomix.Relaxed.CompareAndSwapInt32(&v, 42, 200) {
		t.Fatal("Relaxed.CompareAndSwapInt32: should fail")
	}
	if v != 100 {
		t.Fatalf("got %d, want 100", v)
	}

	// AcqRel
	if !atomix.AcqRel.CompareAndSwapInt32(&v, 100, 200) {
		t.Fatal("AcqRel.CompareAndSwapInt32: should succeed")
	}
}

// Test CompareExchange operations
func TestMemoryOrderCompareExchange(t *testing.T) {
	var v int32 = 42

	prev := atomix.Relaxed.CompareExchangeInt32(&v, 42, 100)
	if prev != 42 {
		t.Fatalf("CompareExchangeInt32: got prev=%d, want 42", prev)
	}
	if v != 100 {
		t.Fatalf("got v=%d, want 100", v)
	}

	// Failure case - returns current value
	prev = atomix.Relaxed.CompareExchangeInt32(&v, 42, 200)
	if prev != 100 {
		t.Fatalf("CompareExchangeInt32 failure: got prev=%d, want 100", prev)
	}
}

// Test Add operations
func TestMemoryOrderAdd(t *testing.T) {
	var v int32 = 42

	// Add returns NEW value (like sync/atomic)
	got := atomix.Relaxed.AddInt32(&v, 10)
	if got != 52 {
		t.Fatalf("AddInt32: got=%d, want 52", got)
	}
	if v != 52 {
		t.Fatalf("AddInt32: v=%d, want 52", v)
	}

	got = atomix.AcqRel.AddInt32(&v, -2)
	if got != 50 {
		t.Fatalf("AddInt32: got=%d, want 50", got)
	}
	if v != 50 {
		t.Fatalf("AddInt32: v=%d, want 50", v)
	}
}

// Test And/Or operations
func TestMemoryOrderAndOr(t *testing.T) {
	var v uint32 = 0xFF

	old := atomix.Relaxed.AndUint32(&v, 0x0F)
	if old != 0xFF {
		t.Fatalf("AndUint32: got old=0x%X, want 0xFF", old)
	}
	if v != 0x0F {
		t.Fatalf("AndUint32: got v=0x%X, want 0x0F", v)
	}

	old = atomix.AcqRel.OrUint32(&v, 0xF0)
	if old != 0x0F {
		t.Fatalf("OrUint32: got old=0x%X, want 0x0F", old)
	}
	if v != 0xFF {
		t.Fatalf("OrUint32: got v=0x%X, want 0xFF", v)
	}
}

// Test Bool operations
func TestMemoryOrderBool(t *testing.T) {
	var v uint32 = 0

	if atomix.Relaxed.LoadBool(&v) {
		t.Fatal("LoadBool: expected false")
	}

	atomix.Release.StoreBool(&v, true)
	if !atomix.Acquire.LoadBool(&v) {
		t.Fatal("LoadBool: expected true")
	}

	old := atomix.AcqRel.SwapBool(&v, false)
	if !old {
		t.Fatal("SwapBool: expected old=true")
	}
	if atomix.Relaxed.LoadBool(&v) {
		t.Fatal("LoadBool: expected false after swap")
	}

	if !atomix.AcqRel.CompareAndSwapBool(&v, false, true) {
		t.Fatal("CompareAndSwapBool: should succeed")
	}
	if !atomix.Relaxed.LoadBool(&v) {
		t.Fatal("LoadBool: expected true")
	}
}

// Test Pointer operations
func TestMemoryOrderPointer(t *testing.T) {
	var x, y int = 42, 100
	var p unsafe.Pointer = unsafe.Pointer(&x)

	got := atomix.Relaxed.LoadPointer(&p)
	if got != unsafe.Pointer(&x) {
		t.Fatal("LoadPointer: wrong value")
	}

	atomix.Release.StorePointer(&p, unsafe.Pointer(&y))
	got = atomix.Acquire.LoadPointer(&p)
	if got != unsafe.Pointer(&y) {
		t.Fatal("LoadPointer: wrong value after store")
	}

	old := atomix.AcqRel.SwapPointer(&p, unsafe.Pointer(&x))
	if old != unsafe.Pointer(&y) {
		t.Fatal("SwapPointer: wrong old value")
	}
}

// Test 128-bit operations
func TestMemoryOrderUint128(t *testing.T) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)

	atomix.Release.StoreUint128(v, 42, 100)
	lo, hi := atomix.Acquire.LoadUint128(v)
	if lo != 42 || hi != 100 {
		t.Fatalf("LoadUint128: got (%d, %d), want (42, 100)", lo, hi)
	}

	oldLo, oldHi := atomix.AcqRel.SwapUint128(v, 1, 2)
	if oldLo != 42 || oldHi != 100 {
		t.Fatalf("SwapUint128: got old=(%d, %d), want (42, 100)", oldLo, oldHi)
	}

	lo, hi = atomix.Relaxed.LoadUint128(v)
	if lo != 1 || hi != 2 {
		t.Fatalf("LoadUint128 after swap: got (%d, %d), want (1, 2)", lo, hi)
	}
}

func TestMemoryOrderUint128Add(t *testing.T) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)

	atomix.Relaxed.StoreUint128(v, 0, 0)

	oldLo, oldHi := atomix.Relaxed.AddUint128(v, 1, 0)
	if oldLo != 0 || oldHi != 0 {
		t.Fatalf("AddUint128: got old=(%d, %d), want (0, 0)", oldLo, oldHi)
	}

	lo, hi := atomix.Relaxed.LoadUint128(v)
	if lo != 1 || hi != 0 {
		t.Fatalf("LoadUint128: got (%d, %d), want (1, 0)", lo, hi)
	}
}

// Contention tests
func TestMemoryOrderAddContention(t *testing.T) {
	var counter int64
	var wg sync.WaitGroup

	const G = 16
	const N = 1000

	for i := 0; i < G; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				atomix.Relaxed.AddInt64(&counter, 1)
			}
		}()
	}
	wg.Wait()

	if counter != G*N {
		t.Fatalf("counter: got %d, want %d", counter, G*N)
	}
}

func TestMemoryOrderAndOrContention(t *testing.T) {
	var flags uint64 = 0
	var wg sync.WaitGroup

	const G = 8

	// Each goroutine sets a different bit
	for i := 0; i < G; i++ {
		wg.Add(1)
		go func(bit int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				atomix.AcqRel.OrUint64(&flags, 1<<bit)
			}
		}(i)
	}
	wg.Wait()

	want := uint64((1 << G) - 1)
	if flags != want {
		t.Fatalf("flags: got 0x%X, want 0x%X", flags, want)
	}
}

func TestMemoryOrderCASContention(t *testing.T) {
	var counter int32
	var wg sync.WaitGroup

	const G = 16
	const N = 100

	for i := 0; i < G; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				for {
					old := atomix.Relaxed.LoadInt32(&counter)
					if atomix.AcqRel.CompareAndSwapInt32(&counter, old, old+1) {
						break
					}
				}
			}
		}()
	}
	wg.Wait()

	if counter != G*N {
		t.Fatalf("counter: got %d, want %d", counter, G*N)
	}
}

// Coverage tests for all orderings and types

func TestMemoryOrderAllOrderingsInt32(t *testing.T) {
	var v int32 = 0

	// Swap all orderings
	atomix.Relaxed.SwapInt32(&v, 1)
	atomix.Acquire.SwapInt32(&v, 2)
	atomix.Release.SwapInt32(&v, 3)
	atomix.AcqRel.SwapInt32(&v, 4)

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapInt32(&v, 4, 5)
	atomix.Acquire.CompareAndSwapInt32(&v, 5, 6)
	atomix.Release.CompareAndSwapInt32(&v, 6, 7)
	atomix.AcqRel.CompareAndSwapInt32(&v, 7, 8)

	// CAX all orderings
	atomix.Relaxed.CompareExchangeInt32(&v, 8, 9)
	atomix.Acquire.CompareExchangeInt32(&v, 9, 10)
	atomix.Release.CompareExchangeInt32(&v, 10, 11)
	atomix.AcqRel.CompareExchangeInt32(&v, 11, 12)

	// Add all orderings
	atomix.Relaxed.AddInt32(&v, 1)
	atomix.Acquire.AddInt32(&v, 1)
	atomix.Release.AddInt32(&v, 1)
	atomix.AcqRel.AddInt32(&v, 1)

	// And/Or
	atomix.Relaxed.AndInt32(&v, -1)
	atomix.AcqRel.AndInt32(&v, -1)
	atomix.Relaxed.OrInt32(&v, 0)
	atomix.AcqRel.OrInt32(&v, 0)
}

func TestMemoryOrderAllOrderingsInt64(t *testing.T) {
	var v int64 = 0

	// Swap all orderings
	atomix.Relaxed.SwapInt64(&v, 1)
	atomix.Acquire.SwapInt64(&v, 2)
	atomix.Release.SwapInt64(&v, 3)
	atomix.AcqRel.SwapInt64(&v, 4)

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapInt64(&v, 4, 5)
	atomix.Acquire.CompareAndSwapInt64(&v, 5, 6)
	atomix.Release.CompareAndSwapInt64(&v, 6, 7)
	atomix.AcqRel.CompareAndSwapInt64(&v, 7, 8)

	// CAX all orderings
	atomix.Relaxed.CompareExchangeInt64(&v, 8, 9)
	atomix.Acquire.CompareExchangeInt64(&v, 9, 10)
	atomix.Release.CompareExchangeInt64(&v, 10, 11)
	atomix.AcqRel.CompareExchangeInt64(&v, 11, 12)

	// Add all orderings
	atomix.Relaxed.AddInt64(&v, 1)
	atomix.Acquire.AddInt64(&v, 1)
	atomix.Release.AddInt64(&v, 1)
	atomix.AcqRel.AddInt64(&v, 1)

	// And/Or
	atomix.Relaxed.AndInt64(&v, -1)
	atomix.AcqRel.AndInt64(&v, -1)
	atomix.Relaxed.OrInt64(&v, 0)
	atomix.AcqRel.OrInt64(&v, 0)
}

func TestMemoryOrderAllOrderingsUint32(t *testing.T) {
	var v uint32 = 0

	// Swap all orderings
	atomix.Relaxed.SwapUint32(&v, 1)
	atomix.Acquire.SwapUint32(&v, 2)
	atomix.Release.SwapUint32(&v, 3)
	atomix.AcqRel.SwapUint32(&v, 4)

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapUint32(&v, 4, 5)
	atomix.Acquire.CompareAndSwapUint32(&v, 5, 6)
	atomix.Release.CompareAndSwapUint32(&v, 6, 7)
	atomix.AcqRel.CompareAndSwapUint32(&v, 7, 8)

	// CAX all orderings
	atomix.Relaxed.CompareExchangeUint32(&v, 8, 9)
	atomix.Acquire.CompareExchangeUint32(&v, 9, 10)
	atomix.Release.CompareExchangeUint32(&v, 10, 11)
	atomix.AcqRel.CompareExchangeUint32(&v, 11, 12)

	// Add all orderings
	atomix.Relaxed.AddUint32(&v, 1)
	atomix.Acquire.AddUint32(&v, 1)
	atomix.Release.AddUint32(&v, 1)
	atomix.AcqRel.AddUint32(&v, 1)
}

func TestMemoryOrderAllOrderingsUint64(t *testing.T) {
	var v uint64 = 0

	// Swap all orderings
	atomix.Relaxed.SwapUint64(&v, 1)
	atomix.Acquire.SwapUint64(&v, 2)
	atomix.Release.SwapUint64(&v, 3)
	atomix.AcqRel.SwapUint64(&v, 4)

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapUint64(&v, 4, 5)
	atomix.Acquire.CompareAndSwapUint64(&v, 5, 6)
	atomix.Release.CompareAndSwapUint64(&v, 6, 7)
	atomix.AcqRel.CompareAndSwapUint64(&v, 7, 8)

	// CAX all orderings
	atomix.Relaxed.CompareExchangeUint64(&v, 8, 9)
	atomix.Acquire.CompareExchangeUint64(&v, 9, 10)
	atomix.Release.CompareExchangeUint64(&v, 10, 11)
	atomix.AcqRel.CompareExchangeUint64(&v, 11, 12)

	// Add all orderings
	atomix.Relaxed.AddUint64(&v, 1)
	atomix.Acquire.AddUint64(&v, 1)
	atomix.Release.AddUint64(&v, 1)
	atomix.AcqRel.AddUint64(&v, 1)

	// And/Or
	atomix.Relaxed.AndUint64(&v, ^uint64(0))
	atomix.AcqRel.AndUint64(&v, ^uint64(0))
}

func TestMemoryOrderAllOrderingsUintptr(t *testing.T) {
	var v uintptr = 0

	// Swap all orderings
	atomix.Relaxed.SwapUintptr(&v, 1)
	atomix.Acquire.SwapUintptr(&v, 2)
	atomix.Release.SwapUintptr(&v, 3)
	atomix.AcqRel.SwapUintptr(&v, 4)

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapUintptr(&v, 4, 5)
	atomix.Acquire.CompareAndSwapUintptr(&v, 5, 6)
	atomix.Release.CompareAndSwapUintptr(&v, 6, 7)
	atomix.AcqRel.CompareAndSwapUintptr(&v, 7, 8)

	// CAX all orderings
	atomix.Relaxed.CompareExchangeUintptr(&v, 8, 9)
	atomix.Acquire.CompareExchangeUintptr(&v, 9, 10)
	atomix.Release.CompareExchangeUintptr(&v, 10, 11)
	atomix.AcqRel.CompareExchangeUintptr(&v, 11, 12)

	// Add all orderings
	atomix.Relaxed.AddUintptr(&v, 1)
	atomix.Acquire.AddUintptr(&v, 1)
	atomix.Release.AddUintptr(&v, 1)
	atomix.AcqRel.AddUintptr(&v, 1)

	// And/Or
	atomix.Relaxed.AndUintptr(&v, ^uintptr(0))
	atomix.AcqRel.AndUintptr(&v, ^uintptr(0))
	atomix.Relaxed.OrUintptr(&v, 0)
	atomix.AcqRel.OrUintptr(&v, 0)
}

func TestMemoryOrderAllOrderingsInt128(t *testing.T) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedInt128(buf, 0)

	// Load/Store
	atomix.Relaxed.StoreInt128(v, 1, 2)
	atomix.Release.StoreInt128(v, 3, 4)
	atomix.Relaxed.LoadInt128(v)
	atomix.Acquire.LoadInt128(v)

	// Swap all orderings
	atomix.Relaxed.SwapInt128(v, 1, 0)
	atomix.Acquire.SwapInt128(v, 2, 0)
	atomix.Release.SwapInt128(v, 3, 0)
	atomix.AcqRel.SwapInt128(v, 4, 0)

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapInt128(v, 4, 0, 5, 0)
	atomix.Acquire.CompareAndSwapInt128(v, 5, 0, 6, 0)
	atomix.Release.CompareAndSwapInt128(v, 6, 0, 7, 0)
	atomix.AcqRel.CompareAndSwapInt128(v, 7, 0, 8, 0)

	// CAX all orderings
	atomix.Relaxed.CompareExchangeInt128(v, 8, 0, 9, 0)
	atomix.Acquire.CompareExchangeInt128(v, 9, 0, 10, 0)
	atomix.Release.CompareExchangeInt128(v, 10, 0, 11, 0)
	atomix.AcqRel.CompareExchangeInt128(v, 11, 0, 12, 0)

	// Add
	atomix.Relaxed.AddInt128(v, 1, 0)
	atomix.AcqRel.AddInt128(v, 1, 0)
}

func TestMemoryOrderAllOrderingsUint128(t *testing.T) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)

	// Swap all orderings
	atomix.Relaxed.SwapUint128(v, 1, 0)
	atomix.Acquire.SwapUint128(v, 2, 0)
	atomix.Release.SwapUint128(v, 3, 0)
	atomix.AcqRel.SwapUint128(v, 4, 0)

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapUint128(v, 4, 0, 5, 0)
	atomix.Acquire.CompareAndSwapUint128(v, 5, 0, 6, 0)
	atomix.Release.CompareAndSwapUint128(v, 6, 0, 7, 0)
	atomix.AcqRel.CompareAndSwapUint128(v, 7, 0, 8, 0)

	// CAX all orderings
	atomix.Relaxed.CompareExchangeUint128(v, 8, 0, 9, 0)
	atomix.Acquire.CompareExchangeUint128(v, 9, 0, 10, 0)
	atomix.Release.CompareExchangeUint128(v, 10, 0, 11, 0)
	atomix.AcqRel.CompareExchangeUint128(v, 11, 0, 12, 0)

	// Add
	atomix.AcqRel.AddUint128(v, 1, 0)
}

func TestMemoryOrderAllOrderingsBool(t *testing.T) {
	var v uint32 = 0

	// Store all orderings
	atomix.Relaxed.StoreBool(&v, true)
	atomix.Release.StoreBool(&v, false)

	// Swap all orderings
	atomix.Relaxed.SwapBool(&v, true)
	atomix.Acquire.SwapBool(&v, false)
	atomix.Release.SwapBool(&v, true)
	atomix.AcqRel.SwapBool(&v, false)

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapBool(&v, false, true)
	atomix.Acquire.CompareAndSwapBool(&v, true, false)
	atomix.Release.CompareAndSwapBool(&v, false, true)
	atomix.AcqRel.CompareAndSwapBool(&v, true, false)
}

func TestMemoryOrderAllOrderingsPointer(t *testing.T) {
	var x, y, z int
	var p unsafe.Pointer

	// Store
	atomix.Relaxed.StorePointer(&p, unsafe.Pointer(&x))
	atomix.Release.StorePointer(&p, unsafe.Pointer(&y))

	// Swap all orderings
	atomix.Relaxed.SwapPointer(&p, unsafe.Pointer(&x))
	atomix.Acquire.SwapPointer(&p, unsafe.Pointer(&y))
	atomix.Release.SwapPointer(&p, unsafe.Pointer(&z))
	atomix.AcqRel.SwapPointer(&p, unsafe.Pointer(&x))

	// CAS all orderings
	atomix.Relaxed.CompareAndSwapPointer(&p, unsafe.Pointer(&x), unsafe.Pointer(&y))
	atomix.Acquire.CompareAndSwapPointer(&p, unsafe.Pointer(&y), unsafe.Pointer(&z))
	atomix.Release.CompareAndSwapPointer(&p, unsafe.Pointer(&z), unsafe.Pointer(&x))
	atomix.AcqRel.CompareAndSwapPointer(&p, unsafe.Pointer(&x), unsafe.Pointer(&y))

	// CAX all orderings
	atomix.Relaxed.CompareExchangePointer(&p, unsafe.Pointer(&y), unsafe.Pointer(&z))
	atomix.Acquire.CompareExchangePointer(&p, unsafe.Pointer(&z), unsafe.Pointer(&x))
	atomix.Release.CompareExchangePointer(&p, unsafe.Pointer(&x), unsafe.Pointer(&y))
	atomix.AcqRel.CompareExchangePointer(&p, unsafe.Pointer(&y), unsafe.Pointer(&z))
}

// Benchmark
func BenchmarkMemoryOrderAddRelaxed(b *testing.B) {
	var counter int64
	for i := 0; i < b.N; i++ {
		atomix.Relaxed.AddInt64(&counter, 1)
	}
}

func BenchmarkMemoryOrderAddAcqRel(b *testing.B) {
	var counter int64
	for i := 0; i < b.N; i++ {
		atomix.AcqRel.AddInt64(&counter, 1)
	}
}

func BenchmarkMemoryOrderCASRelaxed(b *testing.B) {
	var v int64
	for i := 0; i < b.N; i++ {
		for {
			old := atomix.Relaxed.LoadInt64(&v)
			if atomix.Relaxed.CompareAndSwapInt64(&v, old, old+1) {
				break
			}
		}
	}
}

// Coverage tests for uncovered branches

// Test Store with all orderings (covers non-Relaxed branch)
func TestMemoryOrderStoreAllOrderings(t *testing.T) {
	var i32 int32
	atomix.Relaxed.StoreInt32(&i32, 1)
	atomix.Release.StoreInt32(&i32, 2)
	atomix.MemoryOrder(99).StoreInt32(&i32, 3) // unknown falls to Release

	var u32 uint32
	atomix.Relaxed.StoreUint32(&u32, 1)
	atomix.Release.StoreUint32(&u32, 2)
	atomix.MemoryOrder(99).StoreUint32(&u32, 3)

	var i64 int64
	atomix.Relaxed.StoreInt64(&i64, 1)
	atomix.Release.StoreInt64(&i64, 2)
	atomix.MemoryOrder(99).StoreInt64(&i64, 3)

	var u64 uint64
	atomix.Relaxed.StoreUint64(&u64, 1)
	atomix.Release.StoreUint64(&u64, 2)
	atomix.MemoryOrder(99).StoreUint64(&u64, 3)

	var uptr uintptr
	atomix.Relaxed.StoreUintptr(&uptr, 1)
	atomix.Release.StoreUintptr(&uptr, 2)
	atomix.MemoryOrder(99).StoreUintptr(&uptr, 3)
}

// Test And/Or with contention to exercise CAS retry loop
func TestMemoryOrderAndOrCASRetry(t *testing.T) {
	var v32 uint32 = 0xFFFFFFFF
	var v64 uint64 = 0xFFFFFFFFFFFFFFFF
	var wg sync.WaitGroup

	const G = 8
	const N = 100

	// Test uint32 And with contention
	for i := 0; i < G; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				atomix.Relaxed.AndUint32(&v32, 0xFFFFFFFF) // no-op but exercises loop
				atomix.AcqRel.AndUint32(&v32, 0xFFFFFFFF)
			}
		}()
	}
	wg.Wait()

	// Test uint32 Or with contention
	v32 = 0
	for i := 0; i < G; i++ {
		wg.Add(1)
		go func(bit int) {
			defer wg.Done()
			for j := 0; j < N; j++ {
				atomix.Relaxed.OrUint32(&v32, 1<<bit)
				atomix.AcqRel.OrUint32(&v32, 1<<bit)
			}
		}(i % 8)
	}
	wg.Wait()

	// Test uint64 Or with contention
	v64 = 0
	for i := 0; i < G; i++ {
		wg.Add(1)
		go func(bit int) {
			defer wg.Done()
			for j := 0; j < N; j++ {
				atomix.Relaxed.OrUint64(&v64, 1<<bit)
				atomix.AcqRel.OrUint64(&v64, 1<<bit)
			}
		}(i % 8)
	}
	wg.Wait()
}

// Test Add128 with carry (lo overflow)
func TestMemoryOrderAdd128Carry(t *testing.T) {
	buf := make([]byte, 64)

	// Test Uint128 carry
	_, vu := atomix.PlaceAlignedUint128(buf, 0)
	atomix.Relaxed.StoreUint128(vu, ^uint64(0), 0) // lo = max, hi = 0

	// Adding 1 should carry to hi
	oldLo, oldHi := atomix.Relaxed.AddUint128(vu, 1, 0)
	if oldLo != ^uint64(0) || oldHi != 0 {
		t.Fatalf("AddUint128 old: got (%d, %d), want (max, 0)", oldLo, oldHi)
	}

	lo, hi := atomix.Relaxed.LoadUint128(vu)
	if lo != 0 || hi != 1 {
		t.Fatalf("AddUint128 carry: got (%d, %d), want (0, 1)", lo, hi)
	}

	// Test with AcqRel ordering
	atomix.Release.StoreUint128(vu, ^uint64(0), 0)
	oldLo, oldHi = atomix.AcqRel.AddUint128(vu, 1, 0)
	if oldLo != ^uint64(0) || oldHi != 0 {
		t.Fatalf("AddUint128 AcqRel old: got (%d, %d)", oldLo, oldHi)
	}

	lo, hi = atomix.Acquire.LoadUint128(vu)
	if lo != 0 || hi != 1 {
		t.Fatalf("AddUint128 AcqRel carry: got (%d, %d), want (0, 1)", lo, hi)
	}

	// Test Int128 carry
	buf2 := make([]byte, 64)
	_, vi := atomix.PlaceAlignedInt128(buf2, 0)
	atomix.Relaxed.StoreInt128(vi, -1, 0) // lo = -1 (all bits set), hi = 0

	oldLoI, oldHiI := atomix.Relaxed.AddInt128(vi, 1, 0)
	if oldLoI != -1 || oldHiI != 0 {
		t.Fatalf("AddInt128 old: got (%d, %d)", oldLoI, oldHiI)
	}

	loI, hiI := atomix.Relaxed.LoadInt128(vi)
	if loI != 0 || hiI != 1 {
		t.Fatalf("AddInt128 carry: got (%d, %d), want (0, 1)", loI, hiI)
	}

	// Test with AcqRel
	atomix.Release.StoreInt128(vi, -1, 0)
	atomix.AcqRel.AddInt128(vi, 1, 0)
	loI, hiI = atomix.Acquire.LoadInt128(vi)
	if loI != 0 || hiI != 1 {
		t.Fatalf("AddInt128 AcqRel carry: got (%d, %d)", loI, hiI)
	}
}

// Test Add128 with contention to exercise CAS retry
func TestMemoryOrderAdd128Contention(t *testing.T) {
	buf := make([]byte, 64)
	_, v := atomix.PlaceAlignedUint128(buf, 0)
	atomix.Relaxed.StoreUint128(v, 0, 0)

	var wg sync.WaitGroup
	const G = 8
	const N = 100

	for i := 0; i < G; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				atomix.Relaxed.AddUint128(v, 1, 0)
			}
		}()
	}
	wg.Wait()

	lo, hi := atomix.Relaxed.LoadUint128(v)
	want := uint64(G * N)
	if lo != want || hi != 0 {
		t.Fatalf("AddUint128 contention: got (%d, %d), want (%d, 0)", lo, hi, want)
	}

	// Test AcqRel under contention
	atomix.Relaxed.StoreUint128(v, 0, 0)
	for i := 0; i < G; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < N; j++ {
				atomix.AcqRel.AddUint128(v, 1, 0)
			}
		}()
	}
	wg.Wait()

	lo, hi = atomix.Relaxed.LoadUint128(v)
	if lo != want || hi != 0 {
		t.Fatalf("AddUint128 AcqRel contention: got (%d, %d), want (%d, 0)", lo, hi, want)
	}
}

// Test Xor operations
func TestMemoryOrderXor(t *testing.T) {
	var v32 uint32 = 0xFF
	old := atomix.Relaxed.XorUint32(&v32, 0x0F)
	if old != 0xFF {
		t.Fatalf("XorUint32: got old=0x%X, want 0xFF", old)
	}
	if v32 != 0xF0 {
		t.Fatalf("XorUint32: got v=0x%X, want 0xF0", v32)
	}

	old = atomix.AcqRel.XorUint32(&v32, 0xF0)
	if old != 0xF0 {
		t.Fatalf("XorUint32 AcqRel: got old=0x%X, want 0xF0", old)
	}
	if v32 != 0 {
		t.Fatalf("XorUint32 AcqRel: got v=0x%X, want 0", v32)
	}

	var v64 uint64 = 0xFF
	old64 := atomix.Relaxed.XorUint64(&v64, 0x0F)
	if old64 != 0xFF {
		t.Fatalf("XorUint64: got old=0x%X, want 0xFF", old64)
	}
	if v64 != 0xF0 {
		t.Fatalf("XorUint64: got v=0x%X, want 0xF0", v64)
	}

	old64 = atomix.AcqRel.XorUint64(&v64, 0xF0)
	if old64 != 0xF0 {
		t.Fatalf("XorUint64 AcqRel: got old=0x%X, want 0xF0", old64)
	}
}

func TestMemoryOrderXorSigned(t *testing.T) {
	var v32 int32 = 0x7F
	old := atomix.Relaxed.XorInt32(&v32, 0x0F)
	if old != 0x7F {
		t.Fatalf("XorInt32: got old=0x%X, want 0x7F", old)
	}
	if v32 != 0x70 {
		t.Fatalf("XorInt32: got v=0x%X, want 0x70", v32)
	}

	old = atomix.AcqRel.XorInt32(&v32, 0x70)
	if old != 0x70 {
		t.Fatalf("XorInt32 AcqRel: got old=0x%X, want 0x70", old)
	}

	var v64 int64 = 0x7F
	old64 := atomix.Relaxed.XorInt64(&v64, 0x0F)
	if old64 != 0x7F {
		t.Fatalf("XorInt64: got old=0x%X, want 0x7F", old64)
	}

	old64 = atomix.AcqRel.XorInt64(&v64, 0x70)
	if old64 != 0x70 {
		t.Fatalf("XorInt64 AcqRel: got old=0x%X, want 0x70", old64)
	}
}

func TestMemoryOrderXorUintptr(t *testing.T) {
	var v uintptr = 0xFF
	old := atomix.Relaxed.XorUintptr(&v, 0x0F)
	if old != 0xFF {
		t.Fatalf("XorUintptr: got old=0x%X, want 0xFF", old)
	}
	if v != 0xF0 {
		t.Fatalf("XorUintptr: got v=0x%X, want 0xF0", v)
	}

	old = atomix.AcqRel.XorUintptr(&v, 0xF0)
	if old != 0xF0 {
		t.Fatalf("XorUintptr AcqRel: got old=0x%X, want 0xF0", old)
	}
}

// Test Max operations
func TestMemoryOrderMax(t *testing.T) {
	var v32 int32 = 10
	old := atomix.Relaxed.MaxInt32(&v32, 20)
	if old != 10 {
		t.Fatalf("MaxInt32: got old=%d, want 10", old)
	}
	if v32 != 20 {
		t.Fatalf("MaxInt32: got v=%d, want 20", v32)
	}

	// Value already max - no change
	old = atomix.AcqRel.MaxInt32(&v32, 15)
	if old != 20 {
		t.Fatalf("MaxInt32 no-op: got old=%d, want 20", old)
	}
	if v32 != 20 {
		t.Fatalf("MaxInt32 no-op: got v=%d, want 20", v32)
	}

	var v64 int64 = 10
	old64 := atomix.Relaxed.MaxInt64(&v64, 20)
	if old64 != 10 {
		t.Fatalf("MaxInt64: got old=%d, want 10", old64)
	}

	old64 = atomix.AcqRel.MaxInt64(&v64, 15)
	if old64 != 20 {
		t.Fatalf("MaxInt64 no-op: got old=%d, want 20", old64)
	}
}

func TestMemoryOrderMaxUnsigned(t *testing.T) {
	var v32 uint32 = 10
	old := atomix.Relaxed.MaxUint32(&v32, 20)
	if old != 10 {
		t.Fatalf("MaxUint32: got old=%d, want 10", old)
	}
	if v32 != 20 {
		t.Fatalf("MaxUint32: got v=%d, want 20", v32)
	}

	old = atomix.AcqRel.MaxUint32(&v32, 15)
	if old != 20 {
		t.Fatalf("MaxUint32 no-op: got old=%d, want 20", old)
	}

	var v64 uint64 = 10
	old64 := atomix.Relaxed.MaxUint64(&v64, 20)
	if old64 != 10 {
		t.Fatalf("MaxUint64: got old=%d, want 10", old64)
	}

	old64 = atomix.AcqRel.MaxUint64(&v64, 15)
	if old64 != 20 {
		t.Fatalf("MaxUint64 no-op: got old=%d, want 20", old64)
	}

	var vptr uintptr = 10
	oldptr := atomix.Relaxed.MaxUintptr(&vptr, 20)
	if oldptr != 10 {
		t.Fatalf("MaxUintptr: got old=%d, want 10", oldptr)
	}

	oldptr = atomix.AcqRel.MaxUintptr(&vptr, 15)
	if oldptr != 20 {
		t.Fatalf("MaxUintptr no-op: got old=%d, want 20", oldptr)
	}
}

// Test Min operations
func TestMemoryOrderMin(t *testing.T) {
	var v32 int32 = 20
	old := atomix.Relaxed.MinInt32(&v32, 10)
	if old != 20 {
		t.Fatalf("MinInt32: got old=%d, want 20", old)
	}
	if v32 != 10 {
		t.Fatalf("MinInt32: got v=%d, want 10", v32)
	}

	// Value already min - no change
	old = atomix.AcqRel.MinInt32(&v32, 15)
	if old != 10 {
		t.Fatalf("MinInt32 no-op: got old=%d, want 10", old)
	}
	if v32 != 10 {
		t.Fatalf("MinInt32 no-op: got v=%d, want 10", v32)
	}

	var v64 int64 = 20
	old64 := atomix.Relaxed.MinInt64(&v64, 10)
	if old64 != 20 {
		t.Fatalf("MinInt64: got old=%d, want 20", old64)
	}

	old64 = atomix.AcqRel.MinInt64(&v64, 15)
	if old64 != 10 {
		t.Fatalf("MinInt64 no-op: got old=%d, want 10", old64)
	}
}

func TestMemoryOrderMinUnsigned(t *testing.T) {
	var v32 uint32 = 20
	old := atomix.Relaxed.MinUint32(&v32, 10)
	if old != 20 {
		t.Fatalf("MinUint32: got old=%d, want 20", old)
	}
	if v32 != 10 {
		t.Fatalf("MinUint32: got v=%d, want 10", v32)
	}

	old = atomix.AcqRel.MinUint32(&v32, 15)
	if old != 10 {
		t.Fatalf("MinUint32 no-op: got old=%d, want 10", old)
	}

	var v64 uint64 = 20
	old64 := atomix.Relaxed.MinUint64(&v64, 10)
	if old64 != 20 {
		t.Fatalf("MinUint64: got old=%d, want 20", old64)
	}

	old64 = atomix.AcqRel.MinUint64(&v64, 15)
	if old64 != 10 {
		t.Fatalf("MinUint64 no-op: got old=%d, want 10", old64)
	}

	var vptr uintptr = 20
	oldptr := atomix.Relaxed.MinUintptr(&vptr, 10)
	if oldptr != 20 {
		t.Fatalf("MinUintptr: got old=%d, want 20", oldptr)
	}

	oldptr = atomix.AcqRel.MinUintptr(&vptr, 15)
	if oldptr != 10 {
		t.Fatalf("MinUintptr no-op: got old=%d, want 10", oldptr)
	}
}

// TestMemoryOrderAndAllOrderings tests And with all memory orderings.
func TestMemoryOrderAndAllOrderings(t *testing.T) {
	// Int32
	var i32 int32 = 0xFF
	old := atomix.Acquire.AndInt32(&i32, 0x0F)
	if old != 0xFF || i32 != 0x0F {
		t.Fatalf("AndInt32 Acquire: got old=0x%X, v=0x%X, want 0xFF, 0x0F", old, i32)
	}
	i32 = 0xFF
	old = atomix.Release.AndInt32(&i32, 0x0F)
	if old != 0xFF || i32 != 0x0F {
		t.Fatalf("AndInt32 Release: got old=0x%X, v=0x%X, want 0xFF, 0x0F", old, i32)
	}

	// Uint32
	var u32 uint32 = 0xFF
	oldu32 := atomix.Acquire.AndUint32(&u32, 0x0F)
	if oldu32 != 0xFF || u32 != 0x0F {
		t.Fatalf("AndUint32 Acquire: got old=0x%X, v=0x%X, want 0xFF, 0x0F", oldu32, u32)
	}
	u32 = 0xFF
	oldu32 = atomix.Release.AndUint32(&u32, 0x0F)
	if oldu32 != 0xFF || u32 != 0x0F {
		t.Fatalf("AndUint32 Release: got old=0x%X, v=0x%X, want 0xFF, 0x0F", oldu32, u32)
	}

	// Int64
	var i64 int64 = 0xFF
	old64 := atomix.Acquire.AndInt64(&i64, 0x0F)
	if old64 != 0xFF || i64 != 0x0F {
		t.Fatalf("AndInt64 Acquire: got old=0x%X, v=0x%X, want 0xFF, 0x0F", old64, i64)
	}
	i64 = 0xFF
	old64 = atomix.Release.AndInt64(&i64, 0x0F)
	if old64 != 0xFF || i64 != 0x0F {
		t.Fatalf("AndInt64 Release: got old=0x%X, v=0x%X, want 0xFF, 0x0F", old64, i64)
	}

	// Uint64
	var u64 uint64 = 0xFF
	oldu64 := atomix.Acquire.AndUint64(&u64, 0x0F)
	if oldu64 != 0xFF || u64 != 0x0F {
		t.Fatalf("AndUint64 Acquire: got old=0x%X, v=0x%X, want 0xFF, 0x0F", oldu64, u64)
	}
	u64 = 0xFF
	oldu64 = atomix.Release.AndUint64(&u64, 0x0F)
	if oldu64 != 0xFF || u64 != 0x0F {
		t.Fatalf("AndUint64 Release: got old=0x%X, v=0x%X, want 0xFF, 0x0F", oldu64, u64)
	}

	// Uintptr
	var uptr uintptr = 0xFF
	oldptr := atomix.Acquire.AndUintptr(&uptr, 0x0F)
	if oldptr != 0xFF || uptr != 0x0F {
		t.Fatalf("AndUintptr Acquire: got old=0x%X, v=0x%X, want 0xFF, 0x0F", oldptr, uptr)
	}
	uptr = 0xFF
	oldptr = atomix.Release.AndUintptr(&uptr, 0x0F)
	if oldptr != 0xFF || uptr != 0x0F {
		t.Fatalf("AndUintptr Release: got old=0x%X, v=0x%X, want 0xFF, 0x0F", oldptr, uptr)
	}
}

// TestMemoryOrderOrAllOrderings tests Or with all memory orderings.
func TestMemoryOrderOrAllOrderings(t *testing.T) {
	// Int32
	var i32 int32 = 0x0F
	old := atomix.Acquire.OrInt32(&i32, 0xF0)
	if old != 0x0F || i32 != 0xFF {
		t.Fatalf("OrInt32 Acquire: got old=0x%X, v=0x%X, want 0x0F, 0xFF", old, i32)
	}
	i32 = 0x0F
	old = atomix.Release.OrInt32(&i32, 0xF0)
	if old != 0x0F || i32 != 0xFF {
		t.Fatalf("OrInt32 Release: got old=0x%X, v=0x%X, want 0x0F, 0xFF", old, i32)
	}

	// Uint32
	var u32 uint32 = 0x0F
	oldu32 := atomix.Acquire.OrUint32(&u32, 0xF0)
	if oldu32 != 0x0F || u32 != 0xFF {
		t.Fatalf("OrUint32 Acquire: got old=0x%X, v=0x%X, want 0x0F, 0xFF", oldu32, u32)
	}
	u32 = 0x0F
	oldu32 = atomix.Release.OrUint32(&u32, 0xF0)
	if oldu32 != 0x0F || u32 != 0xFF {
		t.Fatalf("OrUint32 Release: got old=0x%X, v=0x%X, want 0x0F, 0xFF", oldu32, u32)
	}

	// Int64
	var i64 int64 = 0x0F
	old64 := atomix.Acquire.OrInt64(&i64, 0xF0)
	if old64 != 0x0F || i64 != 0xFF {
		t.Fatalf("OrInt64 Acquire: got old=0x%X, v=0x%X, want 0x0F, 0xFF", old64, i64)
	}
	i64 = 0x0F
	old64 = atomix.Release.OrInt64(&i64, 0xF0)
	if old64 != 0x0F || i64 != 0xFF {
		t.Fatalf("OrInt64 Release: got old=0x%X, v=0x%X, want 0x0F, 0xFF", old64, i64)
	}

	// Uint64
	var u64 uint64 = 0x0F
	oldu64 := atomix.Acquire.OrUint64(&u64, 0xF0)
	if oldu64 != 0x0F || u64 != 0xFF {
		t.Fatalf("OrUint64 Acquire: got old=0x%X, v=0x%X, want 0x0F, 0xFF", oldu64, u64)
	}
	u64 = 0x0F
	oldu64 = atomix.Release.OrUint64(&u64, 0xF0)
	if oldu64 != 0x0F || u64 != 0xFF {
		t.Fatalf("OrUint64 Release: got old=0x%X, v=0x%X, want 0x0F, 0xFF", oldu64, u64)
	}

	// Uintptr
	var uptr uintptr = 0x0F
	oldptr := atomix.Acquire.OrUintptr(&uptr, 0xF0)
	if oldptr != 0x0F || uptr != 0xFF {
		t.Fatalf("OrUintptr Acquire: got old=0x%X, v=0x%X, want 0x0F, 0xFF", oldptr, uptr)
	}
	uptr = 0x0F
	oldptr = atomix.Release.OrUintptr(&uptr, 0xF0)
	if oldptr != 0x0F || uptr != 0xFF {
		t.Fatalf("OrUintptr Release: got old=0x%X, v=0x%X, want 0x0F, 0xFF", oldptr, uptr)
	}
}

// TestMemoryOrderXorAllOrderings tests Xor with all memory orderings.
func TestMemoryOrderXorAllOrderings(t *testing.T) {
	// Int32
	var i32 int32 = 0xFF
	old := atomix.Acquire.XorInt32(&i32, 0x0F)
	if old != 0xFF || i32 != 0xF0 {
		t.Fatalf("XorInt32 Acquire: got old=0x%X, v=0x%X, want 0xFF, 0xF0", old, i32)
	}
	i32 = 0xF0
	old = atomix.Release.XorInt32(&i32, 0x0F)
	if old != 0xF0 || i32 != 0xFF {
		t.Fatalf("XorInt32 Release: got old=0x%X, v=0x%X, want 0xF0, 0xFF", old, i32)
	}

	// Uint32
	var u32 uint32 = 0xFF
	oldu32 := atomix.Acquire.XorUint32(&u32, 0x0F)
	if oldu32 != 0xFF || u32 != 0xF0 {
		t.Fatalf("XorUint32 Acquire: got old=0x%X, v=0x%X, want 0xFF, 0xF0", oldu32, u32)
	}
	u32 = 0xF0
	oldu32 = atomix.Release.XorUint32(&u32, 0x0F)
	if oldu32 != 0xF0 || u32 != 0xFF {
		t.Fatalf("XorUint32 Release: got old=0x%X, v=0x%X, want 0xF0, 0xFF", oldu32, u32)
	}

	// Int64
	var i64 int64 = 0xFF
	old64 := atomix.Acquire.XorInt64(&i64, 0x0F)
	if old64 != 0xFF || i64 != 0xF0 {
		t.Fatalf("XorInt64 Acquire: got old=0x%X, v=0x%X, want 0xFF, 0xF0", old64, i64)
	}
	i64 = 0xF0
	old64 = atomix.Release.XorInt64(&i64, 0x0F)
	if old64 != 0xF0 || i64 != 0xFF {
		t.Fatalf("XorInt64 Release: got old=0x%X, v=0x%X, want 0xF0, 0xFF", old64, i64)
	}

	// Uint64
	var u64 uint64 = 0xFF
	oldu64 := atomix.Acquire.XorUint64(&u64, 0x0F)
	if oldu64 != 0xFF || u64 != 0xF0 {
		t.Fatalf("XorUint64 Acquire: got old=0x%X, v=0x%X, want 0xFF, 0xF0", oldu64, u64)
	}
	u64 = 0xF0
	oldu64 = atomix.Release.XorUint64(&u64, 0x0F)
	if oldu64 != 0xF0 || u64 != 0xFF {
		t.Fatalf("XorUint64 Release: got old=0x%X, v=0x%X, want 0xF0, 0xFF", oldu64, u64)
	}

	// Uintptr
	var uptr uintptr = 0xFF
	oldptr := atomix.Acquire.XorUintptr(&uptr, 0x0F)
	if oldptr != 0xFF || uptr != 0xF0 {
		t.Fatalf("XorUintptr Acquire: got old=0x%X, v=0x%X, want 0xFF, 0xF0", oldptr, uptr)
	}
	uptr = 0xF0
	oldptr = atomix.Release.XorUintptr(&uptr, 0x0F)
	if oldptr != 0xF0 || uptr != 0xFF {
		t.Fatalf("XorUintptr Release: got old=0x%X, v=0x%X, want 0xF0, 0xFF", oldptr, uptr)
	}
}
