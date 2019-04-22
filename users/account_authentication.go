package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"strings"
)

//Check if the authentication code is valid
func AuthenticateAccount(request *pb.AccountAuthenRequest) *pb.AccountAuthenReply {
	connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	var statusCode int32
	var authenMessage string
	var accessToken *helpers.AccessToken
	var refreshToken *helpers.RefreshToken

	//Get the authentication result back from the database
	databaseResult, err := getAuthenCode(connectedDB, request.GetEmail())
	//Check if there an error that had occurred when retrieving information for the database server
	if err != nil {
		return &pb.AccountAuthenReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "",
			Email: "", Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error."}}
	}

	//Check if the authentication code is still valid. If it's valid then the code will be check
	if helpers.IsTimeValid(databaseResult.ExpiredTime) {
		//Check if the user provided authentication code match. If it match then update the user's profile in the database
		//that the user's had authenticated their account
		if strings.Compare(databaseResult.Code, request.GetAuthenCode()) == 0 {

			//Set authenticated in user database from false to true
			updateUserAuthenStatus(connectedDB, request.GetEmail())
			//Remove the user authentication code from the database
			go removeAuthenCode(connectedDB, request.GetEmail())
			//Create access token
			accessToken = helpers.GenerateAccessToken(request.GetEmail(), "", "")
			//Save access token to database
			saveAccessToken(connectedDB, accessToken)
			//Create refresh token
			refreshToken = helpers.GenerateRefreshToken(request.GetEmail())
			//Save refresh token to database
			saveRefreshToken(connectedDB, refreshToken)

			statusCode = helpers.GetOkStatusCode()
			authenMessage = "Account has been successfully authenticated"

		} else {
			return &pb.AccountAuthenReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "",
				Email: "", Status: &pb.Status{Code: helpers.GetInvalidAuthenCode(), Message: "Invalid authen code"}}
		}
	}
	return &pb.AccountAuthenReply{AccessToken: accessToken.ID, RefreshToken: refreshToken.ID, FirstName: accessToken.FirstName,
		LastName: accessToken.LastName, Email: accessToken.Email, Status: &pb.Status{Code: statusCode, Message: authenMessage}}
}

//Change the user authentication status
func updateUserAuthenStatus(connectedDB *mongo.Database, email string) {
	//Set the filter to find the user
	filter := bson.D{{"email", email}}
	//Update the IsAuthenticated field from false to true
	update := bson.D{{"$set", bson.D{{"isauthenticated", true}}}}
	//Connect to the database and update the information
	_, err := helpers.ConnectToCollection(connectedDB, helpers.GetAccountInfoCollection()).UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
	}
}

//Remove the authentication code from the server
func removeAuthenCode(connectedDB *mongo.Database, email string) {
	//Find the
	filter := bson.D{{"email", email}}
	_, err := helpers.ConnectToCollection(connectedDB, helpers.GetAccountAuthenticationCollection()).DeleteOne(context.Background(), filter)

	if err != nil {
		log.Println(err)
	}

}
