package bufpool

import (
	"bytes"
	"fmt"
	"sync"
)

type BufPoolLock struct {
	lock                     sync.Mutex
	free                     []*bytes.Buffer
	chSize, defSize, maxSize int
}

func NewLock(setters ...func(*BufPoolLock) error) (*BufPoolLock, error) {
	bp := &BufPoolLock{
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
	bp.free = make([]*bytes.Buffer, 0, bp.chSize)
	return bp, nil
}

func (bp *BufPoolLock) Get() *bytes.Buffer {
	bp.lock.Lock()
	defer bp.lock.Unlock()

	if len(bp.free) == 0 {
		return bytes.NewBuffer(make([]byte, 0, bp.defSize))
	}
	var bb *bytes.Buffer
	bb, bp.free = bp.free[len(bp.free)-1], bp.free[:len(bp.free)-1]
	return bb
}

func (bp *BufPoolLock) Put(bb *bytes.Buffer) {
	if cap(bb.Bytes()) > bp.maxSize {
		return // drop buffer on floor if too big
	}

	bp.lock.Lock()
	defer bp.lock.Unlock()

	if len(bp.free) == cap(bp.free) {
		return // drop buffer on floor if already have enough
	}
	bb.Reset()
	bp.free = append(bp.free, bb)
}

func (bp *BufPoolLock) Reset() {
	bp.lock.Lock()
	defer bp.lock.Unlock()

	bp.free = bp.free[:0]
}
