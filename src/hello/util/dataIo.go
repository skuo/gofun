package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// ====================
// bufio
func tryBufio() {
	fmt.Println("\n--- in tryBufio() ---\n")
	// write with bytes
	var books bytes.Buffer
	books.WriteString("Mercury 4879 0 No\n")
	books.WriteString("Venus 12104 0 No\n")
	books.WriteString("Earth 12756 1 No\n")
	books.WriteString("Mars 6792 2 No\n")
	books.WriteString("Jupiter 142984 67 Yes\n")
	books.WriteString("Saturn 120536 62 Yes\n")
	books.WriteString("Uranus 51118 27 Yes\n")
	books.WriteString("Neptune 49528 14 Yes\n")
	books.WriteString("Pluto 2370 5 No\n")

	fout, err := os.Create("output/planets.txt")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		return
	}
	defer fout.Close()
	books.WriteTo(fout)

	//
	fin, err := os.Open("output/planets.txt")
	if err != nil {
		fmt.Println("Unable to open file:", err)
		return
	}
	defer fin.Close()

	fmt.Printf(
		"%-10s %-10s %-6s %-6s\n",
		"Planet", "Diameter", "Moons", "Ring?",
	)
	scanner := bufio.NewScanner(fin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		fmt.Printf(
			"%-10s %-10s %-6s %-6s\n",
			fields[0], fields[1], fields[2], fields[3],
		)
	}

}

// ====================
// fmt
type metalloid struct {
	name   string
	number int32
	weight float64
}

func tryFmt() {
	fmt.Println("\n--- in tryFmt() ---\n")
	var metalloids = []metalloid{
		{"Boron", 5, 10.81},
		{"Silicon", 14, 28.085},
		{"Germanium", 32, 74.63},
		{"Arsenic", 33, 74.921},
		{"Antimony", 51, 121.760},
		{"Tellerium", 52, 127.60},
		{"Polonium", 84, 209.0},
	}
	file, err := os.Create("output/metalloids.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	for _, m := range metalloids {
		fmt.Fprintf(
			file,
			"%-10s %-10d %-10.3f\n",
			m.name, m.number, m.weight,
		)
	}

	var name string
	var number int32
	var weight float64

	data, err := os.Open("output/metalloids.txt")
	if err != nil {
		fmt.Println("Unable to open metalloid data:", err)
		return
	}
	defer data.Close()

	fmt.Printf(
		"%-10s %-10s %-6s\n",
		"Metalloid", "number", "weight",
	)
	for {
		_, err := fmt.Fscanf(
			data,
			"%s %d %f\n",
			&name, &number, &weight,
		)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Scan error:", err)
				return
			}
		}
		fmt.Printf(
			"%-10s %-10d %-6.3f\n",
			name, number, weight,
		)
	}

}

// ====================
func TryDataIo() {
	tryBufio()
	tryFmt()
}
