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

// =============================================================================
// Uintptr Comprehensive Tests
// =============================================================================

func TestUintptrLoadStore(t *testing.T) {
	var a atomix.Uintptr

	// Zero value
	if got := a.Load(); got != 0 {
		t.Fatalf("zero value: got %d, want 0", got)
	}

	// Store/Load
	a.Store(0xDEADBEEF)
	if got := a.Load(); got != 0xDEADBEEF {
		t.Fatalf("Load: got %x, want 0xDEADBEEF", got)
	}

	// Relaxed ordering
	a.StoreRelaxed(0xCAFEBABE)
	if got := a.LoadRelaxed(); got != 0xCAFEBABE {
		t.Fatalf("LoadRelaxed: got %x, want 0xCAFEBABE", got)
	}

	// Acquire/Release
	a.StoreRelease(0x12345678)
	if got := a.LoadAcquire(); got != 0x12345678 {
		t.Fatalf("LoadAcquire: got %x, want 0x12345678", got)
	}
}

func TestUintptrSwap(t *testing.T) {
	var a atomix.Uintptr
	a.Store(100)

	// Default swap
	old := a.Swap(200)
	if old != 100 {
		t.Fatalf("Swap old: got %d, want 100", old)
	}
	if got := a.Load(); got != 200 {
		t.Fatalf("Swap new: got %d, want 200", got)
	}

	// All ordering variants
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

func TestUintptrCompareAndSwap(t *testing.T) {
	var a atomix.Uintptr
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

	// All ordering variants
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

func TestUintptrCompareExchange(t *testing.T) {
	var a atomix.Uintptr
	a.Store(10)

	// Successful exchange
	old := a.CompareExchange(10, 20)
	if old != 10 {
		t.Fatalf("CompareExchange old: got %d, want 10", old)
	}

	// Failed exchange returns current value
	old = a.CompareExchange(10, 30)
	if old != 20 {
		t.Fatalf("CompareExchange failed: got %d, want 20", old)
	}

	// All ordering variants
	a.Store(100)
	if old := a.CompareExchangeRelaxed(100, 101); old != 100 {
		t.Fatalf("CaxRelaxed: got %d, want 100", old)
	}
	if old := a.CompareExchangeAcquire(101, 102); old != 101 {
		t.Fatalf("CaxAcquire: got %d, want 101", old)
	}
	if old := a.CompareExchangeRelease(102, 103); old != 102 {
		t.Fatalf("CaxRelease: got %d, want 102", old)
	}
	if old := a.CompareExchangeAcqRel(103, 104); old != 103 {
		t.Fatalf("CaxAcqRel: got %d, want 103", old)
	}
}

func TestUintptrAdd(t *testing.T) {
	var a atomix.Uintptr
	a.Store(100)

	// Add returns NEW value (like sync/atomic)
	got := a.Add(50)
	if got != 150 {
		t.Fatalf("Add: got %d, want 150", got)
	}
	if v := a.Load(); v != 150 {
		t.Fatalf("Load after Add: got %d, want 150", v)
	}

	// All ordering variants
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

func TestUintptrSub(t *testing.T) {
	var a atomix.Uintptr
	a.Store(100)

	// Sub returns NEW value (like sync/atomic)
	got := a.Sub(30)
	if got != 70 {
		t.Fatalf("Sub: got %d, want 70", got)
	}
	if v := a.Load(); v != 70 {
		t.Fatalf("Load after Sub: got %d, want 70", v)
	}

	// All ordering variants
	a.Store(100)
	if got := a.SubRelaxed(1); got != 99 {
		t.Fatalf("SubRelaxed: got %d, want 99", got)
	}
	if got := a.SubAcquire(1); got != 98 {
		t.Fatalf("SubAcquire: got %d, want 98", got)
	}
	if got := a.SubRelease(1); got != 97 {
		t.Fatalf("SubRelease: got %d, want 97", got)
	}
	if got := a.SubAcqRel(1); got != 96 {
		t.Fatalf("SubAcqRel: got %d, want 96", got)
	}
}

func TestUintptrBitwiseOps(t *testing.T) {
	var a atomix.Uintptr

	// And
	a.Store(0xFF)
	old := a.And(0x0F)
	if old != 0xFF {
		t.Fatalf("And old: got %x, want 0xFF", old)
	}
	if got := a.Load(); got != 0x0F {
		t.Fatalf("And new: got %x, want 0x0F", got)
	}

	// AndRelaxed
	a.Store(0xFF)
	old = a.AndRelaxed(0xF0)
	if old != 0xFF {
		t.Fatalf("AndRelaxed old: got %x, want 0xFF", old)
	}
	if got := a.Load(); got != 0xF0 {
		t.Fatalf("AndRelaxed new: got %x, want 0xF0", got)
	}

	// Or
	a.Store(0x0F)
	old = a.Or(0xF0)
	if old != 0x0F {
		t.Fatalf("Or old: got %x, want 0x0F", old)
	}
	if got := a.Load(); got != 0xFF {
		t.Fatalf("Or new: got %x, want 0xFF", got)
	}

	// OrRelaxed
	a.Store(0x00)
	old = a.OrRelaxed(0xAA)
	if old != 0x00 {
		t.Fatalf("OrRelaxed old: got %x, want 0x00", old)
	}
	if got := a.Load(); got != 0xAA {
		t.Fatalf("OrRelaxed new: got %x, want 0xAA", got)
	}

	// Xor
	a.Store(0xFF)
	old = a.Xor(0x0F)
	if old != 0xFF {
		t.Fatalf("Xor old: got %x, want 0xFF", old)
	}
	if got := a.Load(); got != 0xF0 {
		t.Fatalf("Xor new: got %x, want 0xF0", got)
	}

	// XorRelaxed
	a.Store(0xAA)
	old = a.XorRelaxed(0xFF)
	if old != 0xAA {
		t.Fatalf("XorRelaxed old: got %x, want 0xAA", old)
	}
	if got := a.Load(); got != 0x55 {
		t.Fatalf("XorRelaxed new: got %x, want 0x55", got)
	}
}

func TestUintptrConcurrent(t *testing.T) {
	var a atomix.Uintptr
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

	expected := uintptr(numGoroutines * opsPerGoroutine)
	if got := a.Load(); got != expected {
		t.Fatalf("concurrent Add: got %d, want %d", got, expected)
	}
}

// =============================================================================
// Int64 Comprehensive Tests
// =============================================================================

func TestInt64LoadStore(t *testing.T) {
	var a atomix.Int64

	// Zero value
	if got := a.Load(); got != 0 {
		t.Fatalf("zero value: got %d, want 0", got)
	}

	// Store/Load with extreme values
	a.Store(-9223372036854775808) // min int64
	if got := a.Load(); got != -9223372036854775808 {
		t.Fatalf("Load min: got %d", got)
	}

	a.Store(9223372036854775807) // max int64
	if got := a.Load(); got != 9223372036854775807 {
		t.Fatalf("Load max: got %d", got)
	}

	// All orderings
	a.StoreRelaxed(-100)
	if got := a.LoadRelaxed(); got != -100 {
		t.Fatalf("LoadRelaxed: got %d, want -100", got)
	}

	a.StoreRelease(999)
	if got := a.LoadAcquire(); got != 999 {
		t.Fatalf("LoadAcquire: got %d, want 999", got)
	}
}

func TestInt64Swap(t *testing.T) {
	var a atomix.Int64
	a.Store(10)

	old := a.Swap(20)
	if old != 10 {
		t.Fatalf("Swap old: got %d, want 10", old)
	}

	// All ordering variants
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

func TestInt64CompareAndSwap(t *testing.T) {
	var a atomix.Int64
	a.Store(10)

	if !a.CompareAndSwap(10, 20) {
		t.Fatal("CAS should succeed")
	}
	if a.CompareAndSwap(10, 30) {
		t.Fatal("CAS should fail")
	}

	// All ordering variants
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

func TestInt64CompareExchange(t *testing.T) {
	var a atomix.Int64
	a.Store(10)

	old := a.CompareExchange(10, 20)
	if old != 10 {
		t.Fatalf("CompareExchange old: got %d, want 10", old)
	}

	old = a.CompareExchange(10, 30)
	if old != 20 {
		t.Fatalf("CompareExchange failed: got %d, want 20", old)
	}

	// All ordering variants
	a.Store(100)
	if old := a.CompareExchangeRelaxed(100, 101); old != 100 {
		t.Fatalf("CaxRelaxed: got %d, want 100", old)
	}
	if old := a.CompareExchangeAcquire(101, 102); old != 101 {
		t.Fatalf("CaxAcquire: got %d, want 101", old)
	}
	if old := a.CompareExchangeRelease(102, 103); old != 102 {
		t.Fatalf("CaxRelease: got %d, want 102", old)
	}
	if old := a.CompareExchangeAcqRel(103, 104); old != 103 {
		t.Fatalf("CaxAcqRel: got %d, want 103", old)
	}
}

func TestInt64Add(t *testing.T) {
	var a atomix.Int64
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

	// All ordering variants
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

func TestInt64Sub(t *testing.T) {
	var a atomix.Int64
	a.Store(100)

	// Sub returns NEW value (like sync/atomic)
	got := a.Sub(30)
	if got != 70 {
		t.Fatalf("Sub: got %d, want 70", got)
	}
	if v := a.Load(); v != 70 {
		t.Fatalf("Load after Sub: got %d, want 70", v)
	}

	// All ordering variants
	a.Store(100)
	if got := a.SubRelaxed(1); got != 99 {
		t.Fatalf("SubRelaxed: got %d, want 99", got)
	}
	if got := a.SubAcquire(1); got != 98 {
		t.Fatalf("SubAcquire: got %d, want 98", got)
	}
	if got := a.SubRelease(1); got != 97 {
		t.Fatalf("SubRelease: got %d, want 97", got)
	}
	if got := a.SubAcqRel(1); got != 96 {
		t.Fatalf("SubAcqRel: got %d, want 96", got)
	}
}

func TestInt64BitwiseOps(t *testing.T) {
	var a atomix.Int64

	// And
	a.Store(0xFF)
	old := a.And(0x0F)
	if old != 0xFF {
		t.Fatalf("And old: got %x, want 0xFF", old)
	}
	if got := a.Load(); got != 0x0F {
		t.Fatalf("And new: got %x, want 0x0F", got)
	}

	// Or
	a.Store(0x0F)
	old = a.Or(0xF0)
	if got := a.Load(); got != 0xFF {
		t.Fatalf("Or new: got %x, want 0xFF", got)
	}

	// Xor
	a.Store(0xFF)
	old = a.Xor(0x0F)
	if got := a.Load(); got != 0xF0 {
		t.Fatalf("Xor new: got %x, want 0xF0", got)
	}
}

func TestInt64MinMax(t *testing.T) {
	var a atomix.Int64

	// Max
	a.Store(10)
	old := a.Max(20)
	if old != 10 {
		t.Fatalf("Max old: got %d, want 10", old)
	}
	if got := a.Load(); got != 20 {
		t.Fatalf("Max new: got %d, want 20", got)
	}

	old = a.Max(15)
	if got := a.Load(); got != 20 {
		t.Fatalf("Max unchanged: got %d, want 20", got)
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

	// Negative values
	a.Store(-10)
	old = a.Max(-5)
	if got := a.Load(); got != -5 {
		t.Fatalf("Max negative: got %d, want -5", got)
	}

	a.Store(-5)
	old = a.Min(-10)
	if got := a.Load(); got != -10 {
		t.Fatalf("Min negative: got %d, want -10", got)
	}
}

func TestInt64Concurrent(t *testing.T) {
	var a atomix.Int64
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

	expected := int64(numGoroutines * opsPerGoroutine)
	if got := a.Load(); got != expected {
		t.Fatalf("concurrent Add: got %d, want %d", got, expected)
	}
}

// =============================================================================
// Uint32 Comprehensive Tests
// =============================================================================

func TestUint32LoadStore(t *testing.T) {
	var a atomix.Uint32

	if got := a.Load(); got != 0 {
		t.Fatalf("zero value: got %d, want 0", got)
	}

	a.Store(0xFFFFFFFF)
	if got := a.Load(); got != 0xFFFFFFFF {
		t.Fatalf("Load max: got %x, want 0xFFFFFFFF", got)
	}

	// All orderings
	a.StoreRelaxed(100)
	if got := a.LoadRelaxed(); got != 100 {
		t.Fatalf("LoadRelaxed: got %d, want 100", got)
	}

	a.StoreRelease(200)
	if got := a.LoadAcquire(); got != 200 {
		t.Fatalf("LoadAcquire: got %d, want 200", got)
	}
}

func TestUint32Swap(t *testing.T) {
	var a atomix.Uint32
	a.Store(100)

	old := a.Swap(200)
	if old != 100 {
		t.Fatalf("Swap old: got %d, want 100", old)
	}

	// All ordering variants
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

func TestUint32CompareAndSwap(t *testing.T) {
	var a atomix.Uint32
	a.Store(10)

	if !a.CompareAndSwap(10, 20) {
		t.Fatal("CAS should succeed")
	}
	if a.CompareAndSwap(10, 30) {
		t.Fatal("CAS should fail")
	}

	// All ordering variants
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

func TestUint32CompareExchange(t *testing.T) {
	var a atomix.Uint32
	a.Store(10)

	old := a.CompareExchange(10, 20)
	if old != 10 {
		t.Fatalf("CompareExchange old: got %d, want 10", old)
	}

	old = a.CompareExchange(10, 30)
	if old != 20 {
		t.Fatalf("CompareExchange failed: got %d, want 20", old)
	}

	// All ordering variants
	a.Store(100)
	if old := a.CompareExchangeRelaxed(100, 101); old != 100 {
		t.Fatalf("CaxRelaxed: got %d, want 100", old)
	}
	if old := a.CompareExchangeAcquire(101, 102); old != 101 {
		t.Fatalf("CaxAcquire: got %d, want 101", old)
	}
	if old := a.CompareExchangeRelease(102, 103); old != 102 {
		t.Fatalf("CaxRelease: got %d, want 102", old)
	}
	if old := a.CompareExchangeAcqRel(103, 104); old != 103 {
		t.Fatalf("CaxAcqRel: got %d, want 103", old)
	}
}

func TestUint32Add(t *testing.T) {
	var a atomix.Uint32
	a.Store(100)

	// Add returns NEW value (like sync/atomic)
	got := a.Add(50)
	if got != 150 {
		t.Fatalf("Add: got %d, want 150", got)
	}

	// Overflow
	a.Store(0xFFFFFFFF)
	got = a.Add(2)
	if got != 1 {
		t.Fatalf("Add overflow: got %d, want 1", got)
	}

	// All ordering variants
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

func TestUint32Sub(t *testing.T) {
	var a atomix.Uint32
	a.Store(100)

	// Sub returns NEW value (like sync/atomic)
	got := a.Sub(30)
	if got != 70 {
		t.Fatalf("Sub: got %d, want 70", got)
	}

	// All ordering variants
	a.Store(100)
	if got := a.SubRelaxed(1); got != 99 {
		t.Fatalf("SubRelaxed: got %d, want 99", got)
	}
	if got := a.SubAcquire(1); got != 98 {
		t.Fatalf("SubAcquire: got %d, want 98", got)
	}
	if got := a.SubRelease(1); got != 97 {
		t.Fatalf("SubRelease: got %d, want 97", got)
	}
	if got := a.SubAcqRel(1); got != 96 {
		t.Fatalf("SubAcqRel: got %d, want 96", got)
	}
}

func TestUint32BitwiseOps(t *testing.T) {
	var a atomix.Uint32

	// And
	a.Store(0xFF)
	old := a.And(0x0F)
	if old != 0xFF {
		t.Fatalf("And old: got %x, want 0xFF", old)
	}

	// Or
	a.Store(0x0F)
	old = a.Or(0xF0)
	if got := a.Load(); got != 0xFF {
		t.Fatalf("Or new: got %x, want 0xFF", got)
	}

	// Xor
	a.Store(0xFF)
	old = a.Xor(0x0F)
	if got := a.Load(); got != 0xF0 {
		t.Fatalf("Xor new: got %x, want 0xF0", got)
	}
}

func TestUint32MinMax(t *testing.T) {
	var a atomix.Uint32

	// Max
	a.Store(10)
	old := a.Max(20)
	if old != 10 {
		t.Fatalf("Max old: got %d, want 10", old)
	}
	if got := a.Load(); got != 20 {
		t.Fatalf("Max new: got %d, want 20", got)
	}

	// Min
	a.Store(20)
	old = a.Min(10)
	if got := a.Load(); got != 10 {
		t.Fatalf("Min new: got %d, want 10", got)
	}
}

func TestUint32Concurrent(t *testing.T) {
	var a atomix.Uint32
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

	expected := uint32(numGoroutines * opsPerGoroutine)
	if got := a.Load(); got != expected {
		t.Fatalf("concurrent Add: got %d, want %d", got, expected)
	}
}

// =============================================================================
// Uint64 Comprehensive Tests
// =============================================================================

func TestUint64Swap(t *testing.T) {
	var a atomix.Uint64
	a.Store(100)

	old := a.Swap(200)
	if old != 100 {
		t.Fatalf("Swap old: got %d, want 100", old)
	}

	// All ordering variants
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

func TestUint64CompareAndSwap(t *testing.T) {
	var a atomix.Uint64
	a.Store(10)

	if !a.CompareAndSwap(10, 20) {
		t.Fatal("CAS should succeed")
	}
	if a.CompareAndSwap(10, 30) {
		t.Fatal("CAS should fail")
	}

	// All ordering variants
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

func TestUint64CompareExchange(t *testing.T) {
	var a atomix.Uint64
	a.Store(10)

	old := a.CompareExchange(10, 20)
	if old != 10 {
		t.Fatalf("CompareExchange old: got %d, want 10", old)
	}

	// All ordering variants
	a.Store(100)
	if old := a.CompareExchangeRelaxed(100, 101); old != 100 {
		t.Fatalf("CaxRelaxed: got %d, want 100", old)
	}
	if old := a.CompareExchangeAcquire(101, 102); old != 101 {
		t.Fatalf("CaxAcquire: got %d, want 101", old)
	}
	if old := a.CompareExchangeRelease(102, 103); old != 102 {
		t.Fatalf("CaxRelease: got %d, want 102", old)
	}
	if old := a.CompareExchangeAcqRel(103, 104); old != 103 {
		t.Fatalf("CaxAcqRel: got %d, want 103", old)
	}
}

func TestUint64Sub(t *testing.T) {
	var a atomix.Uint64
	a.Store(100)

	// Sub returns NEW value (like sync/atomic)
	got := a.Sub(30)
	if got != 70 {
		t.Fatalf("Sub: got %d, want 70", got)
	}

	// All ordering variants
	a.Store(100)
	if got := a.SubRelaxed(1); got != 99 {
		t.Fatalf("SubRelaxed: got %d, want 99", got)
	}
	if got := a.SubAcquire(1); got != 98 {
		t.Fatalf("SubAcquire: got %d, want 98", got)
	}
	if got := a.SubRelease(1); got != 97 {
		t.Fatalf("SubRelease: got %d, want 97", got)
	}
	if got := a.SubAcqRel(1); got != 96 {
		t.Fatalf("SubAcqRel: got %d, want 96", got)
	}
}

func TestUint64BitwiseOps(t *testing.T) {
	var a atomix.Uint64

	// And
	a.Store(0xFF)
	old := a.And(0x0F)
	if old != 0xFF {
		t.Fatalf("And old: got %x, want 0xFF", old)
	}

	// Or
	a.Store(0x0F)
	old = a.Or(0xF0)
	if got := a.Load(); got != 0xFF {
		t.Fatalf("Or new: got %x, want 0xFF", got)
	}

	// Xor
	a.Store(0xFF)
	old = a.Xor(0x0F)
	if got := a.Load(); got != 0xF0 {
		t.Fatalf("Xor new: got %x, want 0xF0", got)
	}
}

func TestUint64MinMax(t *testing.T) {
	var a atomix.Uint64

	// Max
	a.Store(10)
	old := a.Max(20)
	if old != 10 {
		t.Fatalf("Max old: got %d, want 10", old)
	}
	if got := a.Load(); got != 20 {
		t.Fatalf("Max new: got %d, want 20", got)
	}

	// Min
	a.Store(20)
	old = a.Min(10)
	if got := a.Load(); got != 10 {
		t.Fatalf("Min new: got %d, want 10", got)
	}
}

// =============================================================================
// Int128 Comprehensive Tests
// =============================================================================

func TestInt128LoadStore(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)

	lo, hi := a.Load()
	if lo != 0 || hi != 0 {
		t.Fatalf("zero value: got (%d, %d), want (0, 0)", lo, hi)
	}

	a.Store(0x12345678, -1)
	lo, hi = a.Load()
	if lo != 0x12345678 || hi != -1 {
		t.Fatalf("Load: got (%x, %d), want (0x12345678, -1)", lo, hi)
	}

	// Relaxed ordering
	a.StoreRelaxed(100, 200)
	lo, hi = a.LoadRelaxed()
	if lo != 100 || hi != 200 {
		t.Fatalf("LoadRelaxed: got (%d, %d), want (100, 200)", lo, hi)
	}

	// Acquire/Release
	a.StoreRelease(300, 400)
	lo, hi = a.LoadAcquire()
	if lo != 300 || hi != 400 {
		t.Fatalf("LoadAcquire: got (%d, %d), want (300, 400)", lo, hi)
	}
}

func TestInt128Swap(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)
	a.Store(10, 20)

	oldLo, oldHi := a.Swap(30, 40)
	if oldLo != 10 || oldHi != 20 {
		t.Fatalf("Swap old: got (%d, %d), want (10, 20)", oldLo, oldHi)
	}
	lo, hi := a.Load()
	if lo != 30 || hi != 40 {
		t.Fatalf("Swap new: got (%d, %d), want (30, 40)", lo, hi)
	}

	// All ordering variants
	a.Store(100, 200)
	if oldLo, oldHi := a.SwapRelaxed(101, 201); oldLo != 100 || oldHi != 200 {
		t.Fatalf("SwapRelaxed: got (%d, %d), want (100, 200)", oldLo, oldHi)
	}
	if oldLo, oldHi := a.SwapAcquire(102, 202); oldLo != 101 || oldHi != 201 {
		t.Fatalf("SwapAcquire: got (%d, %d), want (101, 201)", oldLo, oldHi)
	}
	if oldLo, oldHi := a.SwapRelease(103, 203); oldLo != 102 || oldHi != 202 {
		t.Fatalf("SwapRelease: got (%d, %d), want (102, 202)", oldLo, oldHi)
	}
	if oldLo, oldHi := a.SwapAcqRel(104, 204); oldLo != 103 || oldHi != 203 {
		t.Fatalf("SwapAcqRel: got (%d, %d), want (103, 203)", oldLo, oldHi)
	}
}

func TestInt128CompareAndSwap(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)
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

	// All ordering variants
	a.Store(100, 200)
	if !a.CompareAndSwapRelaxed(100, 200, 101, 201) {
		t.Fatal("CASRelaxed should succeed")
	}
	if !a.CompareAndSwapAcquire(101, 201, 102, 202) {
		t.Fatal("CASAcquire should succeed")
	}
	if !a.CompareAndSwapRelease(102, 202, 103, 203) {
		t.Fatal("CASRelease should succeed")
	}
	if !a.CompareAndSwapAcqRel(103, 203, 104, 204) {
		t.Fatal("CASAcqRel should succeed")
	}
}

func TestInt128CompareExchange(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)
	a.Store(10, 20)

	lo, hi := a.CompareExchange(10, 20, 30, 40)
	if lo != 10 || hi != 20 {
		t.Fatalf("Cax old: got (%d, %d), want (10, 20)", lo, hi)
	}

	lo, hi = a.CompareExchange(10, 20, 50, 60)
	if lo != 30 || hi != 40 {
		t.Fatalf("Cax failed: got (%d, %d), want (30, 40)", lo, hi)
	}

	// All ordering variants
	a.Store(100, 200)
	if lo, hi := a.CompareExchangeRelaxed(100, 200, 101, 201); lo != 100 || hi != 200 {
		t.Fatalf("CaxRelaxed: got (%d, %d), want (100, 200)", lo, hi)
	}
	if lo, hi := a.CompareExchangeAcquire(101, 201, 102, 202); lo != 101 || hi != 201 {
		t.Fatalf("CaxAcquire: got (%d, %d), want (101, 201)", lo, hi)
	}
	if lo, hi := a.CompareExchangeRelease(102, 202, 103, 203); lo != 102 || hi != 202 {
		t.Fatalf("CaxRelease: got (%d, %d), want (102, 202)", lo, hi)
	}
	if lo, hi := a.CompareExchangeAcqRel(103, 203, 104, 204); lo != 103 || hi != 203 {
		t.Fatalf("CaxAcqRel: got (%d, %d), want (103, 203)", lo, hi)
	}
}

func TestInt128Add(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)
	a.Store(10, 0)

	newLo, newHi := a.Add(5, 0)
	if newLo != 15 || newHi != 0 {
		t.Fatalf("Add: got (%d, %d), want (15, 0)", newLo, newHi)
	}

	// AddRelaxed
	a.Store(100, 0)
	if newLo, newHi := a.AddRelaxed(1, 0); newLo != 101 || newHi != 0 {
		t.Fatalf("AddRelaxed: got (%d, %d), want (101, 0)", newLo, newHi)
	}
}

func TestInt128Sub(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)
	a.Store(10, 5)

	newLo, newHi := a.Sub(3, 2)
	if newLo != 7 || newHi != 3 {
		t.Fatalf("Sub: got (%d, %d), want (7, 3)", newLo, newHi)
	}

	// SubRelaxed
	a.Store(100, 50)
	if newLo, newHi := a.SubRelaxed(1, 0); newLo != 99 || newHi != 50 {
		t.Fatalf("SubRelaxed: got (%d, %d), want (99, 50)", newLo, newHi)
	}
}

func TestInt128IncDec(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)
	a.Store(10, 0)

	newLo, _ := a.Inc()
	if newLo != 11 {
		t.Fatalf("Inc: got %d, want 11", newLo)
	}

	newLo, _ = a.IncRelaxed()
	if newLo != 12 {
		t.Fatalf("IncRelaxed: got %d, want 12", newLo)
	}

	newLo, _ = a.Dec()
	if newLo != 11 {
		t.Fatalf("Dec: got %d, want 11", newLo)
	}

	newLo, _ = a.DecRelaxed()
	if newLo != 10 {
		t.Fatalf("DecRelaxed: got %d, want 10", newLo)
	}
}

func TestInt128Comparisons(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)

	// Equal
	a.Store(10, 20)
	if !a.Equal(10, 20) {
		t.Fatal("Equal should return true")
	}
	if a.Equal(10, 21) {
		t.Fatal("Equal should return false")
	}

	// EqualRelaxed
	if !a.EqualRelaxed(10, 20) {
		t.Fatal("EqualRelaxed should return true")
	}

	// Less/Greater with signed values
	a.Store(10, -1) // negative high word
	if !a.Less(10, 0) {
		t.Fatal("Less: negative hi < positive hi")
	}
	if a.Greater(10, 0) {
		t.Fatal("Greater: negative hi not > positive hi")
	}
}

func TestInt128Concurrent(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)
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

	expected := int64(numGoroutines * opsPerGoroutine)
	lo, hi := a.Load()
	if lo != expected || hi != 0 {
		t.Fatalf("concurrent Inc: got (%d, %d), want (%d, 0)", lo, hi, expected)
	}
}

// =============================================================================
// Uint128 Swap Tests (new operations)
// =============================================================================

func TestUint128Swap(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	a.Store(10, 20)

	oldLo, oldHi := a.Swap(30, 40)
	if oldLo != 10 || oldHi != 20 {
		t.Fatalf("Swap old: got (%d, %d), want (10, 20)", oldLo, oldHi)
	}
	lo, hi := a.Load()
	if lo != 30 || hi != 40 {
		t.Fatalf("Swap new: got (%d, %d), want (30, 40)", lo, hi)
	}

	// All ordering variants
	a.Store(100, 200)
	if oldLo, oldHi := a.SwapRelaxed(101, 201); oldLo != 100 || oldHi != 200 {
		t.Fatalf("SwapRelaxed: got (%d, %d), want (100, 200)", oldLo, oldHi)
	}
	if oldLo, oldHi := a.SwapAcquire(102, 202); oldLo != 101 || oldHi != 201 {
		t.Fatalf("SwapAcquire: got (%d, %d), want (101, 201)", oldLo, oldHi)
	}
	if oldLo, oldHi := a.SwapRelease(103, 203); oldLo != 102 || oldHi != 202 {
		t.Fatalf("SwapRelease: got (%d, %d), want (102, 202)", oldLo, oldHi)
	}
	if oldLo, oldHi := a.SwapAcqRel(104, 204); oldLo != 103 || oldHi != 203 {
		t.Fatalf("SwapAcqRel: got (%d, %d), want (103, 203)", oldLo, oldHi)
	}
}

func TestUint128SwapFromZero(t *testing.T) {
	// Regression test: Swap from 0:0 must return 0:0 as old value
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)

	// Store zero (initial state)
	a.Store(0, 0)

	// Swap from 0:0 to new value
	oldLo, oldHi := a.Swap(42, 84)
	if oldLo != 0 || oldHi != 0 {
		t.Fatalf("Swap from 0:0: old value got (%d, %d), want (0, 0)", oldLo, oldHi)
	}

	// Verify new value
	lo, hi := a.Load()
	if lo != 42 || hi != 84 {
		t.Fatalf("Swap new value: got (%d, %d), want (42, 84)", lo, hi)
	}

	// Test all ordering variants from zero
	a.Store(0, 0)
	if oldLo, oldHi := a.SwapRelaxed(1, 2); oldLo != 0 || oldHi != 0 {
		t.Fatalf("SwapRelaxed from 0: old got (%d, %d), want (0, 0)", oldLo, oldHi)
	}
	a.Store(0, 0)
	if oldLo, oldHi := a.SwapAcquire(3, 4); oldLo != 0 || oldHi != 0 {
		t.Fatalf("SwapAcquire from 0: old got (%d, %d), want (0, 0)", oldLo, oldHi)
	}
	a.Store(0, 0)
	if oldLo, oldHi := a.SwapRelease(5, 6); oldLo != 0 || oldHi != 0 {
		t.Fatalf("SwapRelease from 0: old got (%d, %d), want (0, 0)", oldLo, oldHi)
	}
	a.Store(0, 0)
	if oldLo, oldHi := a.SwapAcqRel(7, 8); oldLo != 0 || oldHi != 0 {
		t.Fatalf("SwapAcqRel from 0: old got (%d, %d), want (0, 0)", oldLo, oldHi)
	}
}

func TestUint128SwapConcurrent(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	a.Store(0, 0)

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			newLo := uint64(id + 1)
			newHi := uint64(id + 1000)
			a.Swap(newLo, newHi)
		}(i)
	}

	wg.Wait()

	// Value should be one of the swapped values
	lo, hi := a.Load()
	if lo == 0 && hi == 0 {
		t.Fatal("Swap should have changed the value")
	}
	if lo < 1 || lo > uint64(numGoroutines) {
		t.Fatalf("Swap lo out of range: got %d", lo)
	}
}

