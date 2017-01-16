package json

import (
	"testing"
)

type GiveClass struct {
	Name string    `json:"name"`
	Age  string    `json:"age"`
	Test string `json:"test"`
}

func TestJson(t *testing.T) {
	cfg, err := Config()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := cfg.Load(`{"name":"a","age":"age"}`,
		`{"age":"reset"}`,
		`{"test":"aa"}`, ); err != nil {
		t.Error(err)
	}
	cfg.Load(`{"user":{"name":"test-Name","address":{"code":"1024"}}}`)
	t.Log(cfg.ToString())

	if err := cfg.Load(`{"name":"zhou","age":"20","test":"ok"}`); err != nil {
		t.Error(err)
	}
	t.Log(cfg.ToString())
	t.Log(cfg.GetString("user.address.code"))
	
	t.Log(cfg.Load(`{"array":[1,23,4,"3",{"test":1}]}`))
	t.Log(cfg.GetSlice("array"))
	
	t.Log(cfg.GetSlice("name"))
}