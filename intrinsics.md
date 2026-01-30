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

**Key files** (relative to `src/cmd/compile/internal/ssa/`):

| File | Purpose |
|------|---------|
| `_gen/genericOps.go` | Architecture-independent SSA operations |
| `_gen/AMD64Ops.go` | x86-64 lowered operations with register constraints |
| `_gen/ARM64Ops.go` | ARM64 lowered operations |
| `_gen/AMD64.rules` | Generic → AMD64 lowering rules |
| `_gen/ARM64.rules` | Generic → ARM64 lowering rules |
| `ssagen/intrinsics.go` | Function name → SSA operation mapping |

---

## Implementation Steps

### Step 1: Define Generic SSA Operations

In `_gen/genericOps.go`, add architecture-independent operations:

```go
// 64-bit atomics with explicit ordering (ARM64 benefits, x86-64 maps to same)
{name: "AtomicAdd64Relaxed", argLength: 3, typ: "(UInt64,Mem)", hasSideEffects: true},
{name: "AtomicAdd64Acquire", argLength: 3, typ: "(UInt64,Mem)", hasSideEffects: true},
{name: "AtomicAdd64Release", argLength: 3, typ: "(UInt64,Mem)", hasSideEffects: true},
// AtomicAdd64AcqRel already exists as AtomicAdd64

// 128-bit CAS: [ptr, oldLo, oldHi, newLo, newHi, mem] → (Bool, Mem)
// NOTE: Planned — not yet implemented. Currently uses assembly stubs.
{name: "AtomicCompareAndSwap128", argLength: 6, typ: "(Bool,Mem)", hasSideEffects: true},
{name: "AtomicCompareAndSwap128Acquire", argLength: 6, typ: "(Bool,Mem)", hasSideEffects: true},
{name: "AtomicCompareAndSwap128Release", argLength: 6, typ: "(Bool,Mem)", hasSideEffects: true},
{name: "AtomicCompareAndSwap128AcqRel", argLength: 6, typ: "(Bool,Mem)", hasSideEffects: true},
```

**Fields:**
- `argLength`: Argument count including memory state
- `typ`: Return type tuple (value, memory state)
- `hasSideEffects`: Prevents dead code elimination and reordering

### Step 2: Define Architecture-Specific Operations

**AMD64** (`_gen/AMD64Ops.go`):

> **128-bit operations below are planned — not yet implemented.** Target instruction: `CMPXCHG16B`.

CMPXCHG16B has fixed register requirements:

```go
cmpxchg16b = regInfo{
    inputs: []regMask{
        gp &^ ax &^ bx &^ cx &^ dx,  // addr (not RAX/RBX/RCX/RDX)
        ax, dx,                       // oldLo (RAX), oldHi (RDX)
        bx, cx,                       // newLo (RBX), newHi (RCX)
        0,                            // mem
    },
    outputs:  []regMask{gp &^ ax &^ bx &^ cx &^ dx, 0},
    clobbers: ax | dx,  // Modified on failure
}

{name: "LoweredAtomicCas128", argLength: 6, reg: cmpxchg16b,
    resultNotInArgs: true, clobberFlags: true, hasSideEffects: true,
    faultOnNilArg0: true},
```

**ARM64** (`_gen/ARM64Ops.go`):

```go
// 64-bit Add with ordering variants
{name: "LoweredAtomicAdd64Relaxed", argLength: 3,
    reg: regInfo{inputs: []regMask{gpspsbg, gpg, 0}, outputs: []regMask{gp, 0}},
    resultNotInArgs: true, hasSideEffects: true, faultOnNilArg0: true},

{name: "LoweredAtomicAdd64Acquire", argLength: 3,
    reg: regInfo{inputs: []regMask{gpspsbg, gpg, 0}, outputs: []regMask{gp, 0}},
    resultNotInArgs: true, hasSideEffects: true, faultOnNilArg0: true},

{name: "LoweredAtomicAdd64Release", argLength: 3,
    reg: regInfo{inputs: []regMask{gpspsbg, gpg, 0}, outputs: []regMask{gp, 0}},
    resultNotInArgs: true, hasSideEffects: true, faultOnNilArg0: true},

// 128-bit CAS using LDXP/STXP or CASP — Planned, not yet implemented
{name: "LoweredAtomicCas128", argLength: 6,
    reg: regInfo{
        inputs:  []regMask{gpspsbg, gpg, gpg, gpg, gpg, 0},
        outputs: []regMask{gp, 0},
    },
    clobberFlags: true, hasSideEffects: true, faultOnNilArg0: true},

{name: "LoweredAtomicCas128Acquire", ...},
{name: "LoweredAtomicCas128Release", ...},
{name: "LoweredAtomicCas128AcqRel", ...},
```

