package main

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"Naturae_Server/users"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct{}

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()
	//users.Login("visalhok123@gmail.com", "ABab1234!@#")
	//Connect to all of the services that is needed to run the server
	//result := users.CreateAccount(&pb.CreateAccountRequest{FirstName: "Visal", LastName: "Hok", Email: "visalhok123@gmail.com", Password : "ABab1234!@#"})
	//if result.Status != nil{
	//	fmt.Println("Error")
	//}
	test := make(chan bool, 2)
	test <- true
	test <- true
	if <-test || <-test {
		fmt.Println("Work")
	}

}

//Initialize all of the variable to be uses
func init() {
	//Initialize global variable in the helper package
	helpers.ConnectToGmailAccount()
	helpers.ConnectToDBAccount()
	//asyncq.StartTaskDispatcher(10)
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
		Message: "Hello " + in.Name,
	}, nil
}

//Create user account
func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountReply, error) {
	result := users.CreateAccount(request)
	return result, nil
}

//Login user
func (s *server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginReply, error) {
	return users.Login(request), nil
}

//Account authentication
func (s *server) AccountAuthentication(ctx context.Context, request *pb.AccountAuthenRequest) (*pb.AccountAuthenReply, error) {
	return users.AuthenticateAccount(request), nil
}

func (s *server) GetNewAccessToken(context.Context, *pb.GetAccessTokenRequest) (*pb.GetAccessTokenReply, error) {
	panic("implement me")
}
