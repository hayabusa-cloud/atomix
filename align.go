// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "unsafe"

// CanPlaceAligned4 reports whether a 4-byte aligned value can be placed
// in p starting at offset off. Returns true if there is sufficient space
// for alignment padding plus 4 bytes.
func CanPlaceAligned4(p []byte, off int) bool {
	if off < 0 || off > len(p) {
		return false
	}
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := (4 - (addr & 3)) & 3
	return off+int(pad)+4 <= len(p)
}

// CanPlaceAligned8 reports whether an 8-byte aligned value can be placed
// in p starting at offset off. Returns true if there is sufficient space
// for alignment padding plus 8 bytes.
func CanPlaceAligned8(p []byte, off int) bool {
	if off < 0 || off > len(p) {
		return false
	}
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := (8 - (addr & 7)) & 7
	return off+int(pad)+8 <= len(p)
}

// CanPlaceAligned16 reports whether a 16-byte aligned value can be placed
// in p starting at offset off. Returns true if there is sufficient space
// for alignment padding plus 16 bytes.
func CanPlaceAligned16(p []byte, off int) bool {
	if off < 0 || off > len(p) {
		return false
	}
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := (16 - (addr & 15)) & 15
	return off+int(pad)+16 <= len(p)
}

// CanPlaceCacheAligned reports whether a cache-line aligned value of the
// given size can be placed in p starting at offset off.
func CanPlaceCacheAligned(p []byte, off, size int) bool {
	if off < 0 || off > len(p) {
		return false
	}
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := (CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1)
	return off+pad+size <= len(p)
}

// PlaceAlignedInt32 places an Int32 at a 4-byte aligned address in p.
// Returns the number of bytes consumed (padding + 4) and a pointer to the Int32.
// Panics if there is insufficient space.
//
//go:nocheckptr
func PlaceAlignedInt32(p []byte, off int) (n int, a *Int32) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((4 - (addr & 3)) & 3)
	n = pad + 4
	if off+n > len(p) {
		panic("atomix: insufficient space for Int32")
	}
	a = (*Int32)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceAlignedUint32 places a Uint32 at a 4-byte aligned address in p.
//
//go:nocheckptr
func PlaceAlignedUint32(p []byte, off int) (n int, a *Uint32) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((4 - (addr & 3)) & 3)
	n = pad + 4
	if off+n > len(p) {
		panic("atomix: insufficient space for Uint32")
	}
	a = (*Uint32)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceAlignedBool places a Bool at a 4-byte aligned address in p.
//
//go:nocheckptr
func PlaceAlignedBool(p []byte, off int) (n int, a *Bool) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((4 - (addr & 3)) & 3)
	n = pad + 4
	if off+n > len(p) {
		panic("atomix: insufficient space for Bool")
	}
	a = (*Bool)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceAlignedInt64 places an Int64 at an 8-byte aligned address in p.
//
//go:nocheckptr
func PlaceAlignedInt64(p []byte, off int) (n int, a *Int64) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((8 - (addr & 7)) & 7)
	n = pad + 8
	if off+n > len(p) {
		panic("atomix: insufficient space for Int64")
	}
	a = (*Int64)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceAlignedUint64 places a Uint64 at an 8-byte aligned address in p.
//
//go:nocheckptr
func PlaceAlignedUint64(p []byte, off int) (n int, a *Uint64) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((8 - (addr & 7)) & 7)
	n = pad + 8
	if off+n > len(p) {
		panic("atomix: insufficient space for Uint64")
	}
	a = (*Uint64)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceAlignedUintptr places a Uintptr at a pointer-aligned address in p.
