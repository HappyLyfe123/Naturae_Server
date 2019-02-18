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

//NewAccount : structure
type newAccount struct {
	AccessToken  string
	RefreshToken string
	ErrorCode    int16
	Error        error
}

//CreateAccount : User want to create an account
//email: user email address
//username: user username
//password: user password
func CreateAccount(email, firstName, lastName, password *string) newAccount {
	//Connect to the Users database
	currDatabase := helpers.ConnectToDB("Users")
	//Set a wait group for multi-threading
	//It will wait for all of the thread process to finish before moving on
	var wg sync.WaitGroup

	//Check if the email, firstName, lastName, and password is in a valid format
	if helpers.IsEmailValid(email, currDatabase, helpers.GetAccountInfoCollection()) && helpers.IsNameValid(firstName) && helpers.IsNameValid(lastName) &&
		helpers.IsPasswordValid(password) {
		//Generate random bytes of data to be use as salt for the password
		salt, err := security.GenerateRandomBytes(helpers.GetSaltLength())
		if err != nil {
			return newAccount{}
		}
		//Generate a hash for the user password
		hashPassword := security.GenerateHash(helpers.ConvertStringToByte(password), salt)

		//Create a new user
		newUser := user{*email, *firstName, *lastName, *helpers.ConvertByteToString(salt),
			*helpers.ConvertByteToString(hashPassword), false}

		//Generate access token and set it to have a life span of one day
		accessToken := security.GenerateToken(email, []int{0, 0, 1})
		//Generate refresh token and set it to have a life span of two months
		refreshToken := security.GenerateToken(email, []int{0, 2, 0})

		//Save the user to the database
		go saveNewUserToDB(&wg, currDatabase, helpers.GetAccountInfoCollection(), &newUser)
		wg.Add(1)
		//Generate and save access token for the user
		go saveTokenToDB(&wg, currDatabase, helpers.GetAccessTokenCollection(), accessToken)
		wg.Add(1)
		//Generate and save refresh token for the user
		go saveTokenToDB(&wg, currDatabase, helpers.GetRefreshTokenCollection(), refreshToken)
		wg.Add(1)
		//Wait until all of the go routine to finish
		wg.Wait()
		//Send the user a welcome message and user authentication number to the provided email address
		SendAuthenticationEmail(email, firstName)
	} else {
		//Either email, firstName, lastName, or password is invalid
		return newAccount{"", "", helpers.GetInvalidArgument(),
			errors.New("invalid input")}
	}

	return newAccount{"", "", helpers.GetCreatedStatusCode(), nil}
}

//Save the user to database
func saveNewUserToDB(wg *sync.WaitGroup, database *mongo.Database, collectionName string, user *user) {
	defer wg.Done()
	//Connect to the users collection in the database
	accountInfoCollection := helpers.ConnectToCollection(database, &collectionName)
	//Save the user into the database
	accountInfoCollection.InsertOne(context.TODO(), user)

}

//Save access token to the database
func saveTokenToDB(wg *sync.WaitGroup, database *mongo.Database, collectionName string, token *security.Token) {
	defer wg.Done()
	//Connect to the database collection
	currCollection := helpers.ConnectToCollection(database, &collectionName)
	//Save token to the database
	currCollection.InsertOne(context.TODO(), *token)

}

//SendAuthenticationEmail : Send a confirmation email to the user to make sure it's the user email address
func SendAuthenticationEmail(userEmail, firstName *string) {
	//The system will be send a 6 digits number to the user provided email
	//This six digits number will be use to ensure that it's the user email
	body := fmt.Sprintf("Hello %s,\nThanks for creating account with Naturae.\n"+
		"Please enter the secure verification code: %d\nThis code will expire in 30 minutes."+
		"\nThank you,\nNature Develper Team", *firstName, security.GenerateRandomNumber(helpers.GetAuthCodeMinNum(),
		helpers.GetAuthCodeMaxNum()))
	//Send the email to the user
	err := helpers.SendEmail(&helpers.Email{*userEmail, "Account Authentication", body})
	if err != nil {
		fmt.Println(err)
	}
}
