package commons

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAsyncTimeout(t *testing.T) {

	err := AsyncTimeout(time.Second, func() interface{} {
		time.Sleep(time.Second * 2)
		return 0
	})
	assert.Equal(t,err,ErrAsyncTimeout)

	err = AsyncTimeout(time.Second*2, func() interface{} {
		time.Sleep(time.Second)
		return 1
	})

	assert.Equal(t,err,1)

}
