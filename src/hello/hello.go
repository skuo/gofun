package main

import (
	"fmt"
	"hello/stringutil"
)

func main() {
	fmt.Printf("hello, world \n")
	fmt.Printf(stringutil.Reverse("hello, world") + "\n")
	var str string = "This is to test eclipse's variables display\n"
	shortVarDec := "Lose var and type"
	fmt.Printf("%s %s\n", str, shortVarDec)
}
