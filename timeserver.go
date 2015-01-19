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
func success(response http.ResponseWriter, request *http.Request) {
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

// Display the 404 error
func serverError(response http.ResponseWriter, request *http.Request) {
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

	http.HandleFunc("/time", success) // handle /time path
	http.HandleFunc("/", serverError) // other path
	err := http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil)

	// Server error, display message and exit the application
	if err != nil {
		fmt.Printf("Server fail: %s\n", err)
		os.Exit(1)
	}
}
