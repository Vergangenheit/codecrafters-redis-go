package app

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Ipv6(t *testing.T) {
	remoteAddr := "[::1]:1234"
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	version := net.ParseIP(ip)
	assert.Nil(t, version.To4())
}

func Test_Ipv4(t *testing.T) {
	remoteAddr := "127.0.0.1:1234"
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, "127.0.0.1", ip)
	version := net.ParseIP(ip)
	assert.NotNil(t, version.To4())
}
