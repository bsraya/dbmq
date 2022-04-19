package handlers

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Get struct {
	logger  *log.Logger
	mongodb *mongo.Client
}

func NewGet(logger *log.Logger, mongoClient *mongo.Client) *Get {
	return &Get{logger, mongoClient}
}

func (get *Get) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