//
//go:nocheckptr
func PlaceAlignedUintptr(p []byte, off int) (n int, a *Uintptr) {
	const size = int(unsafe.Sizeof(uintptr(0)))
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((uintptr(size) - (addr & uintptr(size-1))) & uintptr(size-1))
	n = pad + size
	if off+n > len(p) {
		panic("atomix: insufficient space for Uintptr")
	}
	a = (*Uintptr)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceAlignedInt128 places an Int128 at a 16-byte aligned address in p.
//
//go:nocheckptr
func PlaceAlignedInt128(p []byte, off int) (n int, a *Int128) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((16 - (addr & 15)) & 15)
	n = pad + 16
	if off+n > len(p) {
		panic("atomix: insufficient space for Int128")
	}
	a = (*Int128)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceAlignedUint128 places a Uint128 at a 16-byte aligned address in p.
//
//go:nocheckptr
func PlaceAlignedUint128(p []byte, off int) (n int, a *Uint128) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((16 - (addr & 15)) & 15)
	n = pad + 16
	if off+n > len(p) {
		panic("atomix: insufficient space for Uint128")
	}
	a = (*Uint128)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceCacheAlignedInt32 places an Int32Padded at a cache-line aligned address.
//
//go:nocheckptr
func PlaceCacheAlignedInt32(p []byte, off int) (n int, a *Int32Padded) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1))
	n = pad + CacheLineSize
	if off+n > len(p) {
		panic("atomix: insufficient space for Int32Padded")
	}
	a = (*Int32Padded)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceCacheAlignedUint32 places a Uint32Padded at a cache-line aligned address.
//
//go:nocheckptr
func PlaceCacheAlignedUint32(p []byte, off int) (n int, a *Uint32Padded) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1))
	n = pad + CacheLineSize
	if off+n > len(p) {
		panic("atomix: insufficient space for Uint32Padded")
	}
	a = (*Uint32Padded)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceCacheAlignedInt64 places an Int64Padded at a cache-line aligned address.
//
//go:nocheckptr
func PlaceCacheAlignedInt64(p []byte, off int) (n int, a *Int64Padded) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1))
	n = pad + CacheLineSize
	if off+n > len(p) {
		panic("atomix: insufficient space for Int64Padded")
	}
	a = (*Int64Padded)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceCacheAlignedUint64 places a Uint64Padded at a cache-line aligned address.
//
//go:nocheckptr
func PlaceCacheAlignedUint64(p []byte, off int) (n int, a *Uint64Padded) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1))
	n = pad + CacheLineSize
	if off+n > len(p) {
		panic("atomix: insufficient space for Uint64Padded")
	}
	a = (*Uint64Padded)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceCacheAlignedUintptr places a UintptrPadded at a cache-line aligned address.
//
//go:nocheckptr
func PlaceCacheAlignedUintptr(p []byte, off int) (n int, a *UintptrPadded) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1))
	n = pad + CacheLineSize
	if off+n > len(p) {
		panic("atomix: insufficient space for UintptrPadded")
	}
	a = (*UintptrPadded)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceCacheAlignedBool places a BoolPadded at a cache-line aligned address.
//
//go:nocheckptr
func PlaceCacheAlignedBool(p []byte, off int) (n int, a *BoolPadded) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1))
	n = pad + CacheLineSize
	if off+n > len(p) {
		panic("atomix: insufficient space for BoolPadded")
	}
	a = (*BoolPadded)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceCacheAlignedInt128 places an Int128Padded at a cache-line aligned address.
//
//go:nocheckptr
func PlaceCacheAlignedInt128(p []byte, off int) (n int, a *Int128Padded) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1))
	n = pad + CacheLineSize
	if off+n > len(p) {
		panic("atomix: insufficient space for Int128Padded")
	}
	a = (*Int128Padded)(unsafe.Pointer(&p[off+pad]))
	return n, a
}

// PlaceCacheAlignedUint128 places a Uint128Padded at a cache-line aligned address.
//
//go:nocheckptr
func PlaceCacheAlignedUint128(p []byte, off int) (n int, a *Uint128Padded) {
	addr := uintptr(unsafe.Pointer(&p[0])) + uintptr(off)
	pad := int((CacheLineSize - int(addr&(CacheLineSize-1))) & (CacheLineSize - 1))
	n = pad + CacheLineSize
	if off+n > len(p) {
		panic("atomix: insufficient space for Uint128Padded")
	}
	a = (*Uint128Padded)(unsafe.Pointer(&p[off+pad]))
	return n, a
}
