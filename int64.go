// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// Load atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Int64) Load() int64 {
	return arch.LoadInt64Relaxed(&a.v)
}

// LoadRelaxed atomically loads and returns the value with relaxed ordering.
//
//go:nosplit
func (a *Int64) LoadRelaxed() int64 {
	return arch.LoadInt64Relaxed(&a.v)
}

// LoadAcquire atomically loads and returns the value with acquire ordering.
//
//go:nosplit
func (a *Int64) LoadAcquire() int64 {
	return arch.LoadInt64Acquire(&a.v)
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Int64) Store(val int64) {
	arch.StoreInt64Relaxed(&a.v, val)
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Int64) StoreRelaxed(val int64) {
	arch.StoreInt64Relaxed(&a.v, val)
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Int64) StoreRelease(val int64) {
	arch.StoreInt64Release(&a.v, val)
}

// Swap atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) Swap(new int64) int64 {
	return arch.SwapInt64AcqRel(&a.v, new)
}

// SwapRelaxed atomically swaps the value and returns the old value with relaxed ordering.
//
//go:nosplit
func (a *Int64) SwapRelaxed(new int64) int64 {
	return arch.SwapInt64Relaxed(&a.v, new)
}

// SwapAcquire atomically swaps the value and returns the old value with acquire ordering.
//
//go:nosplit
func (a *Int64) SwapAcquire(new int64) int64 {
	return arch.SwapInt64Acquire(&a.v, new)
}

// SwapRelease atomically swaps the value and returns the old value with release ordering.
//
//go:nosplit
func (a *Int64) SwapRelease(new int64) int64 {
	return arch.SwapInt64Release(&a.v, new)
}

// SwapAcqRel atomically swaps the value and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) SwapAcqRel(new int64) int64 {
	return arch.SwapInt64AcqRel(&a.v, new)
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Int64) CompareAndSwap(old, new int64) bool {
	return arch.CasInt64AcqRel(&a.v, old, new)
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Int64) CompareAndSwapRelaxed(old, new int64) bool {
	return arch.CasInt64Relaxed(&a.v, old, new)
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Int64) CompareAndSwapAcquire(old, new int64) bool {
	return arch.CasInt64Acquire(&a.v, old, new)
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Int64) CompareAndSwapRelease(old, new int64) bool {
	return arch.CasInt64Release(&a.v, old, new)
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Int64) CompareAndSwapAcqRel(old, new int64) bool {
	return arch.CasInt64AcqRel(&a.v, old, new)
}

// CompareExchange atomically compares and swaps, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int64) CompareExchange(old, new int64) int64 {
	return arch.CaxInt64AcqRel(&a.v, old, new)
}

// CompareExchangeRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Int64) CompareExchangeRelaxed(old, new int64) int64 {
	return arch.CaxInt64Relaxed(&a.v, old, new)
}

// CompareExchangeAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Int64) CompareExchangeAcquire(old, new int64) int64 {
	return arch.CaxInt64Acquire(&a.v, old, new)
}

// CompareExchangeRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Int64) CompareExchangeRelease(old, new int64) int64 {
	return arch.CaxInt64Release(&a.v, old, new)
}

// CompareExchangeAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Int64) CompareExchangeAcqRel(old, new int64) int64 {
	return arch.CaxInt64AcqRel(&a.v, old, new)
}

// Add atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) Add(delta int64) int64 {
	return arch.AddInt64AcqRel(&a.v, delta)
}

// AddRelaxed atomically adds delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Int64) AddRelaxed(delta int64) int64 {
	return arch.AddInt64Relaxed(&a.v, delta)
}

// AddAcquire atomically adds delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Int64) AddAcquire(delta int64) int64 {
	return arch.AddInt64Acquire(&a.v, delta)
}

// AddRelease atomically adds delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Int64) AddRelease(delta int64) int64 {
	return arch.AddInt64Release(&a.v, delta)
}

// AddAcqRel atomically adds delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) AddAcqRel(delta int64) int64 {
	return arch.AddInt64AcqRel(&a.v, delta)
}

// Sub atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) Sub(delta int64) int64 {
	return arch.AddInt64AcqRel(&a.v, -delta)
}

// SubRelaxed atomically subtracts delta and returns the new value with relaxed ordering.
//
//go:nosplit
func (a *Int64) SubRelaxed(delta int64) int64 {
	return arch.AddInt64Relaxed(&a.v, -delta)
}

// SubAcquire atomically subtracts delta and returns the new value with acquire ordering.
//
//go:nosplit
func (a *Int64) SubAcquire(delta int64) int64 {
	return arch.AddInt64Acquire(&a.v, -delta)
}

