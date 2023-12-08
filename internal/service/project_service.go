package service

import (
	"context"

	"github.com/niluwats/task-service/api/pb"
	"github.com/niluwats/task-service/internal/domain"
)

type ProjectService interface {
	Create(ctx context.Context, project domain.Project) (*pb.ProjectResponse, error)
	Update(ctx context.Context, project domain.Project) (*pb.ProjectResponse, error)
	ViewByID(ctx context.Context, projectId string) (*pb.ProjectResponse, error)
	ViewAll(ctx context.Context) ([]pb.ProjectsResponse, error)
	Remove(ctx context.Context, projectId string) error
}
