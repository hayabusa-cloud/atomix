// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// Load atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Uint32) Load() uint32 {
	return arch.LoadUint32Relaxed(&a.v)
}

// LoadRelaxed atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Uint32) LoadRelaxed() uint32 {
	return arch.LoadUint32Relaxed(&a.v)
}

// LoadAcquire atomically loads and returns the value with acquire ordering.
//
//go:nosplit
func (a *Uint32) LoadAcquire() uint32 {
	return arch.LoadUint32Acquire(&a.v)
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Uint32) Store(val uint32) {
	arch.StoreUint32Relaxed(&a.v, val)
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Uint32) StoreRelaxed(val uint32) {
	arch.StoreUint32Relaxed(&a.v, val)
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Uint32) StoreRelease(val uint32) {
	arch.StoreUint32Release(&a.v, val)
}

// Swap atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) Swap(new uint32) uint32 {
	return arch.SwapUint32AcqRel(&a.v, new)
}

// SwapRelaxed atomically swaps the value and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Uint32) SwapRelaxed(new uint32) uint32 {
	return arch.SwapUint32Relaxed(&a.v, new)
}

// SwapAcquire atomically swaps the value and returns the old value with acquire ordering.
//
//go:nosplit
func (a *Uint32) SwapAcquire(new uint32) uint32 {
	return arch.SwapUint32Acquire(&a.v, new)
}

// SwapRelease atomically swaps the value and returns the old value with release ordering.
//
//go:nosplit
func (a *Uint32) SwapRelease(new uint32) uint32 {
	return arch.SwapUint32Release(&a.v, new)
}

// SwapAcqRel atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) SwapAcqRel(new uint32) uint32 {
	return arch.SwapUint32AcqRel(&a.v, new)
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Uint32) CompareAndSwap(old, new uint32) bool {
	return arch.CasUint32AcqRel(&a.v, old, new)
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Uint32) CompareAndSwapRelaxed(old, new uint32) bool {
	return arch.CasUint32Relaxed(&a.v, old, new)
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Uint32) CompareAndSwapAcquire(old, new uint32) bool {
	return arch.CasUint32Acquire(&a.v, old, new)
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Uint32) CompareAndSwapRelease(old, new uint32) bool {
	return arch.CasUint32Release(&a.v, old, new)
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) CompareAndSwapAcqRel(old, new uint32) bool {
	return arch.CasUint32AcqRel(&a.v, old, new)
}

// CompareExchange atomically compares and swaps, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint32) CompareExchange(old, new uint32) uint32 {
	return arch.CaxUint32AcqRel(&a.v, old, new)
}

// CompareExchangeRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Uint32) CompareExchangeRelaxed(old, new uint32) uint32 {
	return arch.CaxUint32Relaxed(&a.v, old, new)
}

// CompareExchangeAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Uint32) CompareExchangeAcquire(old, new uint32) uint32 {
	return arch.CaxUint32Acquire(&a.v, old, new)
}

// CompareExchangeRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Uint32) CompareExchangeRelease(old, new uint32) uint32 {
	return arch.CaxUint32Release(&a.v, old, new)
}

// CompareExchangeAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) CompareExchangeAcqRel(old, new uint32) uint32 {
	return arch.CaxUint32AcqRel(&a.v, old, new)
}

// Add atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) Add(delta uint32) uint32 {
	return arch.AddUint32AcqRel(&a.v, delta)
}

// AddRelaxed atomically adds delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Uint32) AddRelaxed(delta uint32) uint32 {
	return arch.AddUint32Relaxed(&a.v, delta)
}

// AddAcquire atomically adds delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Uint32) AddAcquire(delta uint32) uint32 {
	return arch.AddUint32Acquire(&a.v, delta)
}

// AddRelease atomically adds delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Uint32) AddRelease(delta uint32) uint32 {
	return arch.AddUint32Release(&a.v, delta)
}

// AddAcqRel atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) AddAcqRel(delta uint32) uint32 {
	return arch.AddUint32AcqRel(&a.v, delta)
}

// Sub atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) Sub(delta uint32) uint32 {
	return arch.AddUint32AcqRel(&a.v, ^(delta - 1))
}

// SubRelaxed atomically subtracts delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Uint32) SubRelaxed(delta uint32) uint32 {
	return arch.AddUint32Relaxed(&a.v, ^(delta - 1))
}

// SubAcquire atomically subtracts delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Uint32) SubAcquire(delta uint32) uint32 {
	return arch.AddUint32Acquire(&a.v, ^(delta - 1))
}

