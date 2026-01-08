// ©Hayabusa Cloud Co., Ltd. 2026. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package atomix_test

import (
	"runtime"
	"sync"
	"testing"

	"code.hybscloud.com/atomix"
)

// =============================================================================
// High Contention Stress Tests
// =============================================================================

// TestStressUint64HighContention tests with maximum goroutines on available cores
func TestStressUint64HighContention(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	numCPU := runtime.NumCPU()
	numGoroutines := numCPU * 8 // 8x oversubscription for high contention
	opsPerGoroutine := 10000

	var a atomix.Uint64
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				a.Add(1)
			}
		}()
	}

	wg.Wait()

	expected := uint64(numGoroutines * opsPerGoroutine)
	if got := a.Load(); got != expected {
		t.Fatalf("lost updates: got %d, want %d (diff: %d)", got, expected, expected-got)
	}

	t.Logf("stress test: %d goroutines x %d ops = %d total ops (PASS)", numGoroutines, opsPerGoroutine, expected)
}

// TestStressUint128HighContention tests 128-bit atomics under extreme contention
func TestStressUint128HighContention(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)

	numCPU := runtime.NumCPU()
	numGoroutines := numCPU * 8
	opsPerGoroutine := 5000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				a.Inc()
			}
		}()
	}

	wg.Wait()

	expected := int64(numGoroutines * opsPerGoroutine)
	lo, hi := a.Load()
	if lo != uint64(expected) || hi != 0 {
		t.Fatalf("lost updates: got (%d, %d), want (%d, 0)", lo, hi, expected)
	}

	t.Logf("stress test: %d goroutines x %d ops = %d total ops (PASS)", numGoroutines, opsPerGoroutine, expected)
}

// TestStressCASContention tests CAS under high contention (high failure rate)
func TestStressCASContention(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	var counter atomix.Uint64
	var casAttempts, casSuccesses atomix.Uint64

	numCPU := runtime.NumCPU()
	numGoroutines := numCPU * 4
	target := uint64(100000)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for {
				old := counter.Load()
				if old >= target {
					return
				}
				casAttempts.Add(1)
				if counter.CompareAndSwap(old, old+1) {
					casSuccesses.Add(1)
				}
			}
		}()
	}

	wg.Wait()

	attempts := casAttempts.Load()
	successes := casSuccesses.Load()
	failureRate := float64(attempts-successes) / float64(attempts) * 100

	if counter.Load() != target {
		t.Fatalf("counter: got %d, want %d", counter.Load(), target)
	}

	t.Logf("CAS contention: %d attempts, %d successes, %.1f%% failure rate",
		attempts, successes, failureRate)
}

// TestStressMixedOperations tests different atomic operations concurrently
func TestStressMixedOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	var a atomix.Uint64
	numCPU := runtime.NumCPU()
	numGoroutines := numCPU * 4
	iterations := 5000

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // 3 types of operations (Add, Sub, Load)

	// Adders
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				a.Add(1)
			}
		}()
	}

	// Subtractors
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				a.Sub(1)
			}
		}()
	}

	// Readers (verify Load works under contention)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = a.Load()
			}
		}()
	}

	wg.Wait()

	// Result should be 0 since adds and subs cancel out
	// Note: Swap is NOT included because Swap(Load()) is a non-atomic
	// read-modify-write that can lose updates from concurrent Add/Sub.
	if got := a.Load(); got != 0 {
		t.Fatalf("mixed ops: got %d, want 0", got)
	}

	t.Logf("mixed operations: %d goroutines × 3 operation types × %d iterations (PASS)",
		numGoroutines, iterations)
}

// TestStressSwapContention tests Swap under high contention
func TestStressSwapContention(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	var a atomix.Uint64
	numCPU := runtime.NumCPU()
	numGoroutines := numCPU * 4
	iterations := 5000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Each goroutine swaps in its ID, we just verify no crashes
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				old := a.Swap(uint64(id))
				// Verify old value is a valid ID (0 to numGoroutines-1) or 0
				if old != 0 && old >= uint64(numGoroutines) {
					t.Errorf("invalid swap result: %d", old)
					return
				}
			}
		}(i)
	}

	wg.Wait()

	// Final value should be one of the goroutine IDs
	final := a.Load()
	if final >= uint64(numGoroutines) {
		t.Fatalf("invalid final value: %d", final)
	}

	t.Logf("swap contention: %d goroutines × %d swaps (PASS)", numGoroutines, iterations)
}

// BenchmarkStressUint64Parallel measures throughput under parallel stress
func BenchmarkStressUint64Parallel(b *testing.B) {
	var a atomix.Uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.Add(1)
		}
	})
}

// BenchmarkStressUint128Parallel measures 128-bit atomic throughput under parallel stress
func BenchmarkStressUint128Parallel(b *testing.B) {
	buf := make([]byte, 64)
	_, a := atomix.PlaceAlignedUint128(buf, 0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.Inc()
		}
	})
}

// BenchmarkStressCASParallel measures CAS throughput (high failure expected)
func BenchmarkStressCASParallel(b *testing.B) {
	var a atomix.Uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for {
				old := a.Load()
				if a.CompareAndSwap(old, old+1) {
					break
				}
			}
		}
	})
}
