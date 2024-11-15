package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Define the server address and port (example: localhost:6379)
	serverAddress := "localhost:6379"

	// Connect to the TCP server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer conn.Close()

	// The slice of bytes to send
	// message := []byte("*3\r\n$3\r\nSET\r\n$6\r\nbanana\r\n$10\r\nstrawberry\r\n")
	// message := []byte("*2\r\n$3\r\nGET\r\n$5\r\nbanana\r\n")
	// message := []byte("*2\r\n$4\r\nECHO\r\n$5\r\napple\r\n")
	// message := []byte("*1\r\n$4\r\nPING\r\n")
	// message := []byte("*5\r\n$3\r\nSET\r\n$6\r\norange\r\n$5\r\ngrape\r\n$2\r\npx\r\n$3\r\n180000\r\n")
	message := []byte("*2\r\n$3\r\nGET\r\n$5\r\norange\r\n")
	// Send the message to the server
	_, err = conn.Write(message)
	if err != nil {
		log.Fatal("Failed to send data:", err)
	}
	fmt.Println("Message sent to server:", string(message))

	// Set a 5-second timeout for reading the response
	// conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	// Read the response from the server
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Fatal("Failed to read conn:", err)
	}

	// Print the server response
	fmt.Printf("Response from server: %s\n", string(buffer))
}
