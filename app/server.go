package main

import (
	"fmt"
	"io"
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

	for {
		var response string
		// parse request
		request, err := RequestParser(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Cannot parse the request %v", err)
		}
		fmt.Println("parsed request ", request)
		response, err = parseResponse(request)
		if err != nil {
			fmt.Println("Error parsing response:", err)
			return
		}
		// Send the response back to the client
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}

}

func isPing(request *Request) bool {
	return request.Command == PING
}

func isEcho(request *Request) (bool, string) {
	if request.Command == ECHO {
		return true, request.Args[0]
	}
	return false, ""
}
