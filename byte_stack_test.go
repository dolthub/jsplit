package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestByteStack(t *testing.T) {
	bs := NewByteStack()

	require.Equal(t, byte(0), bs.Peek())
	require.Equal(t, byte(0), bs.Pop())

	bs.Push(byte('0'))
	bs.Push(byte('1'))
	bs.Push(byte('2'))

	require.Equal(t, byte('2'), bs.Peek())
	require.Equal(t, byte('2'), bs.Pop())
	require.Equal(t, byte('1'), bs.Peek())
	require.Equal(t, byte('1'), bs.Pop())
	require.Equal(t, byte('0'), bs.Peek())
	require.Equal(t, byte('0'), bs.Pop())

	require.Equal(t, byte(0), bs.Peek())
	require.Equal(t, byte(0), bs.Pop())
}
