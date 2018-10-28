package search

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Leveldb 数据库支持
type Leveldb struct {
	db *leveldb.DB
}

// Close 关闭数据库
func (l *Leveldb) Close() error {
	return l.db.Close()
}

// Get 获取
func (l *Leveldb) Get(key []byte) ([]byte, error) {
	return l.db.Get(key, nil)
}

// Put 修改
func (l *Leveldb) Put(key, value []byte) error {
	return l.db.Put(key, value, nil)
}

// Has 包含
func (l *Leveldb) Has(key []byte) (bool, error) {
	return l.db.Has(key, nil)
}

// Delete 删除
func (l *Leveldb) Delete(key []byte) error {
	return l.db.Delete(key, nil)
}

// Iterator 迭代获取数据
func (l *Leveldb) Iterator(prefix []byte, f func(key, value []byte)) error {
	iter := l.db.NewIterator(util.BytesPrefix(prefix), nil)
	for iter.Next() {
		f(iter.Key(), iter.Value())
	}
	iter.Release()
	return iter.Error()
}

// IteratorKey 迭代Key
func (l *Leveldb) IteratorKey(prefix []byte, f func(key []byte)) error {
	iter := l.db.NewIterator(util.BytesPrefix(prefix), nil)
	for iter.Next() {
		f(iter.Key())
	}
	iter.Release()
	return iter.Error()
}
