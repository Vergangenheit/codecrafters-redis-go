package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("cannot accept a connection")
		}
		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connected to client:", conn.RemoteAddr())

	var response string
	// parse request
	request, err := RequestParser(conn)
	if err != nil {
		fmt.Printf("Cannot parse the request %v", err)
	}
	if isPing(request) {
		response = "+PONG\r\n"
	}
	if ok, resp := isEcho(request); ok {
		response = fmt.Sprintf("+%s\r\n", resp)
	}
	// Send the response back to the client
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error sending response:", err)
		return
	}

}

func isPing(request *Request) bool {
	return request.Command == PING
}

func isEcho(request *Request) (bool, string) {
	if request.Command == ECHO {
		return true, request.Args[1]
	}
	return false, ""
}
