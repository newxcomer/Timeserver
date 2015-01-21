// Simple program to run a server showing a web page
// displaying the current time of day
// Personalize with using cookie for the user name
// Example usage: http://localhost:8080/time
// Default port number is 8080
// Optional flag -V printing version number to console output
// Optional flag -port for port number
//
//	Author: Binh Nguyen
//	Date: 1/20/2015
//	Version: 0.0015

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	version = "Timeserver version 0.0015"
	users   = map[string]string{}
	mutex   = &sync.Mutex{}
)

// Display the current time
func timeHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("In time handler")
	// Check the path
	if request.URL.Path != "/time" && request.URL.Path != "/time/" {
		errorHandler(response, request)
		return
	}

	// Current time format
	const timeFormat = "3:04:05 PM"

	// Print the response
	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<head>")
	fmt.Fprintln(response, "<style>")
	fmt.Fprintln(response, "p {font-size: xx-large}")
	fmt.Fprintln(response, "span.time {color: red}")
	fmt.Fprintln(response, "</style>")
	fmt.Fprintln(response, "</head>")
	fmt.Fprintln(response, "<body>")

	// Check the cookie from the client
	cookie, _ := request.Cookie("UserCookie")
	user := strings.TrimPrefix(cookie.String(), "UserCookie=")

	if users[user] != "" {
		fmt.Fprintln(response, "The time is now <span class=\"time\">"+
			time.Now().Format(timeFormat)+"</span>, "+users[user]+".</p>")
	} else {
		fmt.Fprintln(response, "The time is now <span class=\"time\">"+
			time.Now().Format(timeFormat)+"</span>.</p>")
	}
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

// General handler for the home page
func generalHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("In general handler")
	fmt.Println(request.URL.Path)
	// Check the path
	if request.URL.Path != "/" && request.URL.Path != "/index.html" {
		fmt.Println("Invalid path general handler")
		errorHandler(response, request)
		return
	}

	// Check the cookie from the client
	cookie, _ := request.Cookie("UserCookie")
	fmt.Println(cookie.String())
	user := strings.TrimPrefix(cookie.String(), "UserCookie=")

	// If user is logged in, print the greeting
	if users[user] != "" {
		fmt.Fprintln(response, "<html>")
		fmt.Fprintln(response, "<body>")
		fmt.Fprintln(response, "Greetings, "+users[user]+".")
		fmt.Fprintln(response, "</body>")
		fmt.Fprintln(response, "</html>")
	} else {
		// Print the login
		fmt.Fprintln(response, "<html>")
		fmt.Fprintln(response, "<body>")
		fmt.Fprintln(response, "<form action=\"login\">")
		fmt.Fprintln(response, " What is your name, Earthling?")
		fmt.Fprintln(response, " <input type=\"text\" name=\"name\" size=\"50\">")
		fmt.Fprintln(response, " <input type=\"submit\">")
		fmt.Fprintln(response, "</form>")
		fmt.Fprintln(response, "</p>")
		fmt.Fprintln(response, "</body>")
		fmt.Fprintln(response, "</html>")
	}
}

// Login handler for the form
func loginHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("In login handler")
	// Check the path
	if request.URL.Path != "/login" {
		fmt.Println("Invalid path login handler")
		errorHandler(response, request)
		return
	}

	// Read the name value from the form
	fmt.Println("Parsing form value")
	// request.ParseForm()
	userName := request.FormValue("name")
	bytesUUID, _ := exec.Command("uuidgen").Output()
	userUUID := fmt.Sprintf("%v", bytesUUID)

	// Synchronization, locking the users for writing
	fmt.Println("Locking")
	mutex.Lock()
	users[userUUID] = userName
	mutex.Unlock()
	fmt.Println("Lock released")

	// Set cookies on the client
	cookie := &http.Cookie{Name: "UserCookie", Value: userUUID, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(response, cookie)
	fmt.Println("Cookie Set")
	http.Redirect(response, request, "./..", 302) // redirect to home
	fmt.Println("Redirected to home")
}

// Logout handler for the form
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("In logout handler")
	// Check the path
	if request.URL.Path != "/logout" {
		errorHandler(response, request)
		return
	}

	// Check the cookie for error
	_, err := request.Cookie("UserCookie")

	// If no error, cleared the cookie
	if err != nil {
		fmt.Printf("Cookie error: %s!", err)
	} else {
		http.SetCookie(response, &(http.Cookie{Name: "UserCookie", Expires: time.Now()}))

		// Print the goodbye
		fmt.Fprintln(response, "<html>")
		fmt.Fprintln(response, "<head>")
		fmt.Fprintln(response, "<META http-equiv=\"refresh\" content=\"10;URL=/\">")
		fmt.Fprintln(response, "<p>Good-bye.</p>")
		fmt.Fprintln(response, "</body>")
		fmt.Fprintln(response, "</html>")
	}
}

// Display the 404 error
func errorHandler(response http.ResponseWriter, request *http.Request) {
	// Return status
	response.WriteHeader(http.StatusNotFound)

	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "<p>These are not the URLs you're looking for.</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

// Start the server
func main() {
	// Declare flags
	portPtr := flag.Int("port", 8080, "port_number")
	VPtr := flag.Bool("V", false, "display_version")

	// Parse the command line arguments
	flag.Parse()

	// Check -V flag
	if *VPtr == true {
		fmt.Printf(version)
		return
	}

	http.HandleFunc("/time", timeHandler)          // handle /time path
	http.HandleFunc("/time/", timeHandler)         // handle /time path
	http.HandleFunc("/", generalHandler)           // general path
	http.HandleFunc("/index.html", generalHandler) // general path
	http.HandleFunc("/login", loginHandler)        // login path
	http.HandleFunc("/logout", logoutHandler)      // logout path
	err := http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil)

	// Server error, display message and exit the application
	if err != nil {
		fmt.Printf("Server fail: %s\n", err)
		os.Exit(1)
	}
}
