package post

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
)

//RetrievePosts : get all of the posts within the radius from the center of the user's map
func GetPostPreview(radius, latitude, longitude float64) *pb.GetPostPreviewReply {
	//Connect to post database
	postDB := helpers.ConnectToDB(helpers.GetPostDatabase())
	//Connect post location
	postCollection := postDB.Collection(helpers.GetStorePostsCollection())
	var results []*pb.PostStruct
	cur, err := postCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Getting database result error: %v", err)
		return &pb.GetPostPreviewReply{Status: internalServerError(), Reply: nil}
	}
	for cur.Next(context.TODO()) {
		var elem *pb.PostStruct
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("Error while decoding get post result: %v", err)
			return &pb.GetPostPreviewReply{Status: internalServerError(), Reply: nil}

		}
		//Convert the post latitude and longitude from degree to radian
		postLatitude := helpers.ConvertDegreeToRadian(float64(elem.GetLatitude()))
		postLongitude := helpers.ConvertDegreeToRadian(float64(elem.GetLongitude()))
		//Check if the longitude and latitude of the post is within the radius
		if math.Acos(math.Sin(latitude)*math.Sin(postLatitude)+math.Cos(latitude)*math.Cos(postLatitude)*
			math.Cos(longitude-postLongitude))*6371 < radius {
			//If the longitude and latitude is within the radius then add the post to the result list
			results = append(results, elem)
		}
	}
	err = cur.Close(context.TODO())
	if err != nil {
		log.Println(err)
		return &pb.GetPostPreviewReply{Status: internalServerError(), Reply: nil}
	}

	return &pb.GetPostPreviewReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "success"}, Reply: results}
}

func GetPost(latitude, longitude float32) *pb.GetPostReply {
	//Connect to post database
	postDB := helpers.ConnectToDB(helpers.GetPostDatabase())
	//Connect post location
	postCollection := postDB.Collection(helpers.GetStorePostsCollection())
	var results []*pb.PostStruct
	cur, err := postCollection.Find(context.Background(), bson.D{{"latitude", latitude}, {"longitude", longitude}})
	if err != nil {
		log.Printf("Getting database result error: %v", err)
		return &pb.GetPostReply{Status: internalServerError(), Reply: nil}
	}
	for cur.Next(context.TODO()) {
		var elem *pb.PostStruct
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("Error while decoding get post result: %v", err)
			return &pb.GetPostReply{Status: internalServerError(), Reply: nil}
		}
		log.Println(elem.GetPostID())
	}
	err = cur.Close(context.TODO())
	if err != nil {
		log.Println(err)
		return &pb.GetPostReply{Status: internalServerError(), Reply: nil}
	}
	return &pb.GetPostReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "success"}, Reply: results}
}

func internalServerError() *pb.Status {
	return &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "server error"}
}
