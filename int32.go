// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// Load atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Int32) Load() int32 {
	return arch.LoadInt32Relaxed(&a.v)
}

// LoadRelaxed atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Int32) LoadRelaxed() int32 {
	return arch.LoadInt32Relaxed(&a.v)
}

// LoadAcquire atomically loads and returns the value with acquire ordering.
//
//go:nosplit
func (a *Int32) LoadAcquire() int32 {
	return arch.LoadInt32Acquire(&a.v)
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Int32) Store(val int32) {
	arch.StoreInt32Relaxed(&a.v, val)
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Int32) StoreRelaxed(val int32) {
	arch.StoreInt32Relaxed(&a.v, val)
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Int32) StoreRelease(val int32) {
	arch.StoreInt32Release(&a.v, val)
}

// Swap atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) Swap(new int32) int32 {
	return arch.SwapInt32AcqRel(&a.v, new)
}

// SwapRelaxed atomically swaps the value and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Int32) SwapRelaxed(new int32) int32 {
	return arch.SwapInt32Relaxed(&a.v, new)
}

// SwapAcquire atomically swaps the value and returns the old value with acquire ordering.
//
//go:nosplit
func (a *Int32) SwapAcquire(new int32) int32 {
	return arch.SwapInt32Acquire(&a.v, new)
}

// SwapRelease atomically swaps the value and returns the old value with release ordering.
//
//go:nosplit
func (a *Int32) SwapRelease(new int32) int32 {
	return arch.SwapInt32Release(&a.v, new)
}

// SwapAcqRel atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) SwapAcqRel(new int32) int32 {
	return arch.SwapInt32AcqRel(&a.v, new)
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Int32) CompareAndSwap(old, new int32) bool {
	return arch.CasInt32AcqRel(&a.v, old, new)
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Int32) CompareAndSwapRelaxed(old, new int32) bool {
	return arch.CasInt32Relaxed(&a.v, old, new)
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Int32) CompareAndSwapAcquire(old, new int32) bool {
	return arch.CasInt32Acquire(&a.v, old, new)
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Int32) CompareAndSwapRelease(old, new int32) bool {
	return arch.CasInt32Release(&a.v, old, new)
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Int32) CompareAndSwapAcqRel(old, new int32) bool {
	return arch.CasInt32AcqRel(&a.v, old, new)
}

// CompareExchange atomically compares and swaps, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int32) CompareExchange(old, new int32) int32 {
	return arch.CaxInt32AcqRel(&a.v, old, new)
}

// CompareExchangeRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Int32) CompareExchangeRelaxed(old, new int32) int32 {
	return arch.CaxInt32Relaxed(&a.v, old, new)
}

// CompareExchangeAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Int32) CompareExchangeAcquire(old, new int32) int32 {
	return arch.CaxInt32Acquire(&a.v, old, new)
}

// CompareExchangeRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Int32) CompareExchangeRelease(old, new int32) int32 {
	return arch.CaxInt32Release(&a.v, old, new)
}

// CompareExchangeAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Int32) CompareExchangeAcqRel(old, new int32) int32 {
	return arch.CaxInt32AcqRel(&a.v, old, new)
}

// Add atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) Add(delta int32) int32 {
	return arch.AddInt32AcqRel(&a.v, delta)
}

// AddRelaxed atomically adds delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Int32) AddRelaxed(delta int32) int32 {
	return arch.AddInt32Relaxed(&a.v, delta)
}

// AddAcquire atomically adds delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Int32) AddAcquire(delta int32) int32 {
	return arch.AddInt32Acquire(&a.v, delta)
}

// AddRelease atomically adds delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Int32) AddRelease(delta int32) int32 {
	return arch.AddInt32Release(&a.v, delta)
}

// AddAcqRel atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) AddAcqRel(delta int32) int32 {
	return arch.AddInt32AcqRel(&a.v, delta)
}

// Sub atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) Sub(delta int32) int32 {
	return arch.AddInt32AcqRel(&a.v, -delta)
}

// SubRelaxed atomically subtracts delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Int32) SubRelaxed(delta int32) int32 {
	return arch.AddInt32Relaxed(&a.v, -delta)
}

// SubAcquire atomically subtracts delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Int32) SubAcquire(delta int32) int32 {
	return arch.AddInt32Acquire(&a.v, -delta)
}

// SubRelease atomically subtracts delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Int32) SubRelease(delta int32) int32 {
	return arch.AddInt32Release(&a.v, -delta)
}

