// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

// MemoryOrder specifies the memory ordering constraint for atomic operations.
//
// MemoryOrder constants can be used as method receivers to perform atomic
// operations on raw pointers:
//
//	var x int64
//	atomix.Relaxed.StoreInt64(&x, 42)
//	val := atomix.Acquire.LoadInt64(&x)
//
// This is useful for operating on shared memory (e.g., io_uring rings)
// where wrapper types cannot be used.
//
// Alignment: Pointers must be naturally aligned (int32 → 4-byte, int64 → 8-byte).
// Go guarantees alignment for declared variables. Unaligned pointers from unsafe
// operations may cause undefined behavior on some architectures.
type MemoryOrder uint8

const (
	// Relaxed provides no ordering constraints; only atomicity is guaranteed.
	// Use for statistics counters and metrics that don't synchronize with
	// other memory operations.
	Relaxed MemoryOrder = iota

	// Acquire ensures subsequent operations cannot be reordered before the
	// atomic operation. Use when loading a pointer before dereferencing it.
	Acquire

	// Release ensures prior operations cannot be reordered after the atomic
	// operation. Use when publishing data after initialization.
	Release

	// AcqRel combines Acquire and Release semantics. Use for read-modify-write
	// operations in lock-free data structures.
	AcqRel
)
