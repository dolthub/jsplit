package main

import (
	"context"
	"io"
	"os"
	"sync/atomic"
)

// AsyncReader reads an io.Reader asynchronously
type AsyncReader struct {
	readCh     chan []byte
	rd         io.Reader
	bufferSize int
	isClosed   int32
}

// AsyncReaderFromFile creates an AsyncReader for reading the specified file
func AsyncReaderFromFile(filename string, bufferSize int) (*AsyncReader, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return AsyncReaderFromReader(f, bufferSize)
}

// AsyncReaderFromReader returns an AsyncReader for reading the supplied io.Reader
func AsyncReaderFromReader(rd io.Reader, bufferSize int) (*AsyncReader, error) {
	return &AsyncReader{
		readCh:     make(chan []byte, 16),
		rd:         rd,
		bufferSize: bufferSize,
	}, nil
}

// Start starts the background reading of the io.Reader
func (afr *AsyncReader) Start(ctx context.Context) context.Context {
	errCtx, cancelFunc := NewErrContextWithCancel(ctx)

	go func() {
		for {
			buf := make([]byte, afr.bufferSize)
			n, err := afr.rd.Read(buf)

			if err != nil && err != io.EOF {
				cancelFunc(err)
				return
			}

			if n > 0 {
				afr.readCh <- buf[:n]
			}

			if err == io.EOF {
				close(afr.readCh)
				atomic.StoreInt32(&afr.isClosed, 1)
				return
			}
		}
	}()

	return errCtx
}

// Read gets the next chunk which has been read from the file.
func (afr *AsyncReader) Read(ctx context.Context) ([]byte, error) {
	select {
	case buf, ok := <-afr.readCh:
		if !ok {
			return nil, io.EOF
		}

		return buf, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// IsClosed is used for testing to verify that the reader and associated channel has been closed.
func (afr *AsyncReader) IsClosed() bool {
	return atomic.LoadInt32(&afr.isClosed) == 1
}
