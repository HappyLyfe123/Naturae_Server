package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
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
	err := userCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return &loginInfo{}, err
	}
	return &result, nil

}

func getAuthenCode(database *mongo.Database, email string) (*userAuthentication, error) {
	var result userAuthentication
	filter := bson.D{{Key: "email", Value: email}}
	//Connect to the collection database
	userCollection := helpers.ConnectToCollection(database, helpers.GetAccountAuthenticationCollection())
	//Make a request to the database
	err := userCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return &userAuthentication{}, err
	}
	return &result, nil

}

//Save access token to database
func saveAccessToken(database *mongo.Database, token *helpers.AccessToken) {
	for saveSuccessful := false; saveSuccessful == false; {
		connectedCollection := helpers.ConnectToCollection(database, helpers.GetAccessTokenCollection())
		_, err := connectedCollection.InsertOne(context.Background(), token)
		//If there an duplicate ID generate a new a new ID and try to save again
		if err != nil {
			//Generate a new token ID
			token.ID = helpers.GenerateTokenID()
		} else {
			log.Println("Save", token.Email, "access token to access token database")
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

//Generate a new access token for the user when their access token is expired
func RefreshAccessToken(request pb.GetAccessTokenRequest){
	var refreshTokenResult helpers.RefreshToken
	currConnectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	connectedCollection := currConnectedDB.Collection(helpers.GetRefreshTokenCollection())

	filter := bson.D{{Key: "id", Value: request.GetRefreshToken()}}
	err := connectedCollection.FindOne(context.Background(), filter).Decode(&refreshTokenResult)

	if err != nil{

	}else{
		//Compare the refresh token id in the database to the one that the request provided. If the two string match
		//then the server will generate a new access token for the user's. If not then the user will return an error
		if strings.Compare(refreshTokenResult.ID, request.GetRefreshToken()) == 0{

		}
	}


}
