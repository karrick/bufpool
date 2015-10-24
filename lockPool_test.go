package bufpool

import "testing"

func TestLockPool(t *testing.T) {
	bp, _ := NewLockPool()
	test(t, bp)
}

func BenchmarkLockPoolLowConcurrency(b *testing.B) {
	bp, _ := NewLockPool()
	bench(b, bp, lowConcurrency)
}

func BenchmarkLockPoolMedConcurrency(b *testing.B) {
	bp, _ := NewLockPool()
	bench(b, bp, medConcurrency)
}

func BenchmarkLockPoolHighConcurrency(b *testing.B) {
	bp, _ := NewLockPool()
	bench(b, bp, highConcurrency)
}
