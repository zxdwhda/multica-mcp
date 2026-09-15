package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const scope = "multica:access"

type Server struct {
	Origin, Prefix, PAT string
	Store               Store
}
type client struct {
	Redirects []string `json:"redirect_uris"`
}
type grant struct {
	Client, Redirect, Challenge, Resource, CSRF, Family string
	Expires                                             int64
}

func random() string {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
func (s *Server) issuer() string   { return s.Origin + s.Prefix }
func (s *Server) resource() string { return s.issuer() + "/mcp" }
func (s *Server) metadata() string {
	return s.Origin + "/.well-known/oauth-protected-resource" + s.Prefix + "/mcp"
}
func jsonOut(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, status int, code string) {
	jsonOut(w, status, map[string]string{"error": code})
}
func (s *Server) Routes(m *http.ServeMux) {
	for _, p := range []string{"/.well-known/oauth-protected-resource" + s.Prefix + "/mcp", s.Prefix + "/.well-known/oauth-protected-resource"} {
		m.HandleFunc("GET "+p, s.resourceMetadata)
	}
	for _, p := range []string{"/.well-known/oauth-authorization-server" + s.Prefix, s.Prefix + "/.well-known/oauth-authorization-server"} {
		m.HandleFunc("GET "+p, s.authMetadata)
	}
	m.HandleFunc("POST "+s.Prefix+"/register", s.register)
	m.HandleFunc("GET "+s.Prefix+"/authorize", s.authorize)
	m.HandleFunc("POST "+s.Prefix+"/authorize", s.approve)
	m.HandleFunc("POST "+s.Prefix+"/token", s.token)
	m.HandleFunc("POST "+s.Prefix+"/revoke", s.revoke)
}
func (s *Server) resourceMetadata(w http.ResponseWriter, r *http.Request) {
	jsonOut(w, 200, map[string]any{"resource": s.resource(), "authorization_servers": []string{s.issuer()}, "scopes_supported": []string{scope}, "bearer_methods_supported": []string{"header"}})
}
func (s *Server) authMetadata(w http.ResponseWriter, r *http.Request) {
	jsonOut(w, 200, map[string]any{"issuer": s.issuer(), "authorization_endpoint": s.issuer() + "/authorize", "token_endpoint": s.issuer() + "/token", "registration_endpoint": s.issuer() + "/register", "revocation_endpoint": s.issuer() + "/revoke", "response_types_supported": []string{"code"}, "grant_types_supported": []string{"authorization_code", "refresh_token"}, "token_endpoint_auth_methods_supported": []string{"none"}, "code_challenge_methods_supported": []string{"S256"}, "scopes_supported": []string{scope}})
}
func validRedirect(raw string) bool {
	u, e := url.Parse(raw)
	return e == nil && u.Scheme == "https" && u.Host != "" && u.User == nil && u.Fragment == "" && len(raw) < 2048
}
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Redirects []string `json:"redirect_uris"`
		Auth      string   `json:"token_endpoint_auth_method"`
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192)).Decode(&in) != nil || len(in.Redirects) < 1 || len(in.Redirects) > 10 || (in.Auth != "" && in.Auth != "none") {
		fail(w, 400, "invalid_client_metadata")
		return
	}
	for _, u := range in.Redirects {
		if !validRedirect(u) {
			fail(w, 400, "invalid_redirect_uri")
			return
		}
	}
	id := random()
	if s.Store.Put(r.Context(), "clients/"+hash(id), client{in.Redirects}) != nil {
		fail(w, 503, "temporarily_unavailable")
		return
	}
	jsonOut(w, 201, map[string]any{"client_id": id, "redirect_uris": in.Redirects, "client_id_issued_at": time.Now().Unix(), "token_endpoint_auth_method": "none", "grant_types": []string{"authorization_code", "refresh_token"}, "response_types": []string{"code"}})
}

var consent = template.Must(template.New("consent").Parse(`<!doctype html><html lang="zh-CN"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>连接 Multica</title><body><main><h1>连接 WildFlow Multica</h1><p>授权此客户端读取及操作当前工作区的项目、任务、评论与 Agent。</p><p>授权后返回：<strong>{{.Redirect}}</strong></p><form method="post" action="{{.Action}}"><input type="hidden" name="request" value="{{.Request}}"><label>当前 Multica PAT <input type="password" name="pat" required autocomplete="off"></label><p>使用已配置的 Multica PAT 确认身份。PAT 不会发送给客户端。</p><button type="submit">授权连接</button></form></main></body></html>`))

