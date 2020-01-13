package commons

func SafeClose(c chan struct{}) {
	defer func() { _ = recover() }()
	close(c)
}

func SafeIClose(c chan interface{}) {
	defer func() { _ = recover() }()
	close(c)
}

func SafeBoolClose(c chan bool) {
	defer func() { _ = recover() }()
	close(c)
}
