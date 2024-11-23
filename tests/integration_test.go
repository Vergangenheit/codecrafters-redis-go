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
			t.Errorf("Server run failed: %v", err)
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
