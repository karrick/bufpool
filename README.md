# bufpool

Go library for using a free list of byte buffers.

### Usage

Documentation is available via
[![GoDoc](https://godoc.org/github.com/karrick/bufpool?status.svg)](https://godoc.org/github.com/karrick/bufpool).

# Example

```Go
    package main
    
    import (
    	"log"
    
    	"github.com/karrick/bufpool"
    )
    
    func main() {
    	bp, err := bufpool.New()
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