// SubRelease atomically subtracts delta and returns the new value with release ordering.
//
//go:nosplit
func (a *Int64) SubRelease(delta int64) int64 {
	return arch.AddInt64Release(&a.v, -delta)
}

// SubAcqRel atomically subtracts delta and returns the new value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) SubAcqRel(delta int64) int64 {
	return arch.AddInt64AcqRel(&a.v, -delta)
}

// And atomically performs bitwise AND and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) And(mask int64) int64 {
	return arch.AndInt64AcqRel(&a.v, mask)
}

// AndRelaxed atomically performs bitwise AND with relaxed ordering.
//
//go:nosplit
func (a *Int64) AndRelaxed(mask int64) int64 {
	return arch.AndInt64Relaxed(&a.v, mask)
}

// AndAcquire atomically performs bitwise AND with acquire ordering.
//
//go:nosplit
func (a *Int64) AndAcquire(mask int64) int64 {
	return arch.AndInt64Acquire(&a.v, mask)
}

// AndRelease atomically performs bitwise AND with release ordering.
//
//go:nosplit
func (a *Int64) AndRelease(mask int64) int64 {
	return arch.AndInt64Release(&a.v, mask)
}

// AndAcqRel atomically performs bitwise AND with acquire-release ordering.
//
//go:nosplit
func (a *Int64) AndAcqRel(mask int64) int64 {
	return arch.AndInt64AcqRel(&a.v, mask)
}

// Or atomically performs bitwise OR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) Or(mask int64) int64 {
	return arch.OrInt64AcqRel(&a.v, mask)
}

// OrRelaxed atomically performs bitwise OR with relaxed ordering.
//
//go:nosplit
func (a *Int64) OrRelaxed(mask int64) int64 {
	return arch.OrInt64Relaxed(&a.v, mask)
}

// OrAcquire atomically performs bitwise OR with acquire ordering.
//
//go:nosplit
func (a *Int64) OrAcquire(mask int64) int64 {
	return arch.OrInt64Acquire(&a.v, mask)
}

// OrRelease atomically performs bitwise OR with release ordering.
//
//go:nosplit
func (a *Int64) OrRelease(mask int64) int64 {
	return arch.OrInt64Release(&a.v, mask)
}

// OrAcqRel atomically performs bitwise OR with acquire-release ordering.
//
//go:nosplit
func (a *Int64) OrAcqRel(mask int64) int64 {
	return arch.OrInt64AcqRel(&a.v, mask)
}

// Xor atomically performs bitwise XOR and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Int64) Xor(mask int64) int64 {
	return arch.XorInt64AcqRel(&a.v, mask)
}

// XorRelaxed atomically performs bitwise XOR with relaxed ordering.
//
//go:nosplit
func (a *Int64) XorRelaxed(mask int64) int64 {
	return arch.XorInt64Relaxed(&a.v, mask)
}

// XorAcquire atomically performs bitwise XOR with acquire ordering.
//
//go:nosplit
func (a *Int64) XorAcquire(mask int64) int64 {
	return arch.XorInt64Acquire(&a.v, mask)
}

// XorRelease atomically performs bitwise XOR with release ordering.
//
//go:nosplit
func (a *Int64) XorRelease(mask int64) int64 {
	return arch.XorInt64Release(&a.v, mask)
}

// XorAcqRel atomically performs bitwise XOR with acquire-release ordering.
//
//go:nosplit
func (a *Int64) XorAcqRel(mask int64) int64 {
	return arch.XorInt64AcqRel(&a.v, mask)
}

// Max atomically stores the maximum of current and val, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int64) Max(val int64) int64 {
	for {
		old := arch.LoadInt64Relaxed(&a.v)
		if old >= val {
			return old
		}
		if arch.CasInt64AcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MaxRelaxed atomically stores the maximum with relaxed ordering.
//
//go:nosplit
func (a *Int64) MaxRelaxed(val int64) int64 {
	for {
		old := arch.LoadInt64Relaxed(&a.v)
		if old >= val {
			return old
		}
		if arch.CasInt64Relaxed(&a.v, old, val) {
			return old
		}
	}
}

// Min atomically stores the minimum of current and val, returning the old value.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Int64) Min(val int64) int64 {
	for {
		old := arch.LoadInt64Relaxed(&a.v)
		if old <= val {
			return old
		}
		if arch.CasInt64AcqRel(&a.v, old, val) {
			return old
		}
	}
}

// MinRelaxed atomically stores the minimum with relaxed ordering.
//
//go:nosplit
func (a *Int64) MinRelaxed(val int64) int64 {
	for {
		old := arch.LoadInt64Relaxed(&a.v)
		if old <= val {
			return old
		}
		if arch.CasInt64Relaxed(&a.v, old, val) {
			return old
		}
	}
}
