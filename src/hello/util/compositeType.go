package util

import (
	"fmt"
	"math/rand"
	"time"
)

// =====================
// Array
const size = 1024 * 1024

type numbers [size]int

func initialize(nums *numbers) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		nums[i] = rand.Intn(10000)
	}
}

func max(nums *numbers) int {
	temp := nums[0]
	for _, val := range nums {
		if val > temp {
			temp = val
		}
	}
	return temp
}

func tryArray() {
	var nums *numbers = new(numbers)
	initialize(nums)
	///fmt.Println(nums)
	fmt.Println("Max num is: ", max(nums))
}

// =========================
// Slice
func scale(factor float64, vector []float64) []float64 {
	for i, _ := range vector {
		vector[i] *= factor
	}
	return vector
}

func scale2(factor float64, vector []float64) {
	for i, _ := range vector {
		vector[i] *= factor
	}
}

func join(v1, v2 []float64) []float64 {
	return append(v1, v2...)
}

func clone(v []float64) (result []float64) {
	result = make([]float64, len(v), cap(v))
	copy(result, v)
	return
}

// change name from add() to addSlice() to avoid conflict
// with the add() in funcs.go.  
func addSlice(v1, v2 []float64) []float64 {
	if len(v1) != len(v2) {
		panic("Size mismatch")
	}
	result := make([]float64, len(v1))
	for i := range v1 {
		result[i] = v1[i] + v2[i]
	}
	return result
}

func trySlice() {
	h := []float64{12.5, 18.4, 7.0}
	h[0] = 15
	fmt.Println(h[0])

	h10 := scale(2.0, h)
	fmt.Println("h10",h10)
	fmt.Println("h",h)

	h2 := clone(h)
	fmt.Println("addSlice(h,h2)",addSlice(h, h2))
	fmt.Println("h2", h2)
	scale2(0.5, h2)
	fmt.Println("h2", h2)
}

// ============================
func TryCompositeType() {
	tryArray()
	trySlice()
}
