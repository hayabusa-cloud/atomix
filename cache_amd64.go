// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build amd64

package atomix

// CacheLineSize is the cache line size on x86-64 processors.
// Intel and AMD processors use 64-byte cache lines.
const CacheLineSize = 64
