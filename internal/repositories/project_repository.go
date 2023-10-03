package repositories

import (
	"context"

	"github.com/niluwats/task-service/internal/domain"
)

//go:generate mockery --name ProjectRepoMock
type ProjectRepository interface {
	Insert(ctx context.Context, project domain.Project) (*domain.Project, error)
	Update(ctx context.Context, ID string, project domain.Project) (*domain.Project, error)
	Delete(ctx context.Context, ID string) error
	FindByID(ctx context.Context, ID string) (*domain.Project, error)
	FindAll(ctx context.Context) ([]domain.Project, error)
}
