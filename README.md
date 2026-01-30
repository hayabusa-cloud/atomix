# atomix

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/atomix.svg)](https://pkg.go.dev/code.hybscloud.com/atomix)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/atomix)](https://goreportcard.com/report/github.com/hayabusa-cloud/atomix)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/atomix/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/atomix)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Languages:** English | [简体中文](README.zh-CN.md) | [日本語](README.ja.md) | [Español](README.es.md) | [Français](README.fr.md)

Atomic operations with explicit memory ordering for Go.

## Overview

Go's `sync/atomic` provides atomic operations with sequential consistency. This library exposes C++11/C11 memory model orderings (Relaxed, Acquire, Release, AcqRel) through architecture-specific implementations.

```go
import "code.hybscloud.com/atomix"

var counter atomix.Int64

// Method-based API with ordering suffix
counter.AddRelaxed(1)    // Relaxed: no synchronization
counter.Add(1)           // AcqRel: default safe ordering

// Pointer-based API for raw memory
var flags int32
atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
```

## Installation

```bash
go get code.hybscloud.com/atomix
```

**Requirements:** Go 1.25+

## Memory Ordering

The library implements four orderings from the C++11 memory model:

| Ordering | Semantics |
|----------|-----------|
| **Relaxed** | Atomicity only. No synchronization or ordering constraints. |
| **Acquire** | Subsequent reads/writes cannot be reordered before this load. Pairs with Release stores. |
| **Release** | Prior reads/writes cannot be reordered after this store. Pairs with Acquire loads. |
| **AcqRel** | Combines Acquire and Release semantics. For read-modify-write operations. |

### Ordering Selection

Default methods (no ordering suffix) use:
- Load operations: Relaxed
- Store operations: Relaxed
- Read-modify-write operations: AcqRel

**Note:** sync/atomic uses acquire for Load and release for Store (sequential consistency on x86). atomix defaults to Relaxed, which maps to different instructions on weakly-ordered architectures (e.g., LDR vs LDAR on ARM64). Use `LoadAcquire`/`StoreRelease` when sync/atomic-equivalent ordering is required.

### When to Use Each Ordering

| Use Case | Ordering | Rationale |
|----------|----------|-----------|
| Statistics counters | Relaxed | No synchronization needed; eventual consistency acceptable |
| Reference counting | AcqRel | Ensures visibility of object state before deallocation |
| Producer-consumer flags | Release/Acquire | Producer releases data, consumer acquires |
| Spinlock acquire | Acquire | Critical section reads must see prior writes |
| Spinlock release | Release | Critical section writes must complete before unlock |
| Sequence locks | AcqRel | Both directions need ordering |

## Types

### Value Types

| Type | Size | Description |
|------|------|-------------|
| `Bool` | 4 bytes | Atomic boolean (backed by uint32) |
| `Int32`, `Uint32` | 4 bytes | 32-bit integers |
| `Int64`, `Uint64` | 8 bytes | 64-bit integers |
| `Uintptr` | 8 bytes | Pointer-sized integer |
| `Pointer[T]` | 8 bytes | Generic atomic pointer |
| `Int128`, `Uint128` | 16 bytes | 128-bit integers (requires 16-byte alignment) |

### Padded Types

Padded variants (`Int64Padded`, `Uint64Padded`, etc.) occupy a full cache line (64 bytes) to prevent false sharing when multiple atomic variables are accessed by different CPU cores.

```go
// Without padding: variables may share cache line, causing contention
var a, b atomix.Int64  // May be adjacent in memory

// With padding: each variable occupies its own cache line
var a, b atomix.Int64Padded  // 64-byte separation guaranteed
```

## Operations

| Operation | Returns | Description |
|-----------|---------|-------------|
| `Load` | value | Atomic read |
| `Store` | — | Atomic write |
| `Swap` | old value | Atomic exchange |
| `CompareAndSwap` | bool | Returns true if exchange occurred |
| `CompareExchange` | old value | Returns previous value regardless of success |
| `Add`, `Sub` | new value | Atomic arithmetic |
| `Inc`, `Dec` | new value | Atomic increment/decrement by 1 |
| `And`, `Or`, `Xor` | old value | Atomic bitwise operations |
| `Max`, `Min` | old value | Atomic maximum/minimum |

**Return value semantics:** Add/Sub/Inc/Dec return the **new** value (like sync/atomic). Swap/And/Or/Xor/Max/Min return the **old** value.

### CompareAndSwap vs CompareExchange

```go
// CompareAndSwap: returns success/failure
if v.CompareAndSwap(old, new) {
    // Success
}

// CompareExchange: returns previous value (enables CAS loops without separate Load)
for {
    old := v.Load()
    new := transform(old)
    if v.CompareExchange(old, new) == old {
        break  // Success
    }
}
```

## Pointer-Based API

For interoperation with memory-mapped regions, shared memory, or io_uring rings:

```go
var flags int32

atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
atomix.Release.CompareAndSwapInt32(&flags, 0, 1)
```

The pointer-based API operates on raw `*int32`, `*int64`, etc., rather than wrapper types. This is useful when atomic variables cannot use wrapper types (e.g., fields in kernel-shared structures).

## 128-bit Operations

128-bit atomics require 16-byte alignment. Use placement helpers for shared memory:

```go
buf := make([]byte, 32)
_, ptr := atomix.PlaceAlignedUint128(buf, 0)
ptr.Store(lo, hi)

var v atomix.Uint128  // Type ensures alignment
v.Store(lo, hi)
```

| Architecture | 128-bit Implementation |
|--------------|------------------------|
| amd64 | `LOCK CMPXCHG16B` |
| arm64 | `LDXP/STXP` (default) or `CASP` (`-tags=lse2`) |
| riscv64, loong64 | Spinlock emulation (LL/SC on low 64 bits) |

**Note:** 128-bit atomics are primarily useful for double-word CAS patterns (e.g., lock-free data structures with version counters).

## Architecture Implementation

### x86-64 (TSO)

x86-64 provides Total Store Ordering (TSO), a strong memory model where:
- All loads have implicit acquire semantics
- All stores have implicit release semantics
- Store-load ordering requires explicit barrier (MFENCE) or locked instruction

Consequently, all ordering variants compile to identical machine code on x86-64. The primary benefit of explicit ordering on x86-64 is documentation and portability.

| Operation | Instruction | Notes |
|-----------|-------------|-------|
| Load | `MOV` | Plain memory access |
| Store | `MOV` | Plain memory access |
| Add | `LOCK XADD` | Returns old value |
| Swap | `XCHG` | Implicit LOCK |
| CAS | `LOCK CMPXCHG` | |
| And/Or/Xor | `LOCK CMPXCHG` loop | Returns old value via CAS loop |
| CAS128 | `LOCK CMPXCHG16B` | |

Load and Store are implemented in pure Go for compiler inlining.

### ARM64 (Weakly Ordered)

ARM64 has a weakly ordered memory model requiring explicit ordering instructions. LSE (Large System Extensions) provides atomic instructions with ordering suffixes:

**Suffix meanings:** No suffix = Relaxed, `A` = Acquire, `L` = Release, `AL` = Acquire-Release

| Operation | Relaxed | Acquire | Release | AcqRel |
|-----------|---------|---------|---------|--------|
| Load | `LDR` | `LDAR` | — | — |
| Store | `STR` | — | `STLR` | — |
| Add | `LDADD` | `LDADDA` | `LDADDL` | `LDADDAL` |
| CAS | `CAS` | `CASA` | `CASL` | `CASAL` |
| Swap | `SWP` | `SWPA` | `SWPL` | `SWPAL` |
| And | `LDCLR`† | `LDCLRA` | `LDCLRL` | `LDCLRAL` |
| Or | `LDSET` | `LDSETA` | `LDSETL` | `LDSETAL` |
| Xor | `LDEOR` | `LDEORA` | `LDEORL` | `LDEORAL` |

† `LDCLR` clears bits (AND with complement). To implement `And(mask)`, pass `~mask`.

Relaxed load/store are implemented in pure Go for inlining. Other orderings use assembly with LSE instructions.

#### 128-bit Operations

| Build Tag | Instructions | Target Hardware |
|-----------|--------------|-----------------|
| (default) | `LDXP/STXP` (LL/SC loop) | All ARMv8+ |
| `-tags=lse2` | `CASP` (single instruction) | ARMv8.4+ with LSE2 |

LL/SC (Load-Link/Store-Conditional) retries on contention. CASP provides single-instruction atomicity but requires newer hardware.

### RISC-V 64-bit

RISC-V RVWMO (Weak Memory Ordering) uses explicit fence instructions:

| Operation | Implementation |
|-----------|----------------|
| Load Relaxed | `LD` |
| Load Acquire | `LD` + `FENCE R,RW` |
| Store Relaxed | `SD` |
| Store Release | `FENCE RW,W` + `SD` |
| RMW | `AMO` instructions with `.aq`/`.rl` modifiers |

128-bit operations use spinlock-based emulation.

### LoongArch 64-bit

LoongArch uses DBAR (data barrier) instructions:

| Operation | Implementation |
|-----------|----------------|
| Load Relaxed | `LD.D` |
| Load Acquire | `LD.D` + `DBAR` |
| Store Relaxed | `ST.D` |
| Store Release | `DBAR` + `ST.D` |
| RMW | `AM*_DB` instructions |

128-bit operations use spinlock-based emulation.

### Fallback

Unsupported architectures use `sync/atomic`, which provides sequential consistency. 128-bit operations on fallback architectures are **not atomic** (two separate 64-bit operations).

## Design Rationale

### Explicit Memory Ordering

1. **Instruction selection on weak architectures**: ARM64/RISC-V select different instructions based on ordering requirements
2. **Documentation**: Ordering suffix documents synchronization intent
3. **Portability**: Code explicitly specifies requirements rather than relying on architecture-specific guarantees
4. **Correctness**: Makes memory ordering decisions explicit and reviewable

### Comparison with sync/atomic

sync/atomic provides sequential consistency, which is:
- **Sufficient** for most use cases
- **Portable** across all architectures
- **Simple** to reason about

Use atomix when:
- Building lock-free data structures
- Interoperating with kernel or hardware interfaces (io_uring, shared memory)
- Porting C/C++ code with explicit memory ordering
- Targeting ARM64/RISC-V where explicit ordering controls instruction selection

## Platform Support

| Platform | Implementation |
|----------|----------------|
| linux/amd64 | Native assembly |
| linux/arm64 | Native assembly with LSE |
| linux/riscv64 | Native assembly (128-bit emulated) |
| linux/loong64 | Native assembly (128-bit emulated) |
| darwin/amd64, darwin/arm64 | Native assembly |
| freebsd/amd64, freebsd/arm64 | Native assembly |
| Other | sync/atomic fallback |

## Compiler Intrinsics

atomix provides a customized Go compiler that emits inline atomic instructions instead of function calls. This transforms function calls into single CPU instructions, eliminating call overhead.

### Quick Start

```bash
# Install the intrinsics-customized compiler
make install-compiler

# Build with intrinsics
make build

# Test with intrinsics
make test

# Verify intrinsics are applied
make verify
```

### What the Compiler Does

The customized compiler adds SSA operations for atomix intrinsics:

| Operation | x86-64 | ARM64 |
|-----------|--------|-------|
| Load (Relaxed) | `MOV` | `LDR` |
| Load (Acquire) | `MOV` | `LDAR` |
| Store (Relaxed) | `MOV` | `STR` |
| Store (Release) | `MOV` | `STLR` |
| Add (AcqRel) | `LOCK XADD` | `LDADDAL` |
| CAS | `LOCK CMPXCHG` | `CASAL` |

**x86-64 TSO optimization:** Release stores use plain `MOV` instead of `XCHG`, leveraging x86-64's Total Store Ordering which provides implicit release semantics for all stores.

### Manual Compiler Setup

If you prefer manual setup over the Makefile:

```bash
# Clone the intrinsics compiler
git clone --branch atomix https://github.com/hayabusa-cloud/go.git ~/github.com/go

# Build the compiler
cd ~/github.com/go/src && ./make.bash

# Use for atomix
GOROOT=~/github.com/go ~/github.com/go/bin/go build ./...
```

See [intrinsics.md](./intrinsics.md) for detailed implementation documentation.

## License

MIT — see [LICENSE](./LICENSE).

©2026 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com/)
