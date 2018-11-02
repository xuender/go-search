package search

import (
	"strings"
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

// Match 文档匹配
func (doc *Document) Match(str string) bool {
	return (doc.Title != "" && strings.Contains(doc.Title, str)) ||
		(doc.Summary != "" && strings.Contains(doc.Summary, str)) ||
		(doc.Content != "" && strings.Contains(doc.Content, str))
}

// Inverted 生成到排索引
func (doc *Document) Inverted(str string) (bool, []int) {
	tp := Position(doc.Title, str)
	sp := Position(doc.Summary, str)
	cp := Position(doc.Content, str)
	if len(tp) == 0 && len(sp) == 0 && len(cp) == 0 {
		return false, nil
	}
	tl := len(doc.Title)
	for i, v := range sp {
		sp[i] = v + tl
	}
	sl := len(doc.Summary)
	for i, v := range cp {
		cp[i] = v + tl + sl
	}
	tp = append(tp, sp...)
	tp = append(tp, cp...)
	return true, tp
}

func toDBKey(key []byte) []byte {
	return utils.PrefixBytes(key, _dbKey2docIDPrefix, '-')
}
