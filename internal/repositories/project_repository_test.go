package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/niluwats/task-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestInsert(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		collection := mt.Coll
		id := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		assignees := make([]int32, 0)
		assignees = append(assignees, 2)
		assignees = append(assignees, 3)

		project := domain.Project{
			ID:          id,
			Name:        "test1",
			Description: "test project",
			Creator:     1,
			Assignees:   assignees,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		projectRepo := NewProjectRepoDb(collection)
		insertedProject, err := projectRepo.Insert(context.Background(), project)

		assert.Nil(t, err)
		assert.Equal(t, &project, insertedProject)
	})
}

func TestUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		collection := mt.Coll

		project := domain.Project{
			ID:          primitive.NewObjectID(),
			Name:        "test101",
			Description: "test project 101",
			Creator:     1,
		}

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: project.ID},
				{Key: "name", Value: project.Name},
				{Key: "description", Value: project.Description},
				{Key: "creator", Value: project.Creator},
			}},
		})

		projectRepo := NewProjectRepoDb(collection)
		updatedProject, err := projectRepo.Update(context.Background(), project.ID.String(), project)

		assert.Nil(t, err)
		assert.Equal(t, &project, updatedProject)
	})
}

func TestDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})
		collection := mt.Coll
		projectRepo := NewProjectRepoDb(collection)
		err := projectRepo.Delete(context.Background(), primitive.NewObjectID().String())
		assert.Nil(t, err)
	})

	mt.Run("no document deleted", func(mt *mtest.T) {
		collection := mt.Coll
		projectRepo := NewProjectRepoDb(collection)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 0}})
		err := projectRepo.Delete(context.Background(), primitive.NewObjectID().String())
		assert.NotNil(t, err)
	})
}

func TestFindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		collection := mt.Coll

		expectedProject := domain.Project{
			ID:          primitive.NewObjectID(),
			Name:        "test1",
			Description: "test1 description",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedProject.ID},
			{Key: "name", Value: expectedProject.Name},
			{Key: "description", Value: expectedProject.Description},
		}))

		projectRepo := NewProjectRepoDb(collection)
		projectResponse, err := projectRepo.FindByID(context.Background(), expectedProject.ID.String())
		assert.Nil(t, err)
		assert.Equal(t, &expectedProject, projectResponse)
	})
}

func TestFindAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		collection := mt.Coll

		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: id1},
			{Key: "name", Value: "test1"},
			{Key: "description", Value: "project test1"},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: id2},
			{Key: "name", Value: "test2"},
			{Key: "description", Value: "project test2"},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		projectRepo := NewProjectRepoDb(collection)
		projects, err := projectRepo.FindAll(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, []domain.Project{
			{ID: id1, Name: "test1", Description: "project test1"},
			{ID: id2, Name: "test2", Description: "project test2"},
		}, projects)
	})
}
