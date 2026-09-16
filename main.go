package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"multica-mcp/internal/app"
	"multica-mcp/internal/config"
	"multica-mcp/internal/logging"
	mcpserver "multica-mcp/internal/mcp"
	"multica-mcp/internal/middleware"
	"multica-mcp/internal/multica"
	"multica-mcp/internal/oauth"
	"multica-mcp/internal/version"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	logging.Setup(cfg.LogLevel)
	slog.Info("starting multica-mcp",
		"version", version.Version,
		"multica_api", version.MulticaAPI,
		"transport", cfg.MCPTransport,
		"read_only", cfg.ReadOnly,
	)

	client := multica.NewClient(cfg.MulticaBaseURL, cfg.MulticaToken, version.Version)

	slug := strings.TrimSpace(cfg.MulticaWorkspaceSlug)
	wsID := strings.TrimSpace(cfg.MulticaWorkspaceID)
	if slug != "" {
		client.SetWorkspaceScope("", slug)
		slog.Info("workspace scope from slug", "slug", slug)
	} else {
		if wsID == "" {
			wsID = resolveWorkspace(client)
		}
		if wsID == "" {
			fmt.Fprintf(os.Stderr, "MULTICA_WORKSPACE_ID or MULTICA_WORKSPACE_SLUG is required, or a single workspace must exist\n")
			os.Exit(1)
		}
		client.SetWorkspaceScope(wsID, "")
		slog.Info("workspace scope from id")
	}

	useCase := app.NewUseCase(client, cfg.ReadOnly)
	mcpSrv := mcpserver.NewServer(useCase, cfg.ReadOnly)
	mcpSrv.RegisterAPI(client, cfg.ReadOnly)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	switch cfg.MCPTransport {
	case "http":
		runHTTP(mcpSrv.GetMCPServer(), cfg, ctx)
	default:
		runStdio(ctx, mcpSrv.GetMCPServer())
	}
}

func resolveWorkspace(client *multica.Client) string {
	workspaces, err := client.ListWorkspaces(context.Background())
	if err != nil {
		slog.Warn("failed to list workspaces for auto-detection", "error", err)
		return ""
	}
	if len(workspaces) == 1 {
		slog.Info("auto-detected single workspace", "workspace_id", workspaces[0].ID, "name", workspaces[0].Name)
		return workspaces[0].ID
	}
	if len(workspaces) > 1 {
		fmt.Fprintf(os.Stderr, "multiple workspaces found; set MULTICA_WORKSPACE_ID or MULTICA_WORKSPACE_SLUG to one of:\n")
		for _, ws := range workspaces {
			slug := ws.Slug
			if slug == "" {
				slug = "(no slug)"
			}
			fmt.Fprintf(os.Stderr, "  %s (%s) slug=%s\n", ws.ID, ws.Name, slug)
		}
	}
	return ""
}

func runStdio(ctx context.Context, mcpServer *mcp.Server) {
	if err := mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("stdio server error", "error", err)
		os.Exit(1)
	}
}

func httpHandler(mcpServer *mcp.Server, cfg *config.Config) (http.Handler, error) {
	streamable := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return mcpServer }, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true})
	mux := http.NewServeMux()
	var handler http.Handler = streamable
	if cfg.OAuthOrigin != "" {
		store, e := oauth.NewOSSStore(cfg.OSSEndpoint, cfg.OSSBucket, cfg.MulticaToken)
		if e != nil {
			return nil, e
		}
		auth := &oauth.Server{Origin: cfg.OAuthOrigin, Prefix: cfg.HTTPPrefix, PAT: cfg.MulticaToken, Store: store, ValidatePAT: oauth.SameAccountPATValidator(cfg.MulticaBaseURL, cfg.MulticaToken)}
		auth.Routes(mux)
		handler = auth.Protect(handler)
	} else if cfg.MCPAPIKey != "" {
		handler = middleware.NewLoginThrottle(5, time.Minute, 15*time.Minute).Wrap(cfg.MCPAPIKey, handler)
	}
	mux.Handle(cfg.HTTPPrefix+"/mcp", handler)
	mux.HandleFunc("GET "+cfg.HTTPPrefix+"/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","version":%q,"revision":%q}`, version.Version, version.Revision)
	})
	return logging.HTTP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 12*1024*1024)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		mux.ServeHTTP(w, r)
	})), nil
}
func runHTTP(mcpServer *mcp.Server, cfg *config.Config, ctx context.Context) {
	handler, err := httpHandler(mcpServer, cfg)
	if err != nil {
		slog.Error("HTTP setup failed", "error", err)
		os.Exit(1)
	}
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", cfg.HTTPPort), Handler: handler, ReadHeaderTimeout: 10 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 90 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 32 * 1024}
	done := make(chan error, 1)
	go func() { done <- httpServer.ListenAndServe() }()
	slog.Info("HTTP server ready", "port", cfg.HTTPPort, "path", cfg.HTTPPrefix+"/mcp", "stateless", true)
	select {
	case err := <-done:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutdown failed", "error", err)
		}
	}
}
