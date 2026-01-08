// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// LoadUint128 atomically loads *addr with the specified memory ordering.
// Returns (lo, hi) where the full value is (hi << 64) | lo.
// addr MUST be 16-byte aligned; use PlaceAlignedUint128 to ensure alignment.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadUint128(addr *Uint128) (lo, hi uint64) {
	if o == Relaxed {
		return arch.LoadUint128Relaxed(&addr.v)
	}
	return arch.LoadUint128Acquire(&addr.v)
}

// StoreUint128 atomically stores (lo, hi) to *addr with the specified memory ordering.
// addr MUST be 16-byte aligned; use PlaceAlignedUint128 to ensure alignment.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StoreUint128(addr *Uint128, lo, hi uint64) {
	if o == Relaxed {
		arch.StoreUint128Relaxed(&addr.v, lo, hi)
		return
	}
	arch.StoreUint128Release(&addr.v, lo, hi)
}

// SwapUint128 atomically stores (newLo, newHi) to *addr and returns the old value.
// addr MUST be 16-byte aligned; use PlaceAlignedUint128 to ensure alignment.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapUint128(addr *Uint128, newLo, newHi uint64) (oldLo, oldHi uint64) {
	switch o {
	case Relaxed:
		return arch.SwapUint128Relaxed(&addr.v, newLo, newHi)
	case Acquire:
		return arch.SwapUint128Acquire(&addr.v, newLo, newHi)
	case Release:
		return arch.SwapUint128Release(&addr.v, newLo, newHi)
	default:
		return arch.SwapUint128AcqRel(&addr.v, newLo, newHi)
	}
}

// CompareAndSwapUint128 atomically compares *addr with (oldLo, oldHi) and swaps if equal.
// Returns true if the swap was performed.
// addr MUST be 16-byte aligned; use PlaceAlignedUint128 to ensure alignment.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapUint128(addr *Uint128, oldLo, oldHi, newLo, newHi uint64) (swapped bool) {
	switch o {
	case Relaxed:
		return arch.CasUint128Relaxed(&addr.v, oldLo, oldHi, newLo, newHi)
	case Acquire:
		return arch.CasUint128Acquire(&addr.v, oldLo, oldHi, newLo, newHi)
	case Release:
		return arch.CasUint128Release(&addr.v, oldLo, oldHi, newLo, newHi)
	default:
		return arch.CasUint128AcqRel(&addr.v, oldLo, oldHi, newLo, newHi)
	}
}

// CompareExchangeUint128 atomically compares *addr with (oldLo, oldHi) and swaps if equal.
// Returns the previous value (enables CAS loops without separate Load).
// addr MUST be 16-byte aligned; use PlaceAlignedUint128 to ensure alignment.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareExchangeUint128(addr *Uint128, oldLo, oldHi, newLo, newHi uint64) (prevLo, prevHi uint64) {
	switch o {
	case Relaxed:
		return arch.CaxUint128Relaxed(&addr.v, oldLo, oldHi, newLo, newHi)
	case Acquire:
		return arch.CaxUint128Acquire(&addr.v, oldLo, oldHi, newLo, newHi)
	case Release:
		return arch.CaxUint128Release(&addr.v, oldLo, oldHi, newLo, newHi)
	default:
		return arch.CaxUint128AcqRel(&addr.v, oldLo, oldHi, newLo, newHi)
	}
}

// AddUint128 atomically adds (deltaLo, deltaHi) to *addr and returns the old value.
// addr MUST be 16-byte aligned; use PlaceAlignedUint128 to ensure alignment.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AddUint128(addr *Uint128, deltaLo, deltaHi uint64) (oldLo, oldHi uint64) {
	if o == Relaxed {
		for {
			ulo, uhi := arch.LoadUint128Relaxed(&addr.v)
			newLo := ulo + deltaLo
			newHi := uhi + deltaHi
			if newLo < ulo { // Carry
				newHi++
			}
			if arch.CasUint128Relaxed(&addr.v, ulo, uhi, newLo, newHi) {
				return ulo, uhi
			}
		}
	}
	for {
		ulo, uhi := arch.LoadUint128Relaxed(&addr.v)
		newLo := ulo + deltaLo
		newHi := uhi + deltaHi
		if newLo < ulo { // Carry
			newHi++
		}
		if arch.CasUint128AcqRel(&addr.v, ulo, uhi, newLo, newHi) {
			return ulo, uhi
		}
	}
}
