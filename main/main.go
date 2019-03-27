package main

import (
	pb "../proto"
	"Naturae_Server/helpers"
	"Naturae_Server/users"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

type server struct {

}

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()
	//users.Login("visalhok123@gmail.com", "ABab1234!@#")
	//Connect to all of the services that is needed to run the server
	email := "visalhok123@gmail.com"
	firstName := "Visal"
	lastName := "Hok"
	password := "ABab1234!@#"
	users.CreateAccount(email, firstName, lastName, password)
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

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest, opts ...grpc.CallOption) (*pb.HelloReply, error){
	host, err := os.Hostname()
	if err != nil{
		log.Fatal(err)
	}
	return &pb.HelloReply{Message: "Hello" + host}, nil
}

//Initialize and start the server
func createServer() {
	listener, err := net.Listen("tcp", "8080")
	if err != nil{
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServerRequestsServer(s, &server{})
	if err:= s.Serve(listener): err != nil{
		log.Fatalf("failed to server: %v", err)
	}
}
