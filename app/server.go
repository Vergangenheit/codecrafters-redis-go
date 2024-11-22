package app

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type server struct {
	Listener      net.Listener
	InMemoryStore InMemoryStore
	Config        *Config
}

func NewServer(listener net.Listener, store InMemoryStore, config *Config) (*server, error) {
	// if dbfilename is valid, check if it should be parsed into inmemory store
	if config.DbFilename != "" && config.Dir != "" {
		// check if path exists
		fullPath := filepath.Join(config.Dir, config.DbFilename)
		if fileExists(fullPath) {
			inMemoryStore, err := ReadRedisDBFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("cannot parse dump file %v", err)
			}
			store = inMemoryStore
		}
	}
	return &server{
		Listener:      listener,
		InMemoryStore: store,
		Config:        config,
	}, nil
}

func RunServer(config *Config) error {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	store := InMemoryStore{}
	// read config

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.Port))
	if err != nil {
		return fmt.Errorf("Failed to bind to port %d %v", config.Port, err)
	}
	server, err := NewServer(l, store, config)
	if config.ReplicaOf != nil {
		fmt.Printf("server is replica of %s", *config.ReplicaOf)
		err := server.handhshakeWithMaster()
		if err != nil {
			return fmt.Errorf("Failed to handshake with master %v", err)
		}
	}
	if err != nil {
		return fmt.Errorf("Failed to instantiate server %v", err)
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

func (s *server) handhshakeWithMaster() error {
	if s.Config.ReplicaOf == nil {
		return errors.New("master server address cannot be null")
	}
	// compose address since replica of is separated
	serverAddr := strings.Join(strings.Split(*s.Config.ReplicaOf, " "), ":")
	// Connect to the TCP server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("Connection failed: %v", err)
	}
	defer conn.Close()
	// start sending PING
	err = sendRequestToMaster(conn, &Request{
		Command: PING,
	})
	if err != nil {
		return fmt.Errorf("Failed to ping master %v", err)
	}
	// send first REPLCONF
	err = sendRequestToMaster(conn, &Request{
		Command: REPLCONF,
		Args:    []string{"listening-port", s.Config.Port},
	})
	if err != nil {
		return fmt.Errorf("Failed to send first replconf to master %v", err)
	}
	// send second REPLCONF
	err = sendRequestToMaster(conn, &Request{
		Command: REPLCONF,
		Args:    []string{"capa", "psync2"},
	})
	if err != nil {
		return fmt.Errorf("Failed to send second replconf to master %v", err)
	}
	// send PSYNC request
	err = sendRequestToMaster(conn, &Request{
		Command: PSYNC,
		Args:    []string{"?", "-1"},
	})
	if err != nil {
		return fmt.Errorf("Failed to send psync to master %v", err)
	}
	return nil
}
