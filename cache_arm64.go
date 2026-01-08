// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build arm64

package atomix

// CacheLineSize is the cache line size on ARM64 processors.
// Most ARM Cortex and Neoverse cores use 64-byte cache lines.
// Apple M-series efficiency cores use 128 bytes, but 64 is safe for both.
const CacheLineSize = 64
