package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
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
	// buffer the conn
	buffer := bufio.NewReader(conn)

	req := &Request{}
	// start reading chunks delimited by newline byte
	for i := 0; i < 5; i++ {
		chunk, err := buffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("reached EOF, exiting loop")
				break
			}
			fmt.Printf("Error reading from connection %v", err)
			return nil, err
		}
		req.nBytes++

		// I'm expecting command at the third byte postion excluding whitespaces and newline bytes
		chunk = strings.TrimSpace(chunk)
		if i == 2 {
			comm, err := toCommand(chunk)
			if err != nil {
				return nil, err
			}
			req.Command = comm
		}
		if i > 2 {
			// register args
			req.Args = append(req.Args, chunk)
		}
	}
	return req, nil
}
