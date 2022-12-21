package jsplit

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsyncReaderReadsAll(t *testing.T) {
	const maxSize = 64 * 1024

	buffer := make([]byte, maxSize)
	n, err := rand.Read(buffer)
	require.NoError(t, err)
	require.Equal(t, n, maxSize)

	tests := []struct {
		name       string
		size       int
		bufferSize int
	}{
		{
			name:       "64k 32 byte buffer",
			size:       maxSize,
			bufferSize: 32,
		},
		{
			name:       "64k 1 byte buffer",
			size:       maxSize,
			bufferSize: 1,
		},
		{
			name:       "64k 16K byte buffer",
			size:       maxSize,
			bufferSize: 16 * 1024,
		},
		{
			name:       "buffer larger than data",
			size:       128,
			bufferSize: 1024,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			testBuffer := buffer[:test.size]
			rd, err := AsyncReaderFromReader(bytes.NewReader(testBuffer), test.bufferSize)
			require.NoError(t, err)
			ctx = rd.Start(ctx)

			read := make([]byte, 0, test.size)
			for {
				newlyRead, err := rd.Read(ctx)
				if err == io.EOF {
					break
				}

				require.NoError(t, err)

				read = append(read, newlyRead...)
			}

			require.Equal(t, read, testBuffer)
			require.NoError(t, ctx.Err())
			require.True(t, rd.IsClosed())
		})
	}
}

type ErroringReader struct {
	io.Reader
	err            error
	errAfterNReads int
	reads          int
}

func (er *ErroringReader) Read(b []byte) (int, error) {
	if er.reads < er.errAfterNReads {
		er.reads++
		return er.Reader.Read(b)
	}

	return 0, er.err
}

func TestReadError(t *testing.T) {
	const size = 16 * 1024
	const errAfterNReads = 4

	buffer := make([]byte, size)
	n, err := rand.Read(buffer)
	require.NoError(t, err)
	require.Equal(t, n, size)

	expectedErr := errors.New("test error")
	er := &ErroringReader{
		Reader:         bytes.NewReader(buffer),
		err:            expectedErr,
		errAfterNReads: errAfterNReads,
	}

	rd, err := AsyncReaderFromReader(er, 32)
	require.NoError(t, err)
	ctx := rd.Start(context.Background())

	for i := 0; i < errAfterNReads; i++ {
		_, err = rd.Read(ctx)
		if err != nil {
			break
		}
	}

	require.Equal(t, err, expectedErr)
}
