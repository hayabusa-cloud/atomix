// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build !amd64 && !arm64 && !riscv64 && !loong64

package atomix

// CacheLineSize is the assumed cache line size on unknown architectures.
// 64 bytes is a common default on most modern processors.
const CacheLineSize = 64
