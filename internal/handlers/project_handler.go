package handlers

import (
	"context"

	"github.com/niluwats/task-service/api/pb"
	"github.com/niluwats/task-service/internal/domain"
	"github.com/niluwats/task-service/internal/errors"
	"github.com/niluwats/task-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectHandlers interface {
	CreateProject(ctx context.Context, req *pb.NewProjectRequest) (*pb.ProjectResponse, error)
	UpdateProject(ctx context.Context, req *pb.UpdaterojectRequest) (*pb.ProjectResponse, error)
	RemoveProject(ctx context.Context) (*pb.CommonResponse, error)
	ViewProject(ctx context.Context, req *pb.ViewProjectRequest) (*pb.ProjectResponse, error)
	ViewAllProjects(ctx context.Context) (*pb.ProjectsResponse, error)
}

type ProjectHandlersImpl struct {
	pb.UnimplementedTaskServiceServer
	service service.ProjectService
}

func NewProjectHandlerImpl(service service.ProjectService) ProjectHandlersImpl {
	return ProjectHandlersImpl{service: service}
}

func (h ProjectHandlersImpl) CreateProject(ctx context.Context, req *pb.NewProjectRequest) (*pb.ProjectResponse, error) {
	if req.Name == "" || req.Description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name or Description can't be empty")
	}

	res, err := h.service.Create(ctx, domain.Project{
		Name:        req.Name,
		Description: req.Description,
		Assignees:   req.Assignees,
		Creator:     req.Creator,
	})

	if err != nil {
		return nil, status.Errorf(getStatusCode(err), err.Error())
	}

	return res, nil
}

func getStatusCode(err error) codes.Code {
	switch err.(type) {
	case *errors.BadRequest:
		return codes.InvalidArgument
	case *errors.Unauthorized:
		return codes.Unauthenticated
	case *errors.ConflictError:
		return codes.AlreadyExists
	case *errors.NotFoundError:
		return codes.NotFound
	default:
		return codes.Internal
	}
}
