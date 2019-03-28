package users

import (
	"Naturae_Server/helpers"
	"bytes"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type loginInfo struct {
	Salt            string
	Password        string
	IsAuthenticated bool
}

type loginResponse struct {
	AccessToken  string
	RefreshToken string
}

//Login : Let the user login into their account
func Login(email, password string) (loginResponse, error) {
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

		return loginResponse{AccessToken: "", RefreshToken: ""}, errors.New(helpers.ErrorFormat(helpers.GetAccountNotVerifyCode(),
			"login invalid", "user haven't verify account"))
	} else {
		//Hash the provide password with the stored salt from the database
		checkPasswordHash := helpers.GenerateHash(helpers.ConvertStringToByte(password), helpers.ConvertStringToByte(databaseResult.Salt))
		//Check if the two hash password match
		if bytes.Compare(helpers.ConvertStringToByte(databaseResult.Password), checkPasswordHash) == 1 {
			accessToken, err := helpers.GetAccessToken(userInfo, email)
			if err != nil {
				log.Println("Getting access token error: ", err)
				return loginResponse{AccessToken: "", RefreshToken: ""}, errors.New(helpers.ErrorFormat(helpers.GetInternalServerErrorStatusCode(),
					"server error", "internal server error"))
			}
			refreshToken, err := helpers.GetRefreshToken(userInfo, email)
			if err != nil {
				log.Println("Getting refresh token error: ", err)
				return loginResponse{AccessToken: "", RefreshToken: ""},
					errors.New(helpers.ErrorFormat(helpers.GetInternalServerErrorStatusCode(),
						"server error", "internal server error"))
			}
			return loginResponse{AccessToken: accessToken.ID, RefreshToken: refreshToken.ID}, nil
		}
	}
	return loginResponse{AccessToken: "", RefreshToken: ""}, errors.New(helpers.ErrorFormat(helpers.GetInvalidLoginCredentialCode(),
		"invalid login", "user email or password is invalid"))
}
