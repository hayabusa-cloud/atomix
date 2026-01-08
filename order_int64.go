// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// LoadInt64 atomically loads *addr with the specified memory ordering.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadInt64(addr *int64) int64 {
	if o == Relaxed {
		return arch.LoadInt64Relaxed(addr)
	}
	return arch.LoadInt64Acquire(addr)
}

// StoreInt64 atomically stores val to *addr with the specified memory ordering.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StoreInt64(addr *int64, val int64) {
	if o == Relaxed {
		arch.StoreInt64Relaxed(addr, val)
		return
	}
	arch.StoreInt64Release(addr, val)
}

// SwapInt64 atomically stores new to *addr and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapInt64(addr *int64, new int64) (old int64) {
	switch o {
	case Relaxed:
		return arch.SwapInt64Relaxed(addr, new)
	case Acquire:
		return arch.SwapInt64Acquire(addr, new)
	case Release:
		return arch.SwapInt64Release(addr, new)
	default:
		return arch.SwapInt64AcqRel(addr, new)
	}
}

// CompareAndSwapInt64 atomically compares *addr with old and swaps if equal.
// Returns true if the swap was performed.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool) {
	switch o {
	case Relaxed:
		return arch.CasInt64Relaxed(addr, old, new)
	case Acquire:
		return arch.CasInt64Acquire(addr, old, new)
	case Release:
		return arch.CasInt64Release(addr, old, new)
	default:
		return arch.CasInt64AcqRel(addr, old, new)
	}
}

// CompareExchangeInt64 atomically compares *addr with old and swaps if equal.
// Returns the previous value (enables CAS loops without separate Load).
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareExchangeInt64(addr *int64, old, new int64) (prev int64) {
	switch o {
	case Relaxed:
		return arch.CaxInt64Relaxed(addr, old, new)
	case Acquire:
		return arch.CaxInt64Acquire(addr, old, new)
	case Release:
		return arch.CaxInt64Release(addr, old, new)
	default:
		return arch.CaxInt64AcqRel(addr, old, new)
	}
}

// AddInt64 atomically adds delta to *addr and returns the new value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AddInt64(addr *int64, delta int64) (new int64) {
	switch o {
	case Relaxed:
		return arch.AddInt64Relaxed(addr, delta)
	case Acquire:
		return arch.AddInt64Acquire(addr, delta)
	case Release:
		return arch.AddInt64Release(addr, delta)
	default:
		return arch.AddInt64AcqRel(addr, delta)
	}
}

// AndInt64 atomically performs *addr &= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) AndInt64(addr *int64, mask int64) (old int64) {
	switch o {
	case Relaxed:
		return arch.AndInt64Relaxed(addr, mask)
	case Acquire:
		return arch.AndInt64Acquire(addr, mask)
	case Release:
		return arch.AndInt64Release(addr, mask)
	default:
		return arch.AndInt64AcqRel(addr, mask)
	}
}

// OrInt64 atomically performs *addr |= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) OrInt64(addr *int64, mask int64) (old int64) {
	switch o {
	case Relaxed:
		return arch.OrInt64Relaxed(addr, mask)
	case Acquire:
		return arch.OrInt64Acquire(addr, mask)
	case Release:
		return arch.OrInt64Release(addr, mask)
	default:
		return arch.OrInt64AcqRel(addr, mask)
	}
}

// XorInt64 atomically performs *addr ^= mask and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) XorInt64(addr *int64, mask int64) (old int64) {
	switch o {
	case Relaxed:
		return arch.XorInt64Relaxed(addr, mask)
	case Acquire:
		return arch.XorInt64Acquire(addr, mask)
	case Release:
		return arch.XorInt64Release(addr, mask)
	default:
		return arch.XorInt64AcqRel(addr, mask)
	}
}

// MaxInt64 atomically stores max(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MaxInt64(addr *int64, val int64) (old int64) {
	if o == Relaxed {
		for {
			old = arch.LoadInt64Relaxed(addr)
			if old >= val || arch.CasInt64Relaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadInt64Relaxed(addr)
		if old >= val || arch.CasInt64AcqRel(addr, old, val) {
			return old
		}
	}
}

// MinInt64 atomically stores min(*addr, val) and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) MinInt64(addr *int64, val int64) (old int64) {
	if o == Relaxed {
		for {
			old = arch.LoadInt64Relaxed(addr)
			if old <= val || arch.CasInt64Relaxed(addr, old, val) {
				return old
			}
		}
	}
	for {
		old = arch.LoadInt64Relaxed(addr)
		if old <= val || arch.CasInt64AcqRel(addr, old, val) {
			return old
		}
	}
}
