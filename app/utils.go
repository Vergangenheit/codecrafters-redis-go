package app

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func expired(res *Resource, currentTime time.Time) bool {
	if res.expired == nil {
		return false
	}
	expiredTs := *res.expired
	if expiredTs.Before(currentTime) {
		return true
	}
	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	// Ensure it's not some other error
	return err == nil && !info.IsDir()
}

func formatMapKeys(m map[string]*Resource) string {
	var builder strings.Builder

	// Start with the number of keys
	builder.WriteString(fmt.Sprintf("*%d\r\n", len(m)))

	// Loop over the map to add each key
	for key := range m {
		keyLen := len(key)
		builder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", keyLen, key))
	}

	return builder.String()
}

func formatBulkString(data []string) string {
	// Join the slice into a single string separated by newlines
	joined := strings.Join(data, "\n")

	// Calculate the byte length of the joined string
	length := len(joined)

	var builder strings.Builder

	// Start with the number of keys
	builder.WriteString(fmt.Sprintf("$%d\r\n", length))
	builder.WriteString(joined)
	builder.WriteString("\r\n")

	return builder.String()
}