### Step 3: Write Lowering Rules

**AMD64.rules** (TSO: all orderings use same instruction):

```
(AtomicAdd64Relaxed ptr val mem) => (LoweredAtomicAdd64 ptr val mem)
(AtomicAdd64Acquire ptr val mem) => (LoweredAtomicAdd64 ptr val mem)
(AtomicAdd64Release ptr val mem) => (LoweredAtomicAdd64 ptr val mem)

// Release atomic stores - use plain MOV (x86 TSO provides store-release ordering)
(AtomicStoreRel32 ptr val mem) => (MOVLstore ptr val mem)
(AtomicStoreRel64 ptr val mem) => (MOVQstore ptr val mem)

// Planned — not yet implemented:
(AtomicCompareAndSwap128 ...) => (LoweredAtomicCas128 ...)
(AtomicCompareAndSwap128Acquire ...) => (LoweredAtomicCas128 ...)
(AtomicCompareAndSwap128Release ...) => (LoweredAtomicCas128 ...)
(AtomicCompareAndSwap128AcqRel ...) => (LoweredAtomicCas128 ...)
```

**ARM64.rules** (different instructions per ordering):

```
(AtomicAdd64Relaxed ptr val mem) => (LoweredAtomicAdd64Relaxed ptr val mem)
(AtomicAdd64Acquire ptr val mem) => (LoweredAtomicAdd64Acquire ptr val mem)
(AtomicAdd64Release ptr val mem) => (LoweredAtomicAdd64Release ptr val mem)

// Release atomic stores - use STLR (store-release)
(AtomicStoreRel32 ptr val mem) => (STLRW <types.TypeMem> ptr val mem)
(AtomicStoreRel64 ptr val mem) => (STLR  <types.TypeMem> ptr val mem)

// Planned — not yet implemented:
(AtomicCompareAndSwap128 ...) => (LoweredAtomicCas128 ...)
(AtomicCompareAndSwap128Acquire ...) => (LoweredAtomicCas128Acquire ...)
(AtomicCompareAndSwap128Release ...) => (LoweredAtomicCas128Release ...)
(AtomicCompareAndSwap128AcqRel ...) => (LoweredAtomicCas128AcqRel ...)
```

### Step 4: Implement Code Generation

**AMD64** (`amd64/ssa.go`) — *128-bit codegen below is planned, not yet implemented:*

```go
case ssa.OpAMD64LoweredAtomicCas128:
    // CMPXCHG16B: compares RDX:RAX with [mem]
    // If equal: stores RCX:RBX to [mem], sets ZF=1
    // If not equal: loads [mem] to RDX:RAX, sets ZF=0
    s.Prog(x86.ALOCK)
    p := s.Prog(x86.ACMPXCHG16B)
    p.To.Type = obj.TYPE_MEM
    p.To.Reg = v.Args[0].Reg()

    // Set result from ZF
    p2 := s.Prog(x86.ASETEQ)
    p2.To.Type = obj.TYPE_REG
    p2.To.Reg = v.Reg0()
```

**ARM64** (`arm64/ssa.go`):

