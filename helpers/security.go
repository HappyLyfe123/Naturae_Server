package helpers

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math/big"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AccessToken struct {
	Email       string
	FirstName   string
	LastName    string
	ID          string
	Admin       bool
	ExpiredTime time.Time
}

//Create a struct for storing token
type RefreshToken struct {
	Email       string
	ID          string
	ExpiredTime time.Time
}

// CheckAppKey : Check if the app key is valid
// @appKey : provided app key to be check
func CheckAppKey(appKey string) bool {
	//Compare the provided app key with current database app key
	if strings.Compare(GetAppKey(), appKey) == 0 {
		return true
	} else {
		return false
	}
}

//GenerateAccessToken : Generate access token
func GenerateAccessToken(email, firstName, lastName string) *AccessToken {
	//Create an access token that have a life span of 12 hours
	return &AccessToken{ID: GenerateTokenID(), Email: email, FirstName: firstName, LastName: lastName,
		Admin: false, ExpiredTime: time.Now().Add(time.Hour * 5000)}

}

//GenerateRefreshToken : Generate refresh token
func GenerateRefreshToken(email string) *RefreshToken {
	//Refresh token have a life span of 200 years
	return &RefreshToken{Email: email, ID: GenerateTokenID(), ExpiredTime: time.Now().AddDate(200, 0, 0)}
}

//GenerateTokenID : Generate an id for a token
func GenerateTokenID() string {
	return ConvertByteToStringBase64(GenerateHash(GenerateRandomBytes(GetTokenLength()), nil))
}

//IsTokenExpired : check if the token had expired already
func IsTokenExpired(expireTime time.Time) bool {
	if time.Now().Before(expireTime) || time.Now().Equal(expireTime) {
		return false
	}
	return true
}

//GenerateAuthenCode : generate an authentication to verify the user
func GenerateAuthenCode() (string, time.Time) {
	//The time that the authentication code is created
	expiredTime := time.Now().Add(time.Minute * 30)
	authenCode := GenerateRandomNumber(GetAuthCodeMinNum(), GetAuthCodeMaxNum()).String()

	return authenCode, expiredTime
}

//GenerateHash : Hash password with salt using sha512
func GenerateHash(originalData, salt []byte) []byte {
	//Add the salt to the end of the user password
	//After adding salt to user password convert the password with
	//salt into byte array
	//Hash the byte array password using sha515 then convert the hash to
	//base64 encoding and return it
	hashPassword := sha512.New().Sum(append(originalData, salt...))
	return hashPassword
}

//GenerateRandomBytes : random data
func GenerateRandomBytes(len int) []byte {
	//Generate an array of the specific byte in length
	newData := make([]byte, len)
	//If there an error it going to loop until there are no error
	for {
		//Generate a pseudo random number
		_, err := rand.Read(newData)
		if err == nil {
			return newData
		} else {
			log.Fatalf("Generate random bytes error: %v", err)
		}
	}

}

//GenerateRandomNumber : generate a random between 100000 and 999999
func GenerateRandomNumber(minNum, maxNum int64) *big.Int {
	//Generate a number between 0 and the given number
	randomNum, err := rand.Int(rand.Reader, big.NewInt(maxNum))
	if err != nil {
		log.Fatalf("Random number generator error: %v", err)
	}
	//Add minNum to the generated number
	return randomNum.Add(randomNum, big.NewInt(minNum))

}

func GetAccessToken(currDatabase *mongo.Database, ID string) (*AccessToken, error) {
	var result AccessToken
	filter := bson.D{{Key: "id", Value: ID}}
	tokenCollection := ConnectToCollection(currDatabase, GetAccessTokenCollection())
	err := tokenCollection.FindOne(context.Background(), filter).Decode(&result)
	//There no token id match token id
	if err != nil {
		return &result, err
	}
	//There already exist a token ID in the database
	return &result, nil
}

func GetRefreshToken(currDatabase *mongo.Database, email string) (*RefreshToken, error) {
	var result RefreshToken
	filter := bson.D{{Key: "email", Value: email}}
	tokenCollection := ConnectToCollection(currDatabase, GetRefreshTokenCollection())
	err := tokenCollection.FindOne(context.Background(), filter).Decode(&result)
	//There no token id match token id
	if err != nil {
		return &result, err
	}
	//There already exist a token ID in the database
	return &result, nil
}
