package search

import (
	"sort"

	"github.com/xuender/go-utils"
)

// Posting 倒排索引项
type Posting struct {
	DocID utils.ID      // 文件标识
	Pos   sort.IntSlice // 位置信息
}

// TF 出现次数
func (p *Posting) TF() int {
	return len(p.Pos)
}
