package main

import (
	"log"

	"github.com/karrick/bufpool"
)

func main() {
	bp, err := bufpool.NewChanPool()
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: silly example with heavy resource contension
	for i := 0; i < 4*bufpool.DefaultPoolSize; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				bb := bp.Get()
				// NOTE: buffer is ready to use
				for k := 0; k < bufpool.DefaultBufSize/2; k++ {
					bb.WriteByte(byte(k % 256))
				}
				// NOTE: no need to reset buffer prior to release
				bp.Put(bb)
			}
		}()
	}
}
