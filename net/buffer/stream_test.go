package buffer

import (
	bytes2 "bytes"
	"testing"
)

func TestStream(t *testing.T) {
	out := bytes2.NewBuffer([]byte{})
	{
		if err := WriteInt8(out, 10); err != nil {
			t.Fatal(err)
		}

		i, err := ReadInt8(out)
		t.Log(i, err)
	}
	{
		if err := WriteInt32(out, 32); err != nil {
			t.Fatal(err)
		}
		i, err := ReadInt32(out)
		t.Log(i, err)
	}

}
