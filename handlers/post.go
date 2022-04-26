package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// const rabbitURL string = "amqp://guest:guest@localhost:5672"

type Post struct {
	Logger      *log.Logger
	MongoClient *mongo.Client
}

// type Post struct {
// 	Logger       *log.Logger
// 	MongoClient  *mongo.Client
// 	RabbitCLient *amqp.Connection
// }

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Students []Student

func NewPost(logger *log.Logger, mongoClient *mongo.Client) *Post {
	return &Post{logger, mongoClient}
}

// func NewPost(logger *log.Logger, mongoClient *mongo.Client, rabbitClient *amqp.Connection) *Post {
// 	return &Post{logger, mongoClient, rabbitClient}
// }

func (post *Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	post.Logger.Println("Received request for 'POST'")
	w.WriteHeader(http.StatusOK) // 200 OK
	w.Write([]byte("Received request for 'POST'"))

	// read the json file sent by users
	var received Student
	json.NewDecoder(r.Body).Decode(&received)

	// reject the request if the received data is empty
	if received == (Student{}) {
		post.Logger.Println("Received empty data")
		return
	}

	post.Logger.Println("Received data: ", received)

	// access "model" collection in "vodascheduler" database
	collection := post.MongoClient.Database("vodascheduler").Collection("model")
	if collection == nil {
		post.Logger.Println("Collection not found")
		return
	}

	filter := bson.M{"id": received.ID}
	var found Student

	// Find instance and decode the found result into the "found" variable.
	collection.FindOne(context.TODO(), filter).Decode(&found)

	if received.ID != found.ID { // Student IDs should be unique
		// If a student ID doesn't exist, insert his/her data into the database
		collection.InsertOne(context.TODO(), received)
		post.Logger.Println("Successfully inserted student:", received)
	} else {
		// else reject the request
		post.Logger.Println("Student ID already exists")
		post.Logger.Println("The end of 'POST' request")
		return
	}

	post.Logger.Println("The end of 'POST' request")
	w.Write([]byte("The end of 'POST' request"))
}
