package security

import (
	"fmt"

	"../helper"
)

// CheckAppKey : Check if the app key is valid
/*
 *
 *
 */
func CheckAppKey() bool {
	fmt.Println("Work")
	return true
}

// CheckAccessToken : Check is access token is valid
/*
 *
 *
 */
func CheckAccessToken() {
	fmt.Println(helper.GetDeniedStatusCode())
}

//Check if the token is valid
func CheckRefreshToken() {

}

//Generate refresh toek
func GenerateAccessToken() {

}

//Generate access token
func GenerateRefreshToken() {

}

/*
 * Generate an access token and refresh token for user
 *
 */
func GenerateToken() {

}
