package main

import (
	"Naturae_Server/helpers"
	"Naturae_Server/users"
	"log"
)

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()

	//Connect to all of the services that is needed to run the server
	initApp()
	email := "visalhok123@gmail.com"
	firstName := "Visal"
	lastName := "Hok"
	password := "ABab1234!@#"
	users.CreateAccount(email, firstName, lastName, password)
	//currDatabase := helpers.ConnectToDB("Users")

}

//Initialize all of the variable to be uses
func initApp() {
	//Initialize global variable in the helper package
	helpers.ConnectToGmailAccount()
	err := helpers.ConnectToDBAccount()
	if err != nil {
		log.Print("Unable to connect")
	}
	//Create listener for server
	//createServer()
}

//Close all of the connection to everything that the server is connected to
func cleanUpServer() {
	err := helpers.CloseConnectionToDatabase(helpers.GetCurrentDBConnection())
	if err != nil {
		log.Println("Close connection to DB error: ", err)
	}
}

// //Initialize and start the server
// func createServer() {
// 	listener, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		log.Fatalf("unable to listen prot 8080: %v", err)
// 	}

// 	srv := grpc.NewServer()
// 	proto.RegisterServerRequestsServer(srv, &Server{})
// 	srv.Serve(listener)
// }
