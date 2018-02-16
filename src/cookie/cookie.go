package main

import (
	"io"
	"net/http"
	"time"
)

// addCookie will apply a new cookie to the response of a http
// request, with the key/value this method is passed.
func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	addCookie(w, "TestCookieName", "TestValue")
	io.WriteString(w, "Hello world!")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
