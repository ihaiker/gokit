package rpc

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/remoting"
	"time"
)

type RpcClient interface {
	Start()

	Shutdown()

	Send(request *Request, timeout time.Duration) (response *Response)

	Async(request *Request, timeout time.Duration, callback func(response *Response))

	Oneway(request *Request, timeout time.Duration)
}

type rpcClient struct {
	client    remoting.Client
	respCache map[uint32]*responseCache
	id        *atomic.AtomicUint32
}

func NewClient(address string, onMessage OnMessage) (RpcClient, error) {
	return NewClientWithConfig(address, remoting.DefaultTCPConfig(), onMessage)
}

func NewClientWithConfig(address string, config *remoting.Config, onMessage OnMessage) (RpcClient, error) {
	client := new(rpcClient)
	client.id = atomic.NewAtomicUint32(0)
	if tcpClient, err := remoting.NewClient(address, config, newHandler(onMessage, client.dealResponse), newCoder()); err != nil {
		return nil, err
	} else {
		client.client = tcpClient
	}
	client.respCache = make(map[uint32]*responseCache)
	return client, nil
}

func (s *rpcClient) Start() {
	s.client.Start()
}

func (s *rpcClient) Shutdown() {
	s.client.Close().Wait()
}

func (s *rpcClient) dealResponse(resp *Response) {
	if cache, has := s.respCache[resp.id]; has {
		delete(s.respCache, resp.id)
		commons.Try(func() {
			if cache.C != nil {
				cache.C <- resp
			} else {
				cache.callback(resp)
			}
		}, func(e error) { //防止并发问题正好删除关闭触发这里
			logger.Warn("dealResponse error:", e)
		})
	} else {
		logger.Debug("ignore response: ", resp.id)
	}
}

func (s *rpcClient) Send(request *Request, timeout time.Duration) *Response {
	request.id = s.id.IncrementAndGet(1)

	logger.Debug("client send :", request.id, " ", request.URL)
	response := new(Response)
	if err := s.client.Send(request, timeout); err != nil {
		response.Error = err
		return response
	}

	rc := &responseCache{C: make(chan *Response)}
	s.respCache[request.id] = rc

	select {
	case resp := <-rc.C:
		response = resp
	case <-time.After(timeout):
		response.Error = ErrRpcTimeout
	}

	delete(s.respCache, request.id)
	rc.Close()

	return response
}

func (s *rpcClient) Async(request *Request, timeout time.Duration, callback func(response *Response)) {
	request.id = s.id.IncrementAndGet(1)

	rc := &responseCache{callback: callback, timeout: time.Now().Add(timeout)}
	if err := s.client.Send(request, timeout); err != nil {
		callback(&Response{Error: err, id: request.id})
	} else {
		s.respCache[request.id] = rc
	}
}

func (s *rpcClient) Oneway(request *Request, timeout time.Duration) {
	s.Async(request, timeout, func(response *Response) {
		if response.Error != nil {
			logger.Error("send oneway error:", response.Error.Error())
		}
	})
}
