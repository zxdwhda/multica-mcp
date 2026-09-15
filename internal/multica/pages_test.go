package multica

import (
	"encoding/json"
	"multica-mcp/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchPageFiltersAndPreservesCursor(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/issues/search" || r.URL.Query().Get("offset") != "50" {
			t.Errorf("unexpected %s", r.URL)
		}
		json.NewEncoder(w).Encode(map[string]any{"total": 102, "issues": []map[string]any{{"id": "wrong-project", "project_id": "other", "status": "todo"}, {"id": "correct", "project_id": "wanted", "status": "todo", "assignee_id": "agent"}}})
	}))
	defer srv.Close()
	c := NewClient(srv.URL, "fake", "")
	c.SetWorkspaceScope("ws", "")
	q := "query"
	a := "agent"
	p, e := c.TasksPage(t.Context(), domain.ListTasksInput{ProjectID: "wanted", Query: &q, Assignee: &a}, 50, 50)
	if e != nil {
		t.Fatal(e)
	}
	if len(p.Items) != 1 || p.Items[0].ID != "correct" || !p.FilteredOnPage || p.NextOffset == nil || *p.NextOffset != 52 || p.SourceTotal != 102 {
		t.Fatalf("%+v", p)
	}
}
