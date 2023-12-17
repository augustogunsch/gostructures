package main

import (
	"fmt"
	"gostructures/tree"
)

func main() {
	t := tree.Create[int](1, 2, 3, 4, 5, 2)
	fmt.Println(t)
	newT, err := t.RemoveWith(2)
	t = newT
	fmt.Println(t, err)
	newT, err = t.RemoveWith(2)
	t = newT
	fmt.Println(t, err)
}
