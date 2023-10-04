package service

import (
	"context"
	"time"

	"github.com/niluwats/task-service/internal/domain"
	"github.com/niluwats/task-service/internal/repositories"
)

type ProjectServiceImpl struct {
	repo repositories.ProjectRepository
}

func NewProjectServiceImpl(repo repositories.ProjectRepository) ProjectServiceImpl {
	return ProjectServiceImpl{repo}
}

func (serv ProjectServiceImpl) Create(ctx context.Context, project domain.Project) (*domain.Project, error) {
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	createdProject, err := serv.repo.Insert(ctx, project)
	if err != nil {
		return nil, err
	}
	return createdProject, nil
}

func (serv ProjectServiceImpl) Update(ctx context.Context, project domain.Project) (*domain.Project, error) {
	project.UpdatedAt = time.Now()

	updatedProject, err := serv.repo.Update(ctx, project.ID.String(), project)
	if err != nil {
		return nil, err
	}

	return updatedProject, nil
}

func (serv ProjectServiceImpl) ViewByID(ctx context.Context, projectID string) (*domain.Project, error) {
	project, err := serv.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (serv ProjectServiceImpl) ViewAll(ctx context.Context) ([]domain.Project, error) {
	projects, err := serv.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (serv ProjectServiceImpl) Remove(ctx context.Context, projectId string) error {
	err := serv.repo.Delete(ctx, projectId)
	if err != nil {
		return err
	}

	return nil
}
