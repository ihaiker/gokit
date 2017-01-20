package commonKit

//Try handler(err)
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}
//Try handler(err) and finally
func TryFinally(fun func(), handler func(interface{}), finallyFn func()) {
	defer finallyFn()
	Try(fun,handler)
}
