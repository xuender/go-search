package search

import (
	"fmt"
)

func ExampleSplit() {
	fmt.Println(Split("中国China山东省,威海市 环翠区。火炬"))

	// Output:
	// [中国 山东省 china 威海市 环翠区 火炬]
}

func ExampleGroup() {
	fmt.Println(Group("太阳当空照"))

	// Output:
	// [太阳 太阳当 太阳当空 太阳当空照 阳当 阳当空 阳当空照 当空 当空照 空照]
}

func ExamplePosition() {
	fmt.Println(Position("abc cba aac", "a"))

	// Output:
	// [0 6 8 9]
}