// SubAcqRel atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) SubAcqRel(delta int32) int32 {
	return arch.AddInt32AcqRel(&a.v, -delta)
}

// And atomically performs bitwise AND and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) And(mask int32) int32 {
	return arch.AndInt32AcqRel(&a.v, mask)
}

// AndRelaxed atomically performs bitwise AND with relaxed ordering.
//
//go:nosplit
func (a *Int32) AndRelaxed(mask int32) int32 {
	return arch.AndInt32Relaxed(&a.v, mask)
}

// AndAcquire atomically performs bitwise AND with acquire ordering.
//
//go:nosplit
func (a *Int32) AndAcquire(mask int32) int32 {
	return arch.AndInt32Acquire(&a.v, mask)
}

// AndRelease atomically performs bitwise AND with release ordering.
//
//go:nosplit
func (a *Int32) AndRelease(mask int32) int32 {
	return arch.AndInt32Release(&a.v, mask)
}

// AndAcqRel atomically performs bitwise AND with acquire-release ordering.
//
//go:nosplit
func (a *Int32) AndAcqRel(mask int32) int32 {
	return arch.AndInt32AcqRel(&a.v, mask)
}

// Or atomically performs bitwise OR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) Or(mask int32) int32 {
	return arch.OrInt32AcqRel(&a.v, mask)
}

// OrRelaxed atomically performs bitwise OR with relaxed ordering.
//
//go:nosplit
func (a *Int32) OrRelaxed(mask int32) int32 {
	return arch.OrInt32Relaxed(&a.v, mask)
}

// OrAcquire atomically performs bitwise OR with acquire ordering.
//
//go:nosplit
func (a *Int32) OrAcquire(mask int32) int32 {
	return arch.OrInt32Acquire(&a.v, mask)
}

// OrRelease atomically performs bitwise OR with release ordering.
//
//go:nosplit
func (a *Int32) OrRelease(mask int32) int32 {
	return arch.OrInt32Release(&a.v, mask)
}

// OrAcqRel atomically performs bitwise OR with acquire-release ordering.
//
//go:nosplit
func (a *Int32) OrAcqRel(mask int32) int32 {
	return arch.OrInt32AcqRel(&a.v, mask)
}

// Xor atomically performs bitwise XOR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int32) Xor(mask int32) int32 {
	return arch.XorInt32AcqRel(&a.v, mask)
}

// XorRelaxed atomically performs bitwise XOR with relaxed ordering.
//
//go:nosplit
func (a *Int32) XorRelaxed(mask int32) int32 {
	return arch.XorInt32Relaxed(&a.v, mask)
}

// XorAcquire atomically performs bitwise XOR with acquire ordering.
//
//go:nosplit
func (a *Int32) XorAcquire(mask int32) int32 {
	return arch.XorInt32Acquire(&a.v, mask)
}

// XorRelease atomically performs bitwise XOR with release ordering.
//
//go:nosplit
func (a *Int32) XorRelease(mask int32) int32 {
	return arch.XorInt32Release(&a.v, mask)
}

// XorAcqRel atomically performs bitwise XOR with acquire-release ordering.
//
//go:nosplit
func (a *Int32) XorAcqRel(mask int32) int32 {
	return arch.XorInt32AcqRel(&a.v, mask)
}

// Max atomically stores the maximum of current and val, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int32) Max(val int32) int32 {
	for {
		old := arch.LoadInt32Relaxed(&a.v)
		if old >= val {
			return old
		}
		if arch.CasInt32AcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MaxRelaxed atomically stores the maximum with relaxed ordering.
//
//go:nosplit
func (a *Int32) MaxRelaxed(val int32) int32 {
	for {
		old := arch.LoadInt32Relaxed(&a.v)
		if old >= val {
			return old
		}
		if arch.CasInt32Relaxed(&a.v, old, val) {
			return old
		}
	}
}

// Min atomically stores the minimum of current and val, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int32) Min(val int32) int32 {
	for {
		old := arch.LoadInt32Relaxed(&a.v)
		if old <= val {
			return old
		}
		if arch.CasInt32AcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MinRelaxed atomically stores the minimum with relaxed ordering.
//
//go:nosplit
func (a *Int32) MinRelaxed(val int32) int32 {
	for {
		old := arch.LoadInt32Relaxed(&a.v)
		if old <= val {
			return old
		}
		if arch.CasInt32Relaxed(&a.v, old, val) {
			return old
		}
	}
}
