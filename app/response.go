package app

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func (s *server) parseResponse(req *Request) (string, error) {
	if req == nil {
		return "", fmt.Errorf("Request is nil")
	}
	switch req.Command {
	case PING:
		return "+PONG\r\n", nil
	case ECHO:
		return fmt.Sprintf("+%s\r\n", req.Args[0]), nil
	case SET:
		err := s.setValue(req.Args)
		if err != nil {
			return "", err
		}
		return "+OK\r\n", nil
	case GET:
		value, ok := s.getValue(req.Args[0])
		if ok {
			return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value), nil
		}
		return "$-1\r\n", nil
	case CONFIG:
		res, err := s.handleConfig(req.Args)
		if err != nil {
			return "", fmt.Errorf("error handling COMMAND %v", err)
		}
		return res, nil
	case KEYS:
		res, err := s.handleKeys(req.Args)
		if err != nil {
			return "", fmt.Errorf("error handling KEYS %v", err)
		}
		return res, nil
	case INFO:
		res, err := s.handleInfo(req.Args)
		if err != nil {
			return "", fmt.Errorf("error handling INFO %v", err)
		}
		return res, nil
	case REPLCONF:
		res, err := s.handleReplConf()
		if err != nil {
			return "", fmt.Errorf("error handling REPLCONF %v", err)
		}
		return res, nil
	case PSYNC:
		res, err := s.handlePsync(req.Args)
		if err != nil {
			return "", fmt.Errorf("error handling PSYNC %v", err)
		}
		return res, nil
	default:
		return "", fmt.Errorf("unknown request command %s", req.Command)
	}
}

func (s *server) setValue(args []string) error {
	// if expire time in args
	for i, arg := range args {
		args[i] = strings.ToLower(arg)
	}
	if slices.Contains(args, "px") {
		expiry, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("cannot convert expiry milliseconds to integer")
		}
		duration := time.Duration(expiry) * time.Millisecond
		expiredTs := time.Now().Add(duration)
		s.InMemoryStore[args[0]] = &Resource{
			value:   args[1],
			expired: &expiredTs,
		}
		return nil
	}
	s.InMemoryStore[args[0]] = &Resource{
		value: args[1],
	}
	return nil

}

func (s *server) getValue(key string) (string, bool) {
	tNow := time.Now()
	res, ok := s.InMemoryStore[key]
	if ok {
		valStr := res.value.(string)
		// does it have expiry?
		if expired(res, tNow) {
			return "", false
		}
		return valStr, ok
	}
	return "", false
}

func (s *server) handleConfig(args []string) (string, error) {
	// config first arg
	switch args[0] {
	case "GET":
		return s.handleConfigGet(args)
	default:
		return "", fmt.Errorf("unrecognized config command")
	}
}

func (s *server) handleConfigGet(args []string) (string, error) {
	switch args[1] {
	case "dir":
		// build RESP bulk string
		bulkStr := fmt.Sprintf("*2\r\n$%d\r\ndir\r\n$%d\r\n%s\r\n", 3, len(s.Config.Dir), s.Config.Dir)
		return bulkStr, nil
	case "dbfilename":
		bulkStr := fmt.Sprintf("*2\r\n$%d\r\ndbfilename\r\n$%d\r\n%s\r\n", 3, len(s.Config.DbFilename), s.Config.DbFilename)
		return bulkStr, nil
	default:
		return "", fmt.Errorf("unrecognized config command, expecting dir or dbfilename")
	}
}

func (s *server) handleKeys(args []string) (string, error) {
	switch args[0] {
	case "*":
		// return all the keys
		// return all the keys
		return formatMapKeys(s.InMemoryStore), nil

	default:
		return "", fmt.Errorf("argument for KEYS is not supported")
	}
}

func (s *server) handleInfo(args []string) (string, error) {
	switch args[0] {
	case "replication":
		keyVal1 := "role:master"
		if s.Config.ReplicaOf != nil {
			keyVal1 = "role:slave"
		}
		keyVal2 := "master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		keyVal3 := "master_repl_offset:0"
		bulkString := formatBulkString([]string{keyVal1, keyVal2, keyVal3})
		fmt.Printf("got bulk string %s", bulkString)
		return bulkString, nil
	default:
		return "", fmt.Errorf("Unrecognized argument %s for INFO command", args[0])
	}
}

func (s *server) handleReplConf() (string, error) {
	// always repond with a simple RESP simple string OK
	return "+OK\r\n", nil
}

func (s *server) handlePsync(args []string) (string, error) {
	return simpleRespString([]string{
		"FULLRESYNC", "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb", "0",
	}), nil
}
