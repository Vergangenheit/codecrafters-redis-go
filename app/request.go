package app

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	receiveBuf = 1024
)

type Command string

const (
	PING     Command = "PING"
	ECHO     Command = "ECHO"
	SET      Command = "SET"
	GET      Command = "GET"
	CONFIG   Command = "CONFIG"
	KEYS     Command = "KEYS"
	INFO     Command = "INFO"
	REPLCONF Command = "REPLCONF"
	PSYNC    Command = "PSYNC"
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
	case "INFO":
		return INFO, nil
	case "REPLCONF":
		return REPLCONF, nil
	case "PSYNC":
		return PSYNC, nil
	default:
		return "", fmt.Errorf("Command %s not recognized", str)
	}
}

func (c *Command) ToStr() string {
	switch *c {
	case PING:
		return "PING"
	case ECHO:
		return "ECHO"
	case SET:
		return "SET"
	case GET:
		return "GET"
	case CONFIG:
		return "CONFIG"
	case KEYS:
		return "KEYS"
	case INFO:
		return "INFO"
	case REPLCONF:
		return "REPLCONF"
	case PSYNC:
		return "PSYNC"
	default:
		return ""
	}
}

type Request struct {
	nBytes  int
	Command Command
	Args    []string
}

func (s *server) requestParser(conn net.Conn) (*Request, error) {
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

func sendRequestToServer(conn net.Conn, req *Request) error {

	switch req.Command {
	case PING:
		// send as a single element resp array
		message := "*1\r\n$4\r\nPING\r\n"
		_, err := conn.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("Failed to send data: %v", err)
		}
		buffer := make([]byte, receiveBuf)
		// start reading chunks delimited by newline byte
		_, err = conn.Read(buffer)
		if err != nil {
			return err
		}
		return nil
	case SET:
		// set shoould not wait for reply
		message := buildRespArray(req)
		_, err := conn.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("Failed to send SET req: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	case REPLCONF:
		switch req.Args[0] {
		case "listening-port":
			respArray := buildRespArray(req)
			_, err := conn.Write([]byte(respArray))
			if err != nil {
				return fmt.Errorf("Failed to send REPLCONF req: %v", err)
			}
			buffer := make([]byte, receiveBuf)
			// start reading chunks delimited by newline byte
			_, err = conn.Read(buffer)
			if err != nil {
				return err
			}
			return nil
		case "capa":
			respArray := buildRespArray(req)
			_, err := conn.Write([]byte(respArray))
			if err != nil {
				return fmt.Errorf("Failed to send REPLCONF req: %v", err)
			}
			buffer := make([]byte, receiveBuf)
			// start reading chunks delimited by newline byte
			_, err = conn.Read(buffer)
			if err != nil {
				return err
			}
			return nil
		default:
			return fmt.Errorf("REPLCONG args not recognized")
		}
	case PSYNC:
		respArray := buildRespArray(req)
		_, err := conn.Write([]byte(respArray))
		if err != nil {
			return fmt.Errorf("Failed to send PSYNC req: %v", err)
		}
		buffer := make([]byte, receiveBuf)
		// start reading chunks delimited by newline byte
		_, err = conn.Read(buffer)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("Cannot send %v requests", req.Command)
	}
}
