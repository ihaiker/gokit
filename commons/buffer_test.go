package commonKit

import (
	"testing"
	"math"
)

func TestBuffer(t *testing.T) {
	t.Log(BYTE8,BYTE16,BYTE32,BYTE64)

	w := NewWriter()
	w.Int8(math.MaxInt8)
	w.Int16(16)
	w.Int32(32)
	w.Int64(64)

	w.Float32(3.4)
	w.String("sb")


	t.Log(w.ToBytes())

	r := w.ToReader()

	t.Log(r.Int8())
	t.Log(r.Int16())
	t.Log(r.Int32())
	t.Log(r.Int64())

	t.Log(r.Float32())

	t.Log(r.String())
}
