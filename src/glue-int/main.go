package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"logit"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Define our struct for authentication
type user_t struct {
	user     string
	password string
	group    string
}
type authenticationMiddleware_t struct {
	tokenUsers map[string]user_t
}

var amw authenticationMiddleware_t

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type session_t struct {
	User  string
	Token string
	Group string
}

// expert flags and constants for logging
const cSHOWREQUESTHDRS int32 = 0x01 // show request details
const cSHOWENDPOINT int32 = 0x02    // show requested endpoint

var myFlags logit.DFlags_t // holds the logger flags
var mStats runtime.MemStats

func main() {
	// Setup the logging and get the initial stats
	err := logit.OpenLog("") // use default config file
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer closeMain()
	logit.GetMyLogInfo(&myFlags)
	// Now run some simple tests
	logit.Infof(&myFlags, "golang version: '%s'", runtime.Version())
	logit.Infof(&myFlags, "Available CPUs:%d", runtime.NumCPU())
	logit.Infof(&myFlags, "HTTP server process id = %d", syscall.Getpid())
	runtime.ReadMemStats(&mStats)
	//
	var dir string
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	logit.Infof(&myFlags, "Starting glue, page directory is '%s'", dir)
	//
	gob.Register(&session_t{})
	// setup the router and subrouters
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/login").Methods("GET").HandlerFunc(loginPageHandler)
	router.PathPrefix("/login").Methods("POST").HandlerFunc(loginAuthenticate)
	// Now setup the sub-routers
	p1 := router.PathPrefix("/static").Subrouter()
	p1.Methods("GET").HandlerFunc(p1Handler)
	p2 := router.PathPrefix("/dynamic").Subrouter()
	p2.Methods("GET").HandlerFunc(p2Handler)

	// attach middleware authentication to router
	//amw := authenticationMiddleware_t{}
	amw.populateUsers() // populate with users and passwords
	router.Use(amw.middlewareAuthorization)
	//
	// make a channel to notify main when to shutdown
	//
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	//
	// Start the server on another thread
	// This will serve files under http://localhost:8080/....
	//
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		logit.Infof(&myFlags, "Starting server, listening on '%s'", srv.Addr)
		err = srv.ListenAndServe() // this will run...
		if err != nil {
			if strings.ToLower(err.Error()) != "http: server closed" {
				logit.Fatalf(&myFlags, "HTTP server returned with error: '%s'", err.Error())
			}
		}
		logit.Infof(&myFlags, "Server stopped, no longer listening on '%s'", srv.Addr)
	}()
	logit.Info(&myFlags, "Main blocked on interrupt (-2) signal.")
	//
	// block on a stop request from signal
	// This does NOT work when in debug mode for some reason
	//
	<-stop
	//time.Sleep(time.Second * 5) // just for testing
	logit.Info(&myFlags, "Stop signal received")
	if err := srv.Shutdown(context.Background()); err != nil {
		logit.Fatalf(&myFlags, "could not shutdown: %v", err)
	}
	time.Sleep(time.Second * 1) // Some time to let the background processes wrap up
}

/*
  closeMain
  post the final stats and close the logger
*/
func closeMain() {
	var mStats2 runtime.MemStats
	runtime.ReadMemStats(&mStats2)
	logit.Debugf(&myFlags, "Memory allocs %d, Frees '%d'",
		mStats2.Mallocs-mStats.Mallocs, mStats2.Frees-mStats.Frees)
	//
	logit.CloseLog()
}

/*
  Authentication middleware will block further penetration into the code
  if the token is not present.
*/
func (amw *authenticationMiddleware_t) middlewareAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//token := r.Header.Get("X-Session-Token")
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logit.Warn(&myFlags, "Problem returning session name from store.")
			return
		}

		if session.IsNew { // no previous session, check if login
			if strings.ToLower(r.RequestURI) == "/login" { // wants login page
				if r.Method == "GET" { // asking for login page
					logit.Debug(&myFlags, "Login page request")
					next.ServeHTTP(w, r)
					return
				} else if r.Method == "POST" { // asking for authentication
					logit.Debug(&myFlags, "Login authentication request")
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Sorry, Forbidden Page", 403)
			logit.Warn(&myFlags, "Request for page from unauthorized source.")
		} else { // a session is present in header
			value := session.Values["session"]
			var session = &session_t{}
			session, sessionOk := value.(*session_t)
			if !sessionOk {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				logit.Warn(&myFlags, "Request for page from unknown session")
				return
			}
			if !amw.checkUser(session.User, session.Token) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				logit.Warnf(&myFlags, "Request for page from unauthorized user: '%s'", session.User)
				return
			}
			// We found the user/token in our map
			logit.Infof(&myFlags, "Request for page from authorized user: '%s'", session.User)
			next.ServeHTTP(w, r)
			return // no error
		}
	})
}

