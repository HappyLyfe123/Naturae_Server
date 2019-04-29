package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"bytes"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type UserInfo struct {
	Email           string
	FirstName       string
	LastName        string
	Salt            string
	Password        string
	Friends         []string
	IsAuthenticated bool
}

//Login : Let the user login into their account
func Login(request *pb.LoginRequest) *pb.LoginReply {
	databaseResult, err := getLoginInfo(request.GetEmail())
	//Database communication error
	if err != nil {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "", Email: "", Status: &pb.Status{
			Code: helpers.GetNotFoundStatusCode(), Message: "No account has been found"}}
		//User had not authenticated the account yet
	} else if !databaseResult.IsAuthenticated {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "", Email: "", Status: &pb.Status{
			Code: helpers.GetAccountNotVerifyCode(), Message: "Account is not verify"}}
	}

	//Hash the user password
	checkHashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(request.GetPassword()),
		helpers.ConvertStringToByte(databaseResult.Salt))

	//Compare the hash stored in the database and the curr hash password
	if bytes.Compare(helpers.ConvertStringToByte(databaseResult.Password), checkHashPassword) == -1 {
		return &pb.LoginReply{AccessToken: "", RefreshToken: "", FirstName: "", LastName: "", Email: "", Status: &pb.Status{
			Code: helpers.GetInvalidLoginCredentialCode(), Message: "Invalid email or password"}}
	} else {
		//Get the user access and refresh token id
		accessToken, refreshToken, status := getUserToken(request.GetEmail())
		if len(accessToken.ID) == 0 || len(refreshToken.ID) == 0 {
			return &pb.LoginReply{AccessToken: "", RefreshToken: "", FirstName: accessToken.FirstName,
				LastName: accessToken.LastName, Email: "", Status: status}
		}
		return &pb.LoginReply{AccessToken: accessToken.ID, RefreshToken: refreshToken.ID, FirstName: accessToken.FirstName,
			LastName: accessToken.LastName, Email: request.GetEmail(), Status: status}
	}

}

//Get the user's refresh and access token from the database
func getUserToken(email string) (*helpers.AccessToken, *helpers.RefreshToken, *pb.Status) {
	userDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	accessTokenChan := make(chan *helpers.AccessToken)
	refreshTokenChan := make(chan *helpers.RefreshToken)
	errorChan := make(chan bool, 2)
	defer close(accessTokenChan)
	defer close(refreshTokenChan)
	defer close(errorChan)
	go func() {
		accessToken, err := helpers.GetAccessTokenEmail(userDB, email)
		if err != nil {
			log.Printf("Login getting access token error: %v", err)
			errorChan <- true
			accessTokenChan <- &helpers.AccessToken{Email: "", FirstName: "", LastName: "", ID: "", ExpiredTime: time.Now()}

		} else {
			if helpers.IsTokenExpired(accessToken.ExpiredTime) {
				//Create a new access token
				accessToken = helpers.GenerateAccessToken(accessToken.Email, accessToken.FirstName, accessToken.LastName)
				//Save access token to database
				saveAccessToken(userDB, accessToken)
			}
			errorChan <- false
			//Save the new token id to the access token id channel
			accessTokenChan <- accessToken

		}
	}()

	go func() {
		refreshToken, err := helpers.GetRefreshToken(userDB, email)
		if err != nil {
			log.Printf("Login getting refresh token error: %v", err)
			errorChan <- true
			refreshTokenChan <- &helpers.RefreshToken{Email: "", ID: "", ExpiredTime: time.Now()}

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
				saveRefreshToken(userDB, refreshToken)
			}
			errorChan <- false
			//Save the new token id to the refresh token id channel
			refreshTokenChan <- refreshToken
		}
	}()
	//Check if there error occurred when trying to retrieve token from the database
	if <-errorChan || <-errorChan {
		return <-accessTokenChan, <-refreshTokenChan, &pb.Status{
			Code: helpers.GetInternalServerErrorStatusCode(), Message: "Server error"}
	}
	return <-accessTokenChan, <-refreshTokenChan, &pb.Status{
		Code: helpers.GetOkStatusCode(), Message: "Login Successful"}

}
