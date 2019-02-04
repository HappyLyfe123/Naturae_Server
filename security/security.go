package security

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"

	"../helpers"
)

// CheckAppKey : Check if the app key is valid
func CheckAppKey() bool {
	fmt.Println("Work")
	return true
}

// CheckAccessToken : Check is access token is valid
func CheckAccessToken() {
	fmt.Println(helpers.GetDeniedStatusCode())
}

//CheckRefreshToken : Check if the token is valid
func CheckRefreshToken() {

}

//GenerateAccessToken : Generate an access token and refresh token for user
func GenerateAccessToken() (time.Time, time.Time) {
	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, 1)

	return startTime, endTime
}

//GenerateRefreshToken : Generate Refresh token
func GenerateRefreshToken() {

}

//Hash password with salt using sha512
func HashPassword(userPassword, salt string) string {
	hasher := sha512.New()
	password := helpers.ConvertStringToByte(salt + userPassword)
	return base64.StdEncoding.EncodeToString(hasher.Sum(password))
}
