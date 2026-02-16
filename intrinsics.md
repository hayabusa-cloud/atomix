# Compiler Intrinsics

Reference documentation for Go compiler intrinsics in atomix. Intrinsics replace function calls with inline CPU instructions, eliminating call overhead (stack frame, register spilling, prologue/epilogue).

**Status:** Implemented for AMD64 and ARM64 (64-bit operations). RISC-V and LoongArch use assembly stubs. 128-bit operations use assembly stubs on all architectures.

## Instruction Mapping

### ARM64 (LSE)

**Status: Intrinsified (LSE-guarded)**

ARM64 Large System Extensions provide atomic instructions with explicit memory ordering suffixes:

| Operation | Relaxed | Acquire | Release | AcqRel |
|-----------|---------|---------|---------|--------|
| Load | `LDR` | `LDAR` | — | — |
| Store | `STR` | — | `STLR` | — |
| Add | `LDADD` | `LDADDA` | `LDADDL` | `LDADDAL` |
| CAS | `CAS` | `CASA` | `CASL` | `CASAL` |
| Swap | `SWP` | `SWPA` | `SWPL` | `SWPAL` |
| And | `LDCLR` † | `LDCLRA` † | `LDCLRL` † | `LDCLRAL` † |
| Or | `LDSET` | `LDSETA` | `LDSETL` | `LDSETAL` |
| Xor | `LDEOR` | `LDEORA` | `LDEORL` | `LDEORAL` |

† **And operation note:** `LDCLR` clears bits: `old = *addr; *addr = old & ~operand`. To implement `And(mask)`, pass `~mask` to LDCLR.

**Suffix meanings:**
- No suffix: Relaxed (no ordering)
- `A`: Acquire (load ordering)
- `L`: Release (store ordering)
- `AL`: Acquire-Release (full RMW ordering)

**Return value note:** All LSE atomic RMW instructions (`LDADD`, `SWP`, `LDSET`, etc.) return the **old** value. atomix's `Add` returns the **new** value, so the intrinsic must compute `new = old + delta` after the instruction. `Swap`/`And`/`Or`/`Xor` return the old value directly (no conversion needed).

**sync/atomic comparison:** Go's sync/atomic uses `AL` variants (sequential consistency). atomix exposes all orderings.

### x86-64 (TSO)

**Status: Intrinsified (TSO)**

x86-64 Total Store Ordering provides implicit acquire/release. All orderings compile to identical code:

| Operation | Instruction | Notes |
|-----------|-------------|-------|
| Load | `MOV` | Implicit acquire |
| Store | `MOV` | Implicit release |
| Add | `LOCK XADD` | Returns old value |
| CAS | `LOCK CMPXCHG` | RAX = expected |
| Swap | `XCHG` | Implicit LOCK |
| And | `LOCK CMPXCHG` loop | CAS loop returns old value ‡ |
| Or | `LOCK CMPXCHG` loop | CAS loop returns old value ‡ |
| Xor | `LOCK CMPXCHG` loop | CAS loop returns old value ‡ |
| CAS128 | `LOCK CMPXCHG16B` | Assembly stub (not yet intrinsified) |

‡ **Bitwise ops note:** x86 `LOCK AND/OR/XOR` modify memory but don't return the old value. atomix requires old value return, so the implementation uses a `LOCK CMPXCHG` CAS loop: load current value, compute bitwise result, attempt CAS, retry on failure.

**Return value note:** `LOCK XADD` returns the **old** value. atomix's `Add` returns the **new** value, so the intrinsic must compute `new = old + delta` after the instruction.

### RISC-V (RVWMO)

**Status: Assembly stubs (not yet intrinsified)**

RISC-V provides AMO instructions with `.aq` (acquire) and `.rl` (release) modifiers:

