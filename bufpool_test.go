package bufpool

import (
	"sync"
	"testing"
)

func testC(bp FreeList, concurrency, loops int) {
	const byteCount = DefaultBufferSize / 2

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

func test(t *testing.T, bp FreeList) {
	const concurrency = 128
	const loops = 128
	testC(bp, concurrency, loops)
}

const lowConcurrency = 16
const medConcurrency = 128
const highConcurrency = 1024

func bench(b *testing.B, bp FreeList, concurrency int) {
	const loops = 1024
	for i := 0; i < b.N; i++ {
		testC(bp, concurrency, loops)
	}
}
