package users

import (
	"Naturae_Server/helpers"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"sync"
	"time"
)

//create a struct for storing user info in database
type userAccount struct {
	Email            string
	First_Name       string
	Last_Name        string
	Salt             string
	Password         string
	Is_Authenticated bool
}

//create authentication structure for user
type userAuthentication struct {
	Email      string
	Code       string
	Start_Time time.Time
}

//create a struct for storing token
type token struct {
	Email      string
	Token_ID   string
	Start_Time time.Time
}

//SaveToken : save token to the database
func saveToken(wg *sync.WaitGroup, database *mongo.Database, collectionName string, token *token) {
	defer wg.Done()
	//If there an error, it will attempt to save the info to the database until the limit is reach
	for numOfAttempt := 0; numOfAttempt < saveAttemptLimit; numOfAttempt++ {
		//Connect to the database collection
		currCollection := helpers.ConnectToCollection(database, collectionName)
		_, err := currCollection.InsertOne(context.TODO(), token)
		if err != nil {
			log.Println("Save token error: ", err)
		} else {
			break
		}
	}
}
