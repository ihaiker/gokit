package atomic

import (
	"testing"
	"sync"
	"fmt"
)

func TestInt32(t *testing.T) {
	a := NewInt32()
	a.Get()
}

func TestNewInt32V(t *testing.T) {
	a := NewInt32V(1)
	a.Get()
}

func TestAddAndGet(t *testing.T) {
	a := NewInt32()
	var group sync.WaitGroup
	for i:=0 ; i< 10 ; i++  {
		group.Add(1)
		go func() {
			defer group.Done()
			for j:=0; j<10 ; j++  {
				a.IncrementAndGet()
			}
		}()
	}
	group.Wait()
	fmt.Println(a.Get())
}

func TestAddAndGet1(t *testing.T) {
	a := NewInt32()
	var group sync.WaitGroup
	for i:=0 ; i< 10 ; i++  {
		group.Add(1)
		go func() {
			defer group.Done()
			for j:=0; j<10 ; j++  {
				a.AddAndGet(1)
			}
		}()
	}
	group.Wait()
	if 100 != a.Get() {
		t.Error("error ")
	}
}

func TestGetAndDecrement(t *testing.T)  {
	a := NewInt32V(100)
	var group sync.WaitGroup
	for i:=0 ; i< 10 ; i++  {
		group.Add(1)
		go func() {
			defer group.Done()
			for j:=0; j<10 ; j++  {
				a.GetAndDecrement()
			}
		}()
	}
	group.Wait()

	if 0 != a.Get() {
		t.Error("error GetAndDecrement ")
	}
}

func TestGetAndIncrement2(t *testing.T)  {
	a := NewInt32()
	var group sync.WaitGroup
	for i:=0 ; i< 10 ; i++  {
		group.Add(1)
		go func() {
			defer group.Done()
			for j:=0; j<5 ; j++  {
				a.AddAndGet(2)
			}
		}()
	}
	group.Wait()

	if v := a.Get(); v != 100 {
		t.Error("error TestGetAndDecrement2 ",v)
	}
}


func TestAtomicInt32_DecrementAndGet(t *testing.T) {
	a := NewInt32V(100)
	if v := a.Get(); v != 100 {
		fmt.Println("Init ",v)
	}
	var group sync.WaitGroup
	for i:=0 ; i< 10 ; i++  {
		group.Add(1)
		go func() {
			defer group.Done()
			for j:=0; j<5 ; j++  {
				a.GetAndAdd(-2)
			}
		}()
	}
	group.Wait()

	if v := a.Get(); v != 0 {
		t.Error("error TestAtomicInt32_DecrementAndGet ",v)
	}
}




