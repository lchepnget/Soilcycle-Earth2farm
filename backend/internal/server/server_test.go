package server

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/soilcycle-earth2farm/backend/internal/config"
)

func TestRecoveryMiddlewareReturns500(t *testing.T) {
	cfg := &config.Config{
		CORSAllowedOrigins: []string{"http://localhost:3000"},
	}

	app := New(nil, cfg)

	app.Get("/__test/recovery", func(c fiber.Ctx) error {
		panic("intentional test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/__test/recovery", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestListenGracefulShutdown(t *testing.T) {
	// Reserve a temporary local port, then release it for Fiber.
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve test port: %v", err)
	}

	address := probe.Addr().String()

	if err := probe.Close(); err != nil {
		t.Fatalf("release test port: %v", err)
	}

	app := fiber.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- Listen(ctx, app, address)
	}()

	// Wait until the server accepts a TCP connection.
	deadline := time.Now().Add(3 * time.Second)
	connected := false

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			connected = true
			break
		}

		select {
		case listenErr := <-done:
			t.Fatalf("server stopped before accepting connections: %v", listenErr)
		default:
		}

		time.Sleep(20 * time.Millisecond)
	}

	if !connected {
		cancel()
		t.Fatal("server did not start listening within 3 seconds")
	}

	// Cancel the context to trigger graceful shutdown.
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("graceful shutdown returned an error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("server did not shut down within 3 seconds")
	}

	// Confirm the server no longer accepts connections.
	conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
	if err == nil {
		_ = conn.Close()
		t.Fatal("server is still accepting connections after shutdown")
	}
}