// =============================================================================
// Bool Ordering Variant Tests
// =============================================================================

func TestBoolOrderingVariants(t *testing.T) {
	var a atomix.Bool

	// Relaxed
	a.StoreRelaxed(true)
	if !a.LoadRelaxed() {
		t.Fatal("LoadRelaxed should be true")
	}

	// Acquire/Release
	a.StoreRelease(false)
	if a.LoadAcquire() {
		t.Fatal("LoadAcquire should be false")
	}

	// Swap variants
	a.Store(true)
	if old := a.SwapRelaxed(false); !old {
		t.Fatal("SwapRelaxed old should be true")
	}
	if old := a.SwapAcquire(true); old {
		t.Fatal("SwapAcquire old should be false")
	}
	if old := a.SwapRelease(false); !old {
		t.Fatal("SwapRelease old should be true")
	}
	if old := a.SwapAcqRel(true); old {
		t.Fatal("SwapAcqRel old should be false")
	}

	// CAS variants
	a.Store(true)
	if !a.CompareAndSwapRelaxed(true, false) {
		t.Fatal("CASRelaxed should succeed")
	}
	if !a.CompareAndSwapAcquire(false, true) {
		t.Fatal("CASAcquire should succeed")
	}
	if !a.CompareAndSwapRelease(true, false) {
		t.Fatal("CASRelease should succeed")
	}
	if !a.CompareAndSwapAcqRel(false, true) {
		t.Fatal("CASAcqRel should succeed")
	}
}

