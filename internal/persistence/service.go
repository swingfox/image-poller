package persistence

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/swingfox/image-poller/cmd/webapp/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
}

func GetCollection() *mongo.Collection {
	dbUsername := config.Registry.GetString("DB.USERNAME")
	dbPassword := config.Registry.GetString("DB.PASSWORD")
	name := config.Registry.GetString("DB.NAME")
	host := config.Registry.GetString("DB.HOST")
	port := config.Registry.GetString("DB.PORT")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://"+dbUsername+":"+dbPassword+"@"+host+":"+port))
	if err != nil {
		log.Fatal("error connecting mongodb", err)
	}
	return client.Database(name).Collection("products")
}
