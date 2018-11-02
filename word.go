package search

import "github.com/xuender/go-utils"

// Word 词索引
type Word []Posting

// Has 包含DocID
func (w Word) Has(id utils.ID) bool {
	for _, p := range w {
		if p.DocID == id {
			return true
		}
	}
	return false
}

// Add 增加倒排索引项
func (w *Word) Add(p Posting) {
	*w = append(*w, p)
}
