// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// Load atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Uint128) Load() (lo, hi uint64) {
	return arch.LoadUint128Relaxed(&a.v)
}

// LoadRelaxed atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Uint128) LoadRelaxed() (lo, hi uint64) {
	return arch.LoadUint128Relaxed(&a.v)
}

// LoadAcquire atomically loads and returns the value with acquire ordering.
//
//go:nosplit
func (a *Uint128) LoadAcquire() (lo, hi uint64) {
	return arch.LoadUint128Acquire(&a.v)
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Uint128) Store(lo, hi uint64) {
	arch.StoreUint128Relaxed(&a.v, lo, hi)
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Uint128) StoreRelaxed(lo, hi uint64) {
	arch.StoreUint128Relaxed(&a.v, lo, hi)
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Uint128) StoreRelease(lo, hi uint64) {
	arch.StoreUint128Release(&a.v, lo, hi)
}

// Swap atomically stores new value and returns the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint128) Swap(newLo, newHi uint64) (oldLo, oldHi uint64) {
	return arch.SwapUint128AcqRel(&a.v, newLo, newHi)
}

// SwapRelaxed atomically stores new value and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Uint128) SwapRelaxed(newLo, newHi uint64) (oldLo, oldHi uint64) {
	return arch.SwapUint128Relaxed(&a.v, newLo, newHi)
}

// SwapAcquire atomically stores new value and returns the old value with acquire ordering.
//
//go:nosplit
func (a *Uint128) SwapAcquire(newLo, newHi uint64) (oldLo, oldHi uint64) {
	return arch.SwapUint128Acquire(&a.v, newLo, newHi)
}

// SwapRelease atomically stores new value and returns the old value with release ordering.
//
//go:nosplit
func (a *Uint128) SwapRelease(newLo, newHi uint64) (oldLo, oldHi uint64) {
	return arch.SwapUint128Release(&a.v, newLo, newHi)
}

