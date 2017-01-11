package yaml

import (
	"testing"
	"fmt"
)

var json_src = `name: name
age: 123
sex: true
open: yes
close: no
address:
    name: 123
    code: 10000
`
var json_dst = `name: "zhou"
sex: true
age: 20
address:
    name: HeibeiHand
    code: 100230
`

func TestMerge(t *testing.T) {
	a := Merger()

	a.SetSrcString(json_src)

	a.SetDstString(json_dst)

	if err := a.Merge(); err != nil {
		t.Error(err)
	}

	fmt.Println(a.ToString())
}