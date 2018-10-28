package search

import (
	"fmt"

	"github.com/xuender/go-utils"
)

func ExampleKeyword() {
	k := NewKeyword("你不高兴ok?")
	utils.ReadLines(
		"./data/中国五十年儿童文学名家作品选.txt",
		// "./data/成语.txt",
		func(l string) {
			r := Segmentation(l)
			if k.results.Contain(r) {
				k.Read(l)
			}
		})

	fmt.Print(k.Get())
	// Output:
	// [高兴 你不 不高 不高兴]
}
