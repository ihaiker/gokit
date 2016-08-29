package atomic

func NewInt32V(value int32) (*AtomicInt32) {
	return &AtomicInt32{value:value}
}

func NewInt32() (*AtomicInt32) {
	return NewInt32V(0)
}


func NewInt64V(value int64) (*AtomicInt64) {
	return &AtomicInt64{value:value}
}

func NewInt64() (*AtomicInt64) {
	return NewInt64V(0)
}