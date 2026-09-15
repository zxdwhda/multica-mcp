package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"multica-mcp/internal/apicatalog"
	"multica-mcp/internal/app"
	"multica-mcp/internal/multica"
)

type fakeAPI struct{ calls []string }

func (f *fakeAPI) CallAPI(_ context.Context, op apicatalog.Operation, in apicatalog.Input) (*apicatalog.Response, error) {
	f.calls = append(f.calls, op.Name)
	return &apicatalog.Response{Status: 200, Body: map[string]any{"operation": op.Name}}, nil
}
func TestEveryCatalogOperationRegisteredAndCallable(t *testing.T) {
	for _, readOnly := range []bool{false, true} {
		s := NewServer(app.NewUseCase(multica.NewClient("http://127.0.0.1:1", "test", "test"), readOnly), readOnly)
		f := &fakeAPI{}
		s.RegisterAPI(f, readOnly)
		a, b := mcp.NewInMemoryTransports()
		ctx := t.Context()
		ss, e := s.GetMCPServer().Connect(ctx, a, nil)
		if e != nil {
			t.Fatal(e)
		}
		defer ss.Close()
		client := mcp.NewClient(&mcp.Implementation{Name: "catalog-test", Version: "1"}, nil)
		cs, e := client.Connect(ctx, b, nil)
		if e != nil {
			t.Fatal(e)
		}
		defer cs.Close()
		listed, e := cs.ListTools(ctx, nil)
		if e != nil {
			t.Fatal(e)
		}
		names := map[string]bool{}
		for _, tool := range listed.Tools {
			names[tool.Name] = true
		}
		count := 0
		for _, op := range apicatalog.Load().Operations {
			if readOnly && !op.ReadOnly {
				if names[op.Name] {
					t.Fatal("write exposed in readonly", op.Name)
				}
				continue
			}
			count++
			if !names[op.Name] {
				t.Fatal("missing", op.Name)
			}
			args := map[string]any{}
			path := map[string]string{}
			for _, part := range strings.Split(op.Path, "/") {
				if strings.HasPrefix(part, "{") {
					path[strings.Trim(part, "{}")] = "test-id"
				}
			}
			if len(path) > 0 {
				args["path"] = path
			}
			r, e := cs.CallTool(ctx, &mcp.CallToolParams{Name: op.Name, Arguments: args})
			if e != nil || r.IsError {
				t.Fatalf("%s: %v %v", op.Name, e, r)
			}
		}
		if len(f.calls) != count {
			t.Fatal("operations not dispatched", len(f.calls), count)
		}
		if !readOnly && len(listed.Tools) != len(apicatalog.Load().Operations)+17 {
			t.Fatal("unexpected full tool count", len(listed.Tools))
		}
	}
}
