// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// Load atomically loads and returns the value with relaxed ordering.
// Returns (lo, hi) where the full value is (hi << 64) | lo.
//
//go:nosplit
func (a *Int128) Load() (lo, hi int64) {
	ulo, uhi := arch.LoadUint128Relaxed(&a.v)
	return int64(ulo), int64(uhi)
}

// LoadRelaxed atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Int128) LoadRelaxed() (lo, hi int64) {
	ulo, uhi := arch.LoadUint128Relaxed(&a.v)
	return int64(ulo), int64(uhi)
}

// LoadAcquire atomically loads and returns the value with acquire ordering.
//
//go:nosplit
func (a *Int128) LoadAcquire() (lo, hi int64) {
	ulo, uhi := arch.LoadUint128Acquire(&a.v)
	return int64(ulo), int64(uhi)
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Int128) Store(lo, hi int64) {
	arch.StoreUint128Relaxed(&a.v, uint64(lo), uint64(hi))
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Int128) StoreRelaxed(lo, hi int64) {
	arch.StoreUint128Relaxed(&a.v, uint64(lo), uint64(hi))
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Int128) StoreRelease(lo, hi int64) {
	arch.StoreUint128Release(&a.v, uint64(lo), uint64(hi))
}

// Swap atomically stores new value and returns the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int128) Swap(newLo, newHi int64) (oldLo, oldHi int64) {
	ulo, uhi := arch.SwapUint128AcqRel(&a.v, uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// SwapRelaxed atomically stores new value and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Int128) SwapRelaxed(newLo, newHi int64) (oldLo, oldHi int64) {
	ulo, uhi := arch.SwapUint128Relaxed(&a.v, uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// SwapAcquire atomically stores new value and returns the old value with acquire ordering.
//
//go:nosplit
func (a *Int128) SwapAcquire(newLo, newHi int64) (oldLo, oldHi int64) {
	ulo, uhi := arch.SwapUint128Acquire(&a.v, uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// SwapRelease atomically stores new value and returns the old value with release ordering.
//
//go:nosplit
func (a *Int128) SwapRelease(newLo, newHi int64) (oldLo, oldHi int64) {
	ulo, uhi := arch.SwapUint128Release(&a.v, uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// SwapAcqRel atomically stores new value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int128) SwapAcqRel(newLo, newHi int64) (oldLo, oldHi int64) {
	ulo, uhi := arch.SwapUint128AcqRel(&a.v, uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Int128) CompareAndSwap(oldLo, oldHi, newLo, newHi int64) bool {
	return arch.CasUint128AcqRel(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Int128) CompareAndSwapRelaxed(oldLo, oldHi, newLo, newHi int64) bool {
	return arch.CasUint128Relaxed(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Int128) CompareAndSwapAcquire(oldLo, oldHi, newLo, newHi int64) bool {
	return arch.CasUint128Acquire(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Int128) CompareAndSwapRelease(oldLo, oldHi, newLo, newHi int64) bool {
	return arch.CasUint128Release(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Int128) CompareAndSwapAcqRel(oldLo, oldHi, newLo, newHi int64) bool {
	return arch.CasUint128AcqRel(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
}

// CompareExchange atomically compares and swaps, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int128) CompareExchange(oldLo, oldHi, newLo, newHi int64) (lo, hi int64) {
	ulo, uhi := arch.CaxUint128AcqRel(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// CompareExchangeRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Int128) CompareExchangeRelaxed(oldLo, oldHi, newLo, newHi int64) (lo, hi int64) {
	ulo, uhi := arch.CaxUint128Relaxed(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// CompareExchangeAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Int128) CompareExchangeAcquire(oldLo, oldHi, newLo, newHi int64) (lo, hi int64) {
	ulo, uhi := arch.CaxUint128Acquire(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// CompareExchangeRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Int128) CompareExchangeRelease(oldLo, oldHi, newLo, newHi int64) (lo, hi int64) {
	ulo, uhi := arch.CaxUint128Release(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// CompareExchangeAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Int128) CompareExchangeAcqRel(oldLo, oldHi, newLo, newHi int64) (lo, hi int64) {
	ulo, uhi := arch.CaxUint128AcqRel(&a.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	return int64(ulo), int64(uhi)
}

// Add atomically adds (deltaLo, deltaHi) and returns the new value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int128) Add(deltaLo, deltaHi int64) (lo, hi int64) {
	for {
		ulo, uhi := arch.LoadUint128Relaxed(&a.v)
		newLo := ulo + uint64(deltaLo)
		newHi := uhi + uint64(deltaHi)
		if newLo < ulo { // Carry
			newHi++
		}
		if arch.CasUint128AcqRel(&a.v, ulo, uhi, newLo, newHi) {
			return int64(newLo), int64(newHi)
		}
	}
}

// AddRelaxed atomically adds (deltaLo, deltaHi) and returns the new value.
//
//go:nosplit
func (a *Int128) AddRelaxed(deltaLo, deltaHi int64) (lo, hi int64) {
	for {
		ulo, uhi := arch.LoadUint128Relaxed(&a.v)
		newLo := ulo + uint64(deltaLo)
		newHi := uhi + uint64(deltaHi)
		if newLo < ulo {
			newHi++
		}
		if arch.CasUint128Relaxed(&a.v, ulo, uhi, newLo, newHi) {
			return int64(newLo), int64(newHi)
		}
	}
}

// Sub atomically subtracts (deltaLo, deltaHi) and returns the new value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int128) Sub(deltaLo, deltaHi int64) (lo, hi int64) {
	for {
		ulo, uhi := arch.LoadUint128Relaxed(&a.v)
		newLo := ulo - uint64(deltaLo)
		newHi := uhi - uint64(deltaHi)
		if newLo > ulo { // Borrow
			newHi--
		}
		if arch.CasUint128AcqRel(&a.v, ulo, uhi, newLo, newHi) {
			return int64(newLo), int64(newHi)
		}
	}
}

// SubRelaxed atomically subtracts (deltaLo, deltaHi) and returns the new value.
//
//go:nosplit
func (a *Int128) SubRelaxed(deltaLo, deltaHi int64) (lo, hi int64) {
	for {
		ulo, uhi := arch.LoadUint128Relaxed(&a.v)
		newLo := ulo - uint64(deltaLo)
		newHi := uhi - uint64(deltaHi)
		if newLo > ulo {
			newHi--
		}
		if arch.CasUint128Relaxed(&a.v, ulo, uhi, newLo, newHi) {
			return int64(newLo), int64(newHi)
		}
	}
}

// Inc atomically increments by 1 and returns the new value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int128) Inc() (lo, hi int64) {
	return a.Add(1, 0)
}

// IncRelaxed atomically increments by 1 and returns the new value.
//
//go:nosplit
func (a *Int128) IncRelaxed() (lo, hi int64) {
	return a.AddRelaxed(1, 0)
}

// Dec atomically decrements by 1 and returns the new value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int128) Dec() (lo, hi int64) {
	return a.Sub(1, 0)
}

// DecRelaxed atomically decrements by 1 and returns the new value.
//
//go:nosplit
func (a *Int128) DecRelaxed() (lo, hi int64) {
	return a.SubRelaxed(1, 0)
}

// Equal atomically loads and compares for equality.
// Uses acquire ordering.
//
//go:nosplit
func (a *Int128) Equal(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Acquire(&a.v)
	return int64(ulo) == lo && int64(uhi) == hi
}

// EqualRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Int128) EqualRelaxed(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Relaxed(&a.v)
	return int64(ulo) == lo && int64(uhi) == hi
}

// Less atomically loads and returns true if value < (lo, hi).
// Uses acquire ordering. Comparison is signed.
//
//go:nosplit
func (a *Int128) Less(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Acquire(&a.v)
	ahi := int64(uhi)
	alo := int64(ulo)
	return ahi < hi || (ahi == hi && uint64(alo) < uint64(lo))
}

// LessRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Int128) LessRelaxed(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Relaxed(&a.v)
	ahi := int64(uhi)
	alo := int64(ulo)
	return ahi < hi || (ahi == hi && uint64(alo) < uint64(lo))
}

// LessOrEqual atomically loads and returns true if value <= (lo, hi).
// Uses acquire ordering. Comparison is signed.
//
//go:nosplit
func (a *Int128) LessOrEqual(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Acquire(&a.v)
	ahi := int64(uhi)
	alo := int64(ulo)
	return ahi < hi || (ahi == hi && uint64(alo) <= uint64(lo))
}

// LessOrEqualRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Int128) LessOrEqualRelaxed(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Relaxed(&a.v)
	ahi := int64(uhi)
	alo := int64(ulo)
	return ahi < hi || (ahi == hi && uint64(alo) <= uint64(lo))
}

// Greater atomically loads and returns true if value > (lo, hi).
// Uses acquire ordering. Comparison is signed.
//
//go:nosplit
func (a *Int128) Greater(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Acquire(&a.v)
	ahi := int64(uhi)
	alo := int64(ulo)
	return ahi > hi || (ahi == hi && uint64(alo) > uint64(lo))
}

// GreaterRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Int128) GreaterRelaxed(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Relaxed(&a.v)
	ahi := int64(uhi)
	alo := int64(ulo)
	return ahi > hi || (ahi == hi && uint64(alo) > uint64(lo))
}

// GreaterOrEqual atomically loads and returns true if value >= (lo, hi).
// Uses acquire ordering. Comparison is signed.
//
//go:nosplit
func (a *Int128) GreaterOrEqual(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Acquire(&a.v)
	ahi := int64(uhi)
	alo := int64(ulo)
	return ahi > hi || (ahi == hi && uint64(alo) >= uint64(lo))
}

// GreaterOrEqualRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Int128) GreaterOrEqualRelaxed(lo, hi int64) bool {
	ulo, uhi := arch.LoadUint128Relaxed(&a.v)
	ahi := int64(uhi)
	alo := int64(ulo)
	return ahi > hi || (ahi == hi && uint64(alo) >= uint64(lo))
}
