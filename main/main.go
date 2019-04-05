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
	"math/rand"
	"net"
	"strconv"
	"time"
)

type server struct{}

var result = make(chan string, 10)

func main() {
	//Close the connection to the database when the server is turn off
	defer cleanUpServer()
	for i := 0; i < 10; i++ {
		go func(i int) {
			test(strconv.Itoa(i))
		}(i)
	}

	//users.Login("visalhok123@gmail.com", "ABab1234!@#")
	//Connect to all of the services that is needed to run the server
	//email := "visalhok123@gmail.com"
	//firstName := "Visal"
	//lastName := "Hok"
	//password := "ABab1234!@#"
	//users.CreateAccount(email, password, firstName, lastName)

}

func test(num string) string {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	fmt.Println(num)
	return num
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
func (s *server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (response *pb.CreateAccountReply, err error) {
	result, status := users.CreateAccount(request)
	return &pb.CreateAccountReply{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken,
		Status: status}, nil

}
