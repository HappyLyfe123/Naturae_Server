package main

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"os"
)

var dbClient *mongo.Client = nil;

func main() {

	//Connect to the database sever

	fmt.Println(os.Getenv("DATABASE_USERNAME"))
	fmt.Println(os.Getenv("DATABASE_PASSWORD"))
	fmt.Println("Work")

}