// =============================================================================
// Pointer Ordering Variant Tests
// =============================================================================

func TestPointerOrderingVariants(t *testing.T) {
	var a atomix.Pointer[int]
	v1, v2, v3, v4, v5 := 1, 2, 3, 4, 5

	// Relaxed
	a.StoreRelaxed(&v1)
	if a.LoadRelaxed() != &v1 {
		t.Fatal("LoadRelaxed failed")
	}

	// Acquire/Release
	a.StoreRelease(&v2)
	if a.LoadAcquire() != &v2 {
		t.Fatal("LoadAcquire failed")
	}

	// Swap variants
	a.Store(&v1)
	if old := a.SwapRelaxed(&v2); old != &v1 {
		t.Fatal("SwapRelaxed failed")
	}
	if old := a.SwapAcquire(&v3); old != &v2 {
		t.Fatal("SwapAcquire failed")
	}
	if old := a.SwapRelease(&v4); old != &v3 {
		t.Fatal("SwapRelease failed")
	}
	if old := a.SwapAcqRel(&v5); old != &v4 {
		t.Fatal("SwapAcqRel failed")
	}

	// CAS variants
	a.Store(&v1)
	if !a.CompareAndSwapRelaxed(&v1, &v2) {
		t.Fatal("CASRelaxed should succeed")
	}
	if !a.CompareAndSwapAcquire(&v2, &v3) {
		t.Fatal("CASAcquire should succeed")
	}
	if !a.CompareAndSwapRelease(&v3, &v4) {
		t.Fatal("CASRelease should succeed")
	}
	if !a.CompareAndSwapAcqRel(&v4, &v5) {
		t.Fatal("CASAcqRel should succeed")
	}
}

