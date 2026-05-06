# Copilot Instructions for atomix

Use this file only as a bottom-line sanity checklist after reading the changed
implementation, tests, GoDocs, and README. It is not an issue oracle. Do not
open a review finding from this file alone.

If this file is broader, older, or more conservative than the current code and
public documentation, treat that as instruction drift. Prefer the current
implementation and public docs, then update this file as guidance.

## Review Use

- Start from the diff, implementation, tests, GoDocs, and README.
- Use this file to ask final verification questions, not to invent findings.
- Report only line-level contradictions against the current implementation,
  public contract, or tests.
- Do not report pattern mismatches, optional hardening, or generalized memory
  model concerns solely because this file mentions them.

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
- Default (no suffix): Load=Relaxed, Store=Relaxed, RMW=AcqRel
- `Relaxed`: Only atomicity, no ordering
- `Acquire`: Load-acquire semantics
- `Release`: Store-release semantics
- `AcqRel`: Acquire-release for RMW operations

### Return Value Semantics

Use the current GoDocs as the authoritative return-value contract. The common bottom line is:
```go
newVal := counter.Add(1)   // Wrapper Add/Sub/Inc/Dec return NEW value
oldVal := counter.Swap(5)  // Swap/And/Or/Xor/Max/Min return OLD value
```

Pointer-based 32/64/uintptr Add/Sub return the new value. Pointer-based MemoryOrder.AddInt128 and MemoryOrder.AddUint128 return the old value by accepted contract.

### 128-bit Atomics

Int128/Uint128 require 16-byte alignment. True 128-bit atomicity on:
- amd64: LOCK CMPXCHG16B
- arm64: LDXP/STXP (default, ARMv8.0+) or CASP (-tags=lse2, ARMv8.4+)

RISC-V and LoongArch use low-word LR/SC or LL/SC emulation for 128-bit operations. Treat those paths as limited emulation, not true 128-bit atomicity. On those architectures, non-relaxed 128-bit Swap variants alias the relaxed swap path.

## Review Sanity Checks

Use these checks after line-level verification against the current code and
docs. They describe common contracts; they are not standalone findings.

### Assembly Files (`internal/arch/*.s`)

1. Compare frame sizes with Go stubs when either side changes
2. Confirm NOSPLIT remains on hot assembly entry points
3. Check barrier placement when non-TSO ordering code changes
4. For amd64 128-bit changes, confirm BX handling around CMPXCHG16B

### Memory Ordering

Verify the touched implementation before reporting an ordering issue.

1. TSO (x86-64): All orderings collapse to same instructions
2. ARM64: LSE instruction suffixes provide ordering (A=Acquire, L=Release, AL=AcqRel)
3. RISC-V: FENCE instructions for load/store; AMO instructions have .aq/.rl modifiers
4. LoongArch: DBAR for load/store; AM*_DB instructions for RMW

### Testing

- Passing normal tests without `-race` is useful evidence unless the failure is
  an environment limitation.
- High-contention tests are useful evidence for CAS retry paths.
- 128-bit alignment tests are useful evidence for 128-bit changes.

## Common Pitfalls

Treat these as verification questions, not automatic findings.

1. **Wrong return value assumption**: Most wrapper Add/Sub and pointer-based 32/64/uintptr Add/Sub operations return the new value; pointer-based 128-bit Add returns the old value
2. **128-bit alignment**: Prefer PlaceAligned* helpers for alignment-sensitive examples; unaligned 128-bit access can SIGBUS
3. **Race detector**: Do not report race output as a product race solely from `-race`; first verify whether the atomic path bypasses Go's race-detector model and whether the non-race contract fails
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
