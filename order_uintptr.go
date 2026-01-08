// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// LoadUintptr atomically loads *addr with the specified memory ordering.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadUintptr(addr *uintptr) uintptr {
	if o == Relaxed {
		return arch.LoadUintptrRelaxed(addr)
	}
	return arch.LoadUintptrAcquire(addr)
}

// StoreUintptr atomically stores val to *addr with the specified memory ordering.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StoreUintptr(addr *uintptr, val uintptr) {
	if o == Relaxed {
		arch.StoreUintptrRelaxed(addr, val)
		return
	}
	arch.StoreUintptrRelease(addr, val)
}

// SwapUintptr atomically stores new to *addr and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapUintptr(addr *uintptr, new uintptr) (old uintptr) {
	switch o {
	case Relaxed:
		return arch.SwapUintptrRelaxed(addr, new)
	case Acquire:
		return arch.SwapUintptrAcquire(addr, new)
	case Release:
		return arch.SwapUintptrRelease(addr, new)
	default:
		return arch.SwapUintptrAcqRel(addr, new)
	}
}

// CompareAndSwapUintptr atomically compares *addr with old and swaps if equal.
// Returns true if the swap was performed.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool) {
	switch o {
	case Relaxed:
		return arch.CasUintptrRelaxed(addr, old, new)
	case Acquire:
		return arch.CasUintptrAcquire(addr, old, new)
	case Release:
		return arch.CasUintptrRelease(addr, old, new)
	default:
		return arch.CasUintptrAcqRel(addr, old, new)
	}
}

// CompareExchangeUintptr atomically compares *addr with old and swaps if equal.
// Returns the previous value (enables CAS loops without separate Load).
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareExchangeUintptr(addr *uintptr, old, new uintptr) (prev uintptr) {
	switch o {
	case Relaxed:
		return arch.CaxUintptrRelaxed(addr, old, new)
	case Acquire:
		return arch.CaxUintptrAcquire(addr, old, new)
	case Release:
		return arch.CaxUintptrRelease(addr, old, new)
	default:
		return arch.CaxUintptrAcqRel(addr, old, new)
	}
}

// AddUintptr atomically adds delta to *addr and returns the new value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AddUintptr(addr *uintptr, delta uintptr) (new uintptr) {
	switch o {
	case Relaxed:
		return arch.AddUintptrRelaxed(addr, delta)
	case Acquire:
		return arch.AddUintptrAcquire(addr, delta)
	case Release:
		return arch.AddUintptrRelease(addr, delta)
	default:
		return arch.AddUintptrAcqRel(addr, delta)
	}
}

// AndUintptr atomically performs *addr &= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AndUintptr(addr *uintptr, mask uintptr) (old uintptr) {
	switch o {
	case Relaxed:
		return arch.AndUintptrRelaxed(addr, mask)
	case Acquire:
		return arch.AndUintptrAcquire(addr, mask)
	case Release:
		return arch.AndUintptrRelease(addr, mask)
	default:
		return arch.AndUintptrAcqRel(addr, mask)
	}
}

// OrUintptr atomically performs *addr |= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) OrUintptr(addr *uintptr, mask uintptr) (old uintptr) {
	switch o {
	case Relaxed:
		return arch.OrUintptrRelaxed(addr, mask)
	case Acquire:
		return arch.OrUintptrAcquire(addr, mask)
	case Release:
		return arch.OrUintptrRelease(addr, mask)
	default:
		return arch.OrUintptrAcqRel(addr, mask)
	}
}

// XorUintptr atomically performs *addr ^= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) XorUintptr(addr *uintptr, mask uintptr) (old uintptr) {
	switch o {
	case Relaxed:
		return arch.XorUintptrRelaxed(addr, mask)
	case Acquire:
		return arch.XorUintptrAcquire(addr, mask)
	case Release:
		return arch.XorUintptrRelease(addr, mask)
	default:
		return arch.XorUintptrAcqRel(addr, mask)
	}
}

// MaxUintptr atomically stores max(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MaxUintptr(addr *uintptr, val uintptr) (old uintptr) {
	if o == Relaxed {
		for {
			old = arch.LoadUintptrRelaxed(addr)
			if old >= val || arch.CasUintptrRelaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadUintptrRelaxed(addr)
		if old >= val || arch.CasUintptrAcqRel(addr, old, val) {
			return old
		}
	}
}

// MinUintptr atomically stores min(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MinUintptr(addr *uintptr, val uintptr) (old uintptr) {
	if o == Relaxed {
		for {
			old = arch.LoadUintptrRelaxed(addr)
			if old <= val || arch.CasUintptrRelaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadUintptrRelaxed(addr)
		if old <= val || arch.CasUintptrAcqRel(addr, old, val) {
			return old
		}
	}
}
