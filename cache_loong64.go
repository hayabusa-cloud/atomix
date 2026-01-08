// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build loong64

package atomix

// CacheLineSize is the cache line size on LoongArch 64-bit processors.
// Loongson 3A5000/3A6000 and later use 64-byte cache lines.
const CacheLineSize = 64
