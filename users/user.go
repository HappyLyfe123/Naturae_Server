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
	connectedCollection := helpers.ConnectToCollection(database, helpers.GetAccessTokenCollection())
	_, err := connectedCollection.InsertOne(context.Background(), token)
	//If there an duplicate ID generate a new a new ID and try to save again
	if err != nil {
		//Generate a new token ID
		token.ID = helpers.GenerateTokenID()
	} else {
		log.Println("Save", token.Email, "access token to access token database")
	}

}

//Save the refresh token to the database
func saveRefreshToken(database *mongo.Database, token *helpers.RefreshToken) {
	connectedCollection := helpers.ConnectToCollection(database, helpers.GetRefreshTokenCollection())
	_, err := connectedCollection.InsertOne(nil, token)
	if err != nil {
		//Generate a new token ID
		token.ID = helpers.GenerateTokenID()
	} else {
		log.Println("Save", token.Email, "refresh token to refresh token database")
	}

}

//Generate a new access token for the user when their access token is expired
func RefreshAccessToken(request *pb.GetAccessTokenRequest) *pb.GetAccessTokenReply {
	var refreshTokenResult helpers.RefreshToken
	var status *pb.Status
	newAccessToken := ""
	currConnectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	connectedCollection := currConnectedDB.Collection(helpers.GetRefreshTokenCollection())

	//Set a filter to find the user refresh token in the database
	filter := bson.D{{Key: "id", Value: request.GetRefreshToken()}}
	//Connect to the database to retrieve the user's refresh token
	err := connectedCollection.FindOne(context.Background(), filter).Decode(&refreshTokenResult)
	if err != nil {
		status = &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error"}
	} else {
		//Compare the refresh token id in the database to the one that the request provided. If the two string match
		//then the server will generate a new access token for the user's. If not then the user will return an error
		if strings.Compare(refreshTokenResult.ID, request.GetRefreshToken()) == 0 {
			accessToken := helpers.GenerateAccessToken(refreshTokenResult.Email)
			//Set the filter to find the user
			filter := bson.D{{"email", refreshTokenResult.Email}}
			//Update the access token id and expired time in the database with the newly generated one
			update := bson.D{{"$set", bson.D{{"id", accessToken.ID}, {"expiredtime", accessToken.ExpiredTime}}}}
			_, err := helpers.ConnectToCollection(currConnectedDB, helpers.GetAccessTokenCollection()).UpdateOne(context.Background(), filter, update)
			if err != nil {
				status = &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error"}
			} else {
				newAccessToken = accessToken.ID
				status = &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Code has been generated"}
			}
		} else {
			status = &pb.Status{Code: helpers.GetInvalidTokenCode(), Message: "Token is invalid"}
		}
	}

	return &pb.GetAccessTokenReply{AccessToken: newAccessToken, Status: status}
}
