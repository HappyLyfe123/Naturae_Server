package post

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"golang.org/x/net/context"
	"log"
	"time"
)

type postDescription struct {
	PostID       string
	OwnerEmail   string
	Title        string
	Species      string
	Description  string
	Latitude     float32
	Longitude    float32
	EncodedImage string
	Date         time.Time
}

func SavePost(request *pb.CreatePostRequest, ownerEmail string) *pb.CreatePostReply {
	connectedDB := helpers.ConnectToDB(helpers.GetPostDatabase())
	//Create a unique id for the post
	postID := helpers.CreateUUID()

	newPost := postDescription{PostID: postID, OwnerEmail: ownerEmail, Title: request.Title, Species: request.Species, Description: request.Description,
		Latitude: request.GetLat(), Longitude: request.GetLng(), EncodedImage: request.GetEncodedImage(), Date: time.Now()}
	postCollection := connectedDB.Collection(helpers.GetStorePostsCollection())
	_, err := postCollection.InsertOne(context.Background(), newPost)
	if err != nil {
		log.Printf("error while saving post: %v", err)
		return &pb.CreatePostReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "server time out"}}
	}

	return &pb.CreatePostReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "post saved"}}
}
