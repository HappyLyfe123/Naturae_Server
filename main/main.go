package main

import (
	"Naturae_Server/helpers"
	"Naturae_Server/users"
	"fmt"
	"log"
)

func main() {
	//Connect to all of the services that is needed to run the server
	initApp()
	email := "visalhok123@gmail.com"
	firstName := "Visal"
	lastName := "Hok"
	password := "ABab1234!@#"
	newAccount := users.CreateAccount(&email, &firstName, &lastName, &password)
	fmt.Println(newAccount.Error)
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
