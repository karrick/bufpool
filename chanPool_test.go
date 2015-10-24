package bufpool

import "testing"

func TestChanPool(t *testing.T) {
	bp, _ := NewChanPool()
	test(t, bp)
}

func BenchmarkChanPoolLowConcurrency(b *testing.B) {
	bp, _ := NewChanPool()
	bench(b, bp, lowConcurrency)
}

func BenchmarkChanPoolMedConcurrency(b *testing.B) {
	bp, _ := NewChanPool()
	bench(b, bp, medConcurrency)
}

func BenchmarkChanPoolHighConcurrency(b *testing.B) {
	bp, _ := NewChanPool()
	bench(b, bp, highConcurrency)
}
