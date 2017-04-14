package atomic

func NewInt32V(value int32) (*AtomicInt32) {
	return &AtomicInt32{value:value}
}

func NewInt32() (*AtomicInt32) {
	return NewInt32V(0)
}

func NewUInt32V(value uint32) (*AtomicUInt32) {
	return &AtomicUInt32{value:value}
}

func NewUInt32() (*AtomicUInt32) {
	return NewUInt32V(0)
}

func NewInt64V(value int64) (*AtomicInt64) {
	return &AtomicInt64{value:value}
}

func NewInt64() (*AtomicInt64) {
	return NewInt64V(0)
}

func NewUInt64V(value uint64) (*AtomicUInt64) {
	return &AtomicUInt64{value:value}
}

func NewUInt64() (*AtomicUInt64) {
	return NewUInt64V(0)
}