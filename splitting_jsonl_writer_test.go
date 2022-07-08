package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

type BufWriteCloser struct {
	*bytes.Buffer
}

func NewBufWriteCloser() *BufWriteCloser {
	return &BufWriteCloser{
		bytes.NewBuffer(nil),
	}
}

func (bwc *BufWriteCloser) Close() error {
	return nil
}

func TestSplittingJSONLWriter(t *testing.T) {
	var buffers []*BufWriteCloser

	createWriter := func() (io.WriteCloser, error) {
		buf := NewBufWriteCloser()
		buffers = append(buffers, buf)
		return buf, nil
	}

	const splitSize = 128
	wr := NewSplittingJsonlWriter(createWriter, splitSize)

	const numItmes = 1024
	item := `{"k": "val", "l": [1,2,3,4,5,6,7,8,9,10,11,12]}`
	itemLen := len(item)
	itemsPerFile := (splitSize + (itemLen - 1)) / itemLen
	var lines []string
	for i := 0; i < itemsPerFile; i++ {
		lines = append(lines, item)
	}
	expectedVal := strings.Join(lines, "\n")

	expectedFileCount := (numItmes + (itemsPerFile - 1)) / itemsPerFile
	for i := 0; i < numItmes; i++ {
		err := wr.Add([]byte(item))
		require.NoError(t, err)
	}

	require.Len(t, buffers, expectedFileCount)

	for i := 0; i < len(buffers)-1; i++ {
		buf := buffers[i]
		bs := buf.Buffer.Bytes()
		require.Equal(t, expectedVal, string(bs))
	}
}
