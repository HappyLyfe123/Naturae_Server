package users

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

//Conversation the document structure of a conversation
type Conversation struct {
	RoomName string
	Users    []string
	Chatlog  []string
}

//CreateConversation builds a chatroom with a unique room ID that belongs to two verified users
//Called by AddFriends in user.go
func CreateConversation(dbConnection *mongo.Database, user1 string, user2 string) error {
	uuid := helpers.CreateUUID()
	users := []string{user1, user2}
	chatlog := []string{}

	//Build the new conversation struct
	newConvo := Conversation{RoomName: uuid, Users: users, Chatlog: chatlog}

	//Connect to the conversation collection
	convoCollection := helpers.ConnectToCollection(dbConnection, helpers.GetConversationsCollection())

	//Save the user into the database
	_, err := convoCollection.InsertOne(context.Background(), newConvo)

	if err != nil {
		log.Println("Failed to create a new conversation error: ", err)
	} else {
		log.Println("Saved a new conversation for the users: " + user1 + " and " + user2)
	}
	return err
}

//RemoveConversation deletes the chatlog belonging to the two users defined in parameters
//Called by RemoveFriends in user.go
func RemoveConversation(dbConnection *mongo.Database, user1 string, user2 string) error {
	//Remove the document that contains the array with specified users
	removeFilter := bson.D{
		{"users", bson.D{{"$all", bson.A{user1, user2}}}},
	}
	convoCollection := helpers.ConnectToCollection(dbConnection, helpers.GetConversationsCollection())
	_, err := convoCollection.DeleteOne(context.Background(), removeFilter)

	return err
}

//GetRoomName retrieves a room name and returns it to the client by checking the two users that own it
func GetRoomName(request *pb.RoomRequest) *pb.RoomReply {
	var result Conversation
	dbConnection := helpers.ConnectToDB(helpers.GetUserDatabase())
	convoCollection := helpers.ConnectToCollection(dbConnection, helpers.GetConversationsCollection())
	//Acquire the document that contains the array with specified users
	roomFilter := bson.D{
		{"users", bson.D{{"$all", bson.A{request.GetUserOwner1(), request.GetUserOwner2()}}}},
	}

	//Acquire single mongo result and decode into struct
	err := convoCollection.FindOne(context.Background(), roomFilter).Decode(&result)

	if err != nil {
		return &pb.RoomReply{RoomName: " ",
			Status: &pb.Status{Code: helpers.GetNotFoundStatusCode(), Message: "Failure"}}
	}

	room := (&result).RoomName

	return &pb.RoomReply{RoomName: room,
		Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "Success"}}

}
