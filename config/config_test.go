package config

import (
	"github.com/jinzhu/configor"
	"os"
	"testing"
)

type UserInfo struct {
	Name    string `json:"name" toml:"name"`
	Address struct {
		Code string `json:"code" toml:"code"`
		Info string `json:"info" toml:"info"`
	} `json:"address" toml:"address" toml:"address"`
}

func TestRegister(t *testing.T) {

	_ = os.Setenv("CONFIGOR_DEBUG_MODE", "true")
	_ = os.Setenv("CONFIGOR_VERBOSE_MODE", "true")
	_ = os.Setenv("CONFIGOR_SILENT_MODE", "true")

	u := new(UserInfo)
	cfg := NewConfigRegister("test", "test")

	cfg = cfg.With(&configor.Config{
		ENVPrefix:            "TEST",
	})

	if err := cfg.MustExitConfig("./test.noufound.json"); err != nil {
		t.Log(err)
	}

	err := cfg.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name: ", u.Name)
	t.Log("address.info:", u.Address.Info)
	t.Log("address.code:", u.Address.Code)
}
