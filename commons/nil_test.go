package commons

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func aa() (a *struct{ Name string }) {
	return
}

func TestNi(t *testing.T) {
	a := aa()
	assert.True(t, IsNil(a))
}
