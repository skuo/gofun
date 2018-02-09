package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"jsonServer/curr1"
)

var currencies = curr1.Load("data/curr1.csv")

// api endpoint for service
func currs(resp http.ResponseWriter, req *http.Request) {
    fmt.Println("URL", req.URL);
	var currRequest curr1.CurrencyRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&currRequest); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	result := curr1.Find(currencies, currRequest.Get)
	enc := json.NewEncoder(resp)
	if err := enc.Encode(&result); err != nil {
		fmt.Println(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// serves HTML gui
func gui(resp http.ResponseWriter, req *http.Request) {
	file, err := os.Open("src/jsonServer/currency.html")
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	io.Copy(resp, file)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", gui)
	mux.HandleFunc("/currency", currs)

    fmt.Println("Starting http server")
	if err := http.ListenAndServe(":4040", mux); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Does it get here?")
}
