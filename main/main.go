package main

import (
	"Naturae_Server/helpers"
	"fmt"
	"log"
	"sync"
)

func main() {
	//Connect to all of the services that is needed to run the server
	//initApp()
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

func testGoRoutine(wg *sync.WaitGroup, message string) {
	defer wg.Done()
	fmt.Println(message)
}
