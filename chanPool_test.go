package bufpool

import "testing"

func TestChanPool(t *testing.T) {
	bp, _ := NewChanPool()
	test(t, bp)
}

func BenchmarkLowChanPool(b *testing.B) {
	bp, _ := NewChanPool()
	bench(b, bp, lowConcurrency)
}

func BenchmarkMedChanPool(b *testing.B) {
	bp, _ := NewChanPool()
	bench(b, bp, medConcurrency)
}

func BenchmarkHighChanPool(b *testing.B) {
	bp, _ := NewChanPool()
	bench(b, bp, highConcurrency)
}

func BenchmarkHighChanPoolRuthless(b *testing.B) {
	bp, _ := NewChanPool()
	benchRuthless(b, bp, highConcurrency)
}
