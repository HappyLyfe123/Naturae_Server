package helpers

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

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

//GetAccountInfo: Get user account info using email
func (email EmailAccess) GetAccountInfo() *UserLoginInfo {
	var authenInfo UserLoginInfo
	filter := bson.D{{Key: "email", Value: email}}
	connectConnection := ConnectToDB(GetUserDatabase()).Collection(GetAccountInfoCollection())
	err := connectConnection.FindOne(context.TODO(), filter).Decode(authenInfo)
	if err != nil {

	}
	return &authenInfo
}

//GetAccountInfo : Get user account info using token
func (tokenID TokenAccess) GetAccountInfo() {
	filter := bson.D{{Key: "token_id", Value: tokenID}}
	connectConnection := ConnectToDB(GetUserDatabase()).Collection(GetAccessTokenCollection())
	connectConnection.FindOne(context.TODO(), filter)
}

//ConnectToDBAccount : connect to the database client
//Return : mongodb client or err
func ConnectToDBAccount() {
	var err error
	for attemptNum := 0; attemptNum < GetConnectionAttemptLimit(); attemptNum++ {
		//Connect to the mongo database server
		dbAccount, err = mongo.Connect(context.TODO(), "mongodb+srv://"+os.Getenv("DATABASE_USERNAME")+
			":"+os.Getenv("DATABASE_PASSWORD")+"@naturae-server-hxywc.mongodb.net/test?retryWrites=true")
		if err != nil {
			//Print out the error message
			log.Println("Connecting to DB account error: ", err)
		} else {
			log.Println("Connect to Naturae DB")
			break
		}
	}
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
	log.Println("Close DB connection:")
	return nil
}
