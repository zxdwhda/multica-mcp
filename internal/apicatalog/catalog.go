package apicatalog

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

//go:embed catalog.json
var source []byte

type Operation struct {
	Name     string         `json:"name"`
	Handler  string         `json:"handler"`
	Method   string         `json:"method"`
	Path     string         `json:"path"`
	ReadOnly bool           `json:"read_only"`
	Input    map[string]any `json:"input_schema"`
	Source   string         `json:"source"`
}
type Catalog struct {
	Revision   string              `json:"revision"`
	Operations []Operation         `json:"operations"`
	Excluded   []map[string]string `json:"excluded"`
}

func Load() Catalog {
	var c Catalog
	if err := json.Unmarshal(source, &c); err != nil {
		panic(err)
	}
	return c
}

type File struct {
	Field       string `json:"field"`
	Filename    string `json:"filename"`
	Data        string `json:"data_base64"`
	ContentType string `json:"content_type"`
}
type Input struct {
	Path    map[string]string `json:"path,omitempty"`
	Query   map[string]any    `json:"query,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    any               `json:"body,omitempty"`
	Files   []File            `json:"files,omitempty"`
}
type Response struct {
	Status      int               `json:"status"`
	ContentType string            `json:"content_type,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        any               `json:"body,omitempty"`
	DataBase64  string            `json:"data_base64,omitempty"`
	Truncated   bool              `json:"truncated,omitempty"`
}
type Caller interface {
	CallAPI(context.Context, Operation, Input) (*Response, error)
}

var placeholders = regexp.MustCompile(`\{([^}:]+)(?::[^}]+)?\}`)

func (op Operation) URL(in Input) (string, error) {
	path := op.Path
	for _, m := range placeholders.FindAllStringSubmatch(path, -1) {
		value := in.Path[m[1]]
		// A single identifier must never become another route after proxy decoding.
		if value == "" || value == "." || value == ".." || strings.ContainsAny(value, "/\\%?#\r\n") {
			return "", fmt.Errorf("invalid path parameter %s", m[1])
		}
		path = strings.ReplaceAll(path, m[0], url.PathEscape(value))
	}
	q := url.Values{}
	for key, v := range in.Query {
		if values, ok := v.([]any); ok {
			for _, x := range values {
				q.Add(key, fmt.Sprint(x))
			}
		} else if v != nil {
			q.Set(key, fmt.Sprint(v))
		}
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	return path, nil
}
