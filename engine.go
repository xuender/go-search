package search

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/xuender/go-utils"
)

// Engine 引擎
type Engine struct {
	db DB // 数据库
}

// IndexNum 索引文件数
func (e *Engine) IndexNum() int {
	c := 0
	e.db.IteratorKey([]byte{_docIDPrefix, '-'}, func(key []byte) {
		c++
	})
	return c
}

// IndexKeys 获取索引文件Key
func (e *Engine) IndexKeys(f func(key []byte)) int {
	c := 0
	e.db.IteratorKey([]byte{_dbKey2docIDPrefix, '-'}, f)
	return c
}

// Close 关闭引擎
func (e *Engine) Close() error {
	return e.db.Close()
}

// Has 文档是否索引过
func (e *Engine) Has(key []byte) (bool, error) {
	return e.db.Has(toDBKey(key))
}

// Get 获取文档
func (e *Engine) Get(key []byte) (doc *Document, err error) {
	var has bool
	has, err = e.Has(key)
	if err != nil {
		return
	}
	if !has {
		err = errors.New("不存在")
		return
	}
	var id utils.ID
	err = e.get(toDBKey(key), &id)
	if err == nil {
		doc = &Document{}
		err = e.get(id[:], doc)
	}
	return
}

// Put 更新文档
func (e *Engine) Put(doc *Document) {
	dbKey := doc.dbKey()
	var docID utils.ID
	if has, _ := e.db.Has(dbKey); has {
		e.get(dbKey, &docID)
	} else {
		// 创建文档
		docID = utils.NewID(_docIDPrefix)
		e.put(dbKey, docID)
		e.put(utils.PrefixBytes(docID[:], _docID2DocKeyPrefix, '-'), doc.Key)
	}
	e.put(docID[:], doc)
	// TODO: 使用已存在的关键字对文档进行索引
}

// searchID 搜索文档ID
func (e *Engine) searchID(str string) []utils.ID {
	ret := utils.IDS{}          // 返回值
	index := map[string]*Word{} // 需要创建的索引
	hit := false                // 命中过

	for _, w := range Split(str) {
		ids := utils.IDS{}
		word := Word{}
		if e.get(utils.PrefixBytes([]byte(w), _wordIDPrefix, '-'), &word) == nil {
			// fmt.Println("命中", w, word)
			for _, w := range word {
				ids.Add(w.DocID)
			}
			if hit {
				ret.Intersect(ids.Distinct()...)
			} else {
				ret.Add(ids.Distinct()...)
				hit = true
			}
		} else {
			index[w] = &word
		}
	}

	if len(index) == 0 || (len(ret) == 0 && hit) {
		return ret
	}

	ids := utils.IDS{} // 没有索引查找的DocID
	e.db.Iterator([]byte{_docIDPrefix, '-'}, func(key, value []byte) {
		doc := Document{}
		if utils.Decode(value, &doc) == nil {
			id := utils.ID{}
			id.ParseBytes(key)
			for w, word := range index {
				if !word.Has(id) {
					if has, pos := doc.Inverted(w); has {
						p := Posting{
							DocID: id,
							Pos:   pos,
						}
						word.Add(p)
						// fmt.Println(word)
						e.put(utils.PrefixBytes([]byte(w), _wordIDPrefix, '-'), word)
					} else {
						return
					}
				}
			}
			ids.Add(id)
		}
	})
	if hit {
		ret.Intersect(ids.Distinct()...)
		return ret.Distinct()
	}
	return ids
}

// Search 搜索文档
func (e *Engine) Search(str string) []*Document {
	ret := []*Document{}
	for _, id := range e.searchID(str) {
		d := &Document{}
		if data, err := e.db.Get(id[:]); err == nil {
			utils.Decode(data, d)
			ret = append(ret, d)
		}
	}
	return ret
}

// SearchKey 搜索文档主键
func (e *Engine) SearchKey(str string) [][]byte {
	ret := [][]byte{}
	for _, id := range e.searchID(str) {
		docKey := []byte{}
		if err := e.get(utils.PrefixBytes(id[:], _docID2DocKeyPrefix, '-'), &docKey); err == nil {
			ret = append(ret, docKey)
		}
	}
	return ret
}

// Delete 删除文档
func (e *Engine) Delete(key []byte) error {
	dbKey := toDBKey(key)
	if has, _ := e.db.Has(dbKey); has {
		docID := utils.ID{}
		e.get(dbKey, &docID)
		if err := e.Delete(dbKey); err != nil {
			return err
		}
		if err := e.Delete(docID[:]); err != nil {
			return nil
		}
		return e.Delete(utils.PrefixBytes(docID[:], _docID2DocKeyPrefix, '-'))
	}
	return errors.New("Not found")
}

// NewLevelEngine leveldb 搜索引擎
func NewLevelEngine(path string) (*Engine, error) {
	db := &Leveldb{}
	var err error
	if db.db, err = leveldb.OpenFile(path, nil); err != nil {
		return nil, err
	}
	return &Engine{
		db: db,
	}, nil
}
