package mathKit

import (
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	t.Log(Subtract(1.1, 2.0))

	t.Log(Subtract(2.0, 1.1))

	t.Log(Divide(1, 3))

	t.Log(Add(1.2,0.54))
}

func TestScale(t *testing.T) {
	t.Log("5.10 down = ", ScaleDown(5.10, 1))
	t.Log("5.14 down = ", ScaleDown(5.14, 1))
	t.Log("5.44 down = ", ScaleDown(5.44, 1))
	t.Log("5.45 down = ", ScaleDown(5.45, 1))
	t.Log("5.46 down = ", ScaleDown(5.46, 1))
	t.Log("5.47 down = ", ScaleDown(5.47, 1))

	t.Log(strings.Repeat("=",20))

	t.Log("5.10 round = ", ScaleRound(5.10, 1))
	t.Log("5.14 round = ", ScaleRound(5.14, 1))
	t.Log("5.44 round = ", ScaleRound(5.44, 1))
	t.Log("5.45 round = ", ScaleRound(5.45, 1))
	t.Log("5.46 round = ", ScaleRound(5.46, 1))
	t.Log("5.47 round = ", ScaleRound(5.47, 1))

	t.Log(strings.Repeat("=",20))

	t.Log("5.10 up 1 = ", ScaleUp(5.100, 1))
	t.Log("5.14 up 1 = ", ScaleUp(5.141, 1))
	t.Log("5.44 up 1 = ", ScaleUp(5.44, 1))
	t.Log("5.45 up 1 = ", ScaleUp(5.45, 1))
	t.Log("5.46 up 1 = ", ScaleUp(5.46, 1))
	t.Log("5.47 up 1 = ", ScaleUp(5.47, 1))
}
