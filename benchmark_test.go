package bufpool

import (
	"bytes"
	"sync"
	"testing"

	"github.com/oxtoacart/bpool"
)

type BufferFreeList interface {
	Get() *bytes.Buffer
	Put(*bytes.Buffer)
}

func bench(b *testing.B, bp BufferFreeList, concurrency int) {
	const byteCount = DefaultBufferSize / 2
	const loops = 1024

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for c := 0; c < concurrency; c++ {
			go func() {
				for j := 0; j < loops; j++ {
					bb := bp.Get()
					max := byteCount
					if j%8 == 0 {
						max = 2 * DefaultMaxBufferSize
					}
					for k := 0; k < max; k++ {
						bb.WriteByte(byte(k % 256))
					}
					bp.Put(bb)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

const lowConcurrency = 16
const medConcurrency = 128
const highConcurrency = 1024

func BenchmarkBpoolLowConcurrency(b *testing.B) {
	// use bufpool default sizes when creating bpools
	bp := bpool.NewSizedBufferPool(DefaultPoolSize, DefaultBufferSize)
	bench(b, bp, lowConcurrency)
}

func BenchmarkBpoolMedConcurrency(b *testing.B) {
	// use bufpool default sizes when creating bpools
	bp := bpool.NewSizedBufferPool(DefaultPoolSize, DefaultBufferSize)
	bench(b, bp, medConcurrency)
}

func BenchmarkBpoolHighConcurrency(b *testing.B) {
	// use bufpool default sizes when creating bpools
	bp := bpool.NewSizedBufferPool(DefaultPoolSize, DefaultBufferSize)
	bench(b, bp, highConcurrency)
}

func BenchmarkBufPoolLowConcurrency(b *testing.B) {
	bp, _ := New() // use defaults
	bench(b, bp, lowConcurrency)
}

func BenchmarkBufPoolMedConcurrency(b *testing.B) {
	bp, _ := New() // use defaults
	bench(b, bp, medConcurrency)
}

func BenchmarkBufPoolHighConcurrency(b *testing.B) {
	bp, _ := New() // use defaults
	bench(b, bp, highConcurrency)
}

func BenchmarkBufPoolLockLowConcurrency(b *testing.B) {
	bp, _ := NewLock() // use defaults
	bench(b, bp, lowConcurrency)
}

func BenchmarkBufPoolLockMedConcurrency(b *testing.B) {
	bp, _ := NewLock() // use defaults
	bench(b, bp, medConcurrency)
}

func BenchmarkBufPoolLockHighConcurrency(b *testing.B) {
	bp, _ := NewLock() // use defaults
	bench(b, bp, highConcurrency)
}
