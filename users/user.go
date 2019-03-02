package users

import (
	"Naturae_Server/helpers"
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//create a struct for storing user info in database
type userAccount struct {
	Email           string
	FirstName       string
	LastName        string
	Salt            string
	Password        string
	IsAuthenticated bool
}

//create authentication structure for user
type userAuthentication struct {
	Email     string
	Code      string
	StartTime time.Time
}

//Create a struct for storing token
type accessToken struct {
	Email              string
	TokenID            string
	ExpiredTime        time.Time
}

type RefreshToken struct {
	Email       string
	TokenID     string
	ExpiredTime time.Time
}

//SaveToken : save token to the database
func SaveToken(wg *sync.WaitGroup, database *mongo.Database, token interface{}) error {
	//defer wg.Done()
	var collectionName string
	//Check to determine what type of toke to save
	switch token.(type) {
	case accessToken:
		collectionName = helpers.GetAccessTokenCollection()
		break
	case RefreshToken:
		collectionName = helpers.GetRefreshTokenCollection()
		break
	default:
		log.Println("Invalid token type")
		return errors.New("Invalid token type")
	}

	//Connect to the database collection
	currCollection := helpers.ConnectToCollection(database, collectionName)
	//Save token to the database
	_, err := currCollection.InsertOne(nil, token)
	if err != nil {
		log.Println("Save token error: ", err)
		return err
	}
	return nil
}


