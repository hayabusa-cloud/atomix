// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// Load atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) Load() uintptr {
	return arch.LoadUintptrRelaxed(&a.v)
}

// LoadRelaxed atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) LoadRelaxed() uintptr {
	return arch.LoadUintptrRelaxed(&a.v)
}

// LoadAcquire atomically loads and returns the value with acquire ordering.
//
//go:nosplit
func (a *Uintptr) LoadAcquire() uintptr {
	return arch.LoadUintptrAcquire(&a.v)
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) Store(val uintptr) {
	arch.StoreUintptrRelaxed(&a.v, val)
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) StoreRelaxed(val uintptr) {
	arch.StoreUintptrRelaxed(&a.v, val)
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Uintptr) StoreRelease(val uintptr) {
	arch.StoreUintptrRelease(&a.v, val)
}

// Swap atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) Swap(new uintptr) uintptr {
	return arch.SwapUintptrAcqRel(&a.v, new)
}

// SwapRelaxed atomically swaps the value and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) SwapRelaxed(new uintptr) uintptr {
	return arch.SwapUintptrRelaxed(&a.v, new)
}

// SwapAcquire atomically swaps the value and returns the old value with acquire ordering.
//
//go:nosplit
func (a *Uintptr) SwapAcquire(new uintptr) uintptr {
	return arch.SwapUintptrAcquire(&a.v, new)
}

// SwapRelease atomically swaps the value and returns the old value with release ordering.
//
//go:nosplit
func (a *Uintptr) SwapRelease(new uintptr) uintptr {
	return arch.SwapUintptrRelease(&a.v, new)
}

// SwapAcqRel atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) SwapAcqRel(new uintptr) uintptr {
	return arch.SwapUintptrAcqRel(&a.v, new)
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Uintptr) CompareAndSwap(old, new uintptr) bool {
	return arch.CasUintptrAcqRel(&a.v, old, new)
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) CompareAndSwapRelaxed(old, new uintptr) bool {
	return arch.CasUintptrRelaxed(&a.v, old, new)
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Uintptr) CompareAndSwapAcquire(old, new uintptr) bool {
	return arch.CasUintptrAcquire(&a.v, old, new)
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Uintptr) CompareAndSwapRelease(old, new uintptr) bool {
	return arch.CasUintptrRelease(&a.v, old, new)
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) CompareAndSwapAcqRel(old, new uintptr) bool {
	return arch.CasUintptrAcqRel(&a.v, old, new)
}

// CompareExchange atomically compares and swaps, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) CompareExchange(old, new uintptr) uintptr {
	return arch.CaxUintptrAcqRel(&a.v, old, new)
}

// CompareExchangeRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) CompareExchangeRelaxed(old, new uintptr) uintptr {
	return arch.CaxUintptrRelaxed(&a.v, old, new)
}

// CompareExchangeAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Uintptr) CompareExchangeAcquire(old, new uintptr) uintptr {
	return arch.CaxUintptrAcquire(&a.v, old, new)
}

// CompareExchangeRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Uintptr) CompareExchangeRelease(old, new uintptr) uintptr {
	return arch.CaxUintptrRelease(&a.v, old, new)
}

// CompareExchangeAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) CompareExchangeAcqRel(old, new uintptr) uintptr {
	return arch.CaxUintptrAcqRel(&a.v, old, new)
}

// Add atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) Add(delta uintptr) uintptr {
	return arch.AddUintptrAcqRel(&a.v, delta)
}

// AddRelaxed atomically adds delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) AddRelaxed(delta uintptr) uintptr {
	return arch.AddUintptrRelaxed(&a.v, delta)
}

// AddAcquire atomically adds delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Uintptr) AddAcquire(delta uintptr) uintptr {
	return arch.AddUintptrAcquire(&a.v, delta)
}

// AddRelease atomically adds delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Uintptr) AddRelease(delta uintptr) uintptr {
	return arch.AddUintptrRelease(&a.v, delta)
}

// AddAcqRel atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) AddAcqRel(delta uintptr) uintptr {
	return arch.AddUintptrAcqRel(&a.v, delta)
}

// Sub atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) Sub(delta uintptr) uintptr {
	return arch.AddUintptrAcqRel(&a.v, ^(delta - 1))
}

// SubRelaxed atomically subtracts delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) SubRelaxed(delta uintptr) uintptr {
	return arch.AddUintptrRelaxed(&a.v, ^(delta - 1))
}

