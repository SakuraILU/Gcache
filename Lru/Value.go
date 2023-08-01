package lru

type Value struct {
	bytes []byte
}

func NewValue(bytes []byte) (v Value) {
	v.bytes = bytes
	return
}

func (v *Value) ToString() (s string) {
	s = string(v.bytes)
	return
}

func (v *Value) ByteSlice() (newbytes []byte) {
	newbytes = make([]byte, len(v.bytes))
	copy(newbytes, v.bytes)
	return
}

func (v *Value) Len() int {
	return len(v.bytes)
}
