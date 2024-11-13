package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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
	return false, ""
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connected to client:", conn.RemoteAddr())

	// Read incoming data
	reader := bufio.NewReader(conn)

	for {
		request, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Println("read request string:", request)
		request = strings.TrimSpace(request)
		fmt.Println("Received request:", request)

		// Process the request and create a response
		var response string
		switch request {
		case "PING":
			response = "+PONG\r\n"
		default:
			continue
		}
		// Send the response back to the client
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}

}
