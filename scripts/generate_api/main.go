// Generate a factual HTTP interface catalog from the official Go router and
// request declarations. No upstream handler implementation is embedded.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

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

var fset = token.NewFileSet()
var types = map[string]ast.Expr{}
var funcs = map[string]*ast.FuncDecl{}
var constants = map[string]string{}
var routerHelpers = map[string]*ast.FuncDecl{}
var catalog Catalog
var params = regexp.MustCompile(`\{([^}:]+)(?::[^}]+)?\}`)
var camel1 = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
var camel2 = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func snake(s string) string {
	return strings.ToLower(camel2.ReplaceAllString(camel1.ReplaceAllString(s, "${1}_${2}"), "${1}_${2}"))
}
func str(e ast.Expr) string {
	if i, ok := e.(*ast.Ident); ok {
		return constants[i.Name]
	}
	if s, ok := e.(*ast.SelectorExpr); ok {
		if p, ok := s.X.(*ast.Ident); ok {
			return constants[p.Name+"."+s.Sel.Name]
		}
	}
	if b, ok := e.(*ast.BasicLit); ok {
		s, _ := strconv.Unquote(b.Value)
		return s
	}
	return ""
}
func shape(e ast.Expr, depth int) map[string]any {
	if depth > 5 {
		return map[string]any{}
	}
	switch x := e.(type) {
	case *ast.StarExpr:
		return map[string]any{"anyOf": []any{shape(x.X, depth+1), map[string]any{"type": "null"}}}
	case *ast.Ident:
		if v, ok := types[x.Name]; ok {
			return shape(v, depth+1)
		}
		switch x.Name {
		case "string":
			return map[string]any{"type": "string"}
		case "bool":
			return map[string]any{"type": "boolean"}
		case "int", "int32", "int64", "uint", "uint32", "uint64":
			return map[string]any{"type": "integer"}
		case "float32", "float64":
			return map[string]any{"type": "number"}
		}
	case *ast.ArrayType:
		return map[string]any{"type": "array", "items": shape(x.Elt, depth+1)}
	case *ast.MapType:
		return map[string]any{"type": "object", "additionalProperties": shape(x.Value, depth+1)}
	case *ast.StructType:
		p := map[string]any{}
		for _, f := range x.Fields.List {
			if f.Tag == nil {
				continue
			}
			tag, _ := strconv.Unquote(f.Tag.Value)
			name := strings.Split(reflect.StructTag(tag).Get("json"), ",")[0]
			if name == "" || name == "-" {
				continue
			}
			p[name] = shape(f.Type, depth+1)
		}
		return map[string]any{"type": "object", "properties": p, "additionalProperties": true}
	case *ast.SelectorExpr:
		if x.Sel.Name == "UUID" || x.Sel.Name == "Time" {
			return map[string]any{"type": "string"}
		}
	}
	return map[string]any{}
}
func schema(handler, path, method string) map[string]any {
	properties := map[string]any{}
	required := []string{}
	pp := map[string]any{}
	pr := []string{}
	for _, m := range params.FindAllStringSubmatch(path, -1) {
		pp[m[1]] = map[string]any{"type": "string", "minLength": 1}
		pr = append(pr, m[1])
	}
	if len(pp) > 0 {
		properties["path"] = map[string]any{"type": "object", "properties": pp, "required": pr, "additionalProperties": false}
		required = append(required, "path")
	}
	properties["query"] = map[string]any{"type": "object", "description": "URL query parameters. Arrays are repeated query keys. Use upstream pagination fields.", "additionalProperties": true}
	properties["headers"] = map[string]any{"type": "object", "description": "Optional If-Match, Idempotency-Key, X-Multica-Plugin-Installation, Range, Accept headers. Authorization and workspace headers are configured by the server.", "additionalProperties": map[string]any{"type": "string"}}
	body := map[string]any{"type": "object", "additionalProperties": true}
	multipart := false
	if fn := funcs[handler]; fn != nil {
		vars := map[string]ast.Expr{}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if v, ok := n.(*ast.ValueSpec); ok {
				for _, name := range v.Names {
					vars[name.Name] = v.Type
				}
			}
			if c, ok := n.(*ast.CallExpr); ok {
				if s, ok := c.Fun.(*ast.SelectorExpr); ok {
					if s.Sel.Name == "FormFile" {
						multipart = true
					}
					if s.Sel.Name == "Decode" && len(c.Args) > 0 {
						if u, ok := c.Args[0].(*ast.UnaryExpr); ok {
							if i, ok := u.X.(*ast.Ident); ok {
								if t := vars[i.Name]; t != nil {
									body = shape(t, 0)
								}
							}
						}
					}
				}
			}
			return true
		})
		// Some handlers decode through shared helpers rather than Decoder.Decode.
		if len(body) == 2 {
			if t := vars["req"]; t != nil {
				body = shape(t, 0)
			}
		}
	}
	if method != "GET" && method != "HEAD" {
		properties["body"] = body
	}
	if multipart || strings.Contains(strings.ToLower(handler), "upload") || strings.Contains(handler, "PublishPlugin") {
		properties["files"] = map[string]any{"type": "array", "description": "Multipart uploads from base64 bytes. Other body fields become multipart text fields.", "items": map[string]any{"type": "object", "properties": map[string]any{"field": map[string]any{"type": "string"}, "filename": map[string]any{"type": "string"}, "data_base64": map[string]any{"type": "string"}, "content_type": map[string]any{"type": "string"}}, "required": []string{"field", "filename", "data_base64"}, "additionalProperties": false}}
	}
	return map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}
}
func authFor(b *ast.BlockStmt, old string) string {
	surface := old
	for _, st := range b.List {
		ast.Inspect(st, func(n ast.Node) bool {
			if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			if c, ok := n.(*ast.CallExpr); ok {
				if s, ok := c.Fun.(*ast.SelectorExpr); ok {
					switch s.Sel.Name {
					case "Auth":
						surface = "pat"
					case "DaemonAuth":
						surface = "daemon"
					}
				}
			}
			if s, ok := n.(*ast.SelectorExpr); ok && s.Sel.Name == "PluginBearerOnly" {
				surface = "plugin"
			}
			return true
		})
	}
	return surface
}
func walk(b *ast.BlockStmt, prefix, surface string) {
	surface = authFor(b, surface)
	ast.Inspect(b, func(n ast.Node) bool {
		c, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if i, ok := c.Fun.(*ast.Ident); ok {
			if helper := routerHelpers[i.Name]; helper != nil {
				walk(helper.Body, prefix, surface)
				return false
			}
		}
		sel, ok := c.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if sel.Sel.Name == "Route" || sel.Sel.Name == "Group" {
			next := prefix
			if sel.Sel.Name == "Route" {
				p := str(c.Args[0])
				if p == "" {
					return false
				}
				next = strings.TrimRight(prefix, "/") + p
			}
			for _, a := range c.Args {
				if fn, ok := a.(*ast.FuncLit); ok {
					walk(fn.Body, next, surface)
				}
			}
			return false
		}
		method := strings.ToUpper(sel.Sel.Name)
		if !strings.Contains("|GET|POST|PUT|PATCH|DELETE|HEAD|", "|"+method+"|") || len(c.Args) < 2 {
			return true
		}
		p := str(c.Args[0])
		if p == "" {
			return true
		}
		path := strings.TrimRight(prefix, "/") + p
		path = strings.TrimRight(path, "/")
		h, ok := c.Args[1].(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if surface != "pat" || !strings.HasPrefix(path, "/api/") {
			catalog.Excluded = append(catalog.Excluded, map[string]string{"method": method, "path": path, "reason": "Non-user API surface: " + surface})
			return false
		}
		name := h.Sel.Name
		read := method == "GET" || method == "HEAD"
		// GET routes that initiate a connection or refresh remote state are writes.
		if strings.Contains(name, "Connect") && !strings.HasPrefix(name, "List") && !strings.HasPrefix(name, "Get") {
			read = false
		}
		catalog.Operations = append(catalog.Operations, Operation{Name: "multica_api_" + snake(name), Handler: name, Method: method, Path: path, ReadOnly: read, Input: schema(name, path, method), Source: fmt.Sprintf("server/cmd/server/router.go:%d", fset.Position(c.Pos()).Line)})
		return false
	})
}
func main() {
	if len(os.Args) != 4 {
		panic("usage: generate_api SOURCE_ROOT REVISION OUTPUT_JSON")
	}
	root := os.Args[1]
	catalog.Revision = os.Args[2]
	files, err := filepath.Glob(filepath.Join(root, "server/internal/handler/*.go"))
	if err != nil {
		panic(err)
	}
	for _, p := range files {
		if strings.HasSuffix(p, "_test.go") {
			continue
		}
		f, e := parser.ParseFile(fset, p, nil, 0)
		if e != nil {
			panic(e)
		}
		for _, d := range f.Decls {
			switch x := d.(type) {
			case *ast.GenDecl:
				for _, s := range x.Specs {
					if t, ok := s.(*ast.TypeSpec); ok {
						types[t.Name.Name] = t.Type
					}
				}
			case *ast.FuncDecl:
				funcs[x.Name.Name] = x
			}
		}
	}
	f, e := parser.ParseFile(fset, filepath.Join(root, "server/cmd/server/router.go"), nil, 0)
	if e != nil {
		panic(e)
	}
	public, e := parser.ParseFile(fset, filepath.Join(root, "server/pkg/publicapi/v1/routes.go"), nil, 0)
	if e != nil {
		panic(e)
	}
	for label, file := range map[string]*ast.File{"": f, "publicapiv1.": public} {
		ast.Inspect(file, func(n ast.Node) bool {
			if v, ok := n.(*ast.ValueSpec); ok {
				for i, name := range v.Names {
					if i < len(v.Values) {
						if s := str(v.Values[i]); s != "" {
							constants[label+name.Name] = s
						}
					}
				}
			}
			return true
		})
	}
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == "registerPluginActionRoutes" {
			routerHelpers[fn.Name.Name] = fn
		}
	}
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == "NewRouterWithOptions" {
			walk(fn.Body, "", "public")
		}
	}
	sort.Slice(catalog.Operations, func(i, j int) bool {
		a, b := catalog.Operations[i], catalog.Operations[j]
		return a.Path+a.Method < b.Path+b.Method
	})
	counts := map[string]int{}
	for _, op := range catalog.Operations {
		counts[op.Name]++
	}
	used := map[string]bool{}
	for i := range catalog.Operations {
		op := &catalog.Operations[i]
		if counts[op.Name] > 1 {
			op.Name += "_" + strings.ToLower(op.Method)
		}
		if used[op.Name] || len(op.Name) > 64 {
			h := sha256.Sum256([]byte(op.Method + op.Path))
			if len(op.Name) > 55 {
				op.Name = op.Name[:55]
			}
			op.Name += fmt.Sprintf("_%x", h[:4])
		}
		if used[op.Name] {
			panic("duplicate operation " + op.Name)
		}
		used[op.Name] = true
	}
	data, e := json.MarshalIndent(catalog, "", "  ")
	if e != nil {
		panic(e)
	}
	if e = os.WriteFile(os.Args[3], append(data, '\n'), 0644); e != nil {
		panic(e)
	}
	fmt.Printf("Generated %d PAT operations; cataloged %d non-user routes\n", len(catalog.Operations), len(catalog.Excluded))
}
