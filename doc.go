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
//   - [Release]: Orders prior accesses before the store (prior writes for default arm64 128-bit stores)
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
//	val := counter.Add(1)           // Default RMW ordering (AcqRel)
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
// Unknown orderings use the fallback orderings Load→Acquire, Store→Release,
// and RMW→AcqRel.
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
// Wrapper and pointer-based operations include Load, Store, Swap,
// CompareAndSwap, CompareExchange, arithmetic, bitwise, min/max, and
// increment/decrement operations where defined for the value kind.
//
// Default methods use: Load=Relaxed, Store=Relaxed, RMW=AcqRel.
// Note: sync/atomic operations are sequentially consistent. atomix defaults to
// Relaxed for Load and Store; use LoadAcquire/StoreRelease when an
// acquire/release synchronization edge is required.
//
// Return value semantics:
//   - Wrapper Add/Sub/Inc/Dec return the NEW value (after the operation)
//   - Pointer-based 32/64/uintptr Add/Sub return the NEW value
//   - Pointer-based AddInt128/AddUint128 return the OLD value
//   - Swap/And/Or/Xor/Max/Min return the OLD value (before the operation)
//
// # Platform Support
//
// Primary (native atomic instructions):
//   - amd64: LOCK-prefixed instructions; TSO provides acquire/release
//   - arm64: LSE for 32/64-bit atomics and CASP, LSE2 for atomic pair loads and stores
//
// Secondary (native with limitations):
//   - riscv64: AMO instructions with .aq/.rl suffixes
//   - loong64: AM*_DB instructions
//
// Fallback: Other architectures use sync/atomic (over-synchronized).
//
// # ARM64 128-bit Build Options
//
// Default arm64 builds target FEAT_LSE2. Pair loads and stores use LDP/STP
// on 16-byte-aligned Normal memory with Inner and Outer Write-Back cacheability,
// as on the Go heap. Compare-and-swap uses LSE CASP with the requested ordering.
// A 128-bit release store publishes earlier writes with DMB ISHST. To order an
// earlier read before the store, use an acquire load for that read.
//
// On older arm64 CPUs with LSE, build with -tags=atomix_llsc. This backend
// requires writable storage because loads and failed compare-exchanges can
// write back the observed pair.
//
// # 128-bit Atomics
//
// [Int128] and [Uint128] require 16-byte alignment. Use [PlaceAlignedInt128]
// or [PlaceAlignedUint128] to ensure proper alignment.
//
// True 128-bit atomicity is only available on:
//   - amd64: LOCK CMPXCHG16B
//   - arm64: LDP/STP and CASP by default, or exclusive pairs with -tags=atomix_llsc
//
// riscv64 and loong64 emulate 128-bit operations through low-word LR/SC
// (riscv64) or LL/SC (loong64) and may exhibit torn reads. Generic fallback
// 128-bit operations are not concurrency-safe without external synchronization.
// On riscv64 and loong64, 128-bit SwapAcquire/SwapRelease/SwapAcqRel
// alias the relaxed low-word swap path; use separate 32/64-bit
// synchronization or external synchronization when acquire/release
// publication is required.
//
// The wrapper 128-bit Add/Sub/Inc/Dec APIs return the new value.
// The pointer-based [MemoryOrder.AddInt128] and [MemoryOrder.AddUint128]
// APIs return the old value.
//
// # Placement Helpers
//
// For embedding atomics in shared memory or custom allocators:
//   - [CanPlaceAligned4], [CanPlaceAligned8], [CanPlaceAligned16]
//   - [PlaceAlignedInt32], [PlaceAlignedInt64], [PlaceAlignedUint128], etc.
//   - [Allocator]: Sequential allocator for building atomic structures
package atomix
