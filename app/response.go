package main

import "fmt"

func parseResponse(req *Request) (string, error) {
	if req == nil {
		return "", fmt.Errorf("Request is nil")
	}
	switch req.Command {
	case PING:
		return "+PONG\r\n", nil
	case ECHO:
		return fmt.Sprintf("+%s\r\n", req.Args[0]), nil
	case SET:
		return "+OK\r\n", nil
	default:
		return "", fmt.Errorf("unknown request command %s", req.Command)
	}
}
