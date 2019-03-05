package users

import (
	"Naturae_Server/helpers"
	"strings"
)

//Check if the authentication code is valid
func AuthenticateAccount(email, userCode string) {
	connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	databaseResult, err := getAuthenCode(connectedDB, email)
	if err != nil {

	}
	//Check if the authentication code is expired
	if helpers.IsTimeValid(databaseResult.ExpiredTime) {
		//Check if the user provided authentication code match
		if strings.Compare(databaseResult.Code, userCode) == 1 {
			upDateUserAuthenStatus(email)
			removeAuthenCode(email)
		} else {

		}
	}
}

//GenerateNewAuthenCode: Generate a new code for the user if the time expired
func GenerateNewAuthenCode(email string) {

}

//Change the user authentication status
func upDateUserAuthenStatus(email string) {

}

//Remove the authentication code from the server
func removeAuthenCode(email string) {

}
