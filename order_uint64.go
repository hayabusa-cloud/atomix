// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// LoadUint64 atomically loads *addr with the specified memory ordering.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadUint64(addr *uint64) uint64 {
	if o == Relaxed {
		return arch.LoadUint64Relaxed(addr)
	}
	return arch.LoadUint64Acquire(addr)
}

// StoreUint64 atomically stores val to *addr with the specified memory ordering.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StoreUint64(addr *uint64, val uint64) {
	if o == Relaxed {
		arch.StoreUint64Relaxed(addr, val)
		return
	}
	arch.StoreUint64Release(addr, val)
}

// SwapUint64 atomically stores new to *addr and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapUint64(addr *uint64, new uint64) (old uint64) {
	switch o {
	case Relaxed:
		return arch.SwapUint64Relaxed(addr, new)
	case Acquire:
		return arch.SwapUint64Acquire(addr, new)
	case Release:
		return arch.SwapUint64Release(addr, new)
	default:
		return arch.SwapUint64AcqRel(addr, new)
	}
}

// CompareAndSwapUint64 atomically compares *addr with old and swaps if equal.
// Returns true if the swap was performed.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool) {
	switch o {
	case Relaxed:
		return arch.CasUint64Relaxed(addr, old, new)
	case Acquire:
		return arch.CasUint64Acquire(addr, old, new)
	case Release:
		return arch.CasUint64Release(addr, old, new)
	default:
		return arch.CasUint64AcqRel(addr, old, new)
	}
}

// CompareExchangeUint64 atomically compares *addr with old and swaps if equal.
// Returns the previous value (enables CAS loops without separate Load).
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareExchangeUint64(addr *uint64, old, new uint64) (prev uint64) {
	switch o {
	case Relaxed:
		return arch.CaxUint64Relaxed(addr, old, new)
	case Acquire:
		return arch.CaxUint64Acquire(addr, old, new)
	case Release:
		return arch.CaxUint64Release(addr, old, new)
	default:
		return arch.CaxUint64AcqRel(addr, old, new)
	}
}

// AddUint64 atomically adds delta to *addr and returns the new value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AddUint64(addr *uint64, delta uint64) (new uint64) {
	switch o {
	case Relaxed:
		return arch.AddUint64Relaxed(addr, delta)
	case Acquire:
		return arch.AddUint64Acquire(addr, delta)
	case Release:
		return arch.AddUint64Release(addr, delta)
	default:
		return arch.AddUint64AcqRel(addr, delta)
	}
}

// AndUint64 atomically performs *addr &= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AndUint64(addr *uint64, mask uint64) (old uint64) {
	switch o {
	case Relaxed:
		return arch.AndUint64Relaxed(addr, mask)
	case Acquire:
		return arch.AndUint64Acquire(addr, mask)
	case Release:
		return arch.AndUint64Release(addr, mask)
	default:
		return arch.AndUint64AcqRel(addr, mask)
	}
}

// OrUint64 atomically performs *addr |= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) OrUint64(addr *uint64, mask uint64) (old uint64) {
	switch o {
	case Relaxed:
		return arch.OrUint64Relaxed(addr, mask)
	case Acquire:
		return arch.OrUint64Acquire(addr, mask)
	case Release:
		return arch.OrUint64Release(addr, mask)
	default:
		return arch.OrUint64AcqRel(addr, mask)
	}
}

// XorUint64 atomically performs *addr ^= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) XorUint64(addr *uint64, mask uint64) (old uint64) {
	switch o {
	case Relaxed:
		return arch.XorUint64Relaxed(addr, mask)
	case Acquire:
		return arch.XorUint64Acquire(addr, mask)
	case Release:
		return arch.XorUint64Release(addr, mask)
	default:
		return arch.XorUint64AcqRel(addr, mask)
	}
}

// MaxUint64 atomically stores max(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MaxUint64(addr *uint64, val uint64) (old uint64) {
	if o == Relaxed {
		for {
			old = arch.LoadUint64Relaxed(addr)
			if old >= val || arch.CasUint64Relaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadUint64Relaxed(addr)
		if old >= val || arch.CasUint64AcqRel(addr, old, val) {
			return old
		}
	}
}

// MinUint64 atomically stores min(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MinUint64(addr *uint64, val uint64) (old uint64) {
	if o == Relaxed {
		for {
			old = arch.LoadUint64Relaxed(addr)
			if old <= val || arch.CasUint64Relaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadUint64Relaxed(addr)
		if old <= val || arch.CasUint64AcqRel(addr, old, val) {
			return old
		}
	}
}
