// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "unsafe"

// Allocator is a sequential allocator for placing atomic values in a buffer.
// It provides methods to allocate properly aligned atomic types.
type Allocator struct {
	buf []byte
	off int
}

// NewAllocator creates a new Allocator backed by the given buffer.
// The caller must keep buf reachable for the lifetime of all allocated values.
func NewAllocator(buf []byte) *Allocator {
	return &Allocator{buf: buf}
}

// Offset returns the current allocation offset.
func (a *Allocator) Offset() int {
	return a.off
}

// Remaining returns the number of bytes remaining in the buffer.
func (a *Allocator) Remaining() int {
	return len(a.buf) - a.off
}

// Reset resets the allocator to the beginning of the buffer.
// Previously allocated values remain valid but may be overwritten.
func (a *Allocator) Reset() {
	a.off = 0
}

// Skip advances the offset by n bytes.
func (a *Allocator) Skip(n int) {
	if a.off+n > len(a.buf) {
		panic("atomix: skip beyond buffer")
	}
	a.off += n
}

// Align advances the offset to the next multiple of alignment.
func (a *Allocator) Align(alignment int) {
	if alignment <= 0 || (alignment&(alignment-1)) != 0 {
		panic("atomix: alignment must be a power of 2")
	}
	addr := uintptr(unsafe.Pointer(&a.buf[0])) + uintptr(a.off)
	pad := int((uintptr(alignment) - (addr & uintptr(alignment-1))) & uintptr(alignment-1))
	if a.off+pad > len(a.buf) {
		panic("atomix: align beyond buffer")
	}
	a.off += pad
}

// Int32 allocates and returns an Int32 at a 4-byte aligned address.
func (a *Allocator) Int32() *Int32 {
	n, ptr := PlaceAlignedInt32(a.buf, a.off)
	a.off += n
	return ptr
}

// Uint32 allocates and returns a Uint32 at a 4-byte aligned address.
func (a *Allocator) Uint32() *Uint32 {
	n, ptr := PlaceAlignedUint32(a.buf, a.off)
	a.off += n
	return ptr
}

// Bool allocates and returns a Bool at a 4-byte aligned address.
func (a *Allocator) Bool() *Bool {
	n, ptr := PlaceAlignedBool(a.buf, a.off)
	a.off += n
	return ptr
}

// Int64 allocates and returns an Int64 at an 8-byte aligned address.
func (a *Allocator) Int64() *Int64 {
	n, ptr := PlaceAlignedInt64(a.buf, a.off)
	a.off += n
	return ptr
}

// Uint64 allocates and returns a Uint64 at an 8-byte aligned address.
func (a *Allocator) Uint64() *Uint64 {
	n, ptr := PlaceAlignedUint64(a.buf, a.off)
	a.off += n
	return ptr
}

// Uintptr allocates and returns a Uintptr at a pointer-aligned address.
func (a *Allocator) Uintptr() *Uintptr {
	n, ptr := PlaceAlignedUintptr(a.buf, a.off)
	a.off += n
	return ptr
}

// Int128 allocates and returns an Int128 at a 16-byte aligned address.
func (a *Allocator) Int128() *Int128 {
	n, ptr := PlaceAlignedInt128(a.buf, a.off)
	a.off += n
	return ptr
}

// Uint128 allocates and returns a Uint128 at a 16-byte aligned address.
func (a *Allocator) Uint128() *Uint128 {
	n, ptr := PlaceAlignedUint128(a.buf, a.off)
	a.off += n
	return ptr
}

// CacheAlignedInt32 allocates and returns an Int32Padded at a cache-line aligned address.
func (a *Allocator) CacheAlignedInt32() *Int32Padded {
	n, ptr := PlaceCacheAlignedInt32(a.buf, a.off)
	a.off += n
	return ptr
}

// CacheAlignedUint32 allocates and returns a Uint32Padded at a cache-line aligned address.
func (a *Allocator) CacheAlignedUint32() *Uint32Padded {
	n, ptr := PlaceCacheAlignedUint32(a.buf, a.off)
	a.off += n
	return ptr
}

// CacheAlignedInt64 allocates and returns an Int64Padded at a cache-line aligned address.
func (a *Allocator) CacheAlignedInt64() *Int64Padded {
	n, ptr := PlaceCacheAlignedInt64(a.buf, a.off)
	a.off += n
	return ptr
}

// CacheAlignedUint64 allocates and returns a Uint64Padded at a cache-line aligned address.
func (a *Allocator) CacheAlignedUint64() *Uint64Padded {
	n, ptr := PlaceCacheAlignedUint64(a.buf, a.off)
	a.off += n
	return ptr
}

// CacheAlignedUintptr allocates and returns a UintptrPadded at a cache-line aligned address.
func (a *Allocator) CacheAlignedUintptr() *UintptrPadded {
	n, ptr := PlaceCacheAlignedUintptr(a.buf, a.off)
	a.off += n
	return ptr
}

// CacheAlignedBool allocates and returns a BoolPadded at a cache-line aligned address.
func (a *Allocator) CacheAlignedBool() *BoolPadded {
	n, ptr := PlaceCacheAlignedBool(a.buf, a.off)
	a.off += n
	return ptr
}

// CacheAlignedInt128 allocates and returns an Int128Padded at a cache-line aligned address.
func (a *Allocator) CacheAlignedInt128() *Int128Padded {
	n, ptr := PlaceCacheAlignedInt128(a.buf, a.off)
	a.off += n
	return ptr
}

// CacheAlignedUint128 allocates and returns a Uint128Padded at a cache-line aligned address.
func (a *Allocator) CacheAlignedUint128() *Uint128Padded {
	n, ptr := PlaceCacheAlignedUint128(a.buf, a.off)
	a.off += n
	return ptr
}
