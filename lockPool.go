package bufpool

import (
	"bytes"
	"fmt"
	"sync"
)

// LockPool maintains a free-list of buffers.
type LockPool struct {
	lock sync.Mutex
	free []*bytes.Buffer
	pc   poolConfig
}

// NewLockPool creates a new FreeList. The pool size, size of new buffers, and max size of buffers
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
//        	bp, err := bufpool.NewLockPool()
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
func NewLockPool(setters ...func(*poolConfig) error) (FreeList, error) {
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
	bp := &LockPool{
		free: make([]*bytes.Buffer, 0, pc.chSize),
		pc:   *pc,
	}
	return bp, nil
}

// Get returns an initialized buffer from the free-list.
func (bp *LockPool) Get() *bytes.Buffer {
	bp.lock.Lock()

	if len(bp.free) == 0 {
		bp.lock.Unlock()
		return bytes.NewBuffer(make([]byte, 0, bp.pc.defSize))
	}
	var bb *bytes.Buffer
	lmo := len(bp.free) - 1
	bb, bp.free = bp.free[lmo], bp.free[:lmo]
	bp.lock.Unlock()
	return bb
}

// Put will return a used buffer back to the free-list. If the capacity of the used buffer grew
// beyond the max buffer size, it will be discarded and its memory returned to the runtime.
func (bp *LockPool) Put(bb *bytes.Buffer) {
	if cap(bb.Bytes()) > bp.pc.maxSize {
		return // drop buffer on floor if too big
	}

	bp.lock.Lock()

	if len(bp.free) == cap(bp.free) {
		bp.lock.Unlock()
		return // drop buffer on floor if already have enough
	}
	bb.Reset()
	bp.free = append(bp.free, bb)
	bp.lock.Unlock()
}

// Reset releases memory for all buffers presently in the free-list back to the runtime. This method
// is typically not called for long-running programs that use a free-list of buffers for a long
// time.
func (bp *LockPool) Reset() {
	bp.lock.Lock()
	bp.free = bp.free[:0]
	bp.lock.Unlock()
}
