package helpers

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

//Current db account connected to
var dbAccount *mongo.Client

//ConnectToDBAccount : connect to the database client
//Return : err if there an error
func ConnectToDBAccount() {
	var err error
	//Connect to the mongo database server
	dbAccount, err = mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s"+
		"@naturae-server-hxywc.gcp.mongodb.net/test", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"))))
	dbAccount.Connect(nil)

	if err != nil {
		//Print out the error message
		log.Fatalf("Connecting to Naturae database account error: %v", err)
	}
	log.Println("Connected to Naturae database account")
}

//ConnectToDB : connect to a database
func ConnectToDB(databaseName string) *mongo.Database {
	return dbAccount.Database(databaseName)
}

//ConnectToCollection : connect to the database collection
func ConnectToCollection(currDB *mongo.Database, collectionName string) *mongo.Collection {
	return currDB.Collection(collectionName)
}

//GetCurrentDBConnection : Get the database that the server is currently connected to
func GetCurrentDBConnection() *mongo.Client {
	return dbAccount
}

//DropCollection : drop the collection that is currently connect it to
func DropCollection(currCollection *mongo.Collection) {
	err := currCollection.Drop(context.Background())
	if err != nil {
		log.Println("Dropping collection failed error: ", err)
	}
}

//DropDatabase : drop the database that is currently connected to
func DropDatabase(currDB *mongo.Database) error {
	err := currDB.Drop(context.Background())
	if err != nil {
		log.Println("Dropping database failed error: ", err)
	}
	return nil
}

//CloseConnectionToDatabaseAccount : close the current collection to the database
func CloseConnectionToDatabaseAccount() error {
	//Disconnect from the database account
	err := dbAccount.Disconnect(nil)
	if err != nil {
		return err
	}
	log.Println("Closed database connection to Naturae account")
	return nil
}
