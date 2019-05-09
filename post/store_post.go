package post

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"bytes"
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/net/context"
	"log"
	"time"
)

type ImageDescription struct {
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

//SavePost : save the user post to the database
func SavePost(request *pb.CreatePostRequest, ownerEmail string) *pb.CreatePostReply {
	connectedDB := helpers.ConnectToDB(helpers.GetPostDatabase())
	//Create a unique id for the post
	postID := helpers.CreateUUID()
	//Save the image
	if savePostToS3(postID, &request.EncodedImage) {
		//Create an image url to be stored
		imageURL := "https://s3-us-west-2.amazonaws.com/naturae-post-photos/public-posts/" + postID
		//Crate a structure to store the post
		newPost := ImageDescription{PostID: postID, OwnerEmail: ownerEmail, Title: request.Title, Species: request.Species, Description: request.Description,
			Latitude: request.GetLat(), Longitude: request.GetLng(), EncodedImage: imageURL, Date: time.Now()}
		postCollection := connectedDB.Collection(helpers.GetStorePostsCollection())
		_, err := postCollection.InsertOne(context.Background(), newPost)
		//Check if the post was able to save
		if err != nil {
			log.Printf("error while saving post: %v", err)
		}
		//Post was able to save
		return &pb.CreatePostReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "post saved"}}
	}
	return &pb.CreatePostReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "server time out"}}
}

// Save the image to AWS S3 server
func savePostToS3(postID string, image *string) bool {
	saveImage, _ := base64.StdEncoding.DecodeString(*image)
	params := &s3.PutObjectInput{
		Bucket:        aws.String("naturae-post-photos"),
		Key:           aws.String("public-posts/" + postID),
		Body:          bytes.NewReader(saveImage),
		ContentLength: aws.Int64(int64(len(saveImage))),
		ContentType:   aws.String("image/jpeg"),
	}
	_, err := helpers.GetS3Session().PutObject(params)
	if err != nil {
		log.Printf("Saving image to S3 bucket error: %v", err)
		return false
	}
	log.Printf("Saving %s to S3 successfully to public post", postID)
	return true
}
