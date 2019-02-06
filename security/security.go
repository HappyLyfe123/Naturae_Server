package security

import (
	"Naturae_Server/helpers"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// CheckAppKey : Check if the app key is valid
// @appKey : provided app key to be check
func CheckAppKey(appKey string) bool {
	if strings.Compare(appKey, "") == 0 {
		return true
	} else {
		return false
	}
}

//CheckAccessToken : Check is access token is valid
func CheckAccessToken(username, tokenID string) {
	fmt.Println(helpers.GetDeniedStatusCode())
}

//CheckRefreshToken : Check if the token is valid
func CheckRefreshToken(username, tokenID string) {

}

//GenerateAccessToken : Generate an access token and refresh token for user
func GenerateAccessToken() (time.Time, time.Time) {
	//Generate a start time for the token
	startTime := time.Now()
	//Generate an end time for the token
	//AccessToken have a 1 day lifespan
	endTime := startTime.AddDate(0, 0, 1)
	return startTime, endTime
}

//GenerateRefreshToken : Generate Refresh token
func GenerateRefreshToken() (*time.Time, *time.Time) {
	//Generate an start time for the refresh token
	startTime := time.Now()
	//Generate an end time for the refresh token
	//Refresh token have 1 month life span
	endTime := startTime.AddDate(0, 1, 0)

	return &startTime, &endTime
}

//HashPassword : Hash password with salt using sha512
func HashPassword(userPassword, salt string) string {
	//Add the salt to the user password
	//After adding salt to user password convert the password with
	//salt into byte array
	password := helpers.ConvertStringToByte(salt + userPassword)
	//Hash the byte array password using sha515 then convert the hash to
	//base64 encoding and return it
	return base64.StdEncoding.EncodeToString(sha512.New().Sum(password))
}
