package main

import (
	"fmt"
	"hello/stringutil"
)

func main() {
	fmt.Printf("hello, world \n")
	fmt.Printf(stringutil.Reverse("hello, world") + "\n")
	var str = "This is to test eclipse's variables display\n"
	fmt.Printf(str)
}
