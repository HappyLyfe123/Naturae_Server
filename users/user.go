package users

import (
	"Naturae_Server/helpers"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//create a struct for storing user info in database
type userAccount struct {
	Email           string
	FirstName       string
	LastName        string
	Salt            string
	Password        string
	IsAuthenticated bool
}

type accessToken struct {
	Email       string
	ID          string
	admin       bool
	ExpiredTime time.Time
}

//Create a struct for storing token
type refreshToken struct {
	Email       string
	ID          string
	ExpiredTime time.Time
}

//The the user account info from the database
func getUserAccountInfo(database *mongo.Database, email string) (*userAccount, error) {

	var result userAccount
	//Set a filter for the database to search through
	filter := bson.D{{Key: "email", Value: email}}
	//Connect to the collection database
	userCollection := helpers.ConnectToCollection(database, helpers.GetAccountInfoCollection())
	//Make a request to the database
	err := userCollection.FindOne(nil, filter).Decode(&result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func getLoginInfo(database *mongo.Database, email string) (*loginInfo, error) {
	var result loginInfo

	filter := bson.D{{Key: "email", Value: email}}
	//Connect to the collection database
	userCollection := helpers.ConnectToCollection(database, helpers.GetAccountInfoCollection())
	//Make a request to the database
	err := userCollection.FindOne(nil, filter).Decode(&result)
	if err != nil {
		return &loginInfo{}, err
	}
	return &result, nil

}

func getAuthenCode(database *mongo.Database, email string) (*userAuthentication, error) {
	var result userAuthentication
	filter := bson.D{{Key: "email", Value: email}}
	//Connect to the collection database
	userCollection := helpers.ConnectToCollection(database, helpers.GetAccountAuthentication())
	//Make a request to the database
	err := userCollection.FindOne(nil, filter).Decode(&result)
	if err != nil {
		return &userAuthentication{}, err
	}
	return &result, nil

}

//Save access token to database
func (token *accessToken) saveToken(wg *sync.WaitGroup, database *mongo.Database) helpers.AppError {
	defer wg.Done()
	connectedCollection := helpers.ConnectToCollection(database, helpers.GetAccessTokenCollection())
	_, err := connectedCollection.InsertOne(nil, token)
	if err != nil {
		return helpers.AppError{Code: helpers.GetDuplicateInfoCode(), Type: "duplicate token id", Description: err.Error()}
	}
	log.Println("Save", token.Email, "access token to access token database")
	return helpers.AppError{Code: helpers.GetNoErrorCode(), Type: "None", Description: "None"}
}

func (token *refreshToken) saveToken(wg *sync.WaitGroup, database *mongo.Database) helpers.AppError {
	defer wg.Done()
	connectedCollection := helpers.ConnectToCollection(database, helpers.GetRefreshTokenCollection())
	_, err := connectedCollection.InsertOne(nil, token)
	if err != nil {
		return helpers.AppError{Code: helpers.GetDuplicateInfoCode(), Type: "duplicate token id", Description: err.Error()}
	}
	log.Println("Save", token.Email, "refresh token to refresh token database")
	return helpers.AppError{Code: helpers.GetNoErrorCode(), Type: "None", Description: "None"}
}

//Update the user first name in the database
func upDateUserFirstName(name string) {

}

//Update the user last name in the database
func upDateUserLastName(name string) {

}
