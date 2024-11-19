package app

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Vergangenheit/rdb-go"
)

func ReadRedisDBFile(filename string) (map[string]*Resource, error) {
	// Open the Redis RDB file
	rdbFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open RDB file: %v", err)
	}
	defer rdbFile.Close()

	// Parse the RDB file and extract key-value pairs
	result := make(map[string]*Resource)

	parser := rdb.NewParser(rdbFile)

	for {
		data, err := parser.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		switch data := data.(type) {
		case *rdb.StringData:
			// add it to store
			result[data.Key] = &Resource{
				value:   data.Value,
				expired: data.Expiry,
			}
		}
	}

	return result, nil
}
