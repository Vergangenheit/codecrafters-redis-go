package main

import "time"

type InMemoryStore map[string]*Resource

type Resource struct {
	value   string
	expired *time.Time
}
