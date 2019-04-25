package post

import (
	"Naturae_Server/helpers"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type test struct {
	Name string
}

func GetPost() {
	connectedDB := helpers.ConnectToDB(helpers.GetPostDatabase())
	postCollection := connectedDB.Collection(helpers.GetStorePostsCollection())
	var results []*ImageDescription
	filter := bson.D{{"latitude", bson.D{{"$gt", 0}}}}
	cur, err := postCollection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Getting database result error: %v", err)
		return
	}
	for cur.Next(context.TODO()) {
		var elem ImageDescription
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("Error while decoding get post result: %v", err)
		}
		results = append(results, &elem)

	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	cur.Close(context.TODO())
	fmt.Print(results[0].OwnerEmail)
}
