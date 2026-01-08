// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// LoadUint32 atomically loads *addr with the specified memory ordering.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadUint32(addr *uint32) uint32 {
	if o == Relaxed {
		return arch.LoadUint32Relaxed(addr)
	}
	return arch.LoadUint32Acquire(addr)
}

// StoreUint32 atomically stores val to *addr with the specified memory ordering.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StoreUint32(addr *uint32, val uint32) {
	if o == Relaxed {
		arch.StoreUint32Relaxed(addr, val)
		return
	}
	arch.StoreUint32Release(addr, val)
}

// SwapUint32 atomically stores new to *addr and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapUint32(addr *uint32, new uint32) (old uint32) {
	switch o {
	case Relaxed:
		return arch.SwapUint32Relaxed(addr, new)
	case Acquire:
		return arch.SwapUint32Acquire(addr, new)
	case Release:
		return arch.SwapUint32Release(addr, new)
	default:
		return arch.SwapUint32AcqRel(addr, new)
	}
}

// CompareAndSwapUint32 atomically compares *addr with old and swaps if equal.
// Returns true if the swap was performed.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool) {
	switch o {
	case Relaxed:
		return arch.CasUint32Relaxed(addr, old, new)
	case Acquire:
		return arch.CasUint32Acquire(addr, old, new)
	case Release:
		return arch.CasUint32Release(addr, old, new)
	default:
		return arch.CasUint32AcqRel(addr, old, new)
	}
}

// CompareExchangeUint32 atomically compares *addr with old and swaps if equal.
// Returns the previous value (enables CAS loops without separate Load).
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareExchangeUint32(addr *uint32, old, new uint32) (prev uint32) {
	switch o {
	case Relaxed:
		return arch.CaxUint32Relaxed(addr, old, new)
	case Acquire:
		return arch.CaxUint32Acquire(addr, old, new)
	case Release:
		return arch.CaxUint32Release(addr, old, new)
	default:
		return arch.CaxUint32AcqRel(addr, old, new)
	}
}

// AddUint32 atomically adds delta to *addr and returns the new value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AddUint32(addr *uint32, delta uint32) (new uint32) {
	switch o {
	case Relaxed:
		return arch.AddUint32Relaxed(addr, delta)
	case Acquire:
		return arch.AddUint32Acquire(addr, delta)
	case Release:
		return arch.AddUint32Release(addr, delta)
	default:
		return arch.AddUint32AcqRel(addr, delta)
	}
}

// AndUint32 atomically performs *addr &= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AndUint32(addr *uint32, mask uint32) (old uint32) {
	switch o {
	case Relaxed:
		return arch.AndUint32Relaxed(addr, mask)
	case Acquire:
		return arch.AndUint32Acquire(addr, mask)
	case Release:
		return arch.AndUint32Release(addr, mask)
	default:
		return arch.AndUint32AcqRel(addr, mask)
	}
}

// OrUint32 atomically performs *addr |= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) OrUint32(addr *uint32, mask uint32) (old uint32) {
	switch o {
	case Relaxed:
		return arch.OrUint32Relaxed(addr, mask)
	case Acquire:
		return arch.OrUint32Acquire(addr, mask)
	case Release:
		return arch.OrUint32Release(addr, mask)
	default:
		return arch.OrUint32AcqRel(addr, mask)
	}
}

// XorUint32 atomically performs *addr ^= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) XorUint32(addr *uint32, mask uint32) (old uint32) {
	switch o {
	case Relaxed:
		return arch.XorUint32Relaxed(addr, mask)
	case Acquire:
		return arch.XorUint32Acquire(addr, mask)
	case Release:
		return arch.XorUint32Release(addr, mask)
	default:
		return arch.XorUint32AcqRel(addr, mask)
	}
}

// MaxUint32 atomically stores max(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MaxUint32(addr *uint32, val uint32) (old uint32) {
	if o == Relaxed {
		for {
			old = arch.LoadUint32Relaxed(addr)
			if old >= val || arch.CasUint32Relaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadUint32Relaxed(addr)
		if old >= val || arch.CasUint32AcqRel(addr, old, val) {
			return old
		}
	}
}

// MinUint32 atomically stores min(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MinUint32(addr *uint32, val uint32) (old uint32) {
	if o == Relaxed {
		for {
			old = arch.LoadUint32Relaxed(addr)
			if old <= val || arch.CasUint32Relaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadUint32Relaxed(addr)
		if old <= val || arch.CasUint32AcqRel(addr, old, val) {
			return old
		}
	}
}
