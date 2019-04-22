package helpers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/smtp"
	"os"
	"regexp"
	"time"
	"unicode"
)

//Declare local global variable
var gmailAccount smtp.Auth

//Email : Emailing structure for sending email
type Email struct {
	Receiver string
	Subject  string
	Body     string
}

type Status struct {
	Code    int32
	Message string
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

//ConvertByteToString : Convert bytes array to string as base64
func ConvertByteToStringBase64(bytesData []byte) string {
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
func IsEmailValid(email string) bool {
	//Generate an Regexp to check user email address
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	//Check if the email format is correct
	if !emailRegexp.MatchString(email) {
		log.Println("Invalid email format")
		//Return false if the email exist in the database
		return false
	}
	return true
}

//EmailExist : check in the database if there already an account with that email address
func EmailExist(email string, database *mongo.Database, collectionName string) bool {
	//Create a struct to get the result
	var result struct {
		email string
	}

	//Set a filter for the database to search through
	filter := bson.D{{Key: "email", Value: email}}
	//Connect to the collection database
	userCollection := ConnectToCollection(database, collectionName)
	//Make a request to the database
	err := userCollection.FindOne(context.TODO(), filter).Decode(&result)

	//If there an error or no email found it will return false
	if err != nil {
		return false
	}
	return true
}

//IsPasswordValid : check if the password matches the guideline
func IsPasswordValid(password string) bool {
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
		case regexp.MustCompile("[!@#$%^&+]").MatchString(string(c)):
			special = true
			character++
		default:
			return false
		}
	}
	//Check if the password is at least 8 characters long
	eightOrMore := character >= 8
	//The password meet all of the guideline
	if number && upper && lower && special && eightOrMore {
		return true
	}
	return false
}

//IsNameValid : Check if the name contain any of not valid characters
func IsNameValid(name string) bool {
	for _, c := range name {
		//If the name doesn't match the guide it will return false
		if !regexp.MustCompile(`[a-zA-Z_ '-]`).MatchString(string(c)) {
			return false
		}
	}
	//Return that the name is valid
	return true
}

//Check if the time is valid
func IsTimeValid(expiredTime time.Time) bool {
	if time.Now().After(expiredTime) {
		return false
	} else {
		return true
	}
}

func CreateUUID() string {
	newID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Creating UUID error: %v", err)
	}
	return newID.String()
}
