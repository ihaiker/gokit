package lldb

import "testing"

func TestReadConfig(t *testing.T) {
	cfg, err := SetConfig("")
	if err != nil {
		t.Error(err)
        t.FailNow()
	}
    t.Log(cfg)
	t.Log(cfg.GetDataPath())
    t.Log(cfg.Locks)
    t.Log(cfg.Options.Compression)
	
}
