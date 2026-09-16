package logging

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type requestKey struct{}

// HTTP records metadata only. Unwrap preserves ResponseController streaming support.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
func (w *responseWriter) WriteHeader(status int) {
	if status >= 100 && status < 200 {
		w.ResponseWriter.WriteHeader(status)
		return
	}
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(200)
	}
	return w.ResponseWriter.Write(b)
}
func HTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		var id [16]byte
		if _, err := rand.Read(id[:]); err != nil {
			panic(err)
		}
		requestID := hex.EncodeToString(id[:])
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestKey{}, requestID)
		rw := &responseWriter{ResponseWriter: w}
		defer func() {
			status := rw.status
			if status == 0 {
				status = 200
			}
			slog.Info("http_request", "event", "http_request", "request_id", requestID, "status", status, "duration_ms", time.Since(started).Milliseconds())
		}()
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// Tools includes isError responses, which may still have HTTP status 200.
// known is populated during registration, before accepting any requests.
func Tools(known map[string]bool) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			if method != "tools/call" {
				return next(ctx, method, req)
			}
			started := time.Now()
			operation := "unknown"
			if p, ok := req.GetParams().(*mcp.CallToolParamsRaw); ok && known[p.Name] {
				operation = p.Name
			}
			result, err := next(ctx, method, req)
			outcome := "ok"
			if err != nil {
				outcome = "error"
			}
			if r, ok := result.(*mcp.CallToolResult); ok && r != nil && r.IsError {
				outcome = "error"
			}
			requestID, _ := ctx.Value(requestKey{}).(string)
			slog.Info("mcp_tool", "event", "mcp_tool", "operation", operation, "outcome", outcome, "request_id", requestID, "duration_ms", time.Since(started).Milliseconds())
			return result, err
		}
	}
}
