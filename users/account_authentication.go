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
	userInfoDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	var statusCode int32
	var authenMessage string
	var accessToken *helpers.AccessToken
	var refreshToken *helpers.RefreshToken

	//Get the authentication result back from the database
	authenResult, err := getAuthenCode(userInfoDB, request.GetEmail())
	//Check if there an error that had occurred when retrieving information for the database server
	if err != nil {
		return &pb.AccountAuthenReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "",
			Email: "", Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error."}}
	}

	//Get the user information from the database
	userInfoResult, _ := getUserInfo(request.GetEmail())

	//Check if the authentication code is still valid. If it's valid then the code will be check
	if helpers.IsTimeValid(authenResult.ExpiredTime) {
		//Check if the user provided authentication code match. If it match then update the user's profile in the database
		//that the user's had authenticated their account
		if strings.Compare(authenResult.Code, request.GetAuthenCode()) == 0 {
			//Set authenticated in user database from false to true
			updateUserAuthenStatus(userInfoDB, request.GetEmail())
			//Remove the user authentication code from the database
			go removeAuthenCode(userInfoDB, request.GetEmail())
			//Create access token
			accessToken = helpers.GenerateAccessToken(request.GetEmail(), userInfoResult.FirstName, userInfoResult.LastName)
			//Save access token to database
			saveAccessToken(userInfoDB, accessToken)
			//Create refresh token
			refreshToken = helpers.GenerateRefreshToken(request.GetEmail())
			//Save refresh token to database
			saveRefreshToken(userInfoDB, refreshToken)

			statusCode = helpers.GetOkStatusCode()
			authenMessage = "Account has been successfully authenticated"

		} else {
			return &pb.AccountAuthenReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "",
				Email: "", Status: &pb.Status{Code: helpers.GetInvalidCode(), Message: "Invalid authen code"}}
		}
	} else {
		//Create a new authen code for the user because the old one expired
		newAuthenCode(userInfoDB, userInfoResult.Email, userInfoResult.FirstName)
		return &pb.AccountAuthenReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "",
			Email: "", Status: &pb.Status{Code: helpers.GetExpiredAuthenCode(), Message: "Code expired"}}
	}
	return &pb.AccountAuthenReply{AccessToken: accessToken.ID, RefreshToken: refreshToken.ID, FirstName: userInfoResult.FirstName,
		LastName: userInfoResult.LastName, Email: userInfoResult.Email, Status: &pb.Status{Code: statusCode, Message: authenMessage}}
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

func newAuthenCode(userDB *mongo.Database, email, firstName string) {
	//Generate authentication code and expired time
	authenCode, expiredTime := helpers.GenerateAuthenCode()
	//Update the authen code database with the new authen code and expiration date
	//Set the filter to find the user
	filter := bson.D{{"email", email}}
	//Update the IsAuthenticated field from false to true
	update := bson.D{{"$set", bson.D{{"code", authenCode}, {"expiredtime", expiredTime}}}}
	//Connect to the database and update the information
	_, err := helpers.ConnectToCollection(userDB, helpers.GetAccountAuthenticationCollection()).UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error occurred when try to create new authen code: %v", err)
	}
	//Send the user authentication code to the user's email
	sendAuthenticationCode(email, firstName, authenCode)
}
