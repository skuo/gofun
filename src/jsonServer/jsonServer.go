package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"jsonServer/curr1"
)

var currencies = curr1.Load("data/curr1.csv")

// api endpoint for service
// input of this form: {"get" : "Yen"}
func currs(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("URL", req.URL)
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
	addCookie(resp, req)
	path := getPwd() + "/src/jsonServer/currency.html"
	file, err := os.Open(path)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	io.Copy(resp, file)
}

func addCookie(resp http.ResponseWriter, req *http.Request) {
	// add cookie
	fmt.Println("Add testcookiename")
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "testcookiename", Value: "testcookievalue", Path: "/", Expires: expire, MaxAge: 86400}

	http.SetCookie(resp, &cookie)
}

func listCookies(resp http.ResponseWriter, req *http.Request) {
	cookies := req.Cookies()
	for _, cookie := range cookies {
		//fmt.Println("name=", cookie.Name, "value=", cookie.Value)
		s := fmt.Sprintf("name=%s, value=%s", cookie.Name, cookie.Value)
		io.WriteString(resp, s)
	}
}

func deleteCookie(resp http.ResponseWriter, req *http.Request) {
	// set cookie's MaxAge to -1
	cookie := http.Cookie{Name: "testcookiename", Path: "/", MaxAge: -1}
	http.SetCookie(resp, &cookie)

	io.WriteString(resp, "cookie testcookiename deleted")
}

func getPwd() string {
	// find cwd and strip out /src/jsonServer if necessary
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pwd)
	index := strings.Index(pwd, "/src/jsonServer")
	if index != -1 {
		pwd = pwd[0:index]
	}
	return pwd
}

func main() {
	mux := http.NewServeMux()
	// starting gui
	mux.HandleFunc("/gui", gui)
	// handle all variations under /currency/
	mux.HandleFunc("/currency/", currs)
	// cookie pages
	mux.HandleFunc("/cookie/delete", deleteCookie)
	mux.HandleFunc("/cookie/list", listCookies)

	fmt.Println("Starting http server")
	if err := http.ListenAndServe(":4040", mux); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Does it get here?")
}
