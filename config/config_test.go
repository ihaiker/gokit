package config

import "testing"

func TestGetStandardConfigurationLocation(t *testing.T) {
	path := GetStandardConfigurationLocation("test", "logs", "json")
	for _, p := range path {
		t.Log(p)
	}
}

type UserInfo struct {
	Name    string `json:"name" toml:"name"`
	Address struct {
		Code string `json:"code" toml:"code"`
		Info string `json:"info" toml:"info"`
	} `json:"address" toml:"address" toml:"address"`
}

func TestRegister(t *testing.T) {
	u := new(UserInfo)
	cfg := NewConfigRegister("testq", "testq")
	err := cfg.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("name: ", u.Name)
	t.Log("address.info:", u.Address.Info);
	t.Log("address.code:", u.Address.Code)
}
