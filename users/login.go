package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"bytes"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type loginInfo struct {
	Salt            string
	Password        string
	IsAuthenticated bool
}

//Login : Let the user login into their account
func Login(request *pb.LoginRequest) *pb.LoginReply {
	userInfo := helpers.ConnectToDB(helpers.GetUserDatabase())
	databaseResult, err := getLoginInfo(userInfo, request.GetEmail())
	if err != nil {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", Status: &pb.Status{
			Code: helpers.GetNotFoundStatusCode(), Message: "No account has been found",
		}}
	} else if !databaseResult.IsAuthenticated {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", Status: &pb.Status{
			Code: helpers.GetAccountNotVerifyCode(), Message: "Account is not verify",
		}}
	}

	checkHashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(request.GetPassword()),
		helpers.ConvertStringToByte(databaseResult.Salt))

	if bytes.Compare(helpers.ConvertStringToByte(databaseResult.Password), checkHashPassword) == 1 {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", Status: &pb.Status{
			Code: helpers.GetInvalidLoginCredentialCode(), Message: "Invalid email or password",
		}}

	} else {
		//Get the user access and refresh token id
		accessTokenID, refreshTokenID, status := getUserToken(userInfo, request.GetEmail())
		return &pb.LoginReply{AccessToken: accessTokenID, RefreshToken: refreshTokenID, Status: status}
	}

}

func getUserToken(connectedDB *mongo.Database, email string) (string, string, *pb.Status) {
	var wg sync.WaitGroup
	accessTokenChanID := make(chan string)
	refreshTokenChanID := make(chan string)
	errorChan := make(chan error, 2)
	wg.Add(2)
	go func(wg sync.WaitGroup) {
		defer wg.Done()
		accessToken, err := helpers.GetAccessToken(connectedDB, email)
		if err != nil {
			accessTokenChanID <- ""
			errorChan <- err
		} else {
			if helpers.IsTokenExpired(accessToken.ExpiredTime) {
				//Create a new access token
				accessToken = helpers.GenerateAccessToken(email)
				//Save access token to database
				saveAccessToken(connectedDB, accessToken)
			}
			//Save the new token id to the access token id channel
			accessTokenChanID <- accessToken.ID

		}
	}(wg)

	go func(wg sync.WaitGroup) {
		defer wg.Done()
		refreshToken, err := helpers.GetRefreshToken(connectedDB, email)
		if err != nil {
			refreshTokenChanID <- ""
			errorChan <- err
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
			//Save the new token id to the refresh token id channel
			refreshTokenChanID <- refreshToken.ID
		}
	}(wg)

	wg.Wait()

	//There an error occurred
	if len(errorChan) != 0 {
		return <-accessTokenChanID, <-refreshTokenChanID, &pb.Status{
			Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error",
		}
	}
	return <-accessTokenChanID, <-refreshTokenChanID, &pb.Status{
		Code: helpers.GetOkStatusCode(), Message: "Login Successful",
	}

}
