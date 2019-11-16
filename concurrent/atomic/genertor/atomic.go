package main

import (
	fileKit "github.com/ihaiker/gokit/files"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const templateAtomic = `
package atomic

import "sync/atomic"

type Atomic{{.Name}} struct {
	value {{.Type}}
}

func NewAtomic{{.Name}}(initValue {{.Type}}) *Atomic{{.Name}} {
	return &Atomic{{.Name}}{value: initValue}
}

func (self *Atomic{{.Name}}) Get() ({{.Type}}) {
	return atomic.Load{{.Name}}(&self.value)
}

func (self *Atomic{{.Name}}) IncrementAndGet(i uint) ({{.Type}}) {
	return atomic.Add{{.Name}}(&self.value, {{.Type}}(i))
}

func (self *Atomic{{.Name}}) GetAndIncrement(i uint) ({{.Type}}) {
	var ret {{.Type}}
	for {
		ret = atomic.Load{{.Name}}(&self.value)
		newValue := ret + {{.Type}}(i)
		if atomic.CompareAndSwap{{.Name}}(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *Atomic{{.Name}}) DecrementAndGet(i uint) ({{.Type}}) {
	var ret {{.Type}}
	for {
		ret = atomic.Load{{.Name}}(&self.value)
		newValue := ret - {{.Type}}(i)
		if atomic.CompareAndSwap{{.Name}}(&self.value, ret, newValue) {
			return newValue
		}
	}
}

func (self *Atomic{{.Name}}) GetAndDecrement(i uint) ({{.Type}}) {
	var ret {{.Type}}
	for ; ; {
		ret = atomic.Load{{.Name}}(&self.value)
		newValue := ret - {{.Type}}(i)
		if atomic.CompareAndSwap{{.Name}}(&self.value, ret, newValue) {
			return ret
		}
	}
}

func (self *Atomic{{.Name}}) Set(i {{.Type}}) {
	atomic.Store{{.Name}}(&self.value, i)
}

func (self *Atomic{{.Name}}) CompareAndSet(expect {{.Type}}, update {{.Type}}) (bool) {
	return atomic.CompareAndSwap{{.Name}}(&self.value, expect, update)
}


`

func g(dir, typeName string) {
	sweaters := struct {
		Type string
		Name string
	}{
		Type: typeName, Name: strings.Title(typeName),
	}
	tmpl, err := template.New("Atomic").Parse(templateAtomic)
	if err != nil {
		panic(err)
	}
	file, _ := fileKit.New(dir + "/" + typeName + ".go").GetWriter(false)
	defer file.Close()
	err = tmpl.Execute(file, sweaters)
	if err != nil {
		panic(err)
	}
}

func main() {
	dir, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	for _, t := range os.Args[2:] {
		g(dir, t)
	}
}
