// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package arch_test

import (
	"testing"
	"unsafe"

	"code.hybscloud.com/atomix/internal/arch"
)

// =============================================================================
// 32-bit Load/Store Tests
// =============================================================================

func TestLoadStoreInt32(t *testing.T) {
	var v int32 = 42

	// Relaxed
	if got := arch.LoadInt32Relaxed(&v); got != 42 {
		t.Fatalf("LoadInt32Relaxed: got %d, want 42", got)
	}
	arch.StoreInt32Relaxed(&v, 100)
	if v != 100 {
		t.Fatalf("StoreInt32Relaxed: got %d, want 100", v)
	}

	// Acquire/Release
	if got := arch.LoadInt32Acquire(&v); got != 100 {
		t.Fatalf("LoadInt32Acquire: got %d, want 100", got)
	}
	arch.StoreInt32Release(&v, 200)
	if v != 200 {
		t.Fatalf("StoreInt32Release: got %d, want 200", v)
	}
}

func TestLoadStoreUint32(t *testing.T) {
	var v uint32 = 42

	// Relaxed
	if got := arch.LoadUint32Relaxed(&v); got != 42 {
		t.Fatalf("LoadUint32Relaxed: got %d, want 42", got)
	}
	arch.StoreUint32Relaxed(&v, 100)
	if v != 100 {
		t.Fatalf("StoreUint32Relaxed: got %d, want 100", v)
	}

	// Acquire/Release
	if got := arch.LoadUint32Acquire(&v); got != 100 {
		t.Fatalf("LoadUint32Acquire: got %d, want 100", got)
	}
	arch.StoreUint32Release(&v, 200)
	if v != 200 {
		t.Fatalf("StoreUint32Release: got %d, want 200", v)
	}
}

// =============================================================================
// 64-bit Load/Store Tests
// =============================================================================

func TestLoadStoreInt64(t *testing.T) {
	var v int64 = 42

	// Relaxed
	if got := arch.LoadInt64Relaxed(&v); got != 42 {
		t.Fatalf("LoadInt64Relaxed: got %d, want 42", got)
	}
	arch.StoreInt64Relaxed(&v, 100)
	if v != 100 {
		t.Fatalf("StoreInt64Relaxed: got %d, want 100", v)
	}

	// Acquire/Release
	if got := arch.LoadInt64Acquire(&v); got != 100 {
		t.Fatalf("LoadInt64Acquire: got %d, want 100", got)
	}
	arch.StoreInt64Release(&v, 200)
	if v != 200 {
		t.Fatalf("StoreInt64Release: got %d, want 200", v)
	}
}

func TestLoadStoreUint64(t *testing.T) {
	var v uint64 = 42

	// Relaxed
	if got := arch.LoadUint64Relaxed(&v); got != 42 {
		t.Fatalf("LoadUint64Relaxed: got %d, want 42", got)
	}
	arch.StoreUint64Relaxed(&v, 100)
	if v != 100 {
		t.Fatalf("StoreUint64Relaxed: got %d, want 100", v)
	}

	// Acquire/Release
	if got := arch.LoadUint64Acquire(&v); got != 100 {
		t.Fatalf("LoadUint64Acquire: got %d, want 100", got)
	}
	arch.StoreUint64Release(&v, 200)
	if v != 200 {
		t.Fatalf("StoreUint64Release: got %d, want 200", v)
	}
}

// =============================================================================
// Uintptr Load/Store Tests
// =============================================================================

func TestLoadStoreUintptr(t *testing.T) {
	var v uintptr = 42

	// Relaxed
	if got := arch.LoadUintptrRelaxed(&v); got != 42 {
		t.Fatalf("LoadUintptrRelaxed: got %d, want 42", got)
	}
	arch.StoreUintptrRelaxed(&v, 100)
	if v != 100 {
		t.Fatalf("StoreUintptrRelaxed: got %d, want 100", v)
	}

	// Acquire/Release
	if got := arch.LoadUintptrAcquire(&v); got != 100 {
		t.Fatalf("LoadUintptrAcquire: got %d, want 100", got)
	}
	arch.StoreUintptrRelease(&v, 200)
	if v != 200 {
		t.Fatalf("StoreUintptrRelease: got %d, want 200", v)
	}
}

// =============================================================================
// Pointer Load/Store Tests
// =============================================================================

