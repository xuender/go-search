package search

import (
	"strings"

	"github.com/xuender/go-utils"
)

// Document 文档
type Document struct {
	Key          []byte   `json:"key"`               // 主键
	Title        string   `json:"title,omitempty"`   // 标题
	Summary      string   `json:"summary,omitempty"` // 摘要
	Content      string   `json:"content"`           // 内容
	TitleIndex   *Results `json:"-"`                 // 标题索引
	SummaryIndex *Results `json:"-"`                 // 摘要索引
	ContentIndex *Results `json:"-"`                 // 内容索引
}

func (doc *Document) dbKey() []byte {
	return toDBKey(doc.Key)
}

func toDBKey(key []byte) []byte {
	return utils.PrefixBytes(key, _dbKey2docIDPrefix, '-')
}

func (doc *Document) index() {
	if doc.Title != "" {
		doc.TitleIndex = Segmentation(doc.Title)
	}
	if doc.Summary != "" {
		doc.SummaryIndex = Segmentation(doc.Summary)
	}
	if doc.Content != "" {
		doc.ContentIndex = Segmentation(doc.Content)
	}
}

// Add 文档增加关键词
func (doc *Document) Add(word string) bool {
	return doc.TitleIndex.AddNum(word, strings.Count(doc.Title, word)) ||
		doc.SummaryIndex.AddNum(word, strings.Count(doc.Summary, word)) ||
		doc.ContentIndex.AddNum(word, strings.Count(doc.Content, word))
}

// AddByByte 增加关键则
func (doc *Document) AddByByte(word string, b byte) (c int) {
	switch b {
	case _titleIDPrefix:
		c = strings.Count(doc.Title, word)
		doc.TitleIndex.AddNum(word, c)
	case _summaryIDPrefix:
		c = strings.Count(doc.Summary, word)
		doc.SummaryIndex.AddNum(word, c)
	default:
		c = strings.Count(doc.Content, word)
		doc.ContentIndex.AddNum(word, c)
	}
	return
}
