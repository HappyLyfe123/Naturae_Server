package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//create a struct for storing user info in database
type userAccount struct {
	Email           string
	FirstName       string
	LastName        string
	Salt            string
	Password        string
	Friends         []string
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

/*
*	Searches the database for matching users and returns the list
* 	@UserSearchRequest - string user, string query
* 	@UserListReply - repeated string users, Status status
 */
func SearchUsers(request *pb.UserSearchRequest) *pb.UserListReply {
	var result userAccount
	var searchResult []string

	dbConnection := helpers.ConnectToDB(helpers.GetUserDatabase())
	userEmail := request.GetUser()

	//Set a filter for the database to search through()
	if len(userEmail) > 0 {
		//Get the user account to retrieve the friend's list from
		userAccount, err := getUserAccountInfo(dbConnection, userEmail)
		searchResult = userAccount.Friends
	} else {
		//Search via user query input instead, get all users with the email
		queryEmail := request.GetQuery()
		queryString := strings.Split(queryEmail, " ")
		querySplice := make([]bson.M, len(queryString))
		for i, queryEmail := range queryString {
			querySplice[i] = bson.M{"email": bson.M{
				"$regex": bson.RegEx{Pattern: ".*" + queryEmail + ".*", Options: "i"},
			}}
		}

		//Connect to Account_Information collections
		userCollection := helpers.ConnectToCollection(dbConnection, helpers.GetAccountInfoCollection())
		//Make the search request
		filter := bson.M{"$querySplice": querySplice}
		//&filter?
		err := userCollection.Find(nil, &filter).Limit(10).All(&searchResult)
		if err != nil {

		}

		return &pb.UserListReply{User: searchResult,
			Status: &pb.Status{Code: 200, Message: "Created"}}
	}
}

//Adds a friend to a user's list of contacts
func AddFriend(request *pb.FriendRequest) *pb.FriendReply {

}

//Removes a friend from a user's list of contacts
func RemoveFriend(request *pb.FriendRequest) *pb.FriendReply {

}
