package commons

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAsyncTimeout(t *testing.T) {

	err := AsyncTimeout(time.Second, func() interface{} {
		time.Sleep(time.Second * 4)
		return 0
	})
	assert.Equal(t, err, ErrAsyncTimeout)

	time.Sleep(time.Second * 5)
}
