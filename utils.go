package search

import (
	"regexp"
	"strings"
)

var (
	// 标点分割
	// \s 空白符号
	// [:punct:] 英文标点
	// \pP unicode标点
	_Punct = regexp.MustCompile(`[\s[:punct:]\pP]+`)
	// 中文
	_Han = regexp.MustCompile(`[\p{Han}]+`)
	// 日文
	_Jan = regexp.MustCompile(`[\p{Hiragana}\p{Katakana}]+`)
	// 其他
	_Other = regexp.MustCompile(`[^\p{Han}\p{Hiragana}\p{Katakana}]+`)
)

// Split 字符串拆分
func Split(str string) []string {
	ret := []string{}
	for _, n := range _Punct.Split(strings.ToLower(str), -1) {
		// 增加中文
		ret = append(ret, _Han.FindAllString(n, -1)...)
		// 增加日文
		ret = append(ret, _Jan.FindAllString(n, -1)...)
		// 增加其他
		ret = append(ret, _Other.FindAllString(n, -1)...)
	}
	return ret
}

// IsNaturalSplit 自然分割的词
func IsNaturalSplit(str string) bool {
	return _Other.MatchString(str)
}

// Group 句子分组
func Group(str string) []string {
	ret := []string{}
	if len(str) < 2 {
		return ret
	}
	strs := strings.Split(str, "")
	for i := 0; i < len(strs)-1; i++ {
		for f := i + 2; f <= len(strs); f++ {
			if f-i > WordMaxLength {
				break
			}
			ret = append(ret, strings.Join(strs[i:f], ""))
		}
	}
	return ret
}

// Position str 在 text 中所有位置
func Position(text, str string) []int {
	ret := []int{}
	for i, f := strings.Index(text, str), 0; i > -1; i = strings.Index(text[f:], str) {
		f += i
		ret = append(ret, f)
		f++
	}
	return ret
}
