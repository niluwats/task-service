package repositories

import (
	"context"
	"errors"

	"github.com/niluwats/task-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepoDb struct {
	dbCollection *mongo.Collection
}

func NewTaskRepoDb(dbClient *mongo.Collection) TaskRepoDb {
	return TaskRepoDb{dbClient}
}

func (repo TaskRepoDb) Insert(ctx context.Context, projectID string, task domain.Task) (*domain.Task, error) {
	task.ID = primitive.NewObjectID()
	obProjectID, _ := primitive.ObjectIDFromHex(projectID)

	var project domain.Project
	filter := bson.M{"_id": obProjectID}
	update := bson.M{"$push": bson.M{"tasks": task}}

	if err := repo.dbCollection.FindOneAndUpdate(ctx, filter, update).
		Decode(&project); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("project ID not found")
		}
		return nil, err
	}

	return &project.Tasks[len(project.Tasks)-1], nil
}

func (repo TaskRepoDb) Update(ctx context.Context, projectID string, taskID string, task domain.Task) error {
	obProjectID, _ := primitive.ObjectIDFromHex(projectID)
	obtaskID, _ := primitive.ObjectIDFromHex(taskID)

	update := bson.M{
		"$set": bson.M{
			"tasks.$[elem].description": task.Description,
			"tasks.$[elem].assignee":    task.Assignee,
			"tasks.$[elem].updated_at":  task.UpdatedAt,
			"tasks.$[elem].status":      task.TaskStatus,
		},
	}

	filter := bson.M{"_id": obProjectID}
	arrayFilters := options.ArrayFilters{Filters: bson.A{bson.M{"elem._id": obtaskID}}}
	opts := options.FindOneAndUpdate().SetArrayFilters(arrayFilters).SetReturnDocument((1))

	var project domain.Project
	err := repo.dbCollection.FindOneAndUpdate(ctx, update, filter, opts).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("task not found")
		}
		return err
	}
	return nil
}

func (repo TaskRepoDb) Delete(ctx context.Context, projectID, taskID string) error {
	obProjectID, _ := primitive.ObjectIDFromHex(projectID)
	obTaskID, _ := primitive.ObjectIDFromHex(taskID)

	filter := bson.M{"_id": obProjectID}
	update := bson.M{"$pull": bson.M{"tasks": bson.M{"_id": obTaskID}}}

	result, err := repo.dbCollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo TaskRepoDb) FindByID(ctx context.Context, projectID, taskID string) (*domain.Task, error) {
	obProjectID, _ := primitive.ObjectIDFromHex(projectID)
	obTaskID, _ := primitive.ObjectIDFromHex(taskID)

	filter := bson.M{"_id": obProjectID, "tasks._id": obTaskID}
	findOptions := options.FindOne().SetProjection(bson.M{"tasks.$": 1})

	var project domain.Project
	if err := repo.dbCollection.FindOne(ctx, filter, findOptions).Decode(&project); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	if len(project.Tasks) > 0 {
		return &project.Tasks[0], nil
	} else {
		return nil, errors.New("task not found in the specified project")
	}
}

func (repo TaskRepoDb) FindAllByProjectID(ctx context.Context, projectID string) ([]domain.Task, error) {
	obProjectID, _ := primitive.ObjectIDFromHex(projectID)

	filter := bson.M{"_id": obProjectID}

	var project domain.Project
	if err := repo.dbCollection.FindOne(ctx, filter).Decode(&project); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	return project.Tasks, nil
}
