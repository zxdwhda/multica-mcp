package oauth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSameAccountPATValidator(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/me" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		switch r.Header.Get("Authorization") {
		case "Bearer durable", "Bearer temporary":
			fmt.Fprint(w, `{"id":"owner"}`)
		case "Bearer other":
			fmt.Fprint(w, `{"id":"other"}`)
		case "Bearer unavailable":
			w.WriteHeader(503)
		case "Bearer malformed":
			fmt.Fprint(w, `{}`)
		case "Bearer redirect":
			http.Redirect(w, r, "/unexpected", 302)
		default:
			w.WriteHeader(401)
		}
	}))
	defer api.Close()
	for _, tc := range []struct {
		candidate, bound string
		want, wantErr    bool
	}{
		{"temporary", "durable", true, false},
		{"durable", "durable", true, false},
		{"other", "durable", false, false},
		{"expired", "durable", false, false},
		{"", "durable", false, false},
		{"unavailable", "durable", false, true},
		{"malformed", "durable", false, true},
		{"redirect", "durable", false, true},
		{"temporary", "expired", false, true},
	} {
		t.Run(tc.candidate+"/"+tc.bound, func(t *testing.T) {
			got, err := SameAccountPATValidator(api.URL, tc.bound)(t.Context(), tc.candidate)
			if got != tc.want || (err != nil) != tc.wantErr {
				t.Fatalf("got %v, %v", got, err)
			}
		})
	}
}
