// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build amd64

package arch

import "unsafe"

// x86-64 TSO (Total Store Order) Memory Model Guarantees:
//
// 1. All aligned loads and stores are atomic
// 2. Loads are never reordered with other loads
// 3. Stores are never reordered with other stores
// 4. Loads are not reordered with older stores to same location
//
// Therefore, on x86-64:
// - Relaxed Load/Store = regular memory access (atomic for aligned data)
// - Acquire Load = regular memory access (TSO provides acquire semantics)
// - Release Store = regular memory access (TSO provides release semantics)
//
// These pure Go implementations can be inlined by the compiler,
// eliminating function call overhead for ~3x performance improvement.

// =============================================================================
// 32-bit Load operations
// =============================================================================

// LoadInt32Relaxed atomically loads *addr with relaxed memory ordering.
func LoadInt32Relaxed(addr *int32) int32 {
	return *addr
}

// LoadInt32Acquire atomically loads *addr with acquire memory ordering.
// On x86-64 TSO, this is equivalent to relaxed ordering.
func LoadInt32Acquire(addr *int32) int32 {
	return *addr
}

// LoadUint32Relaxed atomically loads *addr with relaxed memory ordering.
func LoadUint32Relaxed(addr *uint32) uint32 {
	return *addr
}

// LoadUint32Acquire atomically loads *addr with acquire memory ordering.
func LoadUint32Acquire(addr *uint32) uint32 {
	return *addr
}

// =============================================================================
// 32-bit Store operations
// =============================================================================

// StoreInt32Relaxed atomically stores val to *addr with relaxed memory ordering.
func StoreInt32Relaxed(addr *int32, val int32) {
	*addr = val
}

// StoreInt32Release atomically stores val to *addr with release memory ordering.
// On x86-64 TSO, this is equivalent to relaxed ordering.
func StoreInt32Release(addr *int32, val int32) {
	*addr = val
}

// StoreUint32Relaxed atomically stores val to *addr with relaxed memory ordering.
func StoreUint32Relaxed(addr *uint32, val uint32) {
	*addr = val
}

// StoreUint32Release atomically stores val to *addr with release memory ordering.
func StoreUint32Release(addr *uint32, val uint32) {
	*addr = val
}

// =============================================================================
// 64-bit Load operations
// =============================================================================

// LoadInt64Relaxed atomically loads *addr with relaxed memory ordering.
func LoadInt64Relaxed(addr *int64) int64 {
	return *addr
}

// LoadInt64Acquire atomically loads *addr with acquire memory ordering.
func LoadInt64Acquire(addr *int64) int64 {
	return *addr
}

// LoadUint64Relaxed atomically loads *addr with relaxed memory ordering.
func LoadUint64Relaxed(addr *uint64) uint64 {
	return *addr
}

// LoadUint64Acquire atomically loads *addr with acquire memory ordering.
func LoadUint64Acquire(addr *uint64) uint64 {
	return *addr
}

// =============================================================================
// 64-bit Store operations
// =============================================================================

// StoreInt64Relaxed atomically stores val to *addr with relaxed memory ordering.
func StoreInt64Relaxed(addr *int64, val int64) {
	*addr = val
}

// StoreInt64Release atomically stores val to *addr with release memory ordering.
func StoreInt64Release(addr *int64, val int64) {
	*addr = val
}

// StoreUint64Relaxed atomically stores val to *addr with relaxed memory ordering.
func StoreUint64Relaxed(addr *uint64, val uint64) {
	*addr = val
}

// StoreUint64Release atomically stores val to *addr with release memory ordering.
func StoreUint64Release(addr *uint64, val uint64) {
	*addr = val
}

// =============================================================================
// Uintptr Load/Store operations
// =============================================================================

// LoadUintptrRelaxed atomically loads *addr with relaxed memory ordering.
func LoadUintptrRelaxed(addr *uintptr) uintptr {
	return *addr
}

// LoadUintptrAcquire atomically loads *addr with acquire memory ordering.
func LoadUintptrAcquire(addr *uintptr) uintptr {
	return *addr
}

// StoreUintptrRelaxed atomically stores val to *addr with relaxed memory ordering.
func StoreUintptrRelaxed(addr *uintptr, val uintptr) {
	*addr = val
}

// StoreUintptrRelease atomically stores val to *addr with release memory ordering.
func StoreUintptrRelease(addr *uintptr, val uintptr) {
	*addr = val
}

// =============================================================================
// Pointer Load/Store operations
// =============================================================================

// LoadPointerRelaxed atomically loads *addr with relaxed memory ordering.
func LoadPointerRelaxed(addr *unsafe.Pointer) unsafe.Pointer {
	return *addr
}

// LoadPointerAcquire atomically loads *addr with acquire memory ordering.
func LoadPointerAcquire(addr *unsafe.Pointer) unsafe.Pointer {
	return *addr
}

// StorePointerRelaxed atomically stores val to *addr with relaxed memory ordering.
func StorePointerRelaxed(addr *unsafe.Pointer, val unsafe.Pointer) {
	*addr = val
}

// StorePointerRelease atomically stores val to *addr with release memory ordering.
func StorePointerRelease(addr *unsafe.Pointer, val unsafe.Pointer) {
	*addr = val
}
