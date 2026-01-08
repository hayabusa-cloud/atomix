// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// Load atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Bool) Load() bool {
	return arch.LoadUint32Relaxed(&a.v) != 0
}

// LoadRelaxed atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Bool) LoadRelaxed() bool {
	return arch.LoadUint32Relaxed(&a.v) != 0
}

// LoadAcquire atomically loads and returns the value with acquire ordering.
//
//go:nosplit
func (a *Bool) LoadAcquire() bool {
	return arch.LoadUint32Acquire(&a.v) != 0
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Bool) Store(val bool) {
	arch.StoreUint32Relaxed(&a.v, b2u(val))
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Bool) StoreRelaxed(val bool) {
	arch.StoreUint32Relaxed(&a.v, b2u(val))
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Bool) StoreRelease(val bool) {
	arch.StoreUint32Release(&a.v, b2u(val))
}

// Swap atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Bool) Swap(new bool) bool {
	return arch.SwapUint32AcqRel(&a.v, b2u(new)) != 0
}

// SwapRelaxed atomically swaps the value and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Bool) SwapRelaxed(new bool) bool {
	return arch.SwapUint32Relaxed(&a.v, b2u(new)) != 0
}

// SwapAcquire atomically swaps the value and returns the old value with acquire ordering.
//
//go:nosplit
func (a *Bool) SwapAcquire(new bool) bool {
	return arch.SwapUint32Acquire(&a.v, b2u(new)) != 0
}

// SwapRelease atomically swaps the value and returns the old value with release ordering.
//
//go:nosplit
func (a *Bool) SwapRelease(new bool) bool {
	return arch.SwapUint32Release(&a.v, b2u(new)) != 0
}

// SwapAcqRel atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Bool) SwapAcqRel(new bool) bool {
	return arch.SwapUint32AcqRel(&a.v, b2u(new)) != 0
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Bool) CompareAndSwap(old, new bool) bool {
	return arch.CasUint32AcqRel(&a.v, b2u(old), b2u(new))
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Bool) CompareAndSwapRelaxed(old, new bool) bool {
	return arch.CasUint32Relaxed(&a.v, b2u(old), b2u(new))
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Bool) CompareAndSwapAcquire(old, new bool) bool {
	return arch.CasUint32Acquire(&a.v, b2u(old), b2u(new))
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Bool) CompareAndSwapRelease(old, new bool) bool {
	return arch.CasUint32Release(&a.v, b2u(old), b2u(new))
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Bool) CompareAndSwapAcqRel(old, new bool) bool {
	return arch.CasUint32AcqRel(&a.v, b2u(old), b2u(new))
}

// b2u converts a bool to uint32.
//
//go:nosplit
func b2u(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