/*
  a GET request has come in for the login page
  return the page to the caller
*/
func loginPageHandler(wtr http.ResponseWriter, rdr *http.Request) {
	const loginPage string = `
	<!DOCTYPE html>
	<html>
		<head>
			<title>Login Page</title>
		</head>
		<body>
			<h3>Welcome! Please enter login credentials.</h3>
			<form action="/login" method="POST">
				Username:&emsp;
				<input type="text" name="username" value=""><br>
				Password:&emsp;&nbsp; 
				<input type="password" name="password" value=""><br><br>
				<input type="submit" name="submit" value="Submit">
			</form>
		</body>
	</html>
	`
	wtr.Write([]byte(loginPage))
}

/*
  a POST request has come in to authenticate the user and password
  return the index.html page or error if not good creds.
*/
func loginAuthenticate(wtr http.ResponseWriter, rdr *http.Request) {
	logit.Debugfx(cSHOWREQUESTHDRS, &myFlags, "Request Header:\n%s", formatRequest(rdr))
	var userName string = ""
	var passWord string = ""
	payload := strings.Split(rdr.Form.Encode(), "&")
	for _, pairs := range payload {
		keyValues := strings.Split(pairs, "=")
		if len(keyValues) == 2 {
			key := strings.ToLower(keyValues[0])
			if key == "username" {
				userName = keyValues[1]
			} else if key == "password" {
				passWord = keyValues[1]
			}
		}
	}
	if (len(userName) == 0) || (len(passWord) == 0) {
		http.Error(wtr, "Not authorized, no user id, and/or password", 401)
		logit.Warnf(&myFlags, "Not authorized, no user id, and/or password in login request")
		return
	}
	if (len(userName) > 0) && (len(passWord) > 0) {
		logit.Debugf(&myFlags, "User:'%s' with password found inside request", userName)
		if amw.checkUser(userName, passWord) {
			logit.Infof(&myFlags, "User:'%s' validated", userName)
			// get a new session for this user
			newSession, err := store.Get(rdr, "session-name")
			if err != nil {
				http.Error(wtr, err.Error(), http.StatusInternalServerError)
				return
			}
			// setup the session structure
			user, _ := amw.tokenUsers[userName]
			var session session_t
			session.Group = user.group
			session.User = user.user
			session.Token = user.password
			newSession.Values["session"] = session
			// Save it before we write to the response/return from the handler.
			err = newSession.Save(rdr, wtr)
			if err != nil {
				http.Error(wtr, err.Error(), http.StatusInternalServerError)
				return
			}
			// get the index.html page
			if user.group == "dev" {
				rdr.Method = "GET"
				rdr.URL.Path = ""
				handler := http.FileServer(http.Dir(""))
				handler.ServeHTTP(wtr, rdr)
			} else { // admin page
				page, err := ioutil.ReadFile("indexA.html")
				if err != nil {
					http.Error(wtr, err.Error(), http.StatusInternalServerError)
					return
				}
				wtr.Write(page)
			}
			return
		}
	}
	http.Error(wtr, "Not authorized, bad user id, or password", 401)
	logit.Warn(&myFlags, "Not authorized, bad user id, or password in login request")
}

func p1Handler(wtr http.ResponseWriter, rdr *http.Request) {
	logit.Debugfx(cSHOWENDPOINT, &myFlags, "P1 Endpoint request:'%s'", rdr.RequestURI)
	handler := http.Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(""))))
	//handler.ServeHTTP(wtr, rdr)
	gzHandler := gziphandler.GzipHandler(handler)
	gzHandler.ServeHTTP(wtr, rdr)
}

func p2Handler(wtr http.ResponseWriter, rdr *http.Request) {
	logit.Debugfx(cSHOWENDPOINT, &myFlags, "P2 Endpoint request:'%s'", rdr.RequestURI)
	handler := http.Handler(http.StripPrefix("/dynamic/", http.FileServer(http.Dir(""))))
	//handler.ServeHTTP(wtr, rdr)
	gzHandler := gziphandler.GzipHandler(handler)
	gzHandler.ServeHTTP(wtr, rdr)
}

/*
  Initialize user/password table
  Just temporary until ldap is setup
*/
func (amw *authenticationMiddleware_t) populateUsers() {
	amw.tokenUsers = make(map[string]user_t)
	amw.tokenUsers["admin"] = user_t{"admin", "admin", "admin"}
	amw.tokenUsers["jeff"] = user_t{"jeff", "1234", "dev"}
	amw.tokenUsers["steve"] = user_t{"steve", "ssss", "dev"}
	amw.tokenUsers["david"] = user_t{"david", "dddd", "dev"}
}

/*
  check to see if the user/password is known
*/
func (amw *authenticationMiddleware_t) checkUser(userName string, token string) bool {
	user, found := amw.tokenUsers[userName]
	if !found {
		return false
	}
	if user.password == token {
		return true
	}
	return false
}

/*
  This generates a string that can be printed or logged.
  WARNING: This will show the user name and password request.
*/
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

/*
func AuthenticationHandler(wtr http.ResponseWriter, rdr *http.Request) {
	_, _, ok := rdr.BasicAuth()
	if !ok { // If the values are not present
		http.Error(wtr, "Not authorized", 401)
		return
	}
}
*/
