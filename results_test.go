package search

import (
	"fmt"
)

func ExampleSegmentation() {
	m := Segmentation("不着急，不上火，不感冒，不发烧。OK?")
	fmt.Println(m.Get("不"))
	fmt.Println(m.Get("ok"))

	// Output:
	// 4
	// 1
}

func ExampleResults_Contain() {
	r := Results{"不着急": 0, "不上火": 0, "不感冒": 0}
	m := Results{"不着急": 0}
	d := Results{"看书": 0}
	fmt.Println(r.Contain(&m))
	fmt.Println(r.Contain(&d))

	// Output:
	// true
	// false
}
