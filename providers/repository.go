package providers

import (
	"context"
	"time"

	"github.com/golang/glog"
	apimongo "github.com/pruh/api/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository to access providers store
type Repository struct {
	mongo *mongo.Client
}

const (
	dbName         = "ApiDB"
	collectionName = "providers"
)

// GetAll returns all providers
func (r *Repository) GetAll() ([]Provider, error) {
	glog.Info("Querying for all providers")
	collection := r.mongo.Database(dbName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		glog.Errorf("Error while performing DB find. %s", err)
		return nil, err
	}
	defer cur.Close(ctx)

	providers := []Provider{}
	for cur.Next(ctx) {
		var prov Provider
		err := cur.Decode(&prov)
		if err != nil {
			glog.Errorf("Cannot decode provider. %s", err)
			return nil, err
		}
		providers = append(providers, prov)
	}
	if err := cur.Err(); err != nil {
		glog.Errorf("Cursor has some error. %s", err)
		return nil, err
	}
	return providers, nil
}

// GetOne returns provider by ID or nil
func (r *Repository) GetOne(uuid apimongo.UUID) (*Provider, error) {
	glog.Infof("Querying for provider with UUID: %s", uuid)

	collection := r.mongo.Database(dbName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var provider Provider
	err := collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		glog.Errorf("Cannot decode provider. %s", err)
		return nil, err
	}

	return &provider, nil
}

// CreateOne creates new provider for specified params
func (r *Repository) CreateOne(provider Provider) bool {
	glog.Infof("Creating new provider: %+v", provider)

	collection := r.mongo.Database(dbName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, provider)
	if err != nil {
		glog.Errorf("Failed to insert provider. %s", err)
		return false
	}

	return true
}

// DeleteOne deletes provider with ID and returns true if record was removed
func (r *Repository) DeleteOne(uuid apimongo.UUID) (bool, error) {
	glog.Infof("Deleting provider with UUID: %s", uuid)

	collection := r.mongo.Database(dbName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": uuid})
	if err != nil {
		glog.Errorf("Failed to delete provider. %s", err)
		return false, err
	}

	if res.DeletedCount == 0 {
		return false, nil
	}

	return true, nil
}

// DeleteAll deletes providers with IDs and returns true if any record was removed
func (r *Repository) DeleteAll(uuids []apimongo.UUID) (bool, error) {
	glog.Infof("Deleting providers with UUID: %v", uuids)

	collection := r.mongo.Database(dbName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": uuids}}
	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		glog.Error("Failed to delete provider. ", err)
		return false, err
	}

	if res.DeletedCount == 0 {
		return false, nil
	}

	return true, nil
}
