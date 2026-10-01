package multica

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"multica-mcp/internal/apicatalog"
)

func TestAPIRequestPreservesValuesAndErrorEnvelope(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" || r.URL.Path != "/api/projects/test-id" || r.Header.Get("Authorization") != "Bearer secret" || r.Header.Get("X-Workspace-ID") != "ws" || r.Header.Get("If-Match") != "revision" {
			t.Errorf("incorrect request")
		}
		if strings.Join(r.URL.Query()["status"], ",") != "done,backlog" {
			t.Errorf("query array lost")
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if v, ok := body["description"]; !ok || v != "" {
			t.Errorf("empty description lost")
		}
		if v, ok := body["due_date"]; !ok || v != nil {
			t.Errorf("null lost")
		}
		w.Header().Set("ETag", "next")
		w.WriteHeader(409)
		io.WriteString(w, `{"code":"revision_conflict","extra":{"revision":2}}`)
	}))
	defer api.Close()
	c := NewClient(api.URL, "secret", "test")
	c.SetWorkspaceScope("ws", "")
	op := apicatalog.Operation{Method: "PATCH", Path: "/api/projects/{id}"}
	result, err := c.CallAPI(t.Context(), op, apicatalog.Input{Path: map[string]string{"id": "test-id"}, Headers: map[string]string{"If-Match": "revision"}, Query: map[string]any{"status": []any{"done", "backlog"}}, Body: map[string]any{"description": "", "due_date": nil}})
	if err != nil || result.Status != 409 || result.Headers["ETag"] != "next" {
		t.Fatal(result, err)
	}
	if result.Body.(map[string]any)["extra"] == nil {
		t.Fatal("error detail lost")
	}
}
func TestAPITransportBoundaries(t *testing.T) {
	calls := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Location", "https://example.invalid/secret")
		w.WriteHeader(302)
	}))
	defer api.Close()
	c := NewClient(api.URL, "secret", "test")
	c.SetWorkspaceScope("ws", "")
	op := apicatalog.Operation{Method: "GET", Path: "/api/issues/{id}"}
	for _, id := range []string{"..", "a/b", "a%2fb", "a?b", "a\\b"} {
		if _, e := c.CallAPI(t.Context(), op, apicatalog.Input{Path: map[string]string{"id": id}}); e == nil {
			t.Fatal("path escape accepted", id)
		}
	}
	if _, e := c.CallAPI(t.Context(), op, apicatalog.Input{Path: map[string]string{"id": "ok"}, Headers: map[string]string{"Authorization": "other"}}); e == nil {
		t.Fatal("credential override accepted")
	}
	if calls != 0 {
		t.Fatal("invalid input reached upstream")
	}
	result, e := c.CallAPI(t.Context(), op, apicatalog.Input{Path: map[string]string{"id": "ok"}})
	if e != nil || result.Status != 302 || calls != 1 {
		t.Fatal("redirect followed", result, e)
	}
}
func TestAPIMultipartAndEmptyResponse(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := r.ParseMultipartForm(1024); e != nil {
			t.Fatal(e)
		}
		defer r.MultipartForm.RemoveAll()
		file, _, e := r.FormFile("file")
		if e != nil {
			t.Fatal(e)
		}
		defer file.Close()
		data, _ := io.ReadAll(file)
		if string(data) != "test bytes" || r.FormValue("description") != "test" {
			t.Error("multipart values lost")
		}
		w.WriteHeader(204)
	}))
	defer api.Close()
	c := NewClient(api.URL, "secret", "test")
	c.SetWorkspaceScope("ws", "")
	result, e := c.CallAPI(t.Context(), apicatalog.Operation{Method: "POST", Path: "/api/upload-file"}, apicatalog.Input{Body: map[string]any{"description": "test"}, Files: []apicatalog.File{{Field: "file", Filename: "test.txt", Data: base64.StdEncoding.EncodeToString([]byte("test bytes"))}}})
	if e != nil || result.Status != 204 || result.Body != nil {
		t.Fatal(result, e)
	}
}

func TestAPICommentPaginationHeaders(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("roots_only") != "true" || r.URL.Query().Get("summary") != "true" {
			t.Error("comment filters lost")
		}
		w.Header().Set("X-Multica-Next-Before", "2026-10-01T00:00:00Z")
		w.Header().Set("X-Multica-Next-Before-Id", "cursor-id")
		w.Header().Set("X-Comments-Truncated", "true")
		w.Header().Set("Set-Cookie", "private=value")
		fmt.Fprint(w, `[{"id":"root","content_truncated":true,"reply_count":42}]`)
	}))
	defer api.Close()
	c := NewClient(api.URL, "secret", "test")
	c.SetWorkspaceScope("ws", "")
	r, err := c.CallAPI(t.Context(), apicatalog.Operation{Method: "GET", Path: "/api/issues/{id}/comments"}, apicatalog.Input{Path: map[string]string{"id": "issue"}, Query: map[string]any{"roots_only": true, "summary": true}})
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range []string{"X-Multica-Next-Before", "X-Multica-Next-Before-Id", "X-Comments-Truncated"} {
		if r.Headers[h] == "" {
			t.Errorf("missing %s", h)
		}
	}
	if r.Headers["Set-Cookie"] != "" || r.Body.([]any)[0].(map[string]any)["reply_count"] != float64(42) {
		t.Fatal("response fidelity or header boundary lost")
	}
}
