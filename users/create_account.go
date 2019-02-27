package users

import (
	"Naturae_Server/helpers"
	"Naturae_Server/security"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
)

//NewAccount : create new account structure
type NewAccount struct {
	AccessToken  string
	RefreshToken string
	ErrorList    map[int16]string
}

var saveAttemptLimit = 10

//CreateAccount : User want to create an account
//email: user email address
//username: user username
//password: user password
func CreateAccount(email, firstName, lastName, password string) NewAccount {

	//Connect to the users database
	connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	errorList := make(map[int16]string)
	//Set a wait group for multi-threading
	//It will wait for all of the thread process to finish before moving on
	var wg sync.WaitGroup

	isEmailValid, err := helpers.IsEmailValid(email)
	errorList[err.ErrorCode] = err.ErrorMessage
	isEmailExist, err := helpers.EmailExist(email, connectedDB, helpers.GetAccountInfoCollection())
	errorList[err.ErrorCode] = err.ErrorMessage
	isFirstNameValid, err := helpers.IsNameValid(firstName)
	errorList[err.ErrorCode] = err.ErrorMessage
	isLastNameValid, err := helpers.IsNameValid(lastName)
	errorList[err.ErrorCode] = err.ErrorMessage
	isPasswordValid, err := helpers.IsPasswordValid(password)
	errorList[err.ErrorCode] = err.ErrorMessage

	//Check if the email, firstName, lastName, and password is in a valid format and there no account with the email
	if isEmailValid && !isEmailExist && isFirstNameValid && isLastNameValid && isPasswordValid {
		//Create a channel for token id and start time
		tokenIDChan := make(chan string)
		startTimeChan := make(chan time.Time)
		//Close the channel
		defer close(tokenIDChan)
		defer close(startTimeChan)

		//Generate random bytes of data to be use as salt for the password
		salt := security.GenerateRandomBytes(helpers.GetSaltLength())
		//Generate a hash for the user password
		hashPassword := security.GenerateHash(helpers.ConvertStringToByte(password), salt)

		//Create a new user
		newUser := userAccount{Email: email, First_Name: firstName, Last_Name: lastName, Salt: helpers.ConvertByteToString(salt),
			Password: helpers.ConvertByteToString(hashPassword), Is_Authenticated: false}

		//Generate access token and set it to have a life span of one day
		go security.GenerateToken(email, 0, tokenIDChan, startTimeChan) //0 For access token
		accessToken := token{email, <-tokenIDChan, <-startTimeChan}

		//Generate refresh token and set it to have a life span of two months
		go security.GenerateToken(email, 1, tokenIDChan, startTimeChan) // 1 For refresh token
		refreshToken := token{email, <-tokenIDChan, <-startTimeChan}

		//Generate authentication Code
		generatedCode, startTime := security.GenerateAuthenCode()
		newAuthenCode := userAuthentication{email, generatedCode, startTime}

		//Save the user to the database
		go saveNewUser(&wg, connectedDB, helpers.GetAccountInfoCollection(), &newUser)
		wg.Add(1)
		//Generate and save access token for the user
		go saveToken(&wg, connectedDB, helpers.GetAccessTokenCollection(), &accessToken)
		wg.Add(1)
		////Generate and save refresh token for the user
		go saveToken(&wg, connectedDB, helpers.GetRefreshTokenCollection(), &refreshToken)
		wg.Add(1)
		//Save authentication Code
		go saveAuthenticationCode(&wg, connectedDB, helpers.GetAccountVerification(), &newAuthenCode)
		wg.Add(1)
		//Wait until all of the go routine to finish
		wg.Wait()
		//Send the user a welcome message and user authentication number to the provided email address
		sendAuthenticationCode(email, firstName, generatedCode)
		fmt.Println("A new account was created for: ", email)
		return NewAccount{AccessToken: accessToken.Token_ID, RefreshToken: refreshToken.Token_ID, ErrorList: nil}
	} else {
		//Either email, firstName, lastName, or password is invalid
		return NewAccount{AccessToken: "", RefreshToken: "", ErrorList: errorList}
	}

}

//SaveNewUser : Save the user to database
func saveNewUser(wg *sync.WaitGroup, database *mongo.Database, collectionName string, user *userAccount) {
	defer wg.Done()
	//If there an error, it will attempt to save the info to the database until the limit is reach
	for numAttempt := 1; numAttempt < saveAttemptLimit; numAttempt++ {
		//Connect to the users collection in the database
		accountInfoCollection := helpers.ConnectToCollection(database, collectionName)
		//Save the user into the database
		_, err := accountInfoCollection.InsertOne(context.TODO(), user)
		if err != nil {
			log.Println("Save user to DB error: ", err)
		} else {
			log.Println("Save ", user.Email, " to the DB")
			break
		}
	}
}

//SaveAuthenCode : Save authentication code to the database
func saveAuthenticationCode(wg *sync.WaitGroup, database *mongo.Database, collectionName string, newAuthenCode *userAuthentication) {
	defer wg.Done()
	//If there an error, it will attempt to save the info to the database until the limit is reach
	for attemptNum := 1; attemptNum < saveAttemptLimit; attemptNum++ {
		//Connect to the database collection
		currCollection := helpers.ConnectToCollection(database, collectionName)
		_, err := currCollection.InsertOne(context.TODO(), newAuthenCode)
		if err != nil {
			log.Println("Save authentication to DB error: ", err)
		} else {
			log.Println("Save ", newAuthenCode.Email, " authentication code to DB")
			//Break out of the for loop
			break
		}
	}
}

//SendAuthenticationEmail : Send a confirmation email to the user to make sure it's the user email address
func sendAuthenticationCode(userEmail, firstName string, authenCode string) {
	//It going to try to keep authentication code if there an error until the limit is reach
	for numOfAttempt := 0; numOfAttempt < saveAttemptLimit; numOfAttempt++ {
		//The system will be send a 6 digits number to the user provided email
		//This six digits number will be use to ensure that it's the user email
		body := fmt.Sprintf("Hello %s,\nThanks for creating account with Naturae.\n"+
			"Please enter the secure verification code: %s\nThis code will expire in 30 minutes."+
			"\nThank you,\nNature Develper Team", firstName, authenCode)
		//Send the email to the user
		newMail := helpers.Email{Receiver: userEmail, Subject: "Account Authentication", Body: body}
		err := helpers.SendEmail(&newMail)
		//If there no error
		if err == nil {
			log.Println("Email ", userEmail, " authentication code")
			break
		}
		log.Println("Send email error: ", err)
	}

}
