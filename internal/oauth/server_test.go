package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

type memStore struct {
	sync.Mutex
	m map[string][]byte
}

func (s *memStore) Put(_ context.Context, k string, v any) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.m[k]; ok {
		return ErrUsed
	}
	s.m[k], _ = json.Marshal(v)
	return nil
}
func (s *memStore) Get(_ context.Context, k string, v any) error {
	s.Lock()
	defer s.Unlock()
	b, ok := s.m[k]
	if !ok {
		return ErrMissing
	}
	return json.Unmarshal(b, v)
}
func (s *memStore) Claim(c context.Context, k string) error { return s.Put(c, "claims/"+k, true) }
func fixture() (*Server, *http.ServeMux) {
	s := &Server{Origin: "https://mcp.example.com", Prefix: "/multica", PAT: "test-pat", Store: &memStore{m: map[string][]byte{}}}
	m := http.NewServeMux()
	s.Routes(m)
	m.Handle("/multica/mcp", s.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })))
	return s, m
}
func post(m http.Handler, path string, f url.Values) *httptest.ResponseRecorder {
	r := httptest.NewRequest("POST", path, strings.NewReader(f.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	return w
}
func TestPKCEReplayRefreshAndRevocation(t *testing.T) {
	s, m := fixture()
	v := strings.Repeat("a", 43)
	h := sha256.Sum256([]byte(v))
	g := grant{Client: "client", Redirect: "https://chatgpt.com/callback", Challenge: base64.RawURLEncoding.EncodeToString(h[:]), Resource: s.resource(), Family: "family", Expires: time.Now().Add(time.Minute).Unix()}
	_ = s.Store.Put(t.Context(), "codes/"+hash("code"), g)
	f := url.Values{"grant_type": {"authorization_code"}, "client_id": {"client"}, "code": {"code"}, "redirect_uri": {g.Redirect}, "code_verifier": {"bad"}, "resource": {s.resource()}}
	if w := post(m, "/multica/token", f); w.Code != 400 {
		t.Fatal(w.Code)
	}
	f.Set("code_verifier", v)
	f.Set("resource", "https://other.example/mcp")
	if w := post(m, "/multica/token", f); w.Code != 400 {
		t.Fatal(w.Code)
	}
	f.Set("resource", s.resource())
	var wg sync.WaitGroup
	statuses := make(chan int, 2)
	tokens := make(chan map[string]any, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := post(m, "/multica/token", f)
			statuses <- w.Code
			if w.Code == 200 {
				var out map[string]any
				_ = json.Unmarshal(w.Body.Bytes(), &out)
				tokens <- out
			}
		}()
	}
	wg.Wait()
	close(statuses)
	success := 0
	for c := range statuses {
		if c == 200 {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("single use successes=%d", success)
	}
	out := <-tokens
	access := out["access_token"].(string)
	request := func(token string) int {
		r := httptest.NewRequest("POST", "/multica/mcp", nil)
		r.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		return w.Code
	}
	if request(access) != 204 {
		t.Fatal("access failed")
	}
	if request("test-pat") != 401 {
		t.Fatal("PAT accepted as MCP token")
	}
	refresh := url.Values{"grant_type": {"refresh_token"}, "client_id": {"client"}, "refresh_token": {out["refresh_token"].(string)}}
	w := post(m, "/multica/token", refresh)
	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
	if post(m, "/multica/token", refresh).Code != 400 {
		t.Fatal("refresh replay accepted")
	}
	if request(access) != 401 {
		t.Fatal("refresh replay did not revoke family")
	}
}
func TestConsentRequiresPATCookieOrigin(t *testing.T) {
	s, m := fixture()
	g := grant{Client: "client", Redirect: "https://chatgpt.com/callback", Resource: s.resource(), CSRF: hash("cookie"), Expires: time.Now().Add(time.Minute).Unix()}
	_ = s.Store.Put(t.Context(), "pending/"+hash("pending"), struct {
		Grant grant
		State string
	}{g, "state"})
	for _, tc := range []struct {
		origin, cookie, pat string
		want                int
	}{{"", "cookie", "test-pat", 400}, {s.Origin, "wrong", "test-pat", 400}, {s.Origin, "cookie", "wrong", 401}, {s.Origin, "cookie", "test-pat", 303}, {s.Origin, "cookie", "test-pat", 400}} {
		values := url.Values{"request": {"pending"}, "pat": {tc.pat}}
		r := httptest.NewRequest("POST", "/multica/authorize", strings.NewReader(values.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Set("Origin", tc.origin)
		r.AddCookie(&http.Cookie{Name: "__Secure-multica-consent", Value: tc.cookie})
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("got %d want %d: %s", w.Code, tc.want, w.Body.String())
		}
	}
}
func TestDCRRejectsUnsafeRedirect(t *testing.T) {
	_, m := fixture()
	for _, u := range []string{"http://chatgpt.com/callback", "javascript:alert(1)", "https://evil.test/#fragment"} {
		b, _ := json.Marshal(map[string]any{"redirect_uris": []string{u}})
		w := httptest.NewRecorder()
		m.ServeHTTP(w, httptest.NewRequest("POST", "/multica/register", strings.NewReader(string(b))))
		if w.Code != 400 {
			t.Fatal(u, w.Code)
		}
	}
}
