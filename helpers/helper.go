package helpers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/smtp"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//Declare local global variable
var dbClient *mongo.Client
var emailClient smtp.Auth
var err error

//Email : Emailing structure for sending email
type Email struct {
	Receiver string
	Subject  string
	Body     string
}

//GenerateRandomBytes : random data
func GenerateRandomBytes(len int16) (string, error) {
	//Generate an array of the specific byte in length
	newData := make([]byte, len)
	//Generate a pseudorandom number
	_, err := rand.Read(newData)
	if err != nil {
		return "", err
	}
	//Convert the new byte of data into base64 string
	return base64.StdEncoding.EncodeToString(newData), nil
}

//ConnectToDBClient : connect to the database client
//Return : mongodb client or err
func ConnectToDBClient() (*mongo.Client, error) {
	//Connect to the mongo database server
	dbClient, err = mongo.Connect(context.TODO(), "mongodb+srv://"+os.Getenv("DATABASE_USERNAME")+
		":"+os.Getenv("DATABASE_PASSWORD")+"@naturae-server-hxywc.mongodb.net/test?retryWrites=true")
	if err != nil {
		//Return error back to the calling methods
		return nil, err
	}
	//Return mongodb client object back to calling method
	return dbClient, nil
}

//ConnectToGmailClient : Connect to the Gmail email client
func ConnectToGmailClient() {
	//Initialize emailClient
	emailClient = smtp.PlainAuth("", GetGmailEmailAdddress(), os.Getenv("GMAIL_PASSWORD"), GetGmailSMTPHost())
}

//ConnectToCollection : connect to a collection in the database
//databaseName : the database name to connect to
//collectionName : the collection name to connect to
//Return : mongodb collection object
func ConnectToCollection(databaseName, collectionName string) *mongo.Collection {
	return dbClient.Database(databaseName).Collection(collectionName)
}

//CloseDatabase : close the current collection to the database
//@dtabase : the databe client that is to be close
func CloseDatabase(database *mongo.Client) error {
	err := database.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

//ConvertStringToByte : convert a string to bytes array
func ConvertStringToByte(stringData string) []byte {
	return []byte(stringData)
}

//SendEmail : send email to the specific user
func SendEmail(sendEmail *Email) error {

	//Format the message
	msg := "From: " + GetAppName() + "\n" +
		"To: " + sendEmail.Receiver + "\n" +
		"Subject: " + sendEmail.Subject + "\n\n" +
		sendEmail.Body

	//Send an email using smtp
	err := smtp.SendMail(GetGmailSMTPServer(), emailClient, GetAppName(), []string{sendEmail.Receiver}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}
