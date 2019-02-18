package main

import (
	"Naturae_Server/helpers"
	"Naturae_Server/users"
	"log"
)

func main() {
	//Connect to all of the services that is needed to run the server
	initApp()
	email := "visalhok123@gmail.com"
	name := "Visal"
	users.SendAuthenticationEmail(&email, &name)
}

//Initialize all of the variable to be uses
func initApp() {
	//Initialize global variable in the helper package
	helpers.ConnectToGmailAccount()
	err := helpers.ConnectToDBAccount()
	if err != nil {
		log.Print("Unable to connect need")
	}
}
