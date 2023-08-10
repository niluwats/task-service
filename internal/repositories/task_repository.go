package repositories

import (
	"context"

	"github.com/niluwats/task-service/internal/domain"
)

type TaskRepository interface {
	Insert(ctx context.Context, projectID string, task domain.Task) (*domain.Task, error)
	Update(ctx context.Context, projectID string, taskID string, task domain.Task) (*domain.Task, error)
	Delete(ctx context.Context, projectID, taskID string) error
	FindByID(ctx context.Context, projectID, taskID string) (*domain.Task, error)
	FindAllByProjectID(ctx context.Context, projectID string) ([]domain.Task, error)
}
