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

// PoolSize specifies the number of buffers to maintain in the pool.
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
func MaxSize(size int) func(*BufPool) error {
	return func(bp *BufPool) error {
		if size <= 0 {
			return fmt.Errorf("max buffer size must be greater than 0: %d", size)
		}
		bp.maxSize = size
		return nil
	}
}
