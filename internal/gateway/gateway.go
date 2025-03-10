package gateway

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/roushou/pocpoc/internal/router"
	"gorm.io/gorm"
)

const defaultAddr = ":8080"
const defaultShutdownTimeout = 5 * time.Second

// Option defines the function signature for gateway options.
type Option func(options *options) error

type options struct {
	addr            string
	shutdownTimeout time.Duration
}

// WithAddr sets the address on which the gateway will listen.
func WithAddr(addr string) Option {
	return func(options *options) error {
		if addr == "" {
			return errors.New("address should not be empty")
		}
		options.addr = addr
		return nil
	}
}

// WithShutdownTimeout sets the timeout for server shutdown.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(options *options) error {
		options.shutdownTimeout = timeout
		return nil
	}
}

type Gateway struct {
	addr            string
	router          *echo.Echo
	shutdownTimeout time.Duration
}

// NewGateway initializes and configures a new Gateway instance with the provided options.
func NewGateway(database *gorm.DB, opts ...Option) (*Gateway, error) {
	options := &options{
		addr:            defaultAddr,
		shutdownTimeout: defaultShutdownTimeout,
	}
	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	return &Gateway{
		router:          router.NewRouter(database),
		addr:            options.addr,
		shutdownTimeout: options.shutdownTimeout,
	}, nil
}

// Serve starts the gateway server and handles graceful shutdown on receiving interrupt signals.
func (gw *Gateway) Serve(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := gw.router.Start(gw.addr)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, gw.shutdownTimeout)
	defer cancel()

	if err := gw.router.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
