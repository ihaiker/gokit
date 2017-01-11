package commonKit

import "io"
func Switch(check, other interface{}) interface{} {
	if check == nil {
		return other
	} else {
		switch check.(type) {
		case string:
			if ( check.(string) == "" ) {
				return other
			}
			return check
		default:
			return check
		}
	}
}
//和三元判断符同义
func IfElse(check bool,ifTrue, ifElse interface{}) interface{} {
	if check {
		return ifTrue
	}else{
		return ifElse
	}
}

//迭代器
type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type CloseIterator interface {
	Iterator
	io.Closer
}