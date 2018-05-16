package bufpool

import "testing"

func TestSyncPool(t *testing.T) {
	bp, _ := NewSyncPool()
	test(t, bp)
}

func BenchmarkLowSyncPool(b *testing.B) {
	bp, _ := NewSyncPool()
	bench(b, bp, lowConcurrency)
}

func BenchmarkMedSyncPool(b *testing.B) {
	bp, _ := NewSyncPool()
	bench(b, bp, medConcurrency)
}

func BenchmarkHighSyncPool(b *testing.B) {
	bp, _ := NewSyncPool()
	bench(b, bp, highConcurrency)
}

func BenchmarkHighSyncPoolRuthless(b *testing.B) {
	bp, _ := NewSyncPool()
	benchRuthless(b, bp, highConcurrency)
}
