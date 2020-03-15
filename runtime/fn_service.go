package runtimeKit


type funcService struct {
	StartFn func() error
	StopFn  func() error
}

func (f *funcService) Start() error {
	if f.StartFn != nil {
		return f.StartFn()
	}
	return nil
}

func (f *funcService) Stop() error {
	if f.StopFn != nil {
		return f.StopFn()
	}
	return nil
}

func StartService(ob interface{}) error {
	if ob != nil {
		if sv, match := ob.(Service); match {
			return sv.Start()
		}
	}
	return nil
}

func StopService(ob interface{}) error {
	if ob != nil {
		if sv, match := ob.(Service); match {
			return sv.Stop()
		}
	}
	return nil
}

