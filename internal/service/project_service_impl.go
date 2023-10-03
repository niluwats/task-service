package service

import (
	"context"

	"github.com/niluwats/task-service/internal/repositories"
)

type ProjectServiceImpl struct {
	repo repositories.ProjectRepository
}

func NewProjectServiceImpl(repo repositories.ProjectRepository) ProjectServiceImpl {
	return ProjectServiceImpl{repo}
}

func (serv ProjectServiceImpl) NewProject(ctx context.Context) {

}
