package app

import "time"

type InMemoryStore map[string]*Resource

type Resource struct {
	value   interface{}
	expired *time.Time
}