| Operation | Relaxed | Acquire | Release | AcqRel |
|-----------|---------|---------|---------|--------|
| Load | `LD` | `LD` + `FENCE R,RW` | — | — |
| Store | `SD` | — | `FENCE RW,W` + `SD` | — |
| Add | `AMOADD.D` | `AMOADD.D.AQ` | `AMOADD.D.RL` | `AMOADD.D.AQRL` |
| Swap | `AMOSWAP.D` | `AMOSWAP.D.AQ` | `AMOSWAP.D.RL` | `AMOSWAP.D.AQRL` |
| And | `AMOAND.D` | `AMOAND.D.AQ` | `AMOAND.D.RL` | `AMOAND.D.AQRL` |
| Or | `AMOOR.D` | `AMOOR.D.AQ` | `AMOOR.D.RL` | `AMOOR.D.AQRL` |
| Xor | `AMOXOR.D` | `AMOXOR.D.AQ` | `AMOXOR.D.RL` | `AMOXOR.D.AQRL` |
| CAS | `LR.D`/`SC.D` | `LR.D.AQ`/`SC.D` | `LR.D`/`SC.D.RL` | `LR.D.AQ`/`SC.D.RL` |

**Return value note:** AMO instructions return the **old** value. atomix's `Add` returns the **new** value, requiring post-instruction addition.

### LoongArch

**Status: Assembly stubs (not yet intrinsified)**

LoongArch uses `AM*_DB` instructions (DB = double-barrier, sequential consistency) and `DBAR` for explicit barriers:

| Operation | Relaxed | Acquire | Release | AcqRel |
|-----------|---------|---------|---------|--------|
| Load | `LD.D` | `LD.D` + `DBAR` | — | — |
| Store | `ST.D` | — | `DBAR` + `ST.D` | — |
| Add | `AMADD.D` | `AMADD_DB.D` | `AMADD_DB.D` | `AMADD_DB.D` |
| Swap | `AMSWAP.D` | `AMSWAP_DB.D` | `AMSWAP_DB.D` | `AMSWAP_DB.D` |
| And | `AMAND.D` | `AMAND_DB.D` | `AMAND_DB.D` | `AMAND_DB.D` |
| Or | `AMOR.D` | `AMOR_DB.D` | `AMOR_DB.D` | `AMOR_DB.D` |
| Xor | `AMXOR.D` | `AMXOR_DB.D` | `AMXOR_DB.D` | `AMXOR_DB.D` |
| CAS | `LL.D`/`SC.D` | + `DBAR` | + `DBAR` | + `DBAR` |

**Return value note:** AM* instructions return the **old** value. atomix's `Add` returns the **new** value, requiring post-instruction addition.

---

## Go Compiler SSA Pipeline

```
Source → AST → SSA (generic) → SSA (arch-specific) → Machine Code
                    ↓                   ↓                  ↓
             genericOps.go        AMD64Ops.go         amd64/ssa.go
                                  ARM64Ops.go         arm64/ssa.go
```

Intrinsics intercept at SSA generation, replacing `CALL` nodes with SSA operations that lower to specific instructions.

**Modified files** (relative to `src/cmd/compile/internal/`):

| File | Purpose |
|------|---------|
| `ssagen/intrinsics.go` | Function name → SSA operation mapping (310 `addF` registrations) |
| `ssagen/intrinsics_test.go` | Whitelist entries for intrinsified functions |
| `ssa/_gen/genericOps.go` | New generic SSA ops: CAX, Xor, relaxed load/store |
| `ssa/_gen/AMD64Ops.go` | New AMD64 ops: `MFENCE`, `CMPXCHGLlockValue`/`QlockValue`, `LoweredAtomicXor32`/`64` |
| `ssa/_gen/ARM64Ops.go` | New ARM64 ops: `LoweredAtomicCax`, `LoweredAtomicXor`, relaxed load/store, `DMB` |
| `ssa/_gen/generic.rules` | Constant folding (identity → Load), CondSelect → math transforms |
| `ssa/_gen/AMD64.rules` | Relaxed/release store → `MOV`, CAX lowering, `MFENCE` coalescing |
| `ssa/_gen/ARM64.rules` | CAX/Xor lowering, relaxed load/store → `MOV`, `DMB` coalescing |
| `amd64/ssa.go` | Code generation: `MFENCE`, CAX (`CMPXCHG`), Xor (CAS loop) |
| `arm64/ssa.go` | Code generation: relaxed load/store, CAX (LL/SC + LSE), Xor (LL/SC + LSE) |

