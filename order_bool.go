// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// LoadBool atomically loads *addr with the specified memory ordering.
// addr points to a uint32 where 0 is false and non-zero is true.
// Unknown orderings fallback to Acquire.
//
//go:nosplit
func (o MemoryOrder) LoadBool(addr *uint32) bool {
	if o == Relaxed {
		return arch.LoadUint32Relaxed(addr) != 0
	}
	return arch.LoadUint32Acquire(addr) != 0
}

// StoreBool atomically stores val to *addr with the specified memory ordering.
// addr points to a uint32 where 0 is false and 1 is true.
// Unknown orderings fallback to Release.
//
//go:nosplit
func (o MemoryOrder) StoreBool(addr *uint32, val bool) {
	var v uint32
	if val {
		v = 1
	}
	if o == Relaxed {
		arch.StoreUint32Relaxed(addr, v)
		return
	}
	arch.StoreUint32Release(addr, v)
}

// SwapBool atomically stores new to *addr and returns the old value.
// addr points to a uint32 where 0 is false and non-zero is true.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) SwapBool(addr *uint32, new bool) (old bool) {
	var v uint32
	if new {
		v = 1
	}
	switch o {
	case Relaxed:
		return arch.SwapUint32Relaxed(addr, v) != 0
	case Acquire:
		return arch.SwapUint32Acquire(addr, v) != 0
	case Release:
		return arch.SwapUint32Release(addr, v) != 0
	default:
		return arch.SwapUint32AcqRel(addr, v) != 0
	}
}

// CompareAndSwapBool atomically compares *addr with old and swaps if equal.
// Returns true if the swap was performed.
// addr points to a uint32 where 0 is false and 1 is true.
// Unknown orderings fallback to AcqRel.
//
//go:nosplit
func (o MemoryOrder) CompareAndSwapBool(addr *uint32, old, new bool) (swapped bool) {
	var oldV, newV uint32
	if old {
		oldV = 1
	}
	if new {
		newV = 1
	}
	switch o {
	case Relaxed:
		return arch.CasUint32Relaxed(addr, oldV, newV)
	case Acquire:
		return arch.CasUint32Acquire(addr, oldV, newV)
	case Release:
		return arch.CasUint32Release(addr, oldV, newV)
	default:
		return arch.CasUint32AcqRel(addr, oldV, newV)
	}
}
