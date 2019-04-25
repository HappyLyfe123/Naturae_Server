package post

import (
	"Naturae_Server/helpers"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func GetPost() {
	connectedDB := helpers.ConnectToDB(helpers.GetPostDatabase())
	postCollection := connectedDB.Collection(helpers.GetStorePostsCollection())
	findOptions := options.Find()
	var results []*ImageDescription
	cur, err := postCollection.Find(context.TODO(), nil, findOptions)
	if err != nil {
		log.Printf("Getting database result error %v: ", err)
	}
	if cur != nil {
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
	} else {
		log.Println("Collection is empty")
	}

	//fmt.Print(results)
}
