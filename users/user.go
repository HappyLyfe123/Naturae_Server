package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"context"
	"log"

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
	var searchResult []string

	dbConnection := helpers.ConnectToDB(helpers.GetUserDatabase())
	userEmail := request.GetUser()

	//if userEmail is not nil, it means to retrieve the friendslist of the logged in user
	if len(userEmail) > 0 {
		userAccount, err := getUserAccountInfo(dbConnection, userEmail)
		searchResult = userAccount.Friends
		if err != nil {
			return &pb.UserListReply{User: nil,
				Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Failed UserEmail Get"}}
		}
	} else {
		//Search for exactly one user with an inputted email
		queryEmail := request.GetQuery()
		userAccount, err := getUserAccountInfo(dbConnection, queryEmail)
		searchResult = userAccount.Friends
		if err != nil {
			return &pb.UserListReply{User: nil,
				Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Failed QueryEmail Get"}}
		}
		/*Search via user query input instead, get all users with similar email(not working)
		queryString := strings.Split(queryEmail, " ")
		querySplice := make([]bson.M, len(queryString))
		for i, queryEmail := range queryString {
			querySplice[i] = bson.M{"email": bson.M{
				"$regex": bson.Regex{Pattern: ".*" + queryEmail + ".*", Options: "i"},
			}}
		}

		}
			//Connect to Account_Information collections
			userCollection := helpers.ConnectToCollection(dbConnection, helpers.GetAccountInfoCollection())
			//Make the search request
			filter := bson.M{"$querySplice": querySplice}
			//&filter?
			err := userCollection.Find(nil, &filter).Limit(10).All(&searchResult)
		*/
	}

	return &pb.UserListReply{User: searchResult,
		Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Okay"}}
}

//Adds a friend to a user's list of contacts
func AddFriend(request *pb.FriendRequest) *pb.FriendReply {
	dbConnection := helpers.ConnectToDB(helpers.GetUserDatabase())
	userCollection := helpers.ConnectToCollection(dbConnection, helpers.GetAccountInfoCollection())
	//Set a filter for the database to search through
	senderFilter := bson.D{{Key: "email", Value: request.GetSender()}}
	receiverFilter := bson.D{{Key: "email", Value: request.GetReceiver()}}

	cherr := make(chan error, 2)
	//Update Interface, Push Friend Usernames to Friendslist of both users
	go func() {
		updateSender := bson.D{
			{"$push", bson.D{
				{"Friends", request.GetReceiver()},
			}},
		}
		_, err := userCollection.UpdateOne(
			context.Background(),
			senderFilter,
			updateSender,
		)
		if err != nil {
			cherr <- err
		}
	}()
	go func() {
		updateReceiver := bson.D{
			{"$push", bson.D{
				{"Friends", request.GetSender()},
			}},
		}
		_, err := userCollection.UpdateOne(
			context.Background(),
			receiverFilter,
			updateReceiver,
		)
		if err != nil {
			cherr <- err
		}
	}()

	if <-cherr != nil || <-cherr != nil {
		return &pb.FriendReply{
			Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Error Adding Friends"}}
	}
	return &pb.FriendReply{
		Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Success"}}
}

//Removes a friend from a user's list of contacts
func RemoveFriend(request *pb.FriendRequest) *pb.FriendReply {
	dbConnection := helpers.ConnectToDB(helpers.GetUserDatabase())
	userCollection := helpers.ConnectToCollection(dbConnection, helpers.GetAccountInfoCollection())
	//Set a filter for the database to search through
	senderFilter := bson.D{{Key: "email", Value: request.GetSender()}}
	receiverFilter := bson.D{{Key: "email", Value: request.GetReceiver()}}

	//Update Interface, Push Friend Usernames to Friendslist of both users
	cherr := make(chan error, 2)

	go func() {
		updateSender := bson.D{
			{"$pull", bson.D{
				{"Friends", request.GetReceiver()},
			}},
		}
		_, err := userCollection.UpdateOne(
			context.Background(),
			senderFilter,
			updateSender,
		)
		if err != nil {
			cherr <- err
		}
	}()
	go func() {
		updateReceiver := bson.D{
			{"$pull", bson.D{
				{"Friends", request.GetSender()},
			}},
		}
		_, err := userCollection.UpdateOne(
			context.Background(),
			receiverFilter,
			updateReceiver,
		)
		if err != nil {
			cherr <- err
		}
	}()

	if <-cherr != nil || <-cherr != nil {
		return &pb.FriendReply{
			Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Error Removing Friends"}}
	}
	return &pb.FriendReply{
		Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Success"}}
}
