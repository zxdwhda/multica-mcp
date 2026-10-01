package mcp

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"multica-mcp/internal/app"
	"multica-mcp/internal/multica"
	"net/http"
	"net/http/httptest"
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

func TestLightweightTaskAndSkillResources(t *testing.T) {
	calls := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/api/issues/id" {
			t.Errorf("unexpected section request %s", r.URL.Path)
		}
		fmt.Fprint(w, `{"id":"id","title":"Task"}`)
	}))
	defer api.Close()
	client := multica.NewClient(api.URL, "test", "test")
	client.SetWorkspaceScope("ws", "")
	s := NewServer(app.NewUseCase(client, true), true)
	a, b := mcp.NewInMemoryTransports()
	ss, err := s.GetMCPServer().Connect(t.Context(), a, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ss.Close()
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil).Connect(t.Context(), b, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()
	r, err := cs.CallTool(t.Context(), &mcp.CallToolParams{Name: "multica_get_task", Arguments: map[string]any{"task_id": "id", "include_comments": false, "include_subtasks": false}})
	if err != nil || r.IsError || calls != 1 {
		t.Fatalf("lightweight read: %+v %v calls=%d", r, err, calls)
	}
	resources, err := cs.ListResources(t.Context(), nil)
	if err != nil || len(resources.Resources) != 13 {
		t.Fatalf("resources: %+v %v", resources, err)
	}
	index, err := cs.ReadResource(t.Context(), &mcp.ReadResourceParams{URI: "skill://multica-mcp/index.json"})
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Skills []struct {
			Resources []struct {
				URI    string
				Digest string
			}
		}
	}
	if err := json.Unmarshal([]byte(index.Contents[0].Text), &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Skills) != 5 {
		t.Fatalf("missing skills: %+v", manifest)
	}
	for _, skill := range manifest.Skills {
		for _, ref := range skill.Resources {
			detail, err := cs.ReadResource(t.Context(), &mcp.ReadResourceParams{URI: ref.URI})
			if err != nil {
				t.Fatal(err)
			}
			if len(detail.Contents) != 1 || fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(detail.Contents[0].Text))) != ref.Digest {
				t.Fatalf("resource content differs from package: %s", ref.URI)
			}
		}
	}
	for _, uri := range []string{"file:///etc/passwd", "skill://multica-mcp/../secrets", "skill://multica-mcp/missing"} {
		if _, err := cs.ReadResource(t.Context(), &mcp.ReadResourceParams{URI: uri}); err == nil {
			t.Fatalf("unregistered resource accepted: %s", uri)
		}
	}
}
