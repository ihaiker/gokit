package lldb

import (
	"github.com/syndtr/goleveldb/leveldb/util"
)

//set key value
//设置指定 key 的值内容.
func (self *LLDBEngine) Set(key string, value []byte) error {
	return self.data.Put(EncodeKV(key), value, self.writeOptions)
}
func (self *LLDBEngine) Get(key string) ([]byte, error) {
	return self.data.Get(EncodeKV(key), self.readOptions)
}

//删除key
func (self *LLDBEngine) Del(key string) (error) {
	return self.data.Delete(EncodeKV(key), self.writeOptions)
}

func (self *LLDBEngine) Has(key string) (bool, error) {
	return self.data.Has(EncodeKV(key), self.readOptions)
}

//mast close KeyIterator.Release
func (self *LLDBEngine) Scan(startKey, endKey string, limit int) (Iterator, error) {
	startRange := EncodeKV(startKey)
	var endRange []byte = nil
	if (endKey != "" ) {
		endRange = EncodeKV(endKey + "\001")
	}
	bp := &util.Range{Start:startRange, Limit:endRange}

	sn,err := self.data.GetSnapshot()
	if err != nil {
		return nil,err
	}
	it := sn.NewIterator(bp, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	it.Seek(startRange) //skip startKey
	if (!it.Valid()) {
		it.Last();
	}

	return newKVIterator(startKey, endKey, limit, FORWARD, it), nil
}

func (self *LLDBEngine) RScan(startKey, endKey string, limit int) (Iterator, error) {
	startRange := EncodeKV(startKey)
	var endRange []byte
	if endKey == "" {
		endRange = EncodeKV("\255")
	} else {
		endRange = EncodeKV(endKey + "\001")
	}
	bp := &util.Range{Start:startRange, Limit:endRange}

	sn,err := self.data.GetSnapshot()
	if err != nil {
		return nil,err
	}
	it := sn.NewIterator(bp, self.readOptions)
	if err := it.Error(); err != nil {
		return nil, err
	}
	it.Last()
	it.Seek(endRange)
	return newKVIterator(startKey, endKey, limit, BACKWARD, it), nil
}