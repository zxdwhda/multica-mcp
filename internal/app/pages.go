package app

import (
	"context"
	"multica-mcp/internal/domain"
	"multica-mcp/internal/multica"
)

func (u *UseCase) ProjectsPage(ctx context.Context, q string, limit, offset int) (*multica.Page[domain.Project], error) {
	return u.client.ProjectsPage(ctx, q, limit, offset)
}
func (u *UseCase) TasksPage(ctx context.Context, in domain.ListTasksInput, limit, offset int) (*multica.Page[domain.Task], error) {
	return u.client.TasksPage(ctx, in, limit, offset)
}
func (u *UseCase) Statuses(ctx context.Context) (any, error) { return u.client.Statuses(ctx) }
