package config

import (
	"context"
	"log"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var MongoDB *mongo.Database
var mc *mongo.Client

func InitMongo() error {
	log.Println("[MONGO] Initializing mongo client...")

	var uri string
	if Envs.Mongo.Username != "" && Envs.Mongo.Password != "" {
		uri = fmt.Sprintf(
			"mongodb://%s:%s@%s:%s",
			Envs.Mongo.Username,
			Envs.Mongo.Password,
			Envs.Mongo.Host,
			Envs.Mongo.Port,
		)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", Envs.Mongo.Host, Envs.Mongo.Port)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Println("[MONGO] Pinging mongo server...")
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	MongoDB = client.Database(Envs.Mongo.Database)
	mc = client
	log.Println("[MONGO] Mongo client initialized successfully.")
	return nil
}

func CloseMongo() {
	if err := mc.Disconnect(context.Background()); err != nil {
		log.Println("[MONGO] Error closing mongo client: ", err)
	}
}
