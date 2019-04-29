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
	IsAuthenticated bool
}

func getLoginInfo(email string) (*UserInfo, error) {
	var result UserInfo
	userInfoDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	filter := bson.D{{Key: "email", Value: email}}
	//Connect to the collection database
	accountInfoCollection := userInfoDB.Collection(helpers.GetAccountInfoCollection())
	//Make a request to the database
	err := accountInfoCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Printf("Getting user info error: %v", err)
		return &UserInfo{}, err
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
	newAccessTokenID := ""
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
			accessToken := helpers.GenerateAccessToken("", "", "")
			//Set the filter to find the user
			filter := bson.D{{"email", refreshTokenResult.Email}}
			//Update the access token id and expired time in the database with the newly generated one
			update := bson.D{{"$set", bson.D{{"id", accessToken.ID}, {"expiredtime", accessToken.ExpiredTime}}}}
			_, err := helpers.ConnectToCollection(currConnectedDB, helpers.GetAccessTokenCollection()).UpdateOne(context.Background(), filter, update)
			if err != nil {
				status = &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error"}
			} else {
				newAccessTokenID = accessToken.ID
				status = &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Code has been generated"}
			}
		} else {
			status = &pb.Status{Code: helpers.GetInvalidTokenCode(), Message: "Token is invalid"}
		}
	}

	return &pb.GetAccessTokenReply{AccessToken: newAccessTokenID, Status: status}
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
		userAccount, err := getLoginInfo(dbConnection, userEmail)
		searchResult = userAccount.Friends
		if err != nil {
			return &pb.UserListReply{Users: nil,
				Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Failed UserEmail Get: " + err.Error()}}
		}
	} else {
		//Search for exactly one user with an inputted email
		queryEmail := request.GetQuery()
		userAccount, err := getLoginInfo(dbConnection, queryEmail)
		searchResult = userAccount.Friends
		if err != nil {
			return &pb.UserListReply{Users: nil,
				Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Failed Query Get: " + err.Error()}}
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

	return &pb.UserListReply{Users: searchResult,
		Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Okay"}}
}

//Adds a friend to a user's list of contacts
func AddFriend(request *pb.FriendRequest) *pb.FriendReply {
	dbConnection := helpers.ConnectToDB(helpers.GetUserDatabase())
	userCollection := helpers.ConnectToCollection(dbConnection, helpers.GetAccountInfoCollection())
	//Set a filter for the database to search through, acquire the document where the email matches the param value
	senderFilter := bson.D{{Key: "email", Value: request.GetSender()}}
	receiverFilter := bson.D{{Key: "email", Value: request.GetReceiver()}}

	cherr := make(chan error, 2)
	//Update Interface, Push Friend Usernames to Friendslist of both users
	go func() {
		updateSender := bson.D{
			{"$push", bson.D{
				{"friends", request.GetReceiver()},
			}},
		}
		_, err := userCollection.UpdateOne(
			context.Background(),
			senderFilter,
			updateSender,
		)
		cherr <- err

	}()
	go func() {
		updateReceiver := bson.D{
			{"$push", bson.D{
				{"friends", request.GetSender()},
			}},
		}
		_, err := userCollection.UpdateOne(
			context.Background(),
			receiverFilter,
			updateReceiver,
		)
		cherr <- err

	}()

	if <-cherr != nil || <-cherr != nil {
		return &pb.FriendReply{
			Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Error Adding Friends"}}
	} else {
		return &pb.FriendReply{
			Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Success"}}
	}

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
				{"friends", request.GetReceiver()},
			}},
		}
		_, err := userCollection.UpdateOne(
			context.Background(),
			senderFilter,
			updateSender,
		)
		cherr <- err
	}()
	go func() {
		updateReceiver := bson.D{
			{"$pull", bson.D{
				{"friends", request.GetSender()},
			}},
		}
		_, err := userCollection.UpdateOne(
			context.Background(),
			receiverFilter,
			updateReceiver,
		)
		cherr <- err

	}()

	if <-cherr != nil || <-cherr != nil {
		return &pb.FriendReply{
			Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Error Removing Friends"}}
	} else {
		return &pb.FriendReply{
			Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Success"}}
	}

}
