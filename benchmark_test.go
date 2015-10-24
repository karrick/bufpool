package bufpool

import (
	"bytes"
	"sync"
	"testing"

	"github.com/oxtoacart/bpool"
)

type GetPuter interface {
	Get() *bytes.Buffer
	Put(*bytes.Buffer)
}

func bench(b *testing.B, bp GetPuter) {
	byteCount := DefaultBufferSize / 2
	const concurrency = 1000
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for a := 0; a < concurrency; a++ {
		go func() {
			for j := 0; j < b.N; j++ {
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

func BenchmarkBufPool(b *testing.B) {
	bp, _ := New() // use defaults
	bench(b, bp)
}

func BenchmarkBpool(b *testing.B) {
	// use bufpool default sizes when creating bpools
	bp := bpool.NewSizedBufferPool(DefaultPoolSize, DefaultBufferSize)
	bench(b, bp)
}
