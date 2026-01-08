// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// Load atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Uint64) Load() uint64 {
	return arch.LoadUint64Relaxed(&a.v)
}

// LoadRelaxed atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Uint64) LoadRelaxed() uint64 {
	return arch.LoadUint64Relaxed(&a.v)
}

// LoadAcquire atomically loads and returns the value with acquire ordering.
//
//go:nosplit
func (a *Uint64) LoadAcquire() uint64 {
	return arch.LoadUint64Acquire(&a.v)
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Uint64) Store(val uint64) {
	arch.StoreUint64Relaxed(&a.v, val)
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Uint64) StoreRelaxed(val uint64) {
	arch.StoreUint64Relaxed(&a.v, val)
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Uint64) StoreRelease(val uint64) {
	arch.StoreUint64Release(&a.v, val)
}

// Swap atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) Swap(new uint64) uint64 {
	return arch.SwapUint64AcqRel(&a.v, new)
}

// SwapRelaxed atomically swaps the value and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Uint64) SwapRelaxed(new uint64) uint64 {
	return arch.SwapUint64Relaxed(&a.v, new)
}

// SwapAcquire atomically swaps the value and returns the old value with acquire ordering.
//
//go:nosplit
func (a *Uint64) SwapAcquire(new uint64) uint64 {
	return arch.SwapUint64Acquire(&a.v, new)
}

// SwapRelease atomically swaps the value and returns the old value with release ordering.
//
//go:nosplit
func (a *Uint64) SwapRelease(new uint64) uint64 {
	return arch.SwapUint64Release(&a.v, new)
}

// SwapAcqRel atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) SwapAcqRel(new uint64) uint64 {
	return arch.SwapUint64AcqRel(&a.v, new)
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Uint64) CompareAndSwap(old, new uint64) bool {
	return arch.CasUint64AcqRel(&a.v, old, new)
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Uint64) CompareAndSwapRelaxed(old, new uint64) bool {
	return arch.CasUint64Relaxed(&a.v, old, new)
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Uint64) CompareAndSwapAcquire(old, new uint64) bool {
	return arch.CasUint64Acquire(&a.v, old, new)
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Uint64) CompareAndSwapRelease(old, new uint64) bool {
	return arch.CasUint64Release(&a.v, old, new)
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) CompareAndSwapAcqRel(old, new uint64) bool {
	return arch.CasUint64AcqRel(&a.v, old, new)
}

// CompareExchange atomically compares and swaps, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint64) CompareExchange(old, new uint64) uint64 {
	return arch.CaxUint64AcqRel(&a.v, old, new)
}

// CompareExchangeRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Uint64) CompareExchangeRelaxed(old, new uint64) uint64 {
	return arch.CaxUint64Relaxed(&a.v, old, new)
}

// CompareExchangeAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Uint64) CompareExchangeAcquire(old, new uint64) uint64 {
	return arch.CaxUint64Acquire(&a.v, old, new)
}

// CompareExchangeRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Uint64) CompareExchangeRelease(old, new uint64) uint64 {
	return arch.CaxUint64Release(&a.v, old, new)
}

// CompareExchangeAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) CompareExchangeAcqRel(old, new uint64) uint64 {
	return arch.CaxUint64AcqRel(&a.v, old, new)
}

// Add atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) Add(delta uint64) uint64 {
	return arch.AddUint64AcqRel(&a.v, delta)
}

// AddRelaxed atomically adds delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Uint64) AddRelaxed(delta uint64) uint64 {
	return arch.AddUint64Relaxed(&a.v, delta)
}

// AddAcquire atomically adds delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Uint64) AddAcquire(delta uint64) uint64 {
	return arch.AddUint64Acquire(&a.v, delta)
}

// AddRelease atomically adds delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Uint64) AddRelease(delta uint64) uint64 {
	return arch.AddUint64Release(&a.v, delta)
}

// AddAcqRel atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) AddAcqRel(delta uint64) uint64 {
	return arch.AddUint64AcqRel(&a.v, delta)
}

// Sub atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) Sub(delta uint64) uint64 {
	return arch.AddUint64AcqRel(&a.v, ^(delta - 1))
}

// SubRelaxed atomically subtracts delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Uint64) SubRelaxed(delta uint64) uint64 {
	return arch.AddUint64Relaxed(&a.v, ^(delta - 1))
}

// SubAcquire atomically subtracts delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Uint64) SubAcquire(delta uint64) uint64 {
	return arch.AddUint64Acquire(&a.v, ^(delta - 1))
}

