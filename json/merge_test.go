package json

import (
	"testing"
	"fmt"
)

var json_src = `
{
	"name":"zhou",
	"age":2,
	"address":["a","b"],
	"email":{
		"pass":100
	},
	"like":[
		{"name":"footboll","level":11},
		{"name":"test","level":1}
	]
}
`
var json_dst = `
{
	"name":"newName",
	"sex":1,
	"age":26,
	"email":{
		"user":{
			"first":"zhou","second":"haichao"
		},
		"pass":"123"
	},
	"address":["c","d"],
	"like":[
		{"name":"mv","level":99},
		{"name":"sv","level":99}
	]
}
`

func TestMerge(t *testing.T) {
	a := JsonMerger{}

	a.SetSrcString(json_src)

	a.SetDstString(json_dst)

	if err := a.Merge(); err != nil {
		t.Error(err)
	}

	fmt.Println(a.ToString())
}

func TestMergeTwo(t *testing.T) {
	a := JsonMerger{}
	a.SetSrcString(`{"one":1,"s":"a"}`)
	a.SetDstString(`{"one":2,"b":"e"}`)

	if err := a.Merge(); err != nil {
		t.Error(err)
	}
	if err := a.NewFromOut(); err != nil {
		t.Error(err)
	}
	a.SetDstString(`{"one":3}`)
	if err := a.Merge(); err != nil {
		t.Error(err)
	}
	fmt.Println(a.ToString())
}