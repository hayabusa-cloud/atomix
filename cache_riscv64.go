// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build riscv64

package atomix

// CacheLineSize is the cache line size on RISC-V 64-bit processors.
// Common RISC-V implementations use 64-byte cache lines.
const CacheLineSize = 64
