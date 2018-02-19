package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"jsonServer/curr1"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
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
	path := getProjDir() + "/src/jsonServer/currency.html"
	file, err := os.Open(path)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	io.Copy(resp, file)
}

// serves HTML plotly
func plotly(resp http.ResponseWriter, req *http.Request) {
	addCookie(resp, req)
	path := getProjDir() + "/src/jsonServer/plotly.html"
	file, err := os.Open(path)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	io.Copy(resp, file)
}

// serves HTML scatter
func scatter(resp http.ResponseWriter, req *http.Request) {
	addCookie(resp, req)
	path := getProjDir() + "/src/jsonServer/scatter.html"
	file, err := os.Open(path)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	io.Copy(resp, file)
}

// serves HTML scatter_gz
func scatter_gz(resp http.ResponseWriter, req *http.Request) {
	addCookie(resp, req)
	path := getProjDir() + "/src/jsonServer/scatter_gz.html"
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
		s := fmt.Sprintf("name=%s, value=%s\n", cookie.Name, cookie.Value)
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
	return pwd
}

func getProjDir() string {
	projDir := getPwd()
	index := strings.Index(projDir, "/src/jsonServer")
	if index != -1 {
		projDir = projDir[0:index]
	}
	return projDir
}

func listHeaders(resp http.ResponseWriter, req *http.Request) {
	for key, value := range req.Header {
		io.WriteString(resp, fmt.Sprintf("%v: %v\n", key, value))
	}
}

func simpleText(resp http.ResponseWriter, req *http.Request) {
	http.ServeFile(resp, req, "simple.txt")
}

// serveFile causes the browser the download the file, not displaying the content
func simpleTextJgz(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Enconding", "gzip")
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	http.ServeFile(resp, req, "simple.txt.jgz")
}

// This also cause the browser the download the file instead of display the content
func simpleGzText(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Enconding", "gzip")
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	b, err := ioutil.ReadFile("simple.gz.txt")
	if err != nil {
		http.NotFound(resp, req)
	}
	resp.Write(b)
}

func makeHandler(withoutGzHandler http.Handler, withGzHandler http.Handler) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// Figure out the requrest filename
		gzFullPath := getPwd() + req.URL.Path + ".gz"
		fmt.Println("gzFullPath=", gzFullPath)
		if _, err := os.Stat(gzFullPath); err == nil {
			// Note.  This does not work.  The content is still gzipped (or even double gzipped)
			// A pre gzipped file exists.  Change req.URL.Path, set header and serve
			req.URL.Path = req.URL.Path + ".gz"
			resp.Header().Add("Content-Enconding", "gzip")
			resp.Header().Add("Content-Type", "text/csv; charset=utf-8")
			//resp.Header().Add("vary", "Accept-Encoding")
			//resp.Header().Add("Transfer-Encoding", "chunked") // chunked
			withoutGzHandler.ServeHTTP(resp, req)
		} else {
			// gzip on the fly
			withGzHandler.ServeHTTP(resp, req)
		}
	}
}

func main() {
	var dir string

	flag.StringVar(&dir, "dir", getPwd(), "the directory to serve files from. Defaults to the pwd")
	flag.Parse()
	fmt.Println("dir=", dir)

	mux := http.NewServeMux()
	// starting gui
	mux.HandleFunc("/gui", gui)
	// handle all variations under /currency/
	mux.HandleFunc("/currency/", currs)
	// cookie pages
	mux.HandleFunc("/cookie/delete", deleteCookie)
	mux.HandleFunc("/cookie/list", listCookies)
	// header page
	mux.HandleFunc("/header/list", listHeaders)
	// simple
	mux.HandleFunc("/simpleText", simpleText)
	mux.HandleFunc("/simpleTextJgz", simpleTextJgz)
	mux.HandleFunc("/simpleGzText", simpleGzText)
	// plotly
	mux.HandleFunc("/plotly", plotly)
	mux.HandleFunc("/scatter", scatter)
	mux.HandleFunc("/scatter_gz", scatter_gz)
	// FileServer
	fs := http.FileServer(http.Dir(dir + "/static"))
	// Gzip handles
	fsWithoutGzHandler := http.StripPrefix("/static", fs)
	//mux.Handle("/static/", fsWithoutGzHandle)
	fsWithGzHandler := gziphandler.GzipHandler(fsWithoutGzHandler)
	staticHandler := makeHandler(fsWithoutGzHandler, fsWithGzHandler)
	//mux.Handle("/static/", fsWithGzHandler)
	mux.Handle("/static/", staticHandler)

	fmt.Println("Starting http server")
	if err := http.ListenAndServe(":4040", mux); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Does it get here?")
}
