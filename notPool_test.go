package bufpool

import (
	"bytes"
	"fmt"
	"testing"
)

type notPool struct {
	pc poolConfig
}

func NewNotPool(setters ...Configurator) (FreeList, error) {
	pc := &poolConfig{
		poolSize: DefaultPoolSize,
		bufSize:  DefaultBufSize,
		maxKeep:  DefaultMaxKeep,
	}
	for _, setter := range setters {
		if err := setter(pc); err != nil {
			return nil, err
		}
	}
	if pc.maxKeep < pc.bufSize {
		return nil, fmt.Errorf("max buffer size must be greater or equal to default buffer size: %d, %d", pc.maxKeep, pc.bufSize)
	}
	bp := &notPool{pc: *pc}
	return bp, nil
}

func (bp *notPool) Get() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, bp.pc.bufSize))
}

func (bp *notPool) Put(*bytes.Buffer) {
}

func TestNotPool(t *testing.T) {
	bp, _ := NewNotPool()
	test(t, bp)
}

func BenchmarkNotPoolLowConcurrency(b *testing.B) {
	bp, _ := NewNotPool()
	bench(b, bp, lowConcurrency)
}

func BenchmarkNotPoolMedConcurrency(b *testing.B) {
	bp, _ := NewNotPool()
	bench(b, bp, medConcurrency)
}

func BenchmarkNotPoolHighConcurrency(b *testing.B) {
	bp, _ := NewNotPool()
	bench(b, bp, highConcurrency)
}
