# atomix

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/atomix.svg)](https://pkg.go.dev/code.hybscloud.com/atomix)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/atomix)](https://goreportcard.com/report/github.com/hayabusa-cloud/atomix)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/atomix/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/atomix)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**语言:** [English](README.md) | 简体中文 | [日本語](README.ja.md) | [Español](README.es.md) | [Français](README.fr.md)

Go 语言的显式内存序原子操作库。

## 概述

Go 的 `sync/atomic` 提供顺序一致性的原子操作。本库通过架构特定实现，暴露 C++11/C11 内存模型的内存序（Relaxed、Acquire、Release、AcqRel）。

```go
import "code.hybscloud.com/atomix"

var counter atomix.Int64

// 带内存序后缀的方法 API
counter.AddRelaxed(1)    // Relaxed：无同步
counter.Add(1)           // AcqRel：默认安全内存序

// 用于原始内存的指针 API
var flags int32
atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
```

## 安装

```bash
go get code.hybscloud.com/atomix
```

**要求:** Go 1.25+

## 内存序

本库实现了 C++11 内存模型的四种内存序：

| 内存序 | 语义 |
|--------|------|
| **Relaxed** | 仅保证原子性，无同步或顺序约束 |
| **Acquire** | 后续读写不能重排到此加载之前，与 Release 存储配对 |
| **Release** | 先前读写不能重排到此存储之后，与 Acquire 加载配对 |
| **AcqRel** | 结合 Acquire 和 Release 语义，用于读-修改-写操作 |

### 内存序选择

默认方法（无内存序后缀）使用：
- 加载操作：Relaxed
- 存储操作：Relaxed
- 读-修改-写操作：AcqRel

**注意：** sync/atomic 的 Load 使用 acquire 语义，Store 使用 release 语义（x86 上为顺序一致性）。atomix 默认使用 Relaxed 以在弱序架构上获得最佳性能。需要与 sync/atomic 等效的语义时，请使用 `LoadAcquire`/`StoreRelease`。

### 各内存序的使用场景

| 使用场景 | 内存序 | 原因 |
|----------|--------|------|
| 统计计数器 | Relaxed | 无需同步，可接受最终一致性 |
| 引用计数 | AcqRel | 确保释放前对象状态可见 |
| 生产者-消费者标志 | Release/Acquire | 生产者释放数据，消费者获取 |
| 自旋锁获取 | Acquire | 临界区读取必须看到之前的写入 |
| 自旋锁释放 | Release | 临界区写入必须在解锁前完成 |
| 序列锁 | AcqRel | 双向都需要顺序保证 |

## 类型

### 值类型

| 类型 | 大小 | 描述 |
|------|------|------|
| `Bool` | 4 字节 | 原子布尔值（底层为 uint32） |
| `Int32`, `Uint32` | 4 字节 | 32 位整数 |
| `Int64`, `Uint64` | 8 字节 | 64 位整数 |
| `Uintptr` | 8 字节 | 指针大小的整数 |
| `Pointer[T]` | 8 字节 | 泛型原子指针 |
| `Int128`, `Uint128` | 16 字节 | 128 位整数（需要 16 字节对齐） |

### 填充类型

填充变体（`Int64Padded`、`Uint64Padded` 等）占用完整的缓存行（64 字节），防止多个 CPU 核心访问多个原子变量时的伪共享。

```go
// 无填充：变量可能共享缓存行，导致争用
var a, b atomix.Int64  // 可能在内存中相邻

// 有填充：每个变量占用独立的缓存行
var a, b atomix.Int64Padded  // 保证 64 字节间隔
```

## 操作

| 操作 | 返回值 | 描述 |
|------|--------|------|
| `Load` | 值 | 原子读取 |
| `Store` | — | 原子写入 |
| `Swap` | 旧值 | 原子交换 |
| `CompareAndSwap` | bool | 交换成功返回 true |
| `CompareExchange` | 旧值 | 无论成功与否都返回原值 |
| `Add`, `Sub` | 新值 | 原子算术运算 |
| `Inc`, `Dec` | 新值 | 原子加减 1 |
| `And`, `Or`, `Xor` | 旧值 | 原子位运算 |
| `Max`, `Min` | 旧值 | 原子最大/最小值 |

**返回值语义:** Add/Sub/Inc/Dec 返回**新值**（与 sync/atomic 一致）。Swap/And/Or/Xor/Max/Min 返回**旧值**。

### CompareAndSwap 与 CompareExchange

```go
// CompareAndSwap：返回成功/失败
if v.CompareAndSwap(old, new) {
    // 成功
}

// CompareExchange：返回原值（可实现无需单独 Load 的 CAS 循环）
for {
    old := v.Load()
    new := transform(old)
    if v.CompareExchange(old, new) == old {
        break  // 成功
    }
}
```

## 指针 API

用于内存映射区域、共享内存或 io_uring 环的交互：

```go
var flags int32

atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
atomix.Release.CompareAndSwapInt32(&flags, 0, 1)
```

指针 API 操作原始的 `*int32`、`*int64` 等，而非包装类型。当原子变量无法使用包装类型时（如内核共享结构中的字段）很有用。

## 128 位操作

128 位原子操作需要 16 字节对齐。对于共享内存使用放置辅助函数：

```go
buf := make([]byte, 32)
_, ptr := atomix.PlaceAlignedUint128(buf, 0)
ptr.Store(lo, hi)

var v atomix.Uint128  // 类型保证对齐
v.Store(lo, hi)
```

| 架构 | 128 位实现 |
|------|-----------|
| amd64 | `LOCK CMPXCHG16B` |
| arm64 | `LDXP/STXP`（默认）或 `CASP`（`-tags=lse2`） |
| riscv64, loong64 | 自旋锁模拟（LL/SC 低 64 位） |

**注意:** 128 位原子操作主要用于双字 CAS 模式（如带版本计数器的无锁数据结构）。

## 架构实现

### x86-64 (TSO)

x86-64 提供全存储序（TSO），这是一种强内存模型：
- 所有加载都有隐式 acquire 语义
- 所有存储都有隐式 release 语义
- 存储-加载序需要显式屏障（MFENCE）或锁定指令

因此，在 x86-64 上所有内存序变体编译为相同的机器码。显式内存序在 x86-64 上的主要好处是文档化和可移植性。

| 操作 | 指令 | 备注 |
|------|------|------|
| Load | `MOV` | 普通内存访问 |
| Store | `MOV` | 普通内存访问 |
| Add | `LOCK XADD` | 返回旧值 |
| Swap | `XCHG` | 隐式 LOCK |
| CAS | `LOCK CMPXCHG` | |
| And/Or/Xor | `LOCK CMPXCHG` 循环 | 通过 CAS 循环返回旧值 |
| CAS128 | `LOCK CMPXCHG16B` | |

Load 和 Store 用纯 Go 实现以便编译器内联。

### ARM64（弱序）

ARM64 是弱序内存模型，需要显式的顺序指令。LSE（大系统扩展）提供带内存序后缀的原子指令：

**后缀含义：** 无后缀 = Relaxed，`A` = Acquire，`L` = Release，`AL` = Acquire-Release

| 操作 | Relaxed | Acquire | Release | AcqRel |
|------|---------|---------|---------|--------|
| Load | `LDR` | `LDAR` | — | — |
| Store | `STR` | — | `STLR` | — |
| Add | `LDADD` | `LDADDA` | `LDADDL` | `LDADDAL` |
| CAS | `CAS` | `CASA` | `CASL` | `CASAL` |
| Swap | `SWP` | `SWPA` | `SWPL` | `SWPAL` |
| And | `LDCLR`† | `LDCLRA` | `LDCLRL` | `LDCLRAL` |
| Or | `LDSET` | `LDSETA` | `LDSETL` | `LDSETAL` |
| Xor | `LDEOR` | `LDEORA` | `LDEORL` | `LDEORAL` |

† `LDCLR` 清除位（与补码进行 AND）。实现 `And(mask)` 需要传入 `~mask`。

Relaxed 加载/存储用纯 Go 实现以便内联。其他内存序使用 LSE 指令的汇编实现。

#### 128 位操作

| 构建标签 | 指令 | 目标硬件 |
|----------|------|----------|
| （默认） | `LDXP/STXP`（LL/SC 循环） | 所有 ARMv8+ |
| `-tags=lse2` | `CASP`（单指令） | ARMv8.4+ with LSE2 |

LL/SC（Load-Link/Store-Conditional）在争用时重试。CASP 提供单指令原子性但需要较新硬件。

### RISC-V 64 位

RISC-V RVWMO（弱内存序）使用显式栅栏指令：

| 操作 | 实现 |
|------|------|
| Load Relaxed | `LD` |
| Load Acquire | `LD` + `FENCE R,RW` |
| Store Relaxed | `SD` |
| Store Release | `FENCE RW,W` + `SD` |
| RMW | `AMO` 指令配合 `.aq`/`.rl` 修饰符 |

128 位操作使用基于自旋锁的模拟。

### LoongArch 64 位

LoongArch 使用 DBAR（数据屏障）指令：

| 操作 | 实现 |
|------|------|
| Load Relaxed | `LD.D` |
| Load Acquire | `LD.D` + `DBAR` |
| Store Relaxed | `ST.D` |
| Store Release | `DBAR` + `ST.D` |
| RMW | `AM*_DB` 指令 |

128 位操作使用基于自旋锁的模拟。

### 回退

不支持的架构使用 `sync/atomic`，提供顺序一致性。回退架构上的 128 位操作**非原子**（两个独立的 64 位操作）。

## 设计原理

### 为什么需要显式内存序？

1. **弱序架构上的性能**: ARM64/RISC-V 可以在不需要完全顺序时使用更弱（更快）的指令
2. **文档化**: 内存序后缀记录同步意图
3. **可移植性**: 代码显式指定需求，而非依赖架构特定保证
4. **正确性**: 使内存序决策显式且可审查

### 为什么不直接用 sync/atomic？

sync/atomic 提供顺序一致性，这是：
- **足够的** 适用于大多数场景
- **可移植的** 跨所有架构
- **简单的** 易于理解

使用 atomix 当：
- 构建高性能无锁数据结构
- 与内核或硬件接口交互（io_uring、共享内存）
- 移植有显式内存序的 C/C++ 代码
- 目标是 ARM64/RISC-V，弱序能带来可测量的收益

## 平台支持

| 平台 | 实现 |
|------|------|
| linux/amd64 | 原生汇编 |
| linux/arm64 | 原生汇编，使用 LSE |
| linux/riscv64 | 原生汇编（128 位模拟） |
| linux/loong64 | 原生汇编（128 位模拟） |
| darwin/amd64, darwin/arm64 | 原生汇编 |
| freebsd/amd64, freebsd/arm64 | 原生汇编 |
| 其他 | sync/atomic 回退 |

## 编译器内联

为充分发挥性能，可将 atomix 与 Go 编译器集成，直接生成内联原子指令，消除函数调用开销。实现方案详见 [intrinsics.md](./intrinsics.md)。

## 许可证

MIT — 见 [LICENSE](./LICENSE)。

©2026 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com/)
