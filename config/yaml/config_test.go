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

	if err := cfg.Load(`name: a`,`age: "reset"`,`test: "aa" `); err != nil {
		t.Error(err)
	}
	t.Log(cfg.ToString())

	if err := cfg.Load(`name: "zhou"
	"age": "20"
	"test": "ok"`); err != nil {
		t.Error(err)
	}
	t.Log(cfg.ToString())
}