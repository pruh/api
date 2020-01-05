package dao

import (
	"context"
	"time"

	"github.com/golang/glog"

	"github.com/pruh/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository accesses notifications store
type Repository struct {
	mongo *mongo.Client
}

const (
	dbName              = "ApiDB"
	notifCollectionName = "notifications"
)

// NewRepository creates new repository and sets up connection to DB
func NewRepository() *Repository {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		glog.Fatalf("Cannot connect to databse. %s", err)
	}

	return &Repository{
		mongo: client,
	}
}

// GetNofitications returns all notifications
func (r *Repository) GetNofitications() ([]models.Notification, error) {
	glog.Info("Querying for all notifications")
	collection := r.mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		glog.Errorf("Error while performing DB find. %s", err)
		return nil, err
	}
	defer cur.Close(ctx)

	var notifs []models.Notification
	for cur.Next(ctx) {
		var notif models.Notification
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

// GetNofitication returns notifications by ID or nil
func (r *Repository) GetNofitication(ID string) (models.Notification, error) {
	glog.Infof("Querying for notification with UUID: %s\n", ID)

	collection := r.mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var notif models.Notification
	uuid, err := models.ParseMongoUUID(ID)
	if err != nil {
		glog.Errorf("Cannot parse ID. %s", err)
		return notif, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": models.MongoUUID(uuid)}).Decode(&notif)
	if err != nil {
		glog.Errorf("Cannot decode notification. %s", err)
		return notif, err
	}

	return notif, nil
}

// CreateNofitication creates new notification for specified params
func (r *Repository) CreateNofitication(notification models.Notification) bool {
	glog.Infof("Creating new notification: %+v\n", notification)

	collection := r.mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, notification)
	if err != nil {
		glog.Errorf("Failed to insert notification. %s", err)
		return false
	}

	return true
}

// DeleteNofitication deletes notifications with ID
func (r *Repository) DeleteNofitication(ID string) bool {
	glog.Infof("Deleting notification with UUID: %s\n", ID)

	collection := r.mongo.Database(dbName).Collection(notifCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uuid, err := models.ParseMongoUUID(ID)
	if err != nil {
		glog.Errorf("Cannot parse ID. %s", err)
		return false
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": models.MongoUUID(uuid)})
	if err != nil {
		glog.Errorf("Failed to delete notification. %s", err)
		return false
	}

	return true
}
