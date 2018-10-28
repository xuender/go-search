package search

import (
	"sort"
	"strings"
)

// Results 统计结果
type Results map[string]int

// Entry map内容
type Entry struct {
	Key   string
	Value int
}

// Top 顺序返回结果大于0的
func (r Results) Top() []string {
	es := []Entry{}
	for k, v := range r {
		if v > 0 {
			es = append(es, Entry{Key: k, Value: v})
		}
	}
	sort.Slice(es, func(i, j int) bool { return es[j].Value < es[i].Value })
	ret := []string{}
	for _, e := range es {
		ret = append(ret, e.Key)
	}
	return ret
}

// Add 增加统计
func (r Results) Add(s string) {
	if num, ok := r[s]; ok {
		r[s] = num + 1
	} else {
		r[s] = 1
	}
}

// Get 获取统计结果
func (r Results) Get(s string) int {
	if num, ok := r[s]; ok {
		return num
	}
	return 0
}

// Has 包含
func (r Results) Has(s string) bool {
	_, has := r[s]
	return has
}

// AddNum 增加数量
func (r Results) AddNum(s string, num int) bool {
	if num != 0 {
		r[s] += num
		return true
	}
	return false
}

// Contain 包含任意一个词
func (r Results) Contain(results *Results) bool {
	for k := range r {
		if results.Has(k) {
			return true
		}
	}
	return false
}

// Segmentation 内容分割
func Segmentation(content string) *Results {
	ret := Results{}
	for _, s := range Split(strings.ToLower(content)) {
		if _Other.MatchString(s) {
			ret.Add(s)
		} else {
			for _, n := range strings.Split(s, "") {
				ret.Add(n)
			}
		}
	}
	return &ret
}
