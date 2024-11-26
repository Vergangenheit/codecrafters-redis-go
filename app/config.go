package app

type Config struct {
	Dir        string
	DbFilename string
	// TODO - should add host alongside port
	Port      string
	ReplicaOf *string
	replicas  []*Replica
}

type Replica struct {
	// TODO - should add host alongside port
	Port string
}
