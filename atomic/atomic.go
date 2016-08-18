package atomic

func NewInt32V(value int32) (*AtomicInt32) {
	return &AtomicInt32{value:value}
}

func NewInt32() (*AtomicInt32) {
	return NewInt32V(0)
}