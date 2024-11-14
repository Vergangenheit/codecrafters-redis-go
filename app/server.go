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

type server struct {
	Listener      net.Listener
	InMemoryStore InMemoryStore
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	store := InMemoryStore{}
	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	server := &server{
		Listener:      l,
		InMemoryStore: store,
	}
	defer l.Close()

	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			fmt.Println("cannot accept a connection")
		}
		// Handle the connection in a new goroutine
		go server.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn net.Conn) {
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
		response, err = s.parseResponse(request)
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
