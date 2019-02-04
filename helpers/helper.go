package helpers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/smtp"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//GenerateRandomBytes : random data
func GenerateRandomBytes(len int16) (string, error) {
	newData := make([]byte, len)
	_, err := rand.Read(newData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(newData), nil
}

//GetUsername : Get the user useranme
func GetUsername(email string) {
	//collection := ConnectToDB().Database("Naturae-Server").Collection("Users")

}

//GetEmail : Get the user email address
func GetEmail(username string) {

}

//ConnectToDBClient : connect to the database client
func ConnectToDBClient(username, password *string)(*mongo.Client, error){
	client, err := mongo.Connect(context.TODO(), "mongodb+srv://"+*username+":"+*password+"@naturae-server-hxywc.mongodb.net/test?retryWrites=true")
	if err != nil{
		return nil, err
	}
	return client, nil
}

//ConnectToDB : Connect to the database
func ConnectToDB(databaseName *string, client *mongo.Client) *mongo.Database{

	database := client.Database(*databaseName)
	return database
}

//ConnectToCollection : connect to a collection in the database
func ConnectToCollection(database *mongo.Database, collectionName *string) *mongo.Collection {
	return database.Collection(*collectionName)
}

//CloseDatabase : close the current collection to the database
func CloseDatabase(database *mongo.Client) error {
	err := database.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

//ConvertStringToByte : convert a string to bytes
func ConvertStringToByte(stringData string) []byte {
	return []byte(stringData)
}

//SendEmail : send email to the specific user
func SendEmail(recieverEmail, body, subject *string) {
	from := "naturae.outdoor@gmail.com"
	pass := "B#sGwqrEp*C17xbmcDMChQYwHQ#wgL"
	to := *recieverEmail
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + *subject + "\n\n" +
		*body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

}
