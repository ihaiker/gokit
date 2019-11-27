package lv

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoder(t *testing.T) {
	lv := NewLVCoder(1024)
	message := []byte("test")
	bs, err := lv.Encode(nil, message)
	assert.Nil(t, err)

	c := bytes.NewBuffer(bs)
	out, err := lv.Decode(nil, c)
	assert.Nil(t, err)
	assert.Equal(t, "test", string(out.([]byte)))
}
