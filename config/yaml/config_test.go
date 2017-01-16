package yaml

import (
	"testing"
)

type GiveClass struct {
	Name string    `json:"name"`
	Age  string    `json:"age"`
	Test string `json:"test"`
}

func TestYaml(t *testing.T) {
	cfg, err := Config()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cfg.Load(
`user:
    name: testname
`)
	if err := cfg.Load(`name: a`,`age: "reset"`,`test: "aa" `); err != nil {
		t.Error(err)
	}
	t.Log(cfg.ToString())

	if err := cfg.Load(`name: "zhou"`,`"test": "ok"`,`age: 23`); err != nil {
		t.Error(err)
	}
	t.Log(cfg.ToString())
	
	t.Log(" =================== ")
	t.Log(cfg.GetString("user.name"))
	t.Log(cfg.GetString("user.name.first"))
	
	t.Log(" ------------------------ ")
	cfg.Load(
`array:
    - 1
    - "zhou"
`)
	t.Log(cfg.GetSlice("array"))
}