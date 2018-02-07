package util

import (
	"fmt"
	"math/rand"
	"time"
)

// =====================
// Array.
// Pass pointer for function param
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
// Slice is already a pointer.
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
	fmt.Println("h10", h10)
	fmt.Println("h", h)

	h2 := clone(h)
	fmt.Println("addSlice(h,h2)", addSlice(h, h2))
	fmt.Println("h2", h2)
	scale2(0.5, h2)
	fmt.Println("h2", h2)
}

// =========================
// Map
// Map is already a pointer

// saves into map, blows up upon dup key
func save(store map[string]int, key string, value int) {
	val, ok := store[key]
	if !ok {
		store[key] = value
	} else {
		panic(fmt.Sprintf("Slot %d taken", val))
	}
}

// removes an entry, error if not found
func remove(store map[string]int, key string) error {
	_, ok := store[key]
	if !ok {
		return fmt.Errorf("Key not found")
	}
	delete(store, key)
	return nil
}

func tryMap() {
	hist := make(map[string]int, 6)
	hist["Jan"] = 100
	hist["Feb"] = 445
	hist["Mar"] = 514
	hist["Apr"] = 233
	hist["May"] = 321
	hist["Jun"] = 644
	hist["Jul"] = 113
	save(hist, "Aug", 734)
	save(hist, "Sep", 553)
	save(hist, "Oct", 344)
	save(hist, "Nov", 831)
	save(hist, "Dec", 312)
	save(hist, "Dec0", 332)
	remove(hist, "Dec0")

	for key, val := range hist {
		adjVal := int(float64(val) * 0.100)
		fmt.Printf("%s (%d):", key, val)
		for i := 0; i < adjVal; i++ {
			fmt.Print(".")
		}
		fmt.Println()
	}
}

// ============================
func TryCompositeType() {
	tryArray()
	trySlice()
	tryMap()
}
