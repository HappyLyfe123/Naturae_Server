package helpers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"regexp"
	"unicode"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

//Declare local global variable
var dbAccount *mongo.Client
var gmailAccount smtp.Auth

type helperError struct {
	ErrorCode    int16
	ErrorMessage string
}

//Email : Emailing structure for sending email
type Email struct {
	Receiver string
	Subject  string
	Body     string
}

//ConnectToGmailAccount : Connect to the Gmail email client
func ConnectToGmailAccount() {
	//Initialize emailClient
	gmailAccount = smtp.PlainAuth(GetAppEmailAddress(), GetAppEmailAddress(), os.Getenv("GMAIL_PASSWORD"), GetGmailSMTPHost())
}

//ConvertStringToByte : convert a string to bytes array
func ConvertStringToByte(stringData string) []byte {
	//Convert string into byte array
	bytesData := []byte(stringData)
	return bytesData
}

//ConvertByteToString : Convert bytes array to string
func ConvertByteToString(bytesData []byte) string {
	//Convert bytes array into base64 string
	stringData := base64.StdEncoding.EncodeToString(bytesData)
	return stringData
}

//RandomizeData : Scramble the given data
func RandomizeData(originalData []byte) []byte {
	data := originalData
	for {
		_, err := rand.Read(data)
		if err != nil {
			log.Println("Randomize data error: ", err)
		} else {
			return data
		}
	}
}

//SendEmail : send email to the specific user
func SendEmail(sendEmail *Email) error {
	//Format the send message
	msg := fmt.Sprintf("From: %s <%s>\nTo: %s\nSubject: %s\n\n%s", GetAppName(), GetAppEmailAddress(),
		sendEmail.Receiver, sendEmail.Subject, sendEmail.Body)
	//Send an email using SMTP
	err := smtp.SendMail(GetGmailSMTPServer(), gmailAccount, GetAppEmailAddress(), []string{sendEmail.Receiver}, []byte(msg))
	if err != nil {
		log.Println("Unable to connect to SMTP error: ", err)
		return err
	}

	return nil
}

//IsEmailValid : Check if the email match the guideline and if there an existing account with that email address
func IsEmailValid(email string) (bool, helperError) {
	//Generate an Regexp to check user email address
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	//Check if the email format is correct
	if !emailRegexp.MatchString(email) {
		log.Println("Invalid email format")
		//Return false if the email exist in the database
		return false, helperError{ErrorCode: GetInvalidEmailCode(), ErrorMessage: "Invalid email address"}
	}
	return true, helperError{}
}

//EmailExist : check in the database if there already an account with that email address
func EmailExist(email string, database *mongo.Database, collectionName string) (bool, helperError) {
	//Create a struct to get the result
	var result struct {
		Email string
	}

	//Set a filter for the database to search through
	filter := bson.D{{Key: "email", Value: email}}
	//Connect to the collection database
	userCollection := ConnectToCollection(database, collectionName)
	//Make a request to the database
	err := userCollection.FindOne(context.TODO(), filter).Decode(&result)

	//If there an error or no email found it will return false
	if err != nil {
		return false, helperError{}
	}
	return true, helperError{ErrorCode: GetEmailExistCode(), ErrorMessage: "Email exist"}

}

//IsPasswordValid : check if the password matches the gideline
func IsPasswordValid(password string) (bool, helperError) {
	//Initialize all the local variables
	var number, upper, lower, special bool
	character := 0
	//Check if the password contain all of the necessary characters
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
			character++
		case unicode.IsUpper(c):
			upper = true
			character++
		case unicode.IsLower(c):
			lower = true
			character++
		case regexp.MustCompile("[!@#$%&^]").MatchString(string(c)):
			special = true
			character++
		default:
			return false, helperError{ErrorCode: GetInvalidPasswordCode(), ErrorMessage: "Invalid password"}
		}
	}
	//Check if the password is at least 8 characters long
	eightOrMore := character >= 8
	//The password meet all of the guideline
	if number && upper && lower && special && eightOrMore {
		return true, helperError{}
	}
	return false, helperError{ErrorCode: GetInvalidPasswordCode(), ErrorMessage: "Invalid passwords"}
}

//IsNameValid : Check if the name contain any of not valid characters
func IsNameValid(name string) (bool, helperError) {
	for _, c := range name {
		//If the name doesn't match the guide it will return false
		if !regexp.MustCompile(`[a-zA-Z_ '-]`).MatchString(string(c)) {
			return false, helperError{ErrorCode: GetInvalidNameCode(), ErrorMessage: "Invalid name"}
		}
	}
	//Return that the name is valid
	return true, helperError{}
}

//GetCurrentDBConnection : Get the database that the server is currently connected to
func GetCurrentDBConnection() *mongo.Client {
	return dbAccount
}
