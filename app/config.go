package app

import "net"

type Config struct {
	Dir        string
	DbFilename string
	Port       string
	ReplicaOf  *string
	replicas   []*Replica
}

type Replica struct {
	Port string
	conn net.Conn
}
