package tests

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Vergangenheit/codecrafters-redis-go/app"
	"github.com/Vergangenheit/codecrafters-redis-go/client"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func Test_BindToPort(t *testing.T) {
	// Test BindToPort
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	cl.Close()

	assert.NoError(t, errC)

	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToPing(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.PING})
	if err != nil {
		t.Fatalf("Failed to send PING: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"PONG"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToMultiplePings(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PINGS
	for i := 0; i < 5; i++ {
		resp, err := cl.Send(&app.Request{Command: app.PING})
		if err != nil {
			t.Fatalf("Failed to send PING: %v", err)
		}

		assert.NoError(t, errC)
		assert.Equal(t, []string{"PONG"}, resp)
	}
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToConcurrentPing(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	const clientCount = 10
	var wg sync.WaitGroup
	responses := make([][]string, clientCount)
	errors := make([]error, clientCount)
	// Step 3: Simulate concurrent clients
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			cl, err := client.NewRedisClient("localhost:6788")
			if err != nil {
				errors[clientID] = fmt.Errorf("Client %d failed to connect: %v", clientID, err)
				return
			}
			defer cl.Close()

			resp, err := cl.Send(&app.Request{Command: app.PING})
			if err != nil {
				errors[clientID] = fmt.Errorf("Client %d failed to send PING: %v", clientID, err)
			} else {
				responses[clientID] = resp
			}
		}(i)
	}

	// Wait for all client goroutines to finish
	wg.Wait()

	// Step 4: Validate results
	for i, err := range errors {
		if err != nil {
			t.Errorf("Error in client %d: %v", i, err)
		}
	}
	for i, resp := range responses {
		assert.Equal(t, []string{"PONG"}, resp, "Client %d response not as expected", i)
	}
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToEcho(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.ECHO, Args: []string{"hello"}})
	if err != nil {
		t.Fatalf("Failed to send ECHO: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"hello"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToSet(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.SET, Args: []string{"foo", "bar"}})
	if err != nil {
		t.Fatalf("Failed to send SET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"OK"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToGetExistingKey(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}
	server.InMemoryStore["foo"] = &app.Resource{
		Value: "bar",
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.GET, Args: []string{"foo"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"bar"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToGetNonExistingKey(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.GET, Args: []string{"foo"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToGetExpiredKey(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}
	expired := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	server.InMemoryStore["foo"] = &app.Resource{
		Value:   "bar",
		Expired: &expired,
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.GET, Args: []string{"foo"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToGetNonExpiredKey(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}
	expired := time.Now().Add(1 * time.Hour)
	server.InMemoryStore["foo"] = &app.Resource{
		Value:   "bar",
		Expired: &expired,
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.GET, Args: []string{"foo"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"bar"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToConfigGetDir(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port:       "6788",
		Dir:        "/tmp/redis-files",
		DbFilename: "dump.rdb",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.CONFIG, Args: []string{"GET", "dir"}})
	if err != nil {
		t.Fatalf("Failed to send CONFIG: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"dir", "/tmp/redis-files"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToConfigGetDbfilename(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port:       "6788",
		Dir:        "/tmp/redis-files",
		DbFilename: "dump.rdb",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.CONFIG, Args: []string{"GET", "dbfilename"}})
	if err != nil {
		t.Fatalf("Failed to send CONFIG: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"dbfilename", "dump.rdb"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToKeys(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port:       "6788",
		Dir:        "/tmp/redis-files",
		DbFilename: "keys_with_expiry.rdb",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}
	server.InMemoryStore["foo"] = &app.Resource{
		Value: "bar",
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.KEYS, Args: []string{"*"}})
	if err != nil {
		t.Fatalf("Failed to send KEYS: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"foo"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToGetFromDumpFile(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port:       "6788",
		Dir:        "../tmp/redis-files",
		DbFilename: "dump.rdb",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.GET, Args: []string{"foo"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"bar"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_ReadKeysFromDumpFile(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port:       "6788",
		Dir:        "../tmp/redis-files",
		DbFilename: "dump.rdb",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.KEYS, Args: []string{"*"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Len(t, resp, 3)
	assert.Contains(t, resp, "foo")
	assert.Contains(t, resp, "bar")
	assert.Contains(t, resp, "exp")
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_GetMultipleValuesFromDumpFile(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port:       "6788",
		Dir:        "../tmp/redis-files",
		DbFilename: "dump.rdb",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send first GET
	resp1, err := cl.Send(&app.Request{Command: app.GET, Args: []string{"foo"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"bar"}, resp1)

	resp2, err := cl.Send(&app.Request{Command: app.GET, Args: []string{"bar"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"foo"}, resp2)

	resp3, err := cl.Send(&app.Request{Command: app.GET, Args: []string{"exp"}})
	if err != nil {
		t.Fatalf("Failed to send GET: %v", err)
	}
	assert.NoError(t, errC)
	assert.Equal(t, []string{}, resp3)

	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToInfoReplication(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6788",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PING
	resp, err := cl.Send(&app.Request{Command: app.INFO, Args: []string{"replication"}})
	if err != nil {
		t.Fatalf("Failed to send INFO: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"role:master\nmaster_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb\nmaster_repl_offset:0"}, resp)
	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToInfoReplicationSlave(t *testing.T) {
	configMaster := &app.Config{
		Port: "6380",
	}
	master, err := app.NewServer(context.Background(), configMaster, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate master server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := master.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	config := &app.Config{
		Port:      "6788",
		ReplicaOf: ToPtr("localhost:6380"),
	}
	slave, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	doneSlave := make(chan struct{})
	go func() {
		err := slave.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(doneSlave)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6788")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send INFO
	resp, err := cl.Send(&app.Request{Command: app.INFO, Args: []string{"replication"}})
	if err != nil {
		t.Fatalf("Failed to send INFO: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"role:slave\nmaster_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb\nmaster_repl_offset:0"}, resp)
	cl.Close()
	// Stop the slave
	slave.Stop() // This method needs to be implemented in your server code

	// Wait for slave to finish
	<-doneSlave

	master.Stop()
	// Wait for master to finish
	<-done
}

func Test_RespondToReplcConf(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6379",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send REPLCONF
	resp1, err := cl.Send(&app.Request{Command: app.REPLCONF, Args: []string{"listening-port", "6788"}})
	if err != nil {
		t.Fatalf("Failed to send REPLCONF: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"OK"}, resp1)

	resp2, err := cl.Send(&app.Request{Command: app.REPLCONF, Args: []string{"capa", "psync2"}})
	if err != nil {
		t.Fatalf("Failed to send REPLCONF: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"OK"}, resp2)

	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_RespondToPsyncEmptyRdbFile(t *testing.T) {
	// Test RespondToPing
	config := &app.Config{
		Port: "6379",
	}
	server, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := server.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send PSYNC
	resp, err := cl.Send(&app.Request{Command: app.PSYNC, Args: []string{"?", "-1"}})
	if err != nil {
		t.Fatalf("Failed to send PSYNC: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"FULLRESYNC 8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb 0"}, resp)

	cl.Close()
	// Stop the server
	server.Stop() // This method needs to be implemented in your server code

	// Wait for server to finish
	<-done
}

func Test_PropagateSet(t *testing.T) {
	configMaster := &app.Config{
		Port: "6381",
	}
	master, err := app.NewServer(context.Background(), configMaster, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate master server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := master.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	config := &app.Config{
		Port:      "6789",
		ReplicaOf: ToPtr("localhost:6381"),
	}
	slave, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	doneSlave := make(chan struct{})
	go func() {
		err := slave.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(doneSlave)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// Test client connection
	cl, errC := client.NewRedisClient("localhost:6381")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	// send INFO
	resp, err := cl.Send(&app.Request{Command: app.SET, Args: []string{"foo", "bar"}})
	if err != nil {
		t.Fatalf("Failed to send INFO: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"OK"}, resp)
	cl.Close()
	// Stop the slave
	slave.Stop() // This method needs to be implemented in your server code

	// Wait for slave to finish
	<-doneSlave

	master.Stop()
	// Wait for master to finish
	<-done
}

func Test_ReplicaRespondToGet(t *testing.T) {
	configMaster := &app.Config{
		Host: "[::1]",
		Port: "6382",
	}
	master, err := app.NewServer(context.Background(), configMaster, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate master server: %v", err)
	}

	// Run server in a separate goroutine
	done := make(chan struct{})
	go func() {
		err := master.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(done)
	}()

	config := &app.Config{
		Host:      "[::1]",
		Port:      "6790",
		ReplicaOf: ToPtr("[::1]:6382"),
	}
	slave, err := app.NewServer(context.Background(), config, hclog.NewNullLogger())
	if err != nil {
		t.Fatalf("Failed to instantiate server: %v", err)
	}

	// Run server in a separate goroutine
	doneSlave := make(chan struct{})
	go func() {
		err := slave.RunServer()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				t.Errorf("Server run failed: %v", err)
			}
		}
		close(doneSlave)
	}()

	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// client connection to master
	cl, errC := client.NewRedisClient("[::1]:6382")
	if err != nil {
		t.Fatalf("Failed to connect to master server: %v", err)
	}
	// send INFO
	resp, err := cl.Send(&app.Request{Command: app.SET, Args: []string{"foo", "bar"}})
	if err != nil {
		t.Fatalf("Failed to send SET: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"OK"}, resp)
	cl.Close()

	time.Sleep(100 * time.Millisecond)

	// send get to slave
	cl2, errC := client.NewRedisClient("[::1]:6790")
	if err != nil {
		t.Fatalf("Failed to connect to slave server: %v", err)
	}
	// send GET
	resp2, err := cl2.Send(&app.Request{Command: app.GET, Args: []string{"foo"}})
	if err != nil {
		t.Fatalf("Failed to send GET to slave: %v", err)
	}

	assert.NoError(t, errC)
	assert.Equal(t, []string{"bar"}, resp2)
	cl2.Close()
	// Stop the slave
	slave.Stop() // This method needs to be implemented in your server code

	// Wait for slave to finish
	<-doneSlave

	master.Stop()
	// Wait for master to finish
	<-done
}
