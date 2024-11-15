package main

import (
	"fmt"
	"slices"
	"strconv"
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
	default:
		return "", fmt.Errorf("unknown request command %s", req.Command)
	}
}

func (s *server) setValue(args []string) error {
	// if expire time in args
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
		valStr := res.value
		// does it have expiry?
		if expired(res, tNow) {
			return "", false
		}
		return valStr, ok
	}
	return "", false
}
