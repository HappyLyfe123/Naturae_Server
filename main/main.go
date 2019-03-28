package main

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/proto"
	"Naturae_Server/users"
	"errors"
	"google.golang.org/grpc/reflection"

	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()
	//users.Login("visalhok123@gmail.com", "ABab1234!@#")
	//Connect to all of the services that is needed to run the server
}

//Initialize all of the variable to be uses
func init() {
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
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("unable to listen on 8080 port: %v", err)
	}
	log.Println("listening on port 8080")

	srv := grpc.NewServer()
	pb.RegisterServerRequestsServer(srv, &server{})
	reflection.Register(srv)
	err = srv.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println(in.Name)
	return &pb.HelloReply{
		Message: "Hello Sam",
	}, nil
}

//Create user account
func (s *server) CreateAccount(ctx context.Context, in *pb.CreateAccountRequest) (response *pb.CreateAccountReply, err error) {
	result, errorList := users.CreateAccount(in.Email, in.Password, in.FirstName, in.LastName)
	if errorList != nil {
		return &pb.CreateAccountReply{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken, ErrorList: errorList},
			errors.New(helpers.ErrorFormat(helpers.GetInvalidInformation(), "invalid information", "the information"+
				" provided contains invalid format"))
	}
	return &pb.CreateAccountReply{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken,
		ErrorList: errorList}, nil
}
