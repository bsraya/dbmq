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

func main() {
	// connect to mongodb
	mongoClient := connectMongo(mongoURL)
	defer mongoClient.Disconnect(context.TODO())
	fmt.Println("Connected to MongoDB")

	logger := log.New(os.Stdout, "vodascheduler ", log.LstdFlags)

	serveMux := http.NewServeMux()
	serveMux.Handle("/post/", handlers.NewPost(logger, mongoClient))
	serveMux.Handle("/get/", handlers.NewGet(logger, mongoClient))
	serveMux.Handle("/delete/", handlers.NewDelete(logger, mongoClient))

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

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
