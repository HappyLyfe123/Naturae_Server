package users

import (
	"Naturae_Server/helpers"
	"context"
	"fmt"
	"log"
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

//CreateAccount : User want to create an account
func CreateAccount(request *pb.CreateAccountRequest) *pb.CreateAccountReply {

	//Connect to the users database
	connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())

	//Create a channel for storing the validity result from checking user input
	checkDataChannel := make(chan bool, 3)
	//Close the channel
	defer close(checkDataChannel)

	if !helpers.IsEmailValid(request.GetEmail()) || helpers.EmailExist(request.GetEmail(), connectedDB, helpers.GetAccountInfoCollection()) {
		return &pb.CreateAccountReply{Status: &pb.Status{Code: int32(helpers.GetEmailExistCode()),
			Message: "email already taken"}}
	}

	//Check if first name is in a valid format
	go func() {
		checkDataChannel <- helpers.IsNameValid(request.GetFirstName())
	}()

	//Check if last name is in a valid format
	go func() {
		checkDataChannel <- helpers.IsNameValid(request.GetLastName())
	}()

	//Check if password is in a valid format
	go func() {
		checkDataChannel <- helpers.IsPasswordValid(request.GetPassword())
	}()

	//Check if the email, firstName, lastName, and password is in a valid format and there no account with the email
	if <-checkDataChannel && <-checkDataChannel && <-checkDataChannel {

		//Generate random bytes of data to be use as salt for the password
		salt := helpers.GenerateRandomBytes(helpers.GetSaltLength())
		//Generate a hash for the user password
		hashPassword := helpers.GenerateHash(helpers.ConvertStringToByte(request.GetPassword()), salt)
		//Create a new user
		newUser := UserInfo{Email: request.GetEmail(), FirstName: request.GetFirstName(), LastName: request.GetLastName(),
			Salt: helpers.ConvertByteToStringBase64(salt), Password: helpers.ConvertByteToStringBase64(hashPassword), IsAuthenticated: false}
		//Save the user to the database
		saveNewUser(connectedDB, helpers.GetAccountInfoCollection(), &newUser)

		//Generate authentication code and expired time
		authenCode, expiredTime := helpers.GenerateAuthenCode()
		//Create a struct for user's authentication
		newAuthenCode := userAuthentication{Email: request.GetEmail(), Code: authenCode, ExpiredTime: expiredTime}
		//Save the user authentication code to the database
		saveAuthenticationCode(connectedDB, helpers.GetAccountAuthenticationCollection(), &newAuthenCode)
		//Send the user authentication code to the user's email
		sendAuthenticationCode(request.GetEmail(), request.GetFirstName(), authenCode)

		//Send the user a welcome message and user authentication number to the provided email address
		log.Println("A new account was created for:", request.GetEmail())
		return &pb.CreateAccountReply{Status: &pb.Status{Code: int32(helpers.GetCreatedStatusCode()), Message: "account created"}}
	} else {
		//Either email, firstName, lastName, or password is invalid
		return &pb.CreateAccountReply{Status: &pb.Status{Code: int32(helpers.GetInvalidInformation()),
			Message: "information provided are invalid"}}
	}

}

//SaveNewUser : Save the user to database
func saveNewUser(database *mongo.Database, collectionName string, user *UserInfo) {
	//Connect to the users collection in the database
	accountInfoCollection := helpers.ConnectToCollection(database, collectionName)
	//Save the user into the database
	_, err := accountInfoCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Save user to DB error: ", err)
	} else {
		log.Println("Save", user.Email, "to the account information collection")
	}
	return
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
		log.Fatalln("Send email error: ", err)
	}
	log.Println("Email", userEmail, "authentication code")

}
