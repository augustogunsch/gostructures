package main

import (
	"fmt"
	"gostructures/tree"
)

func main() {
	t := tree.Create[int](1, 2, 3, 4, 5, 2)
	node, err := t.Find(8)
	fmt.Println(node, err)
}
