package main

import (
	"context"
	"fmt"
	"log"
	"mqdb/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURL string = "mongodb://localhost:27017"

// const rabbitURL string = "amqp://guest:guest@localhost:5672"

func connectMongo(url string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return client
}

// func connectRabbitMQ(url string) *amqp.Connection {
// 	conn, err := amqp.Dial(url)

// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}

// 	return conn
// }

// curl http://localhost:10000/post/ simpan di database
// curl http://localhost:10000/get/ dapatkan dari database
// curl http://localhost:10000/delete/ hapus dari database

func main() {
	// connect to mongodb
	mongoClient := connectMongo(mongoURL)
	defer mongoClient.Disconnect(context.TODO())
	fmt.Println("Connected to MongoDB")

	logger := log.New(os.Stdout, "vodascheduler", log.LstdFlags)

	postHandler := handlers.NewPost(logger, mongoClient)
	getHandler := handlers.NewGet(logger, mongoClient)
	deleteHandler := handlers.NewDelete(logger, mongoClient)

	// /post -> http://localhost:10000/post/
	// /get -> http://localhost:10000/get/
	// /delete -> http://localhost:10000/delete/

	serveMux := http.NewServeMux()
	serveMux.Handle("/post/", postHandler)
	serveMux.Handle("/get/", getHandler)
	serveMux.Handle("/delete/", deleteHandler)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// curl -d "id=1&name=Tukul" -H "Content-Type: application/x-www-form-urlencoded" -X POST http://localhost:9090/post/

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}() // run in background

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	signal := <-signalChannel
	logger.Println("Received termination, gracefully shutdown", signal)

	tc, _ := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	server.Shutdown(tc)

	// // store students to Mongodb
	// collection := mongoClient.Database("vodascheduler").Collection("model")
	// for _, student := range students {
	// 	// check if student already exists
	// 	filter := bson.M{"id": student.ID}
	// 	var result Student
	// 	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	// 	if err == mongo.ErrNoDocuments {
	// 		// student does not exist, insert
	// 		collection.InsertOne(context.TODO(), student)
	// 		fmt.Println("Successfully inserted student:", student)
	// 	} else {
	// 		fmt.Println("Instance already exist")
	// 		continue
	// 	}
	// }

	// // connect to RabbitMQ
	// rmqConnection := connectRabbitMQ(rabbitURL)
	// defer rmqConnection.Close()
	// fmt.Println("Connected to RabbitMQ")

	// // create channel
	// channel, err := rmqConnection.Channel()
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// defer channel.Close()

	// // create queue declare for students
	// queue, err := channel.QueueDeclare(
	// 	"students", // name
	// 	false,      // durable
	// 	false,      // delete when unused
	// 	false,      // exclusive
	// 	false,      // no-wait
	// 	nil,        // arguments
	// )

	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

	// fmt.Println(queue)

	// // push multiple messages into the queue
	// for _, student := range students {
	// 	err = channel.Publish(
	// 		"",
	// 		"students",
	// 		false,
	// 		false,
	// 		amqp.Publishing{
	// 			ContentType: "text/plain",
	// 			// body should include ID and Name
	// 			Body: []byte(fmt.Sprintf("%d %s", student.ID, student.Name)),
	// 		},
	// 	)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		panic(err)
	// 	}
	// 	fmt.Println("Successfully published messages to queue!")
	// }
	// fmt.Println(queue)
}
