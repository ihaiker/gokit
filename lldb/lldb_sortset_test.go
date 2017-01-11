package lldb

import (
	"testing"
	"strings"
	"strconv"
	"math"
)

func TestEncodeSortedSet(t *testing.T) {
	{
		key := "test"
		value := []byte("OVER")
		bs := EncodeSortedSet(key, value)
		k, v := DecodeSortedSet(bs)
		t.Log(k, string(v))
	}
	//------------------------------------
	{
		key := "test"
		value := []byte(strings.Repeat("b", 300) + "OVER")
		bs := EncodeSortedSet(key, value)
		k, v := DecodeSortedSet(bs)
		t.Log(k, string(v))
	}

	{
		key := "test"
		value := []byte("OVER")
		score := uint64(10000)
		bs := EncodeSortedSetScore(key, value, score)
		k, v, c := DecodeSortedSetScore(bs)
		t.Log(k, string(v), c)
	}
}

func TestSortedSet(t *testing.T) {
	store := newDB(t)
	defer store.Close()

	t.Log(store.ZAdd("test", []byte("t1"), 1))
	t.Log("\n", store.toTest())

	t.Log(store.ZAdd("test", []byte("t1"), 2))
	t.Log("\n", store.toTest())

	t.Log(store.ZAdd("test", []byte("t1"), 3))
	t.Log("\n", store.toTest())

	t.Log(store.ZAdd("test", []byte("t2"), 5))

	t.Log(store.ZSize("test"))
}

func TestSortedSetIncr(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	k := "test"
	v := []byte("t1")

	t.Log(store.ZAdd(k, v, 1))
	t.Log(store.ZScore(k, v))

	t.Log(store.ZIncrBy(k, v, 2))
	t.Log(store.ZScore(k, v))

	t.Log(store.ZDel(k, v))

	t.Log(store.ZDel(k, v))

	t.Log(store.ZIncrBy(k, v, 10))

	t.Log(store.ZScore(k, v))
}

func TestSortedSetCount(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	k := "test"
	for i := 0; i < 1000; i++ {
		store.ZAdd(k, []byte(strconv.Itoa(i)), uint64(i))
	}
	t.Log(store.ZSize(k))
	t.Log(store.ZScore(k, []byte("0")))
	t.Log(store.ZCount(k, 1, 2))

	it, _ := store.ZList("", "", 10)
	defer it.Release()

	for ; it.Next(); {
		t.Log(it.Get())
	}

}

func TestSortedSetZlist(t *testing.T) {
	store := newDB(t)
	defer store.Close()
	for i := 0; i < 1000; i++ {
		store.ZAdd(strconv.Itoa(i), []byte(""), uint64(i))
	}
	it, _ := store.ZList("1", "107", 20)
	defer it.Release()
	for ; it.Next(); {
		t.Log(it.Get())
	}
}

func TestSortedSetZDel(t *testing.T) {
	store := newDB(t)
	defer store.Close()

	for i := 0; i < 1000; i++ {
		store.ZAdd("a", []byte(strconv.Itoa(i)), uint64(i))
	}
	t.Log(store.ZSize("a"))
	for i := 0; i < 1000; i++ {
		store.ZDel("a", []byte(strconv.Itoa(i)))
	}
	t.Log(store.ZSize("a"))
}

func TestSortedSetZDelByScore(t *testing.T) {
	store := newDB(t)
	defer store.Close()

	for i := 0; i < 1000; i++ {
		store.ZAdd("a", []byte(strconv.Itoa(i)), uint64(i))
	}
	t.Log(store.ZSize("a"))

	store.ZDelByScore("a", 0, 99)
	t.Log(store.ZSize("a"))

	store.ZDelByRank("a", 0, 9)
	t.Log(store.ZSize("a"))
}

func TestSortedZScan(t *testing.T) {
	store := newDB(t)
	defer store.Close()

	for i := 0; i < 1020; i++ {
		store.ZAdd("a", []byte(strconv.Itoa(i)), uint64(i))
	}

	t.Log(store.ZCount("a", 0, math.MaxUint64))

	t.Log(store.ZRank("a", []byte("20")))

	it, err := store.ZRange("a", 20, 1000, 1)
	if err != nil {
		t.Error(err)
	}
	defer it.Release()

	for ; it.Next(); {
		t.Log(it.Get(), it.Score(), it.Value())
	}

	t.Log(store.ZClear("a"))
	t.Log(store.ZSize("a"))

	t.Log("==== to test ====")
	t.Log(store.toTest())
}