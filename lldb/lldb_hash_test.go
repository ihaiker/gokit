package lldb

import (
	"testing"
	"strconv"
)

func TestHashEncode(t *testing.T) {
	bs := EncodeHash("[123]", "<name>")
	t.Log(string(bs))
	key, label := DecodeHash(bs)
	t.Log(key, label)
}

func TestHash(t *testing.T) {
	store := newDB(t)

	t.Log(store.HSet("1002", "name", []byte("zhou")))

	v, er := store.HGet("1002", "name")
	t.Log(string(v), er)

	t.Log(store.HSet("1002", "name", []byte("zhou")))

	i, er := store.HDel("1002", "name")
	t.Log(i, er)

	v, er = store.HGet("1002", "name")
	t.Log(string(v), er)

	store.Close()
}

func TestHashScan(t *testing.T) {
	store := newDB(t)
	for i := 0; i < 19; i++ {
		if _, err := store.HSet("1002", "nam2" + strconv.Itoa(i), []byte("zhou")); err != nil {
			t.Error(err)
		}
	}
	it, _ := store.HScan("1002", "nam24", "nam28", 20)
	for i := 0; it.Next(); i++ {
		t.Log(i, it.Get(), string(it.Value()))
	}
	it.Release()
}

func TestHashHList(t *testing.T) {
	store := newDB(t)
	store.HSet("a000", "age", []byte("1`5"))
	store.HSet("1001", "age", []byte("1`5"))
	store.HSet("1002", "age", []byte("1`5"))
	store.HSet("1003", "age", []byte("1`5"))

	{
		it, err := store.HList("", "", 10)
		defer it.Release()
		if err != nil {
			t.Error(err)
		}
		for ; it.Next(); {
			t.Log(it.Get())
		}
	}
	t.Log("hrlist =========== ")
	{
		it, err := store.HRList("1001", "1003", 10)
		defer it.Release()
		if err != nil {
			t.Error(err)
		}
		for ; it.Next(); {
			t.Log(it.Get())
		}
	}
}

func TestHashGetAll(t *testing.T) {
	store := newDB(t)
	it, _ := store.HGetAll("1002")
	for ; it.Next(); {
		t.Log(it.Get(), string(it.Value()))
	}
}


func TestToTest(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	t.Log(store.toTest())
}
