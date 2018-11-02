package search

import (
	"encoding/json"
	"fmt"
)

func ExampleDocument() {
	doc := Document{
		Key:     []byte("123"),
		Title:   "标题",
		Summary: "简介",
		Content: "内容",
	}
	j, _ := json.Marshal(doc)

	fmt.Println(string(j))

	// Output:
	// {"key":"MTIz","title":"标题","summary":"简介","content":"内容","modified":"0001-01-01T00:00:00Z"}
}

func ExampleDocument_Inverted() {
	doc := Document{
		Key:     []byte("123"),
		Title:   "abca",
		Summary: "",
		Content: "cca",
	}

	fmt.Println(doc.Inverted("a"))

	// Output:
	// true [0 3 6]
}
