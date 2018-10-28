package search

import utils "github.com/xuender/go-utils"

// DB 数据库接口
type DB interface {
	// Close 数据库关闭
	Close() error
	// Get 获取
	Get(key []byte) ([]byte, error)
	// Put 写入
	Put(key, value []byte) error
	// Has 是否包含
	Has(key []byte) (bool, error)
	// Delete 删除
	Delete(key []byte) error
	// Iterator 迭代数据
	Iterator(prefix []byte, f func(key, value []byte)) error
	// IteratorKey 迭代Key
	IteratorKey(prefix []byte, f func(key []byte)) error
}

// put 写入数据库
func (e *Engine) put(key []byte, i interface{}) error {
	bs, err := utils.Encode(i)
	if err != nil {
		return err
	}
	return e.db.Put(key, bs)
}

// has 数据库是否包含
func (e *Engine) has(key []byte) (bool, error) {
	return e.db.Has(key)
}

// get 读取数据
func (e Engine) get(key []byte, i interface{}) error {
	data, err := e.db.Get(key)
	if err != nil {
		return err
	}
	utils.Decode(data, i)
	return nil
}
