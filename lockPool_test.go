package bufpool

import "testing"

func TestLockPool(t *testing.T) {
	bp, _ := NewLockPool()
	test(t, bp)
}

func BenchmarkLowLockPool(b *testing.B) {
	bp, _ := NewLockPool()
	bench(b, bp, lowConcurrency)
}

func BenchmarkMedLockPool(b *testing.B) {
	bp, _ := NewLockPool()
	bench(b, bp, medConcurrency)
}

func BenchmarkHighLockPool(b *testing.B) {
	bp, _ := NewLockPool()
	bench(b, bp, highConcurrency)
}

func BenchmarkHighLockPoolRuthless(b *testing.B) {
	bp, _ := NewLockPool()
	benchRuthless(b, bp, highConcurrency)
}
