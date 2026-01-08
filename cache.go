// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "unsafe"

// CacheLineSize is defined in cache_*.go per architecture.

// Int32Padded is an Int32 padded to cache line size.
type Int32Padded struct {
	Int32
	_ [CacheLineSize - 4]byte
}

// Uint32Padded is a Uint32 padded to cache line size.
type Uint32Padded struct {
	Uint32
	_ [CacheLineSize - 4]byte
}

// Int64Padded is an Int64 padded to cache line size.
type Int64Padded struct {
	Int64
	_ [CacheLineSize - 8]byte
}

// Uint64Padded is a Uint64 padded to cache line size.
type Uint64Padded struct {
	Uint64
	_ [CacheLineSize - 8]byte
}

// UintptrPadded is a Uintptr padded to cache line size.
type UintptrPadded struct {
	Uintptr
	_ [CacheLineSize - unsafe.Sizeof(uintptr(0))]byte
}

// BoolPadded is a Bool padded to cache line size.
type BoolPadded struct {
	Bool
	_ [CacheLineSize - 4]byte
}

// Int128Padded is an Int128 padded to cache line size.
type Int128Padded struct {
	Int128
	_ [CacheLineSize - 16]byte
}

// Uint128Padded is a Uint128 padded to cache line size.
type Uint128Padded struct {
	Uint128
	_ [CacheLineSize - 16]byte
}
