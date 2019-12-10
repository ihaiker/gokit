package rpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest(t *testing.T) {
	request := new(Request)
	request.id = 1000
	request.URL = "/test"
	request.Header("test", "test")
	request.Header("key", "b")
	request.Body = []byte(",testbody")
	bs, _ := request.Encode()

	req := new(Request)
	err := req.Decode(bs)
	assert.Nil(t, err)

	assert.Equal(t, req.id, request.id)
	assert.Equal(t, req.URL, request.URL)
	assert.Equal(t, req.Body, request.Body)
	assert.Equal(t, req.Headers, request.Headers)
}

func TestRequest2(t *testing.T) {
	request := new(Request)
	request.id = 1000
	request.URL = "/test"
	request.Body = []byte(",testbody")
	request.Headers = map[string]string{}
	bs, _ := request.Encode()

	req := new(Request)
	err := req.Decode(bs)
	assert.Nil(t, err)

	assert.Equal(t, req.id, request.id)
	assert.Equal(t, req.URL, request.URL)
	assert.Equal(t, req.Body, request.Body)
	assert.Equal(t, req.Headers, request.Headers)
}


func TestRequestNoBody(t *testing.T) {
	request := new(Request)
	request.id = 1000
	request.URL = "/test"
	bs, _ := request.Encode()

	req := new(Request)
	err := req.Decode(bs)
	assert.Nil(t, err)

	assert.Equal(t, req.id, request.id)
	assert.Equal(t, req.URL, request.URL)
	assert.Equal(t, req.Body, request.Body)
	assert.Equal(t, req.Headers, request.Headers)
}
