package mcp

import (
	"context"
	"encoding/json"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"multica-mcp/internal/app"
	"multica-mcp/internal/multica"
	"testing"
)

func TestRegisteredSchemaRejectsFractionalStage(t *testing.T) {
	s := NewServer(app.NewUseCase(multica.NewClient("http://127.0.0.1:1", "test", "test"), false), false)
	a, b := mcp.NewInMemoryTransports()
	ctx := context.Background()
	session, e := s.GetMCPServer().Connect(ctx, a, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer session.Close()
	c := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	cs, e := c.Connect(ctx, b, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer cs.Close()
	r, e := cs.CallTool(ctx, &mcp.CallToolParams{Name: "multica_create_task", Arguments: map[string]any{"title": "test", "description": "test", "stage": 1.5, "dry_run": true}})
	if e != nil {
		t.Fatal(e)
	}
	if !r.IsError {
		t.Fatal("fractional stage silently accepted")
	}
}
func TestEmptyDescriptionIsPresent(t *testing.T) {
	r := &mcp.CallToolRequest{Params: &mcp.CallToolParamsRaw{Arguments: json.RawMessage(`{"description":""}`)}}
	if p := argsGetOptionalStringPtr(r, "description"); p == nil || *p != "" {
		t.Fatal("empty description lost")
	}
}