func (s *Server) authorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var c client
	if s.Store.Get(r.Context(), "clients/"+hash(q.Get("client_id")), &c) != nil {
		fail(w, 400, "invalid_client")
		return
	}
	redirect := q.Get("redirect_uri")
	found := false
	for _, u := range c.Redirects {
		if u == redirect {
			found = true
		}
	}
	challenge, e := base64.RawURLEncoding.DecodeString(q.Get("code_challenge"))
	if !found || q.Get("response_type") != "code" || q.Get("resource") != s.resource() || q.Get("code_challenge_method") != "S256" || e != nil || len(challenge) != 32 || (q.Get("scope") != "" && q.Get("scope") != scope) {
		fail(w, 400, "invalid_request")
		return
	}
	nonce, csrf := random(), random()
	g := grant{Client: q.Get("client_id"), Redirect: redirect, Challenge: q.Get("code_challenge"), Resource: s.resource(), CSRF: hash(csrf), Expires: time.Now().Add(10 * time.Minute).Unix()}
	// Store state separately from redirect URI; never reflect unvalidated redirect destinations.
	record := struct {
		Grant grant
		State string
	}{g, q.Get("state")}
	if s.Store.Put(r.Context(), "pending/"+hash(nonce), record) != nil {
		fail(w, 503, "temporarily_unavailable")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "__Secure-multica-consent", Value: csrf, Path: s.Prefix + "/authorize", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: 600})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; form-action 'self'; frame-ancestors 'none'")
	w.Header().Set("Referrer-Policy", "no-referrer")
	_ = consent.Execute(w, map[string]string{"Redirect": redirect, "Action": s.Prefix + "/authorize", "Request": nonce})
}
func (s *Server) approve(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16384)
	if r.ParseForm() != nil || r.Header.Get("Origin") != s.Origin {
		fail(w, 400, "invalid_request")
		return
	}
	var record struct {
		Grant grant
		State string
	}
	id := hash(r.PostForm.Get("request"))
	cookie, e := r.Cookie("__Secure-multica-consent")
	if e != nil || s.Store.Get(r.Context(), "pending/"+id, &record) != nil || record.Grant.Expires <= time.Now().Unix() || subtle.ConstantTimeCompare([]byte(hash(cookie.Value)), []byte(record.Grant.CSRF)) != 1 {
		fail(w, 400, "invalid_request")
		return
	}
	supplied, wanted := sha256.Sum256([]byte(r.PostForm.Get("pat"))), sha256.Sum256([]byte(s.PAT))
	if subtle.ConstantTimeCompare(supplied[:], wanted[:]) != 1 {
		fail(w, 401, "access_denied")
		return
	}
	if s.Store.Claim(r.Context(), "consent/"+id) != nil {
		fail(w, 400, "invalid_request")
		return
	}
	code := random()
	g := record.Grant
	g.CSRF = ""
	g.Expires = time.Now().Add(5 * time.Minute).Unix()
	g.Family = random()
	if s.Store.Put(r.Context(), "codes/"+hash(code), g) != nil {
		fail(w, 503, "temporarily_unavailable")
		return
	}
	u, _ := url.Parse(g.Redirect)
	q := u.Query()
	q.Set("code", code)
	q.Set("state", record.State)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusSeeOther)
}
func (s *Server) active(r *http.Request, g grant) bool {
	var marker any
	e := s.Store.Get(r.Context(), "revoked/"+hash(g.Family), &marker)
	return g.Expires > time.Now().Unix() && g.Resource == s.resource() && errors.Is(e, ErrMissing)
}
func (s *Server) token(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16384)
	if r.ParseForm() != nil {
		fail(w, 400, "invalid_request")
		return
	}
	f := r.PostForm
	kind := f.Get("grant_type")
	path := ""
	switch kind {
	case "authorization_code":
		path = "codes/" + hash(f.Get("code"))
	case "refresh_token":
		path = "refresh/" + hash(f.Get("refresh_token"))
	default:
		fail(w, 400, "unsupported_grant_type")
		return
	}
	var g grant
	if s.Store.Get(r.Context(), path, &g) != nil || !s.active(r, g) || f.Get("client_id") != g.Client || (f.Get("resource") != "" && f.Get("resource") != g.Resource) {
		fail(w, 400, "invalid_grant")
		return
	}
	if kind == "authorization_code" {
		v := f.Get("code_verifier")
		h := sha256.Sum256([]byte(v))
		if len(v) < 43 || len(v) > 128 || base64.RawURLEncoding.EncodeToString(h[:]) != g.Challenge || f.Get("redirect_uri") != g.Redirect {
			fail(w, 400, "invalid_grant")
			return
		}
	}
	if e := s.Store.Claim(r.Context(), path); e != nil {
		if errors.Is(e, ErrUsed) && kind == "refresh_token" {
			_ = s.Store.Put(r.Context(), "revoked/"+hash(g.Family), true)
		}
		fail(w, 400, "invalid_grant")
		return
	}
	access, refresh := random(), random()
	g.Challenge = ""
	g.Redirect = ""
	g.Expires = time.Now().Add(time.Hour).Unix()
	if s.Store.Put(r.Context(), "access/"+hash(access), g) != nil {
		fail(w, 503, "temporarily_unavailable")
		return
	}
	g.Expires = time.Now().Add(30 * 24 * time.Hour).Unix()
	if s.Store.Put(r.Context(), "refresh/"+hash(refresh), g) != nil {
		fail(w, 503, "temporarily_unavailable")
		return
	}
	jsonOut(w, 200, map[string]any{"access_token": access, "refresh_token": refresh, "token_type": "Bearer", "expires_in": 3600, "scope": scope})
}
func (s *Server) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		var g grant
		if !strings.HasPrefix(h, "Bearer ") || s.Store.Get(r.Context(), "access/"+hash(strings.TrimPrefix(h, "Bearer ")), &g) != nil || !s.active(r, g) {
			w.Header().Set("WWW-Authenticate", `Bearer resource_metadata="`+s.metadata()+`", scope="`+scope+`"`)
			fail(w, 401, "invalid_token")
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (s *Server) revoke(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if r.ParseForm() != nil {
		fail(w, 400, "invalid_request")
		return
	}
	for _, kind := range []string{"access/", "refresh/"} {
		var g grant
		if s.Store.Get(r.Context(), kind+hash(r.PostForm.Get("token")), &g) == nil {
			if g.Client != r.PostForm.Get("client_id") {
				fail(w, 400, "invalid_client")
				return
			}
			e := s.Store.Put(r.Context(), "revoked/"+hash(g.Family), true)
			if e != nil {
				var marker any
				if s.Store.Get(r.Context(), "revoked/"+hash(g.Family), &marker) != nil {
					fail(w, 503, "temporarily_unavailable")
					return
				}
			}
		}
	}
	jsonOut(w, 200, map[string]any{})
}
