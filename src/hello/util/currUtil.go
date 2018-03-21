package util

import (
	"fmt"
	"math/rand"
	"strings"

	"hello/constant"
)

func Find(name string) {
	for i := 0; i < 10; i++ {
		c := constant.Currencies[i]
		switch {
		case strings.Contains(c.Currency, name),
			strings.Contains(c.Name, name),
			strings.Contains(c.Country, name):
			fmt.Println("Found", c)
		}
	}
}

func FindNumber(num int) {
	for _, curr := range constant.Currencies {
		if curr.Number == num {
			fmt.Println("Found", curr)
		}
	}
}

func FindAny(val interface{}) {
	switch i := val.(type) {
	case int:
		FindNumber(i)
	case string:
		Find(i)
	default:
		fmt.Printf("Unable to search with type %T\n", val)
	}
}

func AssertEuro(c constant.Curr) bool {
	switch name, curr := "Euro", "EUR"; {
	case c.Name == name:
		return true
	case c.Currency == curr:
		return true
	}
	return false
}

func Shuffle() {
	fmt.Println("... deferred Shuffle()")
	n := len(constant.Currencies)
	for i := range constant.Currencies {
		next := rand.Intn(n)
		temp := constant.Currencies[i]
		constant.Currencies[i] = constant.Currencies[next]
		constant.Currencies[next] = temp
	}
}

func Print() {
	defer Shuffle()
	fmt.Println("---- Currencies ----")
	for i, v := range constant.Currencies {
		fmt.Printf("%d: %v\n", i, v)
	}
}

func Sort() {
	fmt.Println("... Sort() ")
	N := len(constant.Currencies)
	for i := 0; i < N-1; i++ {
		currMin := i
		for k := i + 1; k < N; k++ {
			if constant.Currencies[k].Currency < constant.Currencies[currMin].Currency {
				currMin = k
			}
		}
		// swap
		if currMin != i {
			temp := constant.Currencies[i]
			constant.Currencies[i] = constant.Currencies[currMin]
			constant.Currencies[currMin] = temp
		}
	}
}
