package multica

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"multica-mcp/internal/apicatalog"
)

// CallAPI executes only a generated operation against the configured origin.
// It never follows redirects or accepts caller-supplied credentials/hosts.
func (c *Client) CallAPI(ctx context.Context, op apicatalog.Operation, in apicatalog.Input) (*apicatalog.Response, error) {
	path, err := op.URL(in)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	contentType := "application/json"
	if len(in.Files) > 0 {
		writer := multipart.NewWriter(&buf)
		fields, ok := in.Body.(map[string]any)
		if in.Body != nil && !ok {
			return nil, fmt.Errorf("multipart body must be an object")
		}
		for k, v := range fields {
			value, ok := v.(string)
			if !ok {
				b, e := json.Marshal(v)
				if e != nil {
					return nil, e
				}
				value = string(b)
			}
			if e := writer.WriteField(k, value); e != nil {
				return nil, e
			}
		}
		for _, file := range in.Files {
			data, e := base64.StdEncoding.DecodeString(file.Data)
			if e != nil {
				return nil, fmt.Errorf("invalid base64 upload")
			}
			if len(data) > 4*1024*1024 {
				return nil, fmt.Errorf("upload exceeds 4 MiB per file")
			}
			if strings.ContainsAny(file.Field+file.Filename+file.ContentType, "\r\n") {
				return nil, fmt.Errorf("invalid multipart metadata")
			}
			header := textproto.MIMEHeader{}
			esc := strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
			header.Set("Content-Disposition", `form-data; name="`+esc.Replace(file.Field)+`"; filename="`+esc.Replace(file.Filename)+`"`)
			if file.ContentType != "" {
				header.Set("Content-Type", file.ContentType)
			}
			part, e := writer.CreatePart(header)
			if e != nil {
				return nil, e
			}
			if _, e = part.Write(data); e != nil {
				return nil, e
			}
		}
		if err = writer.Close(); err != nil {
			return nil, err
		}
		contentType = writer.FormDataContentType()
	} else if in.Body != nil {
		if err = json.NewEncoder(&buf).Encode(in.Body); err != nil {
			return nil, err
		}
	}
	if buf.Len() > 8*1024*1024 {
		return nil, fmt.Errorf("request exceeds 8 MiB")
	}
	req, err := http.NewRequestWithContext(ctx, op.Method, c.baseURL+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	for k, v := range in.Headers {
		switch http.CanonicalHeaderKey(k) {
		case "If-Match", "Idempotency-Key", "Range", "Accept", "X-Multica-Plugin-Installation":
			req.Header.Set(k, v)
		default:
			return nil, fmt.Errorf("unsupported request header %s", k)
		}
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-Client-Platform", "mcp")
	req.Header.Set("X-Client-Version", c.clientVersion)
	if err = c.attachWorkspaceHeaders(req); err != nil {
		return nil, err
	}
	client := *c.httpClient
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upstream request failed")
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024+1))
	if err != nil {
		return nil, fmt.Errorf("upstream response interrupted")
	}
	out := &apicatalog.Response{Status: resp.StatusCode, ContentType: resp.Header.Get("Content-Type"), Headers: map[string]string{}}
	for _, h := range []string{"ETag", "Location", "Retry-After", "Content-Range"} {
		if v := resp.Header.Get(h); v != "" {
			out.Headers[h] = v
		}
	}
	if len(raw) > 8*1024*1024 {
		return nil, fmt.Errorf("response exceeds 8 MiB; use pagination or Range")
	}
	if len(raw) == 0 {
		return out, nil
	}
	if json.Valid(raw) {
		if err = json.Unmarshal(raw, &out.Body); err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(out.ContentType, "text/") {
		out.Body = string(raw)
	} else {
		out.DataBase64 = base64.StdEncoding.EncodeToString(raw)
	}
	return out, nil
}
