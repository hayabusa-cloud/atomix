// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64

package arch

import "unsafe"

// ARM64 Memory Model:
//
// ARM64 is weakly ordered - loads and stores can be reordered freely.
// For ordering guarantees:
//   - Acquire: LDAR instruction (load-acquire)
//   - Release: STLR instruction (store-release)
//   - Relaxed: Plain LDR/STR (no ordering, just atomicity)
//
// Pure Go `return *addr` compiles to LDR (relaxed load).
// Pure Go `*addr = val` compiles to STR (relaxed store).
//
// These pure Go implementations can be inlined by the compiler,
// providing significant performance improvement over assembly calls.
//
// IMPORTANT: Acquire loads and Release stores MUST stay in assembly
// because Go cannot generate LDAR/STLR instructions from pure Go.

// =============================================================================
// 32-bit Relaxed Load operations (inlinable)
// =============================================================================

// LoadInt32Relaxed atomically loads *addr with relaxed memory ordering.
// Compiles to: LDR (no acquire barrier)
//
//go:nosplit
func LoadInt32Relaxed(addr *int32) int32 {
	return *addr
}

// LoadUint32Relaxed atomically loads *addr with relaxed memory ordering.
//
//go:nosplit
func LoadUint32Relaxed(addr *uint32) uint32 {
	return *addr
}

// =============================================================================
// 32-bit Relaxed Store operations (inlinable)
// =============================================================================

// StoreInt32Relaxed atomically stores val to *addr with relaxed memory ordering.
// Compiles to: STR (no release barrier)
//
//go:nosplit
func StoreInt32Relaxed(addr *int32, val int32) {
	*addr = val
}

// StoreUint32Relaxed atomically stores val to *addr with relaxed memory ordering.
//
//go:nosplit
func StoreUint32Relaxed(addr *uint32, val uint32) {
	*addr = val
}

// =============================================================================
// 64-bit Relaxed Load operations (inlinable)
// =============================================================================

// LoadInt64Relaxed atomically loads *addr with relaxed memory ordering.
//
//go:nosplit
func LoadInt64Relaxed(addr *int64) int64 {
	return *addr
}

// LoadUint64Relaxed atomically loads *addr with relaxed memory ordering.
//
//go:nosplit
func LoadUint64Relaxed(addr *uint64) uint64 {
	return *addr
}

// =============================================================================
// 64-bit Relaxed Store operations (inlinable)
// =============================================================================

// StoreInt64Relaxed atomically stores val to *addr with relaxed memory ordering.
//
//go:nosplit
func StoreInt64Relaxed(addr *int64, val int64) {
	*addr = val
}

// StoreUint64Relaxed atomically stores val to *addr with relaxed memory ordering.
//
//go:nosplit
func StoreUint64Relaxed(addr *uint64, val uint64) {
	*addr = val
}

// =============================================================================
// Uintptr Relaxed Load/Store operations (inlinable)
// =============================================================================

// LoadUintptrRelaxed atomically loads *addr with relaxed memory ordering.
//
//go:nosplit
func LoadUintptrRelaxed(addr *uintptr) uintptr {
	return *addr
}

// StoreUintptrRelaxed atomically stores val to *addr with relaxed memory ordering.
//
//go:nosplit
func StoreUintptrRelaxed(addr *uintptr, val uintptr) {
	*addr = val
}

// =============================================================================
// Pointer Relaxed Load/Store operations (inlinable)
// =============================================================================

// LoadPointerRelaxed atomically loads *addr with relaxed memory ordering.
//
//go:nosplit
func LoadPointerRelaxed(addr *unsafe.Pointer) unsafe.Pointer {
	return *addr
}

// StorePointerRelaxed atomically stores val to *addr with relaxed memory ordering.
//
//go:nosplit
func StorePointerRelaxed(addr *unsafe.Pointer, val unsafe.Pointer) {
	*addr = val
}
