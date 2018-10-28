package search

import (
	"fmt"

	"github.com/xuender/go-utils"
)

func ExampleNewLevelEngine() {
	e, err := NewLevelEngine("/tmp/search_db")
	defer e.Close()

	fmt.Println(err)

	// Output:
	// <nil>
}
func ExampleEngine_Has() {
	e, _ := NewLevelEngine("/tmp/search_db")
	defer e.Close()

	ok, err := e.Has([]byte("asd"))
	fmt.Println(ok, err)

	// Output:
	// false <nil>
}

func ExampleEngine_Put() {
	key := []byte("11")
	e, _ := NewLevelEngine("/tmp/search_db")
	defer e.Close()
	doc := &Document{
		Key:     key,
		Title:   "不生气",
		Summary: "生气引发的疾病",
		Content: "关于生气有，我有些话要说",
	}
	e.Put(doc)

	ok, err := e.Has(key)
	fmt.Println(ok, err)

	nd, err := e.Get(key)
	fmt.Println(nd.Title == doc.Title, err)

	// Output:
	// true <nil>
	// true <nil>
}

func ExampleEngine_Search() {
	e, _ := NewLevelEngine("/tmp/search_db")
	defer e.Close()

	doc := &Document{
		Key:     []byte("11"),
		Title:   "不生气",
		Summary: "生气引发的疾病",
		Content: "关于生气，我有些话要说",
	}
	e.Put(doc)

	fmt.Println(len(e.Search("生气")))
	fmt.Println(len(e.Search("高兴")))

	// Output:
	// 1
	// 0
}

func ExampleEngine_SearchKey() {
	e, _ := NewLevelEngine(fmt.Sprintf("/tmp/%s", utils.NewID('t')))
	defer e.Close()

	doc := &Document{
		Key:     []byte("11"),
		Title:   "不生气",
		Summary: "生气引发的疾病",
		Content: "关于生气，我有些话要说",
	}
	e.Put(doc)

	keys := e.SearchKey("疾病")
	fmt.Println(keys)

	keys = e.SearchKey("发病")
	fmt.Println(keys)

	// Output:
	// [[49 49]]
	// []
}
