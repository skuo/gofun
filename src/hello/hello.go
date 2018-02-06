package main

import (
	"fmt"

	"hello/constant"
	"hello/stringutil"
)

func main() {
	fmt.Printf("hello, world \n")
	fmt.Printf(stringutil.Reverse("hello, world") + "\n")
	var str string = "This is to test eclipse's variables display\n"
	shortVarDec := "Lose var and type\n"
	fmt.Printf("%s%s\n", str, shortVarDec)

	// print out constants
	fmt.Printf("StarHyperGiant = %v [%T]\n", constant.StarHyperGiant, constant.StarHyperGiant)
	fmt.Printf("StarSuperGiant = %v [%T]\n", constant.StarSuperGiant, constant.StarSuperGiant)
	fmt.Printf("StarBrightGiant = %v [%T]\n", constant.StarBrightGiant, constant.StarBrightGiant)
	fmt.Printf("StarDwarf = %v [%T]\n", constant.StarDwarf, constant.StarDwarf)
	fmt.Printf("StarRedDwarf = %v [%T]\n", constant.StarRedDwarf, constant.StarRedDwarf)
	fmt.Printf("StarBrownDwarf = %v [%T]\n", constant.StarBrownDwarf, constant.StarBrownDwarf)

}
