package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	receiveBuf = 1024
)

type Command string

const (
	PING Command = "PING"
	ECHO Command = "ECHO"
)

func toCommand(str string) (Command, error) {
	switch str {
	case "PING":
		return PING, nil
	case "ECHO":
		return ECHO, nil
	default:
		return "", fmt.Errorf("Command %s not recognized", str)
	}
}

type Request struct {
	nBytes  int
	Command Command
	Args    []string
}

func RequestParser(conn net.Conn) (*Request, error) {
	fmt.Println("parsing request")
	// buffer the conn
	buffer := make([]byte, receiveBuf)

	req := &Request{}
	// start reading chunks delimited by newline byte

	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	req.nBytes = n
	parsedArray := parseBulkBytes(buffer)
	comm, err := toCommand(parsedArray[2])
	if err != nil {
		return nil, fmt.Errorf("command not recognized %s", parsedArray[2])
	}
	req.Command = comm
	if len(parsedArray) > 5 {
		req.Args = parsedArray[4:]
	}

	return req, nil
}

func parseBulkBytes(input []byte) []string {
	// Convert byte slice to string
	str := string(input)

	// Trim any trailing \r\n to prevent an extra empty element after splitting
	str = strings.TrimSuffix(str, "\r\n")

	// Split the string by "\r\n"
	parts := strings.Split(str, "\r\n")

	return parts
}
