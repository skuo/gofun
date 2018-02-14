package main

import (
	"log"
	"os"

	"addressbook/util"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage:  %s ADDRESS_BOOK_FILE list|add \n", os.Args[0])
	}
	fname := os.Args[1]
	option := os.Args[2]

	if option == "add" {
        // Add an address
        util.AddAddress(fname)
    } else if option == "list" {
        // List people in an address book
        util.ListPeople(fname)
    }
}
