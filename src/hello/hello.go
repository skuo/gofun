package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"hello/constant"
	"hello/stringutil"
	"hello/util"
)

func tryReverse() {
	fmt.Printf("--- tryReverse() ---\n")
	fmt.Printf("hello, world \n")
	fmt.Printf(stringutil.Reverse("hello, world") + "\n")
	var str string = "This is to test eclipse's variables display\n"
	shortVarDec := "Lose var and type\n"
	fmt.Printf("%s%s", str, shortVarDec)
}

func tryConstant() {
	fmt.Printf("\n\n--- tryConstant() ---\n")
	fmt.Printf("StarHyperGiant = %v [%T]\n", constant.StarHyperGiant, constant.StarHyperGiant)
	fmt.Printf("StarSuperGiant = %v [%T]\n", constant.StarSuperGiant, constant.StarSuperGiant)
	fmt.Printf("StarBrightGiant = %v [%T]\n", constant.StarBrightGiant, constant.StarBrightGiant)
	fmt.Printf("StarDwarf = %v [%T]\n", constant.StarDwarf, constant.StarDwarf)
	fmt.Printf("StarRedDwarf = %v [%T]\n", constant.StarRedDwarf, constant.StarRedDwarf)
	fmt.Printf("StarBrownDwarf = %v [%T]\n", constant.StarBrownDwarf, constant.StarBrownDwarf)
}

func tryNummap() {
	fmt.Printf("\n\n--- tryNummap() ---\n")
	// create and store bimap in "nummap.txt"
	max := 10
	fileMode := 4000
	mapFileName := "output/nummap.txt"

	nummap := util.MakeBitMap(max)
	err := ioutil.WriteFile(mapFileName, []byte(nummap), os.FileMode(fileMode))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Bitmap file", mapFileName, "created OK")

	// Read from mapFileName, strconv i to string, write to numbersFile
	nums, err := util.LoadNumberMap(mapFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	numbersFile := "output/nums.txt"
	err = ioutil.WriteFile(numbersFile, nums.Bytes(), os.FileMode(fileMode))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Created numers data file", numbersFile, "OK.")

}

func tryGoto() {
	fmt.Printf("\b\n--- tryGoto() ---\n")
	var a string
Start:
	for {
		switch {
		case a < "aaa":
			goto A
		case a >= "aaa" && a < "aaabbb":
			goto B
		case a == "aaabbb":
			break Start
		}
	A:
		a += "a"
		continue Start
	B:
		a += "b"
		continue Start
	}
	fmt.Println(a)
}

func tryCurr() {
	fmt.Printf("\b\n--- tryCurr() ---\n")
	// find
	util.Find("Dinar")
	util.Find("HTG")
	util.Find("Hong Kong")

	util.FindAny("Peso")
	util.FindAny(404)
	util.FindAny(978)
	util.FindAny(false)

	curr1 := constant.Curr{"EUR", "Euro", "Italy", 978}
	if util.AssertEuro(curr1) {
		fmt.Println(curr1, "is Euro")
	}
	
	// sort, defer shuffle, print
	util.Sort()
	util.Print()
	util.Print()
}

func main() {
	tryReverse()
	tryConstant()
	tryNummap()
	tryGoto()
	tryCurr()
}
