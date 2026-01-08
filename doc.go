// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package atomix provides atomic primitives with explicit memory ordering.
//
// # Architecture
//
// atomix is a foundation layer with no external dependencies.
// It provides atomic operations for the I/O stack, used by higher-level
// modules like uring and sox.
//
// # Memory Ordering
//
// The package exposes four memory orderings:
//
//   - [Relaxed]: Only atomicity guaranteed; no ordering constraints
//   - [Acquire]: Subsequent operations cannot reorder before the load
//   - [Release]: Prior operations cannot reorder after the store
//   - [AcqRel]: Acquire + Release; for read-modify-write operations
//
// Unlike sync/atomic which provides sequential consistency, this package
// allows choosing the minimal ordering required for weakly-ordered
// architectures (ARM, RISC-V).
//
// # Two APIs
//
// Type-based API for embedding in structs:
//
//	var counter atomix.Int64
//	counter.Store(0)
//	val := counter.Add(1)           // AcqRel ordering (safe default)
//	val = counter.AddRelaxed(1)     // Explicit relaxed ordering
//
// Pointer-based API for raw memory (shared memory, io_uring):
//
//	var flags int32
//	atomix.Release.StoreInt32(&flags, 1)
//	val := atomix.Acquire.LoadInt32(&flags)
//	atomix.Relaxed.CompareAndSwapInt32(&flags, 0, 1)
//
// The pointer-based API uses [MemoryOrder] constants as method receivers.
// Unknown orderings fall back to safe defaults (Load→Acquire, Store→Release,
// RMW→AcqRel).
//
// # Types
//
// Core atomic types:
//   - [Bool]: Atomic boolean (backed by uint32)
//   - [Int32], [Uint32]: 32-bit integers
//   - [Int64], [Uint64]: 64-bit integers
//   - [Uintptr]: Pointer-sized integer
//   - [Pointer]: Generic atomic pointer
//   - [Int128], [Uint128]: 128-bit integers (requires 16-byte alignment)
//
// Cache-line padded variants prevent false sharing:
//   - [Int32Padded], [Uint32Padded], [Int64Padded], [Uint64Padded]
//   - [UintptrPadded], [BoolPadded], [Int128Padded], [Uint128Padded]
//
// All types are safe for concurrent use. The zero value is valid (0 or nil).
//
// # Operations
//
// All types support Load, Store, Swap, CompareAndSwap, CompareExchange,
// Add, Sub, And, Or, Xor, Max, Min, Inc, Dec with explicit ordering suffixes.
//
// Default methods use: Load=Relaxed, Store=Relaxed, RMW=AcqRel.
// Note: sync/atomic uses acquire for Load and release for Store.
// Use LoadAcquire/StoreRelease for sync/atomic-equivalent ordering.
//
// Return value semantics match sync/atomic:
//   - Add/Sub/Inc/Dec return the NEW value (after the operation)
//   - Swap/And/Or/Xor/Max/Min return the OLD value (before the operation)
//
// # Platform Support
//
// Primary (native atomic instructions):
//   - amd64: LOCK-prefixed instructions; TSO provides acquire/release
//   - arm64: LSE atomics (ARMv8.1+) for 32/64-bit; LL/SC or CASP for 128-bit
//
// Secondary (native with limitations):
//   - riscv64: AMO instructions with .aq/.rl suffixes
//   - loong64: AM*_DB instructions
//
// Fallback: Other architectures use sync/atomic (over-synchronized).
//
// # ARM64 128-bit Build Options
//
// ARM64 128-bit atomics support two implementations via build tags:
//
//   - Default (!lse2): LL/SC using LDXP/STXP instructions
//   - -tags=lse2: CASP instruction (LSE2)
//
// LL/SC has lower microarchitectural overhead than CASP in uncontended cases.
// Use -tags=lse2 for ARMv8.4+ hardware with high-contention workloads.
//
// # 128-bit Atomics
//
// [Int128] and [Uint128] require 16-byte alignment. Use [PlaceAlignedInt128]
// or [PlaceAlignedUint128] to ensure proper alignment.
//
// True 128-bit atomicity is only available on:
//   - amd64: LOCK CMPXCHG16B
//   - arm64: LDXP/STXP (default) or CASP (-tags=lse2)
//
// Other architectures provide mutual exclusion but may exhibit torn reads.
//
// # Placement Helpers
//
// For embedding atomics in shared memory or custom allocators:
//   - [CanPlaceAligned4], [CanPlaceAligned8], [CanPlaceAligned16]
//   - [PlaceAlignedInt32], [PlaceAlignedInt64], [PlaceAlignedUint128], etc.
//   - [Allocator]: Sequential allocator for building atomic structures
package atomix
