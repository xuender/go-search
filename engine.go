package search

import (
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/xuender/go-utils"
)

// Engine 引擎
type Engine struct {
	db DB // 数据库
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
	if has, _ := e.db.Has(dbKey); has {
		var oldID utils.ID
		e.get(dbKey, &oldID)
		old := Document{}
		e.get(oldID[:], &old)
		// TODO: 起新线程删除old一切记录
		fmt.Println("删除", old)
	}
	docID := utils.NewID(_docIDPrefix)
	e.put(dbKey, docID)
	e.put(utils.PrefixBytes(docID[:], _docID2DocKeyPrefix, '-'), doc.Key)
	e.index(doc, docID)
	e.put(docID[:], doc)
}

func (e *Engine) index(doc *Document, id utils.ID) {
	doc.index()
	words := []string{}
	e.get(_words, &words)
	for _, w := range words {
		doc.Add(w)
	}

	if doc.TitleIndex != nil {
		for k, v := range *(doc.TitleIndex) {
			e.putIndex(_titleIDPrefix, k, v, id)
		}
	}
	if doc.SummaryIndex != nil {
		for k, v := range *(doc.SummaryIndex) {
			e.putIndex(_summaryIDPrefix, k, v, id)
		}
	}
	if doc.ContentIndex != nil {
		for k, v := range *(doc.ContentIndex) {
			e.putIndex(_contentIDPrefix, k, v, id)
		}
	}
}

func (e *Engine) putIndex(b byte, k string, v int, id utils.ID) {
	key := utils.PrefixBytes([]byte(k), b, '-')
	i := e.getIndex(key)
	i.Add(id, v)
	e.put(key, i)
}

// searchID 搜索文档ID
func (e *Engine) searchID(str string) []utils.ID {
	// 生成搜索关键词
	keyword := NewKeyword(str)
	e.loadKeyword(keyword)
	words := keyword.Get()

	ret := []utils.ID{}
	// 去重
	m := map[utils.ID]bool{}
	for _, b := range _IndexPrefixs {
		var i Index
		for _, k := range words {
			key := utils.PrefixBytes([]byte(k), b, '-')
			if i.Docs == nil {
				i = *e.getIndex(key)
			} else {
				i.Intersection(*e.getIndex(key))
			}
		}
		if i.Docs != nil {
			for id := range i.Docs {
				if has, _ := m[id]; !has {
					ret = append(ret, id)
				}
			}
		}
	}
	return ret
}

func (e *Engine) loadKeyword(k *Keyword) {
	// 是否有索引
	for w := range *(k.words) {
		for _, b := range _IndexPrefixs {
			key := utils.PrefixBytes([]byte(w), b, '-')
			if has, _ := e.db.Has(key); has {
				index := e.getIndex(key)
				k.words.AddNum(w, len(index.Docs))
				break
			}
		}
	}
	// 未创建过索引的关键词
	zero := []string{}
	for k, n := range *(k.words) {
		if n == 0 {
			zero = append(zero, k)
		}
	}

	// 遍历文档
	e.db.Iterator([]byte{_docIDPrefix, '-'}, func(key, value []byte) {
		doc := Document{}
		if utils.Decode(value, &doc) != nil {
			return
		}
		for _, b := range _IndexPrefixs {
			for _, w := range zero {
				if c := doc.AddByByte(w, b); c > 0 {
					e.put(key, doc)
					k.words.Add(w)
					indexKey := utils.PrefixBytes([]byte(w), b, '-')
					i := e.getIndex(indexKey)
					id := utils.ID{}
					id.ParseBytes(key)
					i.Add(id, c)
					e.put(indexKey, *i)
				}
			}
		}
	})
}

func (e *Engine) getIndex(key []byte) *Index {
	if has, _ := e.has(key); has {
		i := &Index{}
		e.get(key, i)
		return i
	}
	return &Index{Docs: map[utils.ID]int{}}
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
func (e *Engine) Delete(docID []byte) error {
	// TODO
	return nil
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