```go
case ssa.OpARM64LoweredAtomicAdd64Relaxed:
    // LDADD: old = *addr; *addr = old + val; return old
    // atomix Add returns NEW value, so we compute: new = old + val
    p := s.Prog(arm64.ALDADD)
    p.From.Type = obj.TYPE_REG
    p.From.Reg = v.Args[1].Reg()  // val
    p.To.Type = obj.TYPE_MEM
    p.To.Reg = v.Args[0].Reg()    // addr
    p.RegTo2 = v.Reg0()           // result (old value)

    // Convert old → new: result = result + val
    p2 := s.Prog(arm64.AADD)
    p2.From.Type = obj.TYPE_REG
    p2.From.Reg = v.Args[1].Reg() // val
    p2.To.Type = obj.TYPE_REG
    p2.To.Reg = v.Reg0()          // result = old + val = new

case ssa.OpARM64LoweredAtomicAdd64Acquire:
    p := s.Prog(arm64.ALDADDA)    // LDADDA (acquire)
    // ... same operand setup, plus ADD for old → new conversion

case ssa.OpARM64LoweredAtomicAdd64Release:
    p := s.Prog(arm64.ALDADDL)    // LDADDL (release)
    // ... same operand setup, plus ADD for old → new conversion

// 128-bit CAS — Planned, not yet implemented:
case ssa.OpARM64LoweredAtomicCas128:
    addr := v.Args[0].Reg()
    oldLo, oldHi := v.Args[1].Reg(), v.Args[2].Reg()
    newLo, newHi := v.Args[3].Reg(), v.Args[4].Reg()
    out := v.Reg0()
    tmp0, tmp1 := int16(arm64.REGTMP), int16(arm64.REGTMP-1)

    // again: LDXP (tmp0, tmp1), [addr]
    again := s.Prog(arm64.ALDXP)
    again.From.Type = obj.TYPE_MEM
    again.From.Reg = addr
    again.To.Type = obj.TYPE_REGREG
    again.To.Reg = tmp0
    again.To.Offset = int64(tmp1)

    // CMP tmp0, oldLo; BNE fail
    // CMP tmp1, oldHi; BNE fail
    // STXP result, (newLo, newHi), [addr]
    // CBNZ result, again  // SC failed, retry
    // MOV $1, out; B done
    // fail: MOV $0, out
    // done:
```

### Step 5: Register Intrinsics

In `ssagen/intrinsics.go`, map function names to SSA operations:

```go
// Internal package path (not the public API)
const atomixArch = "code.hybscloud.com/atomix/internal/arch"

// 64-bit Add with ordering (ARM64 only - x86-64 TSO handles all orderings)
addF(atomixArch, "AddInt64Relaxed",
    func(s *state, n *ir.CallExpr, args []*ssa.Value) *ssa.Value {
        v := s.newValue3(ssa.OpAtomicAdd64Relaxed,
            types.NewTuple(types.Types[types.TUINT64], types.TypeMem),
            args[0], args[1], s.mem())
        s.vars[memVar] = s.newValue1(ssa.OpSelect1, types.TypeMem, v)
        return s.newValue1(ssa.OpSelect0, types.Types[types.TUINT64], v)
    },
    sys.ARM64)

addF(atomixArch, "AddInt64Acquire",
    func(s *state, n *ir.CallExpr, args []*ssa.Value) *ssa.Value {
        v := s.newValue3(ssa.OpAtomicAdd64Acquire,
            types.NewTuple(types.Types[types.TUINT64], types.TypeMem),
            args[0], args[1], s.mem())
        s.vars[memVar] = s.newValue1(ssa.OpSelect1, types.TypeMem, v)
        return s.newValue1(ssa.OpSelect0, types.Types[types.TUINT64], v)
    },
    sys.ARM64)

// 128-bit CAS (both architectures) — Planned, not yet implemented
addF(atomixArch, "CasUint128",
    func(s *state, n *ir.CallExpr, args []*ssa.Value) *ssa.Value {
        v := s.newValue6(ssa.OpAtomicCompareAndSwap128,
            types.NewTuple(types.Types[types.TBOOL], types.TypeMem),
            args[0], args[1], args[2], args[3], args[4], s.mem())
        s.vars[memVar] = s.newValue1(ssa.OpSelect1, types.TypeMem, v)
        return s.newValue1(ssa.OpSelect0, types.Types[types.TBOOL], v)
    },
    sys.ARM64, sys.AMD64)
```

**Critical:** Use `addF` (not `add`) when specifying an architecture list.

### Step 6: Build and Verify

```bash
# Regenerate SSA from _gen files
cd src/cmd/compile/internal/ssa/_gen && go run .

# Build the modified compiler
cd ../../../../../ && ./make.bash

# Verify intrinsics are applied (should see instructions, not CALL)
GOROOT=$(pwd) ./bin/go build -gcflags='-S' code.hybscloud.com/atomix 2>&1 | \
    grep -E "LDADDA|LDADDAL|CMPXCHG16B|CASAL"

# Verify no function calls to internal/arch (intrinsics not applied)
GOROOT=$(pwd) ./bin/go build -gcflags='-S' code.hybscloud.com/atomix 2>&1 | \
    grep "CALL.*internal/arch"
```

---

## Troubleshooting

