package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// Authenticate against Multica using an ephemeral PAT, while keeping all MCP
// calls bound to the configured account and durable upstream credential.
func SameAccountPATValidator(baseURL, boundPAT string) func(context.Context, string) (bool, error) {
	client := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	identity := func(ctx context.Context, pat string) (string, error) {
		if pat == "" {
			return "", nil
		}
		req, e := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+"/api/me", nil)
		if e != nil {
			return "", errors.New("identity endpoint invalid")
		}
		req.Header.Set("Authorization", "Bearer "+pat)
		resp, e := client.Do(req)
		if e != nil {
			return "", errors.New("identity service unavailable")
		}
		defer resp.Body.Close()
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return "", nil
		}
		if resp.StatusCode != 200 {
			return "", errors.New("identity service unavailable")
		}
		var user struct {
			ID string `json:"id"`
		}
		if json.NewDecoder(io.LimitReader(resp.Body, 128*1024)).Decode(&user) != nil || user.ID == "" {
			return "", errors.New("invalid identity response")
		}
		return user.ID, nil
	}
	return func(ctx context.Context, candidate string) (bool, error) {
		id, e := identity(ctx, candidate)
		if e != nil || id == "" {
			return false, e
		}
		owner, e := identity(ctx, boundPAT)
		if e != nil {
			return false, e
		}
		if owner == "" {
			return false, errors.New("configured PAT is no longer valid")
		}
		return id == owner, nil
	}
}
