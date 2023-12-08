package service

import (
	"context"
	"testing"

	"github.com/niluwats/task-service/internal/domain"
	repos "github.com/niluwats/task-service/internal/mocks/repos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreate(t *testing.T) {
	repo := new(repos.ProjectRepository)
	service := NewProjectServiceImpl(repo)

	expectedProject := domain.Project{
		Name:        "project1",
		Description: "this is the project1",
		Creator:     1,
		Tasks:       make([]domain.Task, 0),
		Assignees:   make([]int32, 0),
	}

	project := domain.Project{
		Name:        "project1",
		Description: "this is the project1",
		Creator:     1,
		Tasks:       make([]domain.Task, 0),
		Assignees:   make([]int32, 0),
	}

	repo.On("Insert", mock.Anything, mock.AnythingOfType("domain.Project")).Return(&expectedProject, nil).Once()

	projectResponse, err := service.Create(context.Background(), project)
	assert.Nil(t, err)
	assert.Equal(t, &project, projectResponse)
}

func TestUpdate(t *testing.T) {
	repo := new(repos.ProjectRepository)
	service := NewProjectServiceImpl(repo)

	expectedProject := domain.Project{
		Name:        "project101",
		Description: "description",
		Creator:     1,
	}

	project := domain.Project{
		Name:        "project101",
		Description: "description",
		Creator:     1,
	}

	repo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Project")).Return(&expectedProject, nil).Once()
	updatedProject, err := service.Update(context.Background(), project)
	assert.Nil(t, err)
	assert.Equal(t, &project, updatedProject)
}

func TestRemove(t *testing.T) {
	repo := new(repos.ProjectRepository)
	service := NewProjectServiceImpl(repo)

	projectId := primitive.NewObjectID()

	repo.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
	err := service.Remove(context.Background(), projectId.String())
	assert.Nil(t, err)
}

func TestViewByID(t *testing.T) {
	repo := new(repos.ProjectRepository)
	service := NewProjectServiceImpl(repo)

	projectId := primitive.NewObjectID()

	expectedProject := domain.Project{
		ID:          projectId,
		Name:        "test1",
		Description: "description",
		Creator:     1,
	}

	repo.On("FindByID", mock.Anything, mock.Anything).Return(&expectedProject, nil).Once()
	project, err := service.ViewByID(context.Background(), projectId.String())
	assert.Nil(t, err)
	assert.Equal(t, &expectedProject, project)
}

func TestViewAll(t *testing.T) {
	repo := new(repos.ProjectRepository)
	service := NewProjectServiceImpl(repo)

	projectId1 := primitive.NewObjectID()
	projectId2 := primitive.NewObjectID()

	expectedProject1 := domain.Project{
		ID:          projectId1,
		Name:        "test1",
		Description: "description",
		Creator:     1,
	}

	expectedProject2 := domain.Project{
		ID:          projectId2,
		Name:        "test2",
		Description: "description",
		Creator:     2,
	}

	repo.On("FindAll", mock.Anything).Return([]domain.Project{expectedProject1, expectedProject2}, nil).Once()
	projects, err := service.ViewAll(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, []domain.Project{expectedProject1, expectedProject2}, projects)
}
