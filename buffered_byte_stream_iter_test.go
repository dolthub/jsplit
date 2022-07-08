package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

type TestByteStream struct {
	bytes    []byte
	readSize int
	pos      int
}

func NewTestByteStream(bytes []byte, readSize int) *TestByteStream {
	return &TestByteStream{
		bytes:    bytes,
		readSize: readSize,
	}
}

func (tbs *TestByteStream) Read(ctx context.Context) ([]byte, error) {
	startPos := tbs.pos
	tbs.pos += tbs.readSize

	if startPos >= len(tbs.bytes) {
		return nil, io.EOF
	}

	if tbs.pos > len(tbs.bytes) {
		tbs.pos = len(tbs.bytes)
	}

	readBytes := tbs.bytes[startPos:tbs.pos]
	return readBytes, nil
}

func (tbs *TestByteStream) RequireEqual(t *testing.T, bytes []byte) {
	require.Equal(t, tbs.bytes, bytes)
}

func TestBufferedByteStreamIter(t *testing.T) {
	testStr := "this is a test. I will use it to validate functionality"
	tbs := NewTestByteStream([]byte(testStr), 8)
	itr := NewBufferedStreamIter(tbs, context.Background())

	splitWords := strings.Split(testStr, " ")
	var words []string
	for {
		ch := itr.Next()
		if ch == byte(' ') {
			itr.Advance(-1)
			words = append(words, string(itr.Value()))

			itr.Next()
			itr.Skip()
		} else if ch == 0 {
			break
		}
	}

	words = append(words, string(itr.Value()))

	require.Equal(t, splitWords, words)
}