// SubRelease atomically subtracts delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Uint64) SubRelease(delta uint64) uint64 {
	return arch.AddUint64Release(&a.v, ^(delta - 1))
}

// SubAcqRel atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) SubAcqRel(delta uint64) uint64 {
	return arch.AddUint64AcqRel(&a.v, ^(delta - 1))
}

// And atomically performs bitwise AND and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) And(mask uint64) uint64 {
	return arch.AndUint64AcqRel(&a.v, mask)
}

// AndRelaxed atomically performs bitwise AND with relaxed ordering.
//
//go:nosplit
func (a *Uint64) AndRelaxed(mask uint64) uint64 {
	return arch.AndUint64Relaxed(&a.v, mask)
}

// AndAcquire atomically performs bitwise AND with acquire ordering.
//
//go:nosplit
func (a *Uint64) AndAcquire(mask uint64) uint64 {
	return arch.AndUint64Acquire(&a.v, mask)
}

// AndRelease atomically performs bitwise AND with release ordering.
//
//go:nosplit
func (a *Uint64) AndRelease(mask uint64) uint64 {
	return arch.AndUint64Release(&a.v, mask)
}

// AndAcqRel atomically performs bitwise AND with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) AndAcqRel(mask uint64) uint64 {
	return arch.AndUint64AcqRel(&a.v, mask)
}

// Or atomically performs bitwise OR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) Or(mask uint64) uint64 {
	return arch.OrUint64AcqRel(&a.v, mask)
}

// OrRelaxed atomically performs bitwise OR with relaxed ordering.
//
//go:nosplit
func (a *Uint64) OrRelaxed(mask uint64) uint64 {
	return arch.OrUint64Relaxed(&a.v, mask)
}

// OrAcquire atomically performs bitwise OR with acquire ordering.
//
//go:nosplit
func (a *Uint64) OrAcquire(mask uint64) uint64 {
	return arch.OrUint64Acquire(&a.v, mask)
}

// OrRelease atomically performs bitwise OR with release ordering.
//
//go:nosplit
func (a *Uint64) OrRelease(mask uint64) uint64 {
	return arch.OrUint64Release(&a.v, mask)
}

// OrAcqRel atomically performs bitwise OR with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) OrAcqRel(mask uint64) uint64 {
	return arch.OrUint64AcqRel(&a.v, mask)
}

// Xor atomically performs bitwise XOR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) Xor(mask uint64) uint64 {
	return arch.XorUint64AcqRel(&a.v, mask)
}

// XorRelaxed atomically performs bitwise XOR with relaxed ordering.
//
//go:nosplit
func (a *Uint64) XorRelaxed(mask uint64) uint64 {
	return arch.XorUint64Relaxed(&a.v, mask)
}

// XorAcquire atomically performs bitwise XOR with acquire ordering.
//
//go:nosplit
func (a *Uint64) XorAcquire(mask uint64) uint64 {
	return arch.XorUint64Acquire(&a.v, mask)
}

// XorRelease atomically performs bitwise XOR with release ordering.
//
//go:nosplit
func (a *Uint64) XorRelease(mask uint64) uint64 {
	return arch.XorUint64Release(&a.v, mask)
}

// XorAcqRel atomically performs bitwise XOR with acquire-release ordering.
//
//go:nosplit
func (a *Uint64) XorAcqRel(mask uint64) uint64 {
	return arch.XorUint64AcqRel(&a.v, mask)
}

// Max atomically stores the maximum of current and val, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint64) Max(val uint64) uint64 {
	for {
		old := arch.LoadUint64Relaxed(&a.v)
		if old >= val {
			return old
		}
		if arch.CasUint64AcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MaxRelaxed atomically stores the maximum with relaxed ordering.
//
//go:nosplit
func (a *Uint64) MaxRelaxed(val uint64) uint64 {
	for {
		old := arch.LoadUint64Relaxed(&a.v)
		if old >= val {
			return old
		}
		if arch.CasUint64Relaxed(&a.v, old, val) {
			return old
		}
	}
}

// Min atomically stores the minimum of current and val, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint64) Min(val uint64) uint64 {
	for {
		old := arch.LoadUint64Relaxed(&a.v)
		if old <= val {
			return old
		}
		if arch.CasUint64AcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MinRelaxed atomically stores the minimum with relaxed ordering.
//
//go:nosplit
func (a *Uint64) MinRelaxed(val uint64) uint64 {
	for {
		old := arch.LoadUint64Relaxed(&a.v)
		if old <= val {
			return old
		}
		if arch.CasUint64Relaxed(&a.v, old, val) {
			return old
		}
	}
}
