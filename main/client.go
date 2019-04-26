package main

import (
	proto "Naturae_Server/naturaeproto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

var client proto.ServerRequestsClient

func main(){

	//cred := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	//conn, err := grpc.Dial("naturae.host:443", grpc.WithTransportCredentials(cred))
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	defer conn.Close()

	if err != nil{
		log.Fatalf("cannot dial server: %v", err)
	}

	client = proto.NewServerRequestsClient(conn)
	client.SayHello(context.Background(), &proto.HelloRequest{
		Name: "Sam",
	})
	
	client.SearchUsers(context.Background(), &proto.UserSearchRequest{
		User: "limstevenlbw@gmail.com",
		Query: "nothing",
	})

	createAccount, _ := client.Login(context.Background(), &proto.LoginRequest{
			AppKey: "fsdfdsfdsfsdfs",
			Email:  "visalhok@yahoo.com",
			Password : "ABCDabcd1!",
	})
	fmt.Println(createAccount.Status.GetMessage())
}
