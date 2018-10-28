package search

import (
	"fmt"

	"github.com/xuender/go-utils"
)

// Index 索引对象
type Index struct {
	Docs map[utils.ID]int // 文档
}

// Intersection 交集
func (i Index) Intersection(index Index) {
	fmt.Println(index)
	for k := range i.Docs {
		if _, ok := index.Docs[k]; !ok {
			delete(i.Docs, k)
		}
	}
}

// Add 索引增加文件
func (i Index) Add(docID utils.ID, num int) {
	i.Docs[docID] = num
}
