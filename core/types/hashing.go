package types

import (
	"bytes"
	"fmt"
	"math"
	"sync"
)

// encodeBufferPool holds temporary encoder buffers for DeriveSha and TX encoding.
var encodeBufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// getPooledBuffer retrieves a buffer from the pool and creates a byte slice of the
// requested size from it.
//
// The caller should return the *bytes.Buffer object back into encodeBufferPool after use!
// The returned byte slice must not be used after returning the buffer.
func getPooledBuffer(size uint64) ([]byte, *bytes.Buffer, error) {
	if size > math.MaxInt {
		return nil, nil, fmt.Errorf("can't get buffer of size %d", size)
	}
	buf := encodeBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Grow(int(size))
	b := buf.Bytes()[:int(size)]
	return b, buf, nil
}