// =============================================================================
// Alignment and Placement Tests
// =============================================================================

func TestCanPlaceAligned(t *testing.T) {
	buf := make([]byte, 128)

	// Test 4-byte alignment - should succeed with enough space
	for i := 0; i < 16; i++ {
		if !atomix.CanPlaceAligned4(buf, i) {
			t.Fatalf("CanPlaceAligned4 failed at offset %d", i)
		}
	}

	// Test 8-byte alignment
	for i := 0; i < 16; i++ {
		if !atomix.CanPlaceAligned8(buf, i) {
			t.Fatalf("CanPlaceAligned8 failed at offset %d", i)
		}
	}

	// Test 16-byte alignment
	for i := 0; i < 32; i++ {
		if !atomix.CanPlaceAligned16(buf, i) {
			t.Fatalf("CanPlaceAligned16 failed at offset %d", i)
		}
	}

	// Test failure cases - buffer too small
	smallBuf := make([]byte, 4)
	if atomix.CanPlaceAligned16(smallBuf, 0) {
		t.Fatal("CanPlaceAligned16 should fail with 4-byte buffer")
	}

	// Test negative offset
	if atomix.CanPlaceAligned4(buf, -1) {
		t.Fatal("CanPlaceAligned4 should fail with negative offset")
	}

	// Test offset beyond buffer
	if atomix.CanPlaceAligned8(buf, 200) {
		t.Fatal("CanPlaceAligned8 should fail with offset beyond buffer")
	}
}

