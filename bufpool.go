package bufpool

import (
	"bytes"
	"fmt"
)

// DefaultPoolSize is the default number of buffers that the free list will maintain.
const DefaultPoolSize = 100

// DefaultBufferSize is the default size used to create new buffers.
const DefaultBufferSize = 4 * 1024

// DefaultMaxBufferSize is the default size used to determine whether to keep buffers returned to
// the pool.
const DefaultMaxBufferSize = 16 * 1024

// BufPool maintains a free list of buffers.
type BufPool struct {
	ch                       chan *bytes.Buffer
	chSize, defSize, maxSize int
}

// New creates a new BufPool instance. The pool size, size of new buffers, and max size of buffers
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
//        	bp, err := bufpool.New() // can have one or more of PoolSize(), BufferSize, and MaxSize()
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
func New(setters ...func(*BufPool) error) (*BufPool, error) {
	bp := &BufPool{
		chSize:  DefaultPoolSize,
		defSize: DefaultBufferSize,
		maxSize: DefaultMaxBufferSize,
	}
	for _, setter := range setters {
		if err := setter(bp); err != nil {
			return nil, err
		}
	}
	if bp.maxSize < bp.defSize {
		return nil, fmt.Errorf("max buffer size must be greater or equal to default buffer size: %d, %d", bp.maxSize, bp.defSize)
	}
	bp.ch = make(chan *bytes.Buffer, bp.chSize)
	return bp, nil
}

// Get returns an initialized buffer from the BufPool free list.
func (bp *BufPool) Get() *bytes.Buffer {
	select {
	case bb := <-bp.ch:
		// reuse buffer
		return bb
	default:
		// empty channel: create new buffer
		return bytes.NewBuffer(make([]byte, 0, bp.defSize))
	}
}

// Put will return a used buffer back to the BufPool free list. If the capacity of the used buffer
// grew beyond the BufPool's max buffer size, it will be discarded and its memory returned to the
// runtime.
func (bp *BufPool) Put(bb *bytes.Buffer) {
	if cap(bb.Bytes()) > bp.maxSize {
		// drop buffer on floor if too big
		return
	}
	bb.Reset()
	select {
	case bp.ch <- bb:
		// queue buffer for reuse
	default:
		// drop on floor if channel full
	}
}

// Reset releases memory for all buffers presently in the BufPool back to the runtime. This method
// is typically not called for long-running programs that use a free list of buffers for a long
// time.
func (bp *BufPool) Reset() {
	for {
		select {
		case _ = <-bp.ch:
			// dropbuffer
		default:
			// empty channel
			return
		}
	}
}

// PoolSize specifies the number of buffers to maintain in the pool.
//        package main
//
//        import (
//        	"log"
//
//        	"github.com/karrick/bufpool"
//        )
//
//        func main() {
//        	bp, err := bufpool.New(bufpool.PoolSize(25))
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
func PoolSize(size int) func(*BufPool) error {
	return func(bp *BufPool) error {
		if size <= 0 {
			return fmt.Errorf("pool size must be greater than 0: %d", size)
		}
		bp.chSize = size
		return nil
	}
}

// BufferSize specifies the size of newly allocated buffers.
//        package main
//
//        import (
//        	"log"
//
//        	"github.com/karrick/bufpool"
//        )
//
//        func main() {
//        	bp, err := bufpool.New(bufpool.BufferSize(1024))
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
func BufferSize(size int) func(*BufPool) error {
	return func(bp *BufPool) error {
		if size <= 0 {
			return fmt.Errorf("default buffer size must be greater than 0: %d", size)
		}
		bp.defSize = size
		return nil
	}
}

// MaxSize specifies the maximum size of buffers that ought to be kept when returned to the free
// list.  Buffers with a capacity larger than this size will be discarded, and their memory returned
// to the runtime.
//        package main
//
//        import (
//        	"log"
//
//        	"github.com/karrick/bufpool"
//        )
//
//        func main() {
//        	bp, err := bufpool.New(bufpool.MaxSize(32 * 1024))
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
func MaxSize(size int) func(*BufPool) error {
	return func(bp *BufPool) error {
		if size <= 0 {
			return fmt.Errorf("max buffer size must be greater than 0: %d", size)
		}
		bp.maxSize = size
		return nil
	}
}
