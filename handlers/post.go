package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Post struct {
	logger      *log.Logger
	MongoClient *mongo.Client
}

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Students []Student

func NewPost(logger *log.Logger, mongoClient *mongo.Client) *Post {
	return &Post{logger, mongoClient}
}

func (post *Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	post.logger.Println("Received request for 'POST'")
	w.WriteHeader(http.StatusOK) // 200 OK
	w.Write([]byte("Received request for 'POST'"))

	// read the json file sent by users
	var received Student
	json.NewDecoder(r.Body).Decode(&received)
	if received == (Student{}) {
		post.logger.Println("Received empty data")
		return
	} else {
		post.logger.Println("Received data: ", received)
	}

	// store students to Mongodb
	collection := post.MongoClient.Database("vodascheduler").Collection("model")
	filter := bson.M{"id": received.ID}
	var result Student

	// Find instance and decode the found result in result variable.
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		// If a student does not exist, insert his/her data into the database
		collection.InsertOne(context.TODO(), received)
		post.logger.Println("Successfully inserted student:", received)
	} else {
		post.logger.Println("Instance already exist")
		return
	}

	post.logger.Println("The end of 'POST' request")
	w.Write([]byte("The end of 'POST' request"))
}
