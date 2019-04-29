package post

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
)

func RetrievePosts(radius int32, latitude, longitude float64) *pb.GetPostReply {
	connectedDB := helpers.ConnectToDB(helpers.GetPostDatabase())
	postCollection := connectedDB.Collection(helpers.GetStorePostsCollection())
	var results []*pb.PostStruct
	cur, err := postCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Getting database result error: %v", err)
		return internalServerError()
	}
	for cur.Next(context.TODO()) {
		var elem *pb.PostStruct
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("Error while decoding get post result: %v", err)
			return internalServerError()

		}
		//Convert the post latitude and longitude from degree to radian
		postLatitude := helpers.ConvertDegreeToRadian(float64(elem.Lat))
		postLongitude := helpers.ConvertDegreeToRadian(float64(elem.Lng))
		//Check if the longitude and latitude of the post is within the radius
		if math.Acos(math.Sin(latitude)*math.Sin(float64(postLatitude))+math.Cos(latitude)*math.Cos(float64(postLatitude))*
			math.Cos(float64(postLongitude)-(longitude)))*6371 < float64(radius) {
			//If the longitude and latitude is within the radius then add the post to the result list
			results = append(results, elem)
		}

	}
	err = cur.Close(context.TODO())
	if err != nil {
		log.Println(err)
		return internalServerError()
	}

	return &pb.GetPostReply{Status: &pb.Status{Code: helpers.GetOkStatusCode(), Message: "success"}, Reply: results}
}

func internalServerError() *pb.GetPostReply {
	return &pb.GetPostReply{Status: &pb.Status{Code: helpers.GetInternalServerErrorStatusCode(), Message: "server error"},
		Reply: nil}
}