// SubAcquire atomically subtracts delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Uintptr) SubAcquire(delta uintptr) uintptr {
	return arch.AddUintptrAcquire(&a.v, ^(delta - 1))
}

// SubRelease atomically subtracts delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Uintptr) SubRelease(delta uintptr) uintptr {
	return arch.AddUintptrRelease(&a.v, ^(delta - 1))
}

// SubAcqRel atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) SubAcqRel(delta uintptr) uintptr {
	return arch.AddUintptrAcqRel(&a.v, ^(delta - 1))
}

// And atomically performs bitwise AND and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) And(mask uintptr) uintptr {
	return arch.AndUintptrAcqRel(&a.v, mask)
}

// AndRelaxed atomically performs bitwise AND with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) AndRelaxed(mask uintptr) uintptr {
	return arch.AndUintptrRelaxed(&a.v, mask)
}

// AndAcquire atomically performs bitwise AND with acquire ordering.
//
//go:nosplit
func (a *Uintptr) AndAcquire(mask uintptr) uintptr {
	return arch.AndUintptrAcquire(&a.v, mask)
}

// AndRelease atomically performs bitwise AND with release ordering.
//
//go:nosplit
func (a *Uintptr) AndRelease(mask uintptr) uintptr {
	return arch.AndUintptrRelease(&a.v, mask)
}

// AndAcqRel atomically performs bitwise AND with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) AndAcqRel(mask uintptr) uintptr {
	return arch.AndUintptrAcqRel(&a.v, mask)
}

// Or atomically performs bitwise OR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) Or(mask uintptr) uintptr {
	return arch.OrUintptrAcqRel(&a.v, mask)
}

// OrRelaxed atomically performs bitwise OR with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) OrRelaxed(mask uintptr) uintptr {
	return arch.OrUintptrRelaxed(&a.v, mask)
}

// OrAcquire atomically performs bitwise OR with acquire ordering.
//
//go:nosplit
func (a *Uintptr) OrAcquire(mask uintptr) uintptr {
	return arch.OrUintptrAcquire(&a.v, mask)
}

// OrRelease atomically performs bitwise OR with release ordering.
//
//go:nosplit
func (a *Uintptr) OrRelease(mask uintptr) uintptr {
	return arch.OrUintptrRelease(&a.v, mask)
}

// OrAcqRel atomically performs bitwise OR with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) OrAcqRel(mask uintptr) uintptr {
	return arch.OrUintptrAcqRel(&a.v, mask)
}

// Xor atomically performs bitwise XOR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) Xor(mask uintptr) uintptr {
	return arch.XorUintptrAcqRel(&a.v, mask)
}

// XorRelaxed atomically performs bitwise XOR with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) XorRelaxed(mask uintptr) uintptr {
	return arch.XorUintptrRelaxed(&a.v, mask)
}

// XorAcquire atomically performs bitwise XOR with acquire ordering.
//
//go:nosplit
func (a *Uintptr) XorAcquire(mask uintptr) uintptr {
	return arch.XorUintptrAcquire(&a.v, mask)
}

// XorRelease atomically performs bitwise XOR with release ordering.
//
//go:nosplit
func (a *Uintptr) XorRelease(mask uintptr) uintptr {
	return arch.XorUintptrRelease(&a.v, mask)
}

// XorAcqRel atomically performs bitwise XOR with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) XorAcqRel(mask uintptr) uintptr {
	return arch.XorUintptrAcqRel(&a.v, mask)
}

// Max atomically stores max(*addr, val) and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) Max(val uintptr) uintptr {
	for {
		old := arch.LoadUintptrRelaxed(&a.v)
		if old >= val || arch.CasUintptrAcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MaxRelaxed atomically stores max(*addr, val) and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) MaxRelaxed(val uintptr) uintptr {
	for {
		old := arch.LoadUintptrRelaxed(&a.v)
		if old >= val || arch.CasUintptrRelaxed(&a.v, old, val) {
			return old
		}
	}
}

// Min atomically stores min(*addr, val) and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Uintptr) Min(val uintptr) uintptr {
	for {
		old := arch.LoadUintptrRelaxed(&a.v)
		if old <= val || arch.CasUintptrAcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MinRelaxed atomically stores min(*addr, val) and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Uintptr) MinRelaxed(val uintptr) uintptr {
	for {
		old := arch.LoadUintptrRelaxed(&a.v)
		if old <= val || arch.CasUintptrRelaxed(&a.v, old, val) {
			return old
		}
	}
}
