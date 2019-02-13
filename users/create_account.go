package users

import (
	"Naturae_Server/helpers"
	"Naturae_Server/security"
	"errors"
	"fmt"
	"sync"

	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
)

//Create a struct for storing user info
type user struct {
	Email           string
	FirstName       string
	LastName        string
	Salt            string
	Password        string
	IsAuthenticated bool
}

type newAccount struct {
	AccessToken  string
	RefreshToken string
	ErrorCode    int16
	Error        error
}

type TokenLifeSpan struct {
	Year  int
	Month int
	Day   int
}

//Create local global variable for all of the database collection
//associate with create account
var accountInfo = "Account_Information"
var accessToken = "Access_Token"
var refreshToken = "Refresh_Token"

//CreateAccount : User want to create an account
/*
 * email: user email address
 * username: user username
 * password: user password
 *
 */
func CreateAccount(email, firstName, lastName, password *string) newAccount {
	//Connect to the Users database
	currDatabase := helpers.ConnectToDB("Users")

	//Check if the email, firstName, lastName, and password is in a valid format
	if helpers.IsEmailValid(email, currDatabase, &accountInfo) && helpers.IsNameValid(firstName) && helpers.IsNameValid(lastName) &&
		helpers.IsPasswordValid(password) {
		//Generate random bytes of data to be use as salt for the password
		salt, err := security.GenerateRandomBytes(helpers.GetSaltLength())
		if err != nil {
			return newAccount{}
		}
		hashPassword := security.GenerateHash(helpers.ConvertStringToByte(password), salt)

		newUser := user{*email, *firstName, *lastName, *helpers.ConvertByteToString(salt),
			*helpers.ConvertByteToString(hashPassword), false}

		//Set a wait group for multi-threading
		//It will wait for all of the thread process to finish before moving on
		var wg sync.WaitGroup
		saveNewUserToDB(&wg, currDatabase, &newUser)
		wg.Add(1)
		saveAccessTokenToDB(&wg, currDatabase)
		wg.Add(1)
		saveRefreshTokenToDB(&wg, currDatabase)
		wg.Add(1)
		//Wait until all of the go routine to finish
		wg.Wait()
		//Send the user a welcome message and user authentication number to the provided email address
		sendConfirmationEmail(email, firstName)
	} else {
		//Either email, firstName, lastName, or password is invalid

		return newAccount{nil, nil, helpers.GetInvalidArgument(),
			errors.New("invalid input")}
	}

	return newAccount{"", "", helpers.GetCreatedStatusCode(), nil}
}

//Save the user to database
func saveNewUserToDB(wg *sync.WaitGroup, database *mongo.Database, user *user) {
	defer wg.Done()
	//Connect to the users collection in the database
	accountInfoCollection := helpers.ConnectToCollection(database, &accountInfo)
	//Save the user into the database
	accountInfoCollection.InsertOne(context.TODO(), user)

}

//Save access token to the database
func saveAccessTokenToDB(wg *sync.WaitGroup, database *mongo.Database) {
	defer wg.Done()
	saveAccessTokenCollection := helpers.ConnectToCollection(database, &accessToken)

	saveAccessTokenCollection.InsertOne(context.TODO(), "")

}

//Save refresh token to the database
func saveRefreshTokenToDB(wg *sync.WaitGroup, database *mongo.Database) {
	defer wg.Done()
	saveAccessTokenCollection := helpers.ConnectToCollection(database, &refreshToken)
	saveAccessTokenCollection.InsertOne(context.TODO(), "")
}

//Send a confirmation email to the user to make sure it's the user email address
func sendConfirmationEmail(userEmail, firstName *string) {
	//The system will be send a 6 digits number to the user provided email
	//This six digits number will be use to ensure that it's the user email
	body := fmt.Sprintf("Hello, %s\nThanks for creating your Naturae account. To continue, please\n"+
		"verify your email address by entering the following code. %d\nThis code will expire in 30 minutes."+
		"Thank you,\nNature Develper Team", *firstName, security.GenerateRandomNumber(helpers.GetAuthCodeMinNum(),
		helpers.GetAuthCodeMaxNum()))

	fmt.Println(body)

}
