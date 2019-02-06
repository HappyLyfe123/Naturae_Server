package main

import "Naturae_Server/helpers"

func main() {
	initApp()

}

func initApp() {
	//Initialize global variable in the helper package
	go helpers.ConnectToDBClient()
	go helpers.ConnectToGmailClient()
}
