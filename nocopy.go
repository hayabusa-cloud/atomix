// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix

// noCopy may be added to structs which must not be copied after first use.
// See https://golang.org/issues/8005#issuecomment-190753527
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock() {}

// Unlock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Unlock() {}
