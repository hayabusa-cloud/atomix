// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

import "code.hybscloud.com/atomix/internal/arch"

// BarrierAcquire issues an acquire memory barrier.
// Subsequent memory operations cannot be reordered before this barrier.
//
//go:nosplit
func BarrierAcquire() {
	arch.BarrierAcquire()
}

// BarrierRelease issues a release memory barrier.
// Prior memory operations cannot be reordered after this barrier.
//
//go:nosplit
func BarrierRelease() {
	arch.BarrierRelease()
}

// BarrierAcqRel issues a full memory barrier (acquire + release).
// Provides both acquire and release semantics.
//
//go:nosplit
func BarrierAcqRel() {
	arch.BarrierAcqRel()
}
