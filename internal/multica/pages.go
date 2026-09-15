package multica

import (
	"context"
	"fmt"
	"multica-mcp/internal/domain"
	"net/url"
	"strconv"
)

type Page[T any] struct {
	Items          []T  `json:"items"`
	SourceTotal    int  `json:"source_total"`
	Offset         int  `json:"offset"`
	NextOffset     *int `json:"next_offset"`
	HasMore        bool `json:"has_more"`
	FilteredOnPage bool `json:"filtered_on_page"`
}

func page[T any](items []T, total, offset, count int) *Page[T] {
	if items == nil {
		items = []T{}
	}
	p := &Page[T]{Items: items, SourceTotal: total, Offset: offset, HasMore: offset+count < total}
	if p.HasMore && count > 0 {
		n := offset + count
		p.NextOffset = &n
	}
	return p
}
func paging(limit, offset, max int) (int, error) {
	if limit == 0 {
		limit = max
	}
	if limit < 1 || limit > max || offset < 0 {
		return 0, fmt.Errorf("limit must be 1..%d and offset >=0", max)
	}
	return limit, nil
}
func (c *Client) ProjectsPage(ctx context.Context, query string, limit, offset int) (*Page[domain.Project], error) {
	max := 100
	if query != "" {
		max = 50
	}
	limit, err := paging(limit, offset, max)
	if err != nil {
		return nil, err
	}
	q := url.Values{"limit": {strconv.Itoa(limit)}, "offset": {strconv.Itoa(offset)}}
	path := "/api/projects"
	if query != "" {
		path += "/search"
		q.Set("q", query)
	}
	var r listProjectsResponse
	if err = c.doGet(ctx, path+"?"+q.Encode(), &r, true); err != nil {
		return nil, err
	}
	return page(r.Projects, r.Total, offset, len(r.Projects)), nil
}
func (c *Client) TasksPage(ctx context.Context, in domain.ListTasksInput, limit, offset int) (*Page[domain.Task], error) {
	search := in.Query != nil && *in.Query != ""
	max := 100
	if search {
		max = 50
	}
	limit, err := paging(limit, offset, max)
	if err != nil {
		return nil, err
	}
	q := url.Values{"limit": {strconv.Itoa(limit)}, "offset": {strconv.Itoa(offset)}}
	path := "/api/issues"
	if search {
		path += "/search"
		q.Set("q", *in.Query)
		q.Set("include_closed", "true")
	} else {
		if in.ProjectID != "" {
			q.Set("project_id", in.ProjectID)
		}
		if in.Status != nil {
			q.Set("status", *in.Status)
		}
		if in.Assignee != nil {
			q.Set("assignee_id", *in.Assignee)
		}
	}
	var r listIssuesResponse
	if err = c.doGet(ctx, path+"?"+q.Encode(), &r, true); err != nil {
		return nil, err
	}
	p := page(r.Issues, r.Total, offset, len(r.Issues))
	if search && (in.ProjectID != "" || in.Status != nil || in.Assignee != nil) {
		p.FilteredOnPage = true
		p.Items = []domain.Task{}
		for _, t := range r.Issues {
			if in.ProjectID != "" && (t.ProjectID == nil || *t.ProjectID != in.ProjectID) {
				continue
			}
			if in.Status != nil && t.Status != *in.Status {
				continue
			}
			if in.Assignee != nil && (t.AssigneeID == nil || *t.AssigneeID != *in.Assignee) {
				continue
			}
			p.Items = append(p.Items, t)
		}
	}
	return p, nil
}
func (c *Client) Statuses(ctx context.Context) (any, error) {
	var r any
	err := c.doGet(ctx, "/api/issue-statuses", &r, true)
	return r, err
}
