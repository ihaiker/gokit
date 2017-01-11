package lldb

import (
	"testing"
	"strconv"
)

func newDB(t *testing.T) *LLDBEngine {
	store, err := Default()
	if err != nil {
		t.Error(err)
	}
	return store
}

func TestInit(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	t.Log("OVER")
}

func TestKV(t *testing.T) {
	key := "test"
	store := newDB(t)
	if has, err := store.Has(key); has {
		t.Log("has :", key)
		t.Log("delete :", store.Del("test"))
	} else {
		t.Log("check key :", key, "error :", err)
	}
	t.Log("set error :", store.Set(key, []byte("name")))
	val, err := store.Get(key)
	t.Log("get ", key, " : ", string(val), err)
}

func TestScan1(t *testing.T) {
	store, _ := Default()
	for i := 0; i < 10; i++ {
		store.Set("b" + strconv.Itoa(i), []byte("b" + strconv.Itoa(i)))
	}

	startKey := ""
	endKey := ""

	keys, err := store.Scan(startKey, endKey, 2)
	if err != nil {
		t.Error(err)
	}
	i := 0
	for ; keys.Next(); {
		startKey = keys.Get()
		t.Log(keys.Get(), keys.Value())
		i++
	}
	keys.Release() //mast
}

func TestScanAll(t *testing.T) {
	store := newDB(t)
	for i := 0; i < 10; i++ {
		store.Set("b" + strconv.Itoa(i), []byte("b" + strconv.Itoa(i)))
	}

	startKey := "a5"
	endKey := "b3"

	for {
		keys, err := store.Scan(startKey, endKey, 2)
		if err != nil {
			t.Error(err)
		}
		i := 0
		for ; keys.Next(); {
			startKey = keys.Get()
			t.Log(keys.Get(), keys.Value())
			i++
		}
		if i < 2 {
			break
		}

		keys.Release() //mast
	}
}
func TestRScan(t *testing.T) {
	store := newDB(t)
	for i := 0; i < 10; i++ {
		store.Set("b" + strconv.Itoa(i), []byte("b" + strconv.Itoa(i)))
	}
	keys, err := store.RScan("b5", "b9", 20)
	if err != nil {
		t.Error(err)
	}
	defer keys.Release() //mast

	for ; keys.Next(); {
		t.Log(keys.Get(), keys.Value())
	}
	t.Log("OVER")
}
