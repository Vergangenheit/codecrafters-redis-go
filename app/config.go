package app

type Config struct {
	Dir        string
	DbFilename string
	Port       string
	ReplicaOf  *string
	replicas   []*Replica
}

type Replica struct {
	Host string
	Port string
}