---

## Implementation Architecture

### Design Principle: SSA Op Reuse + Targeted Extensions

The atomix intrinsics reuse Go's existing SSA ops where possible (Load, Store, Add, Swap, CAS, And, Or) and add new ops only for operations that Go's standard library does not provide:

- **CAX** (CompareAndExchange returning old value): `AtomicCompareAndExchange32`/`64` + Variant
- **Xor** (atomic XOR returning old value): `AtomicXor32value`/`64value` + Variant
- **Relaxed load/store** (no ordering): `AtomicLoad8Relaxed`/`32Relaxed`/`64Relaxed`/`PtrRelaxed`, `AtomicStore8Relaxed`/`32Relaxed`/`64Relaxed`/`PtrRelaxedNoWB`
- **Barriers**: `AMD64.MFENCE`, `ARM64.DMB`

The ordering distinction manifests differently per architecture:
- **AMD64 (TSO):** All orderings map to the same SSA op (hardware provides implicit ordering)
- **ARM64 Load/Store:** Different SSA ops exist for Relaxed vs Acquire/Release (e.g., `OpAtomicLoad64Relaxed` vs `OpAtomicLoad64`)
- **ARM64 RMW:** All orderings map to the same `Variant` SSA op, which lowers to the AcqRel instruction (see [Known Limitation](#known-limitation-arm64-rmw-over-ordering))

### SSA Op Mapping

**RMW operations** (Swap, CAS, Add, And, Or, Xor):

| Operation | Generic SSA Op | ARM64 Variant SSA Op |
|-----------|---------------|---------------------|
| Swap32 | `OpAtomicExchange32` | `OpAtomicExchange32Variant` |
| Swap64 | `OpAtomicExchange64` | `OpAtomicExchange64Variant` |
| CAS32 | `OpAtomicCompareAndSwap32` | `OpAtomicCompareAndSwap32Variant` |
| CAS64 | `OpAtomicCompareAndSwap64` | `OpAtomicCompareAndSwap64Variant` |
| CAX32 | `OpAtomicCompareAndExchange32` | (direct, no Variant) |
| CAX64 | `OpAtomicCompareAndExchange64` | (direct, no Variant) |
| Add32 | `OpAtomicAdd32` | `OpAtomicAdd32Variant` |
| Add64 | `OpAtomicAdd64` | `OpAtomicAdd64Variant` |
| And32 | `OpAtomicAnd32value` | `OpAtomicAnd32valueVariant` |
| And64 | `OpAtomicAnd64value` | `OpAtomicAnd64valueVariant` |
| Or32 | `OpAtomicOr32value` | `OpAtomicOr32valueVariant` |
| Or64 | `OpAtomicOr64value` | `OpAtomicOr64valueVariant` |
| Xor32 | `OpAtomicXor32value` | `OpAtomicXor32valueVariant` |
| Xor64 | `OpAtomicXor64value` | `OpAtomicXor64valueVariant` |

**Load operations** (ordering-differentiated SSA ops):

| Ordering | 32-bit | 64-bit | Pointer |
|----------|--------|--------|---------|
| Relaxed | `OpAtomicLoad32Relaxed` | `OpAtomicLoad64Relaxed` | `OpAtomicLoadPtrRelaxed` |
| Acquire | `OpAtomicLoad32` | `OpAtomicLoad64` | `OpAtomicLoadPtr` |

**Store operations** (ordering-differentiated SSA ops):

| Ordering | 32-bit | 64-bit | Pointer |
|----------|--------|--------|---------|
| Relaxed | `OpAtomicStore32Relaxed` | `OpAtomicStore64Relaxed` | `OpAtomicStorePtrRelaxedNoWB` |
| Release | `OpAtomicStoreRel32` | `OpAtomicStoreRel64` | `OpAtomicStorePtrNoWB` |

**Barrier operations** (architecture-specific SSA ops):

| Ordering | AMD64 | ARM64 |
|----------|-------|-------|
| Acquire | no-op | `OpARM64DMB` (0x9 = ISHLD) |
| Release | no-op | `OpARM64DMB` (0xA = ISHST) |
| AcqRel | `OpAMD64MFENCE` | `OpARM64DMB` (0xB = ISH) |

### Intrinsic Registration Patterns

All intrinsics are registered in `ssagen/intrinsics.go` using `addF` with architecture filtering.

**AMD64 pattern** — direct helper functions emitting ordering-agnostic SSA ops:

```go
// All orderings share the same helper (TSO makes them identical)
atomixAdd64 := func(s *state, n *ir.CallExpr, args []*ssa.Value) *ssa.Value {
    v := s.newValue3(ssa.OpAtomicAdd64,
        types.NewTuple(types.Types[types.TUINT64], types.TypeMem),
        args[0], args[1], s.mem())
    s.vars[memVar] = s.newValue1(ssa.OpSelect1, types.TypeMem, v)
    return s.newValue1(ssa.OpSelect0, types.Types[types.TINT64], v)
}

addF("code.hybscloud.com/atomix/internal/arch", "AddInt64Relaxed", atomixAdd64, sys.AMD64)
addF("code.hybscloud.com/atomix/internal/arch", "AddInt64Acquire", atomixAdd64, sys.AMD64)
addF("code.hybscloud.com/atomix/internal/arch", "AddInt64Release", atomixAdd64, sys.AMD64)
addF("code.hybscloud.com/atomix/internal/arch", "AddInt64AcqRel",  atomixAdd64, sys.AMD64)
```

**ARM64 RMW pattern** — `makeAtomicGuardedIntrinsicARM64` with LSE Variant ops:

```go
// All orderings share the same Variant SSA op (over-ordered to AcqRel)
addF("code.hybscloud.com/atomix/internal/arch", "AddInt64Relaxed",
    makeAtomicGuardedIntrinsicARM64(ssa.OpAtomicAdd64, ssa.OpAtomicAdd64Variant,
        types.TINT64, atomicEmitterARM64), sys.ARM64)
addF("code.hybscloud.com/atomix/internal/arch", "AddInt64Acquire",
    makeAtomicGuardedIntrinsicARM64(ssa.OpAtomicAdd64, ssa.OpAtomicAdd64Variant,
        types.TINT64, atomicEmitterARM64), sys.ARM64)
addF("code.hybscloud.com/atomix/internal/arch", "AddInt64Release",
    makeAtomicGuardedIntrinsicARM64(ssa.OpAtomicAdd64, ssa.OpAtomicAdd64Variant,
        types.TINT64, atomicEmitterARM64), sys.ARM64)
addF("code.hybscloud.com/atomix/internal/arch", "AddInt64AcqRel",
    makeAtomicGuardedIntrinsicARM64(ssa.OpAtomicAdd64, ssa.OpAtomicAdd64Variant,
        types.TINT64, atomicEmitterARM64), sys.ARM64)
```

**CAX pattern** — direct on both architectures (no LSE-guarded wrapper):

```go
atomixCax64 := func(s *state, n *ir.CallExpr, args []*ssa.Value) *ssa.Value {
    v := s.newValue4(ssa.OpAtomicCompareAndExchange64,
        types.NewTuple(types.Types[types.TUINT64], types.TypeMem),
        args[0], args[1], args[2], s.mem())
    s.vars[memVar] = s.newValue1(ssa.OpSelect1, types.TypeMem, v)
    return s.newValue1(ssa.OpSelect0, types.Types[types.TINT64], v)
}

// Registered on both archs in a single call
addF("code.hybscloud.com/atomix/internal/arch", "CaxInt64Relaxed", atomixCax64, sys.AMD64, sys.ARM64)
```

**Load/Store pattern** — ordering-specific SSA ops, registered on both architectures:

```go
// Relaxed load (plain MOV on ARM64, no acquire semantics)
atomixLoad64Relaxed := func(s *state, n *ir.CallExpr, args []*ssa.Value) *ssa.Value {
    v := s.newValue2(ssa.OpAtomicLoad64Relaxed, ...)
    // ...
}

// Acquire load (LDAR on ARM64)
atomixLoad64Acquire := func(s *state, n *ir.CallExpr, args []*ssa.Value) *ssa.Value {
    v := s.newValue2(ssa.OpAtomicLoad64, ...)
    // ...
}

addF("...arch", "LoadInt64Relaxed", atomixLoad64Relaxed, sys.AMD64, sys.ARM64)
addF("...arch", "LoadInt64Acquire", atomixLoad64Acquire, sys.AMD64, sys.ARM64)
```

**Critical:** Use `addF` (not `add`) when specifying an architecture list.

### Known Limitation: ARM64 RMW Over-Ordering

All ARM64 RMW intrinsics (Swap, CAS, Add, And, Or, Xor) use the same `Variant` SSA op regardless of the requested ordering. The `Variant` ops lower to AcqRel instructions (`SWPAL`, `LDADDAL`, `CASAL`, etc.). This means:

- `AddInt64Relaxed` on ARM64 emits `LDADDAL` (AcqRel), not `LDADD` (Relaxed)
- `SwapInt32Release` on ARM64 emits `SWPAL` (AcqRel), not `SWPL` (Release)

This is a correctness-preserving trade-off: AcqRel is strictly stronger than any weaker ordering. The performance cost is bounded by the difference between `LDADD` and `LDADDAL`, which is typically small on modern ARM64 implementations (Graviton3/4, Apple M-series).

Per-ordering instruction selection (e.g., `LDADD` for Relaxed, `LDADDA` for Acquire) would require ordering-specific SSA ops and lowering rules, which is planned for a future iteration.

### 128-bit CAS Intrinsics

**Status: Planned — not yet implemented.** Currently uses assembly stubs (`LOCK CMPXCHG16B` on AMD64, `LDXP/STXP` on ARM64).

### Build and Verify

```bash
# Build the modified compiler
cd src && ./make.bash

# Verify intrinsics are applied (should see instructions, not CALL)
GOROOT=$(pwd)/.. ./bin/go build -gcflags='-S' code.hybscloud.com/atomix 2>&1 | \
    grep -E "LDADDAL|SWPAL|CASAL|LDAR|STLR"

# Verify no function calls to internal/arch (intrinsics not applied)
GOROOT=$(pwd)/.. ./bin/go build -gcflags='-S' code.hybscloud.com/atomix 2>&1 | \
    grep "CALL.*internal/arch"
```

---

## Troubleshooting

| Problem | Cause | Solution |
|---------|-------|----------|
| Intrinsic not applied | Wrong package path | Use `internal/arch`, not public API path |
| Function name mismatch | Naming convention | Match exact function name in `internal/arch` |
| Wrong compiler being used | GOROOT not set | Set `GOROOT` to modified compiler path |
| Sees CALL instructions | Standard compiler used | Rebuild with `make install-compiler` |
| ARM64 RMW over-ordered | Design choice | Expected — Variant ops always emit AcqRel |

**Common mistakes:**

```go
// WRONG: Public package path
addF("code.hybscloud.com/atomix", "AddInt64Relaxed", ...)

// CORRECT: Internal package where low-level functions are defined
addF("code.hybscloud.com/atomix/internal/arch", "AddInt64Relaxed", ...)
```

```go
// WRONG: Using add() which doesn't support arch filtering
add("...", "AddInt64Relaxed", ..., sys.ARM64)

// CORRECT: Using addF() for architecture-specific intrinsics
addF("...", "AddInt64Relaxed", ..., sys.ARM64)
```

---

## Compiler Optimizations

The Go compiler applies several optimizations to atomix intrinsics:

### Constant Folding

Identity operations are optimized to simpler loads at compile time via rules in `generic.rules`:

| Pattern | Optimization | Rationale |
|---------|--------------|-----------|
| `Add(ptr, 0)` | `Load(ptr)` | Adding zero returns current value |
| `Or(ptr, 0)` | `Load(ptr)` | OR with zero is identity |
| `And(ptr, ^0)` | `Load(ptr)` | AND with all-ones is identity |
| `Xor(ptr, 0)` | `Load(ptr)` | XOR with zero is identity |

Both base and Variant SSA ops are folded (e.g., `AtomicAdd64` and `AtomicAdd64Variant` both fold). Applied in the SSA generic rewrite pass before architecture-specific lowering.

### Relaxed Memory Ordering Optimization

On ARM64, relaxed operations use faster instructions:

| Operation | Relaxed | Acquire/Release |
|-----------|---------|-----------------|
| Load | `MOV` (plain load) | `LDAR` (load-acquire) |
| Store | `MOV` (plain store) | `STLR` (store-release) |

**Ordering cost:** `LDAR`/`STLR` enforce ordering constraints that `MOV` does not. The latency difference is microarchitecture-dependent. Use relaxed ordering when no inter-thread visibility guarantees are needed.

On x86-64 (TSO), both relaxed and release stores use plain `MOV`:

| Operation | Relaxed | Release |
|-----------|---------|---------|
| Store | `MOV` | `MOV` |

**Rationale:** x86-64 TSO guarantees that all stores have implicit release semantics. A plain `MOV` instruction provides store-release ordering without requiring `XCHG` or `LOCK` prefix. This optimization was implemented via dedicated `AtomicStoreRel32`/`AtomicStoreRel64` SSA operations that lower directly to `MOVLstore`/`MOVQstore`.

### Fence Coalescing

Adjacent memory barriers are merged via rules in `ARM64.rules` and `AMD64.rules`:

**ARM64 DMB coalescing:**
```
DMB x; DMB x → DMB x             (duplicate elimination, any kind)
DMB ISH; DMB ISHLD → DMB ISH     (full barrier subsumes acquire)
DMB ISH; DMB ISHST → DMB ISH     (full barrier subsumes release)
DMB ISHLD; DMB ISH → DMB ISH     (commutative)
DMB ISHST; DMB ISH → DMB ISH     (commutative)
DMB ISHLD; DMB ISHST → DMB ISH   (acquire + release = full)
DMB ISHST; DMB ISHLD → DMB ISH   (commutative)
```

DMB constants: `0x9` = ISHLD (acquire), `0xA` = ISHST (release), `0xB` = ISH (full).

**x86-64 MFENCE coalescing:**
```
MFENCE; MFENCE → MFENCE          (duplicate elimination)
```

### LSE Instruction Selection (ARM64)

The atomix intrinsics compiler always emits LSE instructions directly. atomix requires ARMv8.4+ with mandatory LSE support — no runtime detection or LL/SC fallback is generated.

```go
// From ssagen/intrinsics.go (makeAtomicGuardedIntrinsicARM64common):
// Always use LSE variant (op1) - atomix requires ARM64 v8.4+
// with mandatory LSE support. No runtime detection needed.
emit(s, n, args, op1, typ, needReturn)
```

Both `sync/atomic` and atomix use the `makeAtomicGuardedIntrinsicARM64` wrapper, but with different behavior. Go's `sync/atomic` uses the standard implementation that generates a runtime `cpu.ARM64.HasLSE` branch with both LSE and LL/SC code paths. The atomix fork modifies this wrapper to always select `op1` (the LSE Variant), eliminating the runtime branch entirely.

| Approach | Runtime Check | Instructions Emitted | Target |
|----------|---------------|---------------------|--------|
| sync/atomic | Yes (`cpu.ARM64.HasLSE`) | Both LSE and LL/SC | All ARMv8 |
| atomix | No | LSE only | ARMv8.4+ |

---

## Performance Tips

### 1. Choose Appropriate Memory Ordering

| Use Case | Ordering | ARM64 Cost |
|----------|----------|------------|
| Statistics counters | Relaxed | Lowest |
| Producer-consumer flag | Acquire/Release | Medium |
| Lock implementation | AcqRel | Highest |

**Rule:** Use the weakest ordering that maintains correctness.

### 2. Prefer Add Over CAS for Counters

```go
// GOOD: Single instruction (LDADDAL on ARM64, LOCK XADD on x86)
counter.Add(1)

// AVOID: CAS loop (multiple instructions, retries under contention)
for {
    old := counter.Load()
    if counter.CompareAndSwap(old, old+1) {
        break
    }
}
```

### 3. Use Relaxed for Thread-Local or Non-Synchronized Access

```go
// Thread-local counter (no synchronization needed)
localCount.AddRelaxed(1)

// Periodic flush to shared counter
sharedCount.AddRelease(localCount.SwapRelaxed(0))
```

### 4. Batch Operations to Amortize Barrier Cost

```go
// AVOID: Multiple barriers
for i := range items {
    process(items[i])
    counter.AddAcqRel(1)  // Barrier per iteration
}

// BETTER: Single barrier
for i := range items {
    process(items[i])
    counter.AddRelaxed(1)
}
BarrierAcqRel()  // One barrier at end
```

### 5. Verify Intrinsics Are Applied

```bash
# Check for direct instructions (intrinsics working)
go build -gcflags='-S' . 2>&1 | grep -E 'LDADDAL|SWPAL|LDAR|STLR'

# Check for function calls (intrinsics NOT working)
go build -gcflags='-S' . 2>&1 | grep 'CALL.*atomix'
```

If you see `CALL` instructions to atomix functions, intrinsics are not being applied. Possible causes:
- Wrong Go version (need modified compiler)
- Indirect function calls (intrinsics only work on direct calls)
- Generic type instantiation in some cases

### 6. ARM64: Build with GOARM64=v8.1 for Servers

```bash
# Graviton2/3, Apple M1+, Ampere Altra all support LSE
GOARM64=v8.1 go build -o app .
```

This eliminates runtime LSE detection branches and guarantees single-instruction atomics.

---

## Intrinsic Count Summary

| Category | Functions | Breakdown | addF Calls | Registration Pattern |
|----------|-----------|-----------|------------|---------------------|
| Load | 12 | 6 types × 2 orderings | 12 | Both archs per call |
| Store | 12 | 6 types × 2 orderings | 12 | Both archs per call |
| Add | 20 | 5 types × 4 orderings | 40 | Separate AMD64/ARM64 |
| Swap | 24 | 6 types × 4 orderings | 48 | Separate AMD64/ARM64 |
| CAS | 24 | 6 types × 4 orderings | 48 | Separate AMD64/ARM64 |
| CAX | 24 | 6 types × 4 orderings | 24 | Both archs per call |
| And | 20 | 5 types × 4 orderings | 40 | Separate AMD64/ARM64 |
| Or | 20 | 5 types × 4 orderings | 40 | Separate AMD64/ARM64 |
| Xor | 20 | 5 types × 4 orderings | 40 | Separate AMD64/ARM64 |
| Barriers | 3 | 3 orderings | 6 | Separate AMD64/ARM64 |
| **Total** | **179** | | **310** | |

**Notes:**
- "Functions" = unique `internal/arch` functions intrinsified (on both AMD64 and ARM64).
- "addF Calls" = total registrations in `intrinsics.go`. "Separate AMD64/ARM64" means each function is registered twice (one `addF` per arch). "Both archs per call" means a single `addF` specifies `sys.AMD64, sys.ARM64`.
- Types: Int32, Uint32, Int64, Uint64, Uintptr (5 types). Swap/CAS/CAX add Pointer (6 types).
- Sub/Inc/Dec are Go-level wrappers over Add; they are not separate intrinsics.
- AMD64 uses direct helper functions for all operations. ARM64 uses `makeAtomicGuardedIntrinsicARM64` for RMW operations (Swap, CAS, Add, And, Or, Xor) and direct helpers for Load, Store, CAX, and Barriers.

---

## References

- [Go SSA documentation](https://github.com/golang/go/tree/master/src/cmd/compile/internal/ssa)
- [Go compiler intrinsics source](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssagen/intrinsics.go)
- [ARM Architecture Reference Manual](https://developer.arm.com/documentation/ddi0487/latest)
- [Intel 64 and IA-32 SDM](https://www.intel.com/content/www/us/en/developer/articles/technical/intel-sdm.html)
- [RISC-V ISA Specification](https://riscv.org/technical/specifications/)
- [LoongArch Reference Manual](https://loongson.github.io/LoongArch-Documentation/)
