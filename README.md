Instructions for running the timeserver program

Install GO and set GOPATH
Navigate to this directory (ie TimeServer) in your command line
Type go run timeserver.go (optional -V and -port) and enter
Open a web browser and type localhost:8080/time in the address bar,
substitue 8080 with your port number if specified the go run command.
To terminate the server, hit Ctrl + C from the command line.

timeserver2 program
Similar to timeserver with personalized message for the user.
Support the following URL:
http://host:port/ or http://host:port/index.html
http://host:port/logout
http://host:port/time