package bufpool

import "testing"

func TestNew(t *testing.T) {
	bp, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4*DefaultPoolSize; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				bb := bp.Get()
				for k := 0; k < 3*DefaultBufferSize; k++ {
					bb.WriteByte(byte(k % 256))
				}
				bp.Put(bb)
			}
		}()
	}
}
