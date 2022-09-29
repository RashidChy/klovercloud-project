package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"project1/config"
	"time"
)

func createConnectionStringDB() string {
	var dbUser = config.DbUser
	var dbPass = config.DbPass
	var clusterEndpoint = config.ClusterEndpoint
	connectionString := "mongodb+srv://" + dbUser + ":" + dbPass + "@" + clusterEndpoint
	return connectionString
}

func InitDBConnection() (*mongo.Client, error) {
	uri := createConnectionStringDB()
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("[ERROR] Unable to create MongoClient: ", err.Error())
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println("[ERROR] Cannot connect MongoClient: ", err.Error())
		return nil, err
	}

	return client, nil
}
