package search

import (
	"strings"
)

// Keyword 关键词
type Keyword struct {
	str     string   // 寻找关键词的输入
	results *Results // 词统计
	words   *Results // 关键词列表
}

// Read 读取文档
func (k *Keyword) Read(doc string) {
	if doc == "" {
		return
	}
	for w := range *(k.words) {
		k.words.AddNum(w, strings.Count(doc, w))
	}
}

// Get 读取关键词
func (k *Keyword) Get() []string {
	return k.words.Top()
}

// NewKeyword 新建关键词
func NewKeyword(str string) *Keyword {
	words := Results{}
	ws := Split(str)
	for _, w := range ws {
		if IsNaturalSplit(w) {
			words[w] = 0
		} else {
			for _, s := range Group(w) {
				words[s] = 0
			}
		}
	}
	return &Keyword{
		str:     str,
		words:   &words,
		results: Segmentation(str),
	}
}
