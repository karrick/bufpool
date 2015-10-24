# bufpool

Go library for using a free list of byte buffers.

## Description

Several free-list algorithms are included to determine which performs
best for various application scenarios.

* NewChanPool -- uses channels to provide concurrent access to internal structures
* NewLockPool -- uses sync.Mutex to provide concurrent access to internal structures
* NewSyncPool -- uses sync.Pool to provide concurrent access to internal structures

### Usage

Documentation is available via
[![GoDoc](https://godoc.org/github.com/karrick/bufpool?status.svg)](https://godoc.org/github.com/karrick/bufpool).

### Example

```Go
    package main
    
    import (
    	"log"
    
    	"github.com/karrick/bufpool"
    )
    
    func main() {
        // bufpool.New*() can all have one or more of PoolSize(), BufferSize, and MaxSize()
        // to customize the FreeList
    	bp, err := bufpool.NewChanPool()
    	if err != nil {
    		log.Fatal(err)
    	}
    	for i := 0; i < 4*bufpool.DefaultPoolSize; i++ {
    		go func() {
    			for j := 0; j < 1000; j++ {
    				bb := bp.Get()
    				for k := 0; k < 3*bufpool.DefaultBufferSize; k++ {
    					bb.WriteByte(byte(k % 256))
    				}
    				bp.Put(bb)
    			}
    		}()
    	}
    }
```

### Performance

Benchmark functions are provided to determine which buffer free-list
algorithm best suits a given application.
