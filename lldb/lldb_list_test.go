package lldb

import (
	"testing"
	"github.com/ihaiker/gokit/commons"
	"sync"
	"fmt"
)

//-- set queue
func TestIndex(t *testing.T) {
	a := commonKit.UInt64(QUEUE_INIT_SEQ)
	b := commonKit.UInt64(QUEUE_MAX_SEQ)
	c := make([]byte, 16)
	copy(c[0:8], a)
	copy(c[8:], b)

	na := commonKit.ToUInt64(a)
	nb := commonKit.ToUInt64(b)
	t.Log(na, QUEUE_INIT_SEQ, na == QUEUE_INIT_SEQ)
	t.Log(nb, QUEUE_MAX_SEQ, nb == QUEUE_MAX_SEQ)

	na = commonKit.ToUInt64(c[0:8])
	nb = commonKit.ToUInt64(c[8:])
	t.Log(na, QUEUE_INIT_SEQ, na == QUEUE_INIT_SEQ)
	t.Log(nb, QUEUE_MAX_SEQ, nb == QUEUE_MAX_SEQ)
}

func TestQueue(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	store.FlushDB()
	key := "myfood"
	store.QPush(key, []byte("a5"))
	store.QRPush(key, []byte("a6"))
	store.QPush(key, []byte("a4"))
	store.QRPush(key, []byte("a7"))

	t.Log(store.toTest())

	v, err := store.QIndex(key, 0)
	t.Log("lindex 0 ", string(v), err)

	v, err = store.QIndex(key, 2)
	t.Log("lindex 2 ", string(v), err)

	v, err = store.QIndex(key, 3)
	t.Log("lindex 3 ", string(v), err)

	v, err = store.QIndex(key, 4)
	t.Log("lindex 4 ", string(v), err)

	v, err = store.QIndex(key, -1)
	t.Log("lindex -1", string(v), err)

	v, err = store.QIndex(key, -3)
	t.Log("lindex -3", string(v), err)

	v, err = store.QIndex(key, -4)
	t.Log("lindex -4", string(v), err)

	v, err = store.QIndex(key, 7)
	t.Log("lindex 7 ", string(v), err)

	v, err = store.QIndex(key, -7)
	t.Log("lindex -7 ", string(v), err)

	v, err = store.QPop(key)
	t.Log("lpop ", string(v), err)

	v, err = store.QPop(key)
	t.Log("lpop ", string(v), err)

	v, err = store.QRPop(key)
	t.Log("rpop ", string(v), err)

	v, err = store.QRPop(key)
	t.Log("rpop ", string(v), err)

	v, err = store.QRPop(key)
	t.Log("rpop ", string(v), err)

	v, err = store.QPop(key)
	t.Log("lpop ", string(v), err)
}

func TestQueueScan(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	store.FlushDB()

	store.QPush("b1", []byte(""))
	store.QPush("b2", []byte(""))
	store.QPush("b3", []byte(""))
	store.QPush("b4", []byte(""))

	store.QPush("t1", []byte(""))
	store.QPush("t2", []byte(""))
	store.QPush("t3", []byte(""))
	store.QPush("t4", []byte(""))
	t.Log(store.toTest())
	{
		t.Log("qlist b1 b4 10")
		it, err := store.QList("", "b4", 10)
		defer it.Release()
		if err != nil {
			t.Error(err)
		}
		for ; it.Next(); {
			t.Log(it.Get())
		}
	}
	{
		t.Log("qrlist b1 b4 10")
		it, err := store.QRList("b1", "b4", 10)
		defer it.Release()
		if err != nil {
			t.Error(err)
		}
		for ; it.Next(); {
			t.Log(it.Get())
		}
	}
}

func TestQueueTrim(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	store.FlushDB()

	store.QPush("b1", []byte("11"))
	{
		size, err := store.QSize("b1")
		t.Log("qsize : =", size, err)
	}
	store.QPush("b1", []byte("12"))
	{
		size, err := store.QSize("b1")
		t.Log("qsize : =", size, err)
	}
	store.QPush("b1", []byte("13"))
	{
		size, err := store.QSize("b1")
		t.Log("qsize : =", size, err)
	}
	store.QPush("b1", []byte("14"))
	{
		size, err := store.QSize("b1")
		t.Log("qsize : =", size, err)
	}
	store.QPush("b1", []byte("15"))
	store.QPush("b1", []byte("16"))
	store.QPush("b1", []byte("17"))
	store.QPush("b1", []byte("18"))

	{
		size, err := store.QSize("b1")
		t.Log("qsize : =", size, err)
	}
	t.Log(store.QRTrim("b1", 4))
	{
		size, err := store.QSize("b1")
		t.Log("qsize : =", size, err)
	}
	{
		size, err := store.QSize("b1")
		t.Log("qsize : =", size, err)
	}
	t.Log(store.QTrim("b1", 4))
	t.Log(store.QSize("b1"))
	t.Log(store.QRTrim("b1", 4))
	t.Log(store.QRTrim("b1", 4))

	t.Log(store.toTest())
}

func TestQueueRange(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	store.FlushDB()

	store.QRPush("b1", []byte("1"))
	store.QRPush("b1", []byte("2"))
	store.QRPush("b1", []byte("3"))
	store.QRPush("b1", []byte("4"))
	store.QRPush("b1", []byte("5"))
	store.QRPush("b1", []byte("6"))
	store.QRPush("b1", []byte("7"))
	store.QRPush("b1", []byte("8"))
	t.Log(store.toTest())

	t.Log("qslice 2 2")
	if it, err := store.QRange("b1", 2, 2); err != nil {
		t.Error(err)
	} else {
		for ; it.Next(); {
			t.Log(string(it.Value()))
		}
	}
	t.Log("qslice 1 -1")
	if it, err := store.QSlice("b1", -5, -1); err != nil {
		t.Error(err)
	} else {
		for ; it.Next(); {
			t.Log(string(it.Value()))
		}
	}

}

func TestCurrentQueue(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	store.FlushDB()

	g := &sync.WaitGroup{}
	for i:=0 ; i<100 ; i++ {
		g.Add(1)
		go func(n int) {
			for j:=0;j<100;j++ {
				store.QPush("a",[]byte(fmt.Sprintf("%02d-%02d",n,j)))
			}
			g.Done()
		}(i)
	}
	g.Wait()
	//t.Log(store.toTest())
	t.Log(store.QSize("a"))
	t.Log("OVER")
}