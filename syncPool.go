package bufpool

import (
	"bytes"
	"fmt"
	"sync"
)

// SyncPool maintains a free-list of buffers.
type SyncPool struct {
	pool sync.Pool
	pc   poolConfig
}

// NewSyncPool creates a new FreeList. The pool size, size of new buffers, and max size of buffers
// to keep when returned to the pool can all be customized.
//
//        package main
//
//        import (
//        	"log"
//
//        	"github.com/karrick/bufpool"
//        )
//
//        func main() {
//        	bp, err := bufpool.NewSyncPool()
//        	if err != nil {
//        		log.Fatal(err)
//        	}
//        	for i := 0; i < 4*bufpool.DefaultPoolSize; i++ {
//        		go func() {
//        			for j := 0; j < 1000; j++ {
//        				bb := bp.Get()
//        				for k := 0; k < 3*bufpool.DefaultBufferSize; k++ {
//        					bb.WriteByte(byte(k % 256))
//        				}
//        				bp.Put(bb)
//        			}
//        		}()
//        	}
//        }
func NewSyncPool(setters ...func(*poolConfig) error) (FreeList, error) {
	pc := &poolConfig{
		chSize:  DefaultPoolSize,
		defSize: DefaultBufferSize,
		maxSize: DefaultMaxBufferSize,
	}
	for _, setter := range setters {
		if err := setter(pc); err != nil {
			return nil, err
		}
	}
	if pc.maxSize < pc.defSize {
		return nil, fmt.Errorf("max buffer size must be greater or equal to default buffer size: %d, %d", pc.maxSize, pc.defSize)
	}
	bp := &SyncPool{pc: *pc}
	bp.pool.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, bp.pc.defSize))
	}
	return bp, nil
}

// Get returns an initialized buffer from the free-list.
func (bp *SyncPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

// Put will return a used buffer back to the free-list. If the capacity of the used buffer grew
// beyond the max buffer size, it will be discarded and its memory returned to the runtime.
func (bp *SyncPool) Put(bb *bytes.Buffer) {
	if cap(bb.Bytes()) > bp.pc.maxSize {
		return // drop buffer on floor if too big
	}
	bb.Reset()
	bp.pool.Put(bb)
}
