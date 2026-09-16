package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolTelemetry(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&output, nil)))
	defer slog.SetDefault(previous)
	for _, tc := range []struct {
		name    string
		result  *mcp.CallToolResult
		err     error
		outcome string
	}{
		{"known", &mcp.CallToolResult{}, nil, "ok"},
		{"known", &mcp.CallToolResult{IsError: true}, nil, "error"},
		{"secret-tool-name", nil, errors.New("secret-error"), "error"},
	} {
		output.Reset()
		handler := Tools(map[string]bool{"known": true})(func(context.Context, string, mcp.Request) (mcp.Result, error) { return tc.result, tc.err })
		req := &mcp.ServerRequest[*mcp.CallToolParamsRaw]{Params: &mcp.CallToolParamsRaw{Name: tc.name, Arguments: json.RawMessage(`{"token":"secret-argument"}`)}}
		result, err := handler(context.WithValue(context.Background(), requestKey{}, "correlation"), "tools/call", req)
		if result != tc.result || err != tc.err {
			t.Fatal("middleware changed result")
		}
		var record map[string]any
		if err := json.Unmarshal(output.Bytes(), &record); err != nil {
			t.Fatal(err)
		}
		if record["outcome"] != tc.outcome || record["request_id"] != "correlation" {
			t.Fatal(record)
		}
		if strings.Contains(output.String(), "secret") {
			t.Fatal("sensitive data logged")
		}
	}
}

func TestHTTPMetadataAndStreaming(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&output, nil)))
	defer slog.SetDefault(previous)
	handler := HTTP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(requestKey{}) == nil {
			t.Error("missing correlation")
		}
		w.WriteHeader(http.StatusAccepted)
		if err := http.NewResponseController(w).Flush(); err != nil {
			t.Error(err)
		}
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest("POST", "/mcp?token=secret", strings.NewReader("secret-body")))
	if recorder.Code != 202 || !recorder.Flushed || recorder.Header().Get("X-Request-ID") == "" {
		t.Fatal("HTTP semantics changed")
	}
	if strings.Contains(output.String(), "secret") || !strings.Contains(output.String(), `"status":202`) {
		t.Fatal(output.String())
	}
}
