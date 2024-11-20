package app

type Config struct {
	Dir        string
	DbFilename string
	Port       string
	ReplicaOf  *string
}