| Problem | Cause | Solution |
|---------|-------|----------|
| "unknown Op" at compile time | Missing lowering rule | Add rule in `*.rules`, run `go generate` |
| Register allocation failure | Incorrect `regInfo` | Check register masks, ensure `clobbers` is complete |
| Intrinsic not applied | Wrong package path | Use `internal/arch`, not public API path |
| Function name mismatch | Naming convention | Match exact function name in internal/arch |
| Wrong compiler being used | GOROOT not set | Set `GOROOT` to modified compiler path |
| Assembler error | Invalid instruction | Verify instruction exists on target arch |

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

Identity operations are optimized to simpler loads at compile time:

| Pattern | Optimization | Rationale |
|---------|--------------|-----------|
| `Add(ptr, 0)` | `Load(ptr)` | Adding zero returns current value |
| `Or(ptr, 0)` | `Load(ptr)` | OR with zero is identity |
| `And(ptr, ^0)` | `Load(ptr)` | AND with all-ones is identity |
| `Xor(ptr, 0)` | `Load(ptr)` | XOR with zero is identity |

These rules are applied in the SSA generic rewrite pass before architecture-specific lowering.

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

Adjacent memory barriers are merged to eliminate redundant instructions:

**ARM64 DMB coalescing:**
```
DMB ISH; DMB ISH → DMB ISH       (duplicate elimination)
DMB ISH; DMB ISHLD → DMB ISH     (full barrier subsumes acquire)
DMB ISH; DMB ISHST → DMB ISH     (full barrier subsumes release)
DMB ISHLD; DMB ISHST → DMB ISH   (acquire + release = full)
```

**x86-64 MFENCE coalescing:**
```
MFENCE; MFENCE → MFENCE          (duplicate elimination)
```

### LSE Instruction Selection (ARM64)

The atomix intrinsics compiler always emits LSE instructions directly. atomix requires ARMv8.1+ with mandatory LSE support — no runtime detection or LL/SC fallback is generated.

```go
// From ssagen/intrinsics.go:
// Always use LSE variant (op1) - atomix requires ARM64 v8.4+
// with mandatory LSE support. No runtime detection needed.
```

This differs from Go's `sync/atomic`, which uses `makeAtomicGuardedIntrinsicARM64` with runtime `cpu.ARM64.HasLSE` branching to support pre-LSE hardware. atomix uses a simplified path that emits only LSE instructions.

| Approach | Runtime Check | Instructions Emitted | Target |
|----------|---------------|---------------------|--------|
| sync/atomic | Yes (`cpu.ARM64.HasLSE`) | Both LSE and LL/SC | All ARMv8 |
| atomix | No | LSE only | ARMv8.1+ |

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

| Category | Functions | Breakdown | addF Calls |
|----------|-----------|-----------|------------|
| Load | 12 | 6 types × 2 orderings | 12 |
| Store | 12 | 6 types × 2 orderings | 12 |
| Add | 20 | 5 types × 4 orderings | 40 |
| Swap | 24 | 6 types × 4 orderings | 48 |
| CAS | 24 | 6 types × 4 orderings | 48 |
| CAX | 24 | 6 types × 4 orderings | 24 |
| And | 20 | 5 types × 4 orderings | 40 |
| Or | 20 | 5 types × 4 orderings | 40 |
| Xor | 20 | 5 types × 4 orderings | 40 |
| Barriers | 3 | 3 orderings | 6 |
| **Total** | **179** | | **310** |

**Notes:**
- "Functions" = unique `internal/arch` functions intrinsified (on both AMD64 and ARM64).
- "addF Calls" = total registrations in `intrinsics.go` (some register both archs per call, others register separately).
- Types: Int32, Uint32, Int64, Uint64, Uintptr (5 types). Swap/CAS/CAX add Pointer (6 types).
- Sub/Inc/Dec are Go-level wrappers over Add; they are not separate intrinsics.

---

## References

- [Go SSA documentation](https://github.com/golang/go/tree/master/src/cmd/compile/internal/ssa)
- [Go compiler intrinsics source](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssagen/intrinsics.go)
- [ARM Architecture Reference Manual](https://developer.arm.com/documentation/ddi0487/latest)
- [Intel 64 and IA-32 SDM](https://www.intel.com/content/www/us/en/developer/articles/technical/intel-sdm.html)
- [RISC-V ISA Specification](https://riscv.org/technical/specifications/)
- [LoongArch Reference Manual](https://loongson.github.io/LoongArch-Documentation/)
