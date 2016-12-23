package lldb

import (
	"testing"
	"strings"
)


func TestEncodeSet(t *testing.T) {
	{
		key := "test"
		value := []byte("OVER")
		bs := EncodeSet(key,value)
		k,v := DecodeSet(bs)
		t.Log(k,string(v))
	}
	//------------------------------------
	{
		key := "test"
		value := []byte(strings.Repeat("b",300) + "OVER")
		bs := EncodeSet(key,value)
		k,v := DecodeSet(bs)
		t.Log(k,string(v))
	}
}

func TestSSet(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	store.FlushDB()

	store.SAdd("t1", []byte("1"))
	store.SAdd("t1", []byte("2"))
	store.SAdd("t1", []byte("2"))
	store.SAdd("t1", []byte("3"))

	store.SAdd("t2", []byte("1"))
	store.SAdd("t2", []byte("2"))

	{
		s, e := store.SSize("t1")
		t.Log("size t1", s, e)
	}
	{
		s, e := store.SDel("t1", []byte("1"))
		t.Log("delete t1 1 ", s, e)
	}
	{
		s, e := store.SSize("t1")
		t.Log("size t1 ", s, e)
	}
	{
		s, e := store.SExits("t1", []byte("1"))
		t.Log("exit t1 1", s, e)
	}
	{
		s, e := store.SExits("t1", []byte("2"))
		t.Log("exit t1 2", s, e)
	}

	{
		t.Log("================== smembers t1")
		it, err := store.SMembers("t1")
		if err != nil {
			t.Error(err); return
		}
		defer it.Release()
		for ; it.Next(); {
			t.Log(it.Get(), string(it.Value()))
		}
	}

	{
		t.Log("================== slist '' '' 10 ")
		it, err := store.SList("a", "t2", 100)
		if err != nil {
			t.Error(err); return
		}
		defer it.Release()
		for ; it.Next(); {
			t.Log(it.Get())
		}
	}
	{
		t.Log("================== SRandmonMember ===")
		{
			v,err := store.SRandomMember("t1")
			t.Log("srandommember t1 ",string(v),err)
		}
		{
			v,err := store.SRandomMember("t1")
			t.Log("srandommember t1 ",string(v),err)
		}
		{
			v,err := store.SPop("t1")
			t.Log("spop t1 ",string(v),err)
		}

		{
			t.Log("================== smembers t1")
			it, err := store.SMembers("t1")
			if err != nil {
				t.Error(err); return
			}
			defer it.Release()
			for ; it.Next(); {
				t.Log(it.Get(), string(it.Value()))
			}
		}

		{
			v,err := store.SPop("t1")
			t.Log("spop t1 ",string(v),err)
		}
		{
			t.Log("================== smembers t1")
			it, err := store.SMembers("t1")
			if err != nil {
				t.Error(err); return
			}
			defer it.Release()
			for ; it.Next(); {
				t.Log(it.Get(), string(it.Value()))
			}
		}
		{
			v,err := store.SPop("t1")
			t.Log("spop t1 ",string(v),err)
		}
	}
}
