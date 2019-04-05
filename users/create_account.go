package users

import (
	"Naturae_Server/helpers"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "Naturae_Server/naturaeproto"
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
	AccessToken  string
	RefreshToken string
}

//CreateAccount : User want to create an account
func CreateAccount(request *pb.CreateAccountRequest) (NewAccount, *pb.Status) {

	//Connect to the users database
	connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())

	checkStatusChannel := make(chan bool, 4)
	wg := sync.WaitGroup{}

	//Check if email is in a valid format
	go func() {
		checkStatusChannel <- helpers.EmailExist(request.GetEmail(), connectedDB, helpers.GetAccountInfoCollection())
	}()

	//Check if first name is in a valid format
	go func() {
		checkStatusChannel <- helpers.IsNameValid(request.GetFirstName())
	}()

	//Check if last name is in a valid format
	go func() {
		checkStatusChannel <- helpers.IsNameValid(request.GetLastName())
	}()

	//Check if password is in a valid format
	go func() {
		checkStatusChannel <- helpers.IsPasswordValid(request.GetPassword())
	}()

	//Check if the email, firstName, lastName, and password is in a valid format and there no account with the email
	if <-checkStatusChannel && <-checkStatusChannel && <-checkStatusChannel && <-checkStatusChannel {
		//Check checkStatusChannel
		close(checkStatusChannel)

		var accessToken *helpers.AccessToken
		var refreshToken *helpers.RefreshToken

		wg.Add(4)
		go func() {
			defer wg.Done()
			//Generate random bytes of data to be use as salt for the password
			salt := helpers.GenerateRandomBytes(helpers.GetSaltLength())
			//Generate a hash for the user password
			hashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(request.GetPassword()), salt)

			//Create a new user
			newUser := userAccount{Email: request.GetEmail(), FirstName: request.GetFirstName(), LastName: request.GetLastName(),
				Salt: helpers.ConvertByteToStringBase64(salt), Password: helpers.ConvertByteToStringBase64(hashPassword), IsAuthenticated: false}

			//Save the user to the database
			saveNewUser(connectedDB, helpers.GetAccountInfoCollection(), &newUser)
		}()

		//Generate access token
		go func() {
			defer wg.Done()
			accessToken = helpers.GenerateAccessToken(request.GetEmail())
			saveAccessToken(connectedDB, accessToken)

		}()

		//Generate refresh token code
		go func() {
			defer wg.Done()
			refreshToken = helpers.GenerateRefreshToken(request.GetEmail())
			saveRefreshToken(connectedDB, refreshToken)
		}()

		//Generate authentication Code
		go func() {
			defer wg.Done()
			authenCode, expiredTime := helpers.GenerateAuthenCode()
			newAuthenCode := userAuthentication{Email: request.GetEmail(), Code: authenCode, ExpiredTime: expiredTime}
			//Save the user authentication code to the database
			saveAuthenticationCode(connectedDB, helpers.GetAccountAuthentication(), &newAuthenCode)
			//Send the user authentication code to the user's email
			sendAuthenticationCode(request.GetEmail(), request.GetFirstName(), authenCode)
		}()
		wg.Wait()
		//Send the user a welcome message and user authentication number to the provided email address
		log.Println("A new account was created for:", request.GetEmail())
		return NewAccount{AccessToken: accessToken.ID, RefreshToken: refreshToken.ID}, &pb.Status{Code: int32(helpers.GetCreatedStatusCode()),
			Message: "account created"}
	} else {
		//Either email, firstName, lastName, or password is invalid
		return NewAccount{AccessToken: "", RefreshToken: ""}, &pb.Status{Code: int32(helpers.GetInvalidInformation()),
			Message: "information provided is invalid"}
	}

}

//SaveNewUser : Save the user to database
func saveNewUser(database *mongo.Database, collectionName string, user *userAccount) {
	//Connect to the users collection in the database
	accountInfoCollection := helpers.ConnectToCollection(database, collectionName)
	//Save the user into the database
	_, err := accountInfoCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Save user to DB error: ", err)
	} else {
		log.Println("Save", user.Email, "to the account information collection")
	}
}

//SaveAuthenCode : Save authentication code to the database
func saveAuthenticationCode(database *mongo.Database, collectionName string, newAuthenCode *userAuthentication) {
	//Connect to the database collection
	currCollection := helpers.ConnectToCollection(database, collectionName)
	_, err := currCollection.InsertOne(context.Background(), newAuthenCode)
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
