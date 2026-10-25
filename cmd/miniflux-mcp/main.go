package main

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	miniflux "miniflux.app/v2/client"

	"github.com/dmingn/miniflux-mcp/internal/config"
	mcpserver "github.com/dmingn/miniflux-mcp/internal/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "miniflux-mcp: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := miniflux.NewClient(cfg.URL, cfg.APIKey)
	if _, err := client.Me(); err != nil {
		return fmt.Errorf("failed to authenticate with Miniflux: %w", err)
	}

	s := mcpserver.New(client)
	return s.Run(context.Background(), &mcp.StdioTransport{})
}
