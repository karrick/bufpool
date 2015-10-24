package bufpool

import "testing"

func TestSyncPool(t *testing.T) {
	bp, _ := NewSyncPool()
	test(t, bp)
}

func BenchmarkSyncPoolLowConcurrency(b *testing.B) {
	bp, _ := NewSyncPool()
	bench(b, bp, lowConcurrency)
}

func BenchmarkSyncPoolMedConcurrency(b *testing.B) {
	bp, _ := NewSyncPool()
	bench(b, bp, medConcurrency)
}

func BenchmarkSyncPoolHighConcurrency(b *testing.B) {
	bp, _ := NewSyncPool()
	bench(b, bp, highConcurrency)
}
