package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"bytes"
	"go.mongodb.org/mongo-driver/mongo"
)

type userInfo struct {
	Email           string
	FirstName       string
	LastName        string
	Salt            string
	Password        string
	IsAuthenticated bool
}

//Login : Let the user login into their account
func Login(request *pb.LoginRequest) *pb.LoginReply {
	userInfo := helpers.ConnectToDB(helpers.GetUserDatabase())
	databaseResult, err := getLoginInfo(userInfo, request.GetEmail())
	//Database communication error
	if err != nil {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", Status: &pb.Status{
			Code: helpers.GetNotFoundStatusCode(), Message: "No account has been found",
		}}
		//User had not authenticated the account yet
	} else if !databaseResult.IsAuthenticated {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "", Email: "", Status: &pb.Status{
			Code: helpers.GetAccountNotVerifyCode(), Message: "Account is not verify",
		}}
	}

	//Hash the user password
	checkHashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(request.GetPassword()),
		helpers.ConvertStringToByte(databaseResult.Salt))

	//Compare the hash stored in the database and the curr hash password
	if bytes.Compare(helpers.ConvertStringToByte(databaseResult.Password), checkHashPassword) == 1 {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "", Email: "", Status: &pb.Status{
			Code: helpers.GetInvalidLoginCredentialCode(), Message: "Invalid email or password",
		}}

	} else {
		//Get the user access and refresh token id
		accessToken, refreshToken, status := getUserToken(userInfo, request.GetEmail())
		return &pb.LoginReply{AccessToken: accessToken.ID, RefreshToken: refreshToken.ID, FirstName: accessToken.FirstName,
			LastName: accessToken.LastName, Email: request.GetEmail(), Status: status}
	}

}

//Get the user's refresh and access token from the database
func getUserToken(connectedDB *mongo.Database, email string) (*helpers.AccessToken, *helpers.RefreshToken, *pb.Status) {
	accessTokenChanID := make(chan *helpers.AccessToken)
	refreshTokenChanID := make(chan *helpers.RefreshToken)
	errorChan := make(chan bool, 2)
	defer close(accessTokenChanID)
	defer close(refreshTokenChanID)
	defer close(errorChan)

	go func() {
		accessToken, err := helpers.GetAccessToken(connectedDB, email)
		if err != nil {
			accessTokenChanID <- nil
			errorChan <- true
		} else {
			if helpers.IsTokenExpired(accessToken.ExpiredTime) {
				//Create a new access token
				accessToken = helpers.GenerateAccessToken(accessToken.Email, accessToken.FirstName, accessToken.LastName)
				//Save access token to database
				saveAccessToken(connectedDB, accessToken)
			}
			errorChan <- false
			//Save the new token id to the access token id channel
			accessTokenChanID <- accessToken

		}
	}()

	go func() {
		refreshToken, err := helpers.GetRefreshToken(connectedDB, email)
		if err != nil {
			refreshTokenChanID <- nil
			errorChan <- true
		} else {
			//Check if the refresh token had expired already
			//If the current time is before or equal to the expired time,
			//then it will go into the if statement. If it's after the current time is after the expired time, then
			//it will go into the else statement. If the refresh token is expired then the user will have to provide their
			//credential again in order to generate a new refresh token
			if helpers.IsTokenExpired(refreshToken.ExpiredTime) {
				//Create a new refresh token
				refreshToken = helpers.GenerateRefreshToken(email)
				//Save refresh token to database
				saveRefreshToken(connectedDB, refreshToken)
			}
			errorChan <- false
			//Save the new token id to the refresh token id channel
			refreshTokenChanID <- refreshToken
		}
	}()

	//Check if there error occurred when trying to retrieve token from the database
	if <-errorChan || <-errorChan {
		return <-accessTokenChanID, <-refreshTokenChanID, &pb.Status{
			Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error",
		}
	}
	return <-accessTokenChanID, <-refreshTokenChanID, &pb.Status{
		Code: helpers.GetOkStatusCode(), Message: "Login Successful",
	}

}
