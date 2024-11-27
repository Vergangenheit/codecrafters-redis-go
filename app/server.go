package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	hclog "github.com/hashicorp/go-hclog"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type server struct {
	Listener      net.Listener
	InMemoryStore InMemoryStore
	Config        *Config
	Logger        hclog.Logger
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewServer(contextBack context.Context, config *Config, logger hclog.Logger) (*server, error) {
	// if dbfilename is valid, check if it should be parsed into inmemory store
	ctx, cancel := context.WithCancel(contextBack)
	store := InMemoryStore{}
	if config.DbFilename != "" && config.Dir != "" {
		// check if path exists
		fullPath := filepath.Join(config.Dir, config.DbFilename)
		if fileExists(fullPath) {
			inMemoryStore, err := ReadRedisDBFile(fullPath)
			if err != nil {
				cancel()
				return nil, fmt.Errorf("cannot parse dump file %v", err)
			}
			store = inMemoryStore
		}
	}
	return &server{
		InMemoryStore: store,
		Config:        config,
		Logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

func (s *server) RunServer() error {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	s.Logger.Info("Logs from your program will appear here!")
	// read config

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", s.Config.Port))
	if err != nil {
		s.Logger.Error("Failed to bind to port %s", s.Config.Port)
		return err
	}
	s.Listener = l
	defer l.Close()
	if s.Config.ReplicaOf != nil {
		fmt.Printf("server is replica of %s", *s.Config.ReplicaOf)
		err := s.handhshakeWithMaster()
		if err != nil {
			s.Logger.Error("Failed to handshake with master %v", err)
			return err
		}
	}

	for {
		select {
		case <-s.ctx.Done():
			s.Logger.Info("Server shutting down")
			return nil
		default:
			conn, err := s.Listener.Accept()
			if err != nil {
				s.Logger.Error("cannot accept a connection")
				return err
			}
			// Handle the connection in a new goroutine
			go s.handleConnection(conn)
		}
	}
}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	s.Logger.Info("Connected to client:", conn.RemoteAddr())

	for {
		responses, err := s.generateResponses(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			s.Logger.Error("Error parsing response:", err)
			return
		}
		s.Logger.Info("Responses generated:", responses)
		for _, response := range responses {
			// Send the response back to the client
			_, err = conn.Write([]byte(response))
			if err != nil {
				s.Logger.Error("Error sending response:", err)
				return
			}
			s.Logger.Info(fmt.Sprintf("Response sent: %s", response))
		}
	}

}

func (s *server) handhshakeWithMaster() error {
	if s.Config.ReplicaOf == nil {
		return errors.New("master server address cannot be null")
	}
	// compose address since replica of is separated
	serverAddr := strings.Join(strings.Split(*s.Config.ReplicaOf, " "), ":")
	// Connect to the TCP server and keep the connection open
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		s.Logger.Error("Connection failed: %v", err)
		return err
	}
	// start sending PING
	err = sendRequestToServer(conn, &Request{
		Command: PING,
	})
	if err != nil {
		s.Logger.Error("Failed to ping master %v", err)
	}
	// send first REPLCONF
	err = sendRequestToServer(conn, &Request{
		Command: REPLCONF,
		Args:    []string{"listening-port", s.Config.Port},
	})
	if err != nil {
		s.Logger.Error("Failed to send first replconf to master %v", err)
	}
	// send second REPLCONF
	err = sendRequestToServer(conn, &Request{
		Command: REPLCONF,
		Args:    []string{"capa", "psync2"},
	})
	if err != nil {
		s.Logger.Error("Failed to send second replconf to master %v", err)
	}
	// send PSYNC request
	err = sendRequestToServer(conn, &Request{
		Command: PSYNC,
		Args:    []string{"?", "-1"},
	})
	if err != nil {
		s.Logger.Error("Failed to send psync to master %v", err)
	}
	return nil
}

func (s *server) propagateToReplicas(request *Request) error {
	if len(s.Config.replicas) == 0 {
		return nil
	}
	for _, replica := range s.Config.replicas {
		// establish connection to replica
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", replica.Host, replica.Port))
		if err != nil {
			return fmt.Errorf("Failed to connect to replica %v", err)
		}
		err = sendRequestToServer(conn, request)
		if err != nil {
			return fmt.Errorf("error propagating to replica %v", err)
		}
		err = conn.Close()
		if err != nil {
			return fmt.Errorf("Failed to close connection to replica %v", err)
		}
	}
	return nil
}

func (s *server) Stop() {
	err := s.Listener.Close()
	if err != nil {
		s.Logger.Error("Failed to close listener %v", err)
	}
	s.cancel() // Signal the server to stop
}
