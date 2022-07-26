package persistence

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/swingfox/image-poller/cmd/webapp/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GetCollection Get MongoDB collection.
func GetCollection(collectionName string) *mongo.Collection {
	dbUsername := config.Registry.GetString("DB.USERNAME")
	dbPassword := config.Registry.GetString("DB.PASSWORD")
	name := config.Registry.GetString("DB.NAME")
	host := config.Registry.GetString("DB.HOST")
	port := config.Registry.GetString("DB.PORT")

	srv := "mongodb://" + dbUsername + ":" + dbPassword + "@" + host + ":" + port + "/" + name + "?authSource=admin"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(srv))
	if err != nil {
		log.Fatal("error connecting mongodb", err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal("error pinging mongodb", err)
	}
	return client.Database(name).Collection(collectionName)
}
