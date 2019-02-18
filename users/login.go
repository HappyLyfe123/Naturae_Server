package users

import "Naturae_Server/helpers"

//Login : Let the user login into their account
func Login(email *string) {
	//Connect to the Users database
	currDatabase := helpers.ConnectToDB("Users")
	currUser := helpers.FindUser(email, currDatabase, helpers.GetAccountInfoColl())
	if currUser != nil {

	}

}
