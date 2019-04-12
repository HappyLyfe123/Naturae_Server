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
	if err != nil {
		//Print out the error message
		log.Fatalf("Creating connection to Naturae database account error: %v", err)
		return
	}
	//Check if the database is connected
	err = dbAccount.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Conneting to Naturae database error: %v", err)
	} else {
		log.Println("Connected to Naturae database account")
	}

}

//ConnectToDB : connect to a the specific database
func ConnectToDB(databaseName string) *mongo.Database {
	return dbAccount.Database(databaseName)
}

//ConnectToCollection : connect to the collection in the database
func ConnectToCollection(currDB *mongo.Database, collectionName string) *mongo.Collection {
	return currDB.Collection(collectionName)
}

//DropCollection : drop the collection that is currently connect it to
func DropCollection(currCollection *mongo.Collection) {
	err := currCollection.Drop(context.Background())
	if err != nil {
		log.Println("Dropping collection failed error: ", err)
	}
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
