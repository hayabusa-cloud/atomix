// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// LoadInt32 atomically loads *addr with the specified memory ordering.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadInt32(addr *int32) int32 {
	if o == Relaxed {
		return arch.LoadInt32Relaxed(addr)
	}
	return arch.LoadInt32Acquire(addr)
}

// StoreInt32 atomically stores val to *addr with the specified memory ordering.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StoreInt32(addr *int32, val int32) {
	if o == Relaxed {
		arch.StoreInt32Relaxed(addr, val)
		return
	}
	arch.StoreInt32Release(addr, val)
}

// SwapInt32 atomically stores new to *addr and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapInt32(addr *int32, new int32) (old int32) {
	switch o {
	case Relaxed:
		return arch.SwapInt32Relaxed(addr, new)
	case Acquire:
		return arch.SwapInt32Acquire(addr, new)
	case Release:
		return arch.SwapInt32Release(addr, new)
	default:
		return arch.SwapInt32AcqRel(addr, new)
	}
}

// CompareAndSwapInt32 atomically compares *addr with old and swaps if equal.
// Returns true if the swap was performed.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool) {
	switch o {
	case Relaxed:
		return arch.CasInt32Relaxed(addr, old, new)
	case Acquire:
		return arch.CasInt32Acquire(addr, old, new)
	case Release:
		return arch.CasInt32Release(addr, old, new)
	default:
		return arch.CasInt32AcqRel(addr, old, new)
	}
}

// CompareExchangeInt32 atomically compares *addr with old and swaps if equal.
// Returns the previous value (enables CAS loops without separate Load).
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareExchangeInt32(addr *int32, old, new int32) (prev int32) {
	switch o {
	case Relaxed:
		return arch.CaxInt32Relaxed(addr, old, new)
	case Acquire:
		return arch.CaxInt32Acquire(addr, old, new)
	case Release:
		return arch.CaxInt32Release(addr, old, new)
	default:
		return arch.CaxInt32AcqRel(addr, old, new)
	}
}

// AddInt32 atomically adds delta to *addr and returns the new value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AddInt32(addr *int32, delta int32) (new int32) {
	switch o {
	case Relaxed:
		return arch.AddInt32Relaxed(addr, delta)
	case Acquire:
		return arch.AddInt32Acquire(addr, delta)
	case Release:
		return arch.AddInt32Release(addr, delta)
	default:
		return arch.AddInt32AcqRel(addr, delta)
	}
}

// AndInt32 atomically performs *addr &= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AndInt32(addr *int32, mask int32) (old int32) {
	switch o {
	case Relaxed:
		return arch.AndInt32Relaxed(addr, mask)
	case Acquire:
		return arch.AndInt32Acquire(addr, mask)
	case Release:
		return arch.AndInt32Release(addr, mask)
	default:
		return arch.AndInt32AcqRel(addr, mask)
	}
}

// OrInt32 atomically performs *addr |= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) OrInt32(addr *int32, mask int32) (old int32) {
	switch o {
	case Relaxed:
		return arch.OrInt32Relaxed(addr, mask)
	case Acquire:
		return arch.OrInt32Acquire(addr, mask)
	case Release:
		return arch.OrInt32Release(addr, mask)
	default:
		return arch.OrInt32AcqRel(addr, mask)
	}
}

// XorInt32 atomically performs *addr ^= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) XorInt32(addr *int32, mask int32) (old int32) {
	switch o {
	case Relaxed:
		return arch.XorInt32Relaxed(addr, mask)
	case Acquire:
		return arch.XorInt32Acquire(addr, mask)
	case Release:
		return arch.XorInt32Release(addr, mask)
	default:
		return arch.XorInt32AcqRel(addr, mask)
	}
}

// MaxInt32 atomically stores max(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MaxInt32(addr *int32, val int32) (old int32) {
	if o == Relaxed {
		for {
			old = arch.LoadInt32Relaxed(addr)
			if old >= val || arch.CasInt32Relaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadInt32Relaxed(addr)
		if old >= val || arch.CasInt32AcqRel(addr, old, val) {
			return old
		}
	}
}

// MinInt32 atomically stores min(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MinInt32(addr *int32, val int32) (old int32) {
	if o == Relaxed {
		for {
			old = arch.LoadInt32Relaxed(addr)
			if old <= val || arch.CasInt32Relaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadInt32Relaxed(addr)
		if old <= val || arch.CasInt32AcqRel(addr, old, val) {
			return old
		}
	}
}
