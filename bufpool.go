package bufpool

import (
	"bytes"
	"fmt"
)

// DefaultPoolSize is the default number of buffers that the free-list will maintain.
const DefaultPoolSize = 100

// DefaultBufferSize is the default size used to create new buffers.
const DefaultBufferSize = 4 * 1024

// DefaultMaxBufferSize is the default size used to determine whether to keep buffers returned to
// the pool.
const DefaultMaxBufferSize = 16 * 1024

// FreeList represents a data structure that maintains a free-list of buffers, accesible via
// Get and Put methods.
type FreeList interface {
	Get() *bytes.Buffer
	Put(*bytes.Buffer)
}

type poolConfig struct {
	chSize, defSize, maxSize int
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
func PoolSize(size int) func(*poolConfig) error {
	return func(pc *poolConfig) error {
		if size <= 0 {
			return fmt.Errorf("pool size must be greater than 0: %d", size)
		}
		pc.chSize = size
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
func BufferSize(size int) func(*poolConfig) error {
	return func(pc *poolConfig) error {
		if size <= 0 {
			return fmt.Errorf("default buffer size must be greater than 0: %d", size)
		}
		pc.defSize = size
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
func MaxSize(size int) func(*poolConfig) error {
	return func(pc *poolConfig) error {
		if size <= 0 {
			return fmt.Errorf("max buffer size must be greater than 0: %d", size)
		}
		pc.maxSize = size
		return nil
	}
}
