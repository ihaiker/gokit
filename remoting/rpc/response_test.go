package rpc

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestResponse(t *testing.T) {
	response := new(Response)
	response.id = 3
	response.Error = nil
	response.Body = []byte("response body")
	bs, _ := response.Encode()

	resp := new(Response)
	err := resp.Decode(bs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.id, resp.id,"id maths")
	assert.Equal(t, response.Error, resp.Error,"error maths")
	assert.Equal(t, response.Body, resp.Body,"body maths")
}


func TestResponse2(t *testing.T) {
	response := new(Response)
	response.id = 3
	response.Error = io.EOF
	bs, _ := response.Encode()

	resp := new(Response)
	err := resp.Decode(bs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.id, resp.id)
	assert.Equal(t, response.Error, resp.Error)
	assert.Equal(t, response.Body, resp.Body)
}
