package bufpool

import (
	"bytes"
	"sync"
	"testing"
)

func testC(bp FreeList, concurrency, loops int) {
	const byteCount = DefaultBufSize / 2

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for c := 0; c < concurrency; c++ {
		go func() {
			for j := 0; j < loops; j++ {
				bb := bp.Get()
				max := byteCount
				if j%8 == 0 {
					max = 2 * DefaultMaxKeep
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
	b.ResetTimer()
	const loops = 1024
	for i := 0; i < b.N; i++ {
		testC(bp, concurrency, loops)
	}
}

func benchRuthless(b *testing.B, bp FreeList, concurrency int) {
	b.ReportAllocs()
	b.ResetTimer()
	const loops = 1024
	for i := 0; i < b.N; i++ {
		testRuthless(bp, concurrency, loops)
	}
}

func testRuthless(bp FreeList, concurrency, loops int) {
	// NOTE: create task that keeps filling free list with new buffers to keep it full
	halt := make(chan struct{})
	go func(halt chan struct{}) {
		for {
			select {
			case <-halt:
				break
			default:
				bp.Put(bytes.NewBuffer(make([]byte, 0, DefaultBufSize)))
			}
		}
	}(halt)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < loops; j++ {
				bb := bp.Get()
				// ensure each buffer sent back is a little too big
				for k := 0; k < DefaultBufSize+1; k++ {
					bb.WriteByte(byte(k % 256))
				}
				bp.Put(bb)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(halt)
}
