package geecache

// create an immutable view of the cached data
type ByteView struct {
	b []byte
}

// return the length of the underlying byte
func (v ByteView) Len() int {
	return len(v.b)
}

// clone a []byte
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// return a copy of our data
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// return a string copy if our data
func (v ByteView) String() string {
	return string(v.b)
}
