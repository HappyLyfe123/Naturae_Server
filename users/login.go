package users

import (
	"Naturae_Server/helpers"
	"bytes"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type loginInfo struct {
	Salt            string
	Password        string
	IsAuthenticated bool
}

type loginResponse struct {
	Success      bool
	AccessToken  string
	RefreshToken string
}

//Login : Let the user login into their account
func Login(email, password string) (loginResponse, helpers.AppError) {
	userInfo := helpers.ConnectToDB(helpers.GetUserDatabase())
	databaseResult, err := getLoginInfo(userInfo, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Getting login info in login error:", err)
		} else {
			log.Println("Getting login info format error: ", err)
		}
		//Check if the user already authenticated the account
	} else if databaseResult.IsAuthenticated == false {

		return loginResponse{Success: false, AccessToken: "", RefreshToken: ""}, helpers.AppError{Code: helpers.GetAccountNotVerifyCode(),
			Type: "Account Verification", Description: "User account was not verify"}
	} else {
		//Hash the provide password with the stored salt from the database
		checkPasswordHash := helpers.GenerateHash(helpers.ConvertStringToByte(password), helpers.ConvertStringToByte(databaseResult.Salt))
		//Check if the two hash password match
		if bytes.Compare(helpers.ConvertStringToByte(databaseResult.Password), checkPasswordHash) == 1 {
			accessToken, err := helpers.GetAccessToken(userInfo, email)
			if err != nil {
				log.Println("Getting access token error: ", err)
				return loginResponse{Success: false, AccessToken: "", RefreshToken: ""},
					helpers.AppError{Code: helpers.GetInternalServerErrorStatusCode(), Type: "Server error",
						Description: "Internal server error"}
			}
			refreshToken, err := helpers.GetRefreshToken(userInfo, email)
			if err != nil {
				log.Println("Getting refresh token error: ", err)
				return loginResponse{Success: false, AccessToken: "", RefreshToken: ""},
					helpers.AppError{Code: helpers.GetInternalServerErrorStatusCode(), Type: "Server error",
						Description: "Internal server error"}
			}
			return loginResponse{Success: true, AccessToken: accessToken.ID, RefreshToken: refreshToken.ID}, helpers.AppError{}
		}
	}
	return loginResponse{Success: false, AccessToken: "", RefreshToken: ""},
		helpers.AppError{Code: helpers.GetInvalidLoginCredentialCode(), Type: "Invalid login credential", Description: "Invalid " +
			"email or password"}
}
