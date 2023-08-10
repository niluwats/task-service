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

type ProjectRepoDb struct {
	dbCollection *mongo.Collection
}

func NewProjectRepoDb(dbClient *mongo.Collection) ProjectRepoDb {
	return ProjectRepoDb{dbClient}
}

func (repo ProjectRepoDb) Insert(ctx context.Context, project domain.Project) (*domain.Project, error) {
	result, err := repo.dbCollection.InsertOne(ctx, project)
	if err != nil {
		er, ok := err.(mongo.WriteException)
		if ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("project with this name already exists")
		}
		return nil, err
	}
	ID := result.InsertedID

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"name": 1}, Options: opt}
	if _, err := repo.dbCollection.Indexes().CreateOne(ctx, index); err != nil {
		return nil, errors.New("error creating index")
	}

	var newProject domain.Project
	if err = repo.dbCollection.FindOne(ctx, bson.M{"_id": ID}).Decode(&newProject); err != nil {
		return nil, err
	}

	return &newProject, nil
}

func (repo ProjectRepoDb) Update(ctx context.Context, ID string, project domain.Project) (*domain.Project, error) {
	obID, _ := primitive.ObjectIDFromHex(ID)
	var updatedDoc domain.Project

	filter := bson.D{{Key: "_id", Value: obID}}
	update := bson.D{{Key: "$set", Value: project}}

	if err := repo.dbCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1)).Decode(&updatedDoc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ID doesn't exists")
		}
		return nil, err
	}
	return &updatedDoc, nil
}

func (repo ProjectRepoDb) Delete(ctx context.Context, ID string) error {
	obID, _ := primitive.ObjectIDFromHex(ID)

	res, err := repo.dbCollection.DeleteOne(ctx, bson.M{"_id": obID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("ID doesn't exists")
	}
	return nil
}

func (repo ProjectRepoDb) FindByID(ctx context.Context, ID string) (*domain.Project, error) {
	obID, _ := primitive.ObjectIDFromHex(ID)

	var project domain.Project

	err := repo.dbCollection.FindOne(ctx, bson.M{"_id": obID}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document")
		}
		return nil, err
	}

	return &project, nil
}

func (repo ProjectRepoDb) FindAll(ctx context.Context) ([]domain.Project, error) {
	var projects []domain.Project
	cursor, err := repo.dbCollection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		project := domain.Project{}
		if err := cursor.Decode(project); err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}
	return projects, nil
}
