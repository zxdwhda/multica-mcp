package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"multica-mcp/internal/apicatalog"
)

func (s *Server) RegisterAPI(caller apicatalog.Caller, readOnly bool) {
	catalog := apicatalog.Load()
	for _, op := range catalog.Operations {
		if readOnly && !op.ReadOnly {
			continue
		}
		destructive, open := !op.ReadOnly, true
		description := fmt.Sprintf("%s. %s %s. Full Multica App API operation; account/workspace permissions and feature availability are enforced by Multica. Body fields not listed in this schema are forwarded unchanged. Inspect status and body for results.", op.Handler, op.Method, op.Path)
		if !op.ReadOnly {
			description += " Changes remote state. Depending on the operation, this may start agent execution, delete data, change membership or affect billing."
		}
		if strings.Contains(op.Path, "/tasks") || strings.Contains(op.Path, "/agents") {
			description += " Tasks here are execution runs; issues are work items."
		}
		tool := &mcp.Tool{Name: op.Name, Description: description, InputSchema: op.Input, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: op.ReadOnly, DestructiveHint: &destructive, OpenWorldHint: &open, IdempotentHint: op.ReadOnly}}
		mcp.AddTool(s.mcpServer, tool, func(ctx context.Context, req *mcp.CallToolRequest, input map[string]any) (*mcp.CallToolResult, any, error) {
			b, e := json.Marshal(input)
			if e != nil {
				return errorResult(op.Name, e), nil, nil
			}
			var in apicatalog.Input
			if e = json.Unmarshal(b, &in); e != nil {
				return errorResult(op.Name, e), nil, nil
			}
			out, e := caller.CallAPI(ctx, op, in)
			if e != nil {
				return errorResult(op.Name, e), nil, nil
			}
			result := jsonResult(out)
			result.IsError = out.Status >= 400
			return result, nil, nil
		})
	}
	s.addTool(newTool("multica_api_catalog", "Find available Multica API operations by name, route or module. Returns input schemas, route provenance and pinned source revision. Empty query returns module counts; use query to find operations. All matching operations are registered as individual tools.", properties(stringProp("query", "Operation name or route fragment")), nil), func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		q := strings.ToLower(argsGetString(req, "query"))
		matches := []apicatalog.Operation{}
		modules := map[string]int{}
		for _, op := range catalog.Operations {
			if readOnly && !op.ReadOnly {
				continue
			}
			parts := strings.Split(op.Path, "/")
			modules[parts[2]]++
			if q != "" && strings.Contains(strings.ToLower(op.Name+" "+op.Path), q) {
				matches = append(matches, op)
			}
		}
		return jsonResult(map[string]any{"revision": catalog.Revision, "modules": modules, "matches": matches}), nil
	})
}