// SwapAcqRel atomically stores new value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uint128) SwapAcqRel(newLo, newHi uint64) (oldLo, oldHi uint64) {
	return arch.SwapUint128AcqRel(&a.v, newLo, newHi)
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Uint128) CompareAndSwap(oldLo, oldHi, newLo, newHi uint64) bool {
	return arch.CasUint128AcqRel(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Uint128) CompareAndSwapRelaxed(oldLo, oldHi, newLo, newHi uint64) bool {
	return arch.CasUint128Relaxed(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Uint128) CompareAndSwapAcquire(oldLo, oldHi, newLo, newHi uint64) bool {
	return arch.CasUint128Acquire(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Uint128) CompareAndSwapRelease(oldLo, oldHi, newLo, newHi uint64) bool {
	return arch.CasUint128Release(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Uint128) CompareAndSwapAcqRel(oldLo, oldHi, newLo, newHi uint64) bool {
	return arch.CasUint128AcqRel(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareExchange atomically compares and swaps, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint128) CompareExchange(oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	return arch.CaxUint128AcqRel(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareExchangeRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Uint128) CompareExchangeRelaxed(oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	return arch.CaxUint128Relaxed(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareExchangeAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Uint128) CompareExchangeAcquire(oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	return arch.CaxUint128Acquire(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareExchangeRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Uint128) CompareExchangeRelease(oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	return arch.CaxUint128Release(&a.v, oldLo, oldHi, newLo, newHi)
}

// CompareExchangeAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Uint128) CompareExchangeAcqRel(oldLo, oldHi, newLo, newHi uint64) (lo, hi uint64) {
	return arch.CaxUint128AcqRel(&a.v, oldLo, oldHi, newLo, newHi)
}

// Add atomically adds (deltaLo, deltaHi) and returns the new value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint128) Add(deltaLo, deltaHi uint64) (newLo, newHi uint64) {
	for {
		oldLo, oldHi := arch.LoadUint128Relaxed(&a.v)
		newLo = oldLo + deltaLo
		newHi = oldHi + deltaHi
		if newLo < oldLo { // Carry
			newHi++
		}
		if arch.CasUint128AcqRel(&a.v, oldLo, oldHi, newLo, newHi) {
			return newLo, newHi
		}
	}
}

// AddRelaxed atomically adds (deltaLo, deltaHi) and returns the new value.
//
//go:nosplit
func (a *Uint128) AddRelaxed(deltaLo, deltaHi uint64) (newLo, newHi uint64) {
	for {
		oldLo, oldHi := arch.LoadUint128Relaxed(&a.v)
		newLo = oldLo + deltaLo
		newHi = oldHi + deltaHi
		if newLo < oldLo {
			newHi++
		}
		if arch.CasUint128Relaxed(&a.v, oldLo, oldHi, newLo, newHi) {
			return newLo, newHi
		}
	}
}

// Sub atomically subtracts (deltaLo, deltaHi) and returns the new value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint128) Sub(deltaLo, deltaHi uint64) (newLo, newHi uint64) {
	for {
		oldLo, oldHi := arch.LoadUint128Relaxed(&a.v)
		newLo = oldLo - deltaLo
		newHi = oldHi - deltaHi
		if newLo > oldLo { // Borrow
			newHi--
		}
		if arch.CasUint128AcqRel(&a.v, oldLo, oldHi, newLo, newHi) {
			return newLo, newHi
		}
	}
}

// SubRelaxed atomically subtracts (deltaLo, deltaHi) and returns the new value.
//
//go:nosplit
func (a *Uint128) SubRelaxed(deltaLo, deltaHi uint64) (newLo, newHi uint64) {
	for {
		oldLo, oldHi := arch.LoadUint128Relaxed(&a.v)
		newLo = oldLo - deltaLo
		newHi = oldHi - deltaHi
		if newLo > oldLo {
			newHi--
		}
		if arch.CasUint128Relaxed(&a.v, oldLo, oldHi, newLo, newHi) {
			return newLo, newHi
		}
	}
}

// Inc atomically increments by 1 and returns the new value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint128) Inc() (newLo, newHi uint64) {
	return a.Add(1, 0)
}

// IncRelaxed atomically increments by 1 and returns the new value.
//
//go:nosplit
func (a *Uint128) IncRelaxed() (newLo, newHi uint64) {
	return a.AddRelaxed(1, 0)
}

// Dec atomically decrements by 1 and returns the new value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uint128) Dec() (newLo, newHi uint64) {
	return a.Sub(1, 0)
}

// DecRelaxed atomically decrements by 1 and returns the new value.
//
//go:nosplit
func (a *Uint128) DecRelaxed() (newLo, newHi uint64) {
	return a.SubRelaxed(1, 0)
}

// Equal atomically loads and compares for equality.
// Uses acquire ordering.
//
//go:nosplit
func (a *Uint128) Equal(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Acquire(&a.v)
	return alo == lo && ahi == hi
}

// EqualRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Uint128) EqualRelaxed(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Relaxed(&a.v)
	return alo == lo && ahi == hi
}

// Less atomically loads and returns true if value < (lo, hi).
// Uses acquire ordering.
//
//go:nosplit
func (a *Uint128) Less(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Acquire(&a.v)
	return ahi < hi || (ahi == hi && alo < lo)
}

// LessRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Uint128) LessRelaxed(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Relaxed(&a.v)
	return ahi < hi || (ahi == hi && alo < lo)
}

// LessOrEqual atomically loads and returns true if value <= (lo, hi).
// Uses acquire ordering.
//
//go:nosplit
func (a *Uint128) LessOrEqual(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Acquire(&a.v)
	return ahi < hi || (ahi == hi && alo <= lo)
}

// LessOrEqualRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Uint128) LessOrEqualRelaxed(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Relaxed(&a.v)
	return ahi < hi || (ahi == hi && alo <= lo)
}

// Greater atomically loads and returns true if value > (lo, hi).
// Uses acquire ordering.
//
//go:nosplit
func (a *Uint128) Greater(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Acquire(&a.v)
	return ahi > hi || (ahi == hi && alo > lo)
}

// GreaterRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Uint128) GreaterRelaxed(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Relaxed(&a.v)
	return ahi > hi || (ahi == hi && alo > lo)
}

// GreaterOrEqual atomically loads and returns true if value >= (lo, hi).
// Uses acquire ordering.
//
//go:nosplit
func (a *Uint128) GreaterOrEqual(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Acquire(&a.v)
	return ahi > hi || (ahi == hi && alo >= lo)
}

// GreaterOrEqualRelaxed atomically loads and compares with relaxed ordering.
//
//go:nosplit
func (a *Uint128) GreaterOrEqualRelaxed(lo, hi uint64) bool {
	alo, ahi := arch.LoadUint128Relaxed(&a.v)
	return ahi > hi || (ahi == hi && alo >= lo)
}