func TestLoadStorePointer(t *testing.T) {
	var x, y int = 1, 2
	var v unsafe.Pointer = unsafe.Pointer(&x)

	// Relaxed
	if got := arch.LoadPointerRelaxed(&v); got != unsafe.Pointer(&x) {
		t.Fatal("LoadPointerRelaxed: got wrong pointer")
	}
	arch.StorePointerRelaxed(&v, unsafe.Pointer(&y))
	if v != unsafe.Pointer(&y) {
		t.Fatal("StorePointerRelaxed: got wrong pointer")
	}

	// Acquire/Release
	if got := arch.LoadPointerAcquire(&v); got != unsafe.Pointer(&y) {
		t.Fatal("LoadPointerAcquire: got wrong pointer")
	}
	arch.StorePointerRelease(&v, unsafe.Pointer(&x))
	if v != unsafe.Pointer(&x) {
		t.Fatal("StorePointerRelease: got wrong pointer")
	}
}

// =============================================================================
// 32-bit Swap Tests
// =============================================================================

func TestSwapInt32(t *testing.T) {
	var v int32 = 10

	if old := arch.SwapInt32Relaxed(&v, 20); old != 10 || v != 20 {
		t.Fatalf("SwapInt32Relaxed: old=%d, v=%d", old, v)
	}
	if old := arch.SwapInt32Acquire(&v, 30); old != 20 || v != 30 {
		t.Fatalf("SwapInt32Acquire: old=%d, v=%d", old, v)
	}
	if old := arch.SwapInt32Release(&v, 40); old != 30 || v != 40 {
		t.Fatalf("SwapInt32Release: old=%d, v=%d", old, v)
	}
	if old := arch.SwapInt32AcqRel(&v, 50); old != 40 || v != 50 {
		t.Fatalf("SwapInt32AcqRel: old=%d, v=%d", old, v)
	}
}

