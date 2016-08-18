package config

import (
	"testing"
	"fmt"
)

type GiveClass struct {
	Name string    `json:"name"`
	Age string    `json:"age"`
	Test string `json:"test"`
}

func TestLoad(t *testing.T) {
	cfg, err := Load(
		`{"name":"a"}`,
		`{"age":"b"}`,
		`{"test":"aa"}`,
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	out := cfg.json_data
	fmt.Println(string(out))

	testa := GiveClass{}

	if err := cfg.Unmarshal(&testa); err != nil {
		t.Error(err)
	}
	fmt.Println(testa)
}