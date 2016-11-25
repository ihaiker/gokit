package lldb

import (
	"testing"
	"strconv"
	"os"
)

func newDB(t *testing.T) *LLDBEngine {
	store, err := New()
	if err != nil {
		t.Error(err)
	}
	return store
}

func TestPath(t *testing.T) {
	t.Log(LSCL_DEFAULT_PATH)
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
	store, _ := New()
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


//-- set queue
func TestIndex(t *testing.T) {
	a := uint642byte(QUEUE_INIT_SEQ)
	b := uint642byte(QUEUE_MAX_SEQ)
	c := make([]byte, 16)
	copy(c[0:8], a)
	copy(c[8:], b)

	na, nae := byte2uint64(a)
	nb, nbe := byte2uint64(b)
	t.Log(nae, nbe)
	t.Log(na, QUEUE_INIT_SEQ, na == QUEUE_INIT_SEQ)
	t.Log(nb, QUEUE_MAX_SEQ, nb == QUEUE_MAX_SEQ)

	na, nae = byte2uint64(c[0:8])
	nb, nbe = byte2uint64(c[8:])
	t.Log(nae, nbe)
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

	store.toTest(os.Stdout)

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
	store.toTest(os.Stdout)
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

	store.toTest(os.Stdout)
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
	store.toTest(os.Stdout)

	t.Log("qslice 2 2")
	if it, err := store.QRange("b1", 2, 2); err !=nil {
		t.Error(err)
	}else{
		for ; it.Next(); {
			t.Log(string(it.Value()))
		}
	}
	t.Log("qslice 1 -1")
	if it, err := store.QSlice("b1", -5, -1); err !=nil {
		t.Error(err)
	}else{
		for ; it.Next(); {
			t.Log(string(it.Value()))
		}
	}

}

func TestToTest(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	store.toTest(os.Stdout)
}
