package main

import (
	"Naturae_Server/helpers"
	"fmt"
	"log"
)

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()
	//users.Login("visalhok123@gmail.com", "ABab1234!@#")
	//Connect to all of the services that is needed to run the server
	fmt.Println(helpers.GetCurrentDBConnection().ListDatabaseNames(nil, nil, nil))
	//email := "visalhok123@gmail.com"
	//firstName := "Visal"
	//lastName := "Hok"
	//password := "ABab1234!@#"
	//users.CreateAccount(email, firstName, lastName, password)
}

//Initialize all of the variable to be uses
func init() {
	//Initialize global variable in the helper package
	helpers.ConnectToGmailAccount()
	helpers.ConnectToDBAccount()
	//Create listener for server
	//createServer()

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

}
