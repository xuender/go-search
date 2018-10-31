package search

import (
	"time"

	"github.com/xuender/go-utils"
)

// Document 文档
type Document struct {
	Key      []byte    `json:"key"`               // 主键
	Title    string    `json:"title,omitempty"`   // 标题
	Summary  string    `json:"summary,omitempty"` // 摘要
	Content  string    `json:"content"`           // 内容
	Modified time.Time `json:"modified"`          // 修改时间
}

func (doc *Document) dbKey() []byte {
	return toDBKey(doc.Key)
}

func toDBKey(key []byte) []byte {
	return utils.PrefixBytes(key, _dbKey2docIDPrefix, '-')
}
