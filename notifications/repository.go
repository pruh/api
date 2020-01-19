package notifications

import (
	"context"
	"time"

	"github.com/golang/glog"
	"github.com/pruh/api/config"
	apimongo "github.com/pruh/api/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository accesses notifications store
type Repository struct {
	Mongo  *mongo.Client
	Config *config.Configuration
}

const (
	dbName              = "ApiDB"
	notifCollectionName = "notifications"
)

// GetAll returns all notifications
func (r *Repository) GetAll() ([]Notification, error) {
	glog.Info("Querying for all notifications")
	collection := r.Mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		glog.Errorf("Error while performing DB find. %s", err)
		return nil, err
	}
	defer cur.Close(ctx)

	notifs := []Notification{}
	for cur.Next(ctx) {
		var notif Notification
		err := cur.Decode(&notif)
		if err != nil {
			glog.Errorf("Cannot decode notification. %s", err)
			return nil, err
		}
		notifs = append(notifs, notif)
	}
	if err := cur.Err(); err != nil {
		glog.Errorf("Cursor has some error. %s", err)
		return nil, err
	}
	return notifs, nil
}

// GetOne returns notifications by ID or nil
func (r *Repository) GetOne(uuid apimongo.UUID) (*Notification, error) {
	glog.Infof("Querying for notification with UUID: %s", uuid)

	collection := r.Mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var notif Notification
	err := collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&notif)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		glog.Errorf("Cannot decode notification. %s", err)
		return nil, err
	}

	return &notif, nil
}

// CreateOne creates new notification for specified params
func (r *Repository) CreateOne(notification Notification) bool {
	glog.Infof("Creating new notification: %+v", notification)

	collection := r.Mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, notification)
	if err != nil {
		glog.Errorf("Failed to insert notification. %s", err)
		return false
	}

	return true
}

// DeleteOne deletes notifications with ID and returns true if record was removed
func (r *Repository) DeleteOne(uuid apimongo.UUID) (bool, error) {
	glog.Infof("Deleting notification with UUID: %s", uuid)

	collection := r.Mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": uuid})
	if err != nil {
		glog.Errorf("Failed to delete notification. %s", err)
		return false, err
	}

	if res.DeletedCount == 0 {
		return false, nil
	}

	return true, nil
}

// DeleteAll deletes notifications with IDs and returns true if any record was removed
func (r *Repository) DeleteAll(uuids []apimongo.UUID) (bool, error) {
	glog.Infof("Deleting notification with UUID: %v", uuids)

	collection := r.Mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": uuids}}
	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		glog.Error("Failed to delete notifications. ", err)
		return false, err
	}

	if res.DeletedCount == 0 {
		return false, nil
	}

	return true, nil
}
