// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix_test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"code.hybscloud.com/atomix"
)

// =============================================================================
// Contention Benchmark: sync/atomic vs atomix
// =============================================================================

// Medium-High Contention: GOMAXPROCS * 2 goroutines
// High Contention: GOMAXPROCS * 4 goroutines

// -----------------------------------------------------------------------------
// Add64 Benchmarks
// -----------------------------------------------------------------------------

func BenchmarkContentionAdd64_SyncAtomic_MediumHigh(b *testing.B) {
	var counter int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionAdd64_Atomix_MediumHigh(b *testing.B) {
	var counter atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				counter.Add(1)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionAdd64_AtomixRelaxed_MediumHigh(b *testing.B) {
	var counter atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				counter.AddRelaxed(1)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionAdd64_SyncAtomic_High(b *testing.B) {
	var counter int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionAdd64_Atomix_High(b *testing.B) {
	var counter atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				counter.Add(1)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionAdd64_AtomixRelaxed_High(b *testing.B) {
	var counter atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				counter.AddRelaxed(1)
			}
		}()
	}
	wg.Wait()
}

// -----------------------------------------------------------------------------
// CAS64 Benchmarks (CAS loop pattern)
// -----------------------------------------------------------------------------

func BenchmarkContentionCAS64_SyncAtomic_MediumHigh(b *testing.B) {
	var counter int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				for {
					old := atomic.LoadInt64(&counter)
					if atomic.CompareAndSwapInt64(&counter, old, old+1) {
						break
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionCAS64_Atomix_MediumHigh(b *testing.B) {
	var counter atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				for {
					old := counter.Load()
					if counter.CompareAndSwap(old, old+1) {
						break
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionCAS64_AtomixRelaxed_MediumHigh(b *testing.B) {
	var counter atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				for {
					old := counter.LoadRelaxed()
					if counter.CompareAndSwapRelaxed(old, old+1) {
						break
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionCAS64_SyncAtomic_High(b *testing.B) {
	var counter int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				for {
					old := atomic.LoadInt64(&counter)
					if atomic.CompareAndSwapInt64(&counter, old, old+1) {
						break
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionCAS64_Atomix_High(b *testing.B) {
	var counter atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				for {
					old := counter.Load()
					if counter.CompareAndSwap(old, old+1) {
						break
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionCAS64_AtomixRelaxed_High(b *testing.B) {
	var counter atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				for {
					old := counter.LoadRelaxed()
					if counter.CompareAndSwapRelaxed(old, old+1) {
						break
					}
				}
			}
		}()
	}
	wg.Wait()
}

// -----------------------------------------------------------------------------
// Swap64 Benchmarks
// -----------------------------------------------------------------------------

func BenchmarkContentionSwap64_SyncAtomic_MediumHigh(b *testing.B) {
	var val int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				atomic.SwapInt64(&val, int64(id))
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkContentionSwap64_Atomix_MediumHigh(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				val.Swap(int64(id))
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkContentionSwap64_AtomixRelaxed_MediumHigh(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				val.SwapRelaxed(int64(id))
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkContentionSwap64_SyncAtomic_High(b *testing.B) {
	var val int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				atomic.SwapInt64(&val, int64(id))
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkContentionSwap64_Atomix_High(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				val.Swap(int64(id))
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkContentionSwap64_AtomixRelaxed_High(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				val.SwapRelaxed(int64(id))
			}
		}(i)
	}
	wg.Wait()
}

// -----------------------------------------------------------------------------
// Load/Store Mixed Benchmarks (readers + writers)
// -----------------------------------------------------------------------------

func BenchmarkContentionLoadStore_SyncAtomic_MediumHigh(b *testing.B) {
	var val int64
	numG := runtime.GOMAXPROCS(0) * 2
	numWriters := numG / 2
	numReaders := numG - numWriters
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	// Writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				atomic.StoreInt64(&val, int64(id))
			}
		}(i)
	}
	// Readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				_ = atomic.LoadInt64(&val)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionLoadStore_Atomix_MediumHigh(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 2
	numWriters := numG / 2
	numReaders := numG - numWriters
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	// Writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				val.Store(int64(id))
			}
		}(i)
	}
	// Readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				_ = val.Load()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionLoadStore_AtomixRelaxed_MediumHigh(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 2
	numWriters := numG / 2
	numReaders := numG - numWriters
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	// Writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				val.StoreRelaxed(int64(id))
			}
		}(i)
	}
	// Readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				_ = val.LoadRelaxed()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionLoadStore_SyncAtomic_High(b *testing.B) {
	var val int64
	numG := runtime.GOMAXPROCS(0) * 4
	numWriters := numG / 2
	numReaders := numG - numWriters
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	// Writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				atomic.StoreInt64(&val, int64(id))
			}
		}(i)
	}
	// Readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				_ = atomic.LoadInt64(&val)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionLoadStore_Atomix_High(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	numWriters := numG / 2
	numReaders := numG - numWriters
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	// Writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				val.Store(int64(id))
			}
		}(i)
	}
	// Readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				_ = val.Load()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkContentionLoadStore_AtomixRelaxed_High(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	numWriters := numG / 2
	numReaders := numG - numWriters
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	// Writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				val.StoreRelaxed(int64(id))
			}
		}(i)
	}
	// Readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerG; j++ {
				_ = val.LoadRelaxed()
			}
		}()
	}
	wg.Wait()
}

// -----------------------------------------------------------------------------
// Or64 Benchmarks (bitwise operations under contention)
// -----------------------------------------------------------------------------

func BenchmarkContentionOr64_SyncAtomic_High(b *testing.B) {
	var val int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(bit int) {
			defer wg.Done()
			mask := int64(1 << (bit % 64))
			for j := 0; j < opsPerG; j++ {
				atomic.OrInt64(&val, mask)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkContentionOr64_Atomix_High(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(bit int) {
			defer wg.Done()
			mask := int64(1 << (bit % 64))
			for j := 0; j < opsPerG; j++ {
				val.Or(mask)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkContentionOr64_AtomixRelaxed_High(b *testing.B) {
	var val atomix.Int64
	numG := runtime.GOMAXPROCS(0) * 4
	var wg sync.WaitGroup
	opsPerG := b.N / numG
	if opsPerG < 1 {
		opsPerG = 1
	}

	b.ResetTimer()
	for i := 0; i < numG; i++ {
		wg.Add(1)
		go func(bit int) {
			defer wg.Done()
			mask := int64(1 << (bit % 64))
			for j := 0; j < opsPerG; j++ {
				val.OrRelaxed(mask)
			}
		}(i)
	}
	wg.Wait()
}
