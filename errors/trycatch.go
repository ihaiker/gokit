package errors

func Catch(fns ...func(re error)) {
	if err := Convert(recover()); err != nil && len(fns) > 0 {
		for _, fn := range fns {
			fn(err)
		}
	}
}

//Try handler(err)
func Try(fun func(), handler ...func(error)) {
	defer Catch(handler...)
	fun()
}

//Try handler(err) and finally
func TryFinally(fun func(), handler func(error), finallyFn func()) {
	defer finallyFn()
	Try(fun, handler)
}

//安全执行如果出错将被拦截
func Safe(fun func() error) error {
	var err error
	Try(func() {
		err = fun()
	}, func(e error) {
		err = e
	})
	return err
}

//此方法解决那些运行时异常但是不会报错的
// 例如： 从已经关闭的channel中读取数据等，golang的异常机制很是恶心啊。
func SafeGet(fun func() interface{}) (ret interface{}, err error) {
	Try(func() {
		ret = fun()
	}, func(e error) {
		err = e
	})
	return
}

func SafeExec(fun func()) (err error) {
	Try(fun, func(e error) {
		err = e
	})
	return err
}