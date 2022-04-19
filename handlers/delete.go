package handlers

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Delete struct {
	logger      *log.Logger
	MongoClient *mongo.Client
}

// http://localhost:10000/post/ udah bisa diakses

func NewDelete(logger *log.Logger, mongoClient *mongo.Client) *Delete {
	return &Delete{logger, mongoClient}
}

func (delete *Delete) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
