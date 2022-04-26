package handlers

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Get struct {
	Logger      *log.Logger
	MongoClient *mongo.Client
}

// type Get struct {
// 	Logger       *log.Logger
// 	MongoClient  *mongo.Client
// 	RabbitClient *amqp.Connection
// }

func NewGet(logger *log.Logger, mongoClient *mongo.Client) *Get {
	return &Get{logger, mongoClient}
}

// func NewGet(logger *log.Logger, mongoClient *mongo.Client, rabbitClient *amqp.Connection) *Get {
// 	return &Get{logger, mongoClient, rabbitClient}
// }

func (get *Get) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	get.Logger.Println("Received request for 'GET'")
	w.WriteHeader(http.StatusOK) // 200 OK
	w.Write([]byte("Received request for 'GET'"))

	// print out all the users in the MongoDB
	collection := get.MongoClient.Database("vodascheduler").Collection("model")
	if collection == nil {
		get.Logger.Println("Collection not found")
		return
	}

	var result Student
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		get.Logger.Println("Failed to query MongoDB:", err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&result)
		if err != nil {
			get.Logger.Println("Failed to decode cursor:", err)
			return
		}
		get.Logger.Println("Student:", result)
	}

	get.Logger.Println("The end of 'GET' request")
	w.Write([]byte("The end of 'GET' request"))
}