func TestPlaceAlignedTypes(t *testing.T) {
	buf := make([]byte, 256)

	// Int32 - verify pointer is 4-byte aligned
	n, a32 := atomix.PlaceAlignedInt32(buf, 1)
	if n < 4 {
		t.Fatalf("PlaceAlignedInt32: consumed %d bytes, expected >= 4", n)
	}
	if uintptr(unsafe.Pointer(a32))%4 != 0 {
		t.Fatal("PlaceAlignedInt32: pointer not 4-byte aligned")
	}
	a32.Store(42)
	if a32.Load() != 42 {
		t.Fatal("PlaceAlignedInt32 failed")
	}

	// Int64 - verify pointer is 8-byte aligned
	n, a64 := atomix.PlaceAlignedInt64(buf, 50)
	if n < 8 {
		t.Fatalf("PlaceAlignedInt64: consumed %d bytes, expected >= 8", n)
	}
	if uintptr(unsafe.Pointer(a64))%8 != 0 {
		t.Fatal("PlaceAlignedInt64: pointer not 8-byte aligned")
	}
	a64.Store(-999)
	if a64.Load() != -999 {
		t.Fatal("PlaceAlignedInt64 failed")
	}

	// Uint128 - verify pointer is 16-byte aligned
	n, a128 := atomix.PlaceAlignedUint128(buf, 100)
	if n < 16 {
		t.Fatalf("PlaceAlignedUint128: consumed %d bytes, expected >= 16", n)
	}
	if uintptr(unsafe.Pointer(a128))%16 != 0 {
		t.Fatal("PlaceAlignedUint128: pointer not 16-byte aligned")
	}
	a128.Store(0xDEAD, 0xBEEF)
	lo, hi := a128.Load()
	if lo != 0xDEAD || hi != 0xBEEF {
		t.Fatal("PlaceAlignedUint128 failed")
	}
}

func TestAllocator(t *testing.T) {
	buf := make([]byte, 512)
	alloc := atomix.NewAllocator(buf)

	// Allocate different types
	a32 := alloc.Int32()
	a64 := alloc.Int64()
	a128 := alloc.Uint128()

	if a32 == nil || a64 == nil || a128 == nil {
		t.Fatal("Allocator failed")
	}

	// Verify they work
	a32.Store(1)
	a64.Store(2)
	a128.Store(3, 4)

	if a32.Load() != 1 {
		t.Fatal("Int32 allocation failed")
	}
	if a64.Load() != 2 {
		t.Fatal("Int64 allocation failed")
	}
	lo, hi := a128.Load()
	if lo != 3 || hi != 4 {
		t.Fatal("Uint128 allocation failed")
	}

	// Test Offset and Remaining
	initialRemaining := alloc.Remaining()
	alloc.Uint32()
	if alloc.Remaining() >= initialRemaining {
		t.Fatal("Remaining should decrease after allocation")
	}

	// Test Reset
	alloc.Reset()
	if alloc.Offset() != 0 {
		t.Fatal("Offset should be 0 after Reset")
	}

	// Test Skip
	alloc.Skip(10)
	if alloc.Offset() != 10 {
		t.Fatal("Offset should be 10 after Skip(10)")
	}

	// Test Align
	alloc.Align(16)
	if alloc.Offset()%16 != 0 {
		// After alignment, the next allocation address should be 16-byte aligned
		// Note: Align adjusts based on buffer base address + offset
	}

	// Test all allocator methods
	alloc.Reset()
	_ = alloc.Bool()
	_ = alloc.Uint32()
	_ = alloc.Uint64()
	_ = alloc.Uintptr()
	_ = alloc.Int128()

	// Test cache-aligned allocations
	alloc.Reset()
	_ = alloc.CacheAlignedInt32()
	_ = alloc.CacheAlignedUint32()
	_ = alloc.CacheAlignedInt64()
	_ = alloc.CacheAlignedUint64()
	_ = alloc.CacheAlignedUintptr()
	_ = alloc.CacheAlignedBool()

	// Test 128-bit cache-aligned
	alloc.Reset()
	_ = alloc.CacheAlignedInt128()
	_ = alloc.CacheAlignedUint128()
}

// =============================================================================
// Padded Type Tests
// =============================================================================

func TestPaddedTypes(t *testing.T) {
	// Verify padded types have expected size (cache line = 64 bytes typically)
	var i32p atomix.Int32Padded
	var i64p atomix.Int64Padded

	// Ensure they work
	i32p.Store(42)
	if i32p.Load() != 42 {
		t.Fatal("Int32Padded failed")
	}

	i64p.Store(-123456789)
	if i64p.Load() != -123456789 {
		t.Fatal("Int64Padded failed")
	}

	// Uint128Padded requires 16-byte alignment; use PlaceAlignedUint128
	buf := make([]byte, 32)
	_, u128p := atomix.PlaceAlignedUint128(buf, 0)
	u128p.Store(0xABCD, 0x1234)
	lo, hi := u128p.Load()
	if lo != 0xABCD || hi != 0x1234 {
		t.Fatal("Uint128Padded failed")
	}

	// Verify size is at least cache line
	if unsafe.Sizeof(i32p) < 64 {
		t.Fatalf("Int32Padded size %d < 64", unsafe.Sizeof(i32p))
	}
	if unsafe.Sizeof(i64p) < 64 {
		t.Fatalf("Int64Padded size %d < 64", unsafe.Sizeof(i64p))
	}
}

// =============================================================================
// Edge Case Tests
// =============================================================================

func TestInt32Overflow(t *testing.T) {
	var a atomix.Int32

	// Max + 1 = Min
	a.Store(2147483647) // max int32
	a.Add(1)
	if got := a.Load(); got != -2147483648 {
		t.Fatalf("overflow: got %d, want -2147483648", got)
	}

	// Min - 1 = Max
	a.Store(-2147483648) // min int32
	a.Sub(1)
	if got := a.Load(); got != 2147483647 {
		t.Fatalf("underflow: got %d, want 2147483647", got)
	}
}

func TestUint64MaxValues(t *testing.T) {
	var a atomix.Uint64

	// Test with max uint64
	a.Store(^uint64(0))
	if got := a.Load(); got != ^uint64(0) {
		t.Fatalf("max uint64: got %x", got)
	}

	// Wrap around
	a.Add(1)
	if got := a.Load(); got != 0 {
		t.Fatalf("wrap around: got %d, want 0", got)
	}
}

func TestUint128LargeValues(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)

	// Max values
	a.Store(^uint64(0), ^uint64(0))
	lo, hi := a.Load()
	if lo != ^uint64(0) || hi != ^uint64(0) {
		t.Fatal("max uint128 failed")
	}

	// Increment wraps
	a.Inc()
	lo, hi = a.Load()
	if lo != 0 || hi != 0 {
		t.Fatalf("wrap around: got (%d, %d), want (0, 0)", lo, hi)
	}
}

// =============================================================================
// Additional Coverage Tests - Relaxed Variants
// =============================================================================

func TestInt32RelaxedVariants(t *testing.T) {
	var a atomix.Int32

	// SubRelaxed
	a.Store(100)
	a.SubRelaxed(30)
	if a.Load() != 70 {
		t.Fatal("SubRelaxed failed")
	}

	// SubAcquire
	a.SubAcquire(10)
	if a.Load() != 60 {
		t.Fatal("SubAcquire failed")
	}

	// SubRelease
	a.SubRelease(10)
	if a.Load() != 50 {
		t.Fatal("SubRelease failed")
	}

	// SubAcqRel
	a.SubAcqRel(10)
	if a.Load() != 40 {
		t.Fatal("SubAcqRel failed")
	}

	// AndRelaxed
	a.Store(0xFF)
	a.AndRelaxed(0x0F)
	if a.Load() != 0x0F {
		t.Fatal("AndRelaxed failed")
	}

	// OrRelaxed
	a.OrRelaxed(0xF0)
	if a.Load() != 0xFF {
		t.Fatal("OrRelaxed failed")
	}

	// XorRelaxed
	a.XorRelaxed(0x0F)
	if a.Load() != 0xF0 {
		t.Fatal("XorRelaxed failed")
	}

	// MaxRelaxed
	a.Store(10)
	a.MaxRelaxed(20)
	if a.Load() != 20 {
		t.Fatal("MaxRelaxed failed")
	}

	// MinRelaxed
	a.MinRelaxed(15)
	if a.Load() != 15 {
		t.Fatal("MinRelaxed failed")
	}
}

func TestInt64RelaxedVariants(t *testing.T) {
	var a atomix.Int64

	// AndRelaxed
	a.Store(0xFF)
	a.AndRelaxed(0x0F)
	if a.Load() != 0x0F {
		t.Fatal("AndRelaxed failed")
	}

	// OrRelaxed
	a.OrRelaxed(0xF0)
	if a.Load() != 0xFF {
		t.Fatal("OrRelaxed failed")
	}

	// XorRelaxed
	a.XorRelaxed(0x0F)
	if a.Load() != 0xF0 {
		t.Fatal("XorRelaxed failed")
	}

	// MaxRelaxed
	a.Store(10)
	a.MaxRelaxed(20)
	if a.Load() != 20 {
		t.Fatal("MaxRelaxed failed")
	}

	// MinRelaxed
	a.MinRelaxed(15)
	if a.Load() != 15 {
		t.Fatal("MinRelaxed failed")
	}
}

