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
	var result bool
	databaseResult, err := getAuthenCode(connectedDB, request.GetEmail())
	if err != nil {
		return &pb.AccountAuthenReply{Result: false, Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error."}}
	}
	//Check if the authentication code is expired
	if helpers.IsTimeValid(databaseResult.ExpiredTime) {
		//Check if the user provided authentication code match
		if strings.Compare(databaseResult.Code, request.GetAuthenCode()) == 1 {
			updateUserAuthenStatus(connectedDB, request.GetEmail())
			//Remove the user authentication code from the database
			go removeAuthenCode(connectedDB, request.GetEmail())
			statusCode = helpers.GetOkStatusCode()
			authenMessage = "Account has been successfully authenticated"
			result = true

		} else {
			statusCode = helpers.GetInvalidAuthenCode()
			authenMessage = "Invalid authen code"
		}
	} else {
		//Generate authentication code and expired time
		authenCode, expiredTime := helpers.GenerateAuthenCode()
		//Create a struct for user's authentication
		newAuthenCode := userAuthentication{Email: request.GetEmail(), Code: authenCode, ExpiredTime: expiredTime}
		//Save the user authentication code to the database
		saveAuthenticationCode(connectedDB, helpers.GetAccountAuthenticationCollection(), &newAuthenCode)
		//Send the user authentication code to the user's email
		sendAuthenticationCode(request.GetEmail(), request.GetFirstName(), authenCode)
	}
	return &pb.AccountAuthenReply{Result: result, Status: &pb.Status{Code: statusCode, Message: authenMessage}}
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
