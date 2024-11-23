package app

import "time"

type InMemoryStore map[string]*Resource

type Resource struct {
	Value   interface{}
	Expired *time.Time
}