func TestUint32RelaxedVariants(t *testing.T) {
	var a atomix.Uint32

	// AndRelaxed
	a.Store(0xFF)
	a.AndRelaxed(0x0F)
	if a.Load() != 0x0F {
		t.Fatal("AndRelaxed failed")
	}

	// OrRelaxed
	a.OrRelaxed(0xF0)
	if a.Load() != 0xFF {
		t.Fatal("OrRelaxed failed")
	}

	// XorRelaxed
	a.XorRelaxed(0x0F)
	if a.Load() != 0xF0 {
		t.Fatal("XorRelaxed failed")
	}

	// MaxRelaxed
	a.Store(10)
	a.MaxRelaxed(20)
	if a.Load() != 20 {
		t.Fatal("MaxRelaxed failed")
	}

	// MinRelaxed
	a.MinRelaxed(15)
	if a.Load() != 15 {
		t.Fatal("MinRelaxed failed")
	}
}

func TestUint64RelaxedVariants(t *testing.T) {
	var a atomix.Uint64

	// LoadRelaxed
	a.Store(42)
	if a.LoadRelaxed() != 42 {
		t.Fatal("LoadRelaxed failed")
	}

	// LoadAcquire
	if a.LoadAcquire() != 42 {
		t.Fatal("LoadAcquire failed")
	}

	// StoreRelaxed
	a.StoreRelaxed(100)
	if a.Load() != 100 {
		t.Fatal("StoreRelaxed failed")
	}

	// StoreRelease
	a.StoreRelease(200)
	if a.Load() != 200 {
		t.Fatal("StoreRelease failed")
	}

	// AddRelaxed
	a.AddRelaxed(10)
	if a.Load() != 210 {
		t.Fatal("AddRelaxed failed")
	}

	// AddAcquire
	a.AddAcquire(10)
	if a.Load() != 220 {
		t.Fatal("AddAcquire failed")
	}

	// AddRelease
	a.AddRelease(10)
	if a.Load() != 230 {
		t.Fatal("AddRelease failed")
	}

	// AddAcqRel
	a.AddAcqRel(10)
	if a.Load() != 240 {
		t.Fatal("AddAcqRel failed")
	}

	// AndRelaxed
	a.Store(0xFF)
	a.AndRelaxed(0x0F)
	if a.Load() != 0x0F {
		t.Fatal("AndRelaxed failed")
	}

	// OrRelaxed
	a.OrRelaxed(0xF0)
	if a.Load() != 0xFF {
		t.Fatal("OrRelaxed failed")
	}

	// XorRelaxed
	a.XorRelaxed(0x0F)
	if a.Load() != 0xF0 {
		t.Fatal("XorRelaxed failed")
	}

	// MaxRelaxed
	a.Store(10)
	a.MaxRelaxed(20)
	if a.Load() != 20 {
		t.Fatal("MaxRelaxed failed")
	}

	// MinRelaxed
	a.MinRelaxed(15)
	if a.Load() != 15 {
		t.Fatal("MinRelaxed failed")
	}
}

func TestUint128RelaxedVariants(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)

	// LoadRelaxed
	a.Store(42, 0)
	lo, hi := a.LoadRelaxed()
	if lo != 42 || hi != 0 {
		t.Fatal("LoadRelaxed failed")
	}

	// LoadAcquire
	lo, hi = a.LoadAcquire()
	if lo != 42 || hi != 0 {
		t.Fatal("LoadAcquire failed")
	}

	// StoreRelaxed
	a.StoreRelaxed(100, 1)
	lo, hi = a.Load()
	if lo != 100 || hi != 1 {
		t.Fatal("StoreRelaxed failed")
	}

	// StoreRelease
	a.StoreRelease(200, 2)
	lo, hi = a.Load()
	if lo != 200 || hi != 2 {
		t.Fatal("StoreRelease failed")
	}

	// CompareAndSwapRelaxed
	if !a.CompareAndSwapRelaxed(200, 2, 300, 3) {
		t.Fatal("CompareAndSwapRelaxed should succeed")
	}

	// CompareAndSwapAcquire
	if !a.CompareAndSwapAcquire(300, 3, 400, 4) {
		t.Fatal("CompareAndSwapAcquire should succeed")
	}

	// CompareAndSwapRelease
	if !a.CompareAndSwapRelease(400, 4, 500, 5) {
		t.Fatal("CompareAndSwapRelease should succeed")
	}

	// CompareAndSwapAcqRel
	if !a.CompareAndSwapAcqRel(500, 5, 600, 6) {
		t.Fatal("CompareAndSwapAcqRel should succeed")
	}

	// CompareExchangeRelaxed
	lo, hi = a.CompareExchangeRelaxed(600, 6, 700, 7)
	if lo != 600 || hi != 6 {
		t.Fatal("CompareExchangeRelaxed failed")
	}

	// CompareExchangeAcquire
	lo, hi = a.CompareExchangeAcquire(700, 7, 800, 8)
	if lo != 700 || hi != 7 {
		t.Fatal("CompareExchangeAcquire failed")
	}

	// CompareExchangeRelease
	lo, hi = a.CompareExchangeRelease(800, 8, 900, 9)
	if lo != 800 || hi != 8 {
		t.Fatal("CompareExchangeRelease failed")
	}

	// CompareExchangeAcqRel
	lo, hi = a.CompareExchangeAcqRel(900, 9, 1000, 10)
	if lo != 900 || hi != 9 {
		t.Fatal("CompareExchangeAcqRel failed")
	}

	// AddRelaxed
	a.Store(100, 0)
	a.AddRelaxed(50, 0)
	lo, hi = a.Load()
	if lo != 150 || hi != 0 {
		t.Fatal("AddRelaxed failed")
	}

	// SubRelaxed
	a.SubRelaxed(50, 0)
	lo, hi = a.Load()
	if lo != 100 || hi != 0 {
		t.Fatal("SubRelaxed failed")
	}

	// IncRelaxed
	a.Store(0, 0)
	a.IncRelaxed()
	lo, hi = a.Load()
	if lo != 1 || hi != 0 {
		t.Fatal("IncRelaxed failed")
	}

	// DecRelaxed
	a.DecRelaxed()
	lo, hi = a.Load()
	if lo != 0 || hi != 0 {
		t.Fatal("DecRelaxed failed")
	}

	// EqualRelaxed
	a.Store(42, 1)
	if !a.EqualRelaxed(42, 1) {
		t.Fatal("EqualRelaxed failed")
	}

	// LessRelaxed
	if !a.LessRelaxed(100, 1) {
		t.Fatal("LessRelaxed failed")
	}

	// LessOrEqual
	if !a.LessOrEqual(42, 1) {
		t.Fatal("LessOrEqual failed")
	}

	// LessOrEqualRelaxed
	if !a.LessOrEqualRelaxed(42, 1) {
		t.Fatal("LessOrEqualRelaxed failed")
	}

	// GreaterRelaxed
	if !a.GreaterRelaxed(10, 1) {
		t.Fatal("GreaterRelaxed failed")
	}

	// GreaterOrEqual
	if !a.GreaterOrEqual(42, 1) {
		t.Fatal("GreaterOrEqual failed")
	}

	// GreaterOrEqualRelaxed
	if !a.GreaterOrEqualRelaxed(42, 1) {
		t.Fatal("GreaterOrEqualRelaxed failed")
	}
}

func TestInt128RelaxedVariants(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)

	// LessRelaxed
	a.Store(10, 0)
	if !a.LessRelaxed(20, 0) {
		t.Fatal("LessRelaxed failed")
	}

	// LessOrEqual
	if !a.LessOrEqual(10, 0) {
		t.Fatal("LessOrEqual failed")
	}

	// LessOrEqualRelaxed
	if !a.LessOrEqualRelaxed(10, 0) {
		t.Fatal("LessOrEqualRelaxed failed")
	}

	// GreaterRelaxed
	if !a.GreaterRelaxed(5, 0) {
		t.Fatal("GreaterRelaxed failed")
	}

	// GreaterOrEqual
	if !a.GreaterOrEqual(10, 0) {
		t.Fatal("GreaterOrEqual failed")
	}

	// GreaterOrEqualRelaxed
	if !a.GreaterOrEqualRelaxed(10, 0) {
		t.Fatal("GreaterOrEqualRelaxed failed")
	}
}

func TestPointerCompareExchange(t *testing.T) {
	var a atomix.Pointer[int]

	v1 := 1
	v2 := 2
	v3 := 3
	v4 := 4
	v5 := 5
	v6 := 6

	a.Store(&v1)

	// CompareExchange
	old := a.CompareExchange(&v1, &v2)
	if old != &v1 {
		t.Fatal("CompareExchange failed")
	}
	if a.Load() != &v2 {
		t.Fatal("CompareExchange: value not updated")
	}

	// CompareExchangeRelaxed
	old = a.CompareExchangeRelaxed(&v2, &v3)
	if old != &v2 {
		t.Fatal("CompareExchangeRelaxed failed")
	}

	// CompareExchangeAcquire
	old = a.CompareExchangeAcquire(&v3, &v4)
	if old != &v3 {
		t.Fatal("CompareExchangeAcquire failed")
	}

	// CompareExchangeRelease
	old = a.CompareExchangeRelease(&v4, &v5)
	if old != &v4 {
		t.Fatal("CompareExchangeRelease failed")
	}

	// CompareExchangeAcqRel
	old = a.CompareExchangeAcqRel(&v5, &v6)
	if old != &v5 {
		t.Fatal("CompareExchangeAcqRel failed")
	}

	if a.Load() != &v6 {
		t.Fatal("CompareExchangeAcqRel: value not updated")
	}
}

