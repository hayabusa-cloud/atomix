// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// LoadInt128 atomically loads *addr with the specified memory ordering.
// Returns (lo, hi) where the full value is (hi << 64) | lo.
// addr MUST be 16-byte aligned; use PlaceAlignedInt128 to ensure alignment.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadInt128(addr *Int128) (lo, hi int64) {
	if o == Relaxed {
		ulo, uhi := arch.LoadUint128Relaxed(&addr.v)
		return int64(ulo), int64(uhi)
	}
	ulo, uhi := arch.LoadUint128Acquire(&addr.v)
	return int64(ulo), int64(uhi)
}

// StoreInt128 atomically stores (lo, hi) to *addr with the specified memory ordering.
// addr MUST be 16-byte aligned; use PlaceAlignedInt128 to ensure alignment.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StoreInt128(addr *Int128, lo, hi int64) {
	if o == Relaxed {
		arch.StoreUint128Relaxed(&addr.v, uint64(lo), uint64(hi))
		return
	}
	arch.StoreUint128Release(&addr.v, uint64(lo), uint64(hi))
}

// SwapInt128 atomically stores (newLo, newHi) to *addr and returns the old value.
// addr MUST be 16-byte aligned; use PlaceAlignedInt128 to ensure alignment.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapInt128(addr *Int128, newLo, newHi int64) (oldLo, oldHi int64) {
	var ulo, uhi uint64
	switch o {
	case Relaxed:
		ulo, uhi = arch.SwapUint128Relaxed(&addr.v, uint64(newLo), uint64(newHi))
	case Acquire:
		ulo, uhi = arch.SwapUint128Acquire(&addr.v, uint64(newLo), uint64(newHi))
	case Release:
		ulo, uhi = arch.SwapUint128Release(&addr.v, uint64(newLo), uint64(newHi))
	default:
		ulo, uhi = arch.SwapUint128AcqRel(&addr.v, uint64(newLo), uint64(newHi))
	}
	return int64(ulo), int64(uhi)
}

// CompareAndSwapInt128 atomically compares *addr with (oldLo, oldHi) and swaps if equal.
// Returns true if the swap was performed.
// addr MUST be 16-byte aligned; use PlaceAlignedInt128 to ensure alignment.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapInt128(addr *Int128, oldLo, oldHi, newLo, newHi int64) (swapped bool) {
	switch o {
	case Relaxed:
		return arch.CasUint128Relaxed(&addr.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	case Acquire:
		return arch.CasUint128Acquire(&addr.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	case Release:
		return arch.CasUint128Release(&addr.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	default:
		return arch.CasUint128AcqRel(&addr.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	}
}

// CompareExchangeInt128 atomically compares *addr with (oldLo, oldHi) and swaps if equal.
// Returns the previous value (enables CAS loops without separate Load).
// addr MUST be 16-byte aligned; use PlaceAlignedInt128 to ensure alignment.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareExchangeInt128(addr *Int128, oldLo, oldHi, newLo, newHi int64) (prevLo, prevHi int64) {
	var ulo, uhi uint64
	switch o {
	case Relaxed:
		ulo, uhi = arch.CaxUint128Relaxed(&addr.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	case Acquire:
		ulo, uhi = arch.CaxUint128Acquire(&addr.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	case Release:
		ulo, uhi = arch.CaxUint128Release(&addr.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	default:
		ulo, uhi = arch.CaxUint128AcqRel(&addr.v, uint64(oldLo), uint64(oldHi), uint64(newLo), uint64(newHi))
	}
	return int64(ulo), int64(uhi)
}

// AddInt128 atomically adds (deltaLo, deltaHi) to *addr and returns the old value.
// addr MUST be 16-byte aligned; use PlaceAlignedInt128 to ensure alignment.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AddInt128(addr *Int128, deltaLo, deltaHi int64) (oldLo, oldHi int64) {
	if o == Relaxed {
		for {
			ulo, uhi := arch.LoadUint128Relaxed(&addr.v)
			newLo := ulo + uint64(deltaLo)
			newHi := uhi + uint64(deltaHi)
			if newLo < ulo { // Carry
				newHi++
			}
			if arch.CasUint128Relaxed(&addr.v, ulo, uhi, newLo, newHi) {
				return int64(ulo), int64(uhi)
			}
		}
	}
	for {
		ulo, uhi := arch.LoadUint128Relaxed(&addr.v)
		newLo := ulo + uint64(deltaLo)
		newHi := uhi + uint64(deltaHi)
		if newLo < ulo { // Carry
			newHi++
		}
		if arch.CasUint128AcqRel(&addr.v, ulo, uhi, newLo, newHi) {
			return int64(ulo), int64(uhi)
		}
	}
}
