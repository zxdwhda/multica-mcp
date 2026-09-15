package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"multica-mcp/internal/config"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPPathsAuthAndStateless(t *testing.T) {
	s := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1"}, &mcp.ServerOptions{GetSessionID: func() string { return "" }})
	h, e := httpHandler(s, &config.Config{HTTPPrefix: "/multica", MCPAPIKey: "key"})
	if e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		path, auth, body string
		want             int
	}{{"/multica/healthz", "", "", 200}, {"/feishu/mcp", "", "", 404}, {"/multica/mcp", "", "{}", 401}, {"/multica/mcp", "Bearer key", `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`, 200}, {"/multica/mcp", "Bearer key", `{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`, 200}} {
		method := "POST"
		if tc.body == "" {
			method = "GET"
		}
		r := httptest.NewRequest(method, tc.path, strings.NewReader(tc.body))
		r.RemoteAddr = "192.0.2.1:1234"
		r.Header.Set("Authorization", tc.auth)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "application/json, text/event-stream")
		r.Header.Set("MCP-Protocol-Version", "2025-06-18")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("%s %d %s", tc.path, w.Code, w.Body.String())
		}
		if w.Header().Get("Mcp-Session-Id") != "" {
			t.Fatal("stateful session")
		}
	}
}
