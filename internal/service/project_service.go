package service

import (
	"context"

	"github.com/niluwats/task-service/internal/domain"
)

type ProjectService interface {
	Create(ctx context.Context, project domain.Project) (*domain.Project, error)
	Update(ctx context.Context, project domain.Project) error
	ViewByID(ctx context.Context, projectID string) (*domain.Project, error)
	ViewAll(ctx context.Context) []domain.Project
}
