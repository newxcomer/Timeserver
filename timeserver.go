// Simple program to run a server showing a web page
// displaying the current time of day
// Example usage: http://localhost:8080/time
// Default port number is 8080
// Optional flag -V printing version number to console output
// Optional flag -port for port number
//
//	Author: Binh Nguyen
//	Date: 1/13/2015
//	Version: 0.0010

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	version = "Timeserver version 0.0010"
)

// Display the current time
func timeHandler(response http.ResponseWriter, request *http.Request) {
	// Current time format
	const timeFormat = "3:04:052 PM"

	// Print the response
	fmt.Fprintln(response, "<html>")
	fmt.Fprintln(response, "<head>")
	fmt.Fprintln(response, "<style>")
	fmt.Fprintln(response, "p {font-size: xx-large}")
	fmt.Fprintln(response, "span.time {color: red}")
	fmt.Fprintln(response, "</style>")
	fmt.Fprintln(response, "</head>")
	fmt.Fprintln(response, "<body>")
	fmt.Fprintln(response, "The time is now <span class=\"time\">"+
		time.Now().Format(timeFormat)+"</span>.</p>")
	fmt.Fprintln(response, "</body>")
	fmt.Fprintln(response, "</html>")
}

// General handler for the home page
func generalHandler(response http.ResponseWriter, request *http.Request) {
	// Check the cookie from the client
	cookie, _ := request.Cookie("UserCookie")

	// If user is logged in, print the greeting
	if cookie.Value != "" {
		fmt.Fprintln(response, "Greetings, "+cookie.Value+".")
	} else { // Redirect to log in
		http.Redirect(response, request, "/login?name=name", http.StatusFound)
	}
}

// Login handler for the form
func loginHandler(response http.ResponseWriter, request *http.Request) {
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

	// Read the name value from the form
	request.ParseForm()
	userName := request.FormValue("name")

	// Set cookies on the client
	cookie := &http.Cookie{Name: "UserCookie", Value: userName, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(response, cookie)
}

// Logout handler for the form
func logoutHandler(response http.ResponseWriter, request *http.Request) {
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

	http.HandleFunc("/time", timeHandler)         // handle /time path
	http.HandleFunc("/", generalHandler)          // general path
	http.HandleFunc("/login?name=", loginHandler) // login path
	http.HandleFunc("/logout", logoutHandler)     // logout path
	err := http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil)

	// Server error, display message and exit the application
	if err != nil {
		fmt.Printf("Server fail: %s\n", err)
		os.Exit(1)
	}
}
