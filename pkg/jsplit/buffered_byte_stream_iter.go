package jsplit

import (
	"context"
	"io"
)

// ByteStream is an interface for reading bytes
type ByteStream interface {
	Read(ctx context.Context) ([]byte, error)
}

// BufferedByteStreamIter attempts to efficiently buffer a stream of bytes to be iterated over sequentially
type BufferedByteStreamIter struct {
	stream ByteStream
	ctx    context.Context
	buffer []byte
	pos    int
}

// NewBufferStreamIter returns a *BufferedByteStreamIter for iterating over the bytes of the given byte stream
func NewBufferedStreamIter(readCtx context.Context, stream ByteStream) *BufferedByteStreamIter {
	return &BufferedByteStreamIter{
		stream: stream,
		buffer: nil,
		pos:    0,
		ctx:    readCtx,
	}
}

// Next read the byte at the current position and move the current position forward.  When all bytes have been iterated
// over a call to next will return 0
func (itr *BufferedByteStreamIter) Next() byte {
	if itr.pos >= len(itr.buffer) {
		err := itr.readMore()
		if err != nil && err != io.EOF {
			panic(err)
		} else if err == io.EOF {
			return 0
		}
	}

	ch := itr.buffer[itr.pos]
	itr.pos++

	return ch
}

// Advance moves the current position forward n places for positive numbers, and back n places for negative numbers
func (itr *BufferedByteStreamIter) Advance(n int) {
	if n > 0 {
		itr.buffer = itr.buffer[n:]
		itr.pos -= n
	} else {
		itr.pos += n
	}
}

// Skip moves the start of the buffer to the current position, and then sets the current position to 0
func (itr *BufferedByteStreamIter) Skip() {
	itr.buffer = itr.buffer[itr.pos:]
	itr.pos = 0
}

// Value gets the slice of bytes from the start of the buffer til the current position and returns it after moving the
// start of the buffer to be the current position, and then sets the current position to 0 and then returns
func (itr *BufferedByteStreamIter) Value() []byte {
	val := itr.buffer[:itr.pos]
	itr.buffer = itr.buffer[itr.pos:]
	itr.pos = 0

	return val
}

func (itr *BufferedByteStreamIter) readMore() error {
	buf, err := itr.stream.Read(itr.ctx)
	if err != nil {
		return err
	}

	if len(itr.buffer) == 0 {
		itr.buffer = buf
		itr.pos = 0
	} else {
		oldLen := len(itr.buffer)
		newBufferLen := oldLen + len(buf)
		newBuf := make([]byte, newBufferLen)
		copy(newBuf, itr.buffer)
		copy(newBuf[oldLen:], buf)
		itr.buffer = newBuf
	}

	return nil
}
