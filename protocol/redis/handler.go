package redis

import (
	"reflect"
	"strings"
	"github.com/ihaiker/gokit/commons/logs"
)

type HandlerFn func(r *Request) (ReplyWriter, error)

func (srv *Server) RegisterFct(key string, f interface{}) error {
	v := reflect.ValueOf(f)
	handlerFn, err := srv.createHandlerFn(f, &v)
	if err != nil {
		return err
	}
	srv.Register(key, handlerFn)
	return nil
}

func (srv *Server) Register(name string, fn HandlerFn) {
	if srv.methods == nil {
		srv.methods = make(map[string]HandlerFn)
	}
	if fn != nil {
		srv.methods[strings.ToLower(name)] = fn
	}
}

func (srv *Server) Apply(r *Request) (ReplyWriter, error) {
	if srv == nil || srv.methods == nil {
		logs.Debugf("The method map is uninitialized")
		return ErrMethodNotSupported, nil
	}
	fn, exists := srv.methods[strings.ToLower(r.Name)]
	if !exists {
		return ErrMethodNotSupported, nil
	}
	return fn(r)
}

func (srv *Server) ApplyString(r *Request) (string, error) {
	reply, err := srv.Apply(r)
	if err != nil {
		return "", err
	}
	return ReplyToString(reply)
}

func (srv *Server) RegisterHandler(handler interface{}) error {
	rh := reflect.TypeOf(handler)
	for i := 0; i < rh.NumMethod(); i++ {
		method := rh.Method(i)
		if method.Name[0] > 'a' && method.Name[0] < 'z' {
			continue
		}
		handlerFn, err := srv.createHandlerFn(handler, &method.Func)
		if err != nil {
			return err
		}
		logs.Infof("registier: %s.%s ", rh.String(), method.Name)
		srv.Register(method.Name, handlerFn)
	}
	return nil
}
