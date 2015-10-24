package bufpool

import (
	"bytes"
	"fmt"
)

const DefaultPoolSize = 100
const DefaultBufferSize = 4 * 1024
const DefaultMaxBufferSize = 16 * 1024

type BufPool struct {
	ch                       chan *bytes.Buffer
	chSize, defSize, maxSize int
}

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

func PoolSize(size int) func(*BufPool) error {
	return func(bp *BufPool) error {
		if size <= 0 {
			return fmt.Errorf("pool size must be greater than 0: %d", size)
		}
		bp.chSize = size
		return nil
	}
}

func BufferSize(size int) func(*BufPool) error {
	return func(bp *BufPool) error {
		if size <= 0 {
			return fmt.Errorf("default buffer size must be greater than 0: %d", size)
		}
		bp.defSize = size
		return nil
	}
}

func MaxSize(size int) func(*BufPool) error {
	return func(bp *BufPool) error {
		if size <= 0 {
			return fmt.Errorf("max buffer size must be greater than 0: %d", size)
		}
		bp.maxSize = size
		return nil
	}
}