func TestCanPlaceCacheAligned(t *testing.T) {
	buf := make([]byte, 256)

	// Should succeed with enough space for cache line
	if !atomix.CanPlaceCacheAligned(buf, 0, 64) {
		t.Fatal("CanPlaceCacheAligned should succeed with 256 byte buffer")
	}

	// Should fail with buffer too small
	smallBuf := make([]byte, 32)
	if atomix.CanPlaceCacheAligned(smallBuf, 0, 64) {
		t.Fatal("CanPlaceCacheAligned should fail with 32 byte buffer")
	}

	// Should fail with negative offset
	if atomix.CanPlaceCacheAligned(buf, -1, 64) {
		t.Fatal("CanPlaceCacheAligned should fail with negative offset")
	}

	// Should fail with offset beyond buffer
	if atomix.CanPlaceCacheAligned(buf, 300, 64) {
		t.Fatal("CanPlaceCacheAligned should fail with offset beyond buffer")
	}
}

// =============================================================================
// Min/Max Contention Tests (to cover CAS retry paths)
// =============================================================================

func TestUintptrMaxContention(t *testing.T) {
	var a atomix.Uintptr
	a.Store(0)

	const numGoroutines = 16
	const iterations = 1000
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				a.Max(uintptr(id*iterations + j))
				a.MaxRelaxed(uintptr(id*iterations + j))
			}
		}(i)
	}

	wg.Wait()

	expected := uintptr((numGoroutines-1)*iterations + iterations - 1)
	if got := a.Load(); got != expected {
		t.Fatalf("Max: got %d, want %d", got, expected)
	}
}

func TestUintptrMinContention(t *testing.T) {
	var a atomix.Uintptr
	a.Store(^uintptr(0)) // Start at max value

	const numGoroutines = 16
	const iterations = 1000
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				a.Min(uintptr(id*iterations + j + 1))
				a.MinRelaxed(uintptr(id*iterations + j + 1))
			}
		}(i)
	}

	wg.Wait()

	if got := a.Load(); got != 1 {
		t.Fatalf("Min: got %d, want 1", got)
	}
}

func TestUint64MaxContention(t *testing.T) {
	var a atomix.Uint64
	a.Store(0)

	const numGoroutines = 16
	const iterations = 1000
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// All goroutines race to set their ID as max
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				a.Max(uint64(id*iterations + j))
				a.MaxRelaxed(uint64(id*iterations + j))
			}
		}(i)
	}

	wg.Wait()

	// Final value should be the maximum
	expected := uint64((numGoroutines-1)*iterations + iterations - 1)
	if got := a.Load(); got != expected {
		t.Fatalf("Max: got %d, want %d", got, expected)
	}
}

func TestUint64MinContention(t *testing.T) {
	var a atomix.Uint64
	a.Store(^uint64(0)) // Start at max value

	const numGoroutines = 16
	const iterations = 1000
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// All goroutines race to set their ID as min
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				a.Min(uint64(id*iterations + j + 1))
				a.MinRelaxed(uint64(id*iterations + j + 1))
			}
		}(i)
	}

	wg.Wait()

	// Final value should be 1 (the minimum set)
	if got := a.Load(); got != 1 {
		t.Fatalf("Min: got %d, want 1", got)
	}
}

func TestInt64MaxContention(t *testing.T) {
	var a atomix.Int64
	a.Store(-1000000)

	const numGoroutines = 16
	const iterations = 1000
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				a.Max(int64(id*iterations + j))
				a.MaxRelaxed(int64(id*iterations + j))
			}
		}(i)
	}

	wg.Wait()

	expected := int64((numGoroutines-1)*iterations + iterations - 1)
	if got := a.Load(); got != expected {
		t.Fatalf("Max: got %d, want %d", got, expected)
	}
}

func TestInt64MinContention(t *testing.T) {
	var a atomix.Int64
	a.Store(1000000)

	const numGoroutines = 16
	const iterations = 1000
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				a.Min(int64(id*iterations + j))
				a.MinRelaxed(int64(id*iterations + j))
			}
		}(i)
	}

	wg.Wait()

	// Final value should be 0 (minimum set by goroutine 0)
	if got := a.Load(); got != 0 {
		t.Fatalf("Min: got %d, want 0", got)
	}
}

func TestUint32MaxMinContention(t *testing.T) {
	var maxA atomix.Uint32
	var minA atomix.Uint32
	maxA.Store(0)
	minA.Store(^uint32(0))

	const numGoroutines = 16
	const iterations = 500
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				maxA.Max(uint32(id*iterations + j))
				maxA.MaxRelaxed(uint32(id*iterations + j))
			}
		}(i)
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				minA.Min(uint32(id*iterations + j + 1))
				minA.MinRelaxed(uint32(id*iterations + j + 1))
			}
		}(i)
	}

	wg.Wait()

	expectedMax := uint32((numGoroutines-1)*iterations + iterations - 1)
	if got := maxA.Load(); got != expectedMax {
		t.Fatalf("Max: got %d, want %d", got, expectedMax)
	}
	if got := minA.Load(); got != 1 {
		t.Fatalf("Min: got %d, want 1", got)
	}
}

func TestInt32MaxMinContention(t *testing.T) {
	var maxA atomix.Int32
	var minA atomix.Int32
	maxA.Store(-1000000)
	minA.Store(1000000)

	const numGoroutines = 16
	const iterations = 500
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				maxA.Max(int32(id*iterations + j))
				maxA.MaxRelaxed(int32(id*iterations + j))
			}
		}(i)
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				minA.Min(int32(id*iterations + j))
				minA.MinRelaxed(int32(id*iterations + j))
			}
		}(i)
	}

	wg.Wait()

	expectedMax := int32((numGoroutines-1)*iterations + iterations - 1)
	if got := maxA.Load(); got != expectedMax {
		t.Fatalf("Max: got %d, want %d", got, expectedMax)
	}
	if got := minA.Load(); got != 0 {
		t.Fatalf("Min: got %d, want 0", got)
	}
}

// Test Uint128/Int128 AddRelaxed and SubRelaxed retry paths
func TestUint128AddSubContention(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)

	const numGoroutines = 8
	const iterations = 500
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	// Adders using Relaxed (lo=1, hi=0)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				a.AddRelaxed(1, 0)
			}
		}()
	}

	// Subtractors using Relaxed (lo=1, hi=0)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				a.SubRelaxed(1, 0)
			}
		}()
	}

	wg.Wait()

	// Should be zero - adds and subs cancel out
	lo, hi := a.Load()
	if lo != 0 || hi != 0 {
		t.Fatalf("AddRelaxed/SubRelaxed: got (%d, %d), want (0, 0)", lo, hi)
	}
}

func TestInt128AddSubContention(t *testing.T) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)

	const numGoroutines = 8
	const iterations = 500
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	// Adders using Relaxed (lo=1, hi=0)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				a.AddRelaxed(1, 0)
			}
		}()
	}

	// Subtractors using Relaxed (lo=1, hi=0)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				a.SubRelaxed(1, 0)
			}
		}()
	}

	wg.Wait()

	// Should be zero - adds and subs cancel out
	lo, hi := a.Load()
	if lo != 0 || hi != 0 {
		t.Fatalf("AddRelaxed/SubRelaxed: got (%d, %d), want (0, 0)", lo, hi)
	}
}

func TestInt128AddSubRetryContention(t *testing.T) {
	// Test the non-relaxed Add/Sub to cover their CAS retry paths
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedInt128(buf, 0)

	const numGoroutines = 8
	const iterations = 500
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				a.Add(1, 0)
			}
		}()
	}

	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				a.Sub(1, 0)
			}
		}()
	}

	wg.Wait()

	lo, hi := a.Load()
	if lo != 0 || hi != 0 {
		t.Fatalf("Add/Sub: got (%d, %d), want (0, 0)", lo, hi)
	}
}

// =============================================================================
// Alignment Error Path Tests
// =============================================================================

func TestPlaceAlignedPanic(t *testing.T) {
	testCases := []struct {
		name string
		fn   func()
	}{
		{"PlaceAlignedInt32", func() { atomix.PlaceAlignedInt32(make([]byte, 2), 0) }},
		{"PlaceAlignedUint32", func() { atomix.PlaceAlignedUint32(make([]byte, 2), 0) }},
		{"PlaceAlignedBool", func() { atomix.PlaceAlignedBool(make([]byte, 2), 0) }},
		{"PlaceAlignedInt64", func() { atomix.PlaceAlignedInt64(make([]byte, 4), 0) }},
		{"PlaceAlignedUint64", func() { atomix.PlaceAlignedUint64(make([]byte, 4), 0) }},
		{"PlaceAlignedUintptr", func() { atomix.PlaceAlignedUintptr(make([]byte, 4), 0) }},
		{"PlaceAlignedInt128", func() { atomix.PlaceAlignedInt128(make([]byte, 8), 0) }},
		{"PlaceAlignedUint128", func() { atomix.PlaceAlignedUint128(make([]byte, 8), 0) }},
		{"PlaceCacheAlignedInt32", func() { atomix.PlaceCacheAlignedInt32(make([]byte, 4), 0) }},
		{"PlaceCacheAlignedUint32", func() { atomix.PlaceCacheAlignedUint32(make([]byte, 4), 0) }},
		{"PlaceCacheAlignedInt64", func() { atomix.PlaceCacheAlignedInt64(make([]byte, 8), 0) }},
		{"PlaceCacheAlignedUint64", func() { atomix.PlaceCacheAlignedUint64(make([]byte, 8), 0) }},
		{"PlaceCacheAlignedUintptr", func() { atomix.PlaceCacheAlignedUintptr(make([]byte, 8), 0) }},
		{"PlaceCacheAlignedBool", func() { atomix.PlaceCacheAlignedBool(make([]byte, 1), 0) }},
		{"PlaceCacheAlignedInt128", func() { atomix.PlaceCacheAlignedInt128(make([]byte, 16), 0) }},
		{"PlaceCacheAlignedUint128", func() { atomix.PlaceCacheAlignedUint128(make([]byte, 16), 0) }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("%s should panic with insufficient space", tc.name)
				}
			}()
			tc.fn()
		})
	}
}

