package helpers

import (
	"context"
	"encoding/base64"
	"fmt"
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

//Email : Emailing structure for sending email
type Email struct {
	Receiver string
	Subject  string
	Body     string
}

//ConnectToDBAccount : connect to the database client
//Return : mongodb client or err
func ConnectToDBAccount() (err error) {
	//Connect to the mongo database server
	dbAccount, err = mongo.Connect(context.TODO(), "mongodb+srv://"+os.Getenv("DATABASE_USERNAME")+
		":"+os.Getenv("DATABASE_PASSWORD")+"@naturae-server-hxywc.mongodb.net/test?retryWrites=true")
	if err != nil {
		//Return error back to the calling methods
		return err
	}
	//Return mongodb client object back to calling method
	return nil
}

//ConnectToDB : connect to a database
//databaseName : the database name to connect to
//collectionName : the collection name to connect to
//Return : mongodb collection object
func ConnectToDB(databaseName string) *mongo.Database {
	return dbAccount.Database(databaseName)

}

//ConnectToCollection : connect to the database collection
func ConnectToCollection(currDB *mongo.Database, collectionName *string) *mongo.Collection {
	return currDB.Collection(*collectionName)
}

//CloseConnectionToDatabase : close the current collection to the database
//@database : the database client that is to be close
func CloseConnectionToDatabase(database *mongo.Client) error {
	//Disconnect from the database account
	err := database.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

//ConnectToGmailAccount : Connect to the Gmail email client
func ConnectToGmailAccount() {
	//Initialize emailClient
	gmailAccount = smtp.PlainAuth("", GetGmailEmailAddress(), os.Getenv("GMAIL_PASSWORD"), GetGmailSMTPHost())
}

//ConvertStringToByte : convert a string to bytes array
func ConvertStringToByte(stringData *string) *[]byte {
	//Convert string into byte array
	bytesData := []byte(*stringData)
	return &bytesData
}

//ConvertByteToString : Convert bytes array to string
func ConvertByteToString(bytesData *[]byte) *string {
	//Convert bytes array into base64 string
	stringData := base64.StdEncoding.EncodeToString(*bytesData)
	return &stringData
}

//FindUser : find the user information the in database
func FindUser(email *string, database *mongo.Database, collectionName string) *mongo.SingleResult {
	findUuserFilter := bson.D{{"Email", *email}}
	userCollection := ConnectToCollection(database, &collectionName)
	//Check if the email exist in the database
	//Return true if the email doesn't exist in the database
	user := userCollection.FindOne(context.TODO(), findUuserFilter)
	return user
}

//SendEmail : send email to the specific user
func SendEmail(sendEmail *Email) error {
	//Format the send message
	msg := fmt.Sprintf("From: %s <%s>\nTo: %s\nSubject: %s\n\n%s", GetAppName(), GetGmailEmailAddress(),
		sendEmail.Receiver, sendEmail.Subject, sendEmail.Body)
	//Send an email using SMTP
	err := smtp.SendMail(GetGmailSMTPServer(), gmailAccount, GetGmailEmailAddress(), []string{sendEmail.Receiver}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

//IsEmailValid : Check if the email match the guideline and if there an existing account with that email address
func IsEmailValid(email *string, database *mongo.Database, collectionName string) bool {
	//Generate an Regexp to check user email address
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	//Check if the email format is correct
	if emailRegexp.MatchString(*email) {
		//Generate a filter to for checking the database
		filter := bson.D{{"Email", *email}}
		//Connect to the collection in the database
		userCollection := ConnectToCollection(database, &collectionName)
		//Check if the email exist in the database
		//Return true if the email doesn't exist in the database
		if userCollection.FindOne(context.TODO(), filter) == nil {
			return true

		}
		//Return false if the email exist in the database
		return false
	}

	return false
}

//IsPasswordValid : check if the password matches the gideline
func IsPasswordValid(password *string) bool {
	//Initialize all the local variables
	var number, upper, lower, special bool
	character := 0
	//Check if the password contain all of the necessary characters
	for _, c := range *password {
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
func IsNameValid(name *string) bool {
	for _, c := range *name {
		//If the name doesn't match the guide it will return false
		if !regexp.MustCompile(`[a-zA-Z0-9._ '-]`).MatchString(string(c)) {
			return false
		}

	}
	//Return that the name is valid
	return true
}
