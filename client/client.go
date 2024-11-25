package client

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Vergangenheit/codecrafters-redis-go/app"
)

type RedisClient struct {
	conn net.Conn
}

func NewRedisClient(serverAddress string) (*RedisClient, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return &RedisClient{conn: conn}, nil
}

func (r *RedisClient) Send(request *app.Request) ([]string, error) {
	message, err := r.requestSerializer(request)
	if err != nil {
		return nil, err
	}
	_, err = r.conn.Write([]byte(message))
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	_, err = r.conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	resp := r.deserializeResponse(buffer)

	return resp, nil
}

func (r *RedisClient) deserializeResponse(buffer []byte) []string {
	// Convert byte slice to string
	str := string(buffer)

	// Trim any trailing \r\n to prevent an extra empty element after splitting
	str = strings.TrimSpace(str)
	str = strings.TrimSuffix(str, "\r\n")

	// Split the string by "\r\n"
	parts := strings.Split(str, "\r\n")

	cleanedParts := []string{}
	// clean +
	for _, part := range parts {
		if strings.HasPrefix(part, "+") {
			part = strings.TrimPrefix(part, "+")
		}
		if strings.HasPrefix(part, "$") || strings.HasPrefix(part, "*") {
			continue
		}
		cleanedParts = append(cleanedParts, part)
	}

	return cleanedParts[:len(cleanedParts)-1]
}

func (r *RedisClient) Close() {
	r.conn.Close()
}

func (r *RedisClient) requestSerializer(request *app.Request) (string, error) {
	switch request.Command {
	case app.PING:
		return "*1\r\n$4\r\nPING\r\n", nil
	case app.ECHO, app.SET, app.GET, app.CONFIG, app.KEYS, app.INFO, app.REPLCONF, app.PSYNC:
		bulkStr := r.buildBulkString(request)
		return bulkStr, nil
	default:
		return "", fmt.Errorf("Command not recognized")
	}
}

func (r *RedisClient) buildBulkString(request *app.Request) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("*%d\r\n", len(request.Args)+1))
	builder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(request.Command), request.Command))

	for _, arg := range request.Args {
		builder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}

	return builder.String()
}

func main() {
	// Define the server address and port (example: localhost:6379)
	serverAddress := "localhost:6789"

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
	// message := []byte("*2\r\n$3\r\nGET\r\n$5\r\norange\r\n")
	// message := []byte("*2\r\n$4\r\nKEYS\r\n$1\r\n*\r\n")
	message := []byte("*2\r\n$4\r\nINFO\r\n$11\r\nreplication\r\n")
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
