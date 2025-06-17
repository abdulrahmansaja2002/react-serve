package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo interface {
	DoSomething(ctx context.Context)
}

type repo struct {
	db *mongo.Database
}

func NewRepo(db *mongo.Database) Repo {
	return &repo{db: db}
}

func (r *repo) DoSomething(ctx context.Context) {}