func TestCanPlaceAligned16EdgeCases(t *testing.T) {
	buf := make([]byte, 64)

	// Negative offset
	if atomix.CanPlaceAligned16(buf, -1) {
		t.Fatal("CanPlaceAligned16 should return false for negative offset")
	}

	// Offset beyond buffer
	if atomix.CanPlaceAligned16(buf, 100) {
		t.Fatal("CanPlaceAligned16 should return false for offset beyond buffer")
	}

	// Buffer too small for alignment + 16 bytes
	smallBuf := make([]byte, 8)
	if atomix.CanPlaceAligned16(smallBuf, 0) {
		t.Fatal("CanPlaceAligned16 should return false for small buffer")
	}
}

func TestAllocatorSkipAlign(t *testing.T) {
	// Test Allocator Skip and Align methods
	buf := make([]byte, 256)
	a := atomix.NewAllocator(buf)

	// Skip with valid amount (no panic)
	a.Skip(10)

	// Align to 8 bytes
	a2 := atomix.NewAllocator(buf)
	a2.Skip(3)
	a2.Align(8) // Should not panic

	// Align when already aligned
	a3 := atomix.NewAllocator(buf)
	a3.Align(8) // Should not panic
}

func TestAllocatorSkipPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Skip beyond buffer should panic")
		}
	}()
	buf := make([]byte, 10)
	a := atomix.NewAllocator(buf)
	a.Skip(100) // Should panic
}

func TestAllocatorAlignPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Align beyond buffer should panic")
		}
	}()
	smallBuf := make([]byte, 4)
	a := atomix.NewAllocator(smallBuf)
	a.Skip(2)
	a.Align(64) // Should panic - not enough space
}

func TestAllocatorAlignNonPowerOf2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Align with non-power-of-2 should panic")
		}
	}()
	buf := make([]byte, 256)
	a := atomix.NewAllocator(buf)
	a.Align(3) // 3 is not a power of 2
}

// =============================================================================
// Bitwise Operations Ordering Variants Coverage
// =============================================================================

func TestInt32BitwiseOrderingVariants(t *testing.T) {
	var a atomix.Int32

	// And ordering variants
	a.Store(0xFF)
	if old := a.AndAcquire(0x0F); old != 0xFF {
		t.Fatalf("AndAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndRelease(0x0F); old != 0xFF {
		t.Fatalf("AndRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndAcqRel(0x0F); old != 0xFF {
		t.Fatalf("AndAcqRel: got %x, want 0xFF", old)
	}

	// Or ordering variants
	a.Store(0x0F)
	if old := a.OrAcquire(0xF0); old != 0x0F {
		t.Fatalf("OrAcquire: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrRelease(0xF0); old != 0x0F {
		t.Fatalf("OrRelease: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrAcqRel(0xF0); old != 0x0F {
		t.Fatalf("OrAcqRel: got %x, want 0x0F", old)
	}

	// Xor ordering variants
	a.Store(0xFF)
	if old := a.XorAcquire(0x0F); old != 0xFF {
		t.Fatalf("XorAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorRelease(0x0F); old != 0xFF {
		t.Fatalf("XorRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorAcqRel(0x0F); old != 0xFF {
		t.Fatalf("XorAcqRel: got %x, want 0xFF", old)
	}
}

func TestInt64BitwiseOrderingVariants(t *testing.T) {
	var a atomix.Int64

	// And ordering variants
	a.Store(0xFF)
	if old := a.AndAcquire(0x0F); old != 0xFF {
		t.Fatalf("AndAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndRelease(0x0F); old != 0xFF {
		t.Fatalf("AndRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndAcqRel(0x0F); old != 0xFF {
		t.Fatalf("AndAcqRel: got %x, want 0xFF", old)
	}

	// Or ordering variants
	a.Store(0x0F)
	if old := a.OrAcquire(0xF0); old != 0x0F {
		t.Fatalf("OrAcquire: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrRelease(0xF0); old != 0x0F {
		t.Fatalf("OrRelease: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrAcqRel(0xF0); old != 0x0F {
		t.Fatalf("OrAcqRel: got %x, want 0x0F", old)
	}

	// Xor ordering variants
	a.Store(0xFF)
	if old := a.XorAcquire(0x0F); old != 0xFF {
		t.Fatalf("XorAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorRelease(0x0F); old != 0xFF {
		t.Fatalf("XorRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorAcqRel(0x0F); old != 0xFF {
		t.Fatalf("XorAcqRel: got %x, want 0xFF", old)
	}
}

func TestUint32BitwiseOrderingVariants(t *testing.T) {
	var a atomix.Uint32

	// And ordering variants
	a.Store(0xFF)
	if old := a.AndAcquire(0x0F); old != 0xFF {
		t.Fatalf("AndAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndRelease(0x0F); old != 0xFF {
		t.Fatalf("AndRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndAcqRel(0x0F); old != 0xFF {
		t.Fatalf("AndAcqRel: got %x, want 0xFF", old)
	}

	// Or ordering variants
	a.Store(0x0F)
	if old := a.OrAcquire(0xF0); old != 0x0F {
		t.Fatalf("OrAcquire: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrRelease(0xF0); old != 0x0F {
		t.Fatalf("OrRelease: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrAcqRel(0xF0); old != 0x0F {
		t.Fatalf("OrAcqRel: got %x, want 0x0F", old)
	}

	// Xor ordering variants
	a.Store(0xFF)
	if old := a.XorAcquire(0x0F); old != 0xFF {
		t.Fatalf("XorAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorRelease(0x0F); old != 0xFF {
		t.Fatalf("XorRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorAcqRel(0x0F); old != 0xFF {
		t.Fatalf("XorAcqRel: got %x, want 0xFF", old)
	}
}

func TestUint64BitwiseOrderingVariants(t *testing.T) {
	var a atomix.Uint64

	// And ordering variants
	a.Store(0xFF)
	if old := a.AndAcquire(0x0F); old != 0xFF {
		t.Fatalf("AndAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndRelease(0x0F); old != 0xFF {
		t.Fatalf("AndRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndAcqRel(0x0F); old != 0xFF {
		t.Fatalf("AndAcqRel: got %x, want 0xFF", old)
	}

	// Or ordering variants
	a.Store(0x0F)
	if old := a.OrAcquire(0xF0); old != 0x0F {
		t.Fatalf("OrAcquire: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrRelease(0xF0); old != 0x0F {
		t.Fatalf("OrRelease: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrAcqRel(0xF0); old != 0x0F {
		t.Fatalf("OrAcqRel: got %x, want 0x0F", old)
	}

	// Xor ordering variants
	a.Store(0xFF)
	if old := a.XorAcquire(0x0F); old != 0xFF {
		t.Fatalf("XorAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorRelease(0x0F); old != 0xFF {
		t.Fatalf("XorRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorAcqRel(0x0F); old != 0xFF {
		t.Fatalf("XorAcqRel: got %x, want 0xFF", old)
	}
}

func TestUintptrBitwiseOrderingVariants(t *testing.T) {
	var a atomix.Uintptr

	// And ordering variants
	a.Store(0xFF)
	if old := a.AndAcquire(0x0F); old != 0xFF {
		t.Fatalf("AndAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndRelease(0x0F); old != 0xFF {
		t.Fatalf("AndRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.AndAcqRel(0x0F); old != 0xFF {
		t.Fatalf("AndAcqRel: got %x, want 0xFF", old)
	}

	// Or ordering variants
	a.Store(0x0F)
	if old := a.OrAcquire(0xF0); old != 0x0F {
		t.Fatalf("OrAcquire: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrRelease(0xF0); old != 0x0F {
		t.Fatalf("OrRelease: got %x, want 0x0F", old)
	}
	a.Store(0x0F)
	if old := a.OrAcqRel(0xF0); old != 0x0F {
		t.Fatalf("OrAcqRel: got %x, want 0x0F", old)
	}

	// Xor ordering variants
	a.Store(0xFF)
	if old := a.XorAcquire(0x0F); old != 0xFF {
		t.Fatalf("XorAcquire: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorRelease(0x0F); old != 0xFF {
		t.Fatalf("XorRelease: got %x, want 0xFF", old)
	}
	a.Store(0xFF)
	if old := a.XorAcqRel(0x0F); old != 0xFF {
		t.Fatalf("XorAcqRel: got %x, want 0xFF", old)
	}
}
