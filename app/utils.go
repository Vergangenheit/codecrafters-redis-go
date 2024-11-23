package app

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func expired(res *Resource, currentTime time.Time) bool {
	if res.Expired == nil {
		return false
	}
	expiredTs := *res.Expired
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

func buildRespArray(req *Request) string {
	data := []string{string(req.Command)}

	data = append(data, req.Args...)

	var builder strings.Builder

	// Start with the number of keys
	builder.WriteString(fmt.Sprintf("*%d\r\n", len(data)))

	// Loop over the map to add each key
	for _, key := range data {
		keyLen := len(key)
		builder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", keyLen, key))
	}

	return builder.String()
}

func simpleRespString(data []string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("+%s", data[0]))
	if len(data) > 1 {
		for _, key := range data[1:] {
			builder.WriteString(fmt.Sprintf(" %s", key))
		}
	}
	builder.WriteString("\r\n")

	return builder.String()
}

func rdbContentResp(rdbBytes []byte) string {
	rdbString := string(rdbBytes)
	return fmt.Sprintf("$%d\r\n%s", len(rdbString), rdbString)
}
