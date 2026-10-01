// Package skillbundle serves the same immutable instructions shipped in the plugin.
package skillbundle

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed bundle.json
var bundle []byte

type resource struct {
	Name     string `json:"name"`
	MIMEType string `json:"mimeType"`
	Text     string `json:"text"`
}

// Register exposes only embedded resources, never filesystem paths or arbitrary URIs.
func Register(server *mcp.Server) {
	var resources map[string]resource
	if err := json.Unmarshal(bundle, &resources); err != nil {
		panic(fmt.Sprintf("invalid embedded skill bundle: %v", err))
	}
	keys := make([]string, 0, len(resources))
	for uri := range resources {
		keys = append(keys, uri)
	}
	sort.Strings(keys)
	for _, uri := range keys {
		item := resources[uri]
		server.AddResource(&mcp.Resource{URI: uri, Name: item.Name, MIMEType: item.MIMEType}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: uri, MIMEType: item.MIMEType, Text: item.Text}}}, nil
		})
	}
}
