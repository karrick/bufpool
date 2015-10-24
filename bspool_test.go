package bufpool

import (
	"testing"

	"github.com/oxtoacart/bpool"
)

// NOTE: when creating bpool instances, use default sized used by bufpool

func TestBspool(t *testing.T) {
	bp := bpool.NewSizedBufferPool(DefaultPoolSize, DefaultBufferSize)
	test(t, bp)
}

func BenchmarkBpoolLowConcurrency(b *testing.B) {
	bp := bpool.NewSizedBufferPool(DefaultPoolSize, DefaultBufferSize)
	bench(b, bp, lowConcurrency)
}

func BenchmarkBpoolMedConcurrency(b *testing.B) {
	bp := bpool.NewSizedBufferPool(DefaultPoolSize, DefaultBufferSize)
	bench(b, bp, medConcurrency)
}

func BenchmarkBpoolHighConcurrency(b *testing.B) {
	bp := bpool.NewSizedBufferPool(DefaultPoolSize, DefaultBufferSize)
	bench(b, bp, highConcurrency)
}