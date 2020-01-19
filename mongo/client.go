package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/pruh/api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient creates new client for mongodb
func NewClient(config *config.Configuration) *mongo.Client {
	var uri string
	if config.MongoUsername != nil && config.MongoPassword != nil {
		uri = fmt.Sprintf("mongodb://%s:%s@mongo:27017", *config.MongoUsername, *config.MongoPassword)
	} else {
		uri = "mongodb://mongo:27017"
	}
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		glog.Fatalf("Cannot connect to databse. %s", err)
	}

	return client
}