func TestSwapUint32(t *testing.T) {
	var v uint32 = 10

	if old := arch.SwapUint32Relaxed(&v, 20); old != 10 || v != 20 {
		t.Fatalf("SwapUint32Relaxed: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUint32Acquire(&v, 30); old != 20 || v != 30 {
		t.Fatalf("SwapUint32Acquire: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUint32Release(&v, 40); old != 30 || v != 40 {
		t.Fatalf("SwapUint32Release: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUint32AcqRel(&v, 50); old != 40 || v != 50 {
		t.Fatalf("SwapUint32AcqRel: old=%d, v=%d", old, v)
	}
}

// =============================================================================
// 64-bit Swap Tests
// =============================================================================

func TestSwapInt64(t *testing.T) {
	var v int64 = 10

	if old := arch.SwapInt64Relaxed(&v, 20); old != 10 || v != 20 {
		t.Fatalf("SwapInt64Relaxed: old=%d, v=%d", old, v)
	}
	if old := arch.SwapInt64Acquire(&v, 30); old != 20 || v != 30 {
		t.Fatalf("SwapInt64Acquire: old=%d, v=%d", old, v)
	}
	if old := arch.SwapInt64Release(&v, 40); old != 30 || v != 40 {
		t.Fatalf("SwapInt64Release: old=%d, v=%d", old, v)
	}
	if old := arch.SwapInt64AcqRel(&v, 50); old != 40 || v != 50 {
		t.Fatalf("SwapInt64AcqRel: old=%d, v=%d", old, v)
	}
}

func TestSwapUint64(t *testing.T) {
	var v uint64 = 10

	if old := arch.SwapUint64Relaxed(&v, 20); old != 10 || v != 20 {
		t.Fatalf("SwapUint64Relaxed: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUint64Acquire(&v, 30); old != 20 || v != 30 {
		t.Fatalf("SwapUint64Acquire: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUint64Release(&v, 40); old != 30 || v != 40 {
		t.Fatalf("SwapUint64Release: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUint64AcqRel(&v, 50); old != 40 || v != 50 {
		t.Fatalf("SwapUint64AcqRel: old=%d, v=%d", old, v)
	}
}

// =============================================================================
// Uintptr Swap Tests
// =============================================================================

func TestSwapUintptr(t *testing.T) {
	var v uintptr = 10

	if old := arch.SwapUintptrRelaxed(&v, 20); old != 10 || v != 20 {
		t.Fatalf("SwapUintptrRelaxed: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUintptrAcquire(&v, 30); old != 20 || v != 30 {
		t.Fatalf("SwapUintptrAcquire: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUintptrRelease(&v, 40); old != 30 || v != 40 {
		t.Fatalf("SwapUintptrRelease: old=%d, v=%d", old, v)
	}
	if old := arch.SwapUintptrAcqRel(&v, 50); old != 40 || v != 50 {
		t.Fatalf("SwapUintptrAcqRel: old=%d, v=%d", old, v)
	}
}

// =============================================================================
// Pointer Swap Tests
// =============================================================================

func TestSwapPointer(t *testing.T) {
	var a, b, c, d, e int
	var v unsafe.Pointer = unsafe.Pointer(&a)

	if old := arch.SwapPointerRelaxed(&v, unsafe.Pointer(&b)); old != unsafe.Pointer(&a) {
		t.Fatal("SwapPointerRelaxed: wrong old")
	}
	if old := arch.SwapPointerAcquire(&v, unsafe.Pointer(&c)); old != unsafe.Pointer(&b) {
		t.Fatal("SwapPointerAcquire: wrong old")
	}
	if old := arch.SwapPointerRelease(&v, unsafe.Pointer(&d)); old != unsafe.Pointer(&c) {
		t.Fatal("SwapPointerRelease: wrong old")
	}
	if old := arch.SwapPointerAcqRel(&v, unsafe.Pointer(&e)); old != unsafe.Pointer(&d) {
		t.Fatal("SwapPointerAcqRel: wrong old")
	}
}

// =============================================================================
// 32-bit CAS Tests
// =============================================================================

func TestCasInt32(t *testing.T) {
	var v int32 = 10

	// Success cases
	if !arch.CasInt32Relaxed(&v, 10, 20) || v != 20 {
		t.Fatal("CasInt32Relaxed: should succeed")
	}
	if !arch.CasInt32Acquire(&v, 20, 30) || v != 30 {
		t.Fatal("CasInt32Acquire: should succeed")
	}
	if !arch.CasInt32Release(&v, 30, 40) || v != 40 {
		t.Fatal("CasInt32Release: should succeed")
	}
	if !arch.CasInt32AcqRel(&v, 40, 50) || v != 50 {
		t.Fatal("CasInt32AcqRel: should succeed")
	}

	// Failure case
	if arch.CasInt32Relaxed(&v, 999, 1000) || v != 50 {
		t.Fatal("CasInt32Relaxed: should fail")
	}
}

func TestCasUint32(t *testing.T) {
	var v uint32 = 10

	// Success cases
	if !arch.CasUint32Relaxed(&v, 10, 20) || v != 20 {
		t.Fatal("CasUint32Relaxed: should succeed")
	}
	if !arch.CasUint32Acquire(&v, 20, 30) || v != 30 {
		t.Fatal("CasUint32Acquire: should succeed")
	}
	if !arch.CasUint32Release(&v, 30, 40) || v != 40 {
		t.Fatal("CasUint32Release: should succeed")
	}
	if !arch.CasUint32AcqRel(&v, 40, 50) || v != 50 {
		t.Fatal("CasUint32AcqRel: should succeed")
	}

	// Failure case
	if arch.CasUint32Relaxed(&v, 999, 1000) || v != 50 {
		t.Fatal("CasUint32Relaxed: should fail")
	}
}

// =============================================================================
// 64-bit CAS Tests
// =============================================================================

func TestCasInt64(t *testing.T) {
	var v int64 = 10

	// Success cases
	if !arch.CasInt64Relaxed(&v, 10, 20) || v != 20 {
		t.Fatal("CasInt64Relaxed: should succeed")
	}
	if !arch.CasInt64Acquire(&v, 20, 30) || v != 30 {
		t.Fatal("CasInt64Acquire: should succeed")
	}
	if !arch.CasInt64Release(&v, 30, 40) || v != 40 {
		t.Fatal("CasInt64Release: should succeed")
	}
	if !arch.CasInt64AcqRel(&v, 40, 50) || v != 50 {
		t.Fatal("CasInt64AcqRel: should succeed")
	}

	// Failure case
	if arch.CasInt64Relaxed(&v, 999, 1000) || v != 50 {
		t.Fatal("CasInt64Relaxed: should fail")
	}
}

func TestCasUint64(t *testing.T) {
	var v uint64 = 10

	// Success cases
	if !arch.CasUint64Relaxed(&v, 10, 20) || v != 20 {
		t.Fatal("CasUint64Relaxed: should succeed")
	}
	if !arch.CasUint64Acquire(&v, 20, 30) || v != 30 {
		t.Fatal("CasUint64Acquire: should succeed")
	}
	if !arch.CasUint64Release(&v, 30, 40) || v != 40 {
		t.Fatal("CasUint64Release: should succeed")
	}
	if !arch.CasUint64AcqRel(&v, 40, 50) || v != 50 {
		t.Fatal("CasUint64AcqRel: should succeed")
	}

	// Failure case
	if arch.CasUint64Relaxed(&v, 999, 1000) || v != 50 {
		t.Fatal("CasUint64Relaxed: should fail")
	}
}

// =============================================================================
// Uintptr CAS Tests
// =============================================================================

func TestCasUintptr(t *testing.T) {
	var v uintptr = 10

	// Success cases
	if !arch.CasUintptrRelaxed(&v, 10, 20) || v != 20 {
		t.Fatal("CasUintptrRelaxed: should succeed")
	}
	if !arch.CasUintptrAcquire(&v, 20, 30) || v != 30 {
		t.Fatal("CasUintptrAcquire: should succeed")
	}
	if !arch.CasUintptrRelease(&v, 30, 40) || v != 40 {
		t.Fatal("CasUintptrRelease: should succeed")
	}
	if !arch.CasUintptrAcqRel(&v, 40, 50) || v != 50 {
		t.Fatal("CasUintptrAcqRel: should succeed")
	}

	// Failure case
	if arch.CasUintptrRelaxed(&v, 999, 1000) || v != 50 {
		t.Fatal("CasUintptrRelaxed: should fail")
	}
}

// =============================================================================
// Pointer CAS Tests
// =============================================================================

func TestCasPointer(t *testing.T) {
	var a, b, c, d, e, f int
	var v unsafe.Pointer = unsafe.Pointer(&a)

	// Success cases
	if !arch.CasPointerRelaxed(&v, unsafe.Pointer(&a), unsafe.Pointer(&b)) {
		t.Fatal("CasPointerRelaxed: should succeed")
	}
	if !arch.CasPointerAcquire(&v, unsafe.Pointer(&b), unsafe.Pointer(&c)) {
		t.Fatal("CasPointerAcquire: should succeed")
	}
	if !arch.CasPointerRelease(&v, unsafe.Pointer(&c), unsafe.Pointer(&d)) {
		t.Fatal("CasPointerRelease: should succeed")
	}
	if !arch.CasPointerAcqRel(&v, unsafe.Pointer(&d), unsafe.Pointer(&e)) {
		t.Fatal("CasPointerAcqRel: should succeed")
	}

	// Failure case
	if arch.CasPointerRelaxed(&v, unsafe.Pointer(&f), unsafe.Pointer(&a)) {
		t.Fatal("CasPointerRelaxed: should fail")
	}
}

// =============================================================================
// 32-bit CAX (Compare-And-Exchange returning old value) Tests
// =============================================================================

func TestCaxInt32(t *testing.T) {
	var v int32 = 10

	// Success: returns old value, updates
	if old := arch.CaxInt32Relaxed(&v, 10, 20); old != 10 || v != 20 {
		t.Fatalf("CaxInt32Relaxed success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxInt32Acquire(&v, 20, 30); old != 20 || v != 30 {
		t.Fatalf("CaxInt32Acquire success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxInt32Release(&v, 30, 40); old != 30 || v != 40 {
		t.Fatalf("CaxInt32Release success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxInt32AcqRel(&v, 40, 50); old != 40 || v != 50 {
		t.Fatalf("CaxInt32AcqRel success: old=%d, v=%d", old, v)
	}

	// Failure: returns current value, no update
	if old := arch.CaxInt32Relaxed(&v, 999, 1000); old != 50 || v != 50 {
		t.Fatalf("CaxInt32Relaxed failure: old=%d, v=%d", old, v)
	}
}

func TestCaxUint32(t *testing.T) {
	var v uint32 = 10

	// Success
	if old := arch.CaxUint32Relaxed(&v, 10, 20); old != 10 || v != 20 {
		t.Fatalf("CaxUint32Relaxed success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUint32Acquire(&v, 20, 30); old != 20 || v != 30 {
		t.Fatalf("CaxUint32Acquire success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUint32Release(&v, 30, 40); old != 30 || v != 40 {
		t.Fatalf("CaxUint32Release success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUint32AcqRel(&v, 40, 50); old != 40 || v != 50 {
		t.Fatalf("CaxUint32AcqRel success: old=%d, v=%d", old, v)
	}

	// Failure
	if old := arch.CaxUint32Relaxed(&v, 999, 1000); old != 50 || v != 50 {
		t.Fatalf("CaxUint32Relaxed failure: old=%d, v=%d", old, v)
	}
}

// =============================================================================
// 64-bit CAX Tests
// =============================================================================

func TestCaxInt64(t *testing.T) {
	var v int64 = 10

	// Success
	if old := arch.CaxInt64Relaxed(&v, 10, 20); old != 10 || v != 20 {
		t.Fatalf("CaxInt64Relaxed success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxInt64Acquire(&v, 20, 30); old != 20 || v != 30 {
		t.Fatalf("CaxInt64Acquire success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxInt64Release(&v, 30, 40); old != 30 || v != 40 {
		t.Fatalf("CaxInt64Release success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxInt64AcqRel(&v, 40, 50); old != 40 || v != 50 {
		t.Fatalf("CaxInt64AcqRel success: old=%d, v=%d", old, v)
	}

	// Failure
	if old := arch.CaxInt64Relaxed(&v, 999, 1000); old != 50 || v != 50 {
		t.Fatalf("CaxInt64Relaxed failure: old=%d, v=%d", old, v)
	}
}

func TestCaxUint64(t *testing.T) {
	var v uint64 = 10

	// Success
	if old := arch.CaxUint64Relaxed(&v, 10, 20); old != 10 || v != 20 {
		t.Fatalf("CaxUint64Relaxed success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUint64Acquire(&v, 20, 30); old != 20 || v != 30 {
		t.Fatalf("CaxUint64Acquire success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUint64Release(&v, 30, 40); old != 30 || v != 40 {
		t.Fatalf("CaxUint64Release success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUint64AcqRel(&v, 40, 50); old != 40 || v != 50 {
		t.Fatalf("CaxUint64AcqRel success: old=%d, v=%d", old, v)
	}

	// Failure
	if old := arch.CaxUint64Relaxed(&v, 999, 1000); old != 50 || v != 50 {
		t.Fatalf("CaxUint64Relaxed failure: old=%d, v=%d", old, v)
	}
}

// =============================================================================
// Uintptr CAX Tests
// =============================================================================

func TestCaxUintptr(t *testing.T) {
	var v uintptr = 10

	// Success
	if old := arch.CaxUintptrRelaxed(&v, 10, 20); old != 10 || v != 20 {
		t.Fatalf("CaxUintptrRelaxed success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUintptrAcquire(&v, 20, 30); old != 20 || v != 30 {
		t.Fatalf("CaxUintptrAcquire success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUintptrRelease(&v, 30, 40); old != 30 || v != 40 {
		t.Fatalf("CaxUintptrRelease success: old=%d, v=%d", old, v)
	}
	if old := arch.CaxUintptrAcqRel(&v, 40, 50); old != 40 || v != 50 {
		t.Fatalf("CaxUintptrAcqRel success: old=%d, v=%d", old, v)
	}

	// Failure
	if old := arch.CaxUintptrRelaxed(&v, 999, 1000); old != 50 || v != 50 {
		t.Fatalf("CaxUintptrRelaxed failure: old=%d, v=%d", old, v)
	}
}

// =============================================================================
// Pointer CAX Tests
// =============================================================================

func TestCaxPointer(t *testing.T) {
	var a, b, c, d, e, f int
	var v unsafe.Pointer = unsafe.Pointer(&a)

	// Success
	if old := arch.CaxPointerRelaxed(&v, unsafe.Pointer(&a), unsafe.Pointer(&b)); old != unsafe.Pointer(&a) {
		t.Fatal("CaxPointerRelaxed success: wrong old")
	}
	if old := arch.CaxPointerAcquire(&v, unsafe.Pointer(&b), unsafe.Pointer(&c)); old != unsafe.Pointer(&b) {
		t.Fatal("CaxPointerAcquire success: wrong old")
	}
	if old := arch.CaxPointerRelease(&v, unsafe.Pointer(&c), unsafe.Pointer(&d)); old != unsafe.Pointer(&c) {
		t.Fatal("CaxPointerRelease success: wrong old")
	}
	if old := arch.CaxPointerAcqRel(&v, unsafe.Pointer(&d), unsafe.Pointer(&e)); old != unsafe.Pointer(&d) {
		t.Fatal("CaxPointerAcqRel success: wrong old")
	}

	// Failure: returns current value
	if old := arch.CaxPointerRelaxed(&v, unsafe.Pointer(&f), unsafe.Pointer(&a)); old != unsafe.Pointer(&e) {
		t.Fatal("CaxPointerRelaxed failure: should return current")
	}
}

// =============================================================================
// 32-bit Add Tests
// =============================================================================

func TestAddInt32(t *testing.T) {
	var v int32 = 10

	// Add returns NEW value (like sync/atomic)
	if got := arch.AddInt32Relaxed(&v, 5); got != 15 || v != 15 {
		t.Fatalf("AddInt32Relaxed: got=%d, v=%d", got, v)
	}
	if got := arch.AddInt32Acquire(&v, 5); got != 20 || v != 20 {
		t.Fatalf("AddInt32Acquire: got=%d, v=%d", got, v)
	}
	if got := arch.AddInt32Release(&v, 5); got != 25 || v != 25 {
		t.Fatalf("AddInt32Release: got=%d, v=%d", got, v)
	}
	if got := arch.AddInt32AcqRel(&v, 5); got != 30 || v != 30 {
		t.Fatalf("AddInt32AcqRel: got=%d, v=%d", got, v)
	}

	// Negative delta
	if got := arch.AddInt32Relaxed(&v, -10); got != 20 || v != 20 {
		t.Fatalf("AddInt32Relaxed negative: got=%d, v=%d", got, v)
	}
}

func TestAddUint32(t *testing.T) {
	var v uint32 = 10

	// Add returns NEW value (like sync/atomic)
	if got := arch.AddUint32Relaxed(&v, 5); got != 15 || v != 15 {
		t.Fatalf("AddUint32Relaxed: got=%d, v=%d", got, v)
	}
	if got := arch.AddUint32Acquire(&v, 5); got != 20 || v != 20 {
		t.Fatalf("AddUint32Acquire: got=%d, v=%d", got, v)
	}
	if got := arch.AddUint32Release(&v, 5); got != 25 || v != 25 {
		t.Fatalf("AddUint32Release: got=%d, v=%d", got, v)
	}
	if got := arch.AddUint32AcqRel(&v, 5); got != 30 || v != 30 {
		t.Fatalf("AddUint32AcqRel: got=%d, v=%d", got, v)
	}
}

// =============================================================================
// 64-bit Add Tests
// =============================================================================

func TestAddInt64(t *testing.T) {
	var v int64 = 10

	// Add returns NEW value (like sync/atomic)
	if got := arch.AddInt64Relaxed(&v, 5); got != 15 || v != 15 {
		t.Fatalf("AddInt64Relaxed: got=%d, v=%d", got, v)
	}
	if got := arch.AddInt64Acquire(&v, 5); got != 20 || v != 20 {
		t.Fatalf("AddInt64Acquire: got=%d, v=%d", got, v)
	}
	if got := arch.AddInt64Release(&v, 5); got != 25 || v != 25 {
		t.Fatalf("AddInt64Release: got=%d, v=%d", got, v)
	}
	if got := arch.AddInt64AcqRel(&v, 5); got != 30 || v != 30 {
		t.Fatalf("AddInt64AcqRel: got=%d, v=%d", got, v)
	}

	// Negative delta
	if got := arch.AddInt64Relaxed(&v, -10); got != 20 || v != 20 {
		t.Fatalf("AddInt64Relaxed negative: got=%d, v=%d", got, v)
	}
}

func TestAddUint64(t *testing.T) {
	var v uint64 = 10

	// Add returns NEW value (like sync/atomic)
	if got := arch.AddUint64Relaxed(&v, 5); got != 15 || v != 15 {
		t.Fatalf("AddUint64Relaxed: got=%d, v=%d", got, v)
	}
	if got := arch.AddUint64Acquire(&v, 5); got != 20 || v != 20 {
		t.Fatalf("AddUint64Acquire: got=%d, v=%d", got, v)
	}
	if got := arch.AddUint64Release(&v, 5); got != 25 || v != 25 {
		t.Fatalf("AddUint64Release: got=%d, v=%d", got, v)
	}
	if got := arch.AddUint64AcqRel(&v, 5); got != 30 || v != 30 {
		t.Fatalf("AddUint64AcqRel: got=%d, v=%d", got, v)
	}
}

// =============================================================================
// Uintptr Add Tests
// =============================================================================

func TestAddUintptr(t *testing.T) {
	var v uintptr = 10

	// Add returns NEW value (like sync/atomic)
	if got := arch.AddUintptrRelaxed(&v, 5); got != 15 || v != 15 {
		t.Fatalf("AddUintptrRelaxed: got=%d, v=%d", got, v)
	}
	if got := arch.AddUintptrAcquire(&v, 5); got != 20 || v != 20 {
		t.Fatalf("AddUintptrAcquire: got=%d, v=%d", got, v)
	}
	if got := arch.AddUintptrRelease(&v, 5); got != 25 || v != 25 {
		t.Fatalf("AddUintptrRelease: got=%d, v=%d", got, v)
	}
	if got := arch.AddUintptrAcqRel(&v, 5); got != 30 || v != 30 {
		t.Fatalf("AddUintptrAcqRel: got=%d, v=%d", got, v)
	}
}

// =============================================================================
// 128-bit Load/Store Tests
// =============================================================================

// newAligned16 allocates a 16-byte aligned [16]byte on the heap.
// Uses the same pattern as atomix.NewInt128 for alignment.
//
//go:nocheckptr
func newAligned16() *[16]byte {
	// Allocate 3 uint64s (24 bytes); at least one 16-byte window will be aligned
	mem := make([]uint64, 3)
	base := uintptr(unsafe.Pointer(unsafe.SliceData(mem)))
	pad := (16 - (base & 15)) & 15 // 0 or 8
	off := pad >> 3                // 0 or 1 (words)
	return (*[16]byte)(unsafe.Pointer(&mem[off]))
}

func TestLoadStoreUint128(t *testing.T) {
	v := newAligned16()

	*(*uint64)(unsafe.Pointer(&v[0])) = 0x1234567890ABCDEF
	*(*uint64)(unsafe.Pointer(&v[8])) = 0xFEDCBA0987654321

	// Load
	lo, hi := arch.LoadUint128Relaxed(v)
	if lo != 0x1234567890ABCDEF || hi != 0xFEDCBA0987654321 {
		t.Fatalf("LoadUint128Relaxed: lo=%x, hi=%x", lo, hi)
	}

	lo, hi = arch.LoadUint128Acquire(v)
	if lo != 0x1234567890ABCDEF || hi != 0xFEDCBA0987654321 {
		t.Fatalf("LoadUint128Acquire: lo=%x, hi=%x", lo, hi)
	}

	// Store
	arch.StoreUint128Relaxed(v, 0xAAAAAAAAAAAAAAAA, 0xBBBBBBBBBBBBBBBB)
	lo = *(*uint64)(unsafe.Pointer(&v[0]))
	hi = *(*uint64)(unsafe.Pointer(&v[8]))
	if lo != 0xAAAAAAAAAAAAAAAA || hi != 0xBBBBBBBBBBBBBBBB {
		t.Fatalf("StoreUint128Relaxed: lo=%x, hi=%x", lo, hi)
	}

	arch.StoreUint128Release(v, 0xCCCCCCCCCCCCCCCC, 0xDDDDDDDDDDDDDDDD)
	lo = *(*uint64)(unsafe.Pointer(&v[0]))
	hi = *(*uint64)(unsafe.Pointer(&v[8]))
	if lo != 0xCCCCCCCCCCCCCCCC || hi != 0xDDDDDDDDDDDDDDDD {
		t.Fatalf("StoreUint128Release: lo=%x, hi=%x", lo, hi)
	}
}

// =============================================================================
// 128-bit Swap Tests
// =============================================================================

func TestSwapUint128(t *testing.T) {
	v := newAligned16()
	*(*uint64)(unsafe.Pointer(&v[0])) = 0x1111111111111111
	*(*uint64)(unsafe.Pointer(&v[8])) = 0x2222222222222222

	oldLo, oldHi := arch.SwapUint128Relaxed(v, 0x3333333333333333, 0x4444444444444444)
	if oldLo != 0x1111111111111111 || oldHi != 0x2222222222222222 {
		t.Fatalf("SwapUint128Relaxed: oldLo=%x, oldHi=%x", oldLo, oldHi)
	}

	oldLo, oldHi = arch.SwapUint128Acquire(v, 0x5555555555555555, 0x6666666666666666)
	if oldLo != 0x3333333333333333 || oldHi != 0x4444444444444444 {
		t.Fatalf("SwapUint128Acquire: oldLo=%x, oldHi=%x", oldLo, oldHi)
	}

	oldLo, oldHi = arch.SwapUint128Release(v, 0x7777777777777777, 0x8888888888888888)
	if oldLo != 0x5555555555555555 || oldHi != 0x6666666666666666 {
		t.Fatalf("SwapUint128Release: oldLo=%x, oldHi=%x", oldLo, oldHi)
	}

	oldLo, oldHi = arch.SwapUint128AcqRel(v, 0x9999999999999999, 0xAAAAAAAAAAAAAAAA)
	if oldLo != 0x7777777777777777 || oldHi != 0x8888888888888888 {
		t.Fatalf("SwapUint128AcqRel: oldLo=%x, oldHi=%x", oldLo, oldHi)
	}
}

// =============================================================================
// 128-bit CAS Tests
// =============================================================================

func TestCasUint128(t *testing.T) {
	v := newAligned16()
	*(*uint64)(unsafe.Pointer(&v[0])) = 0x1111111111111111
	*(*uint64)(unsafe.Pointer(&v[8])) = 0x2222222222222222

	// Success
	if !arch.CasUint128Relaxed(v, 0x1111111111111111, 0x2222222222222222, 0x3333333333333333, 0x4444444444444444) {
		t.Fatal("CasUint128Relaxed: should succeed")
	}
	if !arch.CasUint128Acquire(v, 0x3333333333333333, 0x4444444444444444, 0x5555555555555555, 0x6666666666666666) {
		t.Fatal("CasUint128Acquire: should succeed")
	}
	if !arch.CasUint128Release(v, 0x5555555555555555, 0x6666666666666666, 0x7777777777777777, 0x8888888888888888) {
		t.Fatal("CasUint128Release: should succeed")
	}
	if !arch.CasUint128AcqRel(v, 0x7777777777777777, 0x8888888888888888, 0x9999999999999999, 0xAAAAAAAAAAAAAAAA) {
		t.Fatal("CasUint128AcqRel: should succeed")
	}

	// Failure
	if arch.CasUint128Relaxed(v, 0x0, 0x0, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF) {
		t.Fatal("CasUint128Relaxed: should fail")
	}
}

// =============================================================================
// 128-bit CAX Tests
// =============================================================================

func TestCaxUint128(t *testing.T) {
	v := newAligned16()
	*(*uint64)(unsafe.Pointer(&v[0])) = 0x1111111111111111
	*(*uint64)(unsafe.Pointer(&v[8])) = 0x2222222222222222

	// Success
	lo, hi := arch.CaxUint128Relaxed(v, 0x1111111111111111, 0x2222222222222222, 0x3333333333333333, 0x4444444444444444)
	if lo != 0x1111111111111111 || hi != 0x2222222222222222 {
		t.Fatalf("CaxUint128Relaxed success: lo=%x, hi=%x", lo, hi)
	}

	lo, hi = arch.CaxUint128Acquire(v, 0x3333333333333333, 0x4444444444444444, 0x5555555555555555, 0x6666666666666666)
	if lo != 0x3333333333333333 || hi != 0x4444444444444444 {
		t.Fatalf("CaxUint128Acquire success: lo=%x, hi=%x", lo, hi)
	}

	lo, hi = arch.CaxUint128Release(v, 0x5555555555555555, 0x6666666666666666, 0x7777777777777777, 0x8888888888888888)
	if lo != 0x5555555555555555 || hi != 0x6666666666666666 {
		t.Fatalf("CaxUint128Release success: lo=%x, hi=%x", lo, hi)
	}

	lo, hi = arch.CaxUint128AcqRel(v, 0x7777777777777777, 0x8888888888888888, 0x9999999999999999, 0xAAAAAAAAAAAAAAAA)
	if lo != 0x7777777777777777 || hi != 0x8888888888888888 {
		t.Fatalf("CaxUint128AcqRel success: lo=%x, hi=%x", lo, hi)
	}

	// Failure: returns current value
	lo, hi = arch.CaxUint128Relaxed(v, 0x0, 0x0, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF)
	if lo != 0x9999999999999999 || hi != 0xAAAAAAAAAAAAAAAA {
		t.Fatalf("CaxUint128Relaxed failure: lo=%x, hi=%x", lo, hi)
	}
}

// =============================================================================
// Barrier Tests
// =============================================================================

func TestBarriers(t *testing.T) {
	// Barriers should not panic
	arch.BarrierAcquire()
	arch.BarrierRelease()
	arch.BarrierAcqRel()
}
