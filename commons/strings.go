package commonKit

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


//迭代器
type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type CloseIterator interface {
	Iterator
	Close() error
}