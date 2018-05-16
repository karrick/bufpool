package bufpool

import "testing"

func TestPreAllocatedPool(t *testing.T) {
	bp, _ := NewPreAllocatedPool()
	test(t, bp)
}

func BenchmarkLowPreAllocatedPool(b *testing.B) {
	bp, _ := NewPreAllocatedPool()
	bench(b, bp, lowConcurrency)
}

func BenchmarkMedPreAllocatedPool(b *testing.B) {
	bp, _ := NewPreAllocatedPool()
	bench(b, bp, medConcurrency)
}

func BenchmarkHighPreAllocatedPool(b *testing.B) {
	bp, _ := NewPreAllocatedPool()
	bench(b, bp, highConcurrency)
}

func BenchmarkHighPreAllocatedPoolRuthless(b *testing.B) {
	b.Skip("test kills pre-allocated pool")
	bp, _ := NewPreAllocatedPool()
	benchRuthless(b, bp, highConcurrency)
}
