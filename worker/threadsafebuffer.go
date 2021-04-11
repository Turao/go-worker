package worker

import (
	"bytes"
	"sync"
)

// threadSafeBuffer decorates a buffer to provide thread-safe read/write operations
type threadSafeBuffer struct {
	mx  *sync.RWMutex
	buf *bytes.Buffer
}

func (b *threadSafeBuffer) Read(p []byte) (n int, err error) {
	b.mx.RLock()
	defer b.mx.RUnlock()
	return b.buf.Read(p)
}

func (b *threadSafeBuffer) Write(p []byte) (n int, err error) {
	b.mx.Lock()
	defer b.mx.Unlock()
	return b.buf.Write(p)
}

func (b *threadSafeBuffer) String() string {
	b.mx.RLock()
	defer b.mx.RUnlock()
	return b.buf.String()
}
