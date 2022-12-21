package jsplit

import (
	"errors"
	"io"
)

var newLineBytes = []byte{byte('\n')}

// CreateWriterFn is used for creating new io.WriteCloser objects for writing the split json files
type CreateWriterFn func() (io.WriteCloser, error)

// SplittingJsonlWriter receives json objects one at a time, and it writes these objects in jsonl format to a series of
// files closing streams and creating new ones any time a size threshold is reached
type SplittingJsonlWriter struct {
	wr           io.WriteCloser
	createWriter CreateWriterFn
	splitSize    uint64
	writtenBytes uint64
	writtenItems int
} // repacked by gopium

// NewSplittingJsonlWriter returns a *SplittingJsonlWriter which creates streams using the supplied function.  These streams
// are closed and new ones created any time a stream has had more than splitSize bytes written to it
func NewSplittingJsonlWriter(createWriter CreateWriterFn, splitSize uint64) *SplittingJsonlWriter {
	return &SplittingJsonlWriter{
		createWriter: createWriter,
		splitSize:    splitSize,
		writtenBytes: 0,
		writtenItems: 0,
	}
}

// Add adds a new json list to be written to the current stream
func (sjwr *SplittingJsonlWriter) Add(item []byte) error {
	if sjwr.wr == nil {
		err := sjwr.newWriter()
		if err != nil {
			return err
		}
	}

	if sjwr.writtenItems != 0 {
		n, err := sjwr.wr.Write(newLineBytes)
		if err != nil {
			return err
		} else if n != 1 {
			return errors.New("failed to write newline")
		}
	}

	written := 0
	toWrite := len(item)

	for written < toWrite {
		n, err := sjwr.wr.Write(item[written:])
		if err != nil {
			return err
		}

		written += n
	}

	sjwr.writtenItems++
	sjwr.writtenBytes += uint64(len(item))

	if sjwr.writtenBytes >= sjwr.splitSize {
		err := sjwr.newWriter()
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes the last stream making sure all the data has been flushed
func (sjwr *SplittingJsonlWriter) Close() error {
	if sjwr.wr != nil {
		err := sjwr.wr.Close()
		if err != nil {
			return err
		}

		sjwr.wr = nil
		sjwr.writtenBytes = 0
		sjwr.writtenItems = 0
	}

	return nil
}

func (sjwr *SplittingJsonlWriter) newWriter() error {
	if sjwr.wr != nil {
		err := sjwr.Close()
		if err != nil {
			return err
		}
	}

	newWr, err := sjwr.createWriter()
	if err != nil {
		return err
	}

	sjwr.wr = newWr
	sjwr.writtenItems = 0
	sjwr.writtenBytes = 0

	return nil
}
