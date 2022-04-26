package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Delete struct {
	Logger      *log.Logger
	MongoClient *mongo.Client
}

// type Delete struct {
// 	Logger       *log.Logger
// 	MongoClient  *mongo.Client
// 	RabbitClient *amqp.Connection
// }

type DeleteUser struct {
	ID int
}

func NewDelete(logger *log.Logger, mongoClient *mongo.Client) *Delete {
	return &Delete{logger, mongoClient}
}

// func NewDelete(logger *log.Logger, mongoClient *mongo.Client, rabbitClient *amqp.Connection) *Delete {
// 	return &Delete{logger, mongoClient, rabbitClient}
// }

// delete using ID
// Send a JSON file containing the ID of the document to be deleted
func (delete *Delete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	delete.Logger.Println("Received request for 'DELETE'")
	w.WriteHeader(http.StatusOK) // 200 OK

	// read the json file sent by users
	var received Student
	json.NewDecoder(r.Body).Decode(&received)

	// reject the request if the received data is empty
	if received == (Student{}) {
		delete.Logger.Println("Received empty data")
		return
	}

	// delete the instance from MongoDB
	collection := delete.MongoClient.Database("vodascheduler").Collection("model")

	if collection == nil {
		delete.Logger.Println("Collection not found")
		return
	}

	filter := bson.M{"id": received.ID}
	var found Student

	// Find instance and decode the found result into the "found" variable.
	collection.FindOne(context.TODO(), filter).Decode(&found)

	if received.ID == found.ID { // Student IDs should be unique
		_, err := collection.DeleteOne(context.TODO(), bson.M{"id": received.ID})
		if err != nil {
			delete.Logger.Println("Error deleting from MongoDB: ", err)
			return
		}

		delete.Logger.Println("Deleted data: ", found)
	} else {
		// else reject the request
		delete.Logger.Println("Student ID doesn't exists")
		delete.Logger.Println("The end of 'POST' request")
		return
	}

	var result Student
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		delete.Logger.Println("Failed to query MongoDB:", err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&result)
		if err != nil {
			delete.Logger.Println("Failed to decode cursor:", err)
			return
		}
		delete.Logger.Println("Student:", result)
	}

	delete.Logger.Println("The end of 'DELETE' request")
	w.Write([]byte("The end of 'DELETE' request"))
}
