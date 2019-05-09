package main

import (
	proto "Naturae_Server/naturaeproto"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

var client proto.ServerRequestsClient

func main() {

	//cred := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	//conn, err := grpc.Dial("naturae.host:443", grpc.WithTransportCredentials(cred))
	//conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(cred))
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}
	client = proto.NewServerRequestsClient(conn)

	/******** GET ROOM NAME RPC TEST ******/
	/*
		replyRoomName, err := client.GetRoomName(context.Background(), &proto.RoomRequest{
			UserOwner1: "buffalo@jam.com",
			UserOwner2: "nanaeaubry@gmail.com",
		})

		if err != nil {
			log.Fatalf("cannot call rpc replyRoomName: %v \n", err)
		} else {
			fmt.Println("say getRoom sent, with status")
			fmt.Println(replyRoomName.GetStatus())

			fmt.Printf("%v\n", replyRoomName.GetRoomName())
		}
	*/
	/******** HELLO WORLD RPC TEST ********/
	/*
			client.SayHello(context.Background(), &proto.HelloRequest{
				Name: "Sam",
			})
			fmt.Println("say hello sent")

		fmt.Println("")
	*/
	/******** Search Friends RPC TEST ********/

	replySearchFriends, err := client.SearchUsers(context.Background(), &proto.UserSearchRequest{
		User:  "nanaeaubry@gmail.com",
		Query: "nothing",
	})

	if err != nil {
		log.Fatalf("cannot call rpc SearchFriends: %v \n", err)
	} else {
		fmt.Println("say Search Friends sent, with status")
		fmt.Println(replySearchFriends.GetStatus())

		fmt.Printf("%v\n", replySearchFriends.GetUsers())
		fmt.Printf("%v\n", replySearchFriends.GetAvatars())
	}

	//fmt.Println("")

	/******** Search USERS RPC TEST ********/

	replySearchUsers, err := client.SearchUsers(context.Background(), &proto.UserSearchRequest{
		User:  "",
		Query: "limstevenlbw@gmail.com",
	})

	if err != nil {
		log.Fatalf("cannot call rpc SearchUsers: %v \n", err)
	} else {
		fmt.Println("say Search Users sent, with status")
		fmt.Println(replySearchUsers.GetStatus())

		fmt.Printf("%v\n", replySearchUsers.GetUsers())
		fmt.Printf("%v\n", replySearchUsers.GetAvatars())
	}

	fmt.Println("")

	/******** Add USERS RPC TEST ********/
	/*
		replyAddFriend, err := client.AddFriend(context.Background(), &proto.FriendRequest{
			Sender:   "limstevenlbw@gmail.com",
			Receiver: "nanae@savage.com",
		})
		if err != nil {
			log.Fatalf("cannot call rpc RemoveFriend: %v", err)
		} else {
			fmt.Println("say Add Friend sent, with status")
			fmt.Println(replyAddFriend.GetStatus())
		}
	*/
	/******** Remove USERS RPC TEST ********/
	/*
		fmt.Println("")
		replyRemoveFriend, err := client.RemoveFriend(context.Background(), &proto.FriendRequest{
			Sender:   "limstevenlbw@gmail.com",
			Receiver: "nanae@savage.com",
		})
		if err != nil {
			log.Fatalf("cannot call rpc AddFriend: %v \n", err)
		} else {
			fmt.Println("say Remove Friend sent, with status")
			fmt.Println(replyRemoveFriend.GetStatus())
		}


			createAccount, _ := client.Login(context.Background(), &proto.LoginRequest{
					AppKey: "fsdfdsfdsfsdfs",
					Email:  "visalhok@yahoo.com",
					Password : "ABCDabcd1!",
			})
			fmt.Println(createAccount.Status.GetMessage())
	*/
}
