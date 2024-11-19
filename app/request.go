package app

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
	PING   Command = "PING"
	ECHO   Command = "ECHO"
	SET    Command = "SET"
	GET    Command = "GET"
	CONFIG Command = "CONFIG"
	KEYS   Command = "KEYS"
)

func toCommand(str string) (Command, error) {
	switch str {
	case "PING":
		return PING, nil
	case "ECHO":
		return ECHO, nil
	case "SET":
		return SET, nil
	case "GET":
		return GET, nil
	case "CONFIG":
		return CONFIG, nil
	case "KEYS":
		return KEYS, nil
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
	// parse command args
	req.Args = extractArgs(parsedArray)

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

func extractArgs(parsedArray []string) []string {
	args := []string{}
	// args start after position two in array
	for _, chunk := range parsedArray[3 : len(parsedArray)-1] {
		// check if first character is $
		if chunk[0] == '$' {
			continue
		}
		args = append(args, chunk)
	}
	return args
}