// SubRelease atomically subtracts delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Uint32) SubRelease(delta uint32) uint32 {
	return arch.AddUint32Release(&a.v, ^(delta - 1))
}

// SubAcqRel atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) SubAcqRel(delta uint32) uint32 {
	return arch.AddUint32AcqRel(&a.v, ^(delta - 1))
}

// And atomically performs bitwise AND and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) And(mask uint32) uint32 {
	return arch.AndUint32AcqRel(&a.v, mask)
}

// AndRelaxed atomically performs bitwise AND with relaxed ordering.
//
//go:nosplit
func (a *Uint32) AndRelaxed(mask uint32) uint32 {
	return arch.AndUint32Relaxed(&a.v, mask)
}

// AndAcquire atomically performs bitwise AND with acquire ordering.
//
//go:nosplit
func (a *Uint32) AndAcquire(mask uint32) uint32 {
	return arch.AndUint32Acquire(&a.v, mask)
}

// AndRelease atomically performs bitwise AND with release ordering.
//
//go:nosplit
func (a *Uint32) AndRelease(mask uint32) uint32 {
	return arch.AndUint32Release(&a.v, mask)
}

// AndAcqRel atomically performs bitwise AND with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) AndAcqRel(mask uint32) uint32 {
	return arch.AndUint32AcqRel(&a.v, mask)
}

// Or atomically performs bitwise OR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) Or(mask uint32) uint32 {
	return arch.OrUint32AcqRel(&a.v, mask)
}

// OrRelaxed atomically performs bitwise OR with relaxed ordering.
//
//go:nosplit
func (a *Uint32) OrRelaxed(mask uint32) uint32 {
	return arch.OrUint32Relaxed(&a.v, mask)
}

// OrAcquire atomically performs bitwise OR with acquire ordering.
//
//go:nosplit
func (a *Uint32) OrAcquire(mask uint32) uint32 {
	return arch.OrUint32Acquire(&a.v, mask)
}

// OrRelease atomically performs bitwise OR with release ordering.
//
//go:nosplit
func (a *Uint32) OrRelease(mask uint32) uint32 {
	return arch.OrUint32Release(&a.v, mask)
}

// OrAcqRel atomically performs bitwise OR with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) OrAcqRel(mask uint32) uint32 {
	return arch.OrUint32AcqRel(&a.v, mask)
}

// Xor atomically performs bitwise XOR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) Xor(mask uint32) uint32 {
	return arch.XorUint32AcqRel(&a.v, mask)
}

// XorRelaxed atomically performs bitwise XOR with relaxed ordering.
//
//go:nosplit
func (a *Uint32) XorRelaxed(mask uint32) uint32 {
	return arch.XorUint32Relaxed(&a.v, mask)
}

// XorAcquire atomically performs bitwise XOR with acquire ordering.
//
//go:nosplit
func (a *Uint32) XorAcquire(mask uint32) uint32 {
	return arch.XorUint32Acquire(&a.v, mask)
}

// XorRelease atomically performs bitwise XOR with release ordering.
//
//go:nosplit
func (a *Uint32) XorRelease(mask uint32) uint32 {
	return arch.XorUint32Release(&a.v, mask)
}

// XorAcqRel atomically performs bitwise XOR with acquire-release ordering.
//
//go:nosplit
func (a *Uint32) XorAcqRel(mask uint32) uint32 {
	return arch.XorUint32AcqRel(&a.v, mask)
}

// Max atomically stores the maximum of current and val, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint32) Max(val uint32) uint32 {
	for {
		old := arch.LoadUint32Relaxed(&a.v)
		if old >= val {
			return old
		}
		if arch.CasUint32AcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MaxRelaxed atomically stores the maximum with relaxed ordering.
//
//go:nosplit
func (a *Uint32) MaxRelaxed(val uint32) uint32 {
	for {
		old := arch.LoadUint32Relaxed(&a.v)
		if old >= val {
			return old
		}
		if arch.CasUint32Relaxed(&a.v, old, val) {
			return old
		}
	}
}

// Min atomically stores the minimum of current and val, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint32) Min(val uint32) uint32 {
	for {
		old := arch.LoadUint32Relaxed(&a.v)
		if old <= val {
			return old
		}
		if arch.CasUint32AcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MinRelaxed atomically stores the minimum with relaxed ordering.
//
//go:nosplit
func (a *Uint32) MinRelaxed(val uint32) uint32 {
	for {
		old := arch.LoadUint32Relaxed(&a.v)
		if old <= val {
			return old
		}
		if arch.CasUint32Relaxed(&a.v, old, val) {
			return old
		}
	}
}
