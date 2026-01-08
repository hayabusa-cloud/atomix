# Copilot Instructions for atomix

## Package Overview

`atomix` provides atomic primitives with explicit memory ordering control for Go.
Unlike `sync/atomic` which uses sequential consistency, `atomix` exposes four
memory orderings: Relaxed, Acquire, Release, and AcqRel.

## Architecture

The package uses assembly:

| Architecture | Implementation |
|--------------|----------------|
| amd64 | LOCK-prefixed instructions (XADD, CMPXCHG, XCHG) |
| arm64 | LSE atomics with ordering suffixes (LDADDA/L/AL, CASA/L/AL, SWPA/L/AL) |
| riscv64 | AMO instructions with .aq/.rl suffixes + FENCE |
| loong64 | AM*_DB instructions + DBAR barriers |

## Key Design Patterns

### Memory Ordering Suffixes

All operations have ordering variants:
- Default (no suffix): Load=Relaxed, Store=Relaxed, RMW=AcqRel (matches sync/atomic)
- `Relaxed`: Only atomicity, no ordering
- `Acquire`: Load-acquire semantics
- `Release`: Store-release semantics
- `AcqRel`: Acquire-release for RMW operations

### Return Value Semantics

Return values match sync/atomic behavior:
```go
newVal := counter.Add(1)   // Add/Sub/Inc/Dec return NEW value
oldVal := counter.Swap(5)  // Swap/And/Or/Xor/Max/Min return OLD value
```

### 128-bit Atomics

Int128/Uint128 require 16-byte alignment. True 128-bit atomicity on:
- amd64: LOCK CMPXCHG16B
- arm64: LDXP/STXP (default, ARMv8.0+) or CASP (-tags=lse2, ARMv8.4+)

RISC-V and LoongArch use spinlock-based emulation which provides mutual
exclusion but may exhibit torn reads under concurrent access.

## Code Review Guidelines

### Assembly Files (`internal/arch/*.s`)

1. Verify frame sizes match Go stubs
2. Check NOSPLIT is present on all functions
3. Verify barrier placement for non-TSO architectures
4. Ensure BX register is saved/restored for CMPXCHG16B

### Memory Ordering

1. TSO (x86-64): All orderings collapse to same instructions
2. ARM64: LSE instruction suffixes provide ordering (A=Acquire, L=Release, AL=AcqRel)
3. RISC-V: FENCE instructions for load/store; AMO instructions have .aq/.rl modifiers
4. LoongArch: DBAR for load/store; AM*_DB instructions for RMW

### Testing

- All tests should pass without `-race` flag (lock-free algorithms)
- High-contention tests exercise CAS retry paths
- 128-bit tests must verify alignment requirements

## Common Pitfalls

1. **Wrong return value assumption**: Add/Sub return NEW value (like sync/atomic)
2. **128-bit alignment**: Must use PlaceAligned* for proper alignment; SIGBUS on unaligned access
3. **Race detector**: False positives expected - atomics bypass Go's sync model
4. **ARM64 requirements**: LSE (ARMv8.1+) for 32/64-bit atomics; 128-bit works on ARMv8.0+ (LDXP/STXP)

## File Structure

```
atomix/
├── doc.go           # Package documentation
├── types.go         # Type definitions (noCopy, Pointer)
├── bool.go          # Bool type
├── int32.go         # Int32 type
├── int64.go         # Int64 type
├── uint32.go        # Uint32 type
├── uint64.go        # Uint64 type
├── uintptr.go       # Uintptr type
├── int128.go        # Int128 type (16-byte aligned)
├── uint128.go       # Uint128 type (16-byte aligned)
├── pointer.go       # Pointer[T] generic type
├── align.go         # Placement helpers
├── alloc.go         # Allocator for building structures
├── barrier.go       # Memory barrier functions
├── cache.go         # Cache line size detection
└── internal/arch/   # Platform-specific assembly
    ├── stubs_*.go   # Go function declarations
    ├── asm_amd64.s  # x86-64 assembly
    ├── asm_arm64.s  # ARM64 assembly
    ├── asm_riscv64.s
    └── asm_loong64.s
```
