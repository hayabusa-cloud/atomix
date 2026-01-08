// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import (
	"unsafe"

	"code.hybscloud.com/atomix/internal/arch"
)

// LoadPointer atomically loads *addr with the specified memory ordering.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadPointer(addr *unsafe.Pointer) unsafe.Pointer {
	if o == Relaxed {
		return arch.LoadPointerRelaxed(addr)
	}
	return arch.LoadPointerAcquire(addr)
}

// StorePointer atomically stores val to *addr with the specified memory ordering.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StorePointer(addr *unsafe.Pointer, val unsafe.Pointer) {
	if o == Relaxed {
		arch.StorePointerRelaxed(addr, val)
		return
	}
	arch.StorePointerRelease(addr, val)
}

// SwapPointer atomically stores new to *addr and returns the old value.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer) {
	switch o {
	case Relaxed:
		return arch.SwapPointerRelaxed(addr, new)
	case Acquire:
		return arch.SwapPointerAcquire(addr, new)
	case Release:
		return arch.SwapPointerRelease(addr, new)
	default:
		return arch.SwapPointerAcqRel(addr, new)
	}
}

// CompareAndSwapPointer atomically compares *addr with old and swaps if equal.
// Returns true if the swap was performed.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool) {
	switch o {
	case Relaxed:
		return arch.CasPointerRelaxed(addr, old, new)
	case Acquire:
		return arch.CasPointerAcquire(addr, old, new)
	case Release:
		return arch.CasPointerRelease(addr, old, new)
	default:
		return arch.CasPointerAcqRel(addr, old, new)
	}
}

// CompareExchangePointer atomically compares *addr with old and swaps if equal.
// Returns the previous value (enables CAS loops without separate Load).
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareExchangePointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (prev unsafe.Pointer) {
	switch o {
	case Relaxed:
		return arch.CaxPointerRelaxed(addr, old, new)
	case Acquire:
		return arch.CaxPointerAcquire(addr, old, new)
	case Release:
		return arch.CaxPointerRelease(addr, old, new)
	default:
		return arch.CaxPointerAcqRel(addr, old, new)
	}
}
