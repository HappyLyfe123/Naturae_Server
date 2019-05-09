package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"bytes"
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//UserInfo create a struct for storing user info in database
type UserInfo struct {
	Email           string
	FirstName       string
	LastName        string
	Salt            string
	Password        string
	IsAuthenticated bool
	ProfileImage    string
	Friends         []string
}

func getUserInfo(email string) (*UserInfo, error) {
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

//SearchUsers searches the database for matching users and returns the list
//@UserSearchRequest - string user, string query
//@UserListReply - repeated string users, Status status
func SearchUsers(request *pb.UserSearchRequest) *pb.UserListReply {
	var searchResult []string
	var friendAvatarList []string
	userEmail := request.GetUser()

	//if userEmail is not nil, it means to retrieve the friendslist of the logged in user
	if len(userEmail) > 0 {
		userAccount, err := getUserInfo(userEmail)
		searchResult = userAccount.Friends
		//For each friend found, retrieve profile image and store
		for i := 0; i < len(searchResult); i++ {
			friendAccount, _ := getUserInfo(searchResult[i])
			friendAvatarList = append(friendAvatarList, friendAccount.ProfileImage)
		}
		if err != nil {
			return &pb.UserListReply{Users: nil, Avatars: nil,
				Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Failed UserEmail Get: " + err.Error()}}
		}
	} else {
		//Search for exactly one user with an inputted email
		queryEmail := request.GetQuery()
		userAccount, err := getUserInfo(queryEmail)
		searchResult = append(searchResult, userAccount.Email)
		friendAvatarList = append(friendAvatarList, userAccount.ProfileImage)
		if err != nil {
			return &pb.UserListReply{Users: nil, Avatars: nil,
				Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "Failed Query Get: " + err.Error()}}
		}

		/*Search via user query input instead, get all users with similar email(not working)
		todo: implement Text-Searching
		*/
	}

	return &pb.UserListReply{Users: searchResult, Avatars: friendAvatarList,
		Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Okay"}}
}

//AddFriend adds a friend to a user's list of contacts
func AddFriend(request *pb.FriendRequest) *pb.FriendReply {
	dbConnection := helpers.ConnectToDB(helpers.GetUserDatabase())
	userCollection := helpers.ConnectToCollection(dbConnection, helpers.GetAccountInfoCollection())
	//Set a filter for the database to search through, acquire the document where the email matches the param value
	senderFilter := bson.D{{Key: "email", Value: request.GetSender()}}
	receiverFilter := bson.D{{Key: "email", Value: request.GetReceiver()}}
	cherr := make(chan error, 2)
	//Update Interface, Push Friend Usernames to Friendslist of both users
	go func() {
		//Add the requested user to the client's friendslist
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

		//If add friends was successful, execute this block
		if err == nil {
			err = CreateConversation(dbConnection, request.GetSender(), request.GetReceiver())
		}
		cherr <- err

	}()
	//Update the document of the other user
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
	}
	return &pb.FriendReply{
		Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Success"}}

}

//RemoveFriend removes a friend from a user's list of contacts
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

		//If remove friend was successful, execute this block
		if err == nil {
			err = RemoveConversation(dbConnection, request.GetSender(), request.GetReceiver())
		}

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
	}
	return &pb.FriendReply{
		Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Success"}}

}

//UpdateProfileImage : update the user info collection in the database with the new profile image url
func UpdateProfileImage(email, imageURL string) *pb.SetProfileImageReply {
	//Set the filter to find the user
	filter := bson.D{{"email", email}}
	//Update the IsAuthenticated field from false to true
	update := bson.D{{"$set", bson.D{{"profileimage", imageURL}}}}
	//Connect to the database and update the user profile
	_, err := helpers.ConnectToDB(helpers.GetUserDatabase()).Collection(helpers.GetAccountInfoCollection()).UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(err)
		return &pb.SetProfileImageReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(),
			Message: "unable to save image"}}
	}
	return &pb.SetProfileImageReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "profile image had been set"}}
}

//SaveProfileImage : save the profile image to AWS S3 bucket
func SaveProfileImage(request *pb.SetProfileImageRequest) (string, bool) {
	saveImage, _ := base64.StdEncoding.DecodeString(request.EncodedImage)
	//Generate a profile image id for the image
	postID := helpers.CreateUUID()
	//Generate a url of were the image is being stored
	imageURL := "https://s3-us-west-2.amazonaws.com/naturae-post-photos/profile-images/" + postID
	params := &s3.PutObjectInput{
		Bucket:        aws.String("naturae-post-photos"),
		Key:           aws.String("profile-images/" + postID),
		Body:          bytes.NewReader(saveImage),
		ContentLength: aws.Int64(int64(len(saveImage))),
		ContentType:   aws.String("image/jpeg"),
	}
	_, err := helpers.GetS3Session().PutObject(params)
	if err != nil {
		log.Printf("Saving profile images to S3 bucket error: %v", err)
		return imageURL, false
	}
	log.Printf("Saving %s to profile images in S3 bucket successful ", postID)
	return imageURL, true
}

func GetProfileImage(email string) *pb.GetProfileImageReply {
	userInfo, err := getUserInfo(email)
	if err != nil {
		log.Printf("Getting user info error in GetProfileImage : %v", err)
		return &pb.GetProfileImageReply{EncodedImage: "", Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "server error"}}
	}
	return &pb.GetProfileImageReply{EncodedImage: userInfo.ProfileImage, Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "success"}}
}
