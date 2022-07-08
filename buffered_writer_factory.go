package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BufferedWriteCloser wraps an io.WriteCloser in a bufio.Writer object and provides an io.WriteCloser implementation
// for the bufio.Writer object
type BufferedWriteCloser struct {
	name  string
	start time.Time
	wr    io.WriteCloser
	bufWr *bufio.Writer
}

// NewBufferedWriteCloser returns a BufferedWriteCloser object which writes to the supplied io.WriteCloser
func NewBufferedWriteCloser(name string, wr io.WriteCloser, bufferSize int) *BufferedWriteCloser {
	bufWr := bufio.NewWriterSize(wr, bufferSize)
	return &BufferedWriteCloser{
		name:  name,
		start: time.Now(),
		wr:    wr,
		bufWr: bufWr,
	}
}

// Write calls write on the bufio.Writer object which wraps the io.WriterCloser
func (bwc *BufferedWriteCloser) Write(p []byte) (n int, err error) {
	return bwc.bufWr.Write(p)
}

// Close makes sure the bufio.Writer object flushes, and the supplied io.WriteCloser is closed
func (bwc *BufferedWriteCloser) Close() error {
	err := bwc.bufWr.Flush()
	if err != nil {
		return err
	}

	fmt.Printf("Closing %s after %f seconds\n", bwc.name, time.Since(bwc.start).Seconds())

	return bwc.wr.Close()
}

// BufferedWriterFactory returns an object which can be used for creating jsonl files
type BufferedWriterFactory struct {
	format     string
	index      int
	bufferSize int
}

// NewBufferedWriterFactory returns a *BufferedWriterFactory instance which creates files in the format [key]_%02d.jsonl
// within the supplied directory.
func NewBufferedWriterFactory(directory, key string, bufferSize int) *BufferedWriterFactory {
	format := filepath.Join(directory, key+"_%02d.jsonl")
	return &BufferedWriterFactory{
		format:     format,
		index:      0,
		bufferSize: bufferSize,
	}
}

// CreateWriter will create a new file within and return an io.WriteCloser for writing to the newly created file
func (bwf *BufferedWriterFactory) CreateWriter() (io.WriteCloser, error) {
	filename := fmt.Sprintf(bwf.format, bwf.index)
	bwf.index++

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return NewBufferedWriteCloser(filename, f, bwf.bufferSize), nil
}
