// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

import "unsafe"

// Load atomically loads and returns the pointer with relaxed ordering.
//
//go:nosplit
func (a *Pointer[T]) Load() *T {
	return (*T)(arch.LoadPointerRelaxed(&a.v))
}

// LoadRelaxed atomically loads and returns the pointer with relaxed ordering.
//
//go:nosplit
func (a *Pointer[T]) LoadRelaxed() *T {
	return (*T)(arch.LoadPointerRelaxed(&a.v))
}

// LoadAcquire atomically loads and returns the pointer with acquire ordering.
//
//go:nosplit
func (a *Pointer[T]) LoadAcquire() *T {
	return (*T)(arch.LoadPointerAcquire(&a.v))
}

// Store atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Pointer[T]) Store(val *T) {
	arch.StorePointerRelaxed(&a.v, unsafe.Pointer(val))
}

// StoreRelaxed atomically stores val with relaxed ordering.
//
//go:nosplit
func (a *Pointer[T]) StoreRelaxed(val *T) {
	arch.StorePointerRelaxed(&a.v, unsafe.Pointer(val))
}

// StoreRelease atomically stores val with release ordering.
//
//go:nosplit
func (a *Pointer[T]) StoreRelease(val *T) {
	arch.StorePointerRelease(&a.v, unsafe.Pointer(val))
}

// Swap atomically swaps the pointer and returns the old value with acquire-release ordering.
//
//go:nosplit
func (a *Pointer[T]) Swap(new *T) *T {
	return (*T)(arch.SwapPointerAcqRel(&a.v, unsafe.Pointer(new)))
}

// SwapRelaxed atomically swaps the pointer with relaxed ordering.
//
//go:nosplit
func (a *Pointer[T]) SwapRelaxed(new *T) *T {
	return (*T)(arch.SwapPointerRelaxed(&a.v, unsafe.Pointer(new)))
}

// SwapAcquire atomically swaps the pointer with acquire ordering.
//
//go:nosplit
func (a *Pointer[T]) SwapAcquire(new *T) *T {
	return (*T)(arch.SwapPointerAcquire(&a.v, unsafe.Pointer(new)))
}

// SwapRelease atomically swaps the pointer with release ordering.
//
//go:nosplit
func (a *Pointer[T]) SwapRelease(new *T) *T {
	return (*T)(arch.SwapPointerRelease(&a.v, unsafe.Pointer(new)))
}

// SwapAcqRel atomically swaps the pointer with acquire-release ordering.
//
//go:nosplit
func (a *Pointer[T]) SwapAcqRel(new *T) *T {
	return (*T)(arch.SwapPointerAcqRel(&a.v, unsafe.Pointer(new)))
}

// CompareAndSwap atomically compares and swaps with acquire-release ordering.
// Returns true if the swap was performed.
//
//go:nosplit
func (a *Pointer[T]) CompareAndSwap(old, new *T) bool {
	return arch.CasPointerAcqRel(&a.v, unsafe.Pointer(old), unsafe.Pointer(new))
}

// CompareAndSwapRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareAndSwapRelaxed(old, new *T) bool {
	return arch.CasPointerRelaxed(&a.v, unsafe.Pointer(old), unsafe.Pointer(new))
}

// CompareAndSwapAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareAndSwapAcquire(old, new *T) bool {
	return arch.CasPointerAcquire(&a.v, unsafe.Pointer(old), unsafe.Pointer(new))
}

// CompareAndSwapRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareAndSwapRelease(old, new *T) bool {
	return arch.CasPointerRelease(&a.v, unsafe.Pointer(old), unsafe.Pointer(new))
}

// CompareAndSwapAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareAndSwapAcqRel(old, new *T) bool {
	return arch.CasPointerAcqRel(&a.v, unsafe.Pointer(old), unsafe.Pointer(new))
}

// CompareExchange atomically compares and swaps, returning the old pointer.
// Uses acquire-release ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareExchange(old, new *T) *T {
	return (*T)(arch.CaxPointerAcqRel(&a.v, unsafe.Pointer(old), unsafe.Pointer(new)))
}

// CompareExchangeRelaxed atomically compares and swaps with relaxed ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareExchangeRelaxed(old, new *T) *T {
	return (*T)(arch.CaxPointerRelaxed(&a.v, unsafe.Pointer(old), unsafe.Pointer(new)))
}

// CompareExchangeAcquire atomically compares and swaps with acquire ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareExchangeAcquire(old, new *T) *T {
	return (*T)(arch.CaxPointerAcquire(&a.v, unsafe.Pointer(old), unsafe.Pointer(new)))
}

// CompareExchangeRelease atomically compares and swaps with release ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareExchangeRelease(old, new *T) *T {
	return (*T)(arch.CaxPointerRelease(&a.v, unsafe.Pointer(old), unsafe.Pointer(new)))
}

// CompareExchangeAcqRel atomically compares and swaps with acquire-release ordering.
//
//go:nosplit
func (a *Pointer[T]) CompareExchangeAcqRel(old, new *T) *T {
	return (*T)(arch.CaxPointerAcqRel(&a.v, unsafe.Pointer(old), unsafe.Pointer(new)))
}
