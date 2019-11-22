package main

import (
	"fmt"
	"github.com/ihaiker/gokit/net/buffer"
)

func main() {
	ary := []uint16{
		2,
		14648,
	}
	bs := []byte{}

	for _, b := range ary {
		ubs := buffer.UInt16(b)
		bs = append(bs, ubs...)
	}

	typeId := buffer.ToUInt16(bs)
///	length := buffer.ToUInt16(bs)
	fmt.Println(typeId)

	fmt.Println(string(bs))
}
