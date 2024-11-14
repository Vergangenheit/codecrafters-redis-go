package main

import "fmt"

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
		s.setValue(req.Args[0], req.Args[1])
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

func (s *server) setValue(key, value string) {
	s.InMemoryStore[key] = value
}

func (s *server) getValue(key string) (string, bool) {
	val, ok := s.InMemoryStore[key]
	if ok {
		valStr := val.(string)
		return valStr, ok
	}
	return "", false
}
