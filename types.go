// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "unsafe"

// Bool represents an atomic boolean.
//
// The zero value is false. Bool is safe for concurrent use.
// Must not be copied after first use.
type Bool struct {
	_ noCopy
	v uint32
}

// Int32 represents an atomic 32-bit signed integer.
//
// The zero value is 0. Int32 is safe for concurrent use.
// Must not be copied after first use.
type Int32 struct {
	_ noCopy
	v int32
}

// Uint32 represents an atomic 32-bit unsigned integer.
//
// The zero value is 0. Uint32 is safe for concurrent use.
// Must not be copied after first use.
type Uint32 struct {
	_ noCopy
	v uint32
}

// Int64 represents an atomic 64-bit signed integer.
//
// The zero value is 0. Int64 is safe for concurrent use.
// Must not be copied after first use.
type Int64 struct {
	_ noCopy
	v int64
}

// Uint64 represents an atomic 64-bit unsigned integer.
//
// The zero value is 0. Uint64 is safe for concurrent use.
// Must not be copied after first use.
type Uint64 struct {
	_ noCopy
	v uint64
}

// Uintptr represents an atomic pointer-sized unsigned integer.
//
// The zero value is 0. Uintptr is safe for concurrent use.
// Must not be copied after first use.
type Uintptr struct {
	_ noCopy
	v uintptr
}

// Pointer represents an atomic pointer to a value of type T.
//
// The zero value is nil. Pointer is safe for concurrent use.
// Must not be copied after first use.
type Pointer[T any] struct {
	_ noCopy
	v unsafe.Pointer
}

// Int128 represents an atomic 128-bit signed integer.
//
// The zero value is 0. Int128 is safe for concurrent use.
// Must not be copied after first use.
//
// Int128 requires 16-byte alignment. Use [PlaceAlignedInt128] to ensure
// proper alignment when embedding in byte slices or shared memory.
type Int128 struct {
	_ noCopy
	v [16]byte
}

// Uint128 represents an atomic 128-bit unsigned integer.
//
// The zero value is 0. Uint128 is safe for concurrent use.
// Must not be copied after first use.
//
// Uint128 requires 16-byte alignment. Use [PlaceAlignedUint128] to ensure
// proper alignment when embedding in byte slices or shared memory.
type Uint128 struct {
	_ noCopy
	v [16]byte
}
