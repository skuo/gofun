package util

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
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
			"%-10s %-10d %-10.3f\n",
			name, number, weight,
		)
	}

}

// ====================
// gob
type Name struct {
	First, Last string
}

type Book struct {
	Title       string    `json:"book_title"`
	PageCount   int       `json:"pages,string"`
	ISBN        string    `json:"-"`
	Authors     []Name    `json:"auths,omniempty"`
	Publisher   string    `json:",omniempty"`
	PublishDate time.Time `json:"pub_date"`
}

var books = []Book{
	Book{
		Title:       "Leaning Go",
		PageCount:   375,
		ISBN:        "9781784395438",
		Authors:     []Name{{"Vladimir", "Vivien"}},
		Publisher:   "Packt",
		PublishDate: time.Date(2016, time.July, 0, 0, 0, 0, 0, time.UTC),
	},
	Book{
		Title:       "The Go Programming Language",
		PageCount:   380,
		ISBN:        "9780134190440",
		Authors:     []Name{{"Alan", "Donavan"}, {"Brian", "Kernighan"}},
		Publisher:   "Addison-Wesley",
		PublishDate: time.Date(2015, time.October, 26, 0, 0, 0, 0, time.UTC),
	},
	Book{
		Title:       "Introducing Go",
		PageCount:   124,
		ISBN:        "978-1491941959",
		Authors:     []Name{{"Caleb", "Doxsey"}},
		Publisher:   "O'Reilly",
		PublishDate: time.Date(2016, time.January, 0, 0, 0, 0, 0, time.UTC),
	},
}

func tryGob() {
	fmt.Println("\n--- in tryGob() ---\n")

	// write books
	file, err := os.Create("output/book.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	enc := gob.NewEncoder(file)
	if err := enc.Encode(books); err != nil {
		fmt.Println(err)
	}

	// read in books
	fin, err := os.Open("output/book.dat")
	if err != nil {
		fmt.Println(err)
		return
	}

	var booksRead []Book
	dec := gob.NewDecoder(fin)
	if err := dec.Decode(&booksRead); err != nil {
		fmt.Println(err)
		return
	}

	for _, book := range booksRead {
		fmt.Println("gob", book)
	}
}

// ====================
// gzip
func tryGzip() {
	fmt.Println("\n--- in tryGzip() ---\n")
	filein, err := os.Open("src/hello/util/dataIo.go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer filein.Close()

	// zip content to output file
	fileout, err := os.Create("output/dataIo.go.gz")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fileout.Close()

	zip := gzip.NewWriter(fileout)
	zip.Name = fileout.Name()
	defer zip.Close()

	if count, err := io.Copy(zip, filein); err == nil {
		fmt.Printf("Gzipd file %s with %d bytes", fileout.Name(), count)
	} else {
		fmt.Println("Gzip failed:", err)
		os.Exit(1)
	}
}

// ====================
// json
func (n *Name) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s, %s\"", n.Last, n.First)), nil
}

func (n *Name) UnmarshalJSON(data []byte) error {
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	parts := strings.Split(name, ", ")
	n.Last, n.First = parts[0], parts[1]
	return nil
}

func tryJson() {
	fmt.Println("\n--- in tryJson() ---\n")
	// write json
	file, err := os.Create("output/book.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	enc := json.NewEncoder(file)
	if err := enc.Encode(books); err != nil {
		fmt.Println(err)
	}

	// read in books
	fin, err := os.Open("output/book.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// read json
	var booksRead []Book
	dec := json.NewDecoder(fin)
	if err := dec.Decode(&booksRead); err != nil {
		fmt.Println(err)
		return
	}

	for _, book := range booksRead {
		fmt.Println("json", book)
	}
}

// ====================
func TryDataIo() {
	tryBufio()
	tryFmt()
	tryGob()
	tryGzip()
	tryJson()
}
