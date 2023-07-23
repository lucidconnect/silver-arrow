package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/helicarrierstudio/silver-arrow/repository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	client *mongo.Client
}

func NewMongoDb(mongo *mongo.Client) *MongoDb {
	return &MongoDb{mongo}
}

func SetupMongoDatabase() (*mongo.Client, error) {
	uri := os.Getenv("DATABASE_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, err
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}

func (m *MongoDb) SetAddress(wallet models.Wallet) error {
	coll := m.client.Database("db").Collection("wallets")
	_, err := coll.InsertOne(context.Background(), wallet)
	return err
}

func (m *MongoDb) AddSubscription(subscriptionData models.Subscription) error {
	coll := m.client.Database("db").Collection("subscriptions")
	_, err := coll.InsertOne(context.Background(), subscriptionData)
	return err
}

func (m *MongoDb) ListSubscriptions(address string) ([]models.Subscription, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (m *MongoDb) RemoveSubscription(id int64) error {
	return errors.New("NOT IMPLEMENTED")
}
