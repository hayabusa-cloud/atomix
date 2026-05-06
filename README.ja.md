# atomix

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/atomix.svg)](https://pkg.go.dev/code.hybscloud.com/atomix)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/atomix)](https://goreportcard.com/report/github.com/hayabusa-cloud/atomix)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/atomix/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/atomix)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**言語:** [English](README.md) | [简体中文](README.zh-CN.md) | 日本語 | [Español](README.es.md) | [Français](README.fr.md)

明示的メモリオーダリングを備えた Go アトミック操作ライブラリ。

## 概要

Go の `sync/atomic` パッケージは逐次一貫性のアトミック操作を提供する。本ライブラリはアーキテクチャ固有の実装を通じて、C++11/C11 メモリモデルのオーダリング（Relaxed、Acquire、Release、AcqRel）を公開する。

```go
import "code.hybscloud.com/atomix"

var counter atomix.Int64

// サフィックス付きメソッド API
counter.AddRelaxed(1)    // Relaxed：同期なし
counter.Add(1)           // AcqRel：デフォルトの RMW オーダリング

// 生メモリ用ポインタ API
var flags int32
atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
```

## インストール

```bash
go get code.hybscloud.com/atomix
```

**動作要件:** Go 1.26+

## メモリオーダリング

本ライブラリは C++11 メモリモデルの 4 つのオーダリングを実装する：

| オーダリング | セマンティクス |
|--------------|----------------|
| **Relaxed** | 原子性のみ。同期や順序制約なし。 |
| **Acquire** | 後続の読み書きはこのロード前に並べ替え不可。Release ストアとペアで使用。 |
| **Release** | 先行する読み書きはこのストア後に並べ替え不可。Acquire ロードとペアで使用。 |
| **AcqRel** | Acquire と Release の両方のセマンティクスを組み合わせ。読み取り-変更-書き込み操作に使用。 |

### オーダリングの選択

デフォルトメソッド（オーダリングサフィックスなし）使用：
- Load 操作：Relaxed
- Store 操作：Relaxed
- 読み取り-変更-書き込み操作：AcqRel

**注意：** sync/atomic 操作は逐次一貫性を持つ。atomix は Load と Store で Relaxed をデフォルトとし、ウィークオーダーアーキテクチャで異なる命令にマッピングされる（例：ARM64 の LDR と LDAR）。acquire/release の同期辺が必要な場合は `LoadAcquire`/`StoreRelease` を使用すること。

### オーダリング選択表

| 使用場面 | オーダリング | 理由 |
|----------|--------------|------|
| 統計カウンタ | Relaxed | 同期不要、最終的一貫性で十分 |
| 参照カウント | AcqRel | 解放前にオブジェクト状態の可視性を保証 |
| プロデューサー-コンシューマーフラグ | Release/Acquire | プロデューサーがデータを release、コンシューマーが acquire |
| スピンロック獲得 | Acquire | クリティカルセクションの読み取りは以前の書き込みを見る必要がある |
| スピンロック解放 | Release | クリティカルセクションの書き込みはアンロック前に完了必須 |
| シーケンスロック | AcqRel | 双方向でオーダリングが必要 |

## 型

### 値型

| 型 | サイズ | 説明 |
|----|--------|------|
| `Bool` | 4 バイト | アトミックブール値（uint32 ベース） |
| `Int32`, `Uint32` | 4 バイト | 32 ビット整数 |
| `Int64`, `Uint64` | 8 バイト | 64 ビット整数 |
| `Uintptr` | 8 バイト | ポインタサイズ整数 |
| `Pointer[T]` | 8 バイト | ジェネリックアトミックポインタ |
| `Int128`, `Uint128` | 16 バイト | 128 ビット整数（16 バイトアラインメント必要） |

### パディング型

パディング変種（`Int64Padded`、`Uint64Padded` など）は完全なキャッシュライン（64 バイト）を占有し、複数のアトミック変数が異なる CPU コアからアクセスされる際の偽共有を防止する。

```go
// パディングなし：変数がキャッシュラインを共有し、競合が発生する可能性
var a, b atomix.Int64  // メモリ上で隣接する可能性あり

// パディングあり：各変数が専用のキャッシュラインを占有
var a, b atomix.Int64Padded  // 64 バイト間隔を保証
```

## 操作

| 操作 | 戻り値 | 説明 |
|------|--------|------|
| `Load` | 値 | アトミック読み取り |
| `Store` | — | アトミック書き込み |
| `Swap` | 旧値 | アトミック交換 |
| `CompareAndSwap` | bool | 交換が発生した場合 true を返す |
| `CompareExchange` | 旧値 | 成功・失敗に関わらず以前の値を返す |
| `Add`, `Sub` | 新値 | アトミック算術演算 |
| `Inc`, `Dec` | 新値 | 1 のアトミックインクリメント/デクリメント |
| `And`, `Or`, `Xor` | 旧値 | アトミックビット演算 |
| `Max`, `Min` | 旧値 | アトミック最大/最小 |

**戻り値の意味:** ラッパー型の Add/Sub/Inc/Dec とポインタベースの 32/64/uintptr Add/Sub は**新値**を返す。ポインタベースの `MemoryOrder.AddInt128` と `MemoryOrder.AddUint128` は**旧値**を返す。Swap/And/Or/Xor/Max/Min は**旧値**を返す。

### CompareAndSwap と CompareExchange

```go
// CompareAndSwap：成功/失敗を返す
if v.CompareAndSwap(old, new) {
    // 成功
}

// CompareExchange：以前の値を返す（別途 Load なしで CAS ループを実現）
for {
    old := v.Load()
    new := transform(old)
    if v.CompareExchange(old, new) == old {
        break  // 成功
    }
}
```

## ポインタ API

メモリマップ領域、共有メモリ、io_uring リングとの相互運用のため：

```go
var flags int32

atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
atomix.Release.CompareAndSwapInt32(&flags, 0, 1)
```

ポインタ API はラッパー型ではなく、生の `*int32`、`*int64` などを操作する。アトミック変数がラッパー型を使用できない場合（カーネル共有構造体のフィールドなど）に有用。

## 128 ビット操作

128 ビットアトミック操作は 16 バイトアラインメントを必要とする。共有メモリには配置ヘルパーを使用：

```go
buf := make([]byte, 32)
_, ptr := atomix.PlaceAlignedUint128(buf, 0)
ptr.Store(lo, hi)

var v atomix.Uint128  // この変数が16バイト境界に配置されることを確認
v.Store(lo, hi)
```

| アーキテクチャ | 128 ビット実装 |
|----------------|----------------|
| amd64 | `LOCK CMPXCHG16B` |
| arm64 | `LDXP/STXP`（デフォルト）または `CASP`（`-tags=lse2`） |
| riscv64、loong64 | riscv64 は下位 64 ビット LR/SC、loong64 は LL/SC エミュレーション。真の 128 ビットアトミック性ではない |

riscv64 と loong64 では、128 ビット値の `SwapAcquire`、`SwapRelease`、`SwapAcqRel` は下位 64 ビットの relaxed swap 経路にエイリアスする。acquire/release の公開が必要な場合は、別の 32/64 ビット同期プリミティブまたは外部同期を使用すること。

**注意:** 128 ビットアトミック操作は主にダブルワード CAS パターン（バージョンカウンタ付きロックフリーデータ構造など）に有用。

## アーキテクチャ実装

### x86-64 (TSO)

x86-64 は Total Store Ordering（TSO）を提供する強いメモリモデル：
- すべてのロードは暗黙の acquire セマンティクスを持つ
- すべてのストアは暗黙の release セマンティクスを持つ
- ストア-ロード順序は明示的バリア（MFENCE）またはロック命令を必要とする

したがって、x86-64 ではすべてのオーダリング変種が同一の機械語にコンパイルされる。x86-64 での明示的オーダリングの主な役割はドキュメント化と移植性。

| 操作 | 命令 | 備考 |
|------|------|------|
| Load | `MOV` | 通常のメモリアクセス |
| Store | `MOV` | 通常のメモリアクセス |
| Add | `LOCK XADD` | 命令は旧オペランドを返す。新値を返す Add API は命令後に計算する |
| Swap | `XCHG` | 暗黙の LOCK |
| CAS | `LOCK CMPXCHG` | |
| And/Or/Xor | `LOCK CMPXCHG` ループ | CAS ループで旧値を返す |
| CAS128 | `LOCK CMPXCHG16B` | |

Load と Store はコンパイラインライン化のため純粋な Go で実装。

### ARM64（弱オーダー）

ARM64 は弱いメモリモデルで、明示的なオーダリング命令を必要とする。atomix は 32/64 ビット LSE アトミックのパッケージベースラインを ARMv8.4+ として文書化する。128 ビット LL/SC 経路は別途 ARMv8.1+ として文書化する。LSE（Large System Extensions）はオーダリングサフィックス付きのアトミック命令を提供：

**サフィックスの意味：** サフィックスなし = Relaxed、`A` = Acquire、`L` = Release、`AL` = Acquire-Release

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

† `LDCLR` はビットをクリア（補数との AND）。`And(mask)` を実装するには `~mask` を渡す。

Relaxed load/store はインライン化のため純粋な Go で実装。その他のオーダリングは LSE 命令を使用したアセンブリ実装。

#### 128 ビット操作

| ビルドタグ | 命令 | 対象ハードウェア |
|------------|------|------------------|
| （デフォルト） | `LDXP/STXP`（LL/SC ループ） | ARMv8.1+ LL/SC 経路 |
| `-tags=lse2` | `CASP`（単一命令） | ARMv8.4+ LSE2 搭載 |

LL/SC（Load-Link/Store-Conditional）は競合時にリトライし、ARMv8.1+ として文書化する。CASP は ARMv8.4+ LSE2 で単一命令のアトミック性を提供する。

### RISC-V 64 ビット

RISC-V RVWMO（弱いメモリオーダリング）は明示的 fence 命令を使用：

| 操作 | 実装 |
|------|------|
| Load Relaxed | `LD` |
| Load Acquire | `LD` + `FENCE R,RW` |
| Store Relaxed | `SD` |
| Store Release | `FENCE RW,W` + `SD` |
| RMW | `.aq`/`.rl` 修飾子付き `AMO` 命令 |

128 ビット操作は下位 64 ビット LR/SC エミュレーションを使用し、ティアリング読み取りが発生する可能性がある。

### LoongArch 64 ビット

LoongArch は DBAR（データバリア）命令を使用：

| 操作 | 実装 |
|------|------|
| Load Relaxed | `LD.D` |
| Load Acquire | `LD.D` + `DBAR` |
| Store Relaxed | `ST.D` |
| Store Release | `DBAR` + `ST.D` |
| RMW | `AM*_DB` 命令 |

128 ビット操作は下位 64 ビット LL/SC エミュレーションを使用し、ティアリング読み取りが発生する可能性がある。

### フォールバック

サポートされていないアーキテクチャは `sync/atomic` を使用し、逐次一貫性を提供。フォールバックアーキテクチャでの 128 ビット操作は**アトミックではない**（2 つの独立した 64 ビット操作）。

## 設計理念

### 明示的メモリオーダリング

1. **弱いアーキテクチャでの命令選択**: ARM64/RISC-V はオーダリング要件に応じて異なる命令を選択
2. **ドキュメント化**: オーダリングサフィックスが同期意図を文書化
3. **移植性**: コードがアーキテクチャ固有の保証に依存せず、要件を明示的に指定
4. **正確性**: メモリオーダリングの決定を明示的かつレビュー可能に

### sync/atomic との比較

sync/atomic はすべての操作に単一の逐次一貫オーダリングを提供する：
- 操作全体で単一のオーダリング契約
- Go 対応アーキテクチャ間で移植可能な挙動
- 操作ごとのオーダリング選択なし

atomix の対象：
- ロックフリーデータ構造
- カーネルやハードウェアインターフェース（io_uring、共有メモリ）との相互運用
- 明示的メモリオーダリングを持つ C/C++ コードの移植
- 明示的なオーダリングで命令選択を制御する ARM64/RISC-V 経路

## プラットフォームサポート

| プラットフォーム | 実装 |
|------------------|------|
| linux/amd64 | ネイティブアセンブリ |
| linux/arm64 | ネイティブアセンブリ + LSE、ARMv8.4+ パッケージベースライン |
| linux/riscv64 | ネイティブアセンブリ（128 ビットエミュレート） |
| linux/loong64 | ネイティブアセンブリ（128 ビットエミュレート） |
| darwin/amd64, darwin/arm64 | ネイティブアセンブリ |
| freebsd/amd64, freebsd/arm64 | ネイティブアセンブリ |
| その他 | sync/atomic フォールバック |

## コンパイラ組み込み関数

atomix はカスタマイズされた Go コンパイラを提供し、関数呼び出しの代わりにインラインアトミック命令を直接生成する。これにより関数呼び出しが単一の CPU 命令に変換され、呼び出しオーバーヘッドが排除される。

### コンパイラワークフロー

```bash
# 組み込み関数対応コンパイラをインストール
make install-compiler

# 組み込み関数対応コンパイラでビルド
make build

# 組み込み関数対応コンパイラでテスト（テストタイムアウト 120s）
make test

# 長めのタイムアウトでベンチマークを実行
make bench

# 組み込み関数が適用されているか検証
make verify
```

### コンパイラの機能

カスタマイズされたコンパイラは atomix 用の SSA 操作を追加：

| 操作 | x86-64 | ARM64 |
|------|--------|-------|
| Load (Relaxed) | `MOV` | `LDR` |
| Load (Acquire) | `MOV` | `LDAR` |
| Store (Relaxed) | `MOV` | `STR` |
| Store (Release) | `MOV` | `STLR` |
| Add (AcqRel) | `LOCK XADD` | `LDADDAL` |
| CAS | `LOCK CMPXCHG` | `CASAL` |

**x86-64 TSO 最適化：** Release ストアは `XCHG` ではなく通常の `MOV` を使用。x86-64 の Total Store Ordering がすべてのストアに暗黙の release セマンティクスを提供するため。

### 手動コンパイラセットアップ

Makefile を使用しない場合の手動セットアップ：

```bash
# 組み込み関数対応コンパイラをクローン
git clone --branch atomix https://github.com/hayabusa-cloud/go.git ~/github.com/go

# コンパイラをビルド
cd ~/github.com/go/src && ./make.bash

# atomix で使用
GOROOT=~/github.com/go ~/github.com/go/bin/go build ./...
```

詳細な実装ドキュメントは [intrinsics.md](./intrinsics.md) を参照。

## ライセンス

MIT — [LICENSE](./LICENSE) を参照。

©2026 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com/)
