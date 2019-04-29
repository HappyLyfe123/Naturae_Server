package post

import (
	"Naturae_Server/helpers"
	pb "Naturae_Server/naturaeproto"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
)

func RetrievePosts(radius, latitude, longitude float64) *pb.GetPostReply {
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
		postLatitude := helpers.ConvertDegreeToRadian(float64(elem.GetLat()))
		postLongitude := helpers.ConvertDegreeToRadian(float64(elem.GetLng()))
		fmt.Printf("%f, %f\t", elem.GetLat(), elem.GetLng())
		//postLatitude := elem.GetLat()
		//Check if the longitude and latitude of the post is within the radius
		if math.Acos(math.Sin(latitude)*math.Sin(postLatitude)+math.Cos(latitude)*math.Cos(postLatitude)*
			math.Cos(longitude-postLongitude))*6371 < radius {
			//If the longitude and latitude is within the radius then add the post to the result list
			results = append(results, elem)
			fmt.Println(elem.Title)
		}
		fmt.Println(math.Acos(math.Sin(latitude)*math.Sin(postLatitude)+math.Cos(latitude)*math.Cos(postLatitude)*
			math.Cos(longitude-postLongitude)) * 6371)
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
