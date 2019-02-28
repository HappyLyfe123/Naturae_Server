package main

import (
	"Naturae_Server/helpers"
	"Naturae_Server/users"
	"fmt"
	"log"
	"reflect"
	"time"
)

type testStruct struct {
}

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()
	//Connect to all of the services that is needed to run the server
	initApp()
	connectedDB := helpers.ConnectToDB(helpers.GetUserDatabase())
	i := users.RefreshToken{"a", "a", time.Now()}
	err := users.SaveToken(nil, connectedDB, i)
	if err != nil {
		fmt.Println(reflect.TypeOf(err))
	}
	/*	email := "visalhok123@gmail.com"
		firstName := "Visal"
		lastName := "Hok"
		password := "ABab1234!@#"
		users.CreateAccount(email, firstName, lastName, password)*/
	//currDatabase := helpers.ConnectToDB("Users")

}

//Initialize all of the variable to be uses
func initApp() {
	//Initialize global variable in the helper package
	helpers.ConnectToGmailAccount()
	helpers.ConnectToDBAccount()
	//Create listener for server
	createServer()
}

//Close all of the connection to everything that the server is connected to
func cleanUpServer() {
	err := helpers.CloseConnectionToDatabaseAccount()
	if err != nil {
		log.Println("Closed connection to DB error: ", err)
	}
}

//Initialize and start the server
func createServer() {
	//listener, err := net.Listen("tcp", ":8080")
	//if err != nil {
	//	log.Fatalf("Unable to listen prot 8080: %v", err)
	//}

}
