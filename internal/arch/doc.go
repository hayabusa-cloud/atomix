// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package arch provides architecture-specific atomic primitives with
// explicit memory ordering.
//
// This is an internal package that implements the low-level atomic operations
// for different CPU architectures. Higher-level types in the parent atomix
// package build on these primitives.
//
// # Architecture Support
//
// The package provides optimized implementations for:
//   - amd64: x86-64 with TSO (Total Store Ordering)
//   - arm64: ARM64 with LSE (Large System Extensions)
//   - riscv64: RISC-V 64-bit with LL/SC atomics
//   - loong64: LoongArch 64-bit with LL/SC atomics
//
// All other architectures fall back to sync/atomic which provides sequential
// consistency (equivalent to AcqRel ordering).
//
// # Memory Ordering
//
// All functions include a memory ordering suffix:
//   - Relaxed: No ordering guarantees, fastest
//   - Acquire: Loads after this see stores before a paired Release
//   - Release: Stores before this are visible after a paired Acquire
//   - AcqRel: Both Acquire and Release semantics (for RMW operations)
//
// # Implementation Strategy
//
// The implementation uses different strategies per architecture:
//
// x86-64 (TSO):
//
//	x86-64's Total Store Ordering provides strong guarantees that make
//	all load/store orderings equivalent. Plain memory access is atomic
//	for aligned values, so Load and Store operations are implemented as
//	pure Go functions for inlining. Only read-modify-write operations
//	(Swap, CAS, Add) require assembly with LOCK prefix or XCHG.
//
// ARM64:
//
//	ARM64 has weak ordering requiring explicit acquire/release instructions.
//	- Relaxed Load/Store: Plain memory access (inlinable pure Go)
//	- Acquire Load: LDAR instruction (assembly)
//	- Release Store: STLR instruction (assembly)
//	- RMW operations: LSE instructions with ordering (assembly)
//
// RISC-V64 / LoongArch64:
//
//	These architectures have weak ordering and all atomic operations require
//	explicit fence instructions. All operations are implemented in assembly:
//	- RISC-V: Uses AMO instructions with .aq/.rl suffixes, or LR/SC loops
//	- LoongArch: Uses AM* instructions, or LL/SC loops with DBAR barriers
//
// # 128-bit Atomics
//
// 128-bit operations are available on all supported architectures:
//   - amd64: CMPXCHG16B instruction (requires 16-byte alignment)
//   - arm64: LDXP/STXP (default) or CASP with -tags=lse2 for ARMv8.4+
//   - riscv64/loong64: Emulated via LL/SC on low 64 bits with full barrier
//
// # Inlining Optimization
//
// To enable inlining for hot-path operations, load/store functions are
// split into separate files:
//   - loadstore_amd64.go: Pure Go implementations for x86-64 TSO
//   - loadstore_arm64.go: Pure Go relaxed implementations for ARM64
//
// Assembly stubs are only used where hardware instructions with ordering
// are required, marked with //go:noescape to prevent escape analysis
// overhead.
//
// # Function Naming Convention
//
//	<Op><Type><Ordering>
//
// Where:
//   - Op: Load, Store, Swap, Cas, Cax, Add, And, Or, Xor
//   - Type: Int32, Uint32, Int64, Uint64, Uintptr, Pointer, Uint128
//   - Ordering: Relaxed, Acquire, Release, AcqRel
//
// Cas returns bool (success), Cax returns old value (compare-exchange).
// Add returns the new value. Swap/And/Or/Xor return the previous value.
package arch
