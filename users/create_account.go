package users

import (
	"Naturae_Server/helpers"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

//create authentication structure for user
type userAuthentication struct {
	Email       string
	Code        string
	ExpiredTime time.Time
}

//NewAccount : create new account structure
type NewAccount struct {
	Success      bool
	AccessToken  string
	RefreshToken string
}

//CreateAccount : User want to create an account
func CreateAccount(email, firstName, lastName, password string) (NewAccount, []helpers.AppError) {

	//Connect to the users database
	connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	var errorList []helpers.AppError
	//Set a wait group for multi-threading
	//It will wait for all of the thread process to finish before moving on
	var wg sync.WaitGroup

	//Check if the provided information meet the app requirement
	isEmailValid, err := helpers.IsEmailValid(email)
	errorList = append(errorList, err)
	isEmailExist, err := helpers.EmailExist(email, connectedDB, helpers.GetAccountInfoCollection())
	errorList = append(errorList, err)
	isFirstNameValid, err := helpers.IsNameValid(firstName)
	errorList = append(errorList, err)
	isLastNameValid, err := helpers.IsNameValid(lastName)
	errorList = append(errorList, err)
	isPasswordValid, err := helpers.IsPasswordValid(password)
	errorList = append(errorList, err)

	//Check if the email, firstName, lastName, and password is in a valid format and there no account with the email
	if isEmailValid && !isEmailExist && isFirstNameValid && isLastNameValid && isPasswordValid {
		//Create a channel for token id and start time
		tokenIDChan := make(chan string)
		endTime := make(chan time.Time)
		//Close the channel
		defer close(tokenIDChan)
		defer close(endTime)

		//Generate random bytes of data to be use as salt for the password
		salt := helpers.GenerateRandomBytes(helpers.GetSaltLength())
		//Generate a hash for the user password
		hashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(password), salt)

		//Create a new user
		newUser := userAccount{Email: email, FirstName: firstName, LastName: lastName, Salt: helpers.ConvertByteToStringBase64(salt),
			Password: helpers.ConvertByteToStringBase64(hashPassword), IsAuthenticated: false}

		//Generate access token and set it to have a life span of 6 hours
		accessToken := helpers.GenerateAccessToken(email)
		refreshToken := helpers.GenerateRefreshToken(email)

		//Generate authentication Code
		generatedCode, startTime := helpers.GenerateAuthenCode()
		newAuthenCode := userAuthentication{email, generatedCode, startTime}

		//Save the user to the database
		go saveNewUser(&wg, connectedDB, helpers.GetAccountInfoCollection(), &newUser)
		wg.Add(1)

		//Generate and save access token for the user
		go saveAccessToken(&wg, connectedDB, &accessToken)
		wg.Add(1)

		//Generate and save refresh token for the user
		go saveRefreshToken(&wg, connectedDB, &refreshToken)
		wg.Add(1)

		//Save authentication Code
		go saveAuthenticationCode(&wg, connectedDB, helpers.GetAccountAuthentication(), &newAuthenCode)
		wg.Add(1)

		//Wait until all of the go routine to finish
		wg.Wait()

		//Send the user a welcome message and user authentication number to the provided email address
		sendAuthenticationCode(email, firstName, generatedCode)
		log.Println("A new account was created for:", email)
		return NewAccount{Success: true, AccessToken: accessToken.ID, RefreshToken: refreshToken.ID}, errorList
	} else {
		//Either email, firstName, lastName, or password is invalid
		return NewAccount{Success: false, AccessToken: "", RefreshToken: ""}, errorList
	}

}

//SaveNewUser : Save the user to database
func saveNewUser(wg *sync.WaitGroup, database *mongo.Database, collectionName string, user *userAccount) {
	defer wg.Done()
	//Connect to the users collection in the database
	accountInfoCollection := helpers.ConnectToCollection(database, collectionName)
	//Save the user into the database
	_, err := accountInfoCollection.InsertOne(nil, user)
	if err != nil {
		log.Println("Save user to DB error: ", err)
	} else {
		log.Println("Save", user.Email, "to the account information collection")
	}
}

//SaveAuthenCode : Save authentication code to the database
func saveAuthenticationCode(wg *sync.WaitGroup, database *mongo.Database, collectionName string, newAuthenCode *userAuthentication) {
	defer wg.Done()
	//Connect to the database collection
	currCollection := helpers.ConnectToCollection(database, collectionName)
	_, err := currCollection.InsertOne(nil, newAuthenCode)
	if err != nil {
		log.Println("Save authentication to DB error: ", err)
	} else {
		log.Println("Save", newAuthenCode.Email, "to authentication code to DB")
		//Break out of the for loop
	}
}

//SendAuthenticationEmail : Send a confirmation email to the user to make sure it's the user email address
func sendAuthenticationCode(userEmail, firstName string, authenCode string) {
	//The system will be send a 6 digits number to the user provided email
	//This six digits number will be use to ensure that it's the user email
	body := fmt.Sprintf("Hello %s,\nThanks for creating account with Naturae.\n"+
		"Please enter the secure verification code: %s\nThis code will expire in 30 minutes."+
		"\nThank you,\nNature Develper Team", firstName, authenCode)
	//Send the email to the user
	newMail := helpers.Email{Receiver: userEmail, Subject: "Account Authentication", Body: body}
	err := helpers.SendEmail(&newMail)
	//If there no error
	if err != nil {
		log.Println("Send email error: ", err)
	}
	log.Println("Email", userEmail, "authentication code")

}
