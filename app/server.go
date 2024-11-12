package main

import (
	"bufio"
	"fmt"
	"io"
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
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connected to client:", conn.RemoteAddr())

	// Read incoming data
	reader := bufio.NewReader(conn)

	request := []string{}
	var response string

	for {
		chunk, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("reached EOF, exiting loop")
				break
			}
			fmt.Println("Error reading from connection:", err)
			return
		}
		if chunk == "\n" {
			break
		}
		chunk = strings.TrimSpace(chunk)
		fmt.Println("Received chunk:", chunk)

		request = append(request, chunk)

	}
	if isPing(request) {
		response = "+PONG\r\n"
	}
	if ok, resp := isEcho(request); ok {
		response = fmt.Sprintf("%s\r\n", resp)
	}
	// Send the response back to the client
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error sending response:", err)
		return
	}

}

func isPing(request []string) bool {
	for _, r := range request {
		if r == "PING" {
			return true
		}
	}
	return false
}

func isEcho(request []string) (bool, string) {
	for i, r := range request {
		if r == "ECHO" {
			return true, request[i+2]
		}
	}
	return false, ""
}
