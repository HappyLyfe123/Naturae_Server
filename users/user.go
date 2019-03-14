package users

import (
	"Naturae_Server/helpers"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
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
func saveAccessToken(database *mongo.Database, token *helpers.AccessToken) {
	for saveSuccessful := false; saveSuccessful == false; {
		connectedCollection := helpers.ConnectToCollection(database, helpers.GetAccessTokenCollection())
		_, err := connectedCollection.InsertOne(nil, token)
		//If there an duplicate ID generate a new a new ID and try to save again
		if err != nil {
			//Generate a new token ID
			token.ID = helpers.GenerateTokenID()
		} else {
			log.Println("Save", token, "access token to access token database")
			saveSuccessful = true
		}
	}

}

func saveRefreshToken(database *mongo.Database, token *helpers.RefreshToken) {

	for saveSuccessful := false; saveSuccessful == false; {
		connectedCollection := helpers.ConnectToCollection(database, helpers.GetRefreshTokenCollection())
		_, err := connectedCollection.InsertOne(nil, token)
		if err != nil {
			//Generate a new token ID
			token.ID = helpers.GenerateTokenID()
		} else {
			log.Println("Save", token.Email, "refresh token to refresh token database")
			saveSuccessful = true
		}

	}

}

//Update the user first name in the database
func upDateUserFirstName(name string) {

}

//Update the user last name in the database
func upDateUserLastName(name string) {

}
