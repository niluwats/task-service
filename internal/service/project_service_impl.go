package service

import (
	"context"
	"time"

	"github.com/niluwats/task-service/api/pb"
	"github.com/niluwats/task-service/internal/domain"
	"github.com/niluwats/task-service/internal/repositories"
	"google.golang.org/grpc/codes"
)

type ProjectServiceImpl struct {
	repo repositories.ProjectRepository
}

func NewProjectServiceImpl(repo repositories.ProjectRepository) ProjectServiceImpl {
	return ProjectServiceImpl{repo}
}

func (serv ProjectServiceImpl) Create(ctx context.Context, project domain.Project) (*pb.ProjectResponse, error) {
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	createdProject, err := serv.repo.Insert(ctx, project)
	if err != nil {
		return nil, err
	}

	return convertToPbProject(createdProject), nil
}

func (serv ProjectServiceImpl) Update(ctx context.Context, project domain.Project) (*pb.ProjectResponse, error) {
	project.UpdatedAt = time.Now()

	updatedProject, err := serv.repo.Update(ctx, project.ID.String(), project)
	if err != nil {
		return nil, err
	}

	return convertToPbProject(updatedProject), nil
}

func (serv ProjectServiceImpl) ViewByID(ctx context.Context, projectID string) (*pb.ProjectResponse, error) {
	project, err := serv.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return convertToPbProject(project), nil
}

func (serv ProjectServiceImpl) ViewAll(ctx context.Context) (*pb.ProjectsResponse, error) {
	projects, err := serv.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return convertToPbProjects(projects), nil
}

func (serv ProjectServiceImpl) Remove(ctx context.Context, projectId string) error {
	err := serv.repo.Delete(ctx, projectId)
	if err != nil {
		return err
	}
	return nil
}

func convertToPbProject(project *domain.Project) *pb.ProjectResponse {
	return &pb.ProjectResponse{
		Project: &pb.Project{
			Id:          project.ID.String(),
			Name:        project.Name,
			Description: project.Description,
			Creator:     project.Creator,
			CreatedOn:   project.CreatedAt.String(),
			UpdatedOn:   project.UpdatedAt.String(),
			Assignees:   project.Assignees,
			Tasks:       convertTaskArray(project.Tasks),
		},
	}
}

func convertToPbProjects(projects []domain.Project) *pb.ProjectsResponse {
	pbProjects := make([]*pb.Project, 0)
	for _, v := range projects {
		project := pb.Project{
			Id:          v.ID.String(),
			Name:        v.Name,
			Description: v.Description,
			Creator:     v.Creator,
			Assignees:   v.Assignees,
			CreatedOn:   v.CreatedAt.String(),
			UpdatedOn:   v.UpdatedAt.String(),
			Tasks:       convertTaskArray(v.Tasks),
		}
		pbProjects = append(pbProjects, &project)
	}
	return &pb.ProjectsResponse{
		Projects: pbProjects,
		CommonResponse: &pb.CommonResponse{
			Status:  int32(codes.OK),
			Message: "success",
		},
	}
}

func convertTaskArray(task []domain.Task) []*pb.Task {
	newTasks := make([]*pb.Task, 0)
	for _, v := range task {
		newTask := pb.Task{
			Id:          v.ID.String(),
			Description: v.Description,
			Creator:     v.Creator,
			Assignee:    v.Assignee,
			TaskStatus:  pb.Status(v.TaskStatus),
		}
		newTasks = append(newTasks, &newTask)
	}
	return newTasks
}
