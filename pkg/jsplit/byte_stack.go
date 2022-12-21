package jsplit

// ByteStack is a simple stack of bytes
type ByteStack struct {
	chars []byte
}

// NewByteStack returns a new ByteStack object
func NewByteStack() *ByteStack {
	return &ByteStack{chars: make([]byte, 0, 64)}
}

// Push pushes a new byte on the stack
func (bs *ByteStack) Push(b byte) {
	bs.chars = append(bs.chars, b)
}

// Pop takes the top value of the top of the stack and returns it
func (bs *ByteStack) Pop() byte {
	l := len(bs.chars)

	if l == 0 {
		return 0
	}

	ch := bs.chars[l-1]
	bs.chars = bs.chars[:l-1]

	return ch
}

// Peek returns the value at the top of the stack
func (bs *ByteStack) Peek() byte {
	l := len(bs.chars)

	if l == 0 {
		return 0
	}

	return bs.chars[l-1]
}
