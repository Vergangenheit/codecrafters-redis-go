package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to the TCP server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server.")

	// Read input from the console and send to the server
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter request (e.g., HELLO or TIME): ")
		request, _ := reader.ReadString('\n')
		request = strings.TrimSpace(request)

		// Send the request to the server
		_, err := conn.Write([]byte(request + "\n"))
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		// Read the response from the server
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Printf("Server response: %s", response)
	}
}
