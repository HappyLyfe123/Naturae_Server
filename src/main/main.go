package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	//CreateAccount("", "g", "g")
	client, err := mongo.Connect(context.TODO(), "mongodb+srv://HappyLyfe:kePmTHH8wyrSEIxL@naturae-server-hxywc.mongodb.net/test?retryWrites=true")
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	fmt.Println("Connected to MongoDB!")
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("Naturae-Server").Collection("Users")
	ash := Trainer{"Ash", 10, "Pallet Town"}
	insertResult, err := collection.InsertOne(context.TODO(), ash)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}
