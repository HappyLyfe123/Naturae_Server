package helpers

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//create a struct for storing user info in database
type UserAccount struct {
	Email            string
	First_Name       string
	Last_Name        string
	Salt             string
	Password         string
	Is_Authenticated bool
}

type GetInfo interface {
	GetAccountInfo()
}

type UserLoginInfo struct {
	Email    string
	Salt     string
	Password string
}

type EmailAccess struct {
	Email string
}

type TokenAccess struct {
	TokenID string
}

//Get the account info using email
func (email EmailAccess) GetAccountInfo() {
	filter := bson.D{{Key: "email", Value: email}}
	connectConnection := ConnectToDB(GetUserDatabase()).Collection(GetAccountInfoCollection())
	connectConnection.FindOne(context.TODO(), filter)
}

func (tokenID TokenAccess) GetAccountInfo() {

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
func ConnectToCollection(currDB *mongo.Database, collectionName string) *mongo.Collection {
	return currDB.Collection(collectionName)
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
